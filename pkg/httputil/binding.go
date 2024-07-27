package httputil

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()
var initValidator = &sync.Once{}

func init() {
	initValidator.Do(func() {
		Validator.SetTagName("binding")
	})
}
