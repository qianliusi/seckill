package config

import (
	"fmt"
	_ "gopkg.in/redis.v4"
	"testing"
	"time"
)

func TestCreateClient(t *testing.T) {
	client, err := initRedisClient(&RedisConf{
		Addr:        "service.qls.com:2379",
		MaxActive:   10,
		IdleTimeout: 10,
	})
	if err != nil {
		t.Error(err)
	}
	for {
		value, err := client.BLPop(3*time.Second, "/qls/backend/seckill/ipblackqueue").Result()
		if err != nil {
			t.Error(err)
		}
		fmt.Println("value: ", value)
	}
}
