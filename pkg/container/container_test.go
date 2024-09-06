package container

import (
	"context"
	"testing"
)

type ServiceA interface {
}

type ServiceB interface {
}

type ServiceC interface {
}

type ServiceD interface {
}

type test struct {
	name string
}

func TestNewFuncInstance(t *testing.T) {
	container := New()

	invokeCount := map[string]int{}

	ContainerRegister(container, func() ServiceA {
		invokeCount["ServiceA"]++
		return test{"ServiceA"}
	})
	ContainerRegister(container, func() ServiceB {
		invokeCount["ServiceB"]++
		return test{"ServiceB"}
	})
	ContainerRegister(container, func(a ServiceB) ServiceC {
		invokeCount["ServiceC"]++
		return test{"ServiceC"}
	})
	ContainerRegister(container, func(a ServiceA, b ServiceB, c ServiceC) ServiceD {
		invokeCount["ServiceD"]++
		return test{"ServiceD"}
	})

	ctx := NewScopedContext(context.TODO())
	_, err := ContainerGet[ServiceD](ctx, container)

	if err != nil {
		t.Error(err)
		return
	}

	for k := range invokeCount {
		if invokeCount[k] != 1 {
			t.Errorf("invokeCount[%s] should be 1, but got %d", k, invokeCount[k])
		}
	}
}
