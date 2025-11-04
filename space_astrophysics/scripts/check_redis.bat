@echo off
echo 🔍 CHECKING REDIS SESSIONS
echo ==========================

echo.
echo 1. 📊 All keys in Redis:
docker exec redis redis-cli keys "*"

echo.
echo 2. 👥 Active user sessions:
docker exec redis redis-cli eval "local keys = redis.call('keys', '*') for i, key in ipairs(keys) do local user_id = redis.call('hget', key, 'user_id') local role = redis.call('hget', key, 'role') if user_id then print('🔑 ' .. key) print('   👤 User: ' .. user_id) print('   🎯 Role: ' .. role) print('   ⏰ TTL: ' .. redis.call('ttl', key)) end end" 0