package entity

import (
	"time"
)

type Authorize struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// TODO
}

func NewAuthorize() (*Authorize, error) {
	entity := new(Authorize)
	return entity, nil
}
