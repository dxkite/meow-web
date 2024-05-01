package main

import (
	"reflect"
	"strings"
	"time"

	"dxkite.cn/meownest/pkg/datasource"
	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/server"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin/binding"
	"github.com/glebarez/sqlite"
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
	db, err := gorm.Open(sqlite.Open("data.db"))
	if err != nil {
		panic(err)
	}

	db = db.Debug()
	db.AutoMigrate(entity.ServerName{}, entity.Certificate{},
		entity.Link{},
		entity.Collection{}, entity.Route{}, entity.Endpoint{})

	linkRepository := repository.NewLink()

	certificateRepository := repository.NewCertificate()
	certificateService := service.NewCertificate(certificateRepository)
	certificateServer := server.NewCertificate(certificateService)

	nameServerRepository := repository.NewServerName()
	serverNameService := service.NewServerName(nameServerRepository, certificateRepository)
	serverNameServer := server.NewServerName(serverNameService)

	routeRepository := repository.NewRoute()
	routeService := service.NewRoute(routeRepository, linkRepository)
	routeServer := server.NewRoute(routeService)

	endpointRepository := repository.NewEndpoint()
	endpointService := service.NewEndpoint(endpointRepository)
	endpointServer := server.NewEndpoint(endpointService)

	collectionRepository := repository.NewCollection()
	collectionService := service.NewCollection(collectionRepository, linkRepository, routeRepository, endpointRepository)
	collectionServer := server.NewCollection(collectionService)

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

	httpServer.Use(datasource.RegisterToGin(datasource.New(db)))

	httpServer.RegisterPrefix("/api/v1", certificateServer)
	httpServer.RegisterPrefix("/api/v1", serverNameServer)
	httpServer.RegisterPrefix("/api/v1", routeServer)
	httpServer.RegisterPrefix("/api/v1", endpointServer)
	httpServer.RegisterPrefix("/api/v1", collectionServer)
	httpServer.Register(server.NewSwagger())
	httpServer.Run(":2333")
}
