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
	Name string `form:"name"`

	// pagination
	Page         int  `json:"page" form:"page"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
}

type ListUserResult struct {
	Data  []*dto.User `json:"data"`
	Total int64       `json:"total,omitempty"`
}

func (s *user) List(ctx context.Context, param *ListUserParam) (*ListUserResult, error) {
	if param.Page == 0 {
		param.Page = 1
	}

	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listRst, err := s.r.List(ctx, &repository.ListUserParam{
		Name:         param.Name,
		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	})
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*dto.User, n)

	for i, v := range listRst.Data {
		items[i] = dto.NewUser(v)
	}

	rst := &ListUserResult{}
	rst.Data = items
	rst.Total = listRst.Total
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
