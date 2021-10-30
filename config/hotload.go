package config

import (
	"dxkite.cn/log"
	"os"
	"sync"
	"time"
)

type ConfigChangeCallback func(cfg interface{})
type ConfigLoadCallback func(p string) error

type HotLoadConfig struct {
	// 更新时间
	modifyTime time.Time
	mtx        sync.Mutex
	LoadConfig ConfigLoadCallback
	changCb    []ConfigChangeCallback
	loadTime   int
	cur        string
	cfg        interface{}
}

func NewHotLoad(cb ConfigLoadCallback) *HotLoadConfig {
	return &HotLoadConfig{
		mtx:        sync.Mutex{},
		LoadConfig: cb,
	}
}

func (cfg *HotLoadConfig) OnChange(cb ConfigChangeCallback) {
	if cfg.changCb == nil {
		cfg.changCb = []ConfigChangeCallback{}
	}
	cfg.changCb = append(cfg.changCb, cb)
}

func (cfg *HotLoadConfig) applyConfig() {
	for _, cb := range cfg.changCb {
		cb(cfg.cfg)
	}
}

func (cfg *HotLoadConfig) SetLastLoadTime(t int) {
	cfg.loadTime = t
}

func (cfg *HotLoadConfig) SetLastLoadFile(p string) {
	cfg.cur = p
}

func (cfg *HotLoadConfig) NotifyModify() {
	go cfg.applyConfig()
}

func (cfg *HotLoadConfig) LoadIfModify(p string) (bool, error) {
	update := true
	if info, err := os.Stat(p); err != nil {
		return false, err
	} else {
		update = info.ModTime().After(cfg.modifyTime)
		cfg.modifyTime = info.ModTime()
	}
	if !update {
		return false, nil
	}
	err := cfg.LoadConfig(p)
	if err == nil {
		cfg.NotifyModify()
	}
	return true, err
}

func (cfg *HotLoadConfig) HotLoadIfModify() {
	go func() {
		log.Info("enable hot load config", cfg.cur, cfg.loadTime)
		ticker := time.NewTicker(time.Duration(cfg.loadTime) * time.Second)
		for range ticker.C {
			if ok, err := cfg.LoadIfModify(cfg.cur); err != nil {
				log.Error("load config error", err)
			} else if ok {
				log.Println("config hot load success")
			}
		}
	}()
}
