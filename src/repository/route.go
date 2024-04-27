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

func (r *route) Create(ctx context.Context, route *entity.Route) (*entity.Route, error) {
	if err := r.dataSource(ctx).Create(&route).Error; err != nil {
		return nil, err
	}
	return route, nil
}

func (r *route) Get(ctx context.Context, id uint64) (*entity.Route, error) {
	var cert entity.Route
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *route) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Route, error) {
	var items []*entity.Route
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *route) dataSource(ctx context.Context) *gorm.DB {
	return DataSource(ctx, r.db)
}
