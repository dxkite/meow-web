package repository

import (
	"context"

	"gorm.io/gorm"
)

type dataSourceKey string

var DataSourceKey dataSourceKey = "repository/data_source"

func WithDataSource(ctx context.Context, source *gorm.DB) context.Context {
	ctx = context.WithValue(ctx, DataSourceKey, source)
	return ctx
}

func DataSource(ctx context.Context, defaultSource *gorm.DB) *gorm.DB {
	if v, ok := ctx.Value(DataSourceKey).(*gorm.DB); ok {
		return v
	}
	return defaultSource
}
