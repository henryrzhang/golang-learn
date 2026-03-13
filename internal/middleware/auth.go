package middleware

import (
	"context"
	"net/http"
	"strings"

	"golang-learn/internal/infra/jwt"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyEmail  contextKey = "email"
)

// Auth 从 Authorization: Bearer <token> 解析 JWT 并注入 context
func Auth(jwtMgr *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := jwtMgr.Parse(tokenStr)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID 从 context 获取当前用户 ID，未登录返回 0
func GetUserID(ctx context.Context) int64 {
	v, _ := ctx.Value(ContextKeyUserID).(int64)
	return v
}

// GetEmail 从 context 获取当前用户邮箱
func GetEmail(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyEmail).(string)
	return v
}

// RequireAuth 要求必须登录，未登录返回 401
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserID(r.Context()) == 0 {
			http.Error(w, `{"detail":"未登录或 token 无效"}`, http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			return
		}
		next.ServeHTTP(w, r)
	})
}
