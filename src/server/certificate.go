package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
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
// @Param        id path string true "证书ID"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} dto.Certificate
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates [post]
func (s *Certificate) Create(c *gin.Context) {
	var param service.CreateCertificateParam

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
// @Success      200  {object} service.ListCertificateserver.Result
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates [get]
func (s *Certificate) List(c *gin.Context) {
	var param service.ListCertificateParam

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
// @Param        body body service.UpdateCertificateParam true "数据"
// @Success      200  {object} service.Certificate
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /certificates/{id} [post]
func (s *Certificate) Update(c *gin.Context) {
	var param service.UpdateCertificateParam
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
// @Router       /certificates/{id} [delete]
func (s *Certificate) Delete(c *gin.Context) {
	var param service.DeleteCertificateParam

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

func (s *Certificate) RegisterToHttp(group gin.IRouter) {
	group.POST("/certificates", s.Create)
	group.GET("/certificates", s.List)
	group.GET("/certificates/:id", s.Get)
	group.POST("/certificates/:id", s.Update)
	group.DELETE("/certificates/:id", s.Delete)

}
