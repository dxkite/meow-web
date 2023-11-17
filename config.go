package suda

import (
	"fmt"
)

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

type Port struct {
	Type string   `yaml:"type"`
	Unix UnixPort `yaml:"unix"`
	Http HttpPort `yaml:"http"`
}

type UnixPort struct {
	Path string `yaml:"path"`
}

type HttpPort struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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

func (port Port) String() string {
	switch port.Type {
	case "unix":
		return fmt.Sprintf("unix://%s", port.Unix.Path)
	case "http":
		return fmt.Sprintf("http://%s:%d", port.Http.Host, port.Http.Port)
	default:
		return port.Type
	}
}
