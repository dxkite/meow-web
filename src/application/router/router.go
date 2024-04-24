package router

import (
	v1 "dxkite.cn/meownest/src/application/router/v1"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	apiV1 := r.Group("/api/v1")

	collections := apiV1.Group("/collections")
	{
		collections.POST("", v1.CreateCollection)
	}

	ports := apiV1.Group("/ports")
	{
		ports.GET("", v1.CreateCollection)
	}

	return r
}
