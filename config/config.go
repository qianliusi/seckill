package config

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	AppConf Config
)

type MysqlConfig struct {
	UserName string
	Passwd   string
	Port     int
	Database string
	Host     string
}
type Config struct {
	AppPath     string
	mysql       MysqlConfig
	etcd        EtcdConf
	log         LogConf
	redis       RedisConf
	SecKillConf SecKillKeyConf
}

type EtcdConf struct {
	Addr    string
	Timeout int
}

type LogConf struct {
	Path  string
	Level string
}

type RedisConf struct {
	Addr        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

func initMysqlConf() (err error) {
	initErr := errors.New("init error")
	username := beego.AppConfig.String("mysql_user_name")
	if len(username) == 0 {
		logs.Error("mysql_user_name is empty")
		return initErr
	}
	AppConf.mysql.UserName = username

	mysqlPwd := beego.AppConfig.String("mysql_passwd")
	if len(mysqlPwd) == 0 {
		logs.Error("mysql_passwd is empty")
		return initErr
	}
	AppConf.mysql.Passwd = mysqlPwd
	mysqlHost := beego.AppConfig.String("mysql_host")
	if len(mysqlHost) == 0 {
		logs.Error("mysql_host is empty")
		return initErr
	}
	AppConf.mysql.Host = mysqlHost
	mysqlDb := beego.AppConfig.String("mysql_database")
	if len(mysqlDb) == 0 {
		logs.Error("mysql_database is empty")
		return initErr
	}
	AppConf.mysql.Database = mysqlDb
	port, err := beego.AppConfig.Int("mysql_port")
	if err != nil {
		logs.Error("load config of mysql_port failed, err:%v", err)
		return
	}
	AppConf.mysql.Port = port
	return
}

func initEtcdConf() (err error) {
	initErr := errors.New("init error")
	addr := beego.AppConfig.String("etcd_addr")
	if len(addr) == 0 {
		logs.Error("etcd_addr is empty")
		return initErr
	}
	AppConf.etcd.Addr = addr
	timeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		logs.Error("load config of etcd_timeout failed, err:%v", err)
		return
	}
	AppConf.etcd.Timeout = timeout
	return
}
func initLogConf() (err error) {
	initErr := errors.New("init error")
	AppConf.AppPath = beego.AppConfig.String("app_path")
	path := beego.AppConfig.String("log_path")
	if len(path) == 0 {
		logs.Error("log_path is empty")
		return initErr
	}
	AppConf.log.Path = path
	level := beego.AppConfig.String("log_level")
	if len(level) == 0 {
		logs.Error("log_level is empty")
		return initErr
	}
	AppConf.log.Level = level
	return
}
func initRedisConf() (err error) {
	initErr := errors.New("init error")
	redisAddr := beego.AppConfig.String("redis_addr")
	if len(redisAddr) == 0 {
		logs.Error("redis_addr is empty")
		return initErr
	}
	AppConf.redis.Addr = redisAddr

	redisMaxIdle, err := beego.AppConfig.Int("redis_idle")
	if err != nil {
		logs.Error("load config of redis_idle failed, error:%v", err)
		return
	}
	AppConf.redis.MaxIdle = redisMaxIdle
	redisMaxActive, err := beego.AppConfig.Int("redis_active")
	if err != nil {
		logs.Error("load config of redis_active failed, error:%v", err)
		return
	}
	AppConf.redis.MaxActive = redisMaxActive
	redisIdleTimeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		logs.Error("load config of redis_active failed, error:%v", err)
		return
	}
	AppConf.redis.IdleTimeout = redisIdleTimeout
	return
}

func initConfig() (err error) {
	err = initLogConf()
	if err != nil {
		return
	}
	err = initLogger()
	if err != nil {
		return
	}
	err = initMysqlConf()
	if err != nil {
		return
	}
	err = initEtcdConf()
	if err != nil {
		return
	}
	err = initRedisConf()
	if err != nil {
		return
	}
	err = initSecKillKeyConf()
	if err != nil {
		return
	}
	logs.Info("init config success:%v", AppConf)
	return
}
