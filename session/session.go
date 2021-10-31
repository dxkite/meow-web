package session

import "time"

type SessionManager interface {
	CreateSession(uin uint64, expires time.Time) error
	CheckSession(uin uint64) bool
	RemoveSession(uin uint64) error
}
