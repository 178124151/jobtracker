#!/bin/bash
# JobTracker Seed Data Import Script
set -e

echo "=== Importing Seed Data ==="

# 等待 PostgreSQL 就绪
echo "Waiting for PostgreSQL..."
sleep 10

# 获取 PostgreSQL Pod 名称
PG_POD=$(kubectl get pod -l component=database -o jsonpath="{.items[0].metadata.name}" 2>/dev/null)

if [ -z "$PG_POD" ]; then
    echo "ERROR: PostgreSQL pod not found"
    exit 1
fi

echo "PostgreSQL pod: $PG_POD"

# 检查是否已有数据
COUNT=$(kubectl exec $PG_POD -- psql -U postgres -d jobtracker -t -c "SELECT COUNT(*) FROM companies;" 2>/dev/null | tr -d ' ')

if [ "$COUNT" -gt 0 ]; then
    echo "Database already has $COUNT companies, skipping seed."
    exit 0
fi

# 导入大厂公司数据
echo "Importing big tech companies..."
kubectl exec $PG_POD -- psql -U postgres -d jobtracker -c "
INSERT INTO companies (id, name, website, industry, \"group\", description, is_preset, health_status) VALUES
(gen_random_uuid(), '字节跳动', 'https://jobs.bytedance.com', '互联网', 'bigtech', '短视频和信息流平台', true, 'GREEN'),
(gen_random_uuid(), '阿里巴巴', 'https://talent.alibaba.com', '互联网', 'bigtech', '电商和云计算集团', true, 'GREEN'),
(gen_random_uuid(), '腾讯', 'https://careers.tencent.com', '互联网', 'bigtech', '社交和游戏科技集团', true, 'GREEN'),
(gen_random_uuid(), '美团', 'https://zhaopin.meituan.com', '互联网', 'bigtech', '本地生活服务平台', true, 'GREEN'),
(gen_random_uuid(), '京东', 'https://zhaopin.jd.com', '电商', 'bigtech', '电商和物流集团', true, 'GREEN'),
(gen_random_uuid(), '百度', 'https://talent.baidu.com', '互联网', 'bigtech', '搜索引擎和AI公司', true, 'GREEN'),
(gen_random_uuid(), '快手', 'https://zhaopin.kuaishou.cn', '短视频', 'bigtech', '短视频社区平台', true, 'GREEN'),
(gen_random_uuid(), '小红书', 'https://job.xiaohongshu.com', '社交', 'bigtech', '生活方式分享平台', true, 'GREEN'),
(gen_random_uuid(), '拼多多', 'https://careers.pddglobalhr.com', '电商', 'bigtech', '新电商平台', true, 'GREEN'),
(gen_random_uuid(), 'bilibili', 'https://jobs.bilibili.com', '视频', 'bigtech', '年轻人文化社区', true, 'GREEN'),
(gen_random_uuid(), '网易', 'https://hr.163.com', '互联网', 'bigtech', '游戏和音乐平台', true, 'GREEN'),
(gen_random_uuid(), '携程', 'https://job.ctrip.com', '旅游', 'bigtech', '在线旅行服务平台', true, 'GREEN'),
(gen_random_uuid(), '滴滴出行', 'https://talent.didiglobal.com', '出行', 'bigtech', '出行服务平台', true, 'GREEN'),
(gen_random_uuid(), '知乎', 'https://www.zhihu.com', '知识', 'bigtech', '问答社区平台', true, 'GREEN'),
(gen_random_uuid(), '米哈游', 'https://www.mihoyo.com', '游戏', 'bigtech', '原神开发商', true, 'GREEN'),
(gen_random_uuid(), '微博', 'https://career.sina.com.cn', '社交媒体', 'bigtech', '社交媒体平台', true, 'GREEN'),
(gen_random_uuid(), '搜狐', 'https://hr.sohu.com', '互联网', 'bigtech', '门户网站', true, 'GREEN'),
(gen_random_uuid(), '小米', 'https://hr.xiaomi.com', '智能硬件', 'bigtech', '智能手机和IoT', true, 'GREEN'),
(gen_random_uuid(), '华为', 'https://career.huawei.com', '科技', 'bigtech', 'ICT基础设施', true, 'GREEN'),
(gen_random_uuid(), '荣耀', 'https://career.honor.com', '智能硬件', 'bigtech', '智能手机品牌', true, 'GREEN');
"

# 检查导入结果
NEW_COUNT=$(kubectl exec $PG_POD -- psql -U postgres -d jobtracker -t -c "SELECT COUNT(*) FROM companies;" 2>/dev/null | tr -d ' ')
echo "Imported $NEW_COUNT companies successfully!"

# 将 SME 数据复制到后端 Pod
echo "Setting up SME data..."
BACKEND_POD=$(kubectl get pod -l component=backend -o jsonpath="{.items[0].metadata.name}" 2>/dev/null)
if [ -n "$BACKEND_POD" ]; then
    kubectl cp data/sme_companies.json $BACKEND_POD:/root/data/sme_companies.json
    echo "SME data copied to backend pod"
fi

echo ""
echo "=== Seed Data Import Complete ==="