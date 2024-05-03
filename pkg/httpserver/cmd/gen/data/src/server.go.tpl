package server

import (
	"net/http"

	"{{ .Pkg }}/pkg/httpserver"
	"{{ .Pkg }}/src/service"
	"github.com/gin-gonic/gin"
)

func New{{ .Name }}(s service.{{ .Name }}) *{{ .Name }} {
	return &{{ .Name }}{s: s}
}

type {{ .Name }} struct {
	s service.{{ .Name }}
}

// Create {{ .Name }}
//
// @Summary      Create {{ .Name }}
// @Description  Create {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }} ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} dto.{{ .Name }}
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /{{ .URI }} [post]
func (s *{{ .Name }}) Create(c *gin.Context) {
	var param service.Create{{ .Name }}Param

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

// Get {{ .Name }}
//
// @Summary      Get {{ .Name }}
// @Description  Get {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }} ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} dto.{{ .Name }}
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /{{ .URI }}/{id} [get]
func (s *{{ .Name }}) Get(c *gin.Context) {
	var param service.Get{{ .Name }}Param

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

// {{ .Name }}列表
//
// @Summary      {{ .Name }}列表
// @Description  {{ .Name }}列表
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        name query string false "{{ .Name }}"
// @Param        limit query int false "限制"
// @Param        starting_after query string false "从当前ID开始"
// @Param        ending_before query string false "从当前ID结束"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.List{{ .Name }}Result
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /{{ .URI }} [get]
func (s *{{ .Name }}) List(c *gin.Context) {
	var param service.List{{ .Name }}Param

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

// Update {{ .Name }}
//
// @Summary      Update {{ .Name }}
// @Description  Update {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }} ID"
// @Param        body body service.Update{{ .Name }}Param true "data"
// @Success      200  {object} service.{{ .Name }}
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /{{ .URI }}/{id} [post]
func (s *{{ .Name }}) Update(c *gin.Context) {
	var param service.Update{{ .Name }}Param
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

// Delete {{ .Name }}
//
// @Summary      Delete {{ .Name }}
// @Description  Delete {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }}ID"
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /{{ .URI }}/{id} [delete]
func (s *{{ .Name }}) Delete(c *gin.Context) {
	var param service.Delete{{ .Name }}Param

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

func (s *{{ .Name }}) RegisterToHttp(route gin.IRouter) {
	route.POST("/{{ .URI }}", s.Create)
	route.GET("/{{ .URI }}", s.List)
	route.GET("/{{ .URI }}/:id", s.Get)
	route.POST("/{{ .URI }}/:id", s.Update)
	route.DELETE("/{{ .URI }}/:id", s.Delete)
}