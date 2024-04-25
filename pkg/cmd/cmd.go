package cmd

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type CommandConstructor func() Command

type Command interface {
	Execute(ctx context.Context) CommandResult
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func Exec(newCommand CommandConstructor) gin.HandlerFunc {
	return func(c *gin.Context) {
		cmd := newCommand()

		if err := c.ShouldBind(&cmd); err != nil {
			rst := NewErrorResult(http.StatusBadRequest, "invalid_parameter", err.Error())
			rst.Write(c)
			return
		}

		rst := cmd.Execute(c)
		if rst != nil {
			rst.Write(c)
			return
		}

		c.Status(http.StatusOK)
	}
}
