package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Authorize interface {
	Create(ctx context.Context, authorize *entity.Authorize) (*entity.Authorize, error)
	Get(ctx context.Context, id uint64) (*entity.Authorize, error)
	Update(ctx context.Context, id uint64, ent *entity.Authorize) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListAuthorizeParam) ([]*entity.Authorize, error)
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
	Name          string
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
}

func (r *authorize) List(ctx context.Context, param *ListAuthorizeParam) ([]*entity.Authorize, error) {
	var items []*entity.Authorize
	db := r.dataSource(ctx).Model(entity.Authorize{})

	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}

	if param.StartingAfter != 0 {
		db = db.Where("id > ?", param.StartingAfter)
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
	return data_source.Get(ctx).(data_source.GormDataSource).Gorm()
}
