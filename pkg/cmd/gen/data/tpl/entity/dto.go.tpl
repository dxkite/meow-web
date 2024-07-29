package {{ .ModuleName }}

import (
	"time"

	"{{ .PackageName }}/pkg/crypto/identity"
)

// {{ .Name }}
type {{ .Name }}Dto struct {
	Id            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func New{{ .Name }}Dto(ent *{{ .Name }}) *{{ .Name }}Dto {
	obj := new({{ .Name }}Dto)
	obj.Id = identity.Format({{ .Name }}Prefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	return obj
}
