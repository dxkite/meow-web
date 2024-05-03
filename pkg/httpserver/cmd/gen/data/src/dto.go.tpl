package dto

import (
	"time"

	"{{ .Pkg }}/pkg/identity"
	"{{ .Pkg }}/src/constant"
	"{{ .Pkg }}/src/entity"
)

// {{ .Name }}
type {{ .Name }} struct {
	Id            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func New{{ .Name }}(ent *entity.{{ .Name }}) *{{ .Name }} {
	obj := new({{ .Name }})
	obj.Id = identity.Format(constant.{{ .Name }}Prefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	return obj
}
