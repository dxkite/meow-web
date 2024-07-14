package user

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
)

// SessionDto
type SessionDto struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSessionDto(ent *entity.Session) *SessionDto {
	obj := new(SessionDto)
	obj.Id = identity.Format(constant.SessionPrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	return obj
}
