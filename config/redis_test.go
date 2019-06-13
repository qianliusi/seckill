package config

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestInitRedis(t *testing.T) {
	redisConfig := RedisConf{
		Addr:        "service.qls.com:6379",
		MaxIdle:     8,
		MaxActive:   16,
		IdleTimeout: 300,
	}
	pool, err := initRedisPool(redisConfig)
	if err != nil {
		t.Errorf("init redis pool failed, err:%v", err)
	}
	conn := pool.Get()
	defer conn.Close()
	for {
		reply, err := conn.Do("BLPOP", "/qls/backend/seckill/ipblackqueue", 3)
		data, err := redis.String(reply, err)
		if err != nil {
			t.Errorf("failed, err:%v", err)
		}
		fmt.Printf("aldksfjdlasfj:%v\n", data)
	}
}
