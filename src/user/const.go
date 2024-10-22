package user

import "dxkite.cn/nebula/pkg/httpx"

const UserPrefix = "user_"
const SessionPrefix = "session_"

const (
	ScopeUserRead  httpx.ScopeName = "user:read"
	ScopeUserWrite httpx.ScopeName = "user:write"
)
