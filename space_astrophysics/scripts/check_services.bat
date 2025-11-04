#!/bin/bash

echo "🔍 CHECKING DOCKER SERVICES"
echo "==========================="

echo ""
echo "1. 🐳 Checking running containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "2. 🗄️ Checking PostgreSQL connection:"
docker exec pg-db psql -U astrouser -d astrodb -c "SELECT version();" 2>/dev/null && echo "✅ PostgreSQL is running" || echo "❌ PostgreSQL connection failed"

echo ""
echo "3. 🔐 Checking Redis connection:"
docker exec redis redis-cli ping 2>/dev/null && echo "✅ Redis is running" || echo "❌ Redis connection failed"

echo ""
echo "4. 📁 Checking MinIO:"
curl -s http://localhost:9001/minio/health/live > /dev/null && echo "✅ MinIO is running" || echo "❌ MinIO connection failed"

echo ""
echo "5. 📊 Checking Adminer:"
curl -s http://localhost:8081 > /dev/null && echo "✅ Adminer is running" || echo "❌ Adminer connection failed"