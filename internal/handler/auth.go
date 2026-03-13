package handler

import (
	"context"
	"errors"
	"net/http"

	"golang-learn/internal/dto"
	"golang-learn/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/auth/register",
		Summary:     "用户注册",
		Tags:        []string{"认证"},
	}, h.RegisterHandler)
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "用户登录",
		Tags:        []string{"认证"},
	}, h.LoginHandler)
}

func (h *AuthHandler) RegisterHandler(ctx context.Context, input *struct {
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

func (h *AuthHandler) LoginHandler(ctx context.Context, input *struct {
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
