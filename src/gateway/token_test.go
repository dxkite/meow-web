package gateway

import (
	"reflect"
	"testing"
	"time"
)

func TestToken_Marshal(t *testing.T) {
	token := &Token{
		Id:              1,
		ExpireAt:        uint64(time.Now().Add(1 * time.Hour).Unix()),
		RefreshExpireAt: uint64(time.Now().Add(24 * time.Hour).Unix()),
	}

	tok := token.EncodeToString()
	token1 := &Token{}
	if err := token1.DecodeString(tok); err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(token1, token) {
		t.Errorf("token not equal %v %v", token, token1)
	}
}
