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
	config      *Config
	burnSystem  constraint.ConstraintSystem
	claimSystem constraint.ConstraintSystem
}

func NewSimulator(config *Config) *Simulator {
	return &Simulator{
		config: config,
	}
}

func (s *Simulator) Run() error {
	log.Info("starting simulator")

	ethClient, err := ethclient.DialContext(context.Background(), s.config.Ethereum.Host)
	if err != nil {
		return fmt.Errorf("dial eth: %w", err)
	}

	burnVerifierContract, err := operator.NewBurnVerifierContract(common.HexToAddress(s.config.Ethereum.BurnVerifierContract), ethClient)
	if err != nil {
		return fmt.Errorf("create verifier contract: %w", err)
	}

	rollupContract, err := operator.NewRollupContract(common.HexToAddress(s.config.Ethereum.RollupContract), ethClient)
	if err != nil {
		return fmt.Errorf("create rollup contract: %w", err)
	}

	file, err := os.OpenFile(s.config.BurnCircuit.ConstraintSystemPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	ccs := groth16.NewCS(ecc.BN254)
	_, err = ccs.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
	}

	var m1, m2 runtime.MemStats

	file, err = os.OpenFile(s.config.BurnCircuit.ProvingKeyPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	runtime.GC()
	runtime.ReadMemStats(&m1)
	pkBurn := groth16.NewProvingKey(ecc.BN254)
	_, err = pkBurn.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
	}
	runtime.ReadMemStats(&m2)
	pkMemoryBurn := m2.TotalAlloc - m1.TotalAlloc

	file, err = os.OpenFile(s.config.BurnCircuit.VerifyingKeyPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	vkBurn := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vkBurn.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
	}

	file, err = os.OpenFile(s.config.ClaimCircuit.ConstraintSystemPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	ccsClaim := groth16.NewCS(ecc.BN254)
	_, err = ccsClaim.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
	}

	file, err = os.OpenFile(s.config.ClaimCircuit.ProvingKeyPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	runtime.GC()
	runtime.ReadMemStats(&m1)
	pkClaim := groth16.NewProvingKey(ecc.BN254)
	_, err = pkClaim.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
	}
	runtime.ReadMemStats(&m2)

	pkMemoryClaim := m2.TotalAlloc - m1.TotalAlloc

	file, err = os.OpenFile(s.config.ClaimCircuit.VerifyingKeyPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	vkClaim := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vkClaim.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
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

	rollup, err := NewRollup()
	if err != nil {
		return fmt.Errorf("create rollup: %w", err)
	}

	for i := 0; i < s.config.Runs; i++ {
		measurement := make([]string, 0)
		measurement = append(measurement, strconv.Itoa(i))
		measurement = append(measurement, strconv.Itoa(operator.BatchSize))

		transfers, err := rollup.GenerateTransfers(operator.BatchSize)
		if err != nil {
			return fmt.Errorf("generate transactions: %w", err)
		}

		witness, err := rollup.Burn(transfers)
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

		runtime.GC()
		runtime.ReadMemStats(&m1)

		start := time.Now()
		proof, err := groth16.Prove(ccs, pkBurn, witness)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime := time.Since(start)
		runtime.ReadMemStats(&m2)

		measurement = append(measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		measurement = append(measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryBurn))))
		log.Infof("Burn Proving Time: %v", provingTime)
		log.Infof("Burn Memory Usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryBurn))

		err = groth16.Verify(proof, vkBurn, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		ethereumProof, err := operator.ProofToEthereumProof(proof)
		if err != nil {
			return fmt.Errorf("convert proof to ethereum proof: %w", err)
		}

		compressedProof, err := burnVerifierContract.CompressProof(nil, ethereumProof)

		ecdsaPrivateKey, err := crypto.HexToECDSA("40a22e3e69ce6e6ebd2267567699b3ea90d1553cda128c2b43af69ac83d9c0ed")
		if err != nil {
			return fmt.Errorf("ecdsa private key: %w", err)
		}

		chainID, err := ethClient.ChainID(context.Background())
		if err != nil {
			return fmt.Errorf("chain id: %w", err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
		if err != nil {
			return fmt.Errorf("new transactor: %w", err)
		}
		auth.GasPrice = big.NewInt(20000000000)

		tx, err := rollupContract.Burn(auth, postStateRoot, transactionsRoot, compressedProof)
		if err != nil {
			return fmt.Errorf("submit proof: %w", err)
		}

		receipt, err := bind.WaitMined(context.Background(), ethClient, tx)
		if err != nil || receipt.Status != 1 {
			return err
		}

		witness, err = rollup.Claim(transfers)
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
		proof, err = groth16.Prove(ccsClaim, pkClaim, witness, proverOption)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime = time.Since(start)
		runtime.ReadMemStats(&m2)

		log.Infof("Claim Proving Time: %v", provingTime)
		log.Infof("Claim Memory Usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryClaim))
		measurement = append(measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		measurement = append(measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryClaim))))

		err = groth16.Verify(proof, vkClaim, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		ethereumProof, err = operator.ProofToEthereumProof(proof)
		if err != nil {
			return fmt.Errorf("convert proof to ethereum proof: %w", err)
		}
		log.Infof("Ethereum Proof: %+v", ethereumProof)

		compressedProof, err = burnVerifierContract.CompressProof(nil, ethereumProof)
		log.Infof("Compressed Ethereum Proof: %+v", compressedProof)

		tx, err = rollupContract.Claim(auth, postStateRoot, transactionsRoot, compressedProof)
		if err != nil {
			return fmt.Errorf("submit proof: %w", err)
		}

		receipt, err = bind.WaitMined(context.Background(), ethClient, tx)
		if err != nil || receipt.Status != 1 {
			return err
		}

		data = append(data, measurement)
	}

	file, err = os.Create(s.config.Dst)
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
