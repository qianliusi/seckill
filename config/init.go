package config

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	etcdClient "go.etcd.io/etcd/clientv3"
	"gopkg.in/redis.v4"
	"seckill/model"
	"time"
)

var (
	Db          *sqlx.DB
	EtcdClient  *etcdClient.Client
	RedisClient *redis.Client
	//RedisPool       *redis.Pool
)

func initDb() (err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", AppConf.mysql.UserName, AppConf.mysql.Passwd,
		AppConf.mysql.Host, AppConf.mysql.Port, AppConf.mysql.Database)
	Db, err = sqlx.Open("mysql", dns)
	if err != nil {
		logs.Error("open mysql failed, err:%v ", err)
		return
	}
	logs.Debug("connect to mysql succ")
	return
}

func initEtcd() (err error) {
	EtcdClient, err = etcdClient.New(etcdClient.Config{
		Endpoints:   []string{AppConf.etcd.Addr},
		DialTimeout: time.Duration(AppConf.etcd.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}
	logs.Debug("init etcd success")
	return
}

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = AppConf.log.Path
	config["level"] = convertLogLevel(AppConf.log.Level)
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}
	err = logs.SetLogger(logs.AdapterFile, string(configStr))
	if err != nil {
		fmt.Println("set logger failed, err:", err)
		return
	}
	return
}

func InitAll() (err error) {
	err = initConfig()
	if err != nil {
		logs.Error("init config failed, err:%v", err)
		return
	}
	err = initDb()
	if err != nil {
		logs.Error("init Db failed, err:%v", err)
		return
	}
	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}
	err = initGoRedis()
	if err != nil {
		logs.Error("init redis failed, err:%v", err)
		return
	}
	model.Init(Db, EtcdClient)
	err = InitSecKill()
	if err != nil {
		logs.Error("init SecKillConf failed, err:%v", err)
		return
	}
	return
}
