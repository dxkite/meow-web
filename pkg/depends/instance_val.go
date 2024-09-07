package depends

type valueInstance struct {
	val interface{}
}

func (i *valueInstance) New(args ...any) (any, error) {
	return i.val, nil
}

func (i *valueInstance) Requires() []InstanceId {
	return nil
}

func makeValInstance[T any](v T) Instance {
	return &valueInstance{
		val: v,
	}
}
