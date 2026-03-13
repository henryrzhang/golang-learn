package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"golang-learn/internal/dto"
	"golang-learn/internal/middleware"
	"golang-learn/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type DramaHandler struct {
	svc *service.DramaService
}

func NewDramaHandler(svc *service.DramaService) *DramaHandler {
	return &DramaHandler{svc: svc}
}

func (h *DramaHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-drama",
		Method:      http.MethodPost,
		Path:        "/drama",
		Summary:     "创建剧集",
		Tags:        []string{"剧集"},
	}, h.CreateHandler)
	huma.Register(api, huma.Operation{
		OperationID: "get-drama-by-id",
		Method:      http.MethodGet,
		Path:        "/drama/{id}",
		Summary:     "根据ID获取剧集",
		Tags:        []string{"剧集"},
	}, h.GetByIDHandler)
	huma.Register(api, huma.Operation{
		OperationID: "get-drama-by-no",
		Method:      http.MethodGet,
		Path:        "/drama/no/{drama_no}",
		Summary:     "根据编号获取剧集",
		Tags:        []string{"剧集"},
	}, h.GetByNoHandler)
	huma.Register(api, huma.Operation{
		OperationID: "list-dramas",
		Method:      http.MethodGet,
		Path:        "/drama",
		Summary:     "剧集列表",
		Tags:        []string{"剧集"},
	}, h.ListHandler)
	huma.Register(api, huma.Operation{
		OperationID: "update-drama",
		Method:      http.MethodPatch,
		Path:        "/drama/{id}",
		Summary:     "更新剧集",
		Tags:        []string{"剧集"},
	}, h.UpdateHandler)
	huma.Register(api, huma.Operation{
		OperationID: "delete-drama",
		Method:      http.MethodDelete,
		Path:        "/drama/{id}",
		Summary:     "删除剧集（软删除）",
		Tags:        []string{"剧集"},
	}, h.DeleteHandler)
}

func (h *DramaHandler) CreateHandler(ctx context.Context, input *struct {
	Body dto.DramaCreate
}) (*struct {
	Body dto.DramaResponse
}, error) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, huma.Error401Unauthorized("请先登录")
	}
	email := middleware.GetEmail(ctx)
	resp, err := h.svc.Create(ctx, &input.Body, userID, email)
	if err != nil {
		if errors.Is(err, service.ErrDramaNoExists) {
			return nil, huma.Error409Conflict("剧集编号已存在")
		}
		return nil, err
	}
	return &struct {
		Body dto.DramaResponse
	}{Body: *resp}, nil
}

func (h *DramaHandler) GetByIDHandler(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*struct {
	Body dto.DramaResponse
}, error) {
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return nil, huma.Error400BadRequest("无效的ID")
	}
	d, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrDramaNotFound) {
			return nil, huma.Error404NotFound("剧集不存在")
		}
		return nil, err
	}
	return &struct {
		Body dto.DramaResponse
	}{Body: *d}, nil
}

func (h *DramaHandler) GetByNoHandler(ctx context.Context, input *struct {
	DramaNo string `path:"drama_no"`
}) (*struct {
	Body dto.DramaResponse
}, error) {
	d, err := h.svc.GetByNo(ctx, input.DramaNo)
	if err != nil {
		if errors.Is(err, service.ErrDramaNotFound) {
			return nil, huma.Error404NotFound("剧集不存在")
		}
		return nil, err
	}
	return &struct {
		Body dto.DramaResponse
	}{Body: *d}, nil
}

func (h *DramaHandler) ListHandler(ctx context.Context, input *struct {
	Limit  int `query:"limit" default:"20" doc:"每页数量"`
	Offset int `query:"offset" default:"0" doc:"偏移量"`
}) (*struct {
	Body dto.DramaListResponse
}, error) {
	limit := int32(input.Limit)
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := int32(input.Offset)
	if offset < 0 {
		offset = 0
	}
	resp, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return &struct {
		Body dto.DramaListResponse
	}{Body: *resp}, nil
}

func (h *DramaHandler) UpdateHandler(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body dto.DramaUpdate
}) (*struct {
	Body dto.DramaResponse
}, error) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, huma.Error401Unauthorized("请先登录")
	}
	email := middleware.GetEmail(ctx)
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return nil, huma.Error400BadRequest("无效的ID")
	}
	d, err := h.svc.Update(ctx, id, &input.Body, email)
	if err != nil {
		if errors.Is(err, service.ErrDramaNotFound) {
			return nil, huma.Error404NotFound("剧集不存在")
		}
		return nil, err
	}
	return &struct {
		Body dto.DramaResponse
	}{Body: *d}, nil
}

func (h *DramaHandler) DeleteHandler(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*struct{}, error) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, huma.Error401Unauthorized("请先登录")
	}
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return nil, huma.Error400BadRequest("无效的ID")
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, service.ErrDramaNotFound) {
			return nil, huma.Error404NotFound("剧集不存在")
		}
		return nil, err
	}
	return &struct{}{}, nil
}
