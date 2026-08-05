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
| 生产部署 | Kubernetes (k3s) + 宿主机 Nginx |
| 本地开发 | Docker Compose（可选） |

## 生产部署（k3s，只暴露 80 端口）

生产环境统一使用 k3s 容器，不依赖 Docker Compose 提供服务。外部请求统一从宿主机 Nginx 的 80 端口进入：

```text
浏览器
  └─ :80 宿主机 Nginx
       ├─ /        -> frontend NodePort 30080  -> 前端静态页 + SPA 路由
       ├─ /api/    -> backend  NodePort 30081  -> Go API
       └─ /grafana -> grafana NodePort 30300   -> 监控大盘
```

1. 初始化 k3s（不需要 Ingress Controller，宿主机 Nginx 负责转发）：

```bash
cd /opt/jobtracker
sudo bash scripts/setup-k3s.sh
```

2. 构建镜像并部署全部 K8s 资源：

```bash
bash scripts/deploy-k8s.sh
```

3. 初始化数据库和种子数据：

```bash
bash scripts/seed-data.sh
```

`seed-data.sh` 已内置完整公司备份；普通模式只在空表时导入，需要覆盖恢复时执行 `bash scripts/seed-data.sh --restore`。

4. 配置宿主机 Nginx + UFW，只开放 22 和 80：

```bash
sudo bash setup-nginx.sh
```

之后访问：

```text
前端:   http://<服务器IP>/
API:    http://<服务器IP>/api/
Grafana: http://<服务器IP>/grafana/
```

K8s 里的 NodePort（30080/30081/30300/30090）以及 PostgreSQL/Redis 端口都会被 UFW 拒绝公网访问，只允许从宿主机本机转发。

## 安全注意事项

- 不要暴露 PostgreSQL（5432）、Redis（6379）以及任何 NodePort 到公网。
- 不要在服务器上执行 `kubectl port-forward --address 0.0.0.0 ... 5432:5432` 这类命令。
- 如果数据库里出现名为 `readme_to_recover` 的库，说明 PostgreSQL 已被公网扫描脚本攻击：立即关闭 5432 公网入口、更换数据库口令和后端 JWT Secret、从备份恢复数据，并按“已入侵”检查集群和宿主机。
- 数据目录 `data/` 属于私有数据，不会被提交到 Git；请定期备份数据库。

## 项目结构

```text
jobtracker/
├── frontend/          # Vue 3 前端应用
├── backend/           # Go 后端 API
├── infra/             # Nginx / Prometheus / Grafana / K8s 配置
├── scripts/           # k3s 部署、种子数据等脚本
├── data/              # 私有数据（不入库）
├── docker-compose.yml # 仅本地开发
├── setup-nginx.sh     # 宿主机 Nginx + UFW（只开放 80）
└── README.md
```

## 本地开发（Docker，可选）

```bash
git clone https://github.com/yourusername/jobtracker.git
cd jobtracker

cp .env.example .env
vim .env
./build.sh

# 前端: http://localhost:5173
# API: http://localhost:8080
# Grafana: http://localhost:3000
```

`docker-compose.yml` 的默认口令来自 `.env`，未设置时使用 `.env.example` 里的随机默认值。生产部署前请务必替换 `DB_PASSWORD`、`POSTGRES_PASSWORD`、`JWT_SECRET` 和 `GRAFANA_PASSWORD`。

## 数据备份

K8s 部署：

```bash
kubectl exec deploy/postgres -- pg_dump -U postgres -d jobtracker > backup.sql
```

建议通过 cron 定时执行，并把备份文件放到独立存储中。

## 开源协议
MIT License
