package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewUser(s service.User) *User {
	return &User{s: s}
}

type User struct {
	s service.User
}

// Create User
//
// @Summary      Create User
// @Description  Create User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body service.CreateUserParam true "User data"
// @Success      200  {object} dto.User
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users [post]
func (s *User) Create(c *gin.Context) {
	var param service.CreateUserParam

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

// Get User
//
// @Summary      Get User
// @Description  Get User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} dto.User
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/{id} [get]
func (s *User) Get(c *gin.Context) {
	var param service.GetUserParam

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

// List User
//
// @Summary      User list
// @Description  User list
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        name query string false "User"
// @Param        limit query int false "size limit"
// @Param        starting_after query string false "get list after id"
// @Param        ending_before query string false "get list before id"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} service.ListUserResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users [get]
func (s *User) List(c *gin.Context) {
	var param service.ListUserParam

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

// Update User
//
// @Summary      Update User
// @Description  Update User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        body body service.UpdateUserParam true "data"
// @Success      200  {object} dto.User
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/{id} [post]
func (s *User) Update(c *gin.Context) {
	var param service.UpdateUserParam
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

// Delete User
//
// @Summary      Delete User
// @Description  Delete User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "UserID"
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/{id} [delete]
func (s *User) Delete(c *gin.Context) {
	var param service.DeleteUserParam

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

// Create User Session
//
// @Summary      Create User Session
// @Description  Create User Session
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body service.CreateUserSessionParam true "data"
// @Success      200  {object} service.CreateSessionResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/session [post]
func (s *User) Session(c *gin.Context) {
	var param service.CreateUserSessionParam

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Session(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.Result(c, http.StatusOK, rst)
}

func (s *User) RegisterToHttp(route gin.IRouter) {
	route.POST("/users", s.Create)
	route.GET("/users", s.List)
	route.GET("/users/:id", s.Get)
	route.POST("/users/:id", s.Update)
	route.DELETE("/users/:id", s.Delete)
	route.POST("/users/session", s.Session)
}
