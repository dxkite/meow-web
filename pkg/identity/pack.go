package identity

import (
	"encoding/base64"
	"encoding/binary"
	"strings"
)

var encodeId = "RSTUVde-fghijklABOPQmnopCDstu_vwxyz012345EFGHIJWXYZabcKLMNqr6789"
var Encoding = base64.NewEncoding(encodeId).WithPadding(base64.NoPadding)
var DefaultMask uint64 = 1723627081864056832

func Encode(id uint64) string {
	id = id ^ DefaultMask
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return Encoding.EncodeToString(b)
}

func Decode(id string) uint64 {
	v, err := Encoding.DecodeString(id)
	if err != nil {
		return 0
	}
	vv := binary.BigEndian.Uint64(v)
	return vv ^ DefaultMask
}

func EncodeMask(id, mask uint64) string {
	id = id ^ mask
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return Encoding.EncodeToString(b)
}

func DecodeMask(id string, mask uint64) uint64 {
	v, err := Encoding.DecodeString(id)
	if err != nil {
		return 0
	}
	vv := binary.BigEndian.Uint64(v)
	return vv ^ mask
}

func Format(prefix string, id uint64) string {
	if id == 0 {
		return ""
	}
	return prefix + EncodeMask(id, Mask(prefix))
}

func Parse(prefix, id string) uint64 {
	if id == "" {
		return 0
	}
	id = strings.TrimPrefix(id, prefix)
	return DecodeMask(id, Mask(prefix))
}

func Mask(key string) uint64 {
	keyFull := make([]byte, 8)
	n := len(key)
	if n >= 8 {
		n = 8
	}
	for i := 0; i < n; i++ {
		keyFull[i] = key[i]
	}
	return binary.BigEndian.Uint64(keyFull) ^ DefaultMask
}
