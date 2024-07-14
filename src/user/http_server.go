package user

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/constant"
	"github.com/gin-gonic/gin"
)

func NewUserHttpServer(s UserService, session string) *UserHttpServer {
	return &UserHttpServer{s: s, session: session}
}

type UserHttpServer struct {
	s       UserService
	session string
}

// Create User
//
// @Summary      Create User
// @Description  Create User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body CreateUserRequest true "User data"
// @Success      200  {object} UserDto
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users [post]
func (s *UserHttpServer) Create(c *gin.Context) {
	var param CreateUserRequest

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
// @Success      200  {object} UserDto
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/{id} [get]
func (s *UserHttpServer) Get(c *gin.Context) {
	var param GetUserRequest

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
// @Param		 include_total query bool false "是否包含total"
// @Param        page query int false "页码"
// @Param        pre_page query int false "每页数量"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} ListUserResponse
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users [get]
func (s *UserHttpServer) List(c *gin.Context) {
	var param ListUserRequest

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
// @Param        body body UpdateUserRequest true "data"
// @Success      200  {object} UserDto
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/{id} [post]
func (s *UserHttpServer) Update(c *gin.Context) {
	var param UpdateUserRequest
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
func (s *UserHttpServer) Delete(c *gin.Context) {
	var param DeleteUserRequest

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

// Create User CreateSession
//
// @Summary      Create User CreateSession
// @Description  Create User CreateSession
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body CreateUserSessionRequest true "data"
// @Success      200  {object} CreateSessionResponse
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/session [post]
func (s *UserHttpServer) CreateSession(c *gin.Context) {
	var param CreateUserSessionRequest

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.CreateSession(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	c.SetCookie(s.session, rst.Token, 360, "", "", true, true)

	httpserver.Result(c, http.StatusOK, rst)
}

// Delete User Session
//
// @Summary      Delete User Session
// @Description  Delete User Session
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /users/session [delete]
func (s *UserHttpServer) DeleteSession(c *gin.Context) {

	userId := httpserver.IdentityFrom(c)

	err := s.s.DeleteSession(c, userId)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *UserHttpServer) API() httpserver.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.POST("/users/session", s.CreateSession)
		route.DELETE("/users/session", httpserver.IdentityRequired(), s.DeleteSession)
		route.POST("/users", httpserver.ScopeRequired(constant.ScopeUserWrite), s.Create)
		route.GET("/users", httpserver.ScopeRequired(constant.ScopeUserRead), s.List)

		route.GET("/users/:id", httpserver.ScopeRequired(constant.ScopeUserRead), s.Get)
		route.POST("/users/:id", httpserver.ScopeRequired(constant.ScopeUserWrite), s.Update)
		route.DELETE("/users/:id", httpserver.ScopeRequired(constant.ScopeUserWrite), s.Delete)
	}
}
