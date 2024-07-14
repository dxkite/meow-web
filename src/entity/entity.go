package entity

import (
	"time"

	"dxkite.cn/meownest/src/enum"
)

type User struct {
	Id uint64 `gorm:"primarykey"`
	// 用户名
	Name string
	// 密码
	Password string
	// 权限
	Scopes []string `gorm:"serializer:json"`
	// 用户状态
	Status enum.UserStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser() *User {
	entity := new(User)
	return entity
}
