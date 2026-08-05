# JobTracker Runbook

## 症状

后端 Pod 显示 `0/1 Running`，`/readyz` 返回 503，`backend-service` 没有 Ready endpoints，前端所有 `/api` 请求失败。

## 快速排查

```bash
kubectl get pods -l component=backend
kubectl describe pod -l component=backend | grep -A5 "Readiness"
kubectl logs -l component=backend --tail=200
kubectl get endpoints backend-service
kubectl exec deploy/jobtracker-backend -- wget -qO- http://127.0.0.1:8080/readyz
```

`/readyz` 检查的是 PostgreSQL 连通性，所以 NotReady 通常意味着数据库不可用、库不存在或后端连不上数据库。

## 数据库检查

```bash
kubectl exec deploy/postgres -- psql -U postgres -l
kubectl exec deploy/postgres -- psql -U postgres -d jobtracker -c "SELECT count(*) FROM companies;"
```

如果 `jobtracker` 库不存在，先确认是不是被公网扫描脚本删除：

```bash
kubectl exec deploy/postgres -- psql -U postgres -c "SELECT datname FROM pg_database ORDER BY 1;"
kubectl exec deploy/postgres -- psql -U postgres -d readme_to_recover -c "\d readme"
```

出现 `readme_to_recover` 库和 `readme` 表说明 PostgreSQL 已被攻击：不要付款、不要执行其内容，立即进入应急处理。

## 应急处理

1. 关闭 5432/6379 公网入口（云安全组、UFW、NodePort、kubectl port-forward）。
2. 更换 PostgreSQL 口令和后端 JWT Secret，更新 `.env` 或 K8s Secret 后重启后端。
3. 从备份恢复数据（`pg_dump` 备份文件）。
4. 检查集群和宿主机是否被植入后门：

```bash
kubectl get pods -A
kubectl get cronjobs -A
sudo docker ps -a
sudo crontab -l
cat ~/.ssh/authorized_keys
```

5. 恢复数据库后重启后端并确认 Ready：

```bash
kubectl exec deploy/postgres -- psql -U postgres -c "CREATE DATABASE jobtracker;"
cat backend/migrations/init.sql | kubectl exec -i deploy/postgres -- psql -U postgres -d jobtracker
cat data/seed_all_companies.sql | kubectl exec -i deploy/postgres -- psql -U postgres -d jobtracker
kubectl rollout restart deployment/jobtracker-backend
kubectl get pods -l component=backend
```
