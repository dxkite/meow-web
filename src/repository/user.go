package repository

import (
	"context"
	"errors"

	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

var ErrUserNotExist = errors.New("user not exist")

type User interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, id uint64) (*entity.User, error)
	Update(ctx context.Context, id uint64, ent *entity.User) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListUserParam) ([]*entity.User, error)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}
		return nil, err
	}
	return &item, nil
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

func (r *user) BatchGet(ctx context.Context, ids []uint64) ([]*entity.User, error) {
	var items []*entity.User
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type ListUserParam struct {
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
}

func (r *user) List(ctx context.Context, param *ListUserParam) ([]*entity.User, error) {
	var items []*entity.User
	db := r.dataSource(ctx).Model(entity.User{})

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
	return data_source.Get(ctx).RawSource().(*gorm.DB)
}
