@echo off
echo 🚀 STARTING GO APPLICATION
echo ==========================

echo.
echo 1. 📦 Building application...
go build -o astro-app.exe main.go

echo.
echo 2. 🌐 Starting application...
start astro-app.exe

echo.
echo 3. ⏳ Waiting for app to start...
timeout /t 5 /nobreak

echo.
echo 🎉 Application is ready!
echo    🌐 App: http://localhost:8080
echo    📚 Swagger: http://localhost:8080/swagger/index.html