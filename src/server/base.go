package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type WithOption func(s *HttpServer)

type RegisterHandler interface {
	RegisterToHttp(c gin.IRouter)
}

func New() *HttpServer {
	g := gin.Default()
	s := &HttpServer{engine: g}
	return s
}

type HttpServer struct {
	engine *gin.Engine
}

func (s *HttpServer) Run(addr ...string) {
	s.engine.Run(addr...)
}

func (s *HttpServer) Register(h RegisterHandler) {
	h.RegisterToHttp(s.engine)
}

func (s *HttpServer) RegisterPrefix(prefix string, h RegisterHandler) {
	h.RegisterToHttp(s.engine.Group(prefix))
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"error":         code,
			"error_message": message,
		},
	})
}

func Result(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func ResultEmpty(c *gin.Context, status int) {
	c.Status(status)
}

func ResultError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(c, http.StatusNotFound, "not_found", err.Error())
		return
	}
	Error(c, http.StatusInternalServerError, "internal_error", err.Error())
}

func ResultErrorBind(c *gin.Context, err error) {
	if e, ok := err.(validator.ValidationErrors); ok {
		errorList := []string{}

		for _, v := range e {
			customErr := fmt.Sprintf("invalid key %s (%s) by validate rule %s (%s)", v.Field(), v.Namespace(), v.Tag(), v.Param())
			errorList = append(errorList, customErr)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"error":         "invalid_parameter",
				"error_details": errorList,
			},
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"error":         "invalid_parameter",
				"error_message": err.Error(),
			},
		})
	}
}
