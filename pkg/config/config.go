package config

const EnvDevelopment = "development"

type Config struct {
	Env                     string `env:"ENV" envDefault:"development"`
	Listen                  string `env:"LISTEN" envDefault:"127.0.0.1:2333"`
	DataPath                string `env:"DATA_PATH" envDefault:"data.db"`
	SessionName             string `env:"SESSION_NAME" envDefault:"session_id"`
	SessionCryptoKey        string `env:"SESSION_CRYPTO_KEY" envDefault:"12345678901234567890123456789012"`
	MonitorInterval         int    `env:"MONITOR_INTERVAL" envDefault:"5"`
	MonitorRollInterval     int    `env:"MONITOR_ROLL_INTERVAL" envDefault:"60"`
	MonitorRealtimeInterval int    `env:"MONITOR_REALTIME_INTERVAL" envDefault:"360"`
}
