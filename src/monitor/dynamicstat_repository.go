package monitor

import (
	"context"

	"dxkite.cn/meow-web/pkg/utils"
	"gorm.io/gorm"
)

type DynamicStatRepository interface {
	Create(ctx context.Context, dynamicStat *DynamicStat) (*DynamicStat, error)
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
	db := utils.DB(ctx)

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

func (r *dynamicStatRepository) Create(ctx context.Context, dynamicStat *DynamicStat) (*DynamicStat, error) {
	if err := utils.DB(ctx).Create(&dynamicStat).Error; err != nil {
		return nil, err
	}
	return dynamicStat, nil
}

func (r *dynamicStatRepository) DeleteBefore(ctx context.Context, timeBefore uint64) error {
	if err := utils.DB(ctx).Where("time < ?", timeBefore).Delete(DynamicStat{}).Error; err != nil {
		return err
	}
	return nil
}
