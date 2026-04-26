# test-disposability-v2.ps1
# Fixed: proper escaping, job output handling, dynamic slot selection

param(
    [string]$Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.sYehCU0rmvOObJnCxCJVEv7gZmwD-A66C4WDhOV8G54",
    [string]$BaseUrl = "http://localhost:3000/api"
)

Write-Host "=== Disposability Test v2 ===" -ForegroundColor Cyan
$headers = @{ Authorization = "Bearer $Token" }

# === 1. PRE-CHECK ===
$containers = docker compose ps -q back
$replicaCount = ($containers | Measure-Object).Count
Write-Host "Backend containers running: $replicaCount" -ForegroundColor $(if($replicaCount -ge 1){"Green"}else{"Red"})

if ($replicaCount -eq 0) {
    Write-Host "ERROR: Start with 'docker compose up -d'" -ForegroundColor Red
    exit 1
}
$targetContainer = $containers[0]

# === 2. LOAD GENERATOR ===
Write-Host "Starting load generator..." -ForegroundColor Gray

$loadJob = Start-Job -ScriptBlock {
    param($url, $h, $durationSec)
    
    $results = @{ success = 0; transition = 0; fail = 0 }
    $endTime = (Get-Date).AddSeconds($durationSec)
    
    while ((Get-Date) -lt $endTime) {
        try {
            $null = Invoke-RestMethod -Uri $url -Headers $h -Method Get -TimeoutSec 3
            $results.success++
        }
        catch {
            $code = 0
            if ($_.Exception.Response) { $code = $_.Exception.Response.StatusCode.value__ }
            
            if ($code -in @(499, 502, 0, 503)) {
                $results.transition++
            } else {
                $results.fail++
            }
        }
        Start-Sleep -Milliseconds 250
    }
    return $results
} -ArgumentList "$BaseUrl/places", $headers, 12

Start-Sleep -Seconds 2

# === 3. KILL CONTAINER ===
Write-Host "Killing container: $($targetContainer.Substring(0,12))..." -ForegroundColor Red
$killTime = Get-Date
docker kill --signal=SIGTERM $targetContainer 2>$null | Out-Null
Start-Sleep -Seconds 15

# === 4. COLLECT RESULTS ===
$loadJob | Wait-Job -Timeout 20 | Out-Null
$loadResults = $loadJob | Receive-Job
Remove-Job $loadJob -Force

Write-Host "`n--- Load Test Results ---" -ForegroundColor Cyan
Write-Host "Successful: $($loadResults.success)" -ForegroundColor Green
Write-Host "Transition (499/502/503): $($loadResults.transition)" -ForegroundColor Yellow  
Write-Host "Unexpected failures: $($loadResults.fail)" -ForegroundColor $(if($loadResults.fail -eq 0){"Green"}else{"Red"})

# === 5. RECOVERY CHECK ===
Write-Host "`n--- Recovery Check ---" -ForegroundColor Cyan
$recoveryStart = Get-Date
$maxWait = 45
$recovered = $false

while (((Get-Date) - $recoveryStart).TotalSeconds -lt $maxWait) {
    try {
        $probe = Invoke-WebRequest -Uri "$BaseUrl/places" -Method Get -TimeoutSec 3 -ErrorAction Stop
        if ($probe.StatusCode -eq 200) {
            $recovered = $true
            break
        }
    } catch {
        # Keep waiting
    }
    Start-Sleep -Seconds 1
}

$recoveryTime = [math]::Round(((Get-Date) - $recoveryStart).TotalSeconds, 1)
$killToRecovery = [math]::Round(((Get-Date) - $killTime).TotalSeconds, 1)

if ($recovered) {
    Write-Host "Service recovered in ${recoveryTime}s (kill->ready: ${killToRecovery}s)" -ForegroundColor Green
} else {
    Write-Host "Service did NOT recover within ${maxWait}s" -ForegroundColor Red
}

# === 6. DATA INTEGRITY ===
Write-Host "`n--- Data Integrity Check ---" -ForegroundColor Cyan

try {
    $slotsResp = Invoke-RestMethod -Uri "$BaseUrl/places/1/slots" -Headers $headers -Method Get
    $availableSlot = $slotsResp | Where-Object { $_.is_available -eq $true } | Select-Object -First 1
    
    if (-not $availableSlot) {
        Write-Host "No available slots for testing" -ForegroundColor Yellow
    } else {
        $bookingBody = @{
            place_id = 1
            slot_id = $availableSlot.id
        } | ConvertTo-Json
        
        $bookingResp = Invoke-RestMethod -Uri "$BaseUrl/bookings" -Method Post `
            -Headers $headers -ContentType "application/json" `
            -Body $bookingBody -ErrorAction Stop
        
        Write-Host "Created booking for slot_id=$($availableSlot.id)" -ForegroundColor Green
        
        $myBookings = Invoke-RestMethod -Uri "$BaseUrl/bookings/my" -Headers $headers -Method Get
        $found = $myBookings | Where-Object { $_.id -eq $bookingResp.id }
        
        if ($found) {
            Write-Host "Data integrity: Booking persisted in PostgreSQL" -ForegroundColor Green
        } else {
            Write-Host "Data integrity: Booking not in list" -ForegroundColor Yellow
        }
    }
}
catch {
    Write-Host "Data integrity check: $($_.Exception.Message)" -ForegroundColor Yellow
}

# === 7. FINAL VERDICT ===
Write-Host "`n=== Test Complete ===" -ForegroundColor Cyan

$pass = $recovered -and $loadResults.fail -eq 0
if ($pass) {
    Write-Host "RESULT: PASS - Disposability requirements met" -ForegroundColor Green
} else {
    Write-Host "RESULT: REVIEW" -ForegroundColor Yellow
    if (-not $recovered) { Write-Host "  - Recovery timeout" }
    if ($loadResults.fail -gt 0) { Write-Host "  - Unexpected failures: $($loadResults.fail)" }
}

Write-Host "`nQuick checks:" -ForegroundColor Gray
Write-Host "  - Container status: docker compose ps back" -ForegroundColor Gray
Write-Host "  - Logs: docker compose logs back --tail=30" -ForegroundColor Gray
Write-Host "  - DB count: docker compose exec postgres psql -U postgres -c 'SELECT COUNT(*) FROM bookings;'" -ForegroundColor Gray