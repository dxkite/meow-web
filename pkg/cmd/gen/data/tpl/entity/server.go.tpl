package {{ .ModuleName }}

import (
	"context"
	"net/http"

	"{{ .PackageName }}/pkg/httputil"
	"{{ .PackageName }}/pkg/httputil/router"
)

func New{{ .Name }}Server(s {{ .Name }}Service) *{{ .Name }}Server {
	return &{{ .Name }}Server{s: s}
}

type {{ .Name }}Server struct {
	s {{ .Name }}Service
}

// Create {{ .Name }}
//
// @Summary      Create {{ .Name }}
// @Description  Create {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        body body Create{{ .Name }}Request true "{{ .Name }} data"
// @Success      200  {object} {{ .Name }}Dto
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .BaseURL }} [post]
func (s *{{ .Name }}Server) Create(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param Create{{ .Name }}Request

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

// Get {{ .Name }}
//
// @Summary      Get {{ .Name }}
// @Description  Get {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }} ID"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} {{ .Name }}Dto
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .BaseURL }}/{id} [get]
func (s *{{ .Name }}Server) Get(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param Get{{ .Name }}Request

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

// List {{ .Name }}
//
// @Summary      {{ .Name }} list
// @Description  {{ .Name }} list
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        name query string false "{{ .Name }}"
// @Param        include_total query bool false "include total count"
// @Param        page query int false "page"
// @Param        pre_page query int false "size per page"
// @Param        expand query []string false "expand attribute list"
// @Success      200  {object} List{{ .Name }}Result
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .BaseURL }} [get]
func (s *{{ .Name }}Server) List(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param List{{ .Name }}Request

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

// Update {{ .Name }}
//
// @Summary      Update {{ .Name }}
// @Description  Update {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }} ID"
// @Param        body body Update{{ .Name }}Request true "data"
// @Success      200  {object} {{ .Name }}Dto
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .BaseURL }}/{id} [post]
func (s *{{ .Name }}Server) Update(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param Update{{ .Name }}Request
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

// Delete {{ .Name }}
//
// @Summary      Delete {{ .Name }}
// @Description  Delete {{ .Name }}
// @Tags         {{ .Name }}
// @Accept       json
// @Produce      json
// @Param        id path string true "{{ .Name }}ID"
// @Success      200
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /{{ .BaseURL }}/{id} [delete]
func (s *{{ .Name }}Server) Delete(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param Delete{{ .Name }}Request

	param.Id = vars["id"]

	err := s.s.Delete(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, nil)
}


func (s *{{ .Name }}Server) Routes() []router.Route {
	return []router.Route{
		router.POST("/{{ .BaseURL }}", s.Create),
		router.GET("/{{ .BaseURL }}", s.List),
		router.GET("/{{ .BaseURL }}/:id", s.Get),
		router.POST("/{{ .BaseURL }}/:id", s.Update),
		router.DELETE("/{{ .BaseURL }}/:id", s.Delete),
	}
}