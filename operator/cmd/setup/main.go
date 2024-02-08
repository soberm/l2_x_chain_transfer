package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/sirupsen/logrus"
	"operator/pkg/operator"
	"os"
	"path/filepath"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.TextFormatter)

	log.Out = os.Stdout

	log.Level = logrus.InfoLevel
}

func main() {
	b := flag.String("b", "./build/", "filename of the config file")
	flag.Parse()

	if _, err := os.Stat(*b); errors.Is(err, os.ErrNotExist) {
		dirPath := filepath.Dir(*b)

		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatalf("mkdir: %w", err)
		}
	}

	err := CompileSetupGenerateBurn(*b)
	if err != nil {
		log.Fatalf("compile setup generate burn circuit: %v", err)
	}

	err = CompileSetupGenerateClaim(*b)
	if err != nil {
		log.Fatalf("compile setup generate claim circuit: %v", err)
	}

}

func CompileSetupGenerateBurn(dst string) error {

	var burnCircuit operator.BurnCircuit
	burnCircuit.AllocateSlicesMerkleProofs()

	constraintSystem, err := CompileCircuit(dst+"burn_circuit.r1cs", &burnCircuit)
	if err != nil {
		log.Fatalf("compile burn circuit: %v", err)
	}

	_, vk, err := Setup(dst+"burn_proving_key", dst+"burn_verifying_key", constraintSystem)
	if err != nil {
		log.Fatalf("setup circuit: %w", err)
	}

	err = GenerateSolidityVerifier(dst+"burn_verifier.sol", vk)
	if err != nil {
		log.Fatalf("generate solidity verifier: %v", err)
	}

	return nil
}

func CompileSetupGenerateClaim(dst string) error {

	var claimCircuit operator.ClaimCircuit
	claimCircuit.AllocateSlicesMerkleProofs()

	claimSystem, err := CompileCircuit(dst+"claim_circuit.r1cs", &claimCircuit)
	if err != nil {
		log.Fatalf("compile claim circuit: %v", err)
	}

	_, claimVk, err := Setup(dst+"claim_proving_key", dst+"claim_verifying_key", claimSystem)
	if err != nil {
		log.Fatalf("setup circuit: %w", err)
	}

	err = GenerateSolidityVerifier(dst+"claim_verifier.sol", claimVk)
	if err != nil {
		log.Fatalf("generate solidity verifier: %v", err)
	}

	return nil
}

func CompileCircuit(dst string, circuit frontend.Circuit) (constraint.ConstraintSystem, error) {

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, fmt.Errorf("compile circuit: %w", err)
	}
	log.Info("circuit compiled")

	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	_, err = ccs.WriteTo(file)
	if err != nil {
		return nil, fmt.Errorf("write to file: %w", err)
	}
	log.Info("circuit saved")

	return ccs, nil
}

func Setup(pkDst string, vkDst string, system constraint.ConstraintSystem) (groth16.ProvingKey, groth16.VerifyingKey, error) {

	provingKey, verifyingKey, err := groth16.Setup(system)
	if err != nil {
		return nil, nil, fmt.Errorf("setup circuit: %w", err)
	}
	log.Info("circuit setup completed")

	pkFile, err := os.OpenFile(pkDst, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer pkFile.Close()

	vkFile, err := os.OpenFile(vkDst, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer pkFile.Close()

	_, err = provingKey.WriteTo(pkFile)
	if err != nil {
		return nil, nil, fmt.Errorf("write to file: %w", err)
	}
	log.Info("proving key saved")

	_, err = verifyingKey.WriteTo(vkFile)
	if err != nil {
		return nil, nil, fmt.Errorf("write to file: %w", err)
	}
	log.Info("verifying key saved")

	return provingKey, verifyingKey, nil
}

func GenerateSolidityVerifier(dst string, vk groth16.VerifyingKey) error {

	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if err := vk.ExportSolidity(file); err != nil {
		return fmt.Errorf("export solidity verifier: %w", err)
	}

	log.Info("solidity verifier saved")

	return nil
}
