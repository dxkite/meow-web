package config

import (
	"dxkite.cn/log"
	"os"
	"sync"
	"time"
)

type ConfigChangeCallback func(cfg interface{})

type HotLoadConfig struct {
	// 更新时间
	modifyTime time.Time
	mtx        sync.Mutex
	LoadConfig func(p string) error
	changCb    []ConfigChangeCallback
	loadTime   int
}

func (cfg *HotLoadConfig) OnChange(cb ConfigChangeCallback) {
	if cfg.changCb == nil {
		cfg.changCb = []ConfigChangeCallback{}
	}
	cfg.changCb = append(cfg.changCb, cb)
}

func (cfg *HotLoadConfig) applyConfig() {
	for _, cb := range cfg.changCb {
		cb(cfg)
	}
}

func (cfg *HotLoadConfig) SetLoadTime(t int) {
	cfg.loadTime = t
}

func (cfg *HotLoadConfig) notifyModify() {
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
	if err != nil {
		cfg.notifyModify()
	}
	return true, err
}

func (cfg *HotLoadConfig) HotLoadIfModify(p string) {
	go func() {
		log.Info("enable hot load config", p)
		ticker := time.NewTicker(time.Duration(cfg.loadTime) * time.Second)
		for range ticker.C {
			if ok, err := cfg.LoadIfModify(p); err != nil {
				log.Error("load config error", err)
			} else if ok {
				log.Println("config hot load success")
			}
		}
	}()
}
