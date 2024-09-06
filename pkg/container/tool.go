package container

import (
	"context"
	"reflect"
)

func ContainerGet[T any](ctx context.Context, container Container) (T, error) {
	var nilInstance T
	id := CreateInstanceId(reflect.TypeOf((*T)(nil)).Elem())
	obj, err := container.Get(ctx, id)
	if err != nil {
		return nilInstance, err
	}
	return obj.(T), nil
}

func ContainerRegister[T any](container Container, obj T) error {
	return container.Register(NewInstance[T](obj))
}
