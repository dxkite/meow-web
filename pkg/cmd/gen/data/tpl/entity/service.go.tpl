package {{ .PrivateName }}

import (
	"context"

	"{{ .Pkg }}/pkg/identity"
)

type {{ .Name }}Service interface {
	Create(ctx context.Context, param *Create{{ .Name }}Request) (*{{ .Name }}Dto, error)
	Update(ctx context.Context, param *Update{{ .Name }}Request) (*{{ .Name }}Dto, error)
	Get(ctx context.Context, param *Get{{ .Name }}Request) (*{{ .Name }}Dto, error)
	Delete(ctx context.Context, param *Delete{{ .Name }}Request) error
	List(ctx context.Context, param *List{{ .Name }}Request) (*List{{ .Name }}Response, error)
}

func New{{ .Name }}Service(r {{ .Name }}Repository) {{ .Name }}Service {
	return &{{ .PrivateName }}Service{r: r}
}

type {{ .PrivateName }}Service struct {
	r {{ .Name }}Repository
}


type Create{{ .Name }}Request struct {
	// TODO
}

func (s *{{ .PrivateName }}Service) Create(ctx context.Context, param *Create{{ .Name }}Request) (*{{ .Name }}Dto, error) {
	ent := New{{ .Name }}()

	resp, err := s.r.Create(ctx, ent)
	if err != nil {
		return nil, err
	}

	return New{{ .Name }}Dto(resp), nil
}

type Get{{ .Name }}Request struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *{{ .PrivateName }}Service) Get(ctx context.Context, param *Get{{ .Name }}Request) (*{{ .Name }}Dto, error) {
	ent, err := s.r.Get(ctx, identity.Parse({{ .Name }}Prefix, param.Id))
	if err != nil {
		return nil, err
	}
	obj := New{{ .Name }}Dto(ent)
	return obj, nil
}

type Delete{{ .Name }}Request struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *{{ .PrivateName }}Service) Delete(ctx context.Context, param *Delete{{ .Name }}Request) error {
	err := s.r.Delete(ctx, identity.Parse({{ .Name }}Prefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type List{{ .Name }}Request struct {
	Page         int  `json:"page" form:"page"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
	Expand        []string `json:"expand" form:"expand"`
}

type List{{ .Name }}Response struct {
	Data    []*{{ .Name }}Dto `json:"data"`
	HasMore bool         `json:"has_more"`
	Total   int64        `json:"total,omitempty"`
}

func (s *{{ .PrivateName }}Service) List(ctx context.Context, param *List{{ .Name }}Request) (*List{{ .Name }}Response, error) {
	if param.Page == 0 {
		param.Page = 1
	}

	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listParam := &List{{ .Name }}Param{
		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	}

	listRst, err := s.r.List(ctx, listParam)
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*{{ .Name }}Dto, n)

	for i, v := range listRst.Data {
		items[i] = New{{ .Name }}Dto(v)
	}

	rst := &List{{ .Name }}Response{}
	rst.Data = items
	rst.HasMore = n == param.PerPage
	rst.Total = listRst.Total
	return rst, nil
}

type Update{{ .Name }}Request struct {
	Id string `json:"id" uri:"id" binding:"required"`
	Create{{ .Name }}Request
}

func (s *{{ .PrivateName }}Service) Update(ctx context.Context, param *Update{{ .Name }}Request) (*{{ .Name }}Dto, error) {
	id := identity.Parse({{ .Name }}Prefix, param.Id)
	ent := New{{ .Name }}()

	err := s.r.Update(ctx, id, ent)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &Get{{ .Name }}Request{Id: param.Id})
}
