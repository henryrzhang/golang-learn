package service

import (
	"context"
	"errors"

	"golang-learn/internal/dto"
	"golang-learn/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("用户不存在")

type UserService struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		db: db,
		q:  sqlc.New(db),
	}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*dto.UserPublic, error) {
	u, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	pub := userToPublic(&u)
	return &pub, nil
}
