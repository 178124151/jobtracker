# JobTracker Start Script
Set-Location $PSScriptRoot
Write-Host "Starting JobTracker..." -ForegroundColor Green
docker compose up -d
Write-Host ""
Write-Host "Access URLs:" -ForegroundColor Cyan
Write-Host "  Frontend:  http://localhost:5173" -ForegroundColor White
Write-Host "  API:       http://localhost:8080" -ForegroundColor White
Write-Host "  Grafana:   http://localhost:3000" -ForegroundColor White