package service

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/utils"
)

type Collection interface {
	Create(ctx context.Context, create *CreateCollectionParam) (*dto.Collection, error)
	Update(ctx context.Context, param *UpdateCollectionParam) (*dto.Collection, error)
	Get(ctx context.Context, param *GetCollectionParam) (*dto.Collection, error)
	Delete(ctx context.Context, param *DeleteCollectionParam) error
	List(ctx context.Context, param *ListCollectionParam) (*ListCollectionResult, error)
}

func NewCollection(r repository.Collection, rr repository.Route, re repository.Endpoint, ra repository.Authorize) Collection {
	return &collection{r: r, rr: rr, re: re, ra: ra}
}

type collection struct {
	r  repository.Collection
	rr repository.Route
	re repository.Endpoint
	ra repository.Authorize
}

type CreateCollectionParam struct {
	// 父级节点
	ParentId string `json:"parent_id" form:"parent_id"`
	// 绑定的域名
	ServerNames []string `json:"server_names" form:"server_names"`
	// 分组名
	Name string `json:"name" form:"name" binding:"required"`
	// 分组描述
	Description string `json:"description" form:"description"`
	// 绑定的后端服务
	EndpointId string `json:"endpoint_id" form:"endpoint_id"`
	// 鉴权配置
	AuthorizeId string `json:"authorize_id" form:"authorize_id"`
}

func (s *collection) Create(ctx context.Context, param *CreateCollectionParam) (*dto.Collection, error) {
	var obj *dto.Collection

	database.Transaction(ctx, func(txCtx context.Context) error {
		item, err := s.r.Create(ctx, &entity.Collection{
			Name:        param.Name,
			Description: param.Description,
			ServerNames: param.ServerNames,
			ParentId:    identity.Parse(constant.CollectionPrefix, param.ParentId),
			AuthorizeId: identity.Parse(constant.AuthorizePrefix, param.AuthorizeId),
			EndpointId:  identity.Parse(constant.EndpointPrefix, param.EndpointId),
		})

		if err != nil {
			return err
		}
		obj = dto.NewCollection(item)
		return nil
	})

	return obj, nil
}

type GetCollectionParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *collection) Get(ctx context.Context, param *GetCollectionParam) (*dto.Collection, error) {
	rst, err := s.r.Get(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return nil, err
	}

	collection := dto.NewCollection(rst)

	if utils.InStringSlice("endpoint", param.Expand) {
		ent, err := s.re.Get(ctx, rst.EndpointId)
		if err != nil {
			return nil, err
		}
		collection.Endpoint = dto.NewEndpoint(ent)
	}

	if utils.InStringSlice("authorize", param.Expand) {
		ent, err := s.ra.Get(ctx, rst.AuthorizeId)
		if err != nil {
			return nil, err
		}
		collection.Authorize = dto.NewAuthorize(ent)
	}

	return collection, nil
}

type ListCollectionParam struct {
	ParentId string `form:"parent_id"`
	Name     string `form:"name"`
	Depth    int    `form:"depth" binding:"max=10"`

	// pagination
	Page         int  `json:"page" form:"page"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
}

type ListCollectionResult struct {
	Data  []*dto.Collection `json:"data"`
	Total int64             `json:"total,omitempty"`
}

func (s *collection) List(ctx context.Context, param *ListCollectionParam) (*ListCollectionResult, error) {
	if param.Page == 0 {
		param.Page = 1
	}

	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listParam := &repository.ListCollectionParam{
		Name:     param.Name,
		ParentId: identity.Parse(constant.CollectionPrefix, param.ParentId),
		Depth:    param.Depth,

		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	}

	listRst, err := s.r.List(ctx, listParam)
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*dto.Collection, n)

	for i, v := range listRst.Data {
		items[i] = dto.NewCollection(v)
	}

	rst := &ListCollectionResult{}
	rst.Data = items
	rst.Total = listRst.Total
	return rst, nil
}

type UpdateCollectionParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateCollectionParam
}

func (s *collection) Update(ctx context.Context, param *UpdateCollectionParam) (*dto.Collection, error) {
	database.Transaction(ctx, func(txCtx context.Context) error {
		id := identity.Parse(constant.CollectionPrefix, param.Id)

		err := s.r.Update(ctx, id, &entity.Collection{
			Name:        param.Name,
			ServerNames: param.ServerNames,
			AuthorizeId: identity.Parse(constant.AuthorizePrefix, param.AuthorizeId),
			EndpointId:  identity.Parse(constant.EndpointPrefix, param.EndpointId),
		})

		if err != nil {
			return err
		}

		return nil
	})

	return s.Get(ctx, &GetCollectionParam{Id: param.Id})
}

type DeleteCollectionParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *collection) Delete(ctx context.Context, param *DeleteCollectionParam) error {
	err := s.r.Delete(ctx, identity.Parse(constant.CollectionPrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}
