package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/model"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/utils"
)

type CreateServerNameParam struct {
	Name          string `json:"name" form:"name" binding:"required"`
	Protocol      string `json:"protocol" form:"protocol" binding:"required"`
	CertificateId string `json:"certificate_id" form:"certificate_id"`
}

type ServerName interface {
	Create(ctx context.Context, create *CreateServerNameParam) (*dto.ServerName, error)
	Get(ctx context.Context, id string, expand []string) (*dto.ServerName, error)
}

func NewServerName(r repository.ServerName, rc repository.Certificate) ServerName {
	return &serverName{r: r, rc: rc}
}

type serverName struct {
	r  repository.ServerName
	rc repository.Certificate
}

func (s *serverName) Create(ctx context.Context, param *CreateServerNameParam) (*dto.ServerName, error) {
	rst, err := s.r.Create(ctx, &model.ServerName{
		Name:          param.Name,
		Protocol:      param.Protocol,
		CertificateId: identity.Parse(constant.CertificatePrefix, param.CertificateId),
	})
	if err != nil {
		return nil, err
	}

	name := &dto.ServerName{
		Id:          identity.Format(constant.ServerNamePrefix, rst.Id),
		Name:        param.Name,
		Protocol:    param.Protocol,
		Certificate: &dto.Certificate{Id: param.CertificateId},
	}
	return name, nil
}

func (s *serverName) Get(ctx context.Context, id string, expand []string) (*dto.ServerName, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.ServerNamePrefix, id))
	if err != nil {
		return nil, err
	}
	name := dto.NewServerName(rst)

	if utils.InStringSlice("certificate", expand) {
		cert, err := s.rc.Get(ctx, rst.CertificateId)
		if err != nil {
			return nil, err
		}
		name.Certificate = dto.NewCertificate(cert)
	}

	return name, nil
}
