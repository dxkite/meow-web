package container

import (
	"context"
	"fmt"
	"reflect"
)

type Container interface {
	Register(instance Instance) error
	Get(ctx context.Context, id InstanceId) (any, error)
}

type InstanceId string

type Instance interface {
	Id() InstanceId
	New(...any) (any, error)
	Requires() []InstanceId
}

type container struct {
	instances map[InstanceId]Instance
}

func New() Container {
	return &container{instances: make(map[InstanceId]Instance)}
}

func (c *container) Register(instance Instance) error {
	c.instances[instance.Id()] = instance
	return nil
}

func (c *container) Get(ctx context.Context, id InstanceId) (any, error) {
	val := ctx.Value(scopeKey{})
	scope, _ := val.(Scope)

	if scope != nil && scope.Get(id) != nil {
		return scope.Get(id), nil
	}

	ins, ok := c.instances[id]
	if !ok {
		return nil, fmt.Errorf("%s not found", id)
	}

	requireIdList := ins.Requires()
	requireLen := len(requireIdList)

	requireArgs := make([]any, requireLen)
	for i, required := range requireIdList {
		obj, err := c.Get(ctx, required)
		if err != nil {
			return nil, err
		}
		if scope != nil {
			scope.Set(required, obj)
		}
		requireArgs[i] = obj
	}

	obj, err := ins.New(requireArgs...)
	if err != nil {
		return nil, err
	}
	if scope != nil {
		scope.Set(id, obj)
	}
	return obj, nil
}

func CreateInstanceId(val reflect.Type) InstanceId {
	switch val.Kind() {
	case reflect.Pointer:
		return "*" + CreateInstanceId(val.Elem())
	}
	return InstanceId(fmt.Sprintf("%s.%s", val.PkgPath(), val.Name()))
}
