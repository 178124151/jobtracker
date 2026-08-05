# JobTracker Build Script

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  JobTracker Build" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

Set-Location $PSScriptRoot

# 清理构建缓存（不删除数据卷）
Write-Host ""
Write-Host "Cleaning build cache..." -ForegroundColor Yellow
docker builder prune -f 2>&1 | Out-Null

# 重新构建
Write-Host ""
Write-Host "Building services..." -ForegroundColor Yellow
docker compose up -d --build

# 等待启动
Write-Host ""
Write-Host "Waiting for services to start (60 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 60

# 检查状态
Write-Host ""
Write-Host "==========================================" -ForegroundColor Green
Write-Host "  Status" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
docker compose ps

# 检查容器是否正常运行
Write-Host ""
$frontendStatus = docker inspect -f '{{.State.Status}}' jobtracker-frontend-1 2>&1
$backendStatus = docker inspect -f '{{.State.Status}}' jobtracker-backend-1 2>&1

if ($frontendStatus -eq "running" -and $backendStatus -eq "running") {
    Write-Host "✓ All services are running!" -ForegroundColor Green
} else {
    Write-Host "⚠ Some services failed. Restarting..." -ForegroundColor Yellow
    docker compose restart
    Start-Sleep -Seconds 30
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Green
Write-Host "  URLs" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host "Frontend:  http://localhost:5173" -ForegroundColor White
Write-Host "API:       http://localhost:8080" -ForegroundColor White
Write-Host "Grafana:   http://localhost:3000" -ForegroundColor White
Write-Host "==========================================" -ForegroundColor Cyan
