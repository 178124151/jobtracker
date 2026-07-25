# JobTracker - 求职资源聚合平台

<div align="center">
**为求职者打造的求职资源聚合工具**

![Vue 3](https://img.shields.io/badge/Vue-3.x-brightgreen)
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?)
![Docker](https://img.shields.io/badge/Docker-Compose-blue)
![License](https://img.shields.io/badge/License-MIT-yellow)

</div>

## 功能特性

- **公司导航** — 聚合互联网公司招聘链接，支持分类筛选
- **小而美企业** — 专精特新企业库，一键直达招聘平台
- **投递进度** — 可视化管理求职投递状态
- **简历制作** — 在线编辑、实时预览、导出 PDF
- **链路监控** — Prometheus + Grafana 监控大盘

## 技术栈

| 层次 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Pinia |
| 后端 | Go + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7.2 |
| 监控 | Prometheus + Grafana + Node Exporter |
| 部署 | Docker Compose + Nginx |

## 快速开始

```bash
# 克隆项目
git clone https://github.com/yourusername/jobtracker.git
cd jobtracker

# 复制环境变量
cp .env.example .env

# 启动服务
docker compose up -d --build

# 访问
# 前端: http://localhost:5173
# API: http://localhost:8080
# Grafana: http://localhost:3000
```

## 项目结构

```
jobtracker/
├── frontend/          # Vue 3 前端
├── backend/           # Go 后端 API
├── infra/             # 基础设施配置
│   ├── nginx/         # Nginx 配置
│   ├── prometheus/    # Prometheus 配置
│   ├── grafana/       # Grafana Dashboard
│   └── k8s/           # Kubernetes 配置
├── data/              # 私有数据（不提交）
├── docker-compose.yml
├── deploy.sh          # 部署脚本
└── README.md
```

## 部署

### 服务器部署

```bash
# 1. 克隆到服务器
git clone https://github.com/yourusername/jobtracker.git /opt/jobtracker
cd /opt/jobtracker

# 2. 配置环境变量
cp .env.example .env
vim .env

# 3. 导入公司数据（需要私有数据文件）
docker cp data/seed_all_companies.sql jobtracker-postgres-1:/tmp/
docker exec -i jobtracker-postgres-1 psql -U postgres -d jobtracker -f /tmp/seed_all_companies.sql

# 4. 启动服务
docker compose up -d --build
```

### 更新部署

```bash
cd /opt/jobtracker
./deploy.sh
```

## 开源协议

MIT License

