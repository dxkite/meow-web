package server

import "github.com/gin-gonic/gin"

type WithOption func(s *HttpServer)

func New(opts ...WithOption) *HttpServer {
	s := &HttpServer{engine: gin.Default()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type HttpServer struct {
	engine *gin.Engine
}

func (s *HttpServer) Run(addr ...string) {
	s.engine.Run(addr...)
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

func Result(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
