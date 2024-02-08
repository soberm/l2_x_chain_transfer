package simulator

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/constraint/solver"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"operator/pkg/operator"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Simulator struct {
	config *Config

	chainID               *big.Int
	ethClient             *ethclient.Client
	burnVerifierContract  *operator.BurnVerifierContract
	claimVerifierContract *operator.ClaimVerifierContract
	rollupContract        *operator.RollupContract

	rollup *Rollup

	burnSystem  constraint.ConstraintSystem
	claimSystem constraint.ConstraintSystem

	burnProvingKey  groth16.ProvingKey
	claimProvingKey groth16.ProvingKey

	burnVerifyingKey  groth16.VerifyingKey
	claimVerifyingKey groth16.VerifyingKey

	burnTotalAlloc  uint64
	claimTotalAlloc uint64

	measurement []string
}

func NewSimulator(config *Config) *Simulator {
	return &Simulator{
		config: config,
	}
}

func (s *Simulator) ConnectEthereum() error {
	var err error
	s.ethClient, err = ethclient.DialContext(context.Background(), s.config.Ethereum.Host)
	if err != nil {
		return fmt.Errorf("dial eth: %w", err)
	}

	s.chainID, err = s.ethClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("chain id: %w", err)
	}

	s.burnVerifierContract, err = operator.NewBurnVerifierContract(common.HexToAddress(s.config.Ethereum.BurnVerifierContract), s.ethClient)
	if err != nil {
		return fmt.Errorf("create verifier contract: %w", err)
	}

	s.rollupContract, err = operator.NewRollupContract(common.HexToAddress(s.config.Ethereum.RollupContract), s.ethClient)
	if err != nil {
		return fmt.Errorf("create rollup contract: %w", err)
	}

	return nil
}

func (s *Simulator) LoadConstraintSystems() error {
	var err error
	s.burnSystem, s.burnProvingKey, s.burnVerifyingKey, err = s.LoadConstraintSystem(
		s.config.BurnCircuit.ConstraintSystemPath,
		s.config.BurnCircuit.ProvingKeyPath,
		s.config.BurnCircuit.VerifyingKeyPath,
	)
	if err != nil {
		return fmt.Errorf("load burn constraint system: %w", err)
	}

	s.claimSystem, s.claimProvingKey, s.claimVerifyingKey, err = s.LoadConstraintSystem(
		s.config.ClaimCircuit.ConstraintSystemPath,
		s.config.ClaimCircuit.ProvingKeyPath,
		s.config.ClaimCircuit.VerifyingKeyPath,
	)
	if err != nil {
		return fmt.Errorf("load claim constraint system: %w", err)
	}

	return nil
}

func (s *Simulator) LoadConstraintSystem(constraintSystemPath, provingKeyPath, verifyingKeyPath string) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, error) {
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

func (s *Simulator) Run() error {
	log.Info("starting simulator")

	err := s.ConnectEthereum()
	if err != nil {
		return fmt.Errorf("connect ethereum: %w", err)
	}

	err = s.LoadConstraintSystems()
	if err != nil {
		return fmt.Errorf("load constraint systems: %w", err)
	}

	s.rollup, err = NewRollup()
	if err != nil {
		return fmt.Errorf("create rollup: %w", err)
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

	for i := 0; i < s.config.Runs; i++ {
		s.measurement = make([]string, 0)
		s.measurement = append(s.measurement, strconv.Itoa(i))
		s.measurement = append(s.measurement, strconv.Itoa(operator.BatchSize))

		transfers, err := s.rollup.GenerateTransfers(operator.BatchSize)
		if err != nil {
			return fmt.Errorf("generate transactions: %w", err)
		}

		witness, err := s.rollup.Burn(transfers)
		if err != nil {
			return fmt.Errorf("update state: %w", err)
		}

		publicWitness, _ := witness.Public()

		witnessVector := publicWitness.Vector()

		preStateRoot := big.NewInt(0)
		witnessVector.(fr.Vector)[0].BigInt(preStateRoot)

		postStateRoot := big.NewInt(0)
		witnessVector.(fr.Vector)[1].BigInt(postStateRoot)

		transactionsRoot := big.NewInt(0)
		witnessVector.(fr.Vector)[2].BigInt(transactionsRoot)

		var m1, m2 runtime.MemStats

		runtime.GC()
		runtime.ReadMemStats(&m1)

		start := time.Now()
		proof, err := groth16.Prove(s.burnSystem, s.burnProvingKey, witness)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime := time.Since(start)
		runtime.ReadMemStats(&m2)

		s.measurement = append(s.measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		//s.measurement = append(s.measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryBurn))))
		log.Infof("Burn Proving Time: %v", provingTime)
		//log.Infof("Burn Memory Usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryBurn))

		err = groth16.Verify(proof, s.burnVerifyingKey, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		ethereumProof, err := operator.ProofToEthereumProof(proof)
		if err != nil {
			return fmt.Errorf("convert proof to ethereum proof: %w", err)
		}

		compressedProof, err := s.burnVerifierContract.CompressProof(nil, ethereumProof)

		ecdsaPrivateKey, err := crypto.HexToECDSA("40a22e3e69ce6e6ebd2267567699b3ea90d1553cda128c2b43af69ac83d9c0ed")
		if err != nil {
			return fmt.Errorf("ecdsa private key: %w", err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, s.chainID)
		if err != nil {
			return fmt.Errorf("new transactor: %w", err)
		}
		auth.GasPrice = big.NewInt(20000000000)

		tx, err := s.rollupContract.Burn(auth, postStateRoot, transactionsRoot, compressedProof)
		if err != nil {
			return fmt.Errorf("submit proof: %w", err)
		}

		_, err = bind.WaitMined(context.Background(), s.ethClient, tx)
		if err != nil {
			return fmt.Errorf("wait mined: %w", err)
		}

		witness, err = s.rollup.Claim(transfers)
		if err != nil {
			return fmt.Errorf("update state: %w", err)
		}

		publicWitness, _ = witness.Public()

		witnessVector = publicWitness.Vector()

		witnessVector.(fr.Vector)[0].BigInt(preStateRoot)
		witnessVector.(fr.Vector)[1].BigInt(postStateRoot)
		witnessVector.(fr.Vector)[2].BigInt(transactionsRoot)

		runtime.GC()
		runtime.ReadMemStats(&m1)

		start = time.Now()
		proverOption := backend.WithSolverOptions(solver.WithHints(operator.Div))
		proof, err = groth16.Prove(s.claimSystem, s.claimProvingKey, witness, proverOption)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime = time.Since(start)
		runtime.ReadMemStats(&m2)

		log.Infof("Claim Proving Time: %v", provingTime)
		//log.Infof("Claim Memory Usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryClaim))
		s.measurement = append(s.measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		//s.measurement = append(s.measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryClaim))))

		err = groth16.Verify(proof, s.claimVerifyingKey, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		ethereumProof, err = operator.ProofToEthereumProof(proof)
		if err != nil {
			return fmt.Errorf("convert proof to ethereum proof: %w", err)
		}

		compressedProof, err = s.burnVerifierContract.CompressProof(nil, ethereumProof)

		tx, err = s.rollupContract.Claim(auth, postStateRoot, transactionsRoot, compressedProof)
		if err != nil {
			return fmt.Errorf("submit proof: %w", err)
		}

		_, err = bind.WaitMined(context.Background(), s.ethClient, tx)
		if err != nil {
			return fmt.Errorf("wait mined: %w", err)
		}

		data = append(data, s.measurement)
	}

	file, err := os.Create(s.config.Dst)
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

	log.Info("stop simulator")
	return nil
}
