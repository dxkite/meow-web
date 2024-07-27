package user

import "dxkite.cn/meownest/pkg/httputil"

const UserPrefix = "user_"
const SessionPrefix = "session_"

const (
	ScopeUserRead  httputil.ScopeName = "user:read"
	ScopeUserWrite httputil.ScopeName = "user:write"
)
