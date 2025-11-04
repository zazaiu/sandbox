#!/bin/bash

echo "👥 INITIALIZING TEST USERS"
echo "=========================="

# Компилируем и запускаем инициализацию
go run cmd/init_test_users.go

echo ""
echo "✅ Users ready for testing:"
echo "   👤 Username: user, Password: 1234, Role: astronaut" 
echo "   👑 Username: moderator, Password: 1234, Role: mission_control"