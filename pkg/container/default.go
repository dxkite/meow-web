package container

import "context"

var Default = New()

func Get[T any](ctx context.Context) (T, error) {
	return ContainerGet[T](ctx, Default)
}

func Register[T any](obj T) error {
	return ContainerRegister[T](Default, obj)
}
