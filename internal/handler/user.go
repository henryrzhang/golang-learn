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

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/user/register",
		Summary:     "用户注册",
		Tags:        []string{"用户"},
	}, h.RegisterHandler)
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/user/login",
		Summary:     "用户登录",
		Tags:        []string{"用户"},
	}, h.LoginHandler)
	huma.Register(api, huma.Operation{
		OperationID: "get-current-user",
		Method:      http.MethodGet,
		Path:        "/user/me",
		Summary:     "获取当前用户信息",
		Tags:        []string{"用户"},
	}, h.GetMeHandler)
	huma.Register(api, huma.Operation{
		OperationID: "get-user-by-id",
		Method:      http.MethodGet,
		Path:        "/user/{id}",
		Summary:     "根据ID获取用户信息",
		Tags:        []string{"用户"},
	}, h.GetByIDHandler)
}

func (h *UserHandler) RegisterHandler(ctx context.Context, input *struct {
	Body dto.RegisterRequest
}) (*struct {
	Body dto.AuthResponse
}, error) {
	resp, err := h.svc.Register(ctx, &input.Body)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			return nil, huma.Error409Conflict("邮箱已注册")
		}
		return nil, err
	}
	return &struct {
		Body dto.AuthResponse
	}{Body: *resp}, nil
}

func (h *UserHandler) LoginHandler(ctx context.Context, input *struct {
	Body dto.LoginRequest
}) (*struct {
	Body dto.AuthResponse
}, error) {
	resp, err := h.svc.Login(ctx, &input.Body)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCred) {
			return nil, huma.Error401Unauthorized("邮箱或密码错误")
		}
		if errors.Is(err, service.ErrUserDisabled) {
			return nil, huma.Error403Forbidden("用户已禁用")
		}
		return nil, err
	}
	return &struct {
		Body dto.AuthResponse
	}{Body: *resp}, nil
}

func (h *UserHandler) GetMeHandler(ctx context.Context, input *struct{}) (*struct {
	Body dto.UserPublic
}, error) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, huma.Error401Unauthorized("请先登录")
	}
	u, err := h.svc.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, huma.Error404NotFound("用户不存在")
		}
		return nil, err
	}
	return &struct {
		Body dto.UserPublic
	}{Body: *u}, nil
}

func (h *UserHandler) GetByIDHandler(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*struct {
	Body dto.UserPublic
}, error) {
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return nil, huma.Error400BadRequest("无效的用户ID")
	}
	u, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, huma.Error404NotFound("用户不存在")
		}
		return nil, err
	}
	return &struct {
		Body dto.UserPublic
	}{Body: *u}, nil
}
