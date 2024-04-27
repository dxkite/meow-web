package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewCertificate(s service.Certificate) *Certificate {
	return &Certificate{s: s}
}

type Certificate struct {
	s service.Certificate
}

func (s *Certificate) Create(c *gin.Context) {
	var param service.CreateCertificateParam

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

func WithCertificate(path string, server *Certificate) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
		}
	}
}
