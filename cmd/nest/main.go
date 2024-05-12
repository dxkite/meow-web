package main

import (
	"context"
	"reflect"
	"strings"
	"time"

	"dxkite.cn/meownest/pkg/agent"
	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/pkg/database/sqlite"
	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/server"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func init() {
	initBinding()
	identity.DefaultMask = 1234627081864056831
}

func initBinding() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			if name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]; name != "" && name != "-" {
				return name
			}
			if name := strings.SplitN(field.Tag.Get("form"), ",", 2)[0]; name != "" && name != "-" {
				return name
			}
			return ""
		})
	}
}

func main() {
	ds, err := sqlite.Open("data.db")
	if err != nil {
		panic(err)
	}

	db := ds.Engine().(*gorm.DB)
	db.AutoMigrate(entity.Certificate{}, entity.User{},
		entity.Collection{}, entity.Route{}, entity.Endpoint{}, entity.Authorize{})

	certificateRepository := repository.NewCertificate()
	certificateService := service.NewCertificate(certificateRepository)
	certificateServer := server.NewCertificate(certificateService)

	userRepository := repository.NewUser()
	userService := service.NewUser(userRepository, []byte("12345678901234567890123456789012"))
	userServer := server.NewUser(userService)

	authorizeRepository := repository.NewAuthorize()
	authorizeService := service.NewAuthorize(authorizeRepository)
	authorizeServer := server.NewAuthorize(authorizeService)

	endpointRepository := repository.NewEndpoint()
	endpointService := service.NewEndpoint(endpointRepository)
	endpointServer := server.NewEndpoint(endpointService)

	routeRepository := repository.NewRoute()
	collectionRepository := repository.NewCollection()
	collectionService := service.NewCollection(
		collectionRepository, routeRepository,
		endpointRepository, authorizeRepository,
	)

	routeService := service.NewRoute(
		routeRepository, endpointRepository,
		collectionRepository, authorizeRepository,
	)
	routeServer := server.NewRoute(routeService)

	collectionServer := server.NewCollection(collectionService)

	ag := agent.New()
	agentService := service.NewAgent(ag,
		routeRepository, collectionRepository,
		endpointRepository, authorizeRepository,
	)
	agentServer := server.NewAgent(agentService)

	httpServer := httpserver.New()

	httpServer.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	httpServer.Use(func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(database.With(ctx.Request.Context(), ds))
	})

	httpServer.RegisterPrefix("/api/v1", certificateServer)
	httpServer.RegisterPrefix("/api/v1", userServer)
	httpServer.RegisterPrefix("/api/v1", routeServer)
	httpServer.RegisterPrefix("/api/v1", endpointServer)
	httpServer.RegisterPrefix("/api/v1", authorizeServer)
	httpServer.RegisterPrefix("/api/v1", collectionServer)
	httpServer.RegisterPrefix("/api/v1", agentServer)
	httpServer.Register(server.NewSwagger())

	go httpServer.Run(":2333")

	agentService.LoadRoute(database.With(context.Background(), ds))
	agentService.Run(":80")
}
