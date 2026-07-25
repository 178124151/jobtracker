# JobTracker Stop Script
Set-Location $PSScriptRoot
Write-Host "Stopping JobTracker..." -ForegroundColor Yellow
docker compose stop
Write-Host "Done!" -ForegroundColor Green