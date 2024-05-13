package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
)

// Session
type Session struct {
	Id            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func NewSession(ent *entity.Session) *Session {
	obj := new(Session)
	obj.Id = identity.Format(constant.SessionPrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	return obj
}
