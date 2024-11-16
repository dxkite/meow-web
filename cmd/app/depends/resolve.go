package depends

import (
	"context"

	"dxkite.cn/nebula/pkg/depends"
)

var Scope = depends.NewScopedContext(context.Background())

func Resolve[T any]() (T, error) {
	return depends.Resolve[T](Scope)
}
