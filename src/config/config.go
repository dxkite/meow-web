package config

import (
	"context"

	"dxkite.cn/meownest/pkg/config"
)

type Config struct {
	DataPath                string `env:"DATA_PATH"`
	SessionName             string `env:"SESSION_NAME" envDefault:"session_id"`
	SessionCryptoKey        string `env:"SESSION_CRYPTO_KEY" envDefault:"12345678901234567890123456789012"`
	MonitorInterval         int    `env:"MONITOR_INTERVAL" envDefault:"5s"`
	MonitorRollInterval     int    `env:"MONITOR_ROLL_INTERVAL" envDefault:"60s"`
	MonitorRealtimeInterval int    `env:"MONITOR_REALTIME_INTERVAL" envDefault:"360s"`
}

func Get(ctx context.Context) *Config {
	cfg := Config{}
	config.Bind(ctx, &cfg)
	return &cfg
}
