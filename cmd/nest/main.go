package main

import (
	"reflect"
	"strings"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/server"
	"dxkite.cn/meownest/src/service"
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

	linkRepository := repository.NewLink(db)

	certificateRepository := repository.NewCertificate(db)
	certificateService := service.NewCertificate(certificateRepository)
	certificateServer := server.NewCertificate(certificateService)

	nameServerRepository := repository.NewServerName(db)
	serverNameService := service.NewServerName(nameServerRepository, certificateRepository, db)
	serverNameServer := server.NewServerName(serverNameService)

	routeRepository := repository.NewRoute(db)
	routeService := service.NewRoute(routeRepository, linkRepository, db)
	routeServer := server.NewRoute(routeService)

	endpointRepository := repository.NewEndpoint(db)
	endpointService := service.NewEndpoint(endpointRepository)
	endpointServer := server.NewEndpoint(endpointService)

	collectionRepository := repository.NewCollection(db)
	collectionService := service.NewCollection(collectionRepository, linkRepository, routeRepository, endpointRepository)
	collectionServer := server.NewCollection(collectionService)

	httpServer := httpserver.New()
	httpServer.RegisterPrefix("/api/v1", certificateServer)
	httpServer.RegisterPrefix("/api/v1", serverNameServer)
	httpServer.RegisterPrefix("/api/v1", routeServer)
	httpServer.RegisterPrefix("/api/v1", endpointServer)
	httpServer.RegisterPrefix("/api/v1", collectionServer)
	httpServer.Register(server.NewSwagger())
	httpServer.Run(":2333")
}
