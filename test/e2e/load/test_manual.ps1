Write-Host "🎯 Final Correct Load Test" -ForegroundColor Green
Write-Host ""

# Test 1: PR Creation - используем одинарные кавычки!
Write-Host "1. PR Creation (30 requests)" -ForegroundColor Yellow
hey -n 30 -c 3 -m POST -H "Content-Type: application/json" -d '{"pull_request_id":"pr-correct-{{.N}}","pull_request_name":"Correct Test {{.N}}","author_id":"u1"}' "http://localhost:8080/pullRequest/create"
Write-Host ""

# Test 2: User Operations
Write-Host "2. User Operations (20 requests)" -ForegroundColor Yellow  
hey -n 20 -c 2 -m POST -H "Content-Type: application/json" -d '{"user_id":"u6","is_active":true}' "http://localhost:8080/users/setIsActive"
Write-Host ""

# Test 3: PR Merge
Write-Host "3. PR Merge (15 requests)" -ForegroundColor Yellow
hey -n 15 -c 2 -m POST -H "Content-Type: application/json" -d '{"pull_request_id":"pr-1008"}' "http://localhost:8080/pullRequest/merge"
Write-Host ""

# Test 4: Reviewer Reassignment
Write-Host "4. Reviewer Reassignment (10 requests)" -ForegroundColor Yellow
hey -n 10 -c 2 -m POST -H "Content-Type: application/json" -d '{"pull_request_id":"pr-1005","old_reviewer_id":"u3"}' "http://localhost:8080/pullRequest/reassign"
Write-Host ""

Write-Host "🎉 LOAD TESTING COMPLETED!" -ForegroundColor Green