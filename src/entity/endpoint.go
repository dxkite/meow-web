package entity

import (
	"dxkite.cn/meownest/src/valueobject"
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
	ForwardRewrite *valueobject.ForwardRewriteOption `gorm:"serializer:json" json:"forward_rewrite"`
	// 请求头转发配置
	ForwardHeader []*valueobject.ForwardHeaderOption `gorm:"serializer:json" json:"forward_header"`
	// 匹配规则
	Matcher []*valueobject.MatcherOption `gorm:"serializer:json" json:"matcher"`
	// 远程服务
	Endpoint *valueobject.ForwardEndpoint `gorm:"serializer:json" json:"endpoint"`
}
