package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB 数据库连接池（纯基建，不含数据访问）
type DB struct {
	Pool *pgxpool.Pool
}

// New 创建数据库连接池
func New(ctx context.Context, connStr string, maxOpen, maxIdle int) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = int32(maxOpen)
	cfg.MinConns = int32(maxIdle)

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

// Close 关闭连接池
func (db *DB) Close() {
	db.Pool.Close()
}
