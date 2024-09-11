package migrate

import "dxkite.cn/meow-web/src/user"

func init() {
	dst = append(dst, user.User{})
	dst = append(dst, user.Session{})
}
