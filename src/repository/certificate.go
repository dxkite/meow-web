package repository

import (
	"context"

	"dxkite.cn/meownest/src/model"
	"gorm.io/gorm"
)

type Certificate interface {
	Create(ctx context.Context, certificate *model.Certificate) (*model.Certificate, error)
}

func NewCertificate(db *gorm.DB) Certificate {
	return &certificate{db: db}
}

type certificate struct {
	db *gorm.DB
}

func (s *certificate) Create(ctx context.Context, certificate *model.Certificate) (*model.Certificate, error) {
	if err := s.db.Create(&certificate).Error; err != nil {
		return nil, err
	}
	return certificate, nil
}
