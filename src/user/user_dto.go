package user

import (
	"time"

	"dxkite.cn/nebula/pkg/crypto/identity"
	"dxkite.cn/nebula/pkg/httputil"
)

// UserDto
type UserDto struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 用户名
	Name string `json:"name"`
	// 用户权限
	Scopes []httputil.ScopeName `json:"scopes"`
	// 用户状态
	Status UserStatus `json:"status"`
}

func NewUserDto(ent *User) *UserDto {
	obj := new(UserDto)
	obj.Id = identity.Format(UserPrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	obj.Name = ent.Name
	obj.Scopes = ent.Scopes
	obj.Status = ent.Status
	return obj
}
