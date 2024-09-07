package config

import "fmt"

type ConfigOption func(p ConfigProvider)

type ConfigCreator func(opts ...ConfigOption) (ConfigProvider, error)

var configProvider = map[string]ConfigCreator{}

func Register(name string, creator ConfigCreator) {
	configProvider[name] = creator
}

func NewProvider(name string, opts ...ConfigOption) (ConfigProvider, error) {
	if creator, ok := configProvider[name]; !ok {
		return nil, fmt.Errorf("unknown config provider %s", name)
	} else {
		return creator(opts...)
	}
}
