package operator

type Config struct {
	Logger   LoggerConfig   `json:"logger"`
	Ethereum EthereumConfig `json:"ethereum"`
}

type EthereumConfig struct {
	Host string `json:"host"`
}

type LoggerConfig struct {
	LogLevel string `json:"logLevel"`
	Format   string `json:"format"`
}
