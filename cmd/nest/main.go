package main

import (
	"reflect"
	"strings"
	"time"

	"dxkite.cn/meownest/pkg/agent"
	"dxkite.cn/meownest/pkg/data_source"
	"dxkite.cn/meownest/pkg/data_source/sqlite"
	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/server"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	db := ds.Gorm()
	db.AutoMigrate(entity.ServerName{}, entity.Certificate{},
		entity.Link{},
		entity.Collection{}, entity.Route{}, entity.Endpoint{}, entity.Authorize{})

	certificateRepository := repository.NewCertificate()
	certificateService := service.NewCertificate(certificateRepository)
	certificateServer := server.NewCertificate(certificateService)

	authorizeRepository := repository.NewAuthorize()
	authorizeService := service.NewAuthorize(authorizeRepository)
	authorizeServer := server.NewAuthorize(authorizeService)

	nameServerRepository := repository.NewServerName()
	serverNameService := service.NewServerName(nameServerRepository, certificateRepository)
	serverNameServer := server.NewServerName(serverNameService)

	endpointRepository := repository.NewEndpoint()
	endpointService := service.NewEndpoint(endpointRepository)
	endpointServer := server.NewEndpoint(endpointService)

	linkRepository := repository.NewLink()

	routeRepository := repository.NewRoute()
	collectionRepository := repository.NewCollection()
	collectionService := service.NewCollection(collectionRepository, linkRepository, routeRepository, endpointRepository, nameServerRepository)

	routeService := service.NewRoute(routeRepository, linkRepository, endpointRepository, collectionRepository)
	routeServer := server.NewRoute(routeService)

	collectionServer := server.NewCollection(collectionService)

	agentServer := server.NewAgent(agent.New())

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

	httpServer.Use(data_source.GinDataSource(ds))

	httpServer.RegisterPrefix("/api/v1", certificateServer)
	httpServer.RegisterPrefix("/api/v1", serverNameServer)
	httpServer.RegisterPrefix("/api/v1", routeServer)
	httpServer.RegisterPrefix("/api/v1", endpointServer)
	httpServer.RegisterPrefix("/api/v1", authorizeServer)
	httpServer.RegisterPrefix("/api/v1", collectionServer)
	httpServer.RegisterPrefix("/api/v1", agentServer)
	httpServer.Register(server.NewSwagger())

	go httpServer.Run(":2333")
	agentServer.Run(":80")
}
