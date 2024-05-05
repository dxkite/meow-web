package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
)

// User
type User struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
}

func NewUser(ent *entity.User) *User {
	obj := new(User)
	obj.Id = identity.Format(constant.UserPrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	obj.Name = ent.Name
	obj.Status = ent.Status
	return obj
}
