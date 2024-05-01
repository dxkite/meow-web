package datasource

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DataSourceKey = "pkg/data_source"

var ErrMissSource = errors.New("missing data source")

type DataSource interface {
	DB() *gorm.DB
}

type dataSource struct {
	db *gorm.DB
}

func (dataSource *dataSource) DB() *gorm.DB {
	return dataSource.db
}

func New(db *gorm.DB) DataSource {
	return &dataSource{db: db}
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

func RegisterToGin(ds DataSource) gin.HandlerFunc {
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
	return ds.DB().Transaction(func(tx *gorm.DB) error {
		txCtx := With(ctx, New(tx))
		return txFn(txCtx)
	})
}
