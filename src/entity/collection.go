package entity

import (
	"time"
)

type Collection struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 树型节点部分
	ParentId uint64 `gorm:"index"`
	Index    string `gorm:"index"`
	Order    int
	Depth    int

	// 集合名
	Name string `gorm:"index"`
	// 集合描述
	Description string
	// 绑定域名
	ServerNames []string `gorm:"serializer:json"`
	// 权限配置ID
	AuthorizeId uint64 `gorm:"index"`
	// 后端服务ID
	EndpointId uint64 `gorm:"index"`
}

func NewCollection() *Collection {
	entity := new(Collection)
	return entity
}
