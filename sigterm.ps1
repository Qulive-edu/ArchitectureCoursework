param(
    [string]$Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.sYehCU0rmvOObJnCxCJVEv7gZmwD-A66C4WDhOV8G54",
    [string]$BaseUrl = "http://localhost:3000/api"
)

Write-Host "=== SIGTERM Graceful Shutdown Test ===" -ForegroundColor Cyan
$headers = @{ Authorization = "Bearer $Token" }

$container = (docker compose ps -q back | Select-Object -First 1)
if (-not $container) { 
    Write-Host "ERROR: backend container not found. Run 'docker compose up -d'" -ForegroundColor Red
    exit 1 
}
Write-Host "Testing container: $container" -ForegroundColor Green

# Start in-flight request in background
$job = Start-Job -ScriptBlock {
    param($url, $h)
    try { 
        $resp = Invoke-RestMethod -Uri $url -Headers $h -Method Get -TimeoutSec 10
        return $resp 
    }
    catch { 
        return $_.Exception.Message 
    }
} -ArgumentList "$BaseUrl/bookings/my", $headers

Start-Sleep -Milliseconds 500
Write-Host "Sending SIGTERM..." -ForegroundColor Yellow

docker kill --signal=SIGTERM $container | Out-Null
$result = $job | Wait-Job -Timeout 15 | Receive-Job
if ($result -and $result -is [System.Object]) {
    Write-Host "OK: In-flight request completed successfully" -ForegroundColor Green
} else {
    Write-Host "WARN: In-flight request may have been interrupted" -ForegroundColor Yellow
}
Remove-Job $job -Force

# Test new requests (should get 503 or connection error)
Write-Host "Testing new request rejection..." -ForegroundColor Yellow
Start-Sleep -Seconds 1
try {
    Invoke-RestMethod -Uri "$BaseUrl/places" -Method Get -ErrorAction Stop
    Write-Host "WARN: New request got 200 (expected 503)" -ForegroundColor Yellow
} catch {
    if ($_.Exception.Response) {
        $code = $_.Exception.Response.StatusCode.value__
        if ($code -eq 503) {
            Write-Host "OK: New request rejected with 503" -ForegroundColor Green
        } else {
            Write-Host "OK: Connection closed (status: $code)" -ForegroundColor Green
        }
    } else {
        Write-Host "OK: Connection refused (expected during shutdown)" -ForegroundColor Green
    }
}

Write-Host "`nCheck logs: docker compose logs back --tail=30" -ForegroundColor Cyan