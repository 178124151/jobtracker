.PHONY: build run test clean docker-build docker-up docker-down

# 变量
APP_NAME = jobtracker
DOCKER_COMPOSE = docker compose

# 开发
run:
	cd backend && go run cmd/server/main.go

# 构建
build:
	cd backend && CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

# 测试
test:
	cd backend && go test -v -coverprofile=coverage.out ./...

# 测试覆盖率
cover:
	cd backend && go tool cover -html=coverage.out -o coverage.html

# Docker
docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-restart:
	$(DOCKER_COMPOSE) restart

# 日志
logs:
	$(DOCKER_COMPOSE) logs -f

logs-backend:
	$(DOCKER_COMPOSE) logs -f backend

# 部署
deploy:
	$(DOCKER_COMPOSE) up -d --build

# 清理
clean:
	$(DOCKER_COMPOSE) down -v
	docker system prune -f