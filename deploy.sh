#!/bin/bash
# JobTracker 部署脚本
# 用法: ./deploy.sh

set -e

echo "=== JobTracker Deploy Script ==="

# 拉取最新代码
echo "[1/4] Pulling latest code..."
cd /opt/jobtracker
git pull origin main

# 重新构建并启动
echo "[2/4] Rebuilding services..."
sudo docker compose up -d --build

# 等待服务启动
echo "[3/4] Waiting for services..."
sleep 10

# 检查状态
echo "[4/4] Checking status..."
sudo docker compose ps

echo ""
echo "=== Deploy Complete ==="
echo "Frontend: http://$(curl -s ifconfig.me):5173"
echo "Backend:  http://$(curl -s ifconfig.me):8080"
echo "Grafana:  http://$(curl -s ifconfig.me):3000"