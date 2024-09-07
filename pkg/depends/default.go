package depends

import "context"

var Default = New()

func Resolve[T any](ctx context.Context) (T, error) {
	return ResolveService[T](ctx, Default)
}

func Register[T any](obj T) error {
	return RegisterService(Default, obj)
}
