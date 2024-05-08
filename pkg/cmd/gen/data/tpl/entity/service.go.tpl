package service

import (
	"context"

	"{{ .Pkg }}/pkg/identity"
	"{{ .Pkg }}/src/constant"
	"{{ .Pkg }}/src/dto"
	"{{ .Pkg }}/src/entity"
	"{{ .Pkg }}/src/repository"
)

type {{ .Name }} interface {
	Create(ctx context.Context, param *Create{{ .Name }}Param) (*dto.{{ .Name }}, error)
	Update(ctx context.Context, param *Update{{ .Name }}Param) (*dto.{{ .Name }}, error)
	Get(ctx context.Context, param *Get{{ .Name }}Param) (*dto.{{ .Name }}, error)
	Delete(ctx context.Context, param *Delete{{ .Name }}Param) error
	List(ctx context.Context, param *List{{ .Name }}Param) (*List{{ .Name }}Result, error)
}

func New{{ .Name }}(r repository.{{ .Name }}) {{ .Name }} {
	return &{{ .PrivateName }}{r: r}
}

type {{ .PrivateName }} struct {
	r repository.{{ .Name }}
}


type Create{{ .Name }}Param struct {
	// TODO
}

func (s *{{ .PrivateName }}) Create(ctx context.Context, param *Create{{ .Name }}Param) (*dto.{{ .Name }}, error) {
	ent := entity.New{{ .Name }}()

	resp, err := s.r.Create(ctx, ent)
	if err != nil {
		return nil, err
	}

	return dto.New{{ .Name }}(resp), nil
}

type Get{{ .Name }}Param struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *{{ .PrivateName }}) Get(ctx context.Context, param *Get{{ .Name }}Param) (*dto.{{ .Name }}, error) {
	ent, err := s.r.Get(ctx, identity.Parse(constant.{{ .Name }}Prefix, param.Id))
	if err != nil {
		return nil, err
	}
	obj := dto.New{{ .Name }}(ent)
	return obj, nil
}

type Delete{{ .Name }}Param struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *{{ .PrivateName }}) Delete(ctx context.Context, param *Delete{{ .Name }}Param) error {
	err := s.r.Delete(ctx, identity.Parse(constant.{{ .Name }}Prefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type List{{ .Name }}Param struct {
	Page         int  `json:"page" form:"page" binding:"min=1"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
	Expand        []string `json:"expand" form:"expand"`
}

type List{{ .Name }}Result struct {
	Data    []*dto.{{ .Name }} `json:"data"`
	HasMore bool         `json:"has_more"`
	Total   int64        `json:"total,omitempty"`
}

func (s *{{ .PrivateName }}) List(ctx context.Context, param *List{{ .Name }}Param) (*List{{ .Name }}Result, error) {
	if param.Page == 0 {
		param.Page = 1
	}

	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listParam := &repository.List{{ .Name }}Param{
		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	}

	listRst, err := s.r.List(ctx, listParam)
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*dto.{{ .Name }}, n)

	for i, v := range listRst.Data {
		items[i] = dto.New{{ .Name }}(v)
	}

	rst := &List{{ .Name }}Result{}
	rst.Data = items
	rst.HasMore = n == param.PerPage
	rst.Total = listRst.Total
	return rst, nil
}

type Update{{ .Name }}Param struct {
	Id string `json:"id" uri:"id" binding:"required"`
	Create{{ .Name }}Param
}

func (s *{{ .PrivateName }}) Update(ctx context.Context, param *Update{{ .Name }}Param) (*dto.{{ .Name }}, error) {
	id := identity.Parse(constant.{{ .Name }}Prefix, param.Id)
	ent := entity.New{{ .Name }}()

	err := s.r.Update(ctx, id, ent)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &Get{{ .Name }}Param{Id: param.Id})
}
