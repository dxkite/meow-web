package depends

import (
	"fmt"
	"reflect"
)

type funcInstance struct {
	val reflect.Value
	typ reflect.Type
}

func (i *funcInstance) New(args ...any) (any, error) {
	in := make([]reflect.Value, len(args))
	for i := range args {
		in[i] = reflect.ValueOf(args[i])
	}
	out := i.val.Call(in)
	if len(out) == 0 {
		return nil, fmt.Errorf("invalid return")
	}

	latest := out[len(out)-1]

	isErr := latest.Type().Implements(reflect.TypeOf((*error)(nil)).Elem())
	if isErr && latest.Interface() != nil {
		return nil, latest.Interface().(error)
	}

	return out[0].Interface(), nil
}

func (i *funcInstance) Requires() []InstanceId {
	required := make([]InstanceId, i.typ.NumIn())
	for j := 0; j < i.typ.NumIn(); j++ {
		required[j] = makeInstanceId(i.typ.In(j))
	}
	return required
}

func makeFunInstance(f any) Instance {
	return &funcInstance{
		val: reflect.ValueOf(f),
		typ: reflect.TypeOf(f),
	}
}
