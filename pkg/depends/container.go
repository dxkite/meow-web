package depends

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type nopLogger struct {
}

func (nopLogger) Println(...any) {}

type Logger interface {
	Println(...any)
}

type Container interface {
	Register(id InstanceId, instance Instance) error
	Resolve(ctx context.Context, id InstanceId) (any, error)
}

type InstanceId string

type Instance interface {
	New(...any) (any, error)
	Requires() []InstanceId
}

type container struct {
	instances map[InstanceId]Instance
	logger    Logger
	mtx       *sync.Mutex
}

func New() Container {
	return &container{
		instances: make(map[InstanceId]Instance),
		logger:    &nopLogger{},
		mtx:       &sync.Mutex{},
	}
}

func (c *container) Register(id InstanceId, instance Instance) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.instances[id] = instance
	c.logger.Println("register", id, instance)
	return nil
}

func (c *container) Resolve(ctx context.Context, id InstanceId) (any, error) {
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
		obj, err := c.Resolve(ctx, required)
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

func makeInstanceId(typ reflect.Type) InstanceId {
	switch typ.Kind() {
	case reflect.Pointer:
		return "*" + makeInstanceId(typ.Elem())
	case reflect.Func:
		if typ.NumOut() == 0 {
			return ""
		}
		return makeInstanceId(typ.Out(0))
	}
	return InstanceId(fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name()))
}

func makeInstance(val any) Instance {
	if reflect.TypeOf(val).Kind() == reflect.Func {
		return makeFunInstance(val)
	}
	return makeValInstance(val)
}

func MakeInstanceID[T any](_ T) InstanceId {
	return makeInstanceId(reflect.TypeOf((*T)(nil)).Elem())
}
