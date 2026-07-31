#!/bin/bash
# JobTracker K8s Deploy Script
set -e

echo "=========================================="
echo "  JobTracker K8s Deploy"
echo "=========================================="

cd /opt/jobtracker

# 1. 拉取镜像
echo "[1/6] Pulling images..."
sudo docker pull postgres:16-alpine
sudo docker pull redis:7.2-alpine
sudo docker pull grafana/grafana:latest

# 2. 构建应用镜像
echo "[2/6] Building application images..."
sudo docker build -t jobtracker-backend:latest ./backend
sudo docker build -t jobtracker-frontend:latest ./frontend

# 3. 保存镜像到tar文件
echo "[3/6] Saving images..."
sudo docker save -o /tmp/postgres.tar postgres:16-alpine
sudo docker save -o /tmp/redis.tar redis:7.2-alpine
sudo docker save -o /tmp/grafana.tar grafana/grafana:latest
sudo docker save -o /tmp/backend.tar jobtracker-backend:latest
sudo docker save -o /tmp/frontend.tar jobtracker-frontend:latest

# 4. 导入到 k3s
echo "[4/6] Importing images to k3s..."
sudo k3s ctr images import /tmp/postgres.tar
sudo k3s ctr images import /tmp/redis.tar
sudo k3s ctr images import /tmp/grafana.tar
sudo k3s ctr images import /tmp/backend.tar
sudo k3s ctr images import /tmp/frontend.tar

# 清理临时文件
rm -f /tmp/postgres.tar /tmp/redis.tar /tmp/grafana.tar /tmp/backend.tar /tmp/frontend.tar

# 5. 部署到 K8s
echo "[5/6] Deploying to K8s..."
kubectl apply -f infra/k8s/base/config.yaml
kubectl apply -f infra/k8s/base/postgres.yaml
kubectl apply -f infra/k8s/base/redis.yaml
kubectl apply -f infra/k8s/base/grafana.yaml
sleep 15
kubectl apply -f infra/k8s/base/backend.yaml
kubectl apply -f infra/k8s/base/frontend.yaml
kubectl apply -f infra/k8s/base/hpa.yaml
kubectl apply -f infra/k8s/base/ingress.yaml

# 6. 等待启动
echo "[6/6] Waiting for pods..."
sleep 30
kubectl get pods

echo ""
echo "=========================================="
echo "  Deploy Complete!"
echo "=========================================="
echo ""
echo "Access URLs:"
echo "  Frontend: http://101.35.246.75:8080"
echo "  Grafana:  http://101.35.246.75:30300"
echo ""
echo "After deployment, run seed script:"
echo "  bash scripts/seed-data.sh"