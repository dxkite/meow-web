package container

import "context"

var Default = New()

func Get[T any](ctx context.Context) (T, error) {
	return ContainerGet[T](ctx, Default)
}

func Register(obj any) error {
	return ContainerRegister(Default, obj)
}
