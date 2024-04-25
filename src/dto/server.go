package dto

import (
	"time"
)

// 服务
type Server struct {
	Id          string         `json:"id"`
	Port        []*Port        `json:"port"`
	RouterGroup []*RouterGroup `json:"router_group"`
	Endpoint    []*Endpoint    `json:"endpoint,omitempty"`
}

type Port struct {
	Id            string        `json:"id"`
	ServerName    []*ServerName `json:"server_name"`
	Listen        []string      `json:"listen"`
	Authorization string        `json:"authorization"`
}

// 域名管理
type ServerName struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Certificate *Certificate `json:"certificate"`
}

// 域名证书
type Certificate struct {
	Id          string    `json:"id"`
	Domain      []string  `json:"domain"`
	ExpireAt    time.Time `json:"expire_at"`
	Key         string    `json:"key"`
	Certificate string    `json:"certificate"`
}

// 路由项
type Router struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Method   []string    `json:"method"`
	Path     string      `json:"path"`
	Endpoint []*Endpoint `json:"endpoints,omitempty"`
}

// 路由组
type RouterGroup struct {
	Id          string         `json:"id"`
	ParentId    string         `json:"parent_id,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Routes      []*Router      `json:"routes"`
	Children    []*RouterGroup `json:"children,omitempty"`
	Endpoint    []*Endpoint    `json:"endpoints,omitempty"`
}

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
