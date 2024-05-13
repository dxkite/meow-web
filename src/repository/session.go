package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Session interface {
	Create(ctx context.Context, session *entity.Session) (*entity.Session, error)
	Get(ctx context.Context, id uint64) (*entity.Session, error)
	Update(ctx context.Context, id uint64, ent *entity.Session) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListSessionParam) (*ListSessionResult, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Session, error)
}

func NewSession() Session {
	return new(session)
}

type session struct {
}

func (r *session) Get(ctx context.Context, id uint64) (*entity.Session, error) {
	var item entity.Session
	if err := r.dataSource(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *session) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Session, error) {
	var items []*entity.Session
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type ListSessionParam struct {
	// TODO external condition

	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListSessionResult struct {
	Data  []*entity.Session
	Total int64
}

func (r *session) List(ctx context.Context, param *ListSessionParam) (*ListSessionResult, error) {
	var items []*entity.Session
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
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

	rst := &ListSessionResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.Session{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *session) Create(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	if err := r.dataSource(ctx).Create(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *session) Update(ctx context.Context, id uint64, ent *entity.Session) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *session) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Session{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *session) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
