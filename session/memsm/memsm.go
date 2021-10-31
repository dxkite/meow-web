package memsm

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/session"
	"time"
)

type MemSm struct {
	sess map[uint64]time.Time
}

func NewMemSM(*config.Config) (session.SessionManager, error) {
	return &MemSm{sess: map[uint64]time.Time{}}, nil
}

func (sm *MemSm) CreateSession(uin uint64, expires time.Time) error {
	sm.sess[uin] = expires
	return nil
}

func (sm *MemSm) CheckSession(uin uint64) bool {
	t, ok := sm.sess[uin]
	if ok && time.Now().After(t) {
		delete(sm.sess, uin)
		return false
	}
	return ok
}

func (sm *MemSm) RemoveSession(uin uint64) error {
	delete(sm.sess, uin)
	return nil
}
