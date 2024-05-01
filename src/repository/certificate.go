package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/datasource"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Certificate interface {
	Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error)
	Get(ctx context.Context, id uint64) (*entity.Certificate, error)
	Update(ctx context.Context, id uint64, ent *entity.Certificate) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListCertificateParam) ([]*entity.Certificate, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Certificate, error)
}

func NewCertificate() Certificate {
	return &certificate{}
}

type certificate struct {
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

type ListCertificateParam struct {
	Name          string
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
}

func (r *certificate) List(ctx context.Context, param *ListCertificateParam) ([]*entity.Certificate, error) {
	var items []*entity.Certificate
	db := r.dataSource(ctx).Model(entity.Certificate{})

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

func (r *certificate) Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error) {
	if err := r.dataSource(ctx).Create(&certificate).Error; err != nil {
		return nil, err
	}
	return certificate, nil
}

func (r *certificate) Update(ctx context.Context, id uint64, ent *entity.Certificate) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *certificate) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Certificate{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *certificate) dataSource(ctx context.Context) *gorm.DB {
	return datasource.Get(ctx).DB()
}
