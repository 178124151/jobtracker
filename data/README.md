# Data Directory

此目录存放私有数据文件，不会提交到 Git。

## 文件说明

- `companies.json` - 大厂公司数据
- `sme_companies.json` - 小而美企业数据
- `*.sql` - 数据库导入脚本

## 使用方式

在服务器上导入数据：

```bash
docker cp data/seed_all_companies.sql jobtracker-postgres-1:/tmp/
docker exec -i jobtracker-postgres-1 psql -U postgres -d jobtracker -f /tmp/seed_all_companies.sql
```