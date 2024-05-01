package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/value"
)

// 域名
type ServerName struct {
	Id            string       `json:"id"`
	Name          string       `json:"name"`                     // 域名
	CertificateId string       `json:"certificate_id,omitempty"` // 证书
	Certificate   *Certificate `json:"certificate,omitempty"`    // 证书
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func NewServerName(item *entity.ServerName) *ServerName {
	obj := &ServerName{
		Id: identity.Format(constant.ServerNamePrefix, item.Id),
	}
	obj.Name = item.Name
	obj.CertificateId = identity.Format(constant.CertificatePrefix, item.CertificateId)
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}

// SSL证书
type Certificate struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Domain      []string  `json:"domain"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Key         string    `json:"key,omitempty"`
	Certificate string    `json:"certificate,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewCertificate(item *entity.Certificate) *Certificate {
	obj := &Certificate{
		Id: identity.Format(constant.CertificatePrefix, item.Id),
	}
	obj.Name = item.Name
	obj.StartTime = item.StartTime
	obj.EndTime = item.EndTime
	obj.Domain = item.Domain
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}

// 路由信息
type Route struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Method      []string               `json:"method"`
	Path        string                 `json:"path"`
	Matcher     []*value.MatcherOption `json:"matcher"`
	Endpoints   []*Endpoint            `json:"endpoints,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func NewRoute(item *entity.Route) *Route {
	obj := &Route{Id: identity.Format(constant.RoutePrefix, item.Id)}
	obj.Name = item.Name
	obj.Description = item.Description
	obj.Method = item.Method
	obj.Path = item.Path
	obj.Matcher = item.Matcher
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}

// 路由组
type Collection struct {
	Id          string `json:"id"`
	ParentId    string `json:"parent_id,omitempty"` // 父级ID
	Name        string `json:"name"`
	Description string `json:"description"`
	// 路由信息
	ServerNames []*ServerName `json:"server_names,omitempty"`
	// 路由信息
	Routes []*Route `json:"routes,omitempty"`
	// 后端服务
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
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

// 后端配置
type Endpoint struct {
	Id string `json:"id"`
	// 后端名
	Name string `json:"name"`
	// 服务类型
	Type string `json:"type"`
	// 重写配置
	ForwardRewrite *value.ForwardRewriteOption `json:"forward_rewrite"`
	// 请求头转发配置
	ForwardHeader []*value.ForwardHeaderOption `json:"forward_header"`
	// 匹配规则
	Matcher []*value.MatcherOption `json:"matcher"`
	// 远程服务
	Endpoint *value.ForwardEndpoint `gorm:"serializer:json" json:"endpoint"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewEndpoint(item *entity.Endpoint) *Endpoint {
	obj := &Endpoint{Id: identity.Format(constant.EndpointPrefix, item.Id)}
	obj.Name = item.Name
	obj.Type = item.Type
	obj.ForwardRewrite = item.ForwardRewrite
	obj.ForwardHeader = item.ForwardHeader
	obj.ForwardRewrite = item.ForwardRewrite
	obj.Endpoint = item.Endpoint
	obj.Matcher = item.Matcher
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}
