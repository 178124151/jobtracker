#!/bin/bash
# JobTracker K8s Deploy Script
set -e

echo "=========================================="
echo "  JobTracker K8s Deploy"
echo "=========================================="

cd /opt/jobtracker

# 1. 预拉取镜像
echo "[1/6] Pulling images..."
sudo docker pull postgres:16-alpine
sudo docker pull redis:7.2-alpine
sudo docker pull grafana/grafana:latest
sudo docker pull prom/node-exporter:latest

# 2. 构建应用镜像
echo "[2/6] Building application images..."
sudo docker build -t jobtracker-backend:latest ./backend
sudo docker build -t jobtracker-frontend:latest ./frontend

# 3. 导入到 k3s
echo "[3/6] Importing images to k3s..."
sudo docker save postgres:16-alpine | sudo k3s ctr images import -
sudo docker save redis:7.2-alpine | sudo k3s ctr images import -
sudo docker save grafana/grafana:latest | sudo k3s ctr images import -
sudo docker save prom/node-exporter:latest | sudo k3s ctr images import -
sudo docker save jobtracker-backend:latest | sudo k3s ctr images import -
sudo docker save jobtracker-frontend:latest | sudo k3s ctr images import -

# 4. 部署到 K8s
echo "[4/6] Deploying to K8s..."
kubectl apply -f infra/k8s/base/config.yaml
kubectl apply -f infra/k8s/base/postgres.yaml
kubectl apply -f infra/k8s/base/redis.yaml
kubectl apply -f infra/k8s/base/grafana.yaml
sleep 15
kubectl apply -f infra/k8s/base/backend.yaml
kubectl apply -f infra/k8s/base/frontend.yaml
kubectl apply -f infra/k8s/base/hpa.yaml
kubectl apply -f infra/k8s/base/ingress.yaml

# 5. 等待启动
echo "[5/6] Waiting for pods..."
sleep 30
kubectl get pods

# 6. 导入种子数据
echo "[6/6] Importing seed data..."
bash scripts/seed-data.sh

echo ""
echo "=========================================="
echo "  Deploy Complete!"
echo "=========================================="
echo ""
echo "Access URLs:"
echo "  Frontend: http://101.35.246.75:8080"
echo "  Grafana:  http://101.35.246.75:30300"