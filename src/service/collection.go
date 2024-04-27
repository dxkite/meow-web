package service

import (
	"context"

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
	LinkRoute(ctx context.Context, param *LinkRouteParam) error
	DeleteRoute(ctx context.Context, param *DeleteRouteParam) error
}

func NewCollection(r repository.Collection, rl repository.Link, rr repository.Route) Collection {
	return &collection{r: r, rr: rr, rl: rl}
}

type collection struct {
	r  repository.Collection
	rl repository.Link
	rr repository.Route
}

type CreateCollectionParam struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	ParentId    string `json:"parent_id" form:"parent_id"`
}

func (s *collection) Create(ctx context.Context, param *CreateCollectionParam) (*dto.Collection, error) {
	rst, err := s.r.Create(ctx, &entity.Collection{
		Name:        param.Name,
		Description: param.Description,
		ParentId:    identity.Parse(constant.CollectionPrefix, param.ParentId),
	})
	if err != nil {
		return nil, err
	}
	return dto.NewCollection(rst), nil
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
		routeIds := []uint64{}

		linked, err := s.rl.LinkOf(ctx, constant.LinkDirectCollectionRoute, rst.Id)
		if err != nil {
			return nil, err
		}

		for _, v := range linked {
			routeIds = append(routeIds, v.LinkedId)
		}

		routes, err := s.rr.BatchGet(ctx, routeIds)
		if err != nil {
			return nil, err
		}

		routeDtos := make([]*dto.Route, len(routes))
		for i, v := range routes {
			routeDtos[i] = dto.NewRoute(v)
		}

		collection.Routes = routeDtos
	}

	return collection, nil
}

type LinkRouteParam struct {
	Id      string   `json:"id" uri:"id" binding:"required"`
	RouteId []string `json:"route_id" form:"route_id" binding:"required"`
}

func (s *collection) LinkRoute(ctx context.Context, param *LinkRouteParam) error {
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

type DeleteRouteParam struct {
	Id      string   `json:"id" uri:"id" binding:"required"`
	RouteId []string `json:"route_id" form:"route_id" binding:"required,max=1000"`
}

func (s *collection) DeleteRoute(ctx context.Context, param *DeleteRouteParam) error {
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
