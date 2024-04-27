package entity

import "time"

type Base struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
