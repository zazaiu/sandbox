@echo off
echo 🛑 STOPPING ALL SERVICES
echo ========================

echo.
echo 1. 🔴 Stopping application...
taskkill /f /im astro-app.exe 2>nul

echo.
echo 2. 🐳 Stopping Docker containers...
docker-compose down

echo.
echo 3. 🧹 Cleaning up...
del astro-app.exe 2>nul

echo.
echo ✅ All services stopped!