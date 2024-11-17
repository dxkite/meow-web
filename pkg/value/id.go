package value

import (
	"context"
	"encoding/json"

	"dxkite.cn/nebula/pkg/crypto/identity"
)

type Identity struct {
	id     uint64
	prefix string
}

func (id *Identity) Set(ctx context.Context, prefix string, v uint64) {
	id.prefix = prefix
	id.id = v
}

func (id *Identity) Get() uint64 {
	return id.id
}
func (id *Identity) MarshalJSON() ([]byte, error) {
	val := identity.Format(id.prefix, id.id)
	return json.Marshal(val)
}
