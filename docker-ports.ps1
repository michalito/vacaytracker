# VacayTracker Docker Port Helper
# Displays the randomly assigned ports for running containers
#
# Usage:
#   ./docker-ports.ps1           # Show dev mode ports
#   ./docker-ports.ps1 -Mode prod  # Show prod mode ports

param(
    [string]$Mode = "dev"
)

Write-Host ""
Write-Host "VacayTracker Docker Ports" -ForegroundColor Cyan
Write-Host "=========================" -ForegroundColor Cyan

if ($Mode -eq "prod") {
    $frontendPort = 3000
    $file = "docker-compose.prod.yml"
    Write-Host "Mode: Production" -ForegroundColor Yellow
} else {
    $frontendPort = 5173
    $file = "docker-compose.yml"
    Write-Host "Mode: Development" -ForegroundColor Yellow
}

Write-Host ""

# Get frontend port
$frontend = docker compose -f $file port frontend $frontendPort 2>$null
if ($frontend) {
    Write-Host "  Frontend: " -NoNewline
    Write-Host "http://$frontend" -ForegroundColor Green
} else {
    Write-Host "  Frontend: " -NoNewline
    Write-Host "Not running" -ForegroundColor Red
}

# Get API port
$api = docker compose -f $file port api 3000 2>$null
if ($api) {
    Write-Host "  API:      " -NoNewline
    Write-Host "http://$api" -ForegroundColor Green
} else {
    Write-Host "  API:      " -NoNewline
    Write-Host "Not running" -ForegroundColor Red
}

Write-Host ""
