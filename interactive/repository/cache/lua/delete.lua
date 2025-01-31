-- 提取 KEYS 数组中的第一个元素作为 key
local key = KEYS[1]

-- 检查 key 是否存在
local exists = redis.call("EXISTS", key)

if exists == 1 then
    -- 存在则删除 key
    redis.call("DEL", key)
    -- 返回 1 表示删除成功
    return 1
else
    -- 不存在则返回 0
    return 0
end