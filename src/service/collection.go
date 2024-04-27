package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/model"
	"dxkite.cn/meownest/src/repository"
)

type CreateCollectionParam struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	ParentId    string `json:"parent_id" form:"parent_id"`
}

type Collection interface {
	Create(ctx context.Context, param *CreateCollectionParam) (*dto.Collection, error)
}

func NewCollection(r repository.Collection) Collection {
	return &collection{r: r}
}

type collection struct {
	r repository.Collection
}

func (s *collection) Create(ctx context.Context, param *CreateCollectionParam) (*dto.Collection, error) {
	rst, err := s.r.Create(ctx, &model.Collection{
		Name:        param.Name,
		Description: param.Description,
		ParentId:    identity.Parse(constant.CollectionPrefix, param.ParentId),
	})
	if err != nil {
		return nil, err
	}
	return dto.NewCollection(rst), nil
}
