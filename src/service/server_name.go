package service

import (
	"context"

	"dxkite.cn/meownest/pkg/id"
	"dxkite.cn/meownest/src/model"
	"dxkite.cn/meownest/src/repository"
)

const ServerNamePrefix = "name_"

type ServerName interface {
	Create(ctx context.Context, create *CreateParam) (*CreateResult, error)
	Get(ctx context.Context, id string) (*CreateResult, error)
}

func NewServerName(repo repository.ServerName) ServerName {
	return &serverName{repo: repo}
}

type serverName struct {
	repo repository.ServerName
}

type CreateParam struct {
	Name          string `json:"name" form:"name" binding:"required"`
	Protocol      string `json:"protocol" form:"protocol" binding:"required"`
	CertificateId string `json:"certificate_id" form:"certificate_id"`
}

type CreateResult struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	CertificateId string `json:"certificate_id,omitempty"`
}

func (s *serverName) Create(ctx context.Context, param *CreateParam) (*CreateResult, error) {
	rst, err := s.repo.Create(ctx, &model.ServerName{
		Name:          param.Name,
		Protocol:      param.Protocol,
		CertificateId: param.CertificateId,
	})
	if err != nil {
		return nil, err
	}
	return &CreateResult{
		Id:            id.Format(ServerNamePrefix, rst.Id),
		Name:          param.Name,
		Protocol:      param.Protocol,
		CertificateId: param.CertificateId,
	}, nil
}

func (s *serverName) Get(ctx context.Context, id string) (*CreateResult, error) {
	return nil, nil
}
