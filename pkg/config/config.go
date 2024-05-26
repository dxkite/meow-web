package config

import (
	"context"
	"errors"
	"strconv"
)

type configProviderKey string

var ConfigProviderKey configProviderKey = "pkg/config_provider"

var ErrMissProvider = errors.New("missing config provider")

type ConfigProvider interface {
	Get(name string) (value string, err error)
	Bind(target interface{}) error
	Engine() interface{}
}

// 从上下文中获取配置
func Get(ctx context.Context) ConfigProvider {
	d := GetDefault(ctx, nil)
	if d == nil {
		panic(ErrMissProvider)
	}
	return d
}

// 从上下文中获取配置，允许设置默认
func GetDefault(ctx context.Context, defaultSource ConfigProvider) ConfigProvider {
	if v, ok := ctx.Value(ConfigProviderKey).(ConfigProvider); ok {
		return v
	}
	return defaultSource
}

// 注入配置到 context
func With(ctx context.Context, ds ConfigProvider) context.Context {
	return context.WithValue(ctx, ConfigProviderKey, ds)
}

func String(ctx context.Context, name string) (string, error) {
	return Get(ctx).Get(name)
}

func Int64(ctx context.Context, name string) (int64, error) {
	v, err := Get(ctx).Get(name)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 64)
}

// 绑定到结构
func Bind(ctx context.Context, val interface{}) error {
	err := Get(ctx).Bind(val)
	if err != nil {
		return err
	}
	return nil
}
