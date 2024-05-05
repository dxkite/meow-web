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

// Create Collection
//
// @Summary      Create Collection
// @Description  Create Collection
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        body body service.CreateCollectionParam true "Collection data"
// @Success      200  {object} dto.Collection
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /collections [post]
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

// Get Collection
//
// @Summary      Get Collection
// @Description  Get Collection
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        id path string true "Collection ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} dto.Collection
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /collections/{id} [get]
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

// Collection列表
//
// @Summary      Collection列表
// @Description  Collection列表
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        parent_id query string false "父级ID"
// @Param		 depth query int false "获取深度"
// @Param        name query string false "Collection"
// @Param        limit query int false "限制"
// @Param        starting_after query string false "从当前ID开始"
// @Param        ending_before query string false "从当前ID结束"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListCollectionResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /collections [get]
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

// Update Collection
//
// @Summary      Update Collection
// @Description  Update Collection
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        id path string true "Collection ID"
// @Param        body body service.UpdateCollectionParam true "data"
// @Success      200  {object} service.Collection
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /collections/{id} [post]
func (s *Collection) Update(c *gin.Context) {
	var param service.UpdateCollectionParam
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

// Delete Collection
//
// @Summary      Delete Collection
// @Description  Delete Collection
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        id path string true "CollectionID"
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /collections/{id} [delete]
func (s *Collection) Delete(c *gin.Context) {
	var param service.DeleteCollectionParam

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

func (s *Collection) RegisterToHttp(route gin.IRouter) {
	route.POST("/collections", s.Create)
	route.GET("/collections", s.List)
	route.GET("/collections/:id", s.Get)
	route.POST("/collections/:id", s.Update)
	route.DELETE("/collections/:id", s.Delete)
}
