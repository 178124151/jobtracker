.PHONY: build run test clean docker-build docker-up docker-down docker-push deploy

# 变量
APP_NAME = jobtracker
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
DOCKER_COMPOSE = docker compose

# Go 编译参数
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# ===== 开发 =====
run:
	cd backend && go run $(LDFLAGS) cmd/server/main.go

# ===== 构建 =====
build:
	cd backend && CGO_ENABLED=0 go build $(LDFLAGS) -o bin/server cmd/server/main.go

# ===== 测试 =====
test:
	cd backend && go test -v -coverprofile=coverage.out ./...

cover:
	cd backend && go tool cover -html=coverage.out -o coverage.html

# ===== Docker =====
docker-build:
	docker build -t $(APP_NAME)-backend:$(VERSION) -t $(APP_NAME)-backend:latest ./backend
	docker build -t $(APP_NAME)-frontend:$(VERSION) -t $(APP_NAME)-frontend:latest ./frontend

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-restart:
	$(DOCKER_COMPOSE) restart

docker-push:
	docker push $(APP_NAME)-backend:$(VERSION)
	docker push $(APP_NAME)-backend:latest
	docker push $(APP_NAME)-frontend:$(VERSION)
	docker push $(APP_NAME)-frontend:latest

# ===== 日志 =====
logs:
	$(DOCKER_COMPOSE) logs -f

logs-backend:
	$(DOCKER_COMPOSE) logs -f backend

# ===== 部署 =====
deploy:
	$(DOCKER_COMPOSE) up -d --build

# ===== 清理 =====
clean:
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

# ===== 版本信息 =====
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Build:   $(BUILD_TIME)"