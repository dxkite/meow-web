package {{ .PrivateName }}

import (
	"net/http"

	"{{ .Pkg }}/pkg/httputil"
	"github.com/gin-gonic/gin"
)

func New{{ .Name }}Server(s {{ .Name }}Service) *{{ .Name }}Server {
	return &{{ .Name }}Server{s: s}
}

type {{ .Name }}Server struct {
	s {{ .Name }}Service
}

// Create {{ .Name }}
//
// @Summary      Create {{ .Name }}
// @Description  Create {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        body body Create{{ .Name }}Request true "{{ .Name }} data"
// @Success      200  {object} dto.{{ .Name }}
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .URI }} [post]
func (s *{{ .Name }}Server) Create(c *gin.Context) {
	var param Create{{ .Name }}Request

	if err := c.ShouldBind(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		httputil.ResultError(c, err)
		return
	}

	httputil.Result(c, http.StatusCreated, rst)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .URI }}/{id} [get]
func (s *{{ .Name }}Server) Get(c *gin.Context) {
	var param Get{{ .Name }}Request

	param.Id = c.Param("id")

	if err := c.ShouldBindQuery(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		httputil.ResultError(c, err)
		return
	}
	httputil.Result(c, http.StatusOK, rst)
}

// List {{ .Name }}
//
// @Summary      {{ .Name }} list
// @Description  {{ .Name }} list
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        name query string false "{{ .Name }}"
// @Param        include_total query bool false "include total count"
// @Param        page query int false "page"
// @Param        pre_page query int false "size per page"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} List{{ .Name }}Result
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .URI }} [get]
func (s *{{ .Name }}Server) List(c *gin.Context) {
	var param List{{ .Name }}Request

	if err := c.ShouldBindQuery(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.List(c, &param)
	if err != nil {
		httputil.ResultError(c, err)
		return
	}

	httputil.Result(c, http.StatusOK, rst)
}

// Update {{ .Name }}
//
// @Summary      Update {{ .Name }}
// @Description  Update {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }} ID"
// @Param        body body Update{{ .Name }}Request true "data"
// @Success      200  {object} dto.{{ .Name }}
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .URI }}/{id} [post]
func (s *{{ .Name }}Server) Update(c *gin.Context) {
	var param Update{{ .Name }}Request
	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Update(c, &param)
	if err != nil {
		httputil.ResultError(c, err)
		return
	}
	httputil.Result(c, http.StatusOK, rst)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .URI }}/{id} [delete]
func (s *{{ .Name }}Server) Delete(c *gin.Context) {
	var param Delete{{ .Name }}Request

	if err := c.ShouldBindUri(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}
	err := s.s.Delete(c, &param)
	if err != nil {
		httputil.ResultError(c, err)
		return
	}

	httputil.ResultEmpty(c, http.StatusOK)
}

func (s *{{ .Name }}Server) API() httputil.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.POST("/{{ .URI }}", s.Create)
		route.GET("/{{ .URI }}", s.List)
		route.GET("/{{ .URI }}/:id", s.Get)
		route.POST("/{{ .URI }}/:id", s.Update)
		route.DELETE("/{{ .URI }}/:id", s.Delete)
	}
}