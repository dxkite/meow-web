package user

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(ctx context.Context, session *Session) (*Session, error)
	Get(ctx context.Context, id uint64) (*Session, error)
	Update(ctx context.Context, id uint64, ent *Session) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListSessionParam) (*ListSessionResult, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*Session, error)
	SetDeletedByUser(ctx context.Context, userId uint64) error
}

func NewSessionRepository() SessionRepository {
	return new(sessionRepository)
}

type sessionRepository struct {
}

func (r *sessionRepository) Get(ctx context.Context, id uint64) (*Session, error) {
	var item Session
	if err := r.dataSource(ctx).Where("id = ? and deleted = 0", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *sessionRepository) BatchGet(ctx context.Context, ids []uint64) ([]*Session, error) {
	var items []*Session
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
	Data  []*Session
	Total int64
}

func (r *sessionRepository) List(ctx context.Context, param *ListSessionParam) (*ListSessionResult, error) {
	var items []*Session
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
		if err := db.Model(Session{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *sessionRepository) Create(ctx context.Context, session *Session) (*Session, error) {
	if err := r.dataSource(ctx).Create(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *sessionRepository) Update(ctx context.Context, id uint64, ent *Session) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) SetDeletedByUser(ctx context.Context, userId uint64) error {
	if err := r.dataSource(ctx).Where("user_id = ?", userId).Updates(Session{Deleted: 1}).Error; err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(Session{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
