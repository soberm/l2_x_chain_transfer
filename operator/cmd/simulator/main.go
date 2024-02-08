package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	configFile := flag.String("c", "./configs/sim.json", "config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)

	viper.AutomaticEnv()

	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading config file, %v", err)
	}

	log.Infof("using config: %s\n", viper.ConfigFileUsed())

	var config simulator.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	switch config.Logger.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text", "":
		log.SetFormatter(&logrus.TextFormatter{})
	default:
		log.Errorf("invalid log format: %s, using default 'text'", config.Logger.Format)
		log.SetFormatter(&logrus.TextFormatter{})
	}

	level, err := logrus.ParseLevel(config.Logger.LogLevel)
	if err != nil {
		log.Errorf("invalid log level: %s, using default 'info'", config.Logger.LogLevel)
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(level)
	}

	simulator.SetLogger(log)

	sim := simulator.NewSimulator(&config)
	err = sim.Run()
	if err != nil {
		log.Fatalf("run simulator: %v", err)
	}
}
