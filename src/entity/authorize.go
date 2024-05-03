package entity

import (
	"time"

	"dxkite.cn/meownest/src/value"
)

type Authorize struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	Type      string
	Attribute *value.AuthorizeAttribute `gorm:"serializer:json"`
}

func NewAuthorize() *Authorize {
	entity := new(Authorize)
	return entity
}
