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
	List(ctx context.Context, param *ListServerNameParam) (*ListServerNameResult, error)
}

func NewServerName(r repository.ServerName, rc repository.Certificate, db *gorm.DB) ServerName {
	return &serverName{r: r, rc: rc, db: db}
}

type serverName struct {
	r  repository.ServerName
	rc repository.Certificate
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

		var certificateId = identity.Parse(constant.CertificatePrefix, param.CertificateId)
		var certificate *dto.Certificate

		if param.Certificate != nil {
			certEntity, err := entity.NewCertificateWithCertificateKey(param.Certificate.Certificate, param.Certificate.Key)
			if err != nil {
				return err
			}

			certEntity.Name = param.Name
			if cert, err := s.rc.Create(ctx, certEntity); err != nil {
				return err
			} else {
				certificateId = cert.Id
				certificate = dto.NewCertificate(certEntity)
			}
		}

		entity, err := s.r.Create(ctx, &entity.ServerName{
			Name:          param.Name,
			Protocol:      param.Protocol,
			CertificateId: certificateId,
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

type ListServerNameParam struct {
	Name          string   `form:"name"`
	Limit         int      `form:"limit" binding:"max=1000"`
	StartingAfter string   `form:"starting_after"`
	EndingBefore  string   `form:"ending_before"`
	Expand        []string `json:"expand" form:"expand"`
}

type ListServerNameResult struct {
	HasMore bool              `json:"has_more"`
	Data    []*dto.ServerName `json:"data"`
}

func (s *serverName) List(ctx context.Context, param *ListServerNameParam) (*ListServerNameResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListServerNameParam{
		Name:          param.Name,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.CollectionPrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.CollectionPrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.ServerName, n)

	for i, v := range entities {
		items[i] = dto.NewServerName(v)
	}

	if utils.InStringSlice("certificate", param.Expand) {
		err := utils.ExpandStruct(
			n,
			func(i int) ([]uint64, error) {
				return []uint64{entities[i].CertificateId}, nil
			},
			func(i int, v []interface{}) error {
				if len(v) > 0 {
					if vv, ok := v[0].(*entity.Certificate); ok {
						items[i].Certificate = dto.NewCertificate(vv)
					}
				}
				return nil
			},
			func(ids []uint64) (map[uint64]interface{}, error) {
				v := map[uint64]interface{}{}
				entities, err := s.rc.BatchGet(ctx, ids)
				if err != nil {
					return nil, err
				}
				for _, e := range entities {
					v[e.Id] = e
				}
				return v, nil
			},
		)
		if err != nil {
			return nil, err
		}
	}

	rst := &ListServerNameResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

func (r *serverName) dataSource(ctx context.Context) *gorm.DB {
	return repository.DataSource(ctx, r.db)
}
