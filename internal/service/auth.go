package service

import (
	"context"
	"errors"

	"golang-learn/internal/dto"
	"golang-learn/internal/infra/jwt"
	"golang-learn/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists  = errors.New("邮箱已注册")
	ErrInvalidCred  = errors.New("邮箱或密码错误")
	ErrUserDisabled = errors.New("用户已禁用")
)

type AuthService struct {
	db  *pgxpool.Pool
	q   *sqlc.Queries
	jwt *jwt.Manager
}

func NewAuthService(db *pgxpool.Pool, jwtMgr *jwt.Manager) *AuthService {
	return &AuthService{
		db:  db,
		q:   sqlc.New(db),
		jwt: jwtMgr,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 检查邮箱是否已存在
	_, err := s.q.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailExists
	}
	// pgx 返回 pgx.ErrNoRows
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	u, err := s.q.CreateUser(ctx, sqlc.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    phone,
		Password: string(hash),
		Status:   1,
	})
	if err != nil {
		return nil, err
	}

	token, err := s.jwt.Generate(u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  userToPublic(&u),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	u, err := s.q.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCred
		}
		return nil, err
	}

	if u.Status != 1 {
		return nil, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCred
	}

	token, err := s.jwt.Generate(u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  userToPublic(&u),
	}, nil
}

func userToPublic(u *sqlc.User) dto.UserPublic {
	return dto.UserPublic{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
