package repository

import (
	"context"
	"strconv"

	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Collection interface {
	Create(ctx context.Context, collection *entity.Collection) (*entity.Collection, error)
	Get(ctx context.Context, id uint64) (*entity.Collection, error)
	Update(ctx context.Context, id uint64, ent *entity.Collection) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListCollectionParam) (*ListCollectionResult, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Collection, error)
	GetChildren(ctx context.Context, id uint64) ([]*entity.Collection, error)
}

func NewCollection() Collection {
	return new(collection)
}

type collection struct {
}

func (r *collection) Create(ctx context.Context, param *entity.Collection) (*entity.Collection, error) {
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

	if err := r.dataSource(ctx).Create(&param).Error; err != nil {
		return nil, err
	}

	return param, nil
}

func (r *collection) Get(ctx context.Context, id uint64) (*entity.Collection, error) {
	var item entity.Collection
	if err := r.dataSource(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *collection) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Collection, error) {
	var items []*entity.Collection
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// 获取所有子级，包括子集的子集
func (r *collection) GetChildren(ctx context.Context, id uint64) ([]*entity.Collection, error) {
	var items []*entity.Collection

	item, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	index := item.Index + strconv.FormatUint(item.Id, 10) + "."
	if err := r.dataSource(ctx).Where("`index` like CONCAT(?, '%')", index).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

type ListCollectionParam struct {
	ParentId uint64
	// depth = 0 只获取当前层级
	// depth > 0 获取当前层级 > depth
	Depth int
	Name  string

	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListCollectionResult struct {
	Data  []*entity.Collection
	Total int64
}

func (r *collection) List(ctx context.Context, param *ListCollectionParam) (*ListCollectionResult, error) {
	var items []*entity.Collection
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
		if param.Name != "" {
			db = db.Where("name like ?", "%"+param.Name+"%")
		}

		deep := param.Depth
		if param.ParentId != 0 {
			db = db.Where("parent_id = ?", param.ParentId)
			if param.Depth != 0 {
				parent, err := r.Get(ctx, param.ParentId)
				if err != nil {
					db.AddError(err)
					return db
				}
				deep = parent.Depth + deep
				db = db.Where("`index` like ?", parent.Index+"%")
			}
		}

		if deep != 0 {
			db = db.Where("depth <= ?", deep)
		}

		return db
	}

	// pagination
	query := db.Scopes(condition)
	if param.Page > 0 && param.PerPage > 0 {
		query.Offset((param.Page - 1) * param.PerPage).Limit(param.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	rst := &ListCollectionResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.Collection{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *collection) Update(ctx context.Context, id uint64, ent *entity.Collection) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *collection) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Collection{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *collection) dataSource(ctx context.Context) *gorm.DB {
	return data_source.Get(ctx).RawSource().(*gorm.DB)
}
