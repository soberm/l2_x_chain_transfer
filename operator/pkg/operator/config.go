package operator

type Config struct {
	Logger   LoggerConfig   `json:"logger"`
	Ethereum EthereumConfig `json:"ethereum"`
	Circuit  CircuitConfig  `json:"circuit"`
}

type EthereumConfig struct {
	Host       string `json:"host"`
	PrivateKey string `json:"privateKey"`
}

type LoggerConfig struct {
	LogLevel string `json:"logLevel"`
	Format   string `json:"format"`
}

type CircuitConfig struct {
	ConstraintSystemPath string `json:"constraintSystemPath"`
	ProvingKeyPath       string `json:"provingKeyPath"`
	VerifyingKeyPath     string `json:"verifyingKeyPath"`
	SolidityVerifierPath string `json:"solidityVerifierPath"`
}
