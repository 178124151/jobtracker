#!/bin/bash
# Nginx 反向代理配置脚本
# 在服务器上执行: sudo bash setup-nginx.sh

set -e

echo "=== 配置 Nginx 反向代理 ==="

# 1. 安装 Nginx
echo "[1/5] 安装 Nginx..."
apt update -qq && apt install -y nginx

# 2. 复制配置文件
echo "[2/5] 写入 Nginx 配置..."
cat > /etc/nginx/sites-available/jobtracker << 'EOF'
server {
    listen 80;
    server_name _;

    # 前端（k3s NodePort）
    location / {
        proxy_pass http://127.0.0.1:30080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后端 API（k3s NodePort）
    location /api/ {
        proxy_pass http://127.0.0.1:30081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Grafana
    location /grafana/ {
        proxy_pass http://127.0.0.1:30300;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Grafana WebSocket
    location /grafana/api/live/ {
        proxy_pass http://127.0.0.1:30300;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
EOF

# 3. 启用配置
echo "[3/5] 启用站点配置..."
ln -sf /etc/nginx/sites-available/jobtracker /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# 4. 测试并重启
echo "[4/5] 测试 Nginx 配置..."
nginx -t

echo "[5/5] 重启 Nginx..."
systemctl restart nginx
systemctl enable nginx

# 5. 配置防火墙：只开放 80 和 22
echo ""
echo "=== 配置防火墙 ==="
ufw allow 22/tcp
ufw allow 80/tcp
ufw deny 5173/tcp
ufw deny 8080/tcp
ufw deny 3000/tcp
ufw deny 9090/tcp
ufw deny 5432/tcp
ufw deny 6379/tcp
ufw deny 30080/tcp
ufw deny 30081/tcp
ufw deny 30300/tcp
ufw deny 30090/tcp
ufw deny 9100/tcp
ufw --force enable

echo ""
echo "=== 配置完成 ==="
echo "现在只能通过 80 端口访问："
echo "  前端: http://服务器IP/"
echo "  API:  http://服务器IP/api/"
echo "  监控: http://服务器IP/grafana/"
echo ""
echo "已关闭的端口: 5173, 8080, 3000, 9090, 9100, 5432, 6379, 30080, 30081, 30300, 30090"
echo "警告: 不要把 PostgreSQL/Redis 端口暴露到公网，否则可能被扫描脚本攻击"
