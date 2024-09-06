package container

type valueInstance struct {
	val interface{}
	id  InstanceId
}

func (i *valueInstance) New(args ...any) (any, error) {
	return i.val, nil
}

func (i *valueInstance) Requires() []InstanceId {
	return nil
}

func NewValueInstance[T any](v T) Instance {
	return &valueInstance{
		val: v,
		id:  makeInstanceId(v),
	}
}

func (i *valueInstance) Id() InstanceId {
	return i.id
}
