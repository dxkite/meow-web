package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/value"
)

type Authorize interface {
	Create(ctx context.Context, create *CreateAuthorizeParam) (*dto.Authorize, error)
	Update(ctx context.Context, param *UpdateAuthorizeParam) (*dto.Authorize, error)
	Get(ctx context.Context, param *GetAuthorizeParam) (*dto.Authorize, error)
	Delete(ctx context.Context, param *DeleteAuthorizeParam) error
	List(ctx context.Context, param *ListAuthorizeParam) (*ListAuthorizeResult, error)
}

type CreateAuthorizeParam struct {
	Name      string                    `json:"name" binding:"required"`
	Type      string                    `json:"type" binding:"required"`
	Attribute *value.AuthorizeAttribute `json:"attribute"  binding:"required"`
}

func NewAuthorize(r repository.Authorize) Authorize {
	return &authorize{r: r}
}

type authorize struct {
	r repository.Authorize
}

func (s *authorize) Create(ctx context.Context, param *CreateAuthorizeParam) (*dto.Authorize, error) {
	ent := entity.NewAuthorize()

	ent.Name = param.Name
	ent.Type = param.Type
	ent.Attribute = param.Attribute

	resp, err := s.r.Create(ctx, ent)
	if err != nil {
		return nil, err
	}

	return dto.NewAuthorize(resp), nil
}

type GetAuthorizeParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *authorize) Get(ctx context.Context, param *GetAuthorizeParam) (*dto.Authorize, error) {
	ent, err := s.r.Get(ctx, identity.Parse(constant.AuthorizePrefix, param.Id))
	if err != nil {
		return nil, err
	}
	obj := dto.NewAuthorize(ent)
	return obj, nil
}

type DeleteAuthorizeParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *authorize) Delete(ctx context.Context, param *DeleteAuthorizeParam) error {
	err := s.r.Delete(ctx, identity.Parse(constant.AuthorizePrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type ListAuthorizeParam struct {
	Name          string   `form:"name"`
	Limit         int      `form:"limit" binding:"max=1000"`
	StartingAfter string   `form:"starting_after"`
	EndingBefore  string   `form:"ending_before"`
	Expand        []string `json:"expand" form:"expand"`
}

type ListAuthorizeResult struct {
	HasMore bool             `json:"has_more"`
	Data    []*dto.Authorize `json:"data"`
}

func (s *authorize) List(ctx context.Context, param *ListAuthorizeParam) (*ListAuthorizeResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListAuthorizeParam{
		Name:          param.Name,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.AuthorizePrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.AuthorizePrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.Authorize, n)

	for i, v := range entities {
		items[i] = dto.NewAuthorize(v)
	}

	rst := &ListAuthorizeResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

type UpdateAuthorizeParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateAuthorizeParam
}

func (s *authorize) Update(ctx context.Context, param *UpdateAuthorizeParam) (*dto.Authorize, error) {
	id := identity.Parse(constant.AuthorizePrefix, param.Id)

	ent := entity.NewAuthorize()
	ent.Name = param.Name
	ent.Type = param.Type
	ent.Attribute = param.Attribute

	err := s.r.Update(ctx, id, ent)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetAuthorizeParam{Id: param.Id})
}
