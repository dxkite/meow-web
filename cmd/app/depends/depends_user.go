package depends

import (
	"dxkite.cn/meow-web/src/user"
	"dxkite.cn/nebula/pkg/depends"
)

func init() {
	depends.Register(user.NewUserRepository)
	depends.Register(user.NewSessionRepository)
	depends.Register(user.NewUserService)
}
