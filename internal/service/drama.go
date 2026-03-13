package service

import (
	"context"
	"errors"
	"time"

	"golang-learn/internal/dto"
	"golang-learn/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDramaNotFound = errors.New("剧集不存在")
	ErrDramaNoExists = errors.New("剧集编号已存在")
)

type DramaService struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewDramaService(db *pgxpool.Pool) *DramaService {
	return &DramaService{
		db: db,
		q:  sqlc.New(db),
	}
}

func (s *DramaService) Create(ctx context.Context, req *dto.DramaCreate, userID int64, email string) (*dto.DramaResponse, error) {
	now := time.Now().UnixMilli()
	status := int16(1)
	if req.Status != nil {
		status = *req.Status
	}

	_, err := s.q.GetDramaByNo(ctx, req.DramaNo)
	if err == nil {
		return nil, ErrDramaNoExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	d, err := s.q.CreateDrama(ctx, sqlc.CreateDramaParams{
		DramaNo:               req.DramaNo,
		Title:                 req.Title,
		Outline:               req.Outline,
		CoverImage:            req.CoverImage,
		Characters:            req.Characters,
		CharacterRelationDesc: req.CharacterRelationDesc,
		Status:                &status,
		TaskNo:                req.TaskNo,
		CreateBy:              email,
		UpdateBy:              email,
		CreateAt:              now,
		UpdateAt:              now,
		Deleted:               false,
	})
	if err != nil {
		return nil, err
	}
	return dramaToResponse(&d), nil
}

func (s *DramaService) GetByID(ctx context.Context, id int64) (*dto.DramaResponse, error) {
	d, err := s.q.GetDramaByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDramaNotFound
		}
		return nil, err
	}
	return dramaToResponse(&d), nil
}

func (s *DramaService) GetByNo(ctx context.Context, dramaNo string) (*dto.DramaResponse, error) {
	d, err := s.q.GetDramaByNo(ctx, dramaNo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDramaNotFound
		}
		return nil, err
	}
	return dramaToResponse(&d), nil
}

func (s *DramaService) List(ctx context.Context, limit, offset int32) (*dto.DramaListResponse, error) {
	items, err := s.q.ListDramas(ctx, sqlc.ListDramasParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	res := make([]dto.DramaResponse, len(items))
	for i := range items {
		res[i] = *dramaToResponse(&items[i])
	}
	return &dto.DramaListResponse{
		Items: res,
		Total: int64(len(res)),
	}, nil
}

func (s *DramaService) Update(ctx context.Context, id int64, req *dto.DramaUpdate, email string) (*dto.DramaResponse, error) {
	now := time.Now().UnixMilli()
	_, err := s.q.GetDramaByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDramaNotFound
		}
		return nil, err
	}

	d, err := s.q.UpdateDrama(ctx, sqlc.UpdateDramaParams{
		ID:                    id,
		Title:                 req.Title,
		Outline:               req.Outline,
		CoverImage:            req.CoverImage,
		Characters:            req.Characters,
		CharacterRelationDesc: req.CharacterRelationDesc,
		Status:                req.Status,
		TaskNo:                req.TaskNo,
		UpdateBy:              &email,
		UpdateAt:              now,
	})
	if err != nil {
		return nil, err
	}
	return dramaToResponse(&d), nil
}

func (s *DramaService) Delete(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	_, err := s.q.SoftDeleteDrama(ctx, sqlc.SoftDeleteDramaParams{ID: id, UpdateAt: now})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDramaNotFound
		}
		return err
	}
	return nil
}

func dramaToResponse(d *sqlc.DramaInfo) *dto.DramaResponse {
	return &dto.DramaResponse{
		ID:                    d.ID,
		DramaNo:               d.DramaNo,
		Title:                 d.Title,
		Outline:               d.Outline,
		CoverImage:            d.CoverImage,
		Characters:            d.Characters,
		CharacterRelationDesc: d.CharacterRelationDesc,
		Status:                d.Status,
		TaskNo:                d.TaskNo,
		CreateBy:              d.CreateBy,
		UpdateBy:              d.UpdateBy,
		CreateAt:              d.CreateAt,
		UpdateAt:              d.UpdateAt,
	}
}
