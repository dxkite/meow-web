package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Certificate interface {
	Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error)
	Get(ctx context.Context, id uint64) (*entity.Certificate, error)
}

func NewCertificate(db *gorm.DB) Certificate {
	return &certificate{db: db}
}

type certificate struct {
	db *gorm.DB
}

func (s *certificate) Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error) {
	if err := s.db.Create(&certificate).Error; err != nil {
		return nil, err
	}
	return certificate, nil
}

func (s *certificate) Get(ctx context.Context, id uint64) (*entity.Certificate, error) {
	var cert entity.Certificate
	if err := s.db.Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}
