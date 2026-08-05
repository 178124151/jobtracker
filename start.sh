#!/bin/bash
# JobTracker Start Script (Linux)

cd "$(dirname "$0")"

echo "=========================================="
echo "  Starting JobTracker"
echo "=========================================="

# 启动所有服务
docker compose up -d

echo ""
echo "=========================================="
echo "  Access URLs"
echo "=========================================="
echo "Frontend:  http://localhost:5173"
echo "API:       http://localhost:8080"
echo "Grafana:   http://localhost:3000"
echo "Prometheus: http://localhost:9090"
echo "=========================================="
