package config

type GatewayConfig struct {
	// 日志文件
	LogFile string `yaml:"log_file"`
	// 日志等级
	LogLevel int `yaml:"log_level"`
	// 监控
	Listen string `yaml:"listen"`
	// 组件配置
	Components []Component `yaml:"components"`
	// Http 鉴权配置
	HttpAuthorization HttpAuthorizationConfig `yaml:"http-authorization"`
	// Http 路由配置
	HttpRouter []HttpRouterGroupConfig `yaml:"http-router"`
}

type HttpAuthorizationConfig struct {
	Type     string          `yaml:"type"`
	Header   string          `yaml:"header"`
	Source   TokenSource     `yaml:"source"`
	AesToken *AesTokenConfig `yaml:"aes-token"`
}

type AesTokenConfig struct {
	Key string `yaml:"key"`
}

type TokenSource struct {
	Query  []string `yaml:"query"`
	Header []string `yaml:"header"`
	Cookie []string `yaml:"cookie"`
}

type Component struct {
	Name string   `yaml:"name"`
	Exec []string `yaml:"exec"`
}

type HttpRouterGroupConfig struct {
	Name          string        `yaml:"name"`
	Hostname      []string      `yaml:"hostname"`
	Authorization bool          `yaml:"authorization"`
	Endpoints     []string      `yaml:"endpoints"`
	Rewrite       RewriteConfig `yaml:"rewrite"`
	Paths         []string      `yaml:"paths"`
}

type RewriteConfig struct {
	Regex   string `yaml:"regex"`
	Replace string `yaml:"replace"`
}
