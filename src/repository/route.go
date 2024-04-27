package repository

import (
	"context"

	"dxkite.cn/meownest/src/model"
	"gorm.io/gorm"
)

type Route interface {
	Create(ctx context.Context, route *model.Route) (*model.Route, error)
	Get(ctx context.Context, id uint64) (*model.Route, error)
}

func NewRoute(db *gorm.DB) Route {
	return &route{db: db}
}

type route struct {
	db *gorm.DB
}

func (s *route) Create(ctx context.Context, route *model.Route) (*model.Route, error) {
	if err := s.db.Create(&route).Error; err != nil {
		return nil, err
	}
	return route, nil
}

func (s *route) Get(ctx context.Context, id uint64) (*model.Route, error) {
	var cert model.Route
	if err := s.db.Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}
