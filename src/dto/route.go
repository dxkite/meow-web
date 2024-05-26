package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/enum"
	"dxkite.cn/meownest/src/value"
)

// 路由信息
type Route struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Method      []string `json:"method"`
	// 路径
	Path string `json:"path"`
	// 路径类型
	PathType enum.RoutePathType `json:"path_type"`
	// 路由的特殊匹配规则
	MatchOptions []*value.MatchOption `json:"match_options"`
	// 路由重写规则
	PathRewrite *value.PathRewrite `json:"path_rewrite,omitempty"`
	// 数据重写规则
	ModifyOptions []*value.ModifyOption `json:"modify_options"`
	// 后端服务
	EndpointId string `json:"endpoint_id"`
	// 路由自定义的后端路由
	Endpoint *Endpoint `json:"endpoint,omitempty"`
	// 鉴权信息ID
	AuthorizeId string `json:"authorize_id"`
	// 鉴权信息
	Authorize *Authorize `json:"authorize,omitempty"`
	// 分组ID
	CollectionId string `json:"collection_id"`
	// 状态
	Status enum.RouteStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewRoute(item *entity.Route) *Route {
	obj := &Route{Id: identity.Format(constant.RoutePrefix, item.Id)}
	obj.Name = item.Name
	obj.Description = item.Description
	obj.Method = item.Method
	obj.Path = item.Path
	obj.PathType = item.PathType
	obj.MatchOptions = item.MatchOptions
	obj.PathRewrite = item.PathRewrite
	obj.ModifyOptions = item.ModifyOptions
	obj.CollectionId = identity.Format(constant.CollectionPrefix, item.CollectionId)
	obj.Status = item.Status
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	obj.EndpointId = identity.Format(constant.EndpointPrefix, item.CollectionId)
	obj.AuthorizeId = identity.Format(constant.AuthorizePrefix, item.AuthorizeId)
	return obj
}
