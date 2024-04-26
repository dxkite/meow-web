package server

import "github.com/gin-gonic/gin"

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
