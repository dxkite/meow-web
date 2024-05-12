package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Certificate interface {
	Create(ctx context.Context, certificate *entity.Certificate) (*entity.Certificate, error)
	Get(ctx context.Context, id uint64) (*entity.Certificate, error)
	Update(ctx context.Context, id uint64, ent *entity.Certificate) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListCertificateParam) (*ListCertificateResult, error)
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
	Name string
	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListCertificateResult struct {
	Data  []*entity.Certificate
	Total int64
}

func (r *certificate) List(ctx context.Context, param *ListCertificateParam) (*ListCertificateResult, error) {
	var items []*entity.Certificate
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

	rst := &ListCertificateResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.Certificate{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
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
	return database.Get(ctx).Engine().(*gorm.DB)
}
