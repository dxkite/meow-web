package suda

type Config struct {
	Addr         string     `yaml:"addr"`
	ModuleConfig string     `yaml:"module_config"`
	Auth         AuthConfig `yaml:"auth"`

	// 日志文件
	LogFile string `yaml:"log_file"`
	// 日志等级
	LogLevel int `yaml:"log_level"`
}

type ModuleConfig struct {
	Name      string        `yaml:"name"`
	EndPoints []string      `yaml:"endpoints"`
	Exec      []string      `yaml:"exec"`
	Routes    []RouteConfig `yaml:"routes"`
}

type RouteConfig struct {
	Name    string        `yaml:"name"`
	Auth    bool          `yaml:"auth"`
	Rewrite RewriteConfig `yaml:"rewrite"`
	Paths   []string      `yaml:"paths"`
}

type RewriteConfig struct {
	Regex   string `yaml:"regex"`
	Replace string `yaml:"replace"`
}

type AuthConfig struct {
	Type   string             `yaml:"type"`
	Header string             `yaml:"header"`
	Source []AuthSourceConfig `yaml:"source"`
	Aes    AuthAesConfig      `yaml:"aes"`
}

type AuthAesConfig struct {
	Key string `yaml:"key"`
}

type AuthSourceConfig struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}
