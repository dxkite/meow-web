package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/model"
	"dxkite.cn/meownest/src/valueobject"
)

// 服务
type Server struct {
	Id          string        `json:"id"`
	ServerName  []*ServerName `json:"server_name"`
	Collections []*Collection `json:"collections"`
	Endpoint    []*Endpoint   `json:"endpoint,omitempty"`
}

// 域名
type ServerName struct {
	Id            string       `json:"id"`
	Name          string       `json:"name"`                     // 域名
	Protocol      string       `json:"protocol"`                 // 协议
	CertificateId string       `json:"certificate_id,omitempty"` // 证书
	Certificate   *Certificate `json:"certificate,omitempty"`    // 证书
}

func NewServerName(cert *model.ServerName) *ServerName {
	rst := &ServerName{
		Id: identity.Format(constant.ServerNamePrefix, cert.Id),
	}
	rst.Name = cert.Name
	rst.Protocol = cert.Protocol
	rst.CertificateId = identity.Format(constant.CertificatePrefix, cert.CertificateId)
	return rst
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
}

func NewCertificate(cert *model.Certificate) *Certificate {
	rst := &Certificate{
		Id: identity.Format(constant.CertificatePrefix, cert.Id),
	}
	rst.Name = cert.Name
	rst.StartTime = cert.StartTime
	rst.EndTime = cert.EndTime
	rst.Domain = cert.Domain
	return rst
}

// 路由信息
type Route struct {
	Id          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Method      []string                     `json:"method"`
	Path        string                       `json:"path"`
	Matcher     []*valueobject.MatcherConfig `json:"matcher"`
	Endpoint    []*Endpoint                  `json:"endpoints,omitempty"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
}

func NewRoute(item *model.Route) *Route {
	obj := &Route{Id: identity.Format(constant.CollectionPrefix, item.Id)}
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
	Id          string        `json:"id"`
	ParentId    string        `json:"parent_id,omitempty"` // 父级ID
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Routes      []*Route      `json:"routes,omitempty"`
	Collections []*Collection `json:"collections,omitempty"`
	Endpoint    []*Endpoint   `json:"endpoints,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func NewCollection(item *model.Collection) *Collection {
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
	// 转发正则
	ForwardRegex string `json:"forward_regex"`
	// 转发配置
	ForwardReplace string `json:"forward_replace"`
	// 请求头转发配置
	ForwardHeader []*ForwardHeaderOption `json:"forward_header"`
	// 匹配规则
	MatchFilter []*MatchOption `json:"match_filter"`
}

type ForwardHeaderOption struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MatchOption struct {
	Source string `json:"source"` // 匹配源
	Name   string `json:"name"`   // 匹配值
	Type   string `json:"type"`   // 匹配方式
	Value  string `json:"value"`  // 匹配内容
}
