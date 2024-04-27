package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Route interface {
	Create(ctx context.Context, route *entity.Route) (*entity.Route, error)
	Get(ctx context.Context, id uint64) (*entity.Route, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Route, error)
}

func NewRoute(db *gorm.DB) Route {
	return &route{db: db}
}

type route struct {
	db *gorm.DB
}

func (s *route) Create(ctx context.Context, route *entity.Route) (*entity.Route, error) {
	if err := s.db.Create(&route).Error; err != nil {
		return nil, err
	}
	return route, nil
}

func (s *route) Get(ctx context.Context, id uint64) (*entity.Route, error) {
	var cert entity.Route
	if err := s.db.Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (s *route) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Route, error) {
	var items []*entity.Route
	if err := s.db.Where("id in ?", ids).First(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
