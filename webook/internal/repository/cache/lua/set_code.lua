-- 验证码在 redis 上的 key
local key = KEYS[1]
-- 使用次数 验证次数。最多重复三次
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
-- key存在但没有过期时间
if ttl == -1 then
    return -1
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    -- 符合预期
    return 0
else
    -- 发送太频繁
    return -2
end