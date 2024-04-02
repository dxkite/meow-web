package bootstrap

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"dxkite.cn/log"
	"dxkite.cn/meownest/src/config"
	"dxkite.cn/meownest/src/executer"
	"dxkite.cn/meownest/src/gateway"
	"dxkite.cn/meownest/src/utils"
	"gopkg.in/yaml.v3"
)

func ServeGateway(ctx context.Context, configPath string) error {
	cfg := &config.GatewayConfig{}
	if err := loadYaml(configPath, &cfg); err != nil {
		return err
	}

	applyLogConfig(ctx, cfg.LogLevel, cfg.LogFile)

	execChain := utils.ExecChain{}

	execChain = append(execChain, func() error {
		return StartGatewayHttpServer(ctx, cfg)
	})

	execChain = append(execChain, func() error {
		return StartComponentList(ctx, cfg.Components)
	})

	return execChain.Run()
}

func StartGatewayHttpServer(ctx context.Context, cfg *config.GatewayConfig) error {
	server := gateway.NewHttpServer()
	server.RegisterAuthorizationHandler(cfg.HttpAuthorization.Header, &gateway.HttpAesHandler{
		Key:    cfg.HttpAuthorization.AesToken.Key,
		Query:  cfg.HttpAuthorization.Source.Query,
		Header: cfg.HttpAuthorization.Source.Header,
		Cookie: cfg.HttpAuthorization.Source.Cookie,
	})

	for _, v := range cfg.HttpRouter {
		entry := &gateway.HttpRouterGroupEntry{
			Name:          v.Name,
			Hostname:      v.Hostname,
			Authorization: v.Authorization,
			Endpoints:     v.Endpoints,
			Paths:         v.Paths,
		}
		if v.Rewrite != nil {
			entry.Rewrite = &gateway.RewriteConfig{
				Regex:   v.Rewrite.Regex,
				Replace: v.Rewrite.Replace,
			}
		}
		if v.Matcher != nil {
			entry.Matcher = &gateway.MatcherConfig{
				Query:  v.Matcher.Query,
				Cookie: v.Matcher.Cookie,
				Header: v.Matcher.Header,
			}
		}
		server.RegisterRouterGroup(entry)
	}

	return server.Serve(cfg.Listen)
}

func StartComponentList(ctx context.Context, list []*config.Component) error {
	execChain := utils.ExecChain{}

	for _, comp := range list {
		execChain = append(execChain, (func(comp *config.Component) func() error {
			return func() error {
				return executer.ExecCommandWithName(comp.Name, comp.Exec)
			}
		})(comp))
	}

	return execChain.Run()
}

func loadYaml(name string, data interface{}) error {
	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, data); err != nil {
		return err
	}

	return nil
}

func applyLogConfig(ctx context.Context, level int, output string) {
	if level != 0 {
		log.SetLevel(log.LogLevel(level))
	}
	if output == "" {
		return
	}
	log.Println("log output file", output)
	filename := output
	var w io.Writer
	if f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err != nil {
		log.Warn("log file open error", filename)
		return
	} else {
		w = f
		if filepath.Ext(filename) == ".json" {
			w = log.NewJsonWriter(w)
		} else {
			w = log.NewTextWriter(w)
		}
		go func() {
			<-ctx.Done()
			_ = f.Close()
		}()
	}
	log.SetOutput(log.MultiWriter(w, log.Writer()))
}