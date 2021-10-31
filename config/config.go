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
	AllowHeader      []string `yaml:"allow_header"`
	AllowOrigin      []string `yaml:"allow_origin"`
	AllowMethod      []string `yaml:"allow_method"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

type SignConfig struct {
	RedirectUrl  string `yaml:"redirect_url"`
	RedirectName string `yaml:"redirect_name"`
}

const AESKeySize = 32

type SessionConfig struct {
	Name      string `yaml:"name"`
	ExpiresIn int    `yaml:"expires_in"`
	Domain    string `yaml:"domain"`
	Secure    bool   `yaml:"secure"`
	HttpOnly  bool   `yaml:"http_only"`
	Path      string `yaml:"path"`

	// 会话加密
	Mode    string `yaml:"mode"`
	AesKey  string `yaml:"aes_key"`
	RsaKey  string `yaml:"rsa_key"`
	RsaCert string `yaml:"rsa_cert"`

	// 严格模式，会话必须在内存中存在
	Strict bool `yaml:"strict"`
	// SLO 单点登出检擦
	SloUrl string `yaml:"slo_url"`
	// SLO 检查超时
	SloTimeout int `yaml:"slo_timeout"`
	// SLO 会话超时
	SloExpiresIn int `yaml:"slo_expires_in"`
}

func (s *SessionConfig) GetName() string {
	if len(s.Name) == 0 {
		return "session"
	}
	return s.Name
}

func (s *SessionConfig) GetPath() string {
	if len(s.Path) != 0 {
		return s.Path
	}
	return "/"
}

func (s *SessionConfig) AesTicketKey() []byte {
	if len(s.AesKey) != AESKeySize {
		return nil
	}
	return []byte(s.AesKey)
}

func (s *SessionConfig) GetExpiresIn() time.Duration {
	expireIn := 24 * time.Hour
	if s.ExpiresIn > 0 {
		expireIn = time.Second * time.Duration(s.ExpiresIn)
	}
	return expireIn
}

func (s *SessionConfig) GetSloExpiresIn() time.Duration {
	expireIn := 1 * time.Hour
	if s.SloExpiresIn > 0 {
		expireIn = time.Second * time.Duration(s.SloExpiresIn)
	}
	return expireIn
}

func (s *SessionConfig) GetSloTimeout() time.Duration {
	expireIn := 100 * time.Millisecond
	if s.SloTimeout > 0 {
		expireIn = time.Millisecond * time.Duration(s.SloTimeout)
	}
	return expireIn
}

type Config struct {
	// TLS配置
	EnableTls bool   `yaml:"enable_tls"`
	TlsCa     string `yaml:"tls_ca"`
	TlsCert   string `yaml:"tls_cert"`
	TlsKey    string `yaml:"tls_key"`
	// 监听地址
	Address string `yaml:"address"`
	// 会话配置
	SessionConfig *SessionConfig `yaml:"session"`
	// UIN字段
	UinHeaderName string `yaml:"uin_header_name"`
	// 日志配置
	LogConfig *LogConfig `yaml:"log"`
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
	*HotLoadConfig
}

type LogConfig struct {
	Path  string       `yaml:"path"`
	Level log.LogLevel `yaml:"level"`
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.HotLoadConfig = NewHotLoad(cfg.LoadFromFile)
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
		cfg.SetLastLoadFile(p)
		return cfg.LoadFrom(b)
	} else {
		return err
	}
}
