package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang-learn/config"
	"golang-learn/internal/handler"
	"golang-learn/internal/infra/db"
	"golang-learn/internal/infra/jwt"
	"golang-learn/internal/infra/logger"
	"golang-learn/internal/infra/redis"
	"golang-learn/internal/router"
	"golang-learn/internal/service"

	"github.com/danielgtaylor/huma/v2"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func main() {
	cfgPath := "config/config.yaml"
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		cfgPath = v
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.L().Info().Str("path", cfgPath).Msg("config loaded")

	ctx := context.Background()

	// 数据库
	database, err := db.New(ctx, cfg.Database.URL, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns)
	if err != nil {
		logger.L().Fatal().Err(err).Msg("connect database failed")
	}
	defer database.Close()
	logger.L().Info().Msg("database connected")

	// Redis
	rdb, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.L().Warn().Err(err).Msg("connect redis failed, cache disabled")
	} else {
		defer rdb.Close()
		logger.L().Info().Str("addr", cfg.Redis.Addr).Msg("redis connected")
	}

	// JWT
	jwtMgr := jwt.New(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	logger.L().Info().Int("expire_hours", cfg.JWT.ExpireHours).Msg("jwt initialized")

	// 路由
	userSvc := service.NewUserService(database.Pool, rdb, jwtMgr)
	dramaSvc := service.NewDramaService(database.Pool)

	r := router.Setup(router.Config{
		APIConfig: huma.DefaultConfig("Golang Learn API", "1.0.0"),
		JWT:       jwtMgr,
		Handlers: []router.Handler{
			handler.NewUserHandler(userSvc),
			handler.NewDramaHandler(dramaSvc),
		},
	})
	logger.L().Info().Msg("router initialized")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		logger.L().Info().Str("addr", addr).Msg("server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L().Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L().Info().Msg("shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Error().Err(err).Msg("server shutdown error")
	}
}
