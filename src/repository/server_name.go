package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type ServerName interface {
	Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error)
	Get(ctx context.Context, id uint64) (*entity.ServerName, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.ServerName, error)
	List(ctx context.Context, param *ListServerNameParam) (*ListServerNameResult, error)
	Update(ctx context.Context, id uint64, ent *entity.ServerName) error
	Delete(ctx context.Context, id uint64) error
}

func NewServerName() ServerName {
	return &serverName{}
}

type serverName struct {
}

func (r *serverName) Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error) {
	if err := r.dataSource(ctx).Create(&serverName).Error; err != nil {
		return nil, err
	}
	return serverName, nil
}

func (r *serverName) Get(ctx context.Context, id uint64) (*entity.ServerName, error) {
	var cert entity.ServerName
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *serverName) BatchGet(ctx context.Context, ids []uint64) ([]*entity.ServerName, error) {
	var items []*entity.ServerName
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *serverName) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.ServerName{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverName) Update(ctx context.Context, id uint64, ent *entity.ServerName) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

type ListServerNameParam struct {
	Name string
	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListServerNameResult struct {
	Data  []*entity.ServerName
	Total int64
}

func (r *serverName) List(ctx context.Context, param *ListServerNameParam) (*ListServerNameResult, error) {
	var items []*entity.ServerName
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

	rst := &ListServerNameResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.ServerName{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *serverName) dataSource(ctx context.Context) *gorm.DB {
	return data_source.Get(ctx).RawSource().(*gorm.DB)
}
