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
	GetSession(ctx context.Context, tokStr string) (uint64, []string, error)
}

func NewUser(r repository.User, rs repository.Session, aseKey []byte) User {
	return &user{r: r, rs: rs, aseKey: aseKey}
}

type user struct {
	r      repository.User
	rs     repository.Session
	aseKey []byte
}

type CreateUserParam struct {
	Name     string   `json:"name"`
	Scopes   []string `json:"scopes"`
	Password string   `json:"password"`
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
	ent.Scopes = param.Scopes

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
	ent.Scopes = param.Scopes

	err := s.r.Update(ctx, id, ent)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetUserParam{Id: param.Id})
}

type CreateUserSessionParam struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"-"`
	Agent    string `json:"-"`
}

type CreateSessionResult struct {
	Type     string    `json:"type"`
	UserId   string    `json:"user_id"`
	Name     string    `json:"name"`
	Token    string    `json:"token"`
	Scopes   []string  `json:"scopes"`
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

	// 创建会话
	ent, err := s.rs.Create(ctx, &entity.Session{UserId: user.Id, Address: param.Address, Agent: param.Agent, ExpireAt: expireAt})
	if err != nil {
		return nil, err
	}

	// 创建 token
	tok := &token.BinaryToken{
		Id:       ent.Id,
		ExpireAt: uint64(expireAt.Unix()),
	}

	rst := &CreateSessionResult{}
	rst.Type = "Bearer"
	rst.Name = user.Name
	rst.UserId = identity.Format(constant.UserPrefix, user.Id)
	rst.ExpireAt = expireAt
	rst.Scopes = user.Scopes
	rst.Token, err = tok.Encrypt(token.NewAesCrypto(s.aseKey))

	if err != nil {
		return nil, err
	}

	return rst, nil
}

func (s *user) GetSession(ctx context.Context, tokStr string) (uint64, []string, error) {
	tok := &token.BinaryToken{}
	err := tok.Decrypt(tokStr, token.NewAesCrypto(s.aseKey))
	if err != nil {
		return 0, nil, err
	}

	if uint64(time.Now().Unix()) > tok.ExpireAt {
		return 0, nil, nil
	}

	session, err := s.rs.Get(ctx, tok.Id)
	if err != nil {
		return 0, nil, nil
	}

	user, err := s.r.Get(ctx, session.UserId)
	if err != nil {
		return 0, nil, nil
	}

	return tok.Id, user.Scopes, nil
}
