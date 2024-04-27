package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"

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
	cert, err := tls.X509KeyPair([]byte(param.Certificate), []byte(param.Key))
	if err != nil {
		return nil, err
	}

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	val := &entity.Certificate{}
	val.Name = param.Name
	val.Key = param.Key
	val.Certificate = param.Certificate
	val.StartTime = leaf.NotBefore
	val.EndTime = leaf.NotAfter
	val.Domain = leaf.DNSNames

	resp, err := s.r.Create(ctx, val)
	if err != nil {
		return nil, err
	}

	return dto.NewCertificate(resp), nil
}
