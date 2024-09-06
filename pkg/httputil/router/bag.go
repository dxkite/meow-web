package router

import "context"

type Constructor func(ctx context.Context) (Collection, error)

type Collection interface {
	Routes() []Route
}

func NewCollectionBag() *CollectionBag {
	return &CollectionBag{routes: []Constructor{}}
}

type CollectionBag struct {
	routes []Constructor
}

func (r *CollectionBag) Add(c Constructor) {
	r.routes = append(r.routes, c)
}

func (r *CollectionBag) Build(ctx context.Context) ([]Route, error) {
	routes := []Route{}
	for _, constructor := range r.routes {
		collection, err := constructor(ctx)
		if err != nil {
			return nil, err
		}
		routes = append(routes, collection.Routes()...)
	}
	return routes, nil
}
