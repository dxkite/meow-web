package entity

import (
	"time"
)

type User struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name     string
	Password string
	Status   string
}

func NewUser() *User {
	entity := new(User)
	return entity
}
