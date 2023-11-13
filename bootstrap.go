package suda

import (
	"context"
	"sync"

	"dxkite.cn/log"
)

func Bootstrap(ctx context.Context, configPath string) error {
	cfg := Config{}
	if err := loadYaml(configPath, &cfg); err != nil {
		return err
	}

	applyLogConfig(ctx, cfg.LogLevel, cfg.LogFile)

	wait := &sync.WaitGroup{}
	wait.Add(len(cfg.Services))

	for _, srvCfg := range cfg.Services {
		go func(srvCfg *ServiceConfig) {
			defer wait.Done()

			srv := new(Service)
			if err := srv.Config(srvCfg); err != nil {
				log.Error(err)
				return
			}

			if err := srv.Run(); err != nil {
				log.Error(err)
			}
		}(&srvCfg)
	}

	wait.Wait()
	return nil
}
