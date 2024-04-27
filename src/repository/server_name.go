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

func (s *serverName) Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error) {
	if err := s.db.Create(&serverName).Error; err != nil {
		return nil, err
	}
	return serverName, nil
}

func (s *serverName) Get(ctx context.Context, id uint64) (*entity.ServerName, error) {
	var cert entity.ServerName
	if err := s.db.Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}
