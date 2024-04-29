package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewCollection(s service.Collection) *Collection {
	return &Collection{s: s}
}

type Collection struct {
	s service.Collection
}

func (s *Collection) Create(c *gin.Context) {
	var param service.CreateCollectionParam

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

func (s *Collection) Get(c *gin.Context) {
	var param service.GetCollectionParam

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

func (s *Collection) List(c *gin.Context) {
	var param service.ListCollectionParam

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

func (s *Collection) LinkRoute(c *gin.Context) {
	var param service.LinkCollectionRouteParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	err := s.s.LinkRoute(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *Collection) DeleteRoute(c *gin.Context) {
	var param service.DeleteCollectionRouteParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	err := s.s.DeleteRoute(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *Collection) LinkEndpoint(c *gin.Context) {
	var param service.LinkCollectionEndpointParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	err := s.s.LinkEndpoint(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *Collection) DeleteEndpoint(c *gin.Context) {
	var param service.DeleteCollectionEndpointParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	err := s.s.DeleteEndpoint(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *Collection) RegisterToHttp(group gin.IRouter) {
	group.GET("/collection", s.List)
	group.POST("/collection", s.Create)
	group.GET("/collection/:id", s.Get)
	group.POST("/collection/:id/route", s.LinkRoute)
	group.DELETE("/collection/:id/route", s.DeleteRoute)
	group.POST("/collection/:id/endpoint", s.LinkEndpoint)
	group.DELETE("/collection/:id/endpoint", s.DeleteEndpoint)
}
