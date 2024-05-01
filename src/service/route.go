package service

import (
	"context"

	"dxkite.cn/meownest/pkg/datasource"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/value"
)

type GetRouteParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

type Route interface {
	Create(ctx context.Context, param *CreateRouteParam) (*dto.Route, error)
	Get(ctx context.Context, param *GetRouteParam) (*dto.Route, error)
	List(ctx context.Context, param *ListRouteParam) (*ListRouteResult, error)
	Update(ctx context.Context, param *UpdateRouteParam) (*dto.Route, error)
	Delete(ctx context.Context, param *DeleteRouteParam) error
}

func NewRoute(r repository.Route, rl repository.Link) Route {
	return &route{r: r, rl: rl}
}

type route struct {
	r  repository.Route
	rl repository.Link
}

type CreateRouteParam struct {
	// 路由名称
	Name string `json:"name" form:"name" binding:"required"`
	// 路由描述
	Description string `json:"description" form:"description"`
	// 路由分组ID
	CollectionId string `json:"collection_id" form:"collection_id" binding:"required"`
	// 支持方法
	Method []string `json:"method" form:"method" binding:"required"`
	// 匹配路径
	Path string `json:"path" form:"path" binding:"required"`
	// 特殊匹配规则
	Matcher []*value.MatcherOption `json:"matcher" form:"matcher" binding:"dive,required"`
}

func (s *route) Create(ctx context.Context, param *CreateRouteParam) (*dto.Route, error) {
	var obj *dto.Route
	err := datasource.Transaction(ctx, func(ctx context.Context) error {

		collId := identity.Parse(constant.CollectionPrefix, param.CollectionId)
		ent, err := s.r.Create(ctx, &entity.Route{
			Name:        param.Name,
			Description: param.Description,
			Method:      param.Method,
			Path:        param.Path,
			Matcher:     param.Matcher,
		})

		if err != nil {
			return err
		}

		err = s.rl.LinkOnce(ctx, constant.LinkDirectCollectionRoute, collId, ent.Id)
		if err != nil {
			return err
		}

		obj = dto.NewRoute(ent)
		return nil
	})
	return obj, err
}

func (s *route) Get(ctx context.Context, param *GetRouteParam) (*dto.Route, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.RoutePrefix, param.Id))
	if err != nil {
		return nil, err
	}
	return dto.NewRoute(rst), nil
}

type ListRouteParam struct {
	Name          string `form:"name"`
	Path          string `form:"path"`
	Limit         int    `form:"limit" binding:"max=1000"`
	StartingAfter string `form:"starting_after"`
	EndingBefore  string `form:"ending_before"`
}

type ListRouteResult struct {
	HasMore bool         `json:"has_more"`
	Data    []*dto.Route `json:"data"`
}

func (s *route) List(ctx context.Context, param *ListRouteParam) (*ListRouteResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListRouteParam{
		Name:          param.Name,
		Path:          param.Path,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.CollectionPrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.CollectionPrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.Route, n)

	for i, v := range entities {
		items[i] = dto.NewRoute(v)
	}

	rst := &ListRouteResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

type DeleteRouteParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *route) Delete(ctx context.Context, param *DeleteRouteParam) error {
	err := s.r.Delete(ctx, identity.Parse(constant.RoutePrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type UpdateRouteParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateRouteParam
}

func (s *route) Update(ctx context.Context, param *UpdateRouteParam) (*dto.Route, error) {
	err := datasource.Transaction(ctx, func(ctx context.Context) error {

		entId := identity.Parse(constant.RoutePrefix, param.Id)

		err := s.r.Update(ctx, entId, &entity.Route{
			Name:        param.Name,
			Description: param.Description,
			Method:      param.Method,
			Path:        param.Path,
			Matcher:     param.Matcher,
		})
		if err != nil {
			return err
		}

		if param.CollectionId != "" {
			collId := identity.Parse(constant.CollectionPrefix, param.CollectionId)
			err = s.rl.LinkOnce(ctx, constant.LinkDirectCollectionRoute, collId, entId)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	obj, err := s.r.Get(ctx, identity.Parse(constant.RoutePrefix, param.Id))
	if err != nil {
		return nil, err
	}
	return dto.NewRoute(obj), nil
}
