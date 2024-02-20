package operator

type Config struct {
	Runs         int            `json:"runs"`
	Dst          string         `json:"dst"`
	Submit       bool           `json:"submit"`
	Logger       LoggerConfig   `json:"logger"`
	Ethereum     EthereumConfig `json:"ethereum"`
	BurnCircuit  CircuitConfig  `json:"burnCircuit"`
	ClaimCircuit CircuitConfig  `json:"claimCircuit"`
}

type EthereumConfig struct {
	Host                  string `json:"host"`
	PrivateKey            string `json:"privateKey"`
	RollupContract        string `json:"rollupContract"`
	BurnVerifierContract  string `json:"burnVerifierContract"`
	ClaimVerifierContract string `json:"claimVerifierContract"`
}

type LoggerConfig struct {
	LogLevel string `json:"logLevel"`
	Format   string `json:"format"`
}

type CircuitConfig struct {
	ConstraintSystemPath string `json:"constraintSystemPath"`
	ProvingKeyPath       string `json:"provingKeyPath"`
	VerifyingKeyPath     string `json:"verifyingKeyPath"`
}
