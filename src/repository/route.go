package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Route interface {
	Create(ctx context.Context, route *entity.Route) (*entity.Route, error)
	Get(ctx context.Context, id uint64) (*entity.Route, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Route, error)
	List(ctx context.Context, param *ListRouteParam) ([]*entity.Route, error)
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, id uint64, ent *entity.Route) error
}

func NewRoute() Route {
	return &route{}
}

type route struct {
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

func (r *route) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Route{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *route) Update(ctx context.Context, id uint64, ent *entity.Route) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

type ListRouteParam struct {
	Name          string
	Path          string
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
	IdIn          []uint64
}

func (r *route) List(ctx context.Context, param *ListRouteParam) ([]*entity.Route, error) {
	var items []*entity.Route
	db := r.dataSource(ctx).Model(entity.Route{})

	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}

	if param.Path != "" {
		db = db.Where("path like ?", "%"+param.Path+"%")
	}

	if param.StartingAfter != 0 {
		db = db.Where("id > ?", param.StartingAfter)
	}

	if len(param.IdIn) > 0 {
		db = db.Where("id in ?", param.IdIn)
	}

	if param.EndingBefore != 0 {
		db = db.Where("id < ?", param.EndingBefore)
	}

	if param.Limit != 0 {
		db = db.Limit(param.Limit)
	}

	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *route) dataSource(ctx context.Context) *gorm.DB {
	return data_source.Get(ctx).DB()
}
