package service

import (
	"context"

	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/utils"
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

func NewRoute(r repository.Route, rl repository.Link, re repository.Endpoint, rc repository.Collection, ra repository.Authorize) Route {
	return &route{r: r, rl: rl, re: re, rc: rc, ra: ra}
}

type route struct {
	r  repository.Route
	rl repository.Link
	re repository.Endpoint
	rc repository.Collection
	ra repository.Authorize
}

type CreateRouteParam struct {
	// 路由名称
	Name string `json:"name" form:"name" binding:"required"`
	// 路由描述
	Description string `json:"description" form:"description"`
	// 支持方法
	Method []string `json:"method" form:"method" binding:"required"`
	// 匹配路径
	Path string `json:"path" form:"path" binding:"required"`
	// 特殊匹配规则
	MatchOptions []*value.MatchOption `json:"match_options" form:"match_options" binding:"dive,required"`
	// 路由分组ID
	CollectionId string `json:"collection_id" form:"collection_id" binding:"required"`
	// 绑定的后端服务
	EndpointId string `json:"endpoint_id" form:"endpoint_id"`
	// 鉴权配置
	AuthorizeId string `json:"authorize_id" form:"authorize_id"`
}

func (s *route) Create(ctx context.Context, param *CreateRouteParam) (*dto.Route, error) {
	var obj *dto.Route
	err := data_source.Transaction(ctx, func(ctx context.Context) error {

		ent, err := s.r.Create(ctx, &entity.Route{
			Name:         param.Name,
			Description:  param.Description,
			Method:       param.Method,
			Path:         param.Path,
			MatchOptions: param.MatchOptions,
			CollectionId: identity.Parse(constant.CollectionPrefix, param.CollectionId),
			AuthorizeId:  identity.Parse(constant.AuthorizePrefix, param.AuthorizeId),
			EndpointId:   identity.Parse(constant.EndpointPrefix, param.EndpointId),
		})

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

	obj := dto.NewRoute(rst)
	if utils.InStringSlice("endpoint", param.Expand) {
		ent, err := s.re.Get(ctx, rst.EndpointId)
		if err != nil {
			return nil, err
		}
		obj.Endpoint = dto.NewEndpoint(ent)
	}

	if utils.InStringSlice("authorize", param.Expand) {
		ent, err := s.ra.Get(ctx, rst.AuthorizeId)
		if err != nil {
			return nil, err
		}
		obj.Authorize = dto.NewAuthorize(ent)
	}
	return obj, nil
}

type ListRouteParam struct {
	Name         string `json:"name" form:"name"`
	Path         string `json:"path" form:"path"`
	CollectionId string `json:"collection_id" form:"collection_id"`

	Page         int  `json:"page" form:"page"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
}

type ListRouteResult struct {
	Data    []*dto.Route `json:"data"`
	HasMore bool         `json:"has_more"`
	Total   int64        `json:"total,omitempty"`
}

func (s *route) List(ctx context.Context, param *ListRouteParam) (*ListRouteResult, error) {
	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listParam := &repository.ListRouteParam{
		Name:         param.Name,
		Path:         param.Path,
		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	}

	if param.CollectionId != "" {
		collId := identity.Parse(constant.CollectionPrefix, param.CollectionId)
		collList, err := s.rc.GetChildren(ctx, collId)
		if err != nil {
			return nil, err
		}

		collIdList := []uint64{collId}

		for _, v := range collList {
			collIdList = append(collIdList, v.Id)
		}

		routeLink, err := s.rl.Linked(ctx, constant.LinkDirectCollectionRoute, collIdList)
		if err != nil {
			return nil, err
		}

		routeIdList := []uint64{}
		for _, v := range routeLink {
			routeIdList = append(routeIdList, v.LinkedId)
		}
		listParam.IdIn = routeIdList
	}

	listRst, err := s.r.List(ctx, listParam)
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*dto.Route, n)

	for i, v := range listRst.Data {
		items[i] = dto.NewRoute(v)
	}

	rst := &ListRouteResult{}
	rst.Data = items
	rst.HasMore = n == param.PerPage
	rst.Total = listRst.Total
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
	err := data_source.Transaction(ctx, func(txCtx context.Context) error {

		entId := identity.Parse(constant.RoutePrefix, param.Id)

		err := s.r.Update(txCtx, entId, &entity.Route{
			Name:         param.Name,
			Description:  param.Description,
			Method:       param.Method,
			Path:         param.Path,
			MatchOptions: param.MatchOptions,
			CollectionId: identity.Parse(constant.CollectionPrefix, param.CollectionId),
			AuthorizeId:  identity.Parse(constant.AuthorizePrefix, param.AuthorizeId),
			EndpointId:   identity.Parse(constant.EndpointPrefix, param.EndpointId),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetRouteParam{Id: param.Id})
}
