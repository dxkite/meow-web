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
		ResultErrorBind(c, err)
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
		ResultErrorBind(c, err)
		return
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		ResultErrorBind(c, err)
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
		ResultErrorBind(c, err)
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
		ResultErrorBind(c, err)
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
		ResultErrorBind(c, err)
		return
	}
	err := s.s.Delete(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func (s *ServerName) RegisterToHttp(group gin.IRouter) {
	group.GET("/server_names", s.List)
	group.POST("/server_names", s.Create)
	group.GET("/server_names/:id", s.Get)
	group.DELETE("/server_names/:id", s.Delete)
	group.POST("/server_names/:id", s.Update)
}
