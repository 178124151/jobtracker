#!/bin/bash
# JobTracker K8s 部署脚本

set -e

echo "=========================================="
echo "  JobTracker K8s Deploy"
echo "=========================================="

cd /opt/jobtracker

# 1. 构建镜像
echo "[1/5] Building images..."
sudo docker build -t jobtracker-backend:latest ./backend
sudo docker build -t jobtracker-frontend:latest ./frontend

# 2. 导入到 k3s
echo "[2/5] Importing images to k3s..."
sudo docker save jobtracker-backend:latest | sudo k3s ctr images import -
sudo docker save jobtracker-frontend:latest | sudo k3s ctr images import -

# 3. 删除旧资源
echo "[3/5] Deleting old resources..."
kubectl delete -f infra/k8s/base/ --ignore-not-found 2>/dev/null || true
sleep 5

# 4. 部署新资源
echo "[4/5] Deploying to K8s..."
kubectl apply -f infra/k8s/base/config.yaml
kubectl apply -f infra/k8s/base/postgres.yaml
kubectl apply -f infra/k8s/base/redis.yaml
sleep 10
kubectl apply -f infra/k8s/base/backend.yaml
kubectl apply -f infra/k8s/base/frontend.yaml
kubectl apply -f infra/k8s/base/hpa.yaml
kubectl apply -f infra/k8s/base/ingress.yaml

# 5. 等待启动
echo "[5/5] Waiting for pods..."
sleep 30
kubectl get pods

echo ""
echo "=========================================="
echo "  Deploy Complete!"
echo "=========================================="
echo ""
echo "Access via port-forward:"
echo "  kubectl port-forward svc/frontend-service 8080:80"
echo "  Then visit http://localhost:8080"