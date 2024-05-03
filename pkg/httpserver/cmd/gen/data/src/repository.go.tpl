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
	List(ctx context.Context, param *List{{ .Name }}Param) ([]*entity.{{ .Name }}, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.{{ .Name }}, error)
}

func New{{ .Name }}() {{ .Name }} {
	return new({{ .PrivateName }})
}

type {{ .PrivateName }} struct {
}

func (r *{{ .PrivateName }}) Get(ctx context.Context, id uint64) (*entity.{{ .Name }}, error) {
	var cert entity.{{ .Name }}
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *{{ .PrivateName }}) BatchGet(ctx context.Context, ids []uint64) ([]*entity.{{ .Name }}, error) {
	var items []*entity.{{ .Name }}
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type List{{ .Name }}Param struct {
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
}

func (r *{{ .PrivateName }}) List(ctx context.Context, param *List{{ .Name }}Param) ([]*entity.{{ .Name }}, error) {
	var items []*entity.{{ .Name }}
	db := r.dataSource(ctx).Model(entity.{{ .Name }}{})

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
	return data_source.Get(ctx).(data_source.GormDataSource).Gorm()
}
