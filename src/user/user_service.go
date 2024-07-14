package user

import (
	"context"
	"errors"
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/pkg/passwd"
	"dxkite.cn/meownest/pkg/token"
)

var ErrNamePasswordError = errors.New("name or password error")
var ErrUserExist = errors.New("user exist")

type UserService interface {
	Create(ctx context.Context, param *CreateUserRequest) (*UserDto, error)
	Update(ctx context.Context, param *UpdateUserRequest) (*UserDto, error)
	Get(ctx context.Context, param *GetUserRequest) (*UserDto, error)
	Delete(ctx context.Context, param *DeleteUserRequest) error
	List(ctx context.Context, param *ListUserRequest) (*ListUserResponse, error)
	CreateSession(ctx context.Context, param *CreateUserSessionRequest) (*CreateSessionResponse, error)
	DeleteSession(ctx context.Context, userId uint64) error
	GetSession(ctx context.Context, tokStr string) (uint64, []string, error)
}

func NewUserService(r UserRepository, rs SessionRepository, aseKey []byte) UserService {
	return &user{r: r, rs: rs, aseKey: aseKey}
}

type user struct {
	r      UserRepository
	rs     SessionRepository
	aseKey []byte
}

type CreateUserRequest struct {
	Name     string   `json:"name"`
	Scopes   []string `json:"scopes"`
	Password string   `json:"password"`
}

func (s *user) Create(ctx context.Context, param *CreateUserRequest) (*UserDto, error) {
	user, err := s.r.GetBy(ctx, GetUserByParam{Name: param.Name})
	if err != nil && !errors.Is(err, ErrUserNotExist) {
		return nil, err
	}

	if user != nil {
		return nil, ErrUserExist
	}

	ent := NewUser()
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

	return NewUserDto(resp), nil
}

type GetUserRequest struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *user) Get(ctx context.Context, param *GetUserRequest) (*UserDto, error) {
	ent, err := s.r.Get(ctx, identity.Parse(UserPrefix, param.Id))
	if err != nil {
		return nil, err
	}
	obj := NewUserDto(ent)
	return obj, nil
}

type DeleteUserRequest struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *user) Delete(ctx context.Context, param *DeleteUserRequest) error {
	err := s.r.Delete(ctx, identity.Parse(UserPrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type ListUserRequest struct {
	Name string `form:"name"`

	// pagination
	Page         int  `json:"page" form:"page"`
	PerPage      int  `json:"per_page" form:"per_page" binding:"max=1000"`
	IncludeTotal bool `json:"include_total" form:"include_total"`
}

type ListUserResponse struct {
	Data  []*UserDto `json:"data"`
	Total int64      `json:"total,omitempty"`
}

func (s *user) List(ctx context.Context, param *ListUserRequest) (*ListUserResponse, error) {
	if param.Page == 0 {
		param.Page = 1
	}

	if param.PerPage == 0 {
		param.PerPage = 10
	}

	listRst, err := s.r.List(ctx, &ListUserParam{
		Name:         param.Name,
		Page:         param.Page,
		PerPage:      param.PerPage,
		IncludeTotal: param.IncludeTotal,
	})
	if err != nil {
		return nil, err
	}

	n := len(listRst.Data)

	items := make([]*UserDto, n)

	for i, v := range listRst.Data {
		items[i] = NewUserDto(v)
	}

	rst := &ListUserResponse{}
	rst.Data = items
	rst.Total = listRst.Total
	return rst, nil
}

type UpdateUserRequest struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateUserRequest
}

func (s *user) Update(ctx context.Context, param *UpdateUserRequest) (*UserDto, error) {
	id := identity.Parse(UserPrefix, param.Id)
	ent := NewUser()
	ent.Scopes = param.Scopes

	err := s.r.Update(ctx, id, ent)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetUserRequest{Id: param.Id})
}

type CreateUserSessionRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"-"`
	Agent    string `json:"-"`
}

type CreateSessionResponse struct {
	Type     string    `json:"type"`
	UserId   string    `json:"user_id"`
	Name     string    `json:"name"`
	Token    string    `json:"token"`
	Scopes   []string  `json:"scopes"`
	ExpireAt time.Time `json:"expire_at"`
}

func (s *user) CreateSession(ctx context.Context, param *CreateUserSessionRequest) (*CreateSessionResponse, error) {
	user, err := s.r.GetBy(ctx, GetUserByParam{Name: param.Name})
	if err != nil {
		return nil, ErrNamePasswordError
	}

	if ok, _ := passwd.VerifyHash(param.Password, user.Password); !ok {
		return nil, ErrNamePasswordError
	}

	// 一小时过期
	expireAt := time.Now().Add(time.Hour)

	// 创建会话
	ent, err := s.rs.Create(ctx, &Session{UserId: user.Id, Address: param.Address, Agent: param.Agent, ExpireAt: expireAt})
	if err != nil {
		return nil, err
	}

	// 创建 token
	tok := &token.BinaryToken{
		Id:       ent.Id,
		ExpireAt: uint64(expireAt.Unix()),
	}

	rst := &CreateSessionResponse{}
	rst.Type = "Bearer"
	rst.Name = user.Name
	rst.UserId = identity.Format(UserPrefix, user.Id)
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

	return user.Id, user.Scopes, nil
}

func (s *user) DeleteSession(ctx context.Context, userId uint64) error {
	err := s.rs.SetDeletedByUser(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
