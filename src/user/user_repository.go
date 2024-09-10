package user

import (
	"context"
	"errors"

	"dxkite.cn/meow-web/pkg/utils"
	"gorm.io/gorm"
)

var ErrUserNotExist = errors.New("user not exist")

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	Get(ctx context.Context, id uint64) (*User, error)
	Update(ctx context.Context, id uint64, ent *User) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *ListUserParam) (*ListUserResult, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*User, error)
	GetBy(ctx context.Context, param GetUserByParam) (*User, error)
}

func NewUserRepository() UserRepository {
	return new(userRepository)
}

type userRepository struct {
}

func (r *userRepository) Get(ctx context.Context, id uint64) (*User, error) {
	var item User
	if err := utils.DB(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, r.wrap(err)
	}
	return &item, nil
}

func (r *userRepository) wrap(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotExist
	}
	return err
}

func (r *userRepository) BatchGet(ctx context.Context, ids []uint64) ([]*User, error) {
	var items []*User
	if err := utils.DB(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type GetUserByParam struct {
	Name string
}

func (r *userRepository) GetBy(ctx context.Context, param GetUserByParam) (*User, error) {
	var item User
	db := utils.DB(ctx)
	if param.Name != "" {
		db = db.Where("name = ?", param.Name)
	}
	if err := db.First(&item).Error; err != nil {
		return nil, r.wrap(err)
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
	Data  []*User
	Total int64
}

func (r *userRepository) List(ctx context.Context, param *ListUserParam) (*ListUserResult, error) {
	var items []*User
	db := utils.DB(ctx)

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
		if err := db.Model(User{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *userRepository) Create(ctx context.Context, userRepository *User) (*User, error) {
	if err := utils.DB(ctx).Create(&userRepository).Error; err != nil {
		return nil, err
	}
	return userRepository, nil
}

func (r *userRepository) Update(ctx context.Context, id uint64, ent *User) error {
	if err := utils.DB(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	if err := utils.DB(ctx).Where("id = ?", id).Delete(User{}).Error; err != nil {
		return err
	}
	return nil
}
