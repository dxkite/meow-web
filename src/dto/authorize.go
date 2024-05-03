package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/value"
)

// Authorize
type Authorize struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name      string                    `json:"name"`
	Type      string                    `json:"type"`
	Attribute *value.AuthorizeAttribute `json:"attribute"`
}

func NewAuthorize(ent *entity.Authorize) *Authorize {
	obj := new(Authorize)
	obj.Id = identity.Format(constant.AuthorizePrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	obj.Name = ent.Name
	obj.Type = ent.Type
	obj.Attribute = ent.Attribute
	return obj
}
