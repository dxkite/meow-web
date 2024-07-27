package monitor

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"gorm.io/gorm"
)

type DynamicStatRepository interface {
	Create(ctx context.Context, dynamicstat *DynamicStat) (*DynamicStat, error)
	List(ctx context.Context, param *ListDynamicStatParam) ([]*DynamicStat, error)
	DeleteBefore(ctx context.Context, timeBefore uint64) error
}

func NewDynamicStatRepository() DynamicStatRepository {
	return new(dynamicStatRepository)
}

type dynamicStatRepository struct {
}

type ListDynamicStatParam struct {
	StartTime uint64
	EndTime   uint64
	Limit     int
}

func (r *dynamicStatRepository) List(ctx context.Context, param *ListDynamicStatParam) ([]*DynamicStat, error) {
	var items []*DynamicStat
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

func (r *dynamicStatRepository) Create(ctx context.Context, dynamicstat *DynamicStat) (*DynamicStat, error) {
	if err := r.dataSource(ctx).Create(&dynamicstat).Error; err != nil {
		return nil, err
	}
	return dynamicstat, nil
}

func (r *dynamicStatRepository) DeleteBefore(ctx context.Context, timeBefore uint64) error {
	if err := r.dataSource(ctx).Where("time < ?", timeBefore).Delete(DynamicStat{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *dynamicStatRepository) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
