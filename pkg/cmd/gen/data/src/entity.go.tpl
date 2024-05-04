package entity

import (
	"time"
)

type {{ .Name }} struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// TODO
}

func New{{ .Name }}() *{{ .Name }} {
	entity := new({{ .Name }})
	return entity
}
