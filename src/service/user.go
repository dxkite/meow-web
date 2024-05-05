package service

import (
	"context"
	"errors"
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/pkg/passwd"
	"dxkite.cn/meownest/pkg/token"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
)

var ErrNamePasswordError = errors.New("name or password error")
var ErrUserExist = errors.New("user exist")

type User interface {
	Create(ctx context.Context, param *CreateUserParam) (*dto.User, error)
	Update(ctx context.Context, param *UpdateUserParam) (*dto.User, error)
	Get(ctx context.Context, param *GetUserParam) (*dto.User, error)
	Delete(ctx context.Context, param *DeleteUserParam) error
	List(ctx context.Context, param *ListUserParam) (*ListUserResult, error)
	Session(ctx context.Context, param *CreateUserSessionParam) (*CreateSessionResult, error)
}

func NewUser(r repository.User, aseKey []byte) User {
	return &user{r: r, aseKey: aseKey}
}

type user struct {
	r      repository.User
	aseKey []byte
}

type CreateUserParam struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *user) Create(ctx context.Context, param *CreateUserParam) (*dto.User, error) {
	user, err := s.r.GetBy(ctx, repository.GetUserByParam{Name: param.Name})
	if err != nil && !errors.Is(err, repository.ErrUserNotExist) {
		return nil, err
	}

	if user != nil {
		return nil, ErrUserExist
	}

	ent := entity.NewUser()
	ent.Name = param.Name

	passwdHash, err := passwd.NewHash(param.Password)
	if err != nil {
		return nil, err
	}

	ent.Password = passwdHash

	resp, err := s.r.Create(ctx, ent)
	if err != nil {
		return nil, err
	}

	return dto.NewUser(resp), nil
}

type GetUserParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *user) Get(ctx context.Context, param *GetUserParam) (*dto.User, error) {
	ent, err := s.r.Get(ctx, identity.Parse(constant.UserPrefix, param.Id))
	if err != nil {
		return nil, err
	}
	obj := dto.NewUser(ent)
	return obj, nil
}

type DeleteUserParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *user) Delete(ctx context.Context, param *DeleteUserParam) error {
	err := s.r.Delete(ctx, identity.Parse(constant.UserPrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type ListUserParam struct {
	Limit         int      `form:"limit" binding:"max=1000"`
	StartingAfter string   `form:"starting_after"`
	EndingBefore  string   `form:"ending_before"`
	Expand        []string `json:"expand" form:"expand"`
}

type ListUserResult struct {
	HasMore bool        `json:"has_more"`
	Data    []*dto.User `json:"data"`
}

func (s *user) List(ctx context.Context, param *ListUserParam) (*ListUserResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListUserParam{
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.UserPrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.UserPrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.User, n)

	for i, v := range entities {
		items[i] = dto.NewUser(v)
	}

	rst := &ListUserResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

type UpdateUserParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateUserParam
}

func (s *user) Update(ctx context.Context, param *UpdateUserParam) (*dto.User, error) {
	id := identity.Parse(constant.UserPrefix, param.Id)
	ent := entity.NewUser()

	err := s.r.Update(ctx, id, ent)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetUserParam{Id: param.Id})
}

type CreateUserSessionParam struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateSessionResult struct {
	Type     string    `json:"type"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

func (s *user) Session(ctx context.Context, param *CreateUserSessionParam) (*CreateSessionResult, error) {
	user, err := s.r.GetBy(ctx, repository.GetUserByParam{Name: param.Name})
	if err != nil {
		return nil, ErrNamePasswordError
	}

	if ok, _ := passwd.VerifyHash(param.Password, user.Password); !ok {
		return nil, ErrNamePasswordError
	}

	// 一小时过期
	expireAt := time.Now().Add(time.Hour)

	tok := &token.BinaryToken{
		Id:       user.Id,
		ExpireAt: uint64(expireAt.Unix()),
	}

	rst := &CreateSessionResult{}
	rst.Type = "Bearer"
	rst.ExpireAt = expireAt
	rst.Token, err = tok.Encrypt(token.NewAesCrypto(s.aseKey))

	if err != nil {
		return nil, err
	}

	return rst, nil
}
