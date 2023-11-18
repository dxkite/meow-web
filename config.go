package suda

type Config struct {
	// 日志文件
	LogFile string `yaml:"log_file"`
	// 日志等级
	LogLevel int `yaml:"log_level"`

	// 服务配置
	Services []ServiceConfig `yaml:"services"`
	// 实例配置
	Instances []InstanceConfig `yaml:"instances"`
}

type AuthAesConfig struct {
	Key string `yaml:"key"`
}

type AuthSourceConfig struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
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

type ServiceConfig struct {
	Name   string        `yaml:"name"`
	Auth   AuthConfig    `yaml:"auth"`
	Ports  []Port        `yaml:"ports"`
	Routes []RouteConfig `yaml:"routes"`
}

type RouteConfig struct {
	Name      string        `yaml:"name"`
	Auth      bool          `yaml:"auth"`
	Rewrite   RewriteConfig `yaml:"rewrite"`
	EndPoints []Port        `yaml:"endpoints"`
	Paths     []string      `yaml:"paths"`
}

type InstanceConfig struct {
	Name string   `yaml:"name"`
	Exec []string `yaml:"exec"`
}
