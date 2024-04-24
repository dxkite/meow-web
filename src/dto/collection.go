package dto

import "dxkite.cn/meownest/src/model"

type Collection struct {
	Id          uint64      `json:"id"`
	ParentId    uint64      `json:"parent_id"`
	Order       int         `json:"order"`
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Ports       []*Port     `json:"ports,omitempty"`
	Endpoints   []*Endpoint `json:"endpoints"`
}

type Port struct {
	Id            uint64   `json:"id"`
	Name          string   `json:"name"`
	Listen        []string `json:"listen"`
	Authorization string   `json:"authorization"`
}

type Endpoint struct {
	Id uint64 `json:"id"`
	// 后端名
	Name string `json:"name"`
	// 服务类型
	Type model.EndpointType `json:"type"`
	// 转发正则
	ForwardRegex string `json:"forward_regex"`
	// 转发配置
	ForwardReplace string `json:"forward_replace"`
	// 请求头转发配置
	ForwardHeader []*model.ForwardHeaderOption `gorm:"serializer:json" json:"forward_header"`
	// 匹配规则
	Matcher []*model.MatcherConfig `json:"matcher"`
}
