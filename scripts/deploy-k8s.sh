#!/bin/bash
# JobTracker K8s Deploy Script
set -e

if ! command -v k3s >/dev/null 2>&1; then
    echo "ERROR: k3s not found. Please run 'sudo bash scripts/setup-k3s.sh' first."
    exit 1
fi

if ! sudo k3s ctr images list | grep -q "rancher/mirrored-pause:3.6"; then
    echo "Importing pause image..."
    sudo docker pull rancher/mirrored-pause:3.6
    sudo docker save rancher/mirrored-pause:3.6 -o /tmp/pause.tar
    sudo k3s ctr images import /tmp/pause.tar
    rm -f /tmp/pause.tar
fi

echo "=========================================="
echo "  JobTracker K8s Deploy"
echo "=========================================="

cd /opt/jobtracker

# 使用当前本地代码构建；需要更新代码时请手动执行 git pull origin main

# 1. 拉取镜像
echo "[1/6] Pulling images..."
sudo docker pull postgres:16-alpine
sudo docker pull redis:7.2-alpine
sudo docker pull grafana/grafana:latest

# 2. 构建应用镜像
echo "[2/6] Building application images... (first build may take a few minutes)"
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

# 更新 SME 数据 ConfigMap（数据文件不入库）
if [ -f data/sme_companies.json ]; then
  kubectl create configmap sme-data --from-file=data/sme_companies.json --dry-run=client -o yaml | kubectl apply --server-side --force-conflicts -f - || echo "Warning: failed to update sme-data ConfigMap"
else
  echo "Warning: data/sme_companies.json not found, skipping SME ConfigMap"
fi

kubectl wait --for=condition=Ready pod -l component=database --timeout=120s || true
sleep 15
kubectl apply -f infra/k8s/base/backend.yaml
kubectl apply -f infra/k8s/base/frontend.yaml
kubectl apply -f infra/k8s/base/hpa.yaml

kubectl rollout status deployment/jobtracker-backend --timeout=120s || true

# 6. 等待启动
echo "[6/6] Waiting for pods..."
sleep 30
kubectl get pods

echo ""
echo "=========================================="
echo "  Deploy Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  1. bash scripts/seed-data.sh"
echo "  2. sudo bash setup-nginx.sh"
echo ""
echo "After that, only port 80 is exposed:"
echo "  Frontend: http://<server-ip>/"
echo "  API:      http://<server-ip>/api/"
echo "  Grafana:  http://<server-ip>/grafana/"
