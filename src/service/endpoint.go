package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/value"
)

type CreateEndpointParam struct {
	Name string `json:"name" form:"name" binding:"required"`
	// 服务类型
	Type string `json:"type" binding:"required"`
	// 重写配置
	ForwardRewrite *value.ForwardRewriteOption `json:"forward_rewrite"`
	// 请求头转发配置
	ForwardHeader []*value.ForwardHeaderOption `json:"forward_header" binding:"dive,required"`
	// 匹配规则
	Matcher []*value.MatchOption `json:"matcher" binding:"dive,required"`
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
		Name:           param.Name,
		Type:           param.Type,
		ForwardRewrite: param.ForwardRewrite,
		ForwardHeader:  param.ForwardHeader,
		Matcher:        param.Matcher,
		Endpoint:       param.Endpoint,
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
	Name          string   `form:"name"`
	Limit         int      `form:"limit" binding:"max=1000"`
	StartingAfter string   `form:"starting_after"`
	EndingBefore  string   `form:"ending_before"`
	Expand        []string `json:"expand" form:"expand"`
}

type ListEndpointResult struct {
	HasMore bool            `json:"has_more"`
	Data    []*dto.Endpoint `json:"data"`
}

func (s *endpoint) List(ctx context.Context, param *ListEndpointParam) (*ListEndpointResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListEndpointParam{
		Name:          param.Name,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.EndpointPrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.EndpointPrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.Endpoint, n)

	for i, v := range entities {
		items[i] = dto.NewEndpoint(v)
	}

	rst := &ListEndpointResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

type UpdateEndpointParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateEndpointParam
}

func (s *endpoint) Update(ctx context.Context, param *UpdateEndpointParam) (*dto.Endpoint, error) {
	id := identity.Parse(constant.EndpointPrefix, param.Id)
	err := s.r.Update(ctx, id, &entity.Endpoint{
		Name: param.Name,
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetEndpointParam{Id: param.Id})
}
