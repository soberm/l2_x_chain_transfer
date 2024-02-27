package operator

import (
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/constraint/solver"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"hash"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"time"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *Config

	chainID    *big.Int
	ethClient  *ethclient.Client
	privateKey *ecdsa.PrivateKey

	rollupContract        *RollupContract
	oracleMockContract    *OracleMockContract
	burnVerifierContract  *BurnVerifierContract
	claimVerifierContract *ClaimVerifierContract

	hFunc hash.Hash

	rollup *Rollup

	accounts    []*Account
	privateKeys []*eddsa.PrivateKey

	burnSystem  constraint.ConstraintSystem
	claimSystem constraint.ConstraintSystem

	burnProvingKey  groth16.ProvingKey
	claimProvingKey groth16.ProvingKey

	burnVerifyingKey  groth16.VerifyingKey
	claimVerifyingKey groth16.VerifyingKey

	burnTotalAlloc  uint64
	claimTotalAlloc uint64

	measurement            []string
	burnSystemMemoryUsage  uint64
	claimSystemMemoryUsage uint64
}

func NewApp(config *Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
		config: config,
		hFunc:  mimc.NewMiMC(),
	}
}

func (a *App) Run() error {
	log.Info("running app...")

	err := a.ConnectEthereum()
	if err != nil {
		return fmt.Errorf("connect ethereum: %w", err)
	}

	a.privateKeys, err = generatePrivateKeys(NumberAccounts)
	if err != nil {
		return fmt.Errorf("generate private keys: %w", err)
	}

	log.Infof("creating %v test accounts...", NumberAccounts)
	a.accounts, err = createAccounts(a.privateKeys)
	if err != nil {
		return fmt.Errorf("create accounts: %w", err)
	}

	state, err := NewState(a.hFunc, a.accounts)
	if err != nil {
		return fmt.Errorf("new state: %w", err)
	}

	a.rollup, err = NewRollup(state)
	if err != nil {
		return fmt.Errorf("create rollup: %w", err)
	}

	err = a.LoadConstraintSystems()
	if err != nil {
		return fmt.Errorf("load constraint systems: %w", err)
	}

	headerRow := []string{
		"run",
		"batchSize",
		"provingTimeBurn",
		"memoryUsageBurn",
		"provingTimeClaim",
		"memoryUsageClaim",
	}

	data := [][]string{
		headerRow,
	}

	for i := 0; i < a.config.Runs; i++ {
		a.measurement = make([]string, 0)
		a.measurement = append(a.measurement, strconv.Itoa(i))
		a.measurement = append(a.measurement, strconv.Itoa(BatchSize))

		transfers, err := a.GenerateTestTransfers(BatchSize)
		if err != nil {
			return fmt.Errorf("generate transactions: %w", err)
		}

		var rollupTransfers [BatchSize]RollupTransfer
		for j, transfer := range transfers {
			rollupTransfers[j] = transfer.RollupTransfer()
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		w, err := a.rollup.Burn(transfers)
		if err != nil {
			return fmt.Errorf("update state: %w", err)
		}

		publicWitness, _ := w.Public()

		start := time.Now()
		proof, err := groth16.Prove(a.burnSystem, a.burnProvingKey, w)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime := time.Since(start)
		runtime.ReadMemStats(&m2)

		a.measurement = append(a.measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		a.measurement = append(a.measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+a.burnSystemMemoryUsage))))
		log.Infof("burn proving time: %v", provingTime)
		log.Infof("burn memory usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+a.burnSystemMemoryUsage))

		err = groth16.Verify(proof, a.burnVerifyingKey, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		if a.config.Submit {
			err = a.SubmitBurnUpdate(proof, w, rollupTransfers)
			if err != nil {
				return fmt.Errorf("submit burn update: %w", err)
			}
			_, _, transactionsRoot := a.ExtractPublicInputs(w)
			err = a.SubmitCrossChainData(transactionsRoot, big.NewInt(NumberAccounts-1))
		}

		runtime.GC()
		runtime.ReadMemStats(&m1)
		w, err = a.rollup.Claim(transfers)
		if err != nil {
			return fmt.Errorf("update state: %w", err)
		}

		publicWitness, _ = w.Public()

		start = time.Now()
		proverOption := backend.WithSolverOptions(solver.WithHints(Div))
		proof, err = groth16.Prove(a.claimSystem, a.claimProvingKey, w, proverOption)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime = time.Since(start)
		runtime.ReadMemStats(&m2)

		log.Infof("claim proving time: %v", provingTime)
		log.Infof("claim memory usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+a.claimSystemMemoryUsage))
		a.measurement = append(a.measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		a.measurement = append(a.measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+a.claimSystemMemoryUsage))))

		err = groth16.Verify(proof, a.claimVerifyingKey, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		if a.config.Submit {
			err = a.SubmitClaimUpdate(proof, w, rollupTransfers)
			if err != nil {
				return fmt.Errorf("submit claim update: %w", err)
			}
		}

		data = append(data, a.measurement)
	}

	file, err := os.Create(a.config.Dst)
	if err != nil {
		return fmt.Errorf("create data file: %w", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	err = csvWriter.WriteAll(data)
	if err != nil {
		return fmt.Errorf("write data to file: %w", err)
	}
	csvWriter.Flush()

	log.Info("stopping app...")
	return nil
}

func (a *App) ConnectEthereum() error {
	var err error
	a.ethClient, err = ethclient.DialContext(a.ctx, a.config.Ethereum.Host)
	if err != nil {
		return fmt.Errorf("dial eth: %w", err)
	}

	a.chainID, err = a.ethClient.ChainID(a.ctx)
	if err != nil {
		return fmt.Errorf("chain id: %w", err)
	}

	a.privateKey, err = crypto.HexToECDSA(a.config.Ethereum.PrivateKey)
	if err != nil {
		return fmt.Errorf("ecdsa private key: %w", err)
	}

	a.burnVerifierContract, err = NewBurnVerifierContract(common.HexToAddress(a.config.Ethereum.BurnVerifierContract), a.ethClient)
	if err != nil {
		return fmt.Errorf("create verifier contract: %w", err)
	}

	a.oracleMockContract, err = NewOracleMockContract(common.HexToAddress(a.config.Ethereum.OracleMockContract), a.ethClient)
	if err != nil {
		return fmt.Errorf("create oracle contract: %w", err)
	}

	a.rollupContract, err = NewRollupContract(common.HexToAddress(a.config.Ethereum.RollupContract), a.ethClient)
	if err != nil {
		return fmt.Errorf("create rollup contract: %w", err)
	}

	return nil
}

func (a *App) LoadConstraintSystems() error {
	var err error

	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)
	a.burnSystem, a.burnProvingKey, a.burnVerifyingKey, err = a.LoadConstraintSystem(
		a.config.BurnCircuit.ConstraintSystemPath,
		a.config.BurnCircuit.ProvingKeyPath,
		a.config.BurnCircuit.VerifyingKeyPath,
	)
	if err != nil {
		return fmt.Errorf("load burn constraint system: %w", err)
	}
	runtime.ReadMemStats(&m2)
	a.burnSystemMemoryUsage = m2.TotalAlloc - m1.TotalAlloc

	runtime.ReadMemStats(&m1)
	a.claimSystem, a.claimProvingKey, a.claimVerifyingKey, err = a.LoadConstraintSystem(
		a.config.ClaimCircuit.ConstraintSystemPath,
		a.config.ClaimCircuit.ProvingKeyPath,
		a.config.ClaimCircuit.VerifyingKeyPath,
	)
	if err != nil {
		return fmt.Errorf("load claim constraint system: %w", err)
	}
	runtime.ReadMemStats(&m2)
	a.claimSystemMemoryUsage = m2.TotalAlloc - m1.TotalAlloc

	return nil
}

func (a *App) LoadConstraintSystem(constraintSystemPath, provingKeyPath, verifyingKeyPath string) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, error) {
	log.Infof("loading constraint system from %v...", constraintSystemPath)

	file, err := os.OpenFile(constraintSystemPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	system := groth16.NewCS(ecc.BN254)
	_, err = system.ReadFrom(file)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read from file: %w", err)
	}

	file, err = os.OpenFile(provingKeyPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(file)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read from file: %w", err)
	}

	file, err = os.OpenFile(verifyingKeyPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vk.ReadFrom(file)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read from file: %w", err)
	}

	return system, pk, vk, nil
}

func (a *App) GenerateTestTransfers(number int) ([]Transfer, error) {
	hFunc := mimc.NewMiMC()
	transfers := make([]Transfer, number)
	transferData := make([]byte, hFunc.Size()*number)

	for i := 0; i < number; i++ {
		sender, err := a.rollup.State.ReadAccount(uint64(i))
		if err != nil {
			return nil, fmt.Errorf("read account: %w", err)
		}
		transfer := NewTransfer(10,
			4,
			a.privateKeys[sender.Index.Uint64()].PublicKey,
			a.privateKeys[sender.Index.Uint64()].PublicKey,
			sender.Nonce.Uint64(),
			0,
		)

		_, msg, err := transfer.Sign(*a.privateKeys[sender.Index.Uint64()], mimc.NewMiMC())
		if err != nil {
			return nil, fmt.Errorf("failed to sign transfer: %v", err)
		}

		transfers[i] = transfer
		copy(transferData[i*hFunc.Size():(i+1)*hFunc.Size()], msg)
	}

	return transfers, nil
}

func (a *App) ExtractPublicInputs(witness witness.Witness) (*big.Int, *big.Int, *big.Int) {

	witnessVector := witness.Vector()

	preStateRoot := big.NewInt(0)
	witnessVector.(fr.Vector)[0].BigInt(preStateRoot)

	postStateRoot := big.NewInt(0)
	witnessVector.(fr.Vector)[1].BigInt(postStateRoot)

	transactionsRoot := big.NewInt(0)
	witnessVector.(fr.Vector)[2].BigInt(transactionsRoot)

	return preStateRoot, postStateRoot, transactionsRoot
}

func (a *App) SubmitBurnUpdate(proof groth16.Proof, w witness.Witness, transfers [BatchSize]RollupTransfer) error {
	log.Infof("submitting burn update...")

	ethereumProof, err := ProofToEthereumProof(proof)
	if err != nil {
		return fmt.Errorf("convert proof to ethereum proof: %w", err)
	}

	compressedProof, err := a.burnVerifierContract.CompressProof(nil, ethereumProof)
	if err != nil {
		return fmt.Errorf("compress proof: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, a.chainID)
	if err != nil {
		return fmt.Errorf("new transactor: %w", err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	_, postStateRoot, transactionsRoot := a.ExtractPublicInputs(w)

	tx, err := a.rollupContract.Burn(auth, postStateRoot, transactionsRoot, compressedProof, transfers)
	if err != nil {
		return fmt.Errorf("submit burn update: %w", err)
	}

	_, err = bind.WaitMined(a.ctx, a.ethClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	return nil
}

func (a *App) SubmitCrossChainData(transactionsRoot *big.Int, operator *big.Int) error {
	log.Infof("submitting cross chain data...")

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, a.chainID)
	if err != nil {
		return fmt.Errorf("new transactor: %w", err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	tx, err := a.oracleMockContract.SubmitTransactionsRoot(auth, transactionsRoot, operator)
	if err != nil {
		return fmt.Errorf("submit transactions root: %w", err)
	}

	_, err = bind.WaitMined(a.ctx, a.ethClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	return nil
}

func (a *App) SubmitClaimUpdate(proof groth16.Proof, w witness.Witness, transfers [BatchSize]RollupTransfer) error {
	log.Infof("submitting claim update...")

	ethereumProof, err := ProofToEthereumProof(proof)
	if err != nil {
		return fmt.Errorf("convert proof to ethereum proof: %w", err)
	}

	compressedProof, err := a.burnVerifierContract.CompressProof(nil, ethereumProof)
	if err != nil {
		return fmt.Errorf("compress proof: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(a.privateKey, a.chainID)
	if err != nil {
		return fmt.Errorf("new transactor: %w", err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	_, postStateRoot, transactionsRoot := a.ExtractPublicInputs(w)

	tx, err := a.rollupContract.Claim(auth, postStateRoot, transactionsRoot, big.NewInt(NumberAccounts-1), compressedProof, transfers)
	if err != nil {
		return fmt.Errorf("submit claim update: %w", err)
	}

	_, err = bind.WaitMined(a.ctx, a.ethClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	return nil
}
