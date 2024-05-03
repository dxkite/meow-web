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

func NewRoute(r repository.Route, rl repository.Link, re repository.Endpoint, rc repository.Collection) Route {
	return &route{r: r, rl: rl, re: re, rc: rc}
}

type route struct {
	r  repository.Route
	rl repository.Link
	re repository.Endpoint
	rc repository.Collection
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
	MatchOptions []*value.MatchOption `json:"match_options" form:"match_options" binding:"dive,required"`
	// 绑定的后端服务
	EndpointId []string `json:"endpoint_id" form:"endpoint_id"`
}

func (s *route) Create(ctx context.Context, param *CreateRouteParam) (*dto.Route, error) {
	var obj *dto.Route
	err := data_source.Transaction(ctx, func(ctx context.Context) error {

		collId := identity.Parse(constant.CollectionPrefix, param.CollectionId)
		ent, err := s.r.Create(ctx, &entity.Route{
			Name:         param.Name,
			Description:  param.Description,
			Method:       param.Method,
			Path:         param.Path,
			MatchOptions: param.MatchOptions,
		})

		if err != nil {
			return err
		}

		err = s.rl.LinkOnce(ctx, constant.LinkDirectCollectionRoute, collId, ent.Id)
		if err != nil {
			return err
		}

		if err := s.batchLinkOnce(ctx, constant.LinkDirectRouteEndpoint, ent.Id, identity.ParseSlice(constant.EndpointPrefix, param.EndpointId)); err != nil {
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

	if utils.InStringSlice("endpoints", param.Expand) {
		entityIds := []uint64{}

		linked, err := s.rl.LinkOf(ctx, constant.LinkDirectRouteEndpoint, []uint64{rst.Id})
		if err != nil {
			return nil, err
		}

		for _, v := range linked {
			entityIds = append(entityIds, v.LinkedId)
		}

		entities, err := s.re.BatchGet(ctx, entityIds)
		if err != nil {
			return nil, err
		}

		items := make([]*dto.Endpoint, len(entities))
		for i, v := range entities {
			items[i] = dto.NewEndpoint(v)
		}

		obj.Endpoints = items
	}

	return obj, nil
}

type ListRouteParam struct {
	Name          string `json:"name" form:"name"`
	Path          string `json:"path" form:"path"`
	Limit         int    `json:"limit" form:"limit" binding:"max=1000"`
	CollectionId  string `json:"collection_id" form:"collection_id"`
	StartingAfter string `json:"starting_after" form:"starting_after"`
	EndingBefore  string `json:"ending_before" form:"ending_before"`
}

type ListRouteResult struct {
	HasMore bool         `json:"has_more"`
	Data    []*dto.Route `json:"data"`
}

func (s *route) List(ctx context.Context, param *ListRouteParam) (*ListRouteResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	listParam := &repository.ListRouteParam{
		Name:          param.Name,
		Path:          param.Path,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.RoutePrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.RoutePrefix, param.EndingBefore),
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

		routeLink, err := s.rl.LinkOf(ctx, constant.LinkDirectCollectionRoute, collIdList)
		if err != nil {
			return nil, err
		}

		routeIdList := []uint64{}
		for _, v := range routeLink {
			routeIdList = append(routeIdList, v.LinkedId)
		}
		listParam.IdIn = routeIdList
	}

	entities, err := s.r.List(ctx, listParam)
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
	err := data_source.Transaction(ctx, func(ctx context.Context) error {

		entId := identity.Parse(constant.RoutePrefix, param.Id)

		err := s.r.Update(ctx, entId, &entity.Route{
			Name:         param.Name,
			Description:  param.Description,
			Method:       param.Method,
			Path:         param.Path,
			MatchOptions: param.MatchOptions,
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

		if err := s.batchLinkOnce(ctx, constant.LinkDirectRouteEndpoint, entId, identity.ParseSlice(constant.EndpointPrefix, param.EndpointId)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetRouteParam{Id: param.Id})
}

func (s *route) batchLinkOnce(ctx context.Context, direct string, id uint64, linkedId []uint64) error {
	return s.batchLink(ctx, direct, id, linkedId, true)
}

func (s *route) batchLink(ctx context.Context, direct string, id uint64, linkedId []uint64, once bool) error {
	if once {
		if err := s.rl.DeleteAllLink(ctx, direct, id); err != nil {
			return err
		}
	}

	if err := s.rl.BatchLink(ctx, direct, id, linkedId); err != nil {
		return err
	}

	return nil
}
