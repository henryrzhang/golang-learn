package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"golang-learn/internal/dto"
	"golang-learn/internal/infra/jwt"
	"golang-learn/internal/infra/redis"
	"golang-learn/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrEmailExists  = errors.New("邮箱已注册")
	ErrInvalidCred  = errors.New("邮箱或密码错误")
	ErrUserDisabled = errors.New("用户已禁用")
)

const userCacheKeyPrefix = "user:info:"
const userCacheTTL = time.Hour

type UserService struct {
	db    *pgxpool.Pool
	q     *sqlc.Queries
	redis *redis.Client
	jwt   *jwt.Manager
}

func NewUserService(db *pgxpool.Pool, rdb *redis.Client, jwtMgr *jwt.Manager) *UserService {
	return &UserService{
		db:    db,
		q:     sqlc.New(db),
		redis: rdb,
		jwt:   jwtMgr,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	_, err := s.q.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailExists
	}
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

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
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

func (s *UserService) GetByID(ctx context.Context, id int64) (*dto.UserPublic, error) {
	if s.redis != nil {
		if pub := s.getFromCache(ctx, id); pub != nil {
			return pub, nil
		}
	}

	u, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	pub := userToPublic(&u)

	if s.redis != nil {
		s.setCache(ctx, id, &pub)
	}
	return &pub, nil
}

func (s *UserService) getFromCache(ctx context.Context, id int64) *dto.UserPublic {
	key := userCacheKeyPrefix + strconv.FormatInt(id, 10)
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	var pub dto.UserPublic
	if err := json.Unmarshal([]byte(val), &pub); err != nil {
		return nil
	}
	return &pub
}

func (s *UserService) setCache(ctx context.Context, id int64, pub *dto.UserPublic) {
	key := userCacheKeyPrefix + strconv.FormatInt(id, 10)
	b, err := json.Marshal(pub)
	if err != nil {
		return
	}
	_ = s.redis.Set(ctx, key, b, userCacheTTL).Err()
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
