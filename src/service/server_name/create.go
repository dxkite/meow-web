package server_name

import (
	"context"
	"net/http"

	"dxkite.cn/meownest/pkg/cmd"
)

func NewCreate() cmd.Command {
	return new(Create)
}

type Create struct {
	Name          string `json:"name" form:"name" binding:"required"`
	CertificateId string `json:"certificate_id" form:"certificate_id"`
	Key           string `json:"key" form:"key"`
	Certificate   string `json:"certificate" form:"certificate"`
}

func (req *Create) Execute(ctx context.Context) cmd.CommandResult {
	return cmd.NewResult(http.StatusCreated, req)
}
