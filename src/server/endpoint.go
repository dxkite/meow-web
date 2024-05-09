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

// Create Endpoint
//
// @Summary      Create Endpoint
// @Description  Create Endpoint
// @Tags         Endpoint
// @Accept       json
// @Produce      json
// @Param        body body service.CreateEndpointParam true "Endpoint data"
// @Success      200  {object} dto.Endpoint
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /endpoints [post]
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

// Get Endpoint
//
// @Summary      Get Endpoint
// @Description  Get Endpoint
// @Tags         Endpoint
// @Accept       json
// @Produce      json
// @Param        id path string true "Endpoint ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} dto.Endpoint
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /endpoints/{id} [get]
func (s *Endpoint) Get(c *gin.Context) {
	var param service.GetEndpointParam

	param.Id = c.Param("id")

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

// List Endpoint
//
// @Summary      Endpoint list
// @Description  Endpoint list
// @Tags         Endpoint
// @Accept       json
// @Produce      json
// @Param        name query string false "Endpoint"
// @Param		 include_total query bool false "是否包含total"
// @Param        page query int false "页码"
// @Param        pre_page query int false "每页数量"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListEndpointResult
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

// Update Endpoint
//
// @Summary      Update Endpoint
// @Description  Update Endpoint
// @Tags         Endpoint
// @Accept       json
// @Produce      json
// @Param        id path string true "Endpoint ID"
// @Param        body body service.UpdateEndpointParam true "data"
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

// Delete Endpoint
//
// @Summary      Delete Endpoint
// @Description  Delete Endpoint
// @Tags         Endpoint
// @Accept       json
// @Produce      json
// @Param        id path string true "EndpointID"
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

func (s *Endpoint) RegisterToHttp(route gin.IRouter) {
	route.POST("/endpoints", s.Create)
	route.GET("/endpoints", s.List)
	route.GET("/endpoints/:id", s.Get)
	route.POST("/endpoints/:id", s.Update)
	route.DELETE("/endpoints/:id", s.Delete)
}
