package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewServerName(s service.ServerName) *ServerName {
	return &ServerName{s: s}
}

type ServerName struct {
	s service.ServerName
}

// 创建域名
//
// @Summary      创建域名
// @Description  创建域名
// @Tags         域名
// @Accept       json
// @Produce      json
// @Param        body body service.CreateServerNameParam true "请求体"
// @Success      201  {object} dto.ServerName
// @Failure      400  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /server_name [post]
func (s *ServerName) Create(c *gin.Context) {
	var param service.CreateServerNameParam

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

// 获取域名
//
// @Summary      获取域名
// @Description  获取域名
// @Tags         域名
// @Accept       json
// @Produce      json
// @Param        id path string true "域名ID"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} dto.ServerName
// @Failure      400  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /server_name/{id} [get]
func (s *ServerName) Get(c *gin.Context) {
	var param service.GetServerNameParam

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

// 域名列表
//
// @Summary      域名列表
// @Description  域名列表
// @Tags         域名
// @Accept       json
// @Produce      json
// @Param        name query string false "域名"
// @Param        limit query int false "限制"
// @Param        starting_after query string false "从当前ID开始"
// @Param        ending_before query string false "从当前ID结束"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListServerNameserver.Result
// @Failure      400  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /server_name [get]
func (s *ServerName) List(c *gin.Context) {
	var param service.ListServerNameParam

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

// 更新域名
//
// @Summary      更新域名
// @Description  更新域名
// @Tags         域名
// @Accept       json
// @Produce      json
// @Param        id path string true "域名ID"
// @Param        body body service.UpdateServerNameParam true "数据"
// @Success      200  {object} service.ServerName
// @Failure      400  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /server_name/{id} [post]
func (s *ServerName) Update(c *gin.Context) {
	var param service.UpdateServerNameParam
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

// 删除域名
//
// @Summary      删除域名
// @Description  删除域名
// @Tags         域名
// @Accept       json
// @Produce      json
// @Param        id path string true "域名ID"
// @Success      200
// @Failure      400  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /server_name/{id} [delete]
func (s *ServerName) Delete(c *gin.Context) {
	var param service.DeleteServerNameParam

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

func (s *ServerName) RegisterToHttp(group gin.IRouter) {
	group.GET("/server_name", s.List)
	group.POST("/server_name", s.Create)
	group.GET("/server_name/:id", s.Get)
	group.DELETE("/server_name/:id", s.Delete)
	group.POST("/server_name/:id", s.Update)
}
