local key = KEYS[1]
-- 用户输入的code
local expectedCode = ARGV[1]
-- redis 存的真正code
local code = redis.call("get", key)
local cntKey = key..":cnt"
-- 转成一个数字
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 一直输错
    return -2
elseif expectedCode == code then
    redis.call("set", cntKey, -1)
    return 0
else
    -- 输入错误
    redis.call("decr", cntKey)
    return -1
end