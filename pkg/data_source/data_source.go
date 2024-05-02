package data_source

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DataSourceKey = "pkg/data_source"

var ErrMissSource = errors.New("missing data source")

type DataSource interface {
	RawSource() interface{}
	Transaction(func(s DataSource) error) error
}

type GormDataSource interface {
	Gorm() *gorm.DB
}

func Get(ctx context.Context) DataSource {
	d := GetDefault(ctx, nil)
	if d == nil {
		panic(ErrMissSource)
	}
	return d
}

func GetDefault(ctx context.Context, defaultSource DataSource) DataSource {
	if v, ok := ctx.Value(DataSourceKey).(DataSource); ok {
		return v
	}
	return defaultSource
}

func GinDataSource(ds DataSource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(DataSourceKey, ds)
	}
}

func With(ctx context.Context, ds DataSource) context.Context {
	if v, ok := ctx.(*gin.Context); ok {
		v.Set(DataSourceKey, ds)
		return v
	}
	return context.WithValue(ctx, DataSourceKey, ds)
}

func Transaction(ctx context.Context, txFn func(txCtx context.Context) error) error {
	ds := Get(ctx)
	return ds.Transaction(func(s DataSource) error {
		txCtx := With(ctx, s)
		return txFn(txCtx)
	})
}
