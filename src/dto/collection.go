package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
)

// 路由组
type Collection struct {
	Id          string `json:"id"`
	ParentId    string `json:"parent_id,omitempty"` // 父级ID
	Name        string `json:"name"`
	Description string `json:"description"`
	// 服务域名
	// 外部服务访问路由的入口
	ServerNames []*ServerName `json:"server_names,omitempty"`
	// 路由信息
	Routes []*Route `json:"routes,omitempty"`
	// 后端服务
	// 集合中没有设置后端服务的路由默认继承集合的后端服务信息
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
	// 鉴权信息
	Authorize *Authorize `json:"authorize,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewCollection(item *entity.Collection) *Collection {
	obj := &Collection{Id: identity.Format(constant.CollectionPrefix, item.Id)}
	obj.Name = item.Name
	obj.Description = item.Description
	obj.ParentId = identity.Format(constant.CollectionPrefix, item.ParentId)
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}
