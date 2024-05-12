package database

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

var DataSourceKey = "pkg/data_source"

var ErrMissSource = errors.New("missing data source")

// 数据源
type DataSource interface {
	// 数据源引擎
	Engine() interface{}
	// 开启事务
	Transaction(func(s DataSource) error) error
}

// 从上下文中获取数据源
func Get(ctx context.Context) DataSource {
	d := GetDefault(ctx, nil)
	if d == nil {
		panic(ErrMissSource)
	}
	return d
}

// 从上下文中获取数据源
// 获取数据源失败则使用默认数据源
func GetDefault(ctx context.Context, defaultSource DataSource) DataSource {
	if v, ok := ctx.Value(DataSourceKey).(DataSource); ok {
		return v
	}
	return defaultSource
}

// 注入数据源到 context
func With(ctx context.Context, ds DataSource) context.Context {
	return context.WithValue(ctx, DataSourceKey, ds)
}

func GinDataSource(ds DataSource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(DataSourceKey, ds)
	}
}

// 开启事务
func Transaction(ctx context.Context, txFn func(txCtx context.Context) error) error {
	ds := Get(ctx)
	return ds.Transaction(func(s DataSource) error {
		txCtx := With(ctx, s)
		return txFn(txCtx)
	})
}
