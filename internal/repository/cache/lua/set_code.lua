-- 验证码的 key
local key = KEYS[1]
-- 验证次数
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then -- 没有过期时间
    return -2
-- key 不存在，或者超过一分钟了，可以重新设置
elseif ttl == 2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
-- 发送太频繁
    return -1
end
