package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/valueobject"
)

type CreateEndpointParam struct {
	Name string `json:"name" form:"name" binding:"required"`
	// 服务类型
	Type string `json:"type" binding:"required"`
	// 重写配置
	ForwardRewrite *valueobject.ForwardRewriteOption `json:"forward_rewrite"`
	// 请求头转发配置
	ForwardHeader []*valueobject.ForwardHeaderOption `json:"forward_header"`
	// 匹配规则
	Matcher []*valueobject.MatcherOption `json:"matcher"`
	// 远程服务
	Endpoint *valueobject.ForwardEndpoint `json:"endpoint" binding:"required"`
}

type GetEndpointParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

type Endpoint interface {
	Create(ctx context.Context, param *CreateEndpointParam) (*dto.Endpoint, error)
	Get(ctx context.Context, param *GetEndpointParam) (*dto.Endpoint, error)
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
