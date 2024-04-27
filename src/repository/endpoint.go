package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Endpoint interface {
	Create(ctx context.Context, endpoint *entity.Endpoint) (*entity.Endpoint, error)
	Get(ctx context.Context, id uint64) (*entity.Endpoint, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Endpoint, error)
}

func NewEndpoint(db *gorm.DB) Endpoint {
	return &endpoint{db: db}
}

type endpoint struct {
	db *gorm.DB
}

func (s *endpoint) Create(ctx context.Context, endpoint *entity.Endpoint) (*entity.Endpoint, error) {
	if err := s.db.Create(&endpoint).Error; err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (s *endpoint) Get(ctx context.Context, id uint64) (*entity.Endpoint, error) {
	var cert entity.Endpoint
	if err := s.db.Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (s *endpoint) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Endpoint, error) {
	var items []*entity.Endpoint
	if err := s.db.Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
