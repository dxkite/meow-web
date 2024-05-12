package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/enum"
	"dxkite.cn/meownest/src/value"
)

// 后端配置
type Endpoint struct {
	Id string `json:"id"`
	// 后端名
	Name string `json:"name"`
	// 服务类型
	Type enum.EndpointType `json:"type"`
	// 远程服务
	Endpoint *value.ForwardEndpoint `gorm:"serializer:json" json:"endpoint"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewEndpoint(item *entity.Endpoint) *Endpoint {
	obj := &Endpoint{Id: identity.Format(constant.EndpointPrefix, item.Id)}
	obj.Name = item.Name
	obj.Type = item.Type
	obj.Endpoint = item.Endpoint
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}
