package main

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"dxkite.cn/log"
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
	initLogger()
	initBinding()
	identity.DefaultMask = 1234627081864056831
}

func initLogger() {
	log.SetOutput(log.NewColorWriter(os.Stdout))
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func initBinding() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			log.Error("[panic error]", r)
			log.Error(string(buf[:n]))
			name := fmt.Sprintf("crash-%s.log", time.Now().Format("20060102150405"))
			panicErr := string(buf[:n])
			_ = os.WriteFile(name, []byte(panicErr), os.ModePerm)
		}
	}()

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
	serverNameService := service.NewServerName(nameServerRepository, certificateRepository, certificateService, db)
	serverNameServer := server.NewServerName(serverNameService)

	routeRepository := repository.NewRoute(db)
	routeService := service.NewRoute(routeRepository)
	routeServer := server.NewRoute(routeService)

	endpointRepository := repository.NewEndpoint(db)
	endpointService := service.NewEndpoint(endpointRepository)
	endpointServer := server.NewEndpoint(endpointService)

	collectionRepository := repository.NewCollection(db)
	collectionService := service.NewCollection(collectionRepository, linkRepository, routeRepository, endpointRepository)
	collectionServer := server.NewCollection(collectionService)

	httpServer := server.New(
		server.WithServerName("/api/v1/server_name", serverNameServer),
		server.WithCertificate("/api/v1/certificate", certificateServer),
		server.WithCollection("/api/v1/collection", collectionServer),
		server.WithRoute("/api/v1/route", routeServer),
		server.WithEndpoint("/api/v1/endpoint", endpointServer),
	)
	httpServer.Run(":2333")
}
