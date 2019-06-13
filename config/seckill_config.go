package config

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"gopkg.in/redis.v4"
	"seckill/model"
	"sync"
	"time"
)

var (
	SecKillConfig  SecKillConf
	ErrCannotBeNil = errors.New("can not be nil")
)

type SecKillKeyConf struct {
	AccessLimitKey               string
	ReferWhiteListKey            string
	ActivityProductKey           string
	UserIdBlackKey               string
	UserIdBlackQueueKey          string
	UserIpBlackKey               string
	UserIpBlackQueueKey          string
	UserRequestQueue             string
	UserResponseQueue            string
	UserRequestWriteGoRoutineNum int
	UserResponseReadGoRoutineNum int
}
type AccessLimitConf struct {
	IPSecAccessLimit   int
	UserSecAccessLimit int
	IPMinAccessLimit   int
	UserMinAccessLimit int
}
type SecKillConf struct {
	sync.RWMutex    `json:"-"`
	AccessLimitConf *AccessLimitConf
	ReferWhiteList  string
	IpBlack         map[string]string
	IdBlack         map[string]string
	SecKillActivity *SecKillActProConf
	UserConnMap     map[string]chan *model.SecResult `json:"-"`
	SecReqChan      chan *model.SecRequest           `json:"-"`
	UserConnMapLock sync.Mutex                       `json:"-"`
	RateLimitMgr    *model.RateLimitMgr              `json:"-"`
}

type SecKillActProConf struct {
	Activity    *model.Activity
	Products    []model.ActivityProduct
	ProductsMap map[int]*model.ActivityProduct `json:"-"`
}

func setSecKillActivity(conf *SecKillActProConf) {
	defer SecKillConfig.Unlock()
	SecKillConfig.Lock()
	if conf != nil {
		pMap := make(map[int]*model.ActivityProduct, len(conf.Products))
		for i, v := range conf.Products {
			pMap[v.ProductId] = &conf.Products[i]
		}
		SecKillConfig.SecKillActivity = conf
		SecKillConfig.SecKillActivity.ProductsMap = pMap
	}
}

func initSecKillKeyConf() (err error) {
	initErr := errors.New("init error")
	writeNum, err := beego.AppConfig.Int("user_request_write_goroutine")
	AppConf.SecKillConf.UserRequestWriteGoRoutineNum = writeNum
	if writeNum == 0 || err != nil {
		logs.Error("user_request_write_goroutine is empty")
		return initErr
	}
	readNum, err := beego.AppConfig.Int("user_response_read_goroutine")
	AppConf.SecKillConf.UserResponseReadGoRoutineNum = readNum
	if readNum == 0 || err != nil {
		logs.Error("user_response_read_goroutine is empty")
		return initErr
	}
	reqQueue := beego.AppConfig.String("user_request_queue_key")
	AppConf.SecKillConf.UserRequestQueue = reqQueue
	if len(reqQueue) == 0 {
		logs.Error("user_request_queue_key is empty")
		return initErr
	}
	resQueue := beego.AppConfig.String("user_response_queue_key")
	AppConf.SecKillConf.UserResponseQueue = resQueue
	if len(resQueue) == 0 {
		logs.Error("user_response_queue_key is empty")
		return initErr
	}
	prefix := beego.AppConfig.String("etcd_sec_key_prefix")
	if len(prefix) == 0 {
		logs.Error("etcd_sec_key_prefix is empty")
		return initErr
	}
	product := beego.AppConfig.String("etcd_product_key")
	if len(product) == 0 {
		logs.Error("etcd_product_key is empty")
		return initErr
	}
	AppConf.SecKillConf.ActivityProductKey = prefix + product

	idBlack := beego.AppConfig.String("userId_black_key")
	AppConf.SecKillConf.UserIdBlackKey = prefix + idBlack
	if len(idBlack) == 0 {
		logs.Error("userId_black_key is empty")
		return initErr
	}
	idBlackQueue := beego.AppConfig.String("userId_black_queue_key")
	AppConf.SecKillConf.UserIdBlackQueueKey = prefix + idBlackQueue
	if len(idBlackQueue) == 0 {
		logs.Error("userId_black_queue_key is empty")
		return initErr
	}
	ipBlack := beego.AppConfig.String("userIp_black_key")
	AppConf.SecKillConf.UserIpBlackKey = prefix + ipBlack
	if len(ipBlack) == 0 {
		logs.Error("userIp_black_key is empty")
		return initErr
	}

	ipBlackQueue := beego.AppConfig.String("userIp_black_queue_key")
	AppConf.SecKillConf.UserIpBlackQueueKey = prefix + ipBlackQueue
	if len(ipBlackQueue) == 0 {
		logs.Error("userIp_black_queue_key is empty")
		return initErr
	}
	accLimit := beego.AppConfig.String("access_limit_key")
	AppConf.SecKillConf.AccessLimitKey = prefix + accLimit
	if len(accLimit) == 0 {
		logs.Error("access_limit_key is empty")
		return initErr
	}
	referWhiteList := beego.AppConfig.String("refer_whitelist_key")
	AppConf.SecKillConf.ReferWhiteListKey = prefix + referWhiteList
	if len(referWhiteList) == 0 {
		logs.Error("refer_whitelist_key is empty")
		return initErr
	}
	return
}

func SycConfToEtcd(key string, value interface{}) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		logs.Error("json failed, ", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = EtcdClient.Put(ctx, key, string(data))
	cancel()
	if err != nil {
		logs.Error("put failed, err:", err)
	}
	return
}

func loadSecConf(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := EtcdClient.Get(ctx, key)
	cancel()
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", AppConf.SecKillConf, err)
		return
	}
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] value[%v]", k, v)
		err = updateSecKillConf(key, v.Value)
	}
	return
}
func updateSecKillBlackConf(key string, data []string) {
	logs.Debug("updateSecKillConf [%v]=[%v]", key, data)
	switch key {
	case AppConf.SecKillConf.UserIpBlackQueueKey:
		SecKillConfig.Lock()
		for _, v := range data {
			SecKillConfig.IpBlack[v] = ""
		}
		SecKillConfig.Unlock()
	case AppConf.SecKillConf.UserIdBlackQueueKey:
		SecKillConfig.Lock()
		for _, v := range data {
			SecKillConfig.IdBlack[v] = ""
		}
		SecKillConfig.Unlock()
	}
	return
}
func updateSecKillConf(key string, data []byte) (err error) {
	switch key {
	case AppConf.SecKillConf.ActivityProductKey:
		conf := &SecKillActProConf{}
		err = json.Unmarshal(data, conf)
		if err != nil {
			logs.Error("Unmarshal ActivityProduct failed, err:%v", err)
			return
		}
		setSecKillActivity(conf)
		logs.Debug("updateSecKillConf [%v]=[%v]", key, conf)
	case AppConf.SecKillConf.AccessLimitKey:
		conf := &AccessLimitConf{}
		err = json.Unmarshal(data, conf)
		if err != nil {
			logs.Error("Unmarshal AccessLimitConf failed, err:%v", err)
			return
		}
		SecKillConfig.Lock()
		SecKillConfig.AccessLimitConf = conf
		SecKillConfig.Unlock()
		logs.Debug("updateSecKillConf [%v]=[%v]", key, conf)
	case AppConf.SecKillConf.ReferWhiteListKey:
		var conf string
		err = json.Unmarshal(data, &conf)
		if err != nil {
			logs.Error("Unmarshal ReferWhiteList failed, err:%v", err)
			return
		}
		SecKillConfig.Lock()
		SecKillConfig.ReferWhiteList = conf
		SecKillConfig.Unlock()
		logs.Debug("updateSecKillConf [%v]=[%v]", key, conf)
	}
	return
}
func InitSecKill() (err error) {
	err = loadSecConf(AppConf.SecKillConf.ActivityProductKey)
	if err != nil {
		return
	}
	err = loadSecConf(AppConf.SecKillConf.ReferWhiteListKey)
	if err != nil {
		return
	}
	err = loadSecConf(AppConf.SecKillConf.AccessLimitKey)
	if err != nil {
		return
	}
	initSecConfWatcher()
	logs.Debug("init SecKillConf success, err:%v", err)
	err = loadBlackList()
	return
}

func initSecConfWatcher() {
	go watchSecKillConf(AppConf.SecKillConf.ActivityProductKey)
	go watchSecKillConf(AppConf.SecKillConf.ReferWhiteListKey)
	go watchSecKillConf(AppConf.SecKillConf.AccessLimitKey)
}
func watchSecKillConf(key string) {
	for {
		watch := EtcdClient.Watch(context.Background(), key)
		var conf interface{}
		for res := range watch {
			for _, ev := range res.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &conf)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						continue
					}
					err = updateSecKillConf(key, ev.Kv.Value)
				}
				logs.Debug("get config from etcd: %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
func loadBlackList() (err error) {
	SecKillConfig.IdBlack, err = RedisClient.HGetAll(AppConf.SecKillConf.UserIdBlackKey).Result()
	if err != nil {
		logs.Warn("hget all idBlack failed, err:%v", err)
		return
	}
	SecKillConfig.IpBlack, err = RedisClient.HGetAll(AppConf.SecKillConf.UserIpBlackKey).Result()
	if err != nil {
		logs.Warn("hget all ipBlack failed, err:%v", err)
		return
	}
	go syncBlack(AppConf.SecKillConf.UserIpBlackQueueKey)
	go syncBlack(AppConf.SecKillConf.UserIdBlackQueueKey)
	logs.Debug("loadBlackList success")
	return
}

/*func loadBlackList() (err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	reply, err := conn.Do("hgetall", AppConf.SecKillConf.UserIdBlackKey)
	SecKillConfig.IdBlack, err = redis.IntMap(reply, err)
	if err != nil {
		logs.Warn("hget all idBlack failed, err:%v", err)
		return
	}
	reply, err = conn.Do("hgetall", AppConf.SecKillConf.UserIpBlackKey)
	SecKillConfig.IpBlack, err = redis.IntMap(reply, err)
	if err != nil {
		logs.Warn("hget all ipBlack failed, err:%v", err)
		return
	}
	go syncBlack(AppConf.SecKillConf.UserIpBlackQueueKey)
	go syncBlack(AppConf.SecKillConf.UserIdBlackQueueKey)
	return
}*/

func RedisHset(key, field, value string) (err error) {
	if key == "" || field == "" {
		err = ErrCannotBeNil
		return
	}
	_, err = RedisClient.HSet(key, field, value).Result()
	if err != nil {
		logs.Error("hset to redis failed, err:%v", err)
	}
	logs.Info("RedisHset [%v]:[%v]:[%v]", key, field, value)
	return
}

func RedisPush(key string, data interface{}) (err error) {
	value, err := json.Marshal(data)
	if err != nil {
		logs.Error("marshal failed, err:%v", err)
		return
	}
	s := string(value)
	if len(s) == 0 {
		err = ErrCannotBeNil
		return
	}
	_, err = RedisClient.RPush(key, s).Result()
	if err != nil {
		logs.Error("rpush to redis failed, err:%v", err)
	}
	logs.Info("RedisRPush[%v]:[%v]", key, s)
	return
}

func syncBlack(key string) {
	logs.Debug("syncBlack type:[%s] start", key)
	var list []string
	lastTime := time.Now().Unix()
	for {
		result, err := RedisClient.BLPop(0, key).Result()
		if err != nil {
			if err != redis.Nil {
				logs.Error("syncBlack error,%v", err)
			}
			continue
		}
		curTime := time.Now().Unix()
		list = append(list, result[1])
		if len(list) > 10 || curTime-lastTime > 3 {
			updateSecKillBlackConf(key, list)
			lastTime = curTime
		}
	}
}

/*func syncBlack(key string) {
	logs.Debug("syncBlack type:[%s] start",key)
	var list []string
	conn := RedisPool.Get()
	defer conn.Close()
	for {
		reply, err := conn.Do("BLPOP", key, 10*time.Second)
		data, err := redis.String(reply, err)
		if err != nil {
			continue
		}
		list = append(list, data)
		updateSecKillBlackConf(key, list)
	}
}*/
