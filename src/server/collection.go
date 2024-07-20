package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/src/constant"
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /collections [post]
func (s *Collection) Create(c *gin.Context) {
	var param service.CreateCollectionParam

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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /collections/{id} [get]
func (s *Collection) Get(c *gin.Context) {
	var param service.GetCollectionParam

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
// @Param		 include_total query bool false "是否包含total"
// @Param        page query int false "页码"
// @Param        pre_page query int false "每页数量"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListCollectionResult
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /collections [get]
func (s *Collection) List(c *gin.Context) {
	var param service.ListCollectionParam

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

// Update Collection
//
// @Summary      Update Collection
// @Description  Update Collection
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        id path string true "Collection ID"
// @Param        body body service.UpdateCollectionParam true "data"
// @Success      200  {object} dto.Collection
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /collections/{id} [post]
func (s *Collection) Update(c *gin.Context) {
	var param service.UpdateCollectionParam
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

// Delete Collection
//
// @Summary      Delete Collection
// @Description  Delete Collection
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        id path string true "CollectionID"
// @Success      200
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /collections/{id} [delete]
func (s *Collection) Delete(c *gin.Context) {
	var param service.DeleteCollectionParam

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

func (s *Collection) API() httputil.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.POST("/collections", httputil.ScopeRequired(constant.ScopeCollectionWrite), s.Create)
		route.GET("/collections", httputil.ScopeRequired(constant.ScopeCollectionRead), s.List)
		route.GET("/collections/:id", httputil.ScopeRequired(constant.ScopeCollectionRead), s.Get)
		route.POST("/collections/:id", httputil.ScopeRequired(constant.ScopeCollectionWrite), s.Update)
		route.DELETE("/collections/:id", httputil.ScopeRequired(constant.ScopeCollectionWrite), s.Delete)
	}
}
