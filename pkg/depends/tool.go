package depends

import (
	"context"
)

func ResolveService[T any](ctx context.Context, container Container) (T, error) {
	var nilInstance T
	id := MakeInstanceID(nilInstance)
	obj, err := container.Resolve(ctx, id)
	if err != nil {
		return nilInstance, err
	}
	return obj.(T), nil
}

func RegisterService[T any](container Container, obj T) error {
	id := MakeInstanceID(obj)
	return container.Register(id, makeInstance(obj))
}
