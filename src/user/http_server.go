package user

import (
	"context"
	"net/http"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/nebula/pkg/httputil"
	"dxkite.cn/nebula/pkg/httputil/router"
)

func NewUserHttpServer(s UserService, cfg *config.Config) *UserHttpServer {
	return &UserHttpServer{s: s, session: cfg.SessionName}
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users [post]
func (s *UserHttpServer) Create(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param CreateUserRequest

	if err := httputil.ReadRequest(ctx, req, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	if err := httputil.Validate(ctx, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	rst, err := s.s.Create(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusCreated, rst)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users/{id} [get]
func (s *UserHttpServer) Get(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param GetUserRequest
	param.Id = vars["id"]

	if err := httputil.ReadQuery(ctx, req, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	if err := httputil.Validate(ctx, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	rst, err := s.s.Get(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, rst)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users [get]
func (s *UserHttpServer) List(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param ListUserRequest

	if err := httputil.ReadQuery(ctx, req, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	if err := httputil.Validate(ctx, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	rst, err := s.s.List(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, rst)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users/{id} [post]
func (s *UserHttpServer) Update(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param UpdateUserRequest
	param.Id = vars["id"]

	if err := httputil.ReadRequest(ctx, req, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	if err := httputil.Validate(ctx, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	rst, err := s.s.Update(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, rst)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users/{id} [delete]
func (s *UserHttpServer) Delete(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param DeleteUserRequest

	param.Id = vars["id"]

	err := s.s.Delete(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, nil)
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users/session [post]
func (s *UserHttpServer) CreateSession(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param CreateUserSessionRequest

	if err := httputil.ReadRequest(ctx, req, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	if err := httputil.Validate(ctx, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	rst, err := s.s.CreateSession(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: s.session, Value: rst.Token,
		MaxAge:   360,
		Secure:   true,
		HttpOnly: true,
	})
	httputil.Result(ctx, w, http.StatusOK, rst)
}

// Delete User Session
//
// @Summary      Delete User Session
// @Description  Delete User Session
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /users/session [delete]
func (s *UserHttpServer) DeleteSession(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	err := s.s.DeleteSession(ctx, 0)

	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, nil)
}

func (s *UserHttpServer) Routes() []router.Route {
	return []router.Route{
		router.POST("/api/v1/users/session", s.CreateSession),
		router.DELETE("/api/v1/users/session", s.DeleteSession, httputil.ScopeRequired()),

		router.POST("/api/v1/users", s.Create, httputil.ScopeRequired(ScopeUserWrite)),
		router.GET("/api/v1/users", s.List, httputil.ScopeRequired(ScopeUserRead)),
		router.GET("/api/v1/users/:id", s.Get, httputil.ScopeRequired(ScopeUserRead)),
		router.POST("/api/v1/users/:id", s.Update, httputil.ScopeRequired(ScopeUserWrite)),
		router.DELETE("/api/v1/users/:id", s.Delete, httputil.ScopeRequired(ScopeUserWrite)),
	}
}
