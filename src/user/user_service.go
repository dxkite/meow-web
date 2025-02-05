package user

import (
	"context"
	"time"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/pkg/value"
	"dxkite.cn/nebula/pkg/crypto/identity"
	"dxkite.cn/nebula/pkg/crypto/passwd"
	"dxkite.cn/nebula/pkg/crypto/token"
	"dxkite.cn/nebula/pkg/errorx"
	"dxkite.cn/nebula/pkg/httpx"
)

var ErrNamePasswordError = errorx.UnprocessableEntity(errorx.New("name or password error"))
var ErrUserExist = errorx.UnprocessableEntity(errorx.New("user exist"))

type UserService interface {
	Create(ctx context.Context, param *CreateUserRequest) (*UserDto, error)
	Update(ctx context.Context, param *UpdateUserRequest) (*UserDto, error)
	Get(ctx context.Context, param *GetUserRequest) (*UserDto, error)
	Delete(ctx context.Context, param *DeleteUserRequest) error
	List(ctx context.Context, param *ListUserRequest) (*ListUserResponse, error)
	CreateSession(ctx context.Context, param *CreateUserSessionRequest) (*CreateSessionResponse, error)
	DeleteSession(ctx context.Context, req *DeleteSessionRequest) error
	GetSession(ctx context.Context, tokStr string) (httpx.ScopeContext, error)
}

func NewUserService(r UserRepository, rs SessionRepository, cfg *config.Config) UserService {
	return &userService{r: r, rs: rs, aseKey: []byte(cfg.SessionCryptoKey)}
}

type userService struct {
	r      UserRepository
	rs     SessionRepository
	aseKey []byte
}

type CreateUserRequest struct {
	Name     string            `json:"name"`
	Scopes   []httpx.ScopeName `json:"scopes"`
	Password string            `json:"password"`
}

func (s *userService) Create(ctx context.Context, param *CreateUserRequest) (*UserDto, error) {
	user, err := s.r.GetBy(ctx, GetUserByParam{Name: param.Name})
	if err != nil && !errorx.Is(err, ErrUserNotExist) {
		return nil, errorx.System(err)
	}

	if user != nil {
		return nil, errorx.UnprocessableEntity(ErrUserExist)
	}

	ent := NewUser()
	ent.Name = param.Name

	passwdHash, err := passwd.NewHash(param.Password)
	if err != nil {
		return nil, errorx.System(err)
	}

	ent.Password = passwdHash
	ent.Scopes = param.Scopes

	resp, err := s.r.Create(ctx, ent)
	if err != nil {
		return nil, errorx.System(err)
	}

	return NewUserDto(resp), nil
}

type GetUserRequest struct {
	Id     string   `json:"id" path:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *userService) Get(ctx context.Context, param *GetUserRequest) (*UserDto, error) {
	ent, err := s.r.Get(ctx, identity.Parse(UserPrefix, param.Id))
	if err != nil {
		return nil, errorx.System(err)
	}
	obj := NewUserDto(ent)
	return obj, nil
}

type DeleteUserRequest struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *userService) Delete(ctx context.Context, param *DeleteUserRequest) error {
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

func (s *userService) List(ctx context.Context, param *ListUserRequest) (*ListUserResponse, error) {
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

func (s *userService) Update(ctx context.Context, param *UpdateUserRequest) (*UserDto, error) {
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
	Type     string                                       `json:"type"`
	UserId   value.Identity                               `json:"user_id"`
	User     value.Expand[uint64, *UserDto, UserExpander] `json:"user"`
	Name     string                                       `json:"name"`
	Token    string                                       `json:"token"`
	Scopes   []httpx.ScopeName                            `json:"scopes"`
	ExpireAt time.Time                                    `json:"expire_at"`
}

func (s *userService) CreateSession(ctx context.Context, param *CreateUserSessionRequest) (*CreateSessionResponse, error) {
	user, err := s.r.GetBy(ctx, GetUserByParam{Name: param.Name})
	if err != nil {
		return nil, errorx.UnprocessableEntity(ErrNamePasswordError)
	}

	if ok, _ := passwd.VerifyHash(param.Password, user.Password); !ok {
		return nil, errorx.UnprocessableEntity(ErrNamePasswordError)
	}

	// 一小时过期
	expireAt := time.Now().Add(time.Hour)

	// 创建会话
	ent, err := s.rs.Create(ctx, &Session{UserId: user.Id, Address: param.Address, Agent: param.Agent, ExpireAt: expireAt})
	if err != nil {
		return nil, errorx.System(err)
	}

	// 创建 token
	tok := &token.BinaryToken{
		Id:       ent.Id,
		ExpireAt: uint64(expireAt.Unix()),
	}

	rst := &CreateSessionResponse{}
	rst.Type = "Bearer"
	rst.Name = user.Name
	rst.User.Set(ctx, user.Id)
	rst.UserId.Set(ctx, UserPrefix, user.Id)
	rst.ExpireAt = expireAt
	rst.Scopes = user.Scopes
	rst.Token, err = tok.Encrypt(token.NewAesCrypto(s.aseKey))

	if err != nil {
		return nil, errorx.System(err)
	}

	return rst, nil
}

func (s *userService) GetSession(ctx context.Context, tokStr string) (httpx.ScopeContext, error) {
	tok := &token.BinaryToken{}
	err := tok.Decrypt(tokStr, token.NewAesCrypto(s.aseKey))
	if err != nil {
		return nil, errorx.InvalidParameter(errorx.Wrap(err, "invalid token"))
	}

	if uint64(time.Now().Unix()) > tok.ExpireAt {
		return httpx.NewScope(0), nil
	}

	session, err := s.rs.Get(ctx, tok.Id)
	if err != nil {
		return httpx.NewScope(0), nil
	}

	user, err := s.r.Get(ctx, session.UserId)
	if err != nil {
		return httpx.NewScope(0), nil
	}

	return httpx.NewScope(user.Id, user.Scopes...), nil
}

type DeleteSessionRequest struct {
	UserId uint64 `scope-context:"identity" json:"-"`
}

func (s *userService) DeleteSession(ctx context.Context, req *DeleteSessionRequest) error {
	err := s.rs.SetDeletedByUser(ctx, req.UserId)
	if err != nil {
		return err
	}
	return nil
}
