package router

import (
	"dxkite.cn/meownest/pkg/cmd"

	"dxkite.cn/meownest/src/service/server_name"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")

	// collections := api.Group("/collections")
	// {
	// 	collections.POST("")
	// }

	serverNames := api.Group("/server_names")
	{
		serverNames.GET("", cmd.Exec(server_name.NewCreate))
		serverNames.POST("", cmd.Exec(server_name.NewCreate))
	}

	return r
}
