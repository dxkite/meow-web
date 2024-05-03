package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	ag "dxkite.cn/meownest/pkg/agent"
	"dxkite.cn/meownest/src/constant"
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
	rs  repository.ServerName
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
	s.svr.Use(route)
	return nil
}

func (s *agent) createForwardItem(ctx context.Context, item *entity.Route) (ag.ForwardHandler, error) {
	collectionIdList, err := s.getCollectionList(ctx, item)
	if err != nil {
		return nil, err
	}

	endpoints, err := s.getEndpoint(ctx, item.Id, collectionIdList)
	if err != nil {
		return nil, err
	}

	if len(endpoints) == 0 {
		return nil, errors.New("missing endpoint")
	}

	endpoint := endpoints[0]
	forwardItem := NewForwardHandler(item, endpoint, nil)
	return forwardItem, nil
}

func NewForwardHandler(item *entity.Route, endpoint *entity.Endpoint, auth *entity.Authorize) ag.ForwardHandler {
	targets := []*ag.EndpointTarget{}
	for _, v := range endpoint.Endpoint.Static.Address {
		targets = append(targets, &ag.EndpointTarget{
			Network: v.Network,
			Address: v.Address,
		})
	}
	matcher := ag.NewBasicMatcher()
	matcher.Path = item.Path
	matcher.Method = item.Method
	matcher.Extra = []*ag.ExtraMatchOption{}
	for _, v := range item.MatchOptions {
		matcher.Extra = append(matcher.Extra, &ag.ExtraMatchOption{
			Source: v.Source,
			Type:   v.Type,
			Name:   v.Name,
			Value:  v.Value,
		})
	}
	handler := ag.NewStaticForwardHandler(targets, endpoint.Endpoint.Static.Timeout)

	var authHandler ag.AuthorizeHandler
	if auth != nil {
		authHandler = NewAuthorizeHandler(auth)
	}

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

func (s *agent) getEndpoint(ctx context.Context, routeId uint64, collectionIdList []uint64) ([]*entity.Endpoint, error) {

	if endpoints, err := s.getEndpointBy(ctx, constant.LinkDirectRouteEndpoint, routeId); err != nil {
		return nil, err
	} else if len(endpoints) > 0 {
		return endpoints, nil
	}

	for _, v := range collectionIdList {
		if endpoints, err := s.getEndpointBy(ctx, constant.LinkDirectCollectionEndpoint, v); err != nil {
			return nil, err
		} else if len(endpoints) > 0 {
			return endpoints, nil
		}
	}

	return nil, nil
}

func (s *agent) getAuthorize(ctx context.Context, routeId uint64, collectionIdList []uint64) (*entity.Authorize, error) {
	if auth, err := s.getAuthorizeBy(ctx, constant.LinkDirectRouteAuthorize, routeId); err != nil {
		return nil, err
	} else if auth != nil {
		return auth, nil
	}

	for _, v := range collectionIdList {
		if auth, err := s.getAuthorizeBy(ctx, constant.LinkDirectCollectionAuthorize, v); err != nil {
			return nil, err
		} else if auth != nil {
			return auth, nil
		}
	}

	return nil, nil
}

func (s *agent) getAuthorizeBy(ctx context.Context, direct string, id uint64) (*entity.Authorize, error) {
	linked, err := s.rl.Linked(ctx, direct, []uint64{id})
	if err != nil {
		return nil, err
	}

	if len(linked) > 0 {
		ent, err := s.ra.Get(ctx, linked[0].Id)
		if err != nil {
			return nil, err
		}
		return ent, nil
	}
	return nil, nil
}

func (s *agent) getEndpointBy(ctx context.Context, direct string, id uint64) ([]*entity.Endpoint, error) {
	var endpointList []*entity.Endpoint
	endpointLink, err := s.rl.Linked(ctx, direct, []uint64{id})

	if err != nil {
		return nil, err
	}

	if len(endpointLink) > 0 {
		linkedId := linkedIds(endpointLink)
		endpoints, err := s.re.BatchGet(ctx, linkedId)
		if err != nil {
			return nil, err
		}
		return endpoints, nil
	}

	return endpointList, nil
}

func (s *agent) getCollectionList(ctx context.Context, item *entity.Route) ([]uint64, error) {
	sources, err := s.rl.LinkedSource(ctx, constant.LinkDirectCollectionRoute, []uint64{item.Id})
	if err != nil {
		return nil, err
	}

	if len(sources) == 0 {
		return nil, errors.New("missing collection")
	}

	collection := []uint64{}

	sourceLink := sources[0]
	source, err := s.rc.Get(ctx, sourceLink.SourceId)
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

func linkedIds(item []*entity.Link) []uint64 {
	idList := []uint64{}
	for _, v := range item {
		idList = append(idList, v.LinkedId)
	}
	return idList
}
