#!/bin/bash
# K8s 一键部署脚本
# 用法: bash deploy-k8s.sh

set -e

echo "=========================================="
echo "  JobTracker K8s 部署"
echo "=========================================="

# 检查 kubectl
if ! command -v kubectl &> /dev/null; then
    echo "错误: kubectl 未安装"
    exit 1
fi

# 检查集群连接
echo "[1/4] 检查集群连接..."
kubectl cluster-info

# 构建镜像（如果使用 k3s）
echo "[2/4] 构建 Docker 镜像..."
docker build -t jobtracker-backend:latest ./backend
docker build -t jobtracker-frontend:latest ./frontend

# 导入镜像到 k3s（如果是 k3s 环境）
if command -v k3s &> /dev/null; then
    echo "导入镜像到 k3s..."
    docker save jobtracker-backend:latest | k3s ctr images import -
    docker save jobtracker-frontend:latest | k3s ctr images import -
fi

# 部署
echo "[3/4] 部署到 K8s..."
kubectl apply -k infra/k8s/base/

# 等待就绪
echo "[4/4] 等待 Pod 就绪..."
kubectl wait --for=condition=Ready pods -l app=jobtracker --timeout=180s

echo ""
echo "=========================================="
echo "  部署完成！"
echo "=========================================="
echo ""
kubectl get pods -l app=jobtracker
echo ""
echo "访问方式:"
echo "  kubectl port-forward svc/frontend-service 8080:80"
echo "  然后访问 http://localhost:8080"