package errors

import (
	"fmt"
	"testing"
)

func TestIsNotFound(t *testing.T) {
	err := NotFound(New("test"))
	fmt.Println(IsNotFound(err))
}
