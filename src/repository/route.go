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
	List(ctx context.Context, param *ListRouteParam) (*ListRouteResult, error)
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, id uint64, ent *entity.Route) error
	Batch(ctx context.Context, batchFn func(item *entity.Route) error) error
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
	Name string
	Path string
	IdIn []uint64

	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListRouteResult struct {
	Data  []*entity.Route
	Total int64
}

func (r *route) List(ctx context.Context, param *ListRouteParam) (*ListRouteResult, error) {
	var items []*entity.Route
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
		if param.Name != "" {
			db.Where("name like ?", "%"+param.Name+"%")
		}

		if param.Path != "" {
			db.Where("path like ?", "%"+param.Path+"%")
		}

		if len(param.IdIn) > 0 {
			db.Where("id in ?", param.IdIn)
		}
		return db
	}

	query := db.Scopes(condition)

	// pagination
	if param.Page > 0 {
		query.Offset((param.Page - 1) * param.PerPage)
	}

	if param.PerPage != 0 {
		query.Limit(param.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	rst := &ListRouteResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.Route{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *route) Batch(ctx context.Context, batchFn func(item *entity.Route) error) error {
	var items []*entity.Route
	if err := r.dataSource(ctx).FindInBatches(&items, 100, func(tx *gorm.DB, batch int) error {
		for i := range items {
			if err := batchFn(items[i]); err != nil {
				return err
			}
		}
		return nil
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *route) dataSource(ctx context.Context) *gorm.DB {
	return data_source.Get(ctx).RawSource().(*gorm.DB)
}
