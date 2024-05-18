package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Monitor interface {
	SaveDynamicStat(ctx context.Context, ent *entity.DynamicStat) (*entity.DynamicStat, error)
	ListDynamicStat(ctx context.Context, param *ListDynamicStatParam) ([]*entity.DynamicStat, error)
	DeleteBefore(ctx context.Context, timeBefore uint64) error
}

func NewMonitor() Monitor {
	m := new(monitor)
	return m
}

type monitor struct {
}

func (r *monitor) SaveDynamicStat(ctx context.Context, ent *entity.DynamicStat) (*entity.DynamicStat, error) {
	if err := r.dataSource(ctx).Create(&ent).Error; err != nil {
		return nil, err
	}
	return ent, nil
}

type ListDynamicStatParam struct {
	StartTime uint64
	EndTime   uint64
	Limit     int
}

func (r *monitor) ListDynamicStat(ctx context.Context, param *ListDynamicStatParam) ([]*entity.DynamicStat, error) {
	var items []*entity.DynamicStat
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
		if param.StartTime > 0 {
			db.Where("time >= ?", param.StartTime)
		}
		if param.EndTime > 0 {
			db.Where("time <= ?", param.EndTime)
		}
		if param.Limit > 0 {
			db.Limit(param.Limit)
		}
		return db
	}

	// pagination
	query := db.Scopes(condition).Order("time DESC")
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *monitor) DeleteBefore(ctx context.Context, timeBefore uint64) error {
	if err := r.dataSource(ctx).Where("time < ?", timeBefore).Delete(entity.DynamicStat{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *monitor) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
