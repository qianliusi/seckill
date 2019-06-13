package config

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func initRedisPool(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.MaxIdle,
		MaxActive:   redisConf.MaxActive,
		IdleTimeout: time.Duration(redisConf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConf.Addr)
		},
	}
	conn := pool.Get()
	defer func() {
		err = conn.Close()
	}()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}
	return
}

/*func initRedis() (err error) {
	RedisPool, err = initRedisPool(AppConf.redis)
	if err != nil {
		logs.Error("init redis pool failed, err:%v", err)
	}
	return
}*/
