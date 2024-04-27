package service

import (
	"context"

	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
)

type Certificate interface {
	Create(ctx context.Context, create *CreateCertificateParam) (*dto.Certificate, error)
}

type CreateCertificateParam struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Key         string `json:"key" form:"key" binding:"required"`
	Certificate string `json:"certificate" form:"key" binding:"required"`
}

func NewCertificate(r repository.Certificate) Certificate {
	return &certificate{r: r}
}

type certificate struct {
	r repository.Certificate
}

func (s *certificate) Create(ctx context.Context, param *CreateCertificateParam) (*dto.Certificate, error) {
	entity, err := entity.NewCertificateWithCertificateKey(param.Certificate, param.Key)
	if err != nil {
		return nil, err
	}

	resp, err := s.r.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return dto.NewCertificate(resp), nil
}
