#!/bin/bash
# JobTracker Build Script (Linux)

cd "$(dirname "$0")"

echo "=========================================="
echo "  Building JobTracker"
echo "=========================================="

# 清理旧容器
echo "Cleaning up..."
docker compose down -v

# 重新构建
echo "Building services..."
docker compose up -d --build

# 等待启动
echo "Waiting for services (30s)..."
sleep 30

# 检查状态
echo ""
echo "=========================================="
echo "  Status"
echo "=========================================="
docker compose ps

echo ""
echo "=========================================="
echo "  Access URLs"
echo "=========================================="
echo "Frontend:  http://localhost:5173"
echo "API:       http://localhost:8080"
echo "Grafana:   http://localhost:3000 (admin/admin)"
echo "=========================================="