package operator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *Config

	ethClient  *ethclient.Client
	privateKey *ecdsa.PrivateKey
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
