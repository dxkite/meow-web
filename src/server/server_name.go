package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type ServerName interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
}

func NewServerName(s service.ServerName) ServerName {
	return &serverName{s: s}
}

type serverName struct {
	s service.ServerName
}

func (s *serverName) Create(c *gin.Context) {
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

func (s *serverName) Get(c *gin.Context) {
	rst, err := s.s.Get(c, c.Param("id"), c.QueryArray("expand"))
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	Result(c, http.StatusCreated, rst)
}
