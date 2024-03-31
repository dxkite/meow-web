package utils

import (
	"math/rand"
	"time"
)

func Intn(v int) int {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rd.Intn(v)
}
