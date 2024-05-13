package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/enum"
)

// User
type User struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 用户名
	Name string `json:"name"`
	// 用户权限
	Scope []string `json:"scope"`
	// 用户状态
	Status enum.UserStatus `json:"status"`
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
