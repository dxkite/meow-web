package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Authorize interface {
	Create(ctx context.Context, authorize *entity.Authorize) (*entity.Authorize, error)
	Get(ctx context.Context, id uint64) (*entity.Authorize, error)
	Update(ctx context.Context, id uint64, ent *entity.Authorize) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListAuthorizeParam) (*ListAuthorizeResult, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Authorize, error)
}

func NewAuthorize() Authorize {
	return new(authorize)
}

type authorize struct {
}

func (r *authorize) Get(ctx context.Context, id uint64) (*entity.Authorize, error) {
	var cert entity.Authorize
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *authorize) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Authorize, error) {
	var items []*entity.Authorize
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type ListAuthorizeParam struct {
	Name string
	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListAuthorizeResult struct {
	Data  []*entity.Authorize
	Total int64
}

func (r *authorize) List(ctx context.Context, param *ListAuthorizeParam) (*ListAuthorizeResult, error) {
	var items []*entity.Authorize
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
		if param.Name != "" {
			db = db.Where("name like ?", "%"+param.Name+"%")
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

	rst := &ListAuthorizeResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.Authorize{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *authorize) Create(ctx context.Context, authorize *entity.Authorize) (*entity.Authorize, error) {
	if err := r.dataSource(ctx).Create(&authorize).Error; err != nil {
		return nil, err
	}
	return authorize, nil
}

func (r *authorize) Update(ctx context.Context, id uint64, ent *entity.Authorize) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *authorize) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Authorize{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *authorize) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
