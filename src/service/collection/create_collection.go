package collection

import (
	"dxkite.cn/meownest/src/dto"
	"github.com/gin-gonic/gin"
)

type CreateCollection struct {
	dto.Collection
}

func (req *CreateCollection) Handle(c *gin.Context) (*dto.Collection, error) {
	return &dto.Collection{}, nil
}
