package suda

import (
	"context"
)

func Bootstrap(ctx context.Context, configPath string) error {
	cfg := Config{}
	if err := loadYaml(configPath, &cfg); err != nil {
		return err
	}

	applyLogConfig(ctx, cfg.LogLevel, cfg.LogFile)

	execChain := ExecChain{}

	execChain = append(execChain, func() error {
		return RunInstance(ctx, cfg.Instances)
	})

	execChain = append(execChain, func() error {
		return RunService(ctx, cfg.Services)
	})

	return execChain.Run()
}

func RunService(ctx context.Context, services []ServiceConfig) error {
	execChain := ExecChain{}

	for _, cfg := range services {
		execChain = append(execChain, (func(cfg ServiceConfig) func() error {
			return func() error {
				srv := new(Service)
				if err := srv.Config(&cfg); err != nil {
					return err
				}
				return srv.Run()
			}
		})(cfg))
	}

	return execChain.Run()
}

func RunInstance(ctx context.Context, instance []InstanceConfig) error {
	execChain := ExecChain{}

	for _, ins := range instance {
		execChain = append(execChain, (func(ins InstanceConfig) func() error {
			return func() error {
				return execInstance(&ins)
			}
		})(ins))
	}

	return execChain.Run()
}
