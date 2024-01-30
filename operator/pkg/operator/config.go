package operator

type Config struct {
	Logger   LoggerConfig   `json:"logger"`
	Ethereum EthereumConfig `json:"ethereum"`
}

type EthereumConfig struct {
	Host       string `json:"host"`
	PrivateKey string `json:"privateKey"`
}

type LoggerConfig struct {
	LogLevel string `json:"logLevel"`
	Format   string `json:"format"`
}
