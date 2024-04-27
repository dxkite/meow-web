package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type ServerName interface {
	Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error)
	Get(ctx context.Context, id uint64) (*entity.ServerName, error)
}

func NewServerName(db *gorm.DB) ServerName {
	return &serverName{db: db}
}

type serverName struct {
	db *gorm.DB
}

func (r *serverName) Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error) {
	if err := r.dataSource(ctx).Create(&serverName).Error; err != nil {
		return nil, err
	}
	return serverName, nil
}

func (r *serverName) Get(ctx context.Context, id uint64) (*entity.ServerName, error) {
	var cert entity.ServerName
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *serverName) dataSource(ctx context.Context) *gorm.DB {
	return DataSource(ctx, r.db)
}
