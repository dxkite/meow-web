package session

type SessionManager interface {
	CreateSession(uin uint64) error
	CheckSession(uin uint64) bool
	RemoveSession(uin uint64) error
}
