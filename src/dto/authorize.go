package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
)

// Authorize
type Authorize struct {
	Id            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func NewAuthorize(ent *entity.Authorize) *Authorize {
	obj := new(Authorize)
	obj.Id = identity.Format(constant.AuthorizePrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	return obj
}
