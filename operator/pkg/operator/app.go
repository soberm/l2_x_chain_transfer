package operator

import (
	"context"
	"time"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Run() {
	log.Infof("running app...")
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
}

func (a *App) Stop() {
	log.Infof("stopping app...")
	a.cancel()
}
