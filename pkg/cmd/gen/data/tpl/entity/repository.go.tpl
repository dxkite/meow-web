package repository

import (
	"context"

	"{{ .Pkg }}/pkg/data_source"
	"{{ .Pkg }}/src/entity"
	"gorm.io/gorm"
)

type {{ .Name }} interface {
	Create(ctx context.Context, {{ .PrivateName }} *entity.{{ .Name }}) (*entity.{{ .Name }}, error)
	Get(ctx context.Context, id uint64) (*entity.{{ .Name }}, error)
	Update(ctx context.Context, id uint64, ent *entity.{{ .Name }}) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *List{{ .Name }}Param) (*List{{ .Name }}Result, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.{{ .Name }}, error)
}

func New{{ .Name }}() {{ .Name }} {
	return new({{ .PrivateName }})
}

type {{ .PrivateName }} struct {
}

func (r *{{ .PrivateName }}) Get(ctx context.Context, id uint64) (*entity.{{ .Name }}, error) {
	var item entity.{{ .Name }}
	if err := r.dataSource(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *{{ .PrivateName }}) BatchGet(ctx context.Context, ids []uint64) ([]*entity.{{ .Name }}, error) {
	var items []*entity.{{ .Name }}
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type List{{ .Name }}Param struct {
	// TODO external condition

	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type List{{ .Name }}Result struct {
	Data  []*entity.{{ .Name }}
	Total int64
}


func (r *{{ .PrivateName }}) List(ctx context.Context, param *List{{ .Name }}Param) (*List{{ .Name }}Result, error) {
	var items []*entity.{{ .Name }}
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

	rst := &List{{ .Name }}Result{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.{{ .Name }}{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *{{ .PrivateName }}) Create(ctx context.Context, {{ .PrivateName }} *entity.{{ .Name }}) (*entity.{{ .Name }}, error) {
	if err := r.dataSource(ctx).Create(&{{ .PrivateName }}).Error; err != nil {
		return nil, err
	}
	return {{ .PrivateName }}, nil
}

func (r *{{ .PrivateName }}) Update(ctx context.Context, id uint64, ent *entity.{{ .Name }}) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .PrivateName }}) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.{{ .Name }}{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .PrivateName }}) dataSource(ctx context.Context) *gorm.DB {
	return data_source.Get(ctx).RawSource().(*gorm.DB)
}
