#!/bin/bash

# JobTracker Health Checker
# 检查公司官网可达性

COMPANIES_FILE="data/companies.json"
RESULTS_FILE="data/health_results.json"

# 检查单个URL
check_url() {
    local url=$1
    local timeout=10
    
    start_time=$(date +%s%N)
    response=$(curl -s -o /dev/null -w "%{http_code}" --max-time $timeout "$url" 2>/dev/null)
    end_time=$(date +%s%N)
    
    duration=$(( (end_time - start_time) / 1000000 ))
    
    if [ "$response" -ge 200 ] && [ "$response" -lt 400 ]; then
        echo "GREEN $duration"
    else
        echo "RED $duration"
    fi
}

# 主函数
main() {
    echo "Starting health check..."
    echo "[]" > "$RESULTS_FILE"
    
    # 从JSON文件读取公司列表并检查
    # 这里简化处理，实际应该解析JSON
    urls=(
        "https://talent.alibaba.com"
        "https://careers.tencent.com"
        "https://jobs.bytedance.com"
        "https://talent.baidu.com"
        "https://zhaopin.meituan.com"
    )
    
    for url in "${urls[@]}"; do
        result=$(check_url "$url")
        status=$(echo $result | awk '{print $1}')
        duration=$(echo $result | awk '{print $2}')
        
        echo "  $url: $status (${duration}ms)"
    done
    
    echo "Health check completed."
}

main "$@"
