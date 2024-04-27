package server

import (
	"net/http"

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
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	Result(c, http.StatusCreated, rst)
}

func (s *Collection) Get(c *gin.Context) {
	var param service.GetCollectionParam

	param.Id = c.Param("id")

	if err := c.ShouldBindQuery(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	Result(c, http.StatusOK, rst)
}

func (s *Collection) LinkRoute(c *gin.Context) {
	var param service.LinkRouteParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	err := s.s.LinkRoute(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func (s *Collection) DeleteRoute(c *gin.Context) {
	var param service.DeleteRouteParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	err := s.s.DeleteRoute(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func (s *Collection) LinkEndpoint(c *gin.Context) {
	var param service.LinkEndpointParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	err := s.s.LinkEndpoint(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func (s *Collection) DeleteEndpoint(c *gin.Context) {
	var param service.DeleteEndpointParam

	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		Error(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	err := s.s.DeleteEndpoint(c, &param)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func WithCollection(path string, server *Collection) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
			group.GET("/:id", server.Get)
			group.POST("/:id/route", server.LinkRoute)
			group.DELETE("/:id/route", server.DeleteRoute)
			group.POST("/:id/endpoint", server.LinkEndpoint)
			group.DELETE("/:id/endpoint", server.DeleteEndpoint)
		}
	}
}
