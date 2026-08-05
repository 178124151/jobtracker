#!/bin/bash
# JobTracker 部署脚本
# 用法: sudo bash deploy.sh

set -e

echo "=== JobTracker Deploy ==="

cd /opt/jobtracker

echo "[1/4] 拉取最新代码..."
git pull origin main

if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
fi

echo "[2/4] 重新构建服务..."
docker compose up -d --build

echo "[3/4] 等待服务启动..."
sleep 15

echo "[4/4] 检查状态..."
docker compose ps

echo ""
echo "=== 部署完成 ==="
echo "前端: http://$(curl -s ifconfig.me)"
echo "监控: http://$(curl -s ifconfig.me)/grafana/"
echo "API:  http://$(curl -s ifconfig.me)/api/"
