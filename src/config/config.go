package config

import (
	"context"

	"dxkite.cn/meownest/pkg/config"
)

type Config struct {
	DataPath         string `env:"DATA_PATH"`
	SessionName      string `env:"SESSION_NAME" envDefault:"session_id"`
	SessionCryptoKey string `env:"SESSION_CRYPTO_KEY" envDefault:"12345678901234567890123456789012"`
}

func Get(ctx context.Context) *Config {
	cfg := Config{}
	config.Bind(ctx, &cfg)
	return &cfg
}
