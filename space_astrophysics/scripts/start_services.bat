@echo off
echo 🚀 STARTING ALL SERVICES
echo ========================

echo.
echo 1. 🧹 Stopping existing services...
docker-compose down

echo.
echo 2. 🐳 Starting Docker containers...
docker-compose up -d

echo.
echo 3. ⏳ Waiting for services to start...
timeout /t 10 /nobreak

echo.
echo 4. ✅ Checking services status...
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo.
echo 🎉 Services are ready!
echo    📊 PostgreSQL: localhost:5433
echo    🔐 Redis: localhost:6379
echo    📁 MinIO: http://localhost:9001
echo    📋 Adminer: http://localhost:8081