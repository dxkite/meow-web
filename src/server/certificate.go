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

func (s *Certificate) RegisterToHttp(group gin.IRouter) {
	group.POST("/certificate", s.Create)
}
