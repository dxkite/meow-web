package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/utils"
	"gorm.io/gorm"
)

type GetServerNameParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

type ServerName interface {
	Create(ctx context.Context, param *CreateServerNameParam) (*dto.ServerName, error)
	Get(ctx context.Context, param *GetServerNameParam) (*dto.ServerName, error)
}

func NewServerName(r repository.ServerName, rc repository.Certificate, sc Certificate, db *gorm.DB) ServerName {
	return &serverName{r: r, rc: rc, sc: sc, db: db}
}

type serverName struct {
	r  repository.ServerName
	rc repository.Certificate
	sc Certificate
	db *gorm.DB
}

type CreateServerNameParam struct {
	Name          string                            `json:"name" form:"name" binding:"required"`
	Protocol      string                            `json:"protocol" form:"protocol" binding:"required"`
	CertificateId string                            `json:"certificate_id" form:"certificate_id"`
	Certificate   *CreateServerNameCertificateParam `json:"certificate" form:"certificate"`
}

type CreateServerNameCertificateParam struct {
	Key         string `json:"key" form:"key" binding:"required"`
	Certificate string `json:"certificate" form:"key" binding:"required"`
}

// 创建域名
// 支持联动创建证书
func (s *serverName) Create(ctx context.Context, param *CreateServerNameParam) (*dto.ServerName, error) {
	var name *dto.ServerName

	err := s.dataSource(ctx).Transaction(func(tx *gorm.DB) error {
		ctx := repository.WithDataSource(ctx, tx)

		var certificateId = param.CertificateId
		var certificate *dto.Certificate

		if param.Certificate != nil {
			if cert, err := s.sc.Create(ctx, &CreateCertificateParam{
				Name:        param.Name,
				Key:         param.Certificate.Key,
				Certificate: param.Certificate.Certificate,
			}); err != nil {
				return err
			} else {
				certificateId = cert.Id
				certificate = cert
			}
		}

		entity, err := s.r.Create(ctx, &entity.ServerName{
			Name:          param.Name,
			Protocol:      param.Protocol,
			CertificateId: identity.Parse(constant.CertificatePrefix, certificateId),
		})
		if err != nil {
			return err
		}

		name = dto.NewServerName(entity)
		name.Certificate = certificate
		return nil
	})

	return name, err
}

func (s *serverName) Get(ctx context.Context, param *GetServerNameParam) (*dto.ServerName, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.ServerNamePrefix, param.Id))
	if err != nil {
		return nil, err
	}
	name := dto.NewServerName(rst)

	if utils.InStringSlice("certificate", param.Expand) {
		cert, err := s.rc.Get(ctx, rst.CertificateId)
		if err != nil {
			return nil, err
		}
		name.Certificate = dto.NewCertificate(cert)
	}

	return name, nil
}

func (r *serverName) dataSource(ctx context.Context) *gorm.DB {
	return repository.DataSource(ctx, r.db)
}
