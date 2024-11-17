package user

import "dxkite.cn/nebula/pkg/depends"

func init() {
	depends.Register(NewUserRepository)
	depends.Register(NewSessionRepository)
	depends.Register(NewUserService)
	depends.Register(NewUserExpander)
}
