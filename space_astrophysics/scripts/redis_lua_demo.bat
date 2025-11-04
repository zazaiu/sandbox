local keys = redis.call('keys', '*')
local result = {}
local user_stats = {}

for i, key in ipairs(keys) do
    local user_id = redis.call('get', key)
    if user_id then
        -- Собираем статистику по пользователям
        user_stats[user_id] = (user_stats[user_id] or 0) + 1
        
        -- Добавляем информацию о сессии
        table.insert(result, "Session: " .. key)
        table.insert(result, "  User ID: " .. user_id)
        table.insert(result, "  TTL: " .. redis.call('ttl', key) .. " seconds")
        table.insert(result, "")
    end
end

-- Добавляем статистику в начало
local stats = {}
table.insert(stats, "=== STATISTICS ===")
table.insert(stats, "Total sessions: " .. #result)
table.insert(stats, "")

table.insert(stats, "Sessions by user:")
for user_id, count in pairs(user_stats) do
    table.insert(stats, "  User " .. user_id .. ": " .. count .. " sessions")
end

table.insert(stats, "")
table.insert(stats, "=== ACTIVE SESSIONS ===")

-- Объединяем статистику и детали
for i, line in ipairs(stats) do
    table.insert(result, 1, line)
end

return table.concat(result, '\\n')