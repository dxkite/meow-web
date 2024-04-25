package cmd

import "github.com/gin-gonic/gin"

type CommandResult interface {
	Write(c *gin.Context)
}

func NewErrorResult(status int, code, message string) CommandResult {
	return NewResult(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

type result struct {
	status int
	data   interface{}
}

func (e *result) Write(c *gin.Context) {
	c.JSON(e.status, e.data)
}

func NewResult(status int, data interface{}) CommandResult {
	return &result{status: status, data: data}
}
