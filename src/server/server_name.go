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
		ResultError(c, err)
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
		ResultError(c, err)
		return
	}
	Result(c, http.StatusOK, rst)
}

func (s *ServerName) List(c *gin.Context) {
	var param service.ListServerNameParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.List(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}

	Result(c, http.StatusOK, rst)
}

func (s *ServerName) Update(c *gin.Context) {
	var param service.UpdateServerNameParam
	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Update(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}
	Result(c, http.StatusOK, rst)
}

func (s *ServerName) Delete(c *gin.Context) {
	var param service.DeleteServerNameParam

	if err := c.ShouldBindUri(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	err := s.s.Delete(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func WithServerName(path string, server *ServerName) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.GET("", server.List)
			group.POST("", server.Create)
			group.GET("/:id", server.Get)
			group.DELETE("/:id", server.Delete)
			group.POST("/:id", server.Update)
		}
	}
}
