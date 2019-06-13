package config

import (
	"github.com/astaxie/beego/logs"
	"gopkg.in/redis.v4"
	"time"
)

func initRedisClient(redisConf *RedisConf) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:        redisConf.Addr,
		IdleTimeout: time.Duration(redisConf.IdleTimeout) * time.Second,
		PoolSize:    redisConf.MaxActive,
		Password:    "",
		DB:          0,
	})
	_, err = client.Ping().Result()
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}
	return
}
func initGoRedis() (err error) {
	RedisClient, err = initRedisClient(&AppConf.redis)
	if err != nil {
		logs.Error("init redis failed, err:%v", err)
	}
	logs.Debug("init redis success")
	return
}
