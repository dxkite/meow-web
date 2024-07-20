package server

import (
	"dxkite.cn/meownest/docs"
	"dxkite.cn/meownest/pkg/httputil"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Swagger struct {
}

func NewSwagger() *Swagger {
	return &Swagger{}
}

func (s *Swagger) API() httputil.RouteHandleFunc {
	return func(route gin.IRouter) {
		docs.SwaggerInfo.Title = "MeowNest Admin API"
		docs.SwaggerInfo.Description = "This is a sample server meow nest server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "192.168.1.105:2333"
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http"}
		route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
