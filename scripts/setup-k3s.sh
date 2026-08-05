#!/bin/bash
# k3s 单节点集群搭建脚本
# 用法: sudo bash setup-k3s.sh

set -e

echo "=========================================="
echo "  k3s 单节点集群搭建"
echo "=========================================="

# 1. 安装 k3s
echo "[1/4] 安装 k3s..."
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn INSTALL_K3S_EXEC="--disable=traefik" sh -

# 2. 配置 kubectl
echo "[2/4] 配置 kubectl..."
mkdir -p ~/.kube
cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
chmod 600 ~/.kube/config

# 3. 安装 kubectl（如果没有）
if ! command -v kubectl &> /dev/null; then
    echo "[3/4] 安装 kubectl..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    mv kubectl /usr/local/bin/
else
    echo "[3/4] kubectl 已存在（k3s 已自动创建 symlink）"
fi

# 4. 等待节点就绪
echo "[4/4] 等待节点就绪..."
for i in {1..60}; do
    if kubectl get nodes --no-headers 2>/dev/null | grep -q .; then
        kubectl wait --for=condition=Ready node --all --timeout=60s && break
    fi
    echo "Waiting for node to register... ($i/60)"
    sleep 5
done
kubectl get nodes

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
echo "  1. bash scripts/deploy-k8s.sh"
echo "  2. bash scripts/seed-data.sh"
echo "  3. sudo bash setup-nginx.sh"
