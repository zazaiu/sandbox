@echo off
echo 🔐 TESTING AUTHENTICATION
echo ========================

echo.
echo 1. 🔑 Testing USER authentication...
curl -s -X POST http://localhost:8080/api/users/login -H "Content-Type: application/json" -d "{\"username\":\"user\",\"password\":\"1234\"}"

echo.
echo 2. 🔑 Testing MODERATOR authentication...
curl -s -X POST http://localhost:8080/api/users/login -H "Content-Type: application/json" -d "{\"username\":\"moderator\",\"password\":\"1234\"}"

echo.
echo 3. 🚫 Testing UNAUTHORIZED access...
curl -s -X GET http://localhost:8080/api/worlds -I