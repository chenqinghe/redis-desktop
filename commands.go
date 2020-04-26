package main

import "strings"

var Commands = []string{
	// string
	"SET", "SETNX", "SETEX", "PSETEX", "GET", "GETSET", "STRLEN", "APPEND", "SETRANGE", "GETRANGE", "INCR", "INCRBY", "INCRBYFLOAT", "DECR", "DECRBY", "MSET", "MSETNX", "MGET",

	// hash
	"HSET", "HSETNX", "HGET", "HEXISTS", "HDEL", "HLEN", "HSTRLEN", "HINCRBY", "HINCRBYFLOAT", "HMSET", "HMGET", "HKEYS", "HVALS", "HGETALL", "HSCAN",

	// list
	"LPUSH", "LPUSHX", "RPUSH", "RPUSHX", "LPOP", "RPOP", "RPOPLPUSH", "LREM", "LLEN", "LINDEX", "LINSERT", "LSET", "LRANGE", "LTRIM", "BLPOP", "BRPOP", "BRPOPLPUSH",

	// set
	"SADD", "SISMEMBER", "SPOP", "SRANDMEMBER", "SREM", "SMOVE", "SCARD", "SMEMBERS", "SSCAN", "SINTER", "SINTERSTORE", "SUNION", "SUNIONSTORE", "SDIFF", "SDIFFSTORE",

	// sorted set
	"ZADD", "ZSCORE", "ZINCRBY", "ZCARD", "ZCOUNT", "ZRANGE", "ZREVRANGE", "ZRANGEBYSCORE", "ZREVRANGEBYSCORE", "ZRANK", "ZREVRANK", "ZREM", "ZREMRANGEBYRANK", "ZREMRANGEBYSCORE", "ZRANGEBYLEX", "ZLEXCOUNT", "ZREMRANGEBYLEX", "ZSCAN", "ZUNIONSTORE", "ZINTERSTORE",

	// geo
	"GEOADD", "GEOPOS", "GEODIST", "GEORADIUS", "GEORADIUSBYMEMBER", "GEOHASH",

	// bitmap
	"SETBIT", "GETBIT", "BITCOUNT", "BITPOS", "BITOP", "BITFIELD",

	// database
	"EXISTS", "TYPE", "RENAME", "RENAMENX", "MOVE", "DEL", "RANDOMKEY", "DBSIZE", "KEYS", "SCAN", "SORT", "FLUSHDB", "FLUSHALL", "SELECT", "SWAPDB",

	// expire
	"EXPIRE", "EXPIREAT", "TTL", "PERSIST", "PEXPIRE", "PEXPIREAT", "PTTL",

	// multi
	"MULTI", "EXEC", "DISCARD", "WATCH", "UNWATCH",

	// lua script
	"EVAL", "EVALSHA", "SCRIPT_LOAD", "SCRIPT_EXISTS", "SCRIPT_FLUSH", "SCRIPT_KILL",

	// persistent
	"SAVE", "BGSAVE", "BGREWRITEAOF", "LASTSAVE",

	// pub & sub
	"PUBLISH", "SUBSCRIBE", "PSUBSCRIBE", "UNSUBSCRIBE", "PUNSUBSCRIBE", "PUBSUB",

	// cluster
	"SLAVEOF", "ROLE",

	// client
	"AUTH", "QUIT", "INFO", "SHUTDOWN", "TIME", "CLIENT_GETNAME", "CLIENT_KILL", "CLIENT_LIST", "CLIENT_SETNAME",

	// config
	"CONFIG_SET", "CONFIG_GET", "CONFIG_RESETSTAT", "CONFIG_REWRITE",

	// debug
	"PING", "ECHO", "OBJECT", "SLOWLOG", "MONITOR", "DEBUG_OBJECT", "DEBUG_SEGFAULT",

	// internal
	"MIGRATE", "DUMP", "RESTORE", "SYNC", "PSYNC",
}

func parseCommand(cmd string) (string, []interface{}) {
	ss := strings.Split(strings.TrimSpace(cmd), " ")
	tmp := make([]string, 0, len(ss))
	for _, v := range ss {
		if t := strings.TrimSpace(v); len(t) > 0 {
			tmp = append(tmp, t)
		}
	}
	if len(tmp) == 0 {
		return "", nil
	}
	args := make([]interface{}, 0, len(tmp))
	for _, v := range tmp[1:] {
		args = append(args, v)
	}
	return strings.ToUpper(tmp[0]), args
}

var CommandManuals = map[string]string{
	`SET`: `
SET key value [EX seconds] [PX milliseconds] [NX|XX]
可用版本： >= 1.0.0
时间复杂度： O(1)
将字符串值 value 关联到 key 。

如果 key 已经持有其他值， SET 就覆写旧值， 无视类型。

当 SET 命令对一个带有生存时间（TTL）的键进行设置之后， 该键原有的 TTL 将被清除。

可选参数
从 Redis 2.6.12 版本开始， SET 命令的行为可以通过一系列参数来修改：

EX seconds ： 将键的过期时间设置为 seconds 秒。 执行 SET key value EX seconds 的效果等同于执行 SETEX key seconds value 。

PX milliseconds ： 将键的过期时间设置为 milliseconds 毫秒。 执行 SET key value PX milliseconds 的效果等同于执行 PSETEX key milliseconds value 。

NX ： 只在键不存在时， 才对键进行设置操作。 执行 SET key value NX 的效果等同于执行 SETNX key value 。

XX ： 只在键已经存在时， 才对键进行设置操作。

Note

因为 SET 命令可以通过参数来实现 SETNX 、 SETEX 以及 PSETEX 命令的效果， 所以 Redis 将来的版本可能会移除并废弃 SETNX 、 SETEX 和 PSETEX 这三个命令。

返回值
在 Redis 2.6.12 版本以前， SET 命令总是返回 OK 。

从 Redis 2.6.12 版本开始， SET 命令只在设置操作成功完成时才返回 OK ； 如果命令使用了 NX 或者 XX 选项， 但是因为条件没达到而造成设置操作未执行， 那么命令将返回空批量回复（NULL Bulk Reply）。
`,
}
