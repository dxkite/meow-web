package leveldbsm

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/session"
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbSM struct {
	db *leveldb.DB
}

func NewLevelDbSM(cfg *config.Config) (session.SessionManager, error) {
	db, err := leveldb.OpenFile(cfg.SessionPath, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDbSM{db: db}, nil
}

func (sm *LevelDbSM) CreateSession(uin uint64) error {
	id := make([]byte, 8)
	binary.BigEndian.PutUint64(id, uin)
	return sm.db.Put(id, []byte{1}, nil)
}

func (sm *LevelDbSM) CheckSession(uin uint64) bool {
	id := make([]byte, 8)
	binary.BigEndian.PutUint64(id, uin)
	_, err := sm.db.Get(id, nil)
	return err == nil
}

func (sm *LevelDbSM) RemoveSession(uin uint64) error {
	id := make([]byte, 8)
	binary.BigEndian.PutUint64(id, uin)
	return sm.db.Delete(id, nil)
}
