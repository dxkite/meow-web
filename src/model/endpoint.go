package model

import "gorm.io/gorm"

type EndpointType string

const (
	EndpointTypeStatic = "static"
)

type EndpointConfig struct {
	gorm.Model
	// 服务类型
	Type EndpointType `json:"type"`
	// 转发正则
	ForwardRegex string `json:"forward_regex"`
	// 转发配置
	ForwardReplace string `json:"forward_replace"`
	// 请求头转发配置
	ForwardHeader []*ForwardHeaderOption `gorm:"serializer:json" json:"forward_header"`
	// 匹配规则
	Matcher []*MatcherConfig `json:"matcher"`
}

type ForwardHeaderOption struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
