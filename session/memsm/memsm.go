package memsm

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/session"
)

type MemSm struct {
	sess map[uint64]struct{}
}

func NewMemSM(*config.Config) (session.SessionManager, error) {
	return &MemSm{sess: map[uint64]struct{}{}}, nil
}

func (sm *MemSm) CreateSession(uin uint64) error {
	sm.sess[uin] = struct{}{}
	return nil
}

func (sm *MemSm) CheckSession(uin uint64) bool {
	_, ok := sm.sess[uin]
	return ok
}

func (sm *MemSm) RemoveSession(uin uint64) error {
	delete(sm.sess, uin)
	return nil
}
