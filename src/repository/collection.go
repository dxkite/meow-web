package repository

import (
	"context"
	"strconv"

	"dxkite.cn/meownest/src/model"
	"gorm.io/gorm"
)

type Collection interface {
	Create(ctx context.Context, param *model.Collection) (*model.Collection, error)
}

func NewCollection(db *gorm.DB) Collection {
	return &collection{db: db}
}

type collection struct {
	db *gorm.DB
}

func (r *collection) Create(ctx context.Context, param *model.Collection) (*model.Collection, error) {
	index := "."
	depth := 1

	if param.ParentId != 0 {
		parent, err := r.Get(ctx, param.ParentId)
		if err != nil {
			return nil, err
		}
		index = parent.Index + strconv.FormatUint(parent.Id, 10) + "."
		depth = parent.Depth + 1
	}

	param.Index = index
	param.Depth = depth

	if err := r.db.Create(&param).Error; err != nil {
		return nil, err
	}

	return param, nil
}

func (r *collection) Get(ctx context.Context, id uint64) (*model.Collection, error) {
	var item model.Collection
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
