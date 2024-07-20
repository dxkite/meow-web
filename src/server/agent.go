package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type Agent struct {
	s service.Agent
}

func NewAgent(s service.Agent) *Agent {
	return &Agent{s: s}
}

// 重载代理服务路由
//
// @Summary      重载代理服务路由
// @Description  重载代理服务路由
// @Tags         Agent
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /agent/reload [post]
func (s *Agent) Reload(c *gin.Context) {
	s.s.LoadRoute(c)
	c.Status(http.StatusOK)
}

func (s *Agent) API() httputil.RouteHandleFunc {
	return func(r gin.IRouter) {
		r.POST("/agent/reload", s.Reload)
	}
}
