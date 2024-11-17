package user

import (
	"context"
)

type userExpander struct {
	r UserRepository
}

type UserExpander interface {
	ExpandItem(ctx context.Context, id uint64) (*UserDto, error)
	ExpandItems(ctx context.Context, id []uint64) ([]*UserDto, error)
}

func (e userExpander) ExpandItem(ctx context.Context, id uint64) (*UserDto, error) {
	if user, err := e.r.Get(ctx, id); err != nil {
		return nil, err
	} else {
		return NewUserDto(user), nil
	}
}

func (e userExpander) ExpandItems(ctx context.Context, ids []uint64) ([]*UserDto, error) {
	if user, err := e.r.BatchGet(ctx, ids); err != nil {
		return nil, err
	} else {
		items := make([]*UserDto, 0, len(user))
		for _, v := range user {
			items = append(items, NewUserDto(v))
		}
		return items, nil
	}
}

func NewUserExpander(repo UserRepository) UserExpander {
	return &userExpander{r: repo}
}
