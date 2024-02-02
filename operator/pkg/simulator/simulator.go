package simulator

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"operator/pkg/operator"
	"runtime"
	"time"
)

type Simulator struct {
}

func NewSimulator() *Simulator {
	return &Simulator{}
}

func (s *Simulator) Run() error {
	log.Info("starting simulator")

	var burnCircuit operator.BurnCircuit
	burnCircuit.AllocateSlicesMerkleProofs()

	var claimCircuit operator.ClaimCircuit
	claimCircuit.AllocateSlicesMerkleProofs()

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &burnCircuit)
	if err != nil {
		return fmt.Errorf("compile burnCircuit: %w", err)
	}
	log.Info("burn circuit compiled")

	/*	claimCS, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &claimCircuit)
		if err != nil {
			return fmt.Errorf("compile burnCircuit: %w", err)
		}
		log.Info("claim circuit compiled")*/

	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf("setup burn circuit: %w", err)
	}
	log.Info("burn circuit setup completed")
	runtime.ReadMemStats(&m2)

	pkMemory := m2.TotalAlloc - m1.TotalAlloc

	/*	_, _, err = groth16.Setup(claimCS)
		if err != nil {
			return fmt.Errorf("setup claim circuit: %w", err)
		}
		log.Info("claim circuit setup completed")*/

	rollup, err := NewRollup()
	if err != nil {
		return fmt.Errorf("create rollup: %w", err)
	}

	transfers, err := rollup.GenerateTransfers(operator.BatchSize)
	if err != nil {
		return fmt.Errorf("generate transactions: %w", err)
	}

	witness, err := rollup.UpdateState(transfers)
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

	log.Infof("Proving Time: %v", provingTime)
	log.Infof("Memory Usage: %v MB", bToMb(m2.TotalAlloc-m1.TotalAlloc+pkMemory))

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return fmt.Errorf("failed to verify proof: %v", err)
	}

	log.Info("stop simulator")
	return nil
}
