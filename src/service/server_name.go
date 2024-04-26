package service

import (
	"context"

	"dxkite.cn/meownest/pkg/id"
	"dxkite.cn/meownest/src/model"
	"dxkite.cn/meownest/src/repository"
)

const ServerNamePrefix = "name_"

type CreateServerNameParam struct {
	Name          string `json:"name" form:"name" binding:"required"`
	Protocol      string `json:"protocol" form:"protocol" binding:"required"`
	CertificateId string `json:"certificate_id" form:"certificate_id"`
}

type CreateServerNameResult struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	CertificateId string `json:"certificate_id,omitempty"`
}

type ServerName interface {
	Create(ctx context.Context, create *CreateServerNameParam) (*CreateServerNameResult, error)
	Get(ctx context.Context, id string) (*CreateServerNameResult, error)
}

func NewServerName(r repository.ServerName) ServerName {
	return &serverName{r: r}
}

type serverName struct {
	r repository.ServerName
}

func (s *serverName) Create(ctx context.Context, param *CreateServerNameParam) (*CreateServerNameResult, error) {
	rst, err := s.r.Create(ctx, &model.ServerName{
		Name:          param.Name,
		Protocol:      param.Protocol,
		CertificateId: param.CertificateId,
	})
	if err != nil {
		return nil, err
	}
	return &CreateServerNameResult{
		Id:            id.Format(ServerNamePrefix, rst.Id),
		Name:          param.Name,
		Protocol:      param.Protocol,
		CertificateId: param.CertificateId,
	}, nil
}

func (s *serverName) Get(ctx context.Context, id string) (*CreateServerNameResult, error) {
	return nil, nil
}
