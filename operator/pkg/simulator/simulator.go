package simulator

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint/solver"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"operator/pkg/operator"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Simulator struct {
	runs int
	dst  string
}

func NewSimulator(runs int, dst string) *Simulator {
	return &Simulator{
		runs: runs,
		dst:  dst,
	}
}

func (s *Simulator) Run() error {
	log.Info("starting simulator")

	ethClient, err := ethclient.DialContext(context.Background(), "ws://127.0.0.1:8545")
	if err != nil {
		return fmt.Errorf("dial eth: %w", err)
	}

	verifierContract, err := operator.NewBurnVerifierContract(common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"), ethClient)
	if err != nil {
		return fmt.Errorf("create verifier contract: %w", err)
	}

	var burnCircuit operator.BurnCircuit
	burnCircuit.AllocateSlicesMerkleProofs()

	var claimCircuit operator.ClaimCircuit
	claimCircuit.AllocateSlicesMerkleProofs()

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &burnCircuit)
	if err != nil {
		return fmt.Errorf("compile burn circuit: %w", err)
	}
	log.Info("burn circuit compiled")

	ccsClaim, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &claimCircuit)
	if err != nil {
		return fmt.Errorf("compile claim circuit: %w", err)
	}
	log.Info("claim circuit compiled")

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

	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf("setup burn circuit: %w", err)
	}
	log.Info("burn circuit setup completed")
	runtime.ReadMemStats(&m2)

	pkMemoryBurn := m2.TotalAlloc - m1.TotalAlloc

	solidityBurnVerifierFile, err := os.OpenFile("./BurnVerifier.sol", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer solidityBurnVerifierFile.Close()
	err = vk.ExportSolidity(solidityBurnVerifierFile)
	if err != nil {
		return fmt.Errorf("export solidity verifier: %w", err)
	}

	runtime.GC()
	runtime.ReadMemStats(&m1)
	pkClaim, vkClaim, err := groth16.Setup(ccsClaim)
	if err != nil {
		return fmt.Errorf("setup claim circuit: %w", err)
	}
	log.Info("claim circuit setup completed")
	runtime.ReadMemStats(&m2)

	pkMemoryClaim := m2.TotalAlloc - m1.TotalAlloc

	solidityClaimVerifierFile, err := os.OpenFile("./ClaimVerifier.sol", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer solidityClaimVerifierFile.Close()
	err = vkClaim.ExportSolidity(solidityClaimVerifierFile)
	if err != nil {
		return fmt.Errorf("export solidity verifier: %w", err)
	}

	rollup, err := NewRollup()
	if err != nil {
		return fmt.Errorf("create rollup: %w", err)
	}

	for i := 0; i < s.runs; i++ {
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

		runtime.GC()
		runtime.ReadMemStats(&m1)

		start := time.Now()
		proof, err := groth16.Prove(ccs, pk, witness)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}
		provingTime := time.Since(start)
		runtime.ReadMemStats(&m2)

		measurement = append(measurement, strconv.Itoa(int(provingTime.Milliseconds())))
		measurement = append(measurement, strconv.Itoa(int(bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryBurn))))
		log.Infof("Burn Proving Time: %v", provingTime)
		log.Infof("Burn Memory Usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemoryBurn))

		err = groth16.Verify(proof, vk, publicWitness)
		if err != nil {
			return fmt.Errorf("failed to verify proof: %v", err)
		}

		ethereumProof, err := operator.ProofToEthereumProof(proof)
		if err != nil {
			return fmt.Errorf("convert proof to ethereum proof: %w", err)
		}
		log.Infof("Ethereum Proof: %+v", ethereumProof)

		compressedProof, err := verifierContract.CompressProof(&bind.CallOpts{
			Pending: true,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, [8]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)})
		if err != nil {
			return fmt.Errorf("compress proof: %w", err)
		}

		log.Infof("Compressed Ethereum Proof: %+v", compressedProof)

		witness, err = rollup.Claim(transfers)
		if err != nil {
			return fmt.Errorf("update state: %w", err)
		}

		publicWitness, _ = witness.Public()

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

		data = append(data, measurement)
	}

	file, err := os.Create(s.dst)
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
