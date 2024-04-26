package repository

import (
	"context"

	"dxkite.cn/meownest/src/model"
	"gorm.io/gorm"
)

type ServerName interface {
	Create(ctx context.Context, serverName *model.ServerName) (*model.ServerName, error)
}

func NewServerName(db *gorm.DB) ServerName {
	return &serverName{db: db}
}

type serverName struct {
	db *gorm.DB
}

var _ ServerName = &serverName{}

func (s *serverName) Create(ctx context.Context, serverName *model.ServerName) (*model.ServerName, error) {
	if err := s.db.Create(&serverName).Error; err != nil {
		return nil, err
	}
	return serverName, nil
}
