package entity

import (
	"dxkite.cn/meownest/src/enum"
	"dxkite.cn/meownest/src/value"
)

const (
	EndpointTypeStatic = "static"
)

type Endpoint struct {
	Base
	// 服务名称
	Name string `json:"name"`
	// 服务类型
	Type enum.EndpointType `json:"type"`
	// 远程服务
	Endpoint *value.ForwardEndpoint `gorm:"serializer:json" json:"endpoint"`
}
