package user

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
)

// SessionDto
type SessionDto struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSessionDto(ent *Session) *SessionDto {
	obj := new(SessionDto)
	obj.Id = identity.Format(SessionPrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	return obj
}
