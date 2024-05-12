package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/enum"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/value"
)

type CreateEndpointParam struct {
	// 服务备注名
	Name string `json:"name" form:"name" binding:"required"`
	// 服务类型
	Type enum.EndpointType `json:"type" binding:"required"`
	// 远程服务
	Endpoint *value.ForwardEndpoint `json:"endpoint" binding:"required"`
}

type GetEndpointParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

type Endpoint interface {
	Create(ctx context.Context, param *CreateEndpointParam) (*dto.Endpoint, error)
	Get(ctx context.Context, param *GetEndpointParam) (*dto.Endpoint, error)
	Delete(ctx context.Context, param *DeleteEndpointParam) error
	List(ctx context.Context, param *ListEndpointParam) (*ListEndpointResult, error)
	Update(ctx context.Context, param *UpdateEndpointParam) (*dto.Endpoint, error)
}

func NewEndpoint(r repository.Endpoint) Endpoint {
	return &endpoint{r: r}
}

type endpoint struct {
	r repository.Endpoint
}

func (s *endpoint) Create(ctx context.Context, param *CreateEndpointParam) (*dto.Endpoint, error) {
	rst, err := s.r.Create(ctx, &entity.Endpoint{
		Name:     param.Name,
		Type:     param.Type,
		Endpoint: param.Endpoint,
	})
	if err != nil {
		return nil, err
	}
	return dto.NewEndpoint(rst), nil
}

func (s *endpoint) Get(ctx context.Context, param *GetEndpointParam) (*dto.Endpoint, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.EndpointPrefix, param.Id))
	if err != nil {
		return nil, err
	}
	return dto.NewEndpoint(rst), nil
}

type DeleteEndpointParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *endpoint) Delete(ctx context.Context, param *DeleteEndpointParam) error {
	err := s.r.Delete(ctx, identity.Parse(constant.EndpointPrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type ListEndpointParam struct {
	Name string `form:"name"`

	// pagination
	Page         int  `json:"page" form:"page"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
}

type ListEndpointResult struct {
	Data  []*dto.Endpoint `json:"data"`
	Total int64           `json:"total,omitempty"`
}

func (s *endpoint) List(ctx context.Context, param *ListEndpointParam) (*ListEndpointResult, error) {
	if param.Page == 0 {
		param.Page = 1
	}

	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listRst, err := s.r.List(ctx, &repository.ListEndpointParam{
		Name:         param.Name,
		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	})
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*dto.Endpoint, n)

	for i, v := range listRst.Data {
		items[i] = dto.NewEndpoint(v)
	}

	rst := &ListEndpointResult{}
	rst.Data = items
	rst.Total = listRst.Total
	return rst, nil
}

type UpdateEndpointParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateEndpointParam
}

func (s *endpoint) Update(ctx context.Context, param *UpdateEndpointParam) (*dto.Endpoint, error) {
	id := identity.Parse(constant.EndpointPrefix, param.Id)
	err := s.r.Update(ctx, id, &entity.Endpoint{
		Name:     param.Name,
		Type:     param.Type,
		Endpoint: param.Endpoint,
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetEndpointParam{Id: param.Id})
}
