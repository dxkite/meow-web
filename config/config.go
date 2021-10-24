package config

import (
	"dxkite.cn/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
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

type CORSConfig struct {
	AllowHeader []string `yaml:"allow_header"`
	AllowOrigin []string `yaml:"allow_origin"`
	AllowMethod []string `yaml:"allow_method"`
}

type SignConfig struct {
	RedirectUrl  string `yaml:"redirect_url"`
	RedirectName string `yaml:"redirect_name"`
}

type SessionConfig struct {
	Name      string `yaml:"name"`
	ExpiresIn int    `yaml:"expires_in"`
	Domain    string `yaml:"domain"`
	Secure    bool   `yaml:"secure"`
}

func (s *SessionConfig) GetName() string {
	if len(s.Name) == 0 {
		return "session"
	}
	return s.Name
}

func (s *SessionConfig) GetExpiresIn() time.Duration {
	expireIn := 24 * time.Hour
	if s.ExpiresIn > 0 {
		expireIn = time.Second * time.Duration(s.ExpiresIn)
	}
	return expireIn
}

type Config struct {
	EnableVerify   bool           `yaml:"enable_verify"`
	Address        string         `yaml:"address"`
	CAPath         string         `yaml:"ca_path"`
	ModuleCertPath string         `yaml:"module_cert_pem"`
	ModuleKeyPath  string         `yaml:"module_key_pem"`
	SessionConfig  *SessionConfig `yaml:"session"`
	UinHeaderName  string         `yaml:"uin_header_name"`
	LogConfig      *LogConfig     `yaml:"log"`
	// 登录配置
	Sign *SignConfig `yaml:"sign_info"`
	// 路由配置
	Routes []Route `yaml:"routes"`
	// HTTP请求头白名单
	HttpAllowHeader []string `yaml:"http_header_allow"`
	// 跨域配置
	Cors *CORSConfig `yaml:"cors_config"`
	// 热更新时间（秒）
	HotLoad int `yaml:"hot_load"`
	HotLoadConfig
}

type LogConfig struct {
	Path  string       `yaml:"path"`
	Level log.LogLevel `yaml:"level"`
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.LoadConfig = cfg.LoadFromFile
	cfg.cfg = cfg
	cfg.Cors = &CORSConfig{}
	cfg.Sign = &SignConfig{}
	cfg.LogConfig = &LogConfig{}
	return cfg
}

func (cfg *Config) Session() *SessionConfig {
	if cfg.SessionConfig == nil {
		cfg.SessionConfig = &SessionConfig{}
	}
	return cfg.SessionConfig
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
