Write-Host "🎯 Final Quick Load Test" -ForegroundColor Green
Write-Host ""

Write-Host "1. PR Creation (50 requests)" -ForegroundColor Yellow
hey -n 50 -c 5 -m POST -H "Content-Type: application/json" -d "{\"pull_request_id\":\"pr-final- { { .N } }\",\"pull_request_name\":\"Final { { .N } }\",\"author_id\":\"u1\"}" "http://localhost:8080/pullRequest/create"
Write-Host ""

Write-Host "2. User Operations (30 requests)" -ForegroundColor Yellow
hey -n 30 -c 3 -m POST -H "Content-Type: application/json" -d "{\"user_id\":\"u6\",\"is_active\":true}" "http://localhost:8080/users/setIsActive"
Write-Host ""

Write-Host "3. Mixed Stress Test (30s)" -ForegroundColor Red
hey -z 30s -c 100 -m GET "http://localhost:8080/health"
Write-Host ""

Write-Host "ALL TESTS COMPLETED!" -ForegroundColor Green
