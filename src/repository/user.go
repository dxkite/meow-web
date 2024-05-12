package repository

import (
	"context"
	"errors"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

var ErrUserNotExist = errors.New("user not exist")

type User interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, id uint64) (*entity.User, error)
	Update(ctx context.Context, id uint64, ent *entity.User) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListUserParam) (*ListUserResult, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.User, error)
	GetBy(ctx context.Context, param GetUserByParam) (*entity.User, error)
}

func NewUser() User {
	return new(user)
}

type user struct {
}

func (r *user) Get(ctx context.Context, id uint64) (*entity.User, error) {
	var item entity.User
	if err := r.dataSource(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *user) BatchGet(ctx context.Context, ids []uint64) ([]*entity.User, error) {
	var items []*entity.User
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type GetUserByParam struct {
	Name string
}

func (r *user) GetBy(ctx context.Context, param GetUserByParam) (*entity.User, error) {
	var item entity.User
	db := r.dataSource(ctx)
	if param.Name != "" {
		db = db.Where("name = ?", param.Name)
	}
	if err := db.First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

type ListUserParam struct {
	// condition
	Name string
	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListUserResult struct {
	Data  []*entity.User
	Total int64
}

func (r *user) List(ctx context.Context, param *ListUserParam) (*ListUserResult, error) {
	var items []*entity.User
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

	rst := &ListUserResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.User{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *user) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.dataSource(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *user) Update(ctx context.Context, id uint64, ent *entity.User) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *user) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *user) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
