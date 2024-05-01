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
)

type Collection interface {
	Create(ctx context.Context, param *CreateCollectionParam) (*dto.Collection, error)
	Get(ctx context.Context, param *GetCollectionParam) (*dto.Collection, error)
	List(ctx context.Context, param *ListCollectionParam) (*ListCollectionResult, error)
	LinkRoute(ctx context.Context, param *LinkCollectionRouteParam) error
	DeleteRoute(ctx context.Context, param *DeleteCollectionRouteParam) error
	LinkEndpoint(ctx context.Context, param *LinkCollectionEndpointParam) error
	DeleteEndpoint(ctx context.Context, param *DeleteCollectionEndpointParam) error
}

func NewCollection(r repository.Collection, rl repository.Link, rr repository.Route, re repository.Endpoint) Collection {
	return &collection{r: r, rr: rr, rl: rl, re: re}
}

type collection struct {
	r  repository.Collection
	rl repository.Link
	rr repository.Route
	re repository.Endpoint
}

type CreateCollectionParam struct {
	// 分组名
	Name string `json:"name" form:"name" binding:"required"`
	// 分组描述
	Description string `json:"description" form:"description"`
	// 父级节点
	ParentId string `json:"parent_id" form:"parent_id"`
	// 绑定的域名
	ServerNameId []string `json:"server_name_id" form:"server_name"`
	// 绑定的后端服务
	EndpointId []string `json:"endpoint_id" form:"endpoint_id"`
}

func (s *collection) Create(ctx context.Context, param *CreateCollectionParam) (*dto.Collection, error) {
	var obj *dto.Collection

	data_source.Transaction(ctx, func(txCtx context.Context) error {
		item, err := s.r.Create(ctx, &entity.Collection{
			Name:        param.Name,
			Description: param.Description,
			ParentId:    identity.Parse(constant.CollectionPrefix, param.ParentId),
		})

		if err != nil {
			return err
		}

		if err := s.batchLink(ctx, constant.LinkDirectCollectionServerName, item.Id, identity.ParseSlice(constant.ServerNamePrefix, param.ServerNameId)); err != nil {
			return err
		}

		if err := s.batchLink(ctx, constant.LinkDirectCollectionEndpoint, item.Id, identity.ParseSlice(constant.EndpointPrefix, param.EndpointId)); err != nil {
			return err
		}

		obj = dto.NewCollection(item)
		return nil
	})

	return obj, nil
}

func (s *collection) batchLink(ctx context.Context, direct string, id uint64, idArray []uint64) error {
	linkIds := []uint64{}
	routes, err := s.rr.BatchGet(ctx, idArray)

	if err != nil {
		return err
	}

	for _, v := range routes {
		linkIds = append(linkIds, v.Id)
	}

	if err := s.rl.DeleteAllLink(ctx, direct, id); err != nil {
		return err
	}

	if err := s.rl.BatchLink(ctx, direct, id, linkIds); err != nil {
		return err
	}

	return nil
}

type GetCollectionParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *collection) Get(ctx context.Context, param *GetCollectionParam) (*dto.Collection, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return nil, err
	}

	collection := dto.NewCollection(rst)

	if utils.InStringSlice("routes", param.Expand) {
		entityIds := []uint64{}

		linked, err := s.rl.LinkOf(ctx, constant.LinkDirectCollectionRoute, rst.Id)
		if err != nil {
			return nil, err
		}

		for _, v := range linked {
			entityIds = append(entityIds, v.LinkedId)
		}

		entities, err := s.rr.BatchGet(ctx, entityIds)
		if err != nil {
			return nil, err
		}

		items := make([]*dto.Route, len(entities))
		for i, v := range entities {
			items[i] = dto.NewRoute(v)
		}

		collection.Routes = items
	}

	if utils.InStringSlice("endpoints", param.Expand) {
		entityIds := []uint64{}

		linked, err := s.rl.LinkOf(ctx, constant.LinkDirectCollectionEndpoint, rst.Id)
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

		collection.Endpoints = items
	}

	return collection, nil
}

type LinkCollectionRouteParam struct {
	Id      string   `json:"id" uri:"id" binding:"required"`
	RouteId []string `json:"route_id" form:"route_id" binding:"required"`
}

func (s *collection) LinkRoute(ctx context.Context, param *LinkCollectionRouteParam) error {
	item, err := s.r.Get(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return err
	}

	linkIds := []uint64{}
	routes, err := s.rr.BatchGet(ctx, identity.ParseSlice(constant.RoutePrefix, param.RouteId))
	if err != nil {
		return err
	}

	for _, v := range routes {
		linkIds = append(linkIds, v.Id)
	}

	return s.rl.BatchLink(ctx, constant.LinkDirectCollectionRoute, item.Id, linkIds)
}

type DeleteCollectionRouteParam struct {
	Id      string   `json:"id" uri:"id" binding:"required"`
	RouteId []string `json:"route_id" form:"route_id" binding:"required,max=1000"`
}

func (s *collection) DeleteRoute(ctx context.Context, param *DeleteCollectionRouteParam) error {
	item, err := s.r.Get(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return err
	}

	linkIds := []uint64{}
	routes, err := s.rr.BatchGet(ctx, identity.ParseSlice(constant.RoutePrefix, param.RouteId))
	if err != nil {
		return err
	}

	for _, v := range routes {
		linkIds = append(linkIds, v.Id)
	}

	return s.rl.BatchDeleteLink(ctx, constant.LinkDirectCollectionRoute, item.Id, linkIds)
}

type LinkCollectionEndpointParam struct {
	Id         string   `json:"id" uri:"id" binding:"required"`
	EndpointId []string `json:"endpoint_id" form:"endpoint_id" binding:"required"`
}

func (s *collection) LinkEndpoint(ctx context.Context, param *LinkCollectionEndpointParam) error {
	item, err := s.r.Get(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return err
	}

	linkIds := []uint64{}
	endpoints, err := s.re.BatchGet(ctx, identity.ParseSlice(constant.EndpointPrefix, param.EndpointId))
	if err != nil {
		return err
	}

	for _, v := range endpoints {
		linkIds = append(linkIds, v.Id)
	}

	return s.rl.BatchLink(ctx, constant.LinkDirectCollectionEndpoint, item.Id, linkIds)
}

type DeleteCollectionEndpointParam struct {
	Id         string   `json:"id" uri:"id" binding:"required"`
	EndpointId []string `json:"endpoint_id" form:"endpoint_id" binding:"required,max=1000"`
}

func (s *collection) DeleteEndpoint(ctx context.Context, param *DeleteCollectionEndpointParam) error {
	item, err := s.r.Get(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return err
	}

	linkIds := []uint64{}
	endpoints, err := s.re.BatchGet(ctx, identity.ParseSlice(constant.EndpointPrefix, param.EndpointId))
	if err != nil {
		return err
	}

	for _, v := range endpoints {
		linkIds = append(linkIds, v.Id)
	}

	return s.rl.BatchDeleteLink(ctx, constant.LinkDirectCollectionEndpoint, item.Id, linkIds)
}

type ListCollectionParam struct {
	ParentId      string `form:"parent_id"`
	Name          string `form:"name"`
	Deep          int    `form:"deep" binding:"max=10"`
	Limit         int    `form:"limit" binding:"max=1000"`
	StartingAfter string `form:"starting_after"`
	EndingBefore  string `form:"ending_before"`
}

type ListCollectionResult struct {
	HasMore bool              `json:"has_more"`
	Data    []*dto.Collection `json:"data"`
}

func (s *collection) List(ctx context.Context, param *ListCollectionParam) (*ListCollectionResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListCollectionParam{
		Name:          param.Name,
		ParentId:      identity.Parse(constant.CollectionPrefix, param.ParentId),
		Deep:          param.Deep,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.CollectionPrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.CollectionPrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.Collection, n)
	for i, v := range entities {
		items[i] = dto.NewCollection(v)
	}

	rst := &ListCollectionResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

type UpdateCollectionParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateCollectionParam
}

func (s *collection) Update(ctx context.Context, param *UpdateCollectionParam) (*dto.Collection, error) {
	data_source.Transaction(ctx, func(txCtx context.Context) error {
		id := identity.Parse(constant.ServerNamePrefix, param.Id)

		err := s.r.Update(ctx, id, &entity.Collection{
			Name: param.Name,
		})

		if err != nil {
			return err
		}

		if err := s.batchLink(ctx, constant.LinkDirectCollectionServerName, id, identity.ParseSlice(constant.ServerNamePrefix, param.ServerNameId)); err != nil {
			return err
		}

		if err := s.batchLink(ctx, constant.LinkDirectCollectionEndpoint, id, identity.ParseSlice(constant.EndpointPrefix, param.EndpointId)); err != nil {
			return err
		}

		return nil
	})

	return s.Get(ctx, &GetCollectionParam{Id: param.Id})
}
