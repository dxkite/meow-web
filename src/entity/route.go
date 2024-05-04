package entity

import "dxkite.cn/meownest/src/value"

type Route struct {
	Base
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Method      []string `json:"method" gorm:"serializer:json"`
	Path        string   `json:"path"`
	// 匹配规则
	MatchOptions []*value.MatchOption `json:"match_options" gorm:"serializer:json"` // 权限配置ID
	// 所属集合ID
	CollectionId uint64 `gorm:"index"`
	// 权限配置ID
	AuthorizeId uint64 `gorm:"index"`
	// 后端服务ID
	EndpointId uint64 `gorm:"index"`
}
