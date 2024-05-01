package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewEndpoint(s service.Endpoint) *Endpoint {
	return &Endpoint{s: s}
}

type Endpoint struct {
	s service.Endpoint
}

func (s *Endpoint) Create(c *gin.Context) {
	var param service.CreateEndpointParam

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.Result(c, http.StatusCreated, rst)
}

func (s *Endpoint) Get(c *gin.Context) {
	var param service.GetEndpointParam

	if err := c.ShouldBindUri(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}
	httpserver.Result(c, http.StatusOK, rst)
}

// 证书列表
//
// @Summary      证书列表
// @Description  证书列表
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        name query string false "证书"
// @Param        limit query int false "限制"
// @Param        starting_after query string false "从当前ID开始"
// @Param        ending_before query string false "从当前ID结束"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListEndpointserver.Result
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /endpoints [get]
func (s *Endpoint) List(c *gin.Context) {
	var param service.ListEndpointParam

	if err := c.ShouldBindQuery(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.List(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.Result(c, http.StatusOK, rst)
}

// 更新证书
//
// @Summary      更新证书
// @Description  更新证书
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        id path string true "证书ID"
// @Param        body body service.UpdateEndpointParam true "数据"
// @Success      200  {object} service.Endpoint
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /endpoints/{id} [post]
func (s *Endpoint) Update(c *gin.Context) {
	var param service.UpdateEndpointParam
	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Update(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}
	httpserver.Result(c, http.StatusOK, rst)
}

// 删除证书
//
// @Summary      删除证书
// @Description  删除证书
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        id path string true "证书ID"
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /endpoints/{id} [delete]
func (s *Endpoint) Delete(c *gin.Context) {
	var param service.DeleteEndpointParam

	if err := c.ShouldBindUri(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}
	err := s.s.Delete(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *Endpoint) RegisterToHttp(group gin.IRouter) {
	group.POST("/endpoints", s.Create)
	group.GET("/endpoints", s.List)
	group.GET("/endpoints/:id", s.Get)
	group.POST("/endpoints/:id", s.Update)
	group.DELETE("/endpoints/:id", s.Delete)
}
