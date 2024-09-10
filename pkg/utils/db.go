package utils

import (
	"context"

	"dxkite.cn/nebula/pkg/database"
	"gorm.io/gorm"
)

func DB(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
