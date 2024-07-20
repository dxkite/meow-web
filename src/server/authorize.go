package server

import (
	"net/http"

	"dxkite.cn/meownest/src/constant"

	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewAuthorize(s service.Authorize) *Authorize {
	return &Authorize{s: s}
}

type Authorize struct {
	s service.Authorize
}

// Create Authorize
//
// @Summary      Create Authorize
// @Description  Create Authorize
// @Tags         Authorize
// @Accept       json
// @Produce      json
// @Param        body body service.CreateAuthorizeParam true "数据"
// @Success      200  {object} dto.Authorize
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /authorizes [post]
func (s *Authorize) Create(c *gin.Context) {
	var param service.CreateAuthorizeParam

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

// Get Authorize
//
// @Summary      Get Authorize
// @Description  Get Authorize
// @Tags         Authorize
// @Accept       json
// @Produce      json
// @Param        id path string true "Authorize ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} dto.Authorize
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /authorizes/{id} [get]
func (s *Authorize) Get(c *gin.Context) {
	var param service.GetAuthorizeParam

	if err := c.ShouldBindUri(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}

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

// Authorize列表
//
// @Summary      Authorize列表
// @Description  Authorize列表
// @Tags         Authorize
// @Accept       json
// @Produce      json
// @Param        name query string false "Authorize"
// @Param		 include_total query bool false "是否包含total"
// @Param        page query int false "页码"
// @Param        pre_page query int false "每页数量"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListAuthorizeResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /authorizes [get]
func (s *Authorize) List(c *gin.Context) {
	var param service.ListAuthorizeParam

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

// Update Authorize
//
// @Summary      Update Authorize
// @Description  Update Authorize
// @Tags         Authorize
// @Accept       json
// @Produce      json
// @Param        id path string true "Authorize ID"
// @Param        body body service.UpdateAuthorizeParam true "data"
// @Success      200  {object} dto.Authorize
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /authorizes/{id} [post]
func (s *Authorize) Update(c *gin.Context) {
	var param service.UpdateAuthorizeParam
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

// Delete Authorize
//
// @Summary      Delete Authorize
// @Description  Delete Authorize
// @Tags         Authorize
// @Accept       json
// @Produce      json
// @Param        id path string true "AuthorizeID"
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /authorizes/{id} [delete]
func (s *Authorize) Delete(c *gin.Context) {
	var param service.DeleteAuthorizeParam

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

func (s *Authorize) API() httputil.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.POST("/authorizes", httputil.ScopeRequired(constant.ScopeAuthorizeWrite), s.Create)
		route.GET("/authorizes", httputil.ScopeRequired(constant.ScopeAuthorizeRead), s.List)
		route.GET("/authorizes/:id", httputil.ScopeRequired(constant.ScopeAuthorizeRead), s.Get)
		route.POST("/authorizes/:id", httputil.ScopeRequired(constant.ScopeAuthorizeWrite), s.Update)
		route.DELETE("/authorizes/:id", httputil.ScopeRequired(constant.ScopeAuthorizeWrite), s.Delete)
	}
}
