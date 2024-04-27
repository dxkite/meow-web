package repository

import (
	"context"

	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Certificate interface {
	Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error)
	Get(ctx context.Context, id uint64) (*entity.Certificate, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Certificate, error)
}

func NewCertificate(db *gorm.DB) Certificate {
	return &certificate{db: db}
}

type certificate struct {
	db *gorm.DB
}

func (r *certificate) Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error) {
	if err := r.dataSource(ctx).Create(&certificate).Error; err != nil {
		return nil, err
	}
	return certificate, nil
}

func (r *certificate) Get(ctx context.Context, id uint64) (*entity.Certificate, error) {
	var cert entity.Certificate
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificate) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Certificate, error) {
	var items []*entity.Certificate
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *certificate) dataSource(ctx context.Context) *gorm.DB {
	return DataSource(ctx, r.db)
}
