.PHONY: build run sqlc migrate test

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

# 需先安装 sqlc: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc:
	cd sqlc && sqlc generate

# 使用 psql 执行建表（需配置 DATABASE_URL 或修改下方连接串）
migrate:
	psql "postgres://postgres:postgres@localhost:5432/golang_learn?sslmode=disable" -f doc/schema.sql

test:
	go test ./...
