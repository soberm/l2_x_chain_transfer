package operator

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *Config

	ethClient  *ethclient.Client
	privateKey *ecdsa.PrivateKey

	ccs          constraint.ConstraintSystem
	provingKey   groth16.ProvingKey
	verifyingKey groth16.VerifyingKey
}

func NewApp(config *Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
		config: config,
	}
}

func (a *App) Run() error {
	log.Infof("running app...")

	var err error
	a.ethClient, err = ethclient.DialContext(a.ctx, a.config.Ethereum.Host)
	if err != nil {
		return fmt.Errorf("dial eth: %w", err)
	}

	a.privateKey, err = crypto.HexToECDSA(a.config.Ethereum.PrivateKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}

	if err = a.CompileCircuit(); err != nil {
		return fmt.Errorf("compile circuit: %w", err)
	}

	if err = a.Setup(); err != nil {
		return fmt.Errorf("setup: %w", err)
	}

	go func() {
		for {
			select {
			case <-a.ctx.Done():
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return nil
}

func (a *App) Stop() {
	log.Infof("stopping app...")
	a.cancel()
}

func (a *App) CompileCircuit() error {

	exists := true
	if _, err := os.Stat(a.config.Circuit.ConstraintSystemPath); errors.Is(err, os.ErrNotExist) {
		dirPath := filepath.Dir(a.config.Circuit.ConstraintSystemPath)

		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}

		exists = false
	}

	file, err := os.OpenFile(a.config.Circuit.ConstraintSystemPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if exists {
		log.Info("constraint system found, loading...")
		a.ccs = groth16.NewCS(ecc.BN254)
		_, err := a.ccs.ReadFrom(file)
		if err != nil {
			return fmt.Errorf("read from file: %w", err)
		}
		log.Info("constraint system loaded")
		return nil
	}
	log.Info("constraint system not found, compiling circuit...")

	var circuit Circuit
	a.ccs, err = frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return fmt.Errorf("compile circuit: %w", err)
	}
	log.Info("circuit compiled")

	_, err = a.ccs.WriteTo(file)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}
	log.Info("circuit saved")

	return nil
}

func (a *App) Setup() error {

	exists := true
	if _, err := os.Stat(a.config.Circuit.ProvingKeyPath); errors.Is(err, os.ErrNotExist) {
		dirPath := filepath.Dir(a.config.Circuit.ProvingKeyPath)

		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}

		exists = false
	}

	provingKeyFile, err := os.OpenFile(a.config.Circuit.ProvingKeyPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer provingKeyFile.Close()

	verifyingKeyFile, err := os.OpenFile(a.config.Circuit.VerifyingKeyPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer verifyingKeyFile.Close()

	solidityVerifierFile, err := os.OpenFile(a.config.Circuit.SolidityVerifierPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer solidityVerifierFile.Close()

	if exists {
		log.Info("proving and verifying key found, loading...")

		a.provingKey = groth16.NewProvingKey(ecc.BN254)
		_, err := a.provingKey.ReadFrom(provingKeyFile)
		if err != nil {
			return fmt.Errorf("read from file: %w", err)
		}

		a.verifyingKey = groth16.NewVerifyingKey(ecc.BN254)
		_, err = a.verifyingKey.ReadFrom(verifyingKeyFile)
		if err != nil {
			return fmt.Errorf("read from file: %w", err)
		}

		log.Info("proving and verifying key loaded")
		return nil
	}
	log.Info("proving and verifying key not found, executing setup...")

	a.provingKey, a.verifyingKey, err = groth16.Setup(a.ccs)
	if err != nil {
		return fmt.Errorf("setup circuit: %w", err)
	}
	log.Info("circuit setup completed")

	_, err = a.provingKey.WriteTo(provingKeyFile)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}
	log.Info("proving key saved")

	_, err = a.verifyingKey.WriteTo(verifyingKeyFile)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}
	log.Info("verifying key saved")

	if err = a.verifyingKey.ExportSolidity(solidityVerifierFile); err != nil {
		return fmt.Errorf("export solidity verifier: %w", err)
	}
	log.Info("solidity verifier saved")

	return nil
}
