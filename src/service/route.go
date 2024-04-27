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

type CreateRouteParam struct {
	Name        string                       `json:"name" form:"name" binding:"required"`
	Description string                       `json:"description" form:"description"`
	Method      []string                     `json:"method" form:"method" binding:"required"`
	Path        string                       `json:"path" form:"path" binding:"required"`
	Matcher     []*valueobject.MatcherOption `json:"matcher" form:"matcher"`
}

type GetRouteParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

type Route interface {
	Create(ctx context.Context, param *CreateRouteParam) (*dto.Route, error)
	Get(ctx context.Context, param *GetRouteParam) (*dto.Route, error)
}

func NewRoute(r repository.Route) Route {
	return &route{r: r}
}

type route struct {
	r repository.Route
}

func (s *route) Create(ctx context.Context, param *CreateRouteParam) (*dto.Route, error) {
	rst, err := s.r.Create(ctx, &entity.Route{
		Name:        param.Name,
		Description: param.Description,
		Method:      param.Method,
		Path:        param.Path,
		Matcher:     param.Matcher,
	})
	if err != nil {
		return nil, err
	}
	return dto.NewRoute(rst), nil
}

func (s *route) Get(ctx context.Context, param *GetRouteParam) (*dto.Route, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.RoutePrefix, param.Id))
	if err != nil {
		return nil, err
	}
	return dto.NewRoute(rst), nil
}
