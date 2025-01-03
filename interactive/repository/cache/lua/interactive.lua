-- cmd.Eval(ctx, setReadLua, []string{key}, "read_count", 1)

-- 获取[]string{key}中的第一个元素，key
local key = KEYS[1]

-- 获取{"read_count", 1}中的第一个元素，read_count
local keyCount = ARGV[1]

-- 获取{"read_count", 1}中的第二个元素，1
local increaseNumber = tonumber(ARGV[2])

-- 检查在redis中key是否存在
local exists = redis.call("EXISTS", key)

-- 判断key是否存在
if exists == 1 then
    -- key存在，keyCount + increaseNumber
    redis.call("HINCRBY", key, keyCount, increaseNumber)
    -- 自增成功
    return 1
else
    -- key不存在，
    -- HSET命令，用于创建map结构，创建keyCount，并初始化keyCount = increaseNumber
    redis.call("HSET", key, keyCount, increaseNumber)
    -- 自增失败
    return 0
end