package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewCertificate(s service.Certificate) *Certificate {
	return &Certificate{s: s}
}

type Certificate struct {
	s service.Certificate
}

// 创建证书
//
// @Summary      创建证书
// @Description  创建证书
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        body body service.CreateCertificateParam true "数据"
// @Success      200  {object} dto.Certificate
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates [post]
func (s *Certificate) Create(c *gin.Context) {
	var param service.CreateCertificateParam

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

// 获取证书
//
// @Summary      获取证书
// @Description  获取证书
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        id path string true "证书ID"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} dto.Certificate
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates/{id} [get]
func (s *Certificate) Get(c *gin.Context) {
	var param service.GetCertificateParam

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

// 证书列表
//
// @Summary      证书列表
// @Description  证书列表
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        name query string false "证书"
// @Param		 include_total query bool false "是否包含total"
// @Param        page query int false "页码"
// @Param        pre_page query int false "每页数量"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListCertificateResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates [get]
func (s *Certificate) List(c *gin.Context) {
	var param service.ListCertificateParam

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

// 更新证书
//
// @Summary      更新证书
// @Description  更新证书
// @Tags         证书
// @Accept       json
// @Produce      json
// @Param        id path string true "证书ID"
// @Param        body body service.UpdateCertificateParam true "数据"
// @Success      200  {object} dto.Authorize
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates/{id} [post]
func (s *Certificate) Update(c *gin.Context) {
	var param service.UpdateCertificateParam
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
// @Router       /certificates/{id} [delete]
func (s *Certificate) Delete(c *gin.Context) {
	var param service.DeleteCertificateParam

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

func (s *Certificate) API() httputil.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.POST("/certificates", httputil.ScopeRequired(constant.ScopeCertificateWrite), s.Create)
		route.GET("/certificates", httputil.ScopeRequired(constant.ScopeAuthorizeRead), s.List)
		route.GET("/certificates/:id", httputil.ScopeRequired(constant.ScopeAuthorizeRead), s.Get)
		route.POST("/certificates/:id", httputil.ScopeRequired(constant.ScopeCertificateWrite), s.Update)
		route.DELETE("/certificates/:id", httputil.ScopeRequired(constant.ScopeCertificateWrite), s.Delete)
	}
}
