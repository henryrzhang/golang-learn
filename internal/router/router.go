package router

import (
	"net/http"

	"golang-learn/internal/infra/jwt"
	"golang-learn/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// Handler 可注册到 huma API 的处理器
type Handler interface {
	Register(api huma.API)
}

// Config 路由配置
type Config struct {
	APIConfig huma.Config
	JWT       *jwt.Manager
	Handlers  []Handler
}

// Setup 创建并配置路由，返回根 Router
func Setup(cfg Config) chi.Router {
	r := chi.NewMux()

	// 健康检查（无鉴权、无业务中间件）
	r.Get("/health", healthHandler)
	r.Get("/ready", readyHandler)

	// API 路由组：/api 前缀，统一中间件
	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Auth(cfg.JWT))

	api := humachi.New(apiRouter, cfg.APIConfig)
	for _, h := range cfg.Handlers {
		h.Register(api)
	}

	r.Mount("/api", apiRouter)
	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
