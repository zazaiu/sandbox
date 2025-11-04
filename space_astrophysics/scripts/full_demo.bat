@echo off
echo 🎯 FULL DEMONSTRATION - LAB 4
echo =============================

echo.
echo 📋 STEP 1: Starting Docker services...
call scripts\start_services.bat

echo.
echo 📋 STEP 2: Initializing test data...
call scripts\init_test_data.bat

echo.
echo 📋 STEP 3: Starting application...
call scripts\start_app.bat

echo.
echo 📋 STEP 4: Testing authentication...
timeout /t 3 /nobreak
call scripts\test_auth.bat

echo.
echo 📋 STEP 5: Checking Redis sessions...
call scripts\check_redis.bat

echo.
echo 🎉 DEMONSTRATION COMPLETED!
echo.
echo 📊 Next steps:
echo    1. Open Swagger: http://localhost:8080/swagger/index.html
echo    2. Test authentication in browser
echo    3. Check cookies in Developer Tools
echo    4. Use Redis commands to view sessions