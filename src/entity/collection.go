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

	// 合辑名
	Name        string `gorm:"index"`
	Description string
}

func NewCollection() *Collection {
	entity := new(Collection)
	return entity
}
