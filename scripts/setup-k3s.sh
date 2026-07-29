#!/bin/bash
# k3s 单节点集群搭建脚本
# 用法: sudo bash setup-k3s.sh

set -e

echo "=========================================="
echo "  k3s 单节点集群搭建"
echo "=========================================="

# 1. 安装 k3s
echo "[1/5] 安装 k3s..."
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik" sh -

# 2. 配置 kubectl
echo "[2/5] 配置 kubectl..."
mkdir -p ~/.kube
cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
chmod 600 ~/.kube/config

# 3. 安装 kubectl（如果没有）
if ! command -v kubectl &> /dev/null; then
    echo "[3/5] 安装 kubectl..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    mv kubectl /usr/local/bin/
fi

# 4. 等待节点就绪
echo "[4/5] 等待节点就绪..."
kubectl wait --for=condition=Ready node --all --timeout=120s

# 5. 安装 Nginx Ingress Controller
echo "[5/5] 安装 Nginx Ingress Controller..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/baremetal/deploy.yaml

echo ""
echo "=========================================="
echo "  k3s 安装完成！"
echo "=========================================="
echo ""
echo "集群状态:"
kubectl get nodes
echo ""
echo "系统组件:"
kubectl get pods -A
echo ""
echo "下一步:"
echo "  1. kubectl apply -f infra/k8s/base/"
echo "  2. kubectl get pods -w"