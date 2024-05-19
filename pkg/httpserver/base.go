package httpserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RouteHandleFunc func(route gin.IRouter)

type HttpError struct {
	status int
	Error  *HttpErrorDetail `json:"error"`
}

type HttpErrorDetail struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func NewHttpError(status int, details *HttpErrorDetail) *HttpError {
	return &HttpError{status: status, Error: details}
}

func (e *HttpError) Respond(c *gin.Context) {
	c.JSON(e.status, e)
}

func New() *HttpServer {
	g := gin.Default()
	g.ContextWithFallback = true
	s := &HttpServer{engine: g}
	return s
}

type HttpServer struct {
	engine *gin.Engine
}

func (s *HttpServer) Run(addr ...string) {
	s.engine.Run(addr...)
}

func (s *HttpServer) Use(middleware ...gin.HandlerFunc) {
	s.engine.Use(middleware...)
}

func (s *HttpServer) Handle(fn RouteHandleFunc) {
	fn(s.engine)
}

func (s *HttpServer) HandlePrefix(prefix string, fn RouteHandleFunc) {
	fn(s.engine.Group(prefix))
}

func Error(c *gin.Context, status int, code, message string) {
	NewHttpError(status, &HttpErrorDetail{
		Code:    code,
		Message: message,
	}).Respond(c)

}

func Result(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func ResultEmpty(c *gin.Context, status int) {
	c.Status(status)
}

func ResultError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(c, http.StatusNotFound, "not_found", err.Error())
		return
	}
	Error(c, http.StatusInternalServerError, "internal_error", err.Error())
}

func ResultErrorBind(c *gin.Context, err error) {
	if e, ok := err.(validator.ValidationErrors); ok {
		errorList := []string{}
		for _, v := range e {
			customErr := fmt.Sprintf("invalid key %s (%s) by validate rule %s (%s)", v.Field(), v.Namespace(), v.Tag(), v.Param())
			errorList = append(errorList, customErr)
		}
		NewHttpError(http.StatusBadRequest, &HttpErrorDetail{
			Code:    "invalid_parameter",
			Message: "validate_error",
			Details: errorList,
		}).Respond(c)
	} else {
		NewHttpError(http.StatusBadRequest, &HttpErrorDetail{
			Code:    "invalid_parameter",
			Message: err.Error(),
		}).Respond(c)
	}
}
