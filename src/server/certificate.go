package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type Certificate interface {
	Create(c *gin.Context)
}

func NewCertificate(s service.Certificate) Certificate {
	return &certificate{s: s}
}

type certificate struct {
	s service.Certificate
}

func (s *certificate) Create(c *gin.Context) {
	var param service.CreateCertificateParam

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

func WithCertificate(path string, server Certificate) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
		}
	}
}
