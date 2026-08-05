# JobTracker - 求职资源聚合平台

<div align="center">
**为求职者打造的求职资源聚合工具**

![Vue 3](https://img.shields.io/badge/Vue-3.x-brightgreen)
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?)
![Docker](https://img.shields.io/badge/Docker-Compose-blue)
![License](https://img.shields.io/badge/License-MIT-yellow)

</div>

## 功能特性
- 公司导航：聚合互联网公司招聘链接，支持分类筛选
- 小而美企业：精选创新企业库，一键直达招聘平台
- 投递进度：可视化管理求职投递状态
- 简历制作：在线编辑、实时预览、导出 PDF
- 链路监控：Prometheus + Grafana 监控大盘

## 技术栈

| 层次 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Pinia |
| 后端 | Go + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7.2 |
| 监控 | Prometheus + Grafana + Node Exporter |
| 部署 | Docker Compose / Kubernetes (k3s) |

## 快速开始（Docker）

```bash
# 克隆项目
git clone https://github.com/yourusername/jobtracker.git
cd jobtracker

# 复制环境变量并修改口令
cp .env.example .env
vim .env

# 构建并启动
./build.sh

# 访问
# 前端: http://localhost:5173
# API: http://localhost:8080
# Grafana: http://localhost:3000
```

`docker-compose.yml` 中的默认口令来自 `.env`，未设置时会使用 `.env.example` 里的随机默认值。**生产环境部署前请务必替换 `DB_PASSWORD`、`POSTGRES_PASSWORD`、`JWT_SECRET` 和 `GRAFANA_PASSWORD`。**

PostgreSQL 端口已只绑定 `127.0.0.1:5432`，不要把它改成 `0.0.0.0:5432` 或直接暴露到公网。

## 安全注意事项

- 永远不要把 PostgreSQL（5432）、Redis（6379）端口暴露到公网。
- 不要在服务器上执行 `kubectl port-forward --address 0.0.0.0 ... 5432:5432` 这类命令。
- 如果数据库里出现名为 `readme_to_recover` 的库，说明 PostgreSQL 已被公网扫描脚本攻击：立即关闭 5432 公网入口、更换数据库口令和后端 JWT Secret、从备份恢复数据，并按“已入侵”检查集群和宿主机。
- 数据目录 `data/` 属于私有数据，不会被提交到 Git；请定期备份数据库。

## 项目结构

```
jobtracker/
├── frontend/          # Vue 3 前端应用
├── backend/           # Go 后端 API
├── infra/             # 基础设施配置（Nginx / Prometheus / Grafana / K8s）
├── scripts/           # k3s 部署、种子数据等脚本
├── data/              # 私有数据（不入库）
├── docker-compose.yml
├── build.sh           # 构建并启动（不删除数据卷）
├── deploy.sh          # Docker 更新部署
├── setup-nginx.sh     # Nginx + UFW 防火墙配置
└── README.md
```

## 部署

### Docker 部署

```bash
cd /opt/jobtracker
cp .env.example .env
vim .env

# 首次构建
./build.sh

# 日常更新（拉取代码 + 重新构建）
./deploy.sh
```

`build.sh` 不会删除数据卷，日常重建不会丢数据。只有确实需要完全重置时才手动执行 `docker compose down -v`。

### Kubernetes / k3s 部署

```bash
cd /opt/jobtracker
bash scripts/deploy-k8s.sh
```

脚本会拉取最新代码、构建并导入镜像到 k3s、更新 `sme-data` ConfigMap，然后应用 K8s 清单。首次部署或数据库被清空后，需要初始化数据：

```bash
bash scripts/seed-data.sh
```

`seed-data.sh` 会：
1. 确保 `jobtracker` 数据库和表结构存在；
2. 在 `companies` 表为空时导入公司数据；
3. 把 `data/sme_companies.json` 写入 `sme-data` ConfigMap 并滚动重启后端，避免数据随 Pod 重启丢失。

K8s 环境下建议给 PostgreSQL 配置 PVC，并保留 `POSTGRES_DB`、`POSTGRES_USER`、`POSTGRES_PASSWORD` 环境变量，否则 Pod 重建会丢失全部数据。

## 数据备份

Docker 部署：

```bash
docker exec -i jobtracker-postgres-1 pg_dump -U postgres -d jobtracker > backup.sql
```

K8s 部署：

```bash
kubectl exec deploy/postgres -- pg_dump -U postgres -d jobtracker > backup.sql
```

建议通过 cron 定时执行，并把备份文件放到独立存储中。

## 开源协议
MIT License
