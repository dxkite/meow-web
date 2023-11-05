package suda

type Config struct {
	Addr         string `yaml:"addr"`
	ModuleConfig string `yaml:"module_config"`
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
