package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type Collection interface {
	Create(c *gin.Context)
}

func NewCollection(s service.Collection) Collection {
	return &collection{s: s}
}

type collection struct {
	s service.Collection
}

func (s *collection) Create(c *gin.Context) {
	var param service.CreateCollectionParam

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	Result(c, http.StatusCreated, rst)
}

func WithCollection(path string, server Collection) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
		}
	}
}
