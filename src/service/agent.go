package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	ag "dxkite.cn/meownest/pkg/agent"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
)

type Agent interface {
	Run(addr string)
	LoadRoute(ctx context.Context) error
}

type agent struct {
	svr *ag.Server
	rc  repository.Collection
	rl  repository.Link
	rr  repository.Route
	re  repository.Endpoint
	ra  repository.Authorize
}

func NewAgent(svr *ag.Server, rr repository.Route, rc repository.Collection, re repository.Endpoint, ra repository.Authorize, rl repository.Link) Agent {
	return &agent{svr: svr, rr: rr, rc: rc, rl: rl, re: re, ra: ra}
}

func (s *agent) Run(addr string) {
	s.svr.Run(addr)
}

func (s *agent) LoadRoute(ctx context.Context) error {
	route := ag.NewHandler()
	if err := s.rr.Batch(ctx, func(item *entity.Route) error {
		forward, err := s.createForwardItem(ctx, item)
		if err != nil {
			printLog("skip route %v %s %s\n", item.Method, item.Path, err.Error())
			return nil
		}
		route.Add(forward)

		printLog("load route %v %s\n", item.Method, item.Path)
		return nil
	}); err != nil {
		return err
	}

	route.Sort()
	s.svr.Use(route)
	return nil
}

func (s *agent) createForwardItem(ctx context.Context, item *entity.Route) (ag.ForwardHandler, error) {
	collectionIdList, err := s.getCollectionList(ctx, item)
	if err != nil {
		return nil, err
	}

	collections, err := s.rc.BatchGet(ctx, collectionIdList)
	if err != nil {
		return nil, err
	}

	collectionMap := map[uint64]*entity.Collection{}
	for i, v := range collections {
		collectionMap[v.Id] = collections[i]
	}

	endpoint, err := s.getEndpoint(ctx, item, collectionIdList, collectionMap)
	if err != nil {
		return nil, err
	}

	if endpoint == nil {
		return nil, errors.New("missing endpoint")
	}

	authorize, err := s.getAuthorize(ctx, item, collectionIdList, collectionMap)
	if err != nil {
		return nil, err
	}

	forwardItem := NewForwardHandler(item, []string{}, endpoint, authorize)
	return forwardItem, nil
}

func NewEndpointForwardHandler(endpoint *entity.Endpoint) ag.RequestForwardHandler {
	targets := []*ag.EndpointTarget{}
	for _, v := range endpoint.Endpoint.Static.Address {
		targets = append(targets, &ag.EndpointTarget{
			Network: v.Network,
			Address: v.Address,
		})
	}

	handler := ag.NewStaticForwardHandler(targets, endpoint.Endpoint.Static.Timeout)
	return handler
}

func NewForwardHandler(item *entity.Route, serverNameList []string, endpoint *entity.Endpoint, auth *entity.Authorize) ag.ForwardHandler {
	matcher := ag.NewBasicMatcher()
	matcher.Path = ag.NewRequestPathMatcher(item.Path)
	matcher.Method = item.Method
	matcher.Extra = []*ag.ExtraMatchOption{}

	matcher.Host = serverNameList

	for _, v := range item.MatchOptions {
		matcher.Extra = append(matcher.Extra, &ag.ExtraMatchOption{
			Source: v.Source,
			Type:   v.Type,
			Name:   v.Name,
			Value:  v.Value,
		})
	}

	var authHandler ag.AuthorizeHandler
	if auth != nil {
		authHandler = NewAuthorizeHandler(auth)
	}

	handler := NewEndpointForwardHandler(endpoint)
	return ag.NewForwardHandler(matcher, handler, authHandler)
}

func NewAuthorizeHandler(auth *entity.Authorize) ag.AuthorizeHandler {
	binary := auth.Attribute.Binary
	source := []*ag.AuthorizeSource{}
	for _, v := range binary.Sources {
		source = append(source, &ag.AuthorizeSource{Source: v.Source, Name: v.Name})
	}
	return ag.NewBinaryAuth(binary.Key, binary.Header, source)
}

func (s *agent) getEndpoint(ctx context.Context, route *entity.Route, collectionIdList []uint64, collectionMap map[uint64]*entity.Collection) (*entity.Endpoint, error) {
	if route.EndpointId != 0 {
		item, err := s.re.Get(ctx, route.EndpointId)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	for _, v := range collectionIdList {
		if coll, ok := collectionMap[v]; ok {
			if coll.EndpointId != 0 {
				endpoint, err := s.re.Get(ctx, coll.EndpointId)
				if err != nil {
					return nil, err
				}
				return endpoint, nil
			}
		}
	}
	return nil, nil
}

func (s *agent) getAuthorize(ctx context.Context, route *entity.Route, collectionIdList []uint64, collectionMap map[uint64]*entity.Collection) (*entity.Authorize, error) {
	if route.AuthorizeId != 0 {
		item, err := s.ra.Get(ctx, route.AuthorizeId)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	for _, v := range collectionIdList {
		if coll, ok := collectionMap[v]; ok {
			if coll.AuthorizeId != 0 {
				item, err := s.ra.Get(ctx, coll.AuthorizeId)
				if err != nil {
					return nil, err
				}
				return item, nil
			}
		}
	}
	return nil, nil
}

func (s *agent) getCollectionList(ctx context.Context, item *entity.Route) ([]uint64, error) {
	collection := []uint64{}
	source, err := s.rc.Get(ctx, item.CollectionId)
	if err != nil {
		return nil, err
	}

	collection = append(collection, source.Id)

	idList := strings.Split(source.Index, ".")

	for i := len(idList) - 1; i >= 0; i-- {
		id, _ := strconv.ParseUint(idList[i], 10, 64)
		if id > 0 {
			collection = append(collection, id)
		}
	}

	return collection, nil
}

func printLog(format string, values ...interface{}) {
	fmt.Printf(format, values...)
}
