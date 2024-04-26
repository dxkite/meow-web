package identity

import (
	"encoding/base64"
	"encoding/binary"
	"strings"
)

var encodeId = "RSTUVde-fghijklABOPQmnopCDstu_vwxyz012345EFGHIJWXYZabcKLMNqr6789"
var IdEncoding = base64.NewEncoding(encodeId).WithPadding(base64.NoPadding)
var mask uint64 = 1723627081864056832

func Encode(id uint64) string {
	id = id ^ mask
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return IdEncoding.EncodeToString(b)
}

func Decode(id string) uint64 {
	v, err := IdEncoding.DecodeString(id)
	if err != nil {
		return 0
	}
	vv := binary.BigEndian.Uint64(v)
	return vv ^ mask
}

func Format(prefix string, id uint64) string {
	return prefix + Encode(id)
}

func Parse(prefix, id string) uint64 {
	id = strings.TrimPrefix(id, prefix)
	return Decode(id)
}
