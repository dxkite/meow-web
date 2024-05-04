package entity

import (
	"dxkite.cn/meownest/src/value"
)

const (
	EndpointTypeStatic = "static"
)

type Endpoint struct {
	Base
	Name string `json:"name"`
	// 服务类型
	Type string `json:"type"`
	// 重写配置
	ForwardRewrite *value.ForwardRewriteOption `gorm:"serializer:json" json:"forward_rewrite"`
	// 请求头转发配置
	ForwardHeader []*value.ForwardHeaderOption `gorm:"serializer:json" json:"forward_header"`
	// 匹配规则
	MatchOptions []*value.MatchOption `gorm:"serializer:json" json:"match_options"`
	// 远程服务
	Endpoint *value.ForwardEndpoint `gorm:"serializer:json" json:"endpoint"`
}
