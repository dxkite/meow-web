package v1

import (
	"net/http"

	"dxkite.cn/meownest/src/service/collection"
	"github.com/gin-gonic/gin"
)

func CreateCollection(c *gin.Context) {
	req := new(collection.CreateCollection)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_parameter",
			"message": err.Error(),
		})
		return
	}

	resp, err := req.Handle(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    "process_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
