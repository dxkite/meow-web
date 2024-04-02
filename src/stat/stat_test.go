package stat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"dxkite.cn/meownest/src/binary"
)

func TestDynamic(t *testing.T) {
	got, err := Dynamic()
	if err != nil {
		t.Error(err)
		return
	}
	val, _ := json.Marshal(got)
	fmt.Println(string(val))

	buf := &bytes.Buffer{}

	if err = binary.Write(buf, got); err != nil {
		t.Error(err)
		return
	}

	fmt.Println(buf.Len(), strconv.Quote(buf.String()))
}

func TestSystem(t *testing.T) {
	got, err := System()
	if err != nil {
		t.Error(err)
		return
	}
	val, _ := json.Marshal(got)
	fmt.Println(string(val))
	time.Sleep(time.Second)
}
