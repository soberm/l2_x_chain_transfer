package main

import (
	"github.com/sirupsen/logrus"
	"operator/pkg/simulator"
	"os"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.TextFormatter)

	log.Out = os.Stdout

	log.Level = logrus.InfoLevel
}

func main() {
	simulator.SetLogger(log)
	sim := simulator.NewSimulator()
	err := sim.Run()
	if err != nil {
		log.Fatalf("run simulator: %v", err)
	}
}
