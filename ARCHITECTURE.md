# JobTracker 架构设计说明

## 核心架构与数据分离设计

本项目采用**核心代码与数据分离**的设计理念，便于开源核心功能，同时保护私有数据。

---

## 目录结构说明

```
jobtracker/
├── frontend/              # ✅ 开源 - Vue3前端应用
├── backend/               # ✅ 开源 - Go后端服务
├── infra/                 # ✅ 开源 - 基础设施配置
├── data/                  # ❌ 私有 - 公司数据库（.gitignore）
│   ├── .gitkeep           # 保持目录结构
│   └── companies.json     # 公司数据（不提交git）
├── .gitignore             # 排除私有数据
├── docker-compose.yml     # ✅ 开源 - Docker编排
├── .env.example           # ✅ 开源 - 环境变量模板
├── README.md              # ✅ 开源 - 项目说明
└── LICENSE                # ✅ 开源 - MIT协议
```

---

## 开源内容清单

### 前端代码（frontend/）
- Vue 3 + Vite + TypeScript
- 组件化架构：Sidebar、Topbar、CompanyCard等
- Pinia状态管理
- Vue Router路由配置

### 后端代码（backend/）
- Go + Gin框架
- 标准分层架构：Handler → Service → Repository
- PostgreSQL + GORM ORM
- JWT认证中间件
- 数据库迁移脚本

### 基础设施（infra/）
- Nginx反向代理配置
- Prometheus监控配置
- Grafana Dashboard模板
- Kubernetes部署配置

---

## 私有数据保护

### 公司数据库（data/companies.json）
包含公司详细信息，不提交到GitHub：
- 公司名称、官网链接
- 行业分类、分组信息
- 招聘页面URL
- 健康检查状态

### 环境变量（.env）
敏感配置信息：
- 数据库密码
- JWT密钥
- API密钥

---

## 部署架构

### 本地开发
```
Browser → Vite Dev Server (5173)
              ↓ (proxy)
         Go API Server (8080)
              ↓
         PostgreSQL + Redis
```

### Docker部署
```
Browser → Nginx (80)
              ↓
    ┌────────┴────────┐
    │                 │
Frontend (80)    Backend (8080)
    │                 │
    └────────┬────────┘
             ↓
      PostgreSQL + Redis
             ↓
    Prometheus + Grafana
```

### Kubernetes部署
```
Browser → Ingress (nginx)
              ↓
    ┌────────┴────────┐
    │                 │
Frontend Service  Backend Service
    │                 │
    └────────┬────────┘
             ↓
      PostgreSQL + Redis (StatefulSet)
```

---

## 数据分离实现

### 1. Git忽略规则
```gitignore
# Data directory (private company data)
data/*.json
!data/.gitkeep
```

### 2. 数据加载机制
后端从`data/companies.json`加载公司数据，但该文件不提交到Git。

用户首次使用时需要：
1. 复制示例数据或创建自己的公司数据
2. 将数据文件放在`data/`目录下

### 3. 环境变量分离
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑.env文件，填入实际配置
vim .env
```

---

## 开源贡献指南

### 添加新功能
1. Fork项目
2. 创建功能分支
3. 提交PR到main分支
4. 通过CI检查后合并

### 数据贡献
公司数据通过单独的PR提交，需人工审核：
1. 在`data/companies.example.json`中添加示例
2. 说明数据来源
3. 等待维护者审核

---

## 技术栈总结

| 层次 | 技术 | 版本 |
|------|------|------|
| 前端框架 | Vue 3 | ≥3.4 |
| 构建工具 | Vite | ≥5.0 |
| 状态管理 | Pinia | ≥2.1 |
| 后端语言 | Go | ≥1.22 |
| Web框架 | Gin | ≥1.9 |
| ORM | GORM | ≥2.0 |
| 数据库 | PostgreSQL | ≥16 |
| 缓存 | Redis | ≥7.2 |
| 容器化 | Docker + Compose | - |
| 监控 | Prometheus + Grafana | - |

---

## 项目优势

1. **架构清晰**：前后端分离，标准分层
2. **易于扩展**：模块化设计，便于添加新功能
3. **部署灵活**：支持Docker、K8s多种部署方式
4. **数据安全**：核心代码与私有数据分离
5. **监控完善**：内置Prometheus + Grafana监控

---

## 下一步计划

- [ ] 实现JWT认证完整流程
- [ ] 添加公司健康检查定时任务
- [ ] 完善Grafana Dashboard
- [ ] 编写CI/CD流水线
- [ ] 添加单元测试
