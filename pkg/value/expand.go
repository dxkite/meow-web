package value

import (
	"context"
	"encoding/json"

	"dxkite.cn/nebula/pkg/depends"
)

type ExpandItemProvider[T any, E any] interface {
	ExpandItem(context.Context, T) (E, error)
}

type ExpandItemsProvider[T any, E any] interface {
	ExpandItems(context.Context, []T) ([]E, error)
}

type Expand[T any, E any, P ExpandItemProvider[T, E]] struct {
	item T
	ctx  context.Context
}

func (e *Expand[T, E, P]) Set(ctx context.Context, item T) {
	e.item = item
	e.ctx = ctx
}

func (e *Expand[T, E, P]) Get() (item T) {
	return e.item
}

func (e *Expand[T, E, P]) GetExpand() (E, error) {
	provider, err := depends.Resolve[P]()
	if err != nil {
		var emptyE E
		return emptyE, err
	}
	return provider.ExpandItem(e.ctx, e.item)
}

func (e *Expand[T, E, P]) MarshalJSON() ([]byte, error) {
	expand, err := e.GetExpand()
	if err != nil {
		return nil, err
	}
	return json.Marshal(expand)
}

type Expands[T any, E any, P ExpandItemsProvider[T, E]] struct {
	item []T
	ctx  context.Context
}

func (e *Expands[T, E, P]) Set(ctx context.Context, item []T) {
	e.item = item
	e.ctx = ctx
}

func (e *Expands[T, E, P]) Get() (item []T) {
	return e.item
}

func (e *Expands[T, E, P]) GetExpand() ([]E, error) {
	provider, err := depends.Resolve[P]()
	if err != nil {
		return nil, err
	}
	return provider.ExpandItems(e.ctx, e.item)
}

func (e *Expands[T, E, P]) MarshalJSON() ([]byte, error) {
	expand, err := e.GetExpand()
	if err != nil {
		return nil, err
	}
	return json.Marshal(expand)
}
