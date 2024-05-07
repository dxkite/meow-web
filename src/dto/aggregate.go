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
	Id string `json:"id"`
	// 证书备注名
	Name string `json:"name"`
	// 证书支持的域名
	DNSNames []string `json:"dns_names"`
	// 证书开启时间
	NotBefore time.Time `json:"not_before"`
	// 证书有效期
	NotAfter time.Time `json:"not_after"`
	// 私钥
	Key string `json:"key,omitempty"`
	// 证书
	Certificate string    `json:"certificate,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewCertificate(item *entity.Certificate) *Certificate {
	obj := &Certificate{
		Id: identity.Format(constant.CertificatePrefix, item.Id),
	}
	obj.Name = item.Name
	obj.NotBefore = item.NotBefore
	obj.NotAfter = item.NotAfter
	obj.DNSNames = item.DNSNames
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}

// 路由信息
type Route struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Method      []string `json:"method"`
	Path        string   `json:"path"`
	// 路由的特殊匹配规则
	Matcher []*value.MatchOption `json:"matcher"`
	// 后端服务
	EndpointId string `json:"endpoint_id,omitempty"`
	// 路由自定义的后端路由
	Endpoint *Endpoint `json:"endpoint,omitempty"`
	// 鉴权信息ID
	AuthorizeId string `json:"authorize_id,omitempty"`
	// 鉴权信息
	Authorize *Authorize `json:"authorize,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewRoute(item *entity.Route) *Route {
	obj := &Route{Id: identity.Format(constant.RoutePrefix, item.Id)}
	obj.Name = item.Name
	obj.Description = item.Description
	obj.Method = item.Method
	obj.Path = item.Path
	obj.Matcher = item.MatchOptions
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}
