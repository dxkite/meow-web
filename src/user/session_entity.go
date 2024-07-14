package user

import (
	"time"
)

type Session struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uint64
	Address   string
	Agent     string
	ExpireAt  time.Time
	Deleted   int
}

func NewSession() *Session {
	entity := new(Session)
	return entity
}
