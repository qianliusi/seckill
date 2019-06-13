package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gopkg.in/redis.v4"
	"seckill/config"
	"seckill/model"
)

func WriteHandle() {
	for {
		req := <-conf.SecReqChan
		data, err := json.Marshal(req)
		if err != nil {
			logs.Error("json.Marshal failed, error:%v req:%v", err, req)
			continue
		}
		_, err = config.RedisClient.RPush(config.AppConf.SecKillConf.UserRequestQueue, string(data)).Result()
		if err != nil {
			logs.Error("rpush failed, err:%v, req:%v", err, req)
			continue
		}
	}
}

func ReadHandle() {
	for {
		data, err := config.RedisClient.BLPop(0, config.AppConf.SecKillConf.UserResponseQueue).Result()
		if err != nil {
			if err != redis.Nil {
				logs.Error("redis error,%v", err)
			}
			continue
		}
		logs.Debug("rpop from redis suc, data:%v", data)
		result := &model.SecResult{}
		err = json.Unmarshal([]byte(data[1]), &result)
		if err != nil {
			logs.Error("json.Unmarshal failed, err:%v", err)
			continue
		}
		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		conf.UserConnMapLock.Lock()
		resultChan, ok := conf.UserConnMap[userKey]
		conf.UserConnMapLock.Unlock()
		if !ok {
			logs.Warn("user not found:%v", userKey)
			continue
		}
		resultChan <- result
	}
}
