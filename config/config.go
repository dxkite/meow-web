package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Route struct {
	Pattern string `yaml:"pattern"`
	// 需要登录才能访问
	Sign bool `yaml:"sign"`
	// 是登录API
	SignIn bool `yaml:"signin"`
	// 是登出API
	SignOut bool `yaml:"signout"`
	// 可用后端
	Backend []string `yaml:"backend"`
}

type Config struct {
	EnableVerify    bool    `yaml:"enable_verify"`
	CAPath          string  `yaml:"ca_path"`
	ModuleCertPath  string  `yaml:"module_cert_pem"`
	ModuleKeyPath   string  `yaml:"module_key_pem"`
	SessionExpireIn int     `yaml:"session_expire_in"`
	SessionPath     string  `yaml:"session_path"`
	CookieName      string  `yaml:"cookie_name"`
	UinHeaderName   string  `yaml:"uin_header_name"`
	Routes          []Route `yaml:"routes"`
	// 热更新时间（秒）
	HotLoad int `yaml:"hot_load"`
	HotLoadConfig
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.LoadConfig = cfg.LoadFromFile
	return cfg
}

func (cfg *Config) LoadFrom(in []byte) error {
	if err := yaml.Unmarshal(in, cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) LoadFromFile(p string) error {
	if b, err := ioutil.ReadFile(p); err == nil {
		return cfg.LoadFrom(b)
	} else {
		return err
	}
}
