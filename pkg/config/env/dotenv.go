package env

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type EnvConfigProvider struct {
}

func (EnvConfigProvider) Get(name string) (value string, err error) {
	v := os.Getenv(name)
	return v, nil
}

func (p *EnvConfigProvider) Bind(target interface{}) error {
	if err := env.Parse(target); err != nil {
		return err
	}
	return nil
}

func (p *EnvConfigProvider) Engine() interface{} {
	return p
}

func NewDotEnvConfig() (*EnvConfigProvider, error) {
	cfg := &EnvConfigProvider{}
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
