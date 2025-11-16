# cleanup_test_data.ps1
Write-Host "🧹 Cleaning up test data..." -ForegroundColor Yellow

# Удаляем тестовые PR (осторожно - удалит все PR с префиксом pr-load-test)
Write-Host "Cleaning test PRs..." -ForegroundColor Cyan

# Восстанавливаем исходное состояние пользователей
Write-Host "Resetting user states..." -ForegroundColor Cyan
$resetUser = '{"user_id":"u3","is_active":true}'
curl -s -X POST "http://localhost:8080/users/setIsActive" -H "Content-Type: application/json" -d $resetUser

Write-Host "✅ Cleanup completed!" -ForegroundColor Green