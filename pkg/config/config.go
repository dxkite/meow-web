package config

import (
	"errors"
)

type configProviderKey string

var ConfigProviderKey configProviderKey = "pkg/config_provider"

var ErrMissProvider = errors.New("missing config provider")

type ConfigProvider interface {
	Get(name string) (value string, err error)
	Bind(target interface{}) error
	Engine() interface{}
}

// 绑定到结构
func Bind(name string, target any, opts ...ConfigOption) error {
	provider, err := NewProvider(name)
	if err != nil {
		return err
	}
	return provider.Bind(target)
}
