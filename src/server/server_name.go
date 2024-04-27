package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewServerName(s service.ServerName) *ServerName {
	return &ServerName{s: s}
}

type ServerName struct {
	s service.ServerName
}

func (s *ServerName) Create(c *gin.Context) {
	var param service.CreateServerNameParam

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

func (s *ServerName) Get(c *gin.Context) {
	var param service.GetServerNameParam

	if err := c.ShouldBindUri(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	Result(c, http.StatusCreated, rst)
}

func WithServerName(path string, server *ServerName) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
			group.GET("/:id", server.Get)
		}
	}
}
