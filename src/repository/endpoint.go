package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Endpoint interface {
	Create(ctx context.Context, endpoint *entity.Endpoint) (*entity.Endpoint, error)
	Get(ctx context.Context, id uint64) (*entity.Endpoint, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Endpoint, error)
	List(ctx context.Context, param *ListEndpointParam) ([]*entity.Endpoint, error)
	Update(ctx context.Context, id uint64, ent *entity.Endpoint) error
	Delete(ctx context.Context, id uint64) error
}

func NewEndpoint(db *gorm.DB) Endpoint {
	return &endpoint{db: db}
}

type endpoint struct {
	db *gorm.DB
}

func (r *endpoint) Create(ctx context.Context, endpoint *entity.Endpoint) (*entity.Endpoint, error) {
	if err := r.dataSource(ctx).Create(&endpoint).Error; err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (r *endpoint) Get(ctx context.Context, id uint64) (*entity.Endpoint, error) {
	var cert entity.Endpoint
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *endpoint) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Endpoint, error) {
	var items []*entity.Endpoint
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type ListEndpointParam struct {
	Name          string
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
}

func (r *endpoint) List(ctx context.Context, param *ListEndpointParam) ([]*entity.Endpoint, error) {
	var items []*entity.Endpoint
	db := r.dataSource(ctx).Model(entity.Endpoint{})

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

func (r *endpoint) Update(ctx context.Context, id uint64, ent *entity.Endpoint) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *endpoint) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Endpoint{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *endpoint) dataSource(ctx context.Context) *gorm.DB {
	return DataSource(ctx, r.db)
}
