@echo off
echo 👥 INITIALIZING TEST DATA
echo ========================

echo.
echo 1. 📦 Creating test users...
go run cmd\init_test_users.go

echo.
echo 2. 📋 Checking users in database...
docker exec pg-db psql -U astrouser -d astrodb -c "SELECT id, username, role FROM users;"

echo.
echo 3. 🪐 Checking planets...
docker exec pg-db psql -U astrouser -d astrodb -c "SELECT id, name, status FROM planets;"

echo.
echo ✅ Test data initialized!
echo    👤 User: user / 1234 (astronaut)
echo    👑 Moderator: moderator / 1234 (mission_control)