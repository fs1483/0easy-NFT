package config

import (
	"github.com/caarlos0/env/v10"
)

// Config represents shared configuration values across microservices.
type Config struct {
	HTTPPort             string `env:"HTTP_PORT" envDefault:"8080"`
	OrderServicePort     string `env:"ORDER_SERVICE_PORT" envDefault:"8081"`
	MatchingServicePort  string `env:"MATCHING_SERVICE_PORT" envDefault:"8082"`
	ExecutionServicePort string `env:"EXECUTION_SERVICE_PORT" envDefault:"8083"`
	IndexerServicePort   string `env:"INDEXER_SERVICE_PORT" envDefault:"8084"`

	PostgresDSN   string `env:"POSTGRES_DSN,notEmpty"`
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"127.0.0.1:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	MarketplaceAddr string `env:"MARKETPLACE_ADDRESS,notEmpty"`
	RPCURL          string `env:"RPC_URL,notEmpty"`
	PrivateKeyHex   string `env:"EXECUTOR_PRIVATE_KEY"`
	ChainID         uint64 `env:"CHAIN_ID,notEmpty"`
}

// Load parses environment variables into Config.
func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
