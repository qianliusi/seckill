package soa

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	etcdClient "go.etcd.io/etcd/clientv3"
	"gopkg.in/redis.v4"
	"seckill/config"
	"seckill/model"
	"sync"
	"time"
)

var (
	secSoaCtx *SecSoaContext
)

func InitSoaConfig() (err error) {
	secSoaCtx = &SecSoaContext{}
	secSoaConf := &SecKillSoaConf{}
	secSoaCtx.secSoaConf = secSoaConf
	secSoaConf.TokenPwd = beego.AppConfig.String("seckill_token_passwd")
	secSoaConf.UserRequestQueue = beego.AppConfig.String("user_request_queue_key")
	secSoaConf.UserResponseQueue = beego.AppConfig.String("user_response_queue_key")
	secSoaConf.WriteGoroutineNum = 4
	secSoaConf.ReadGoroutineNum = 4
	secSoaConf.HandleKillGoroutineNum = 4
	secSoaConf.Read2HandleChanSize = 1024
	secSoaConf.Handle2WriteChanSize = 1024
	secSoaConf.MaxRequestWaitTimeout = 10
	secSoaConf.SendToWriteChanTimeout = 100
	secSoaConf.SendToHandleChanTimeout = 100
	secSoaConf.SecProductInfoMap = config.SecKillConfig.SecKillActivity.ProductsMap
	secSoaCtx.Read2HandleChan = make(chan *SecSoaRequest, secSoaCtx.secSoaConf.Read2HandleChanSize)
	secSoaCtx.Handle2WriteChan = make(chan *SecSoaResponse, secSoaCtx.secSoaConf.Handle2WriteChanSize)
	secSoaCtx.HistoryMap = make(map[int]*UserBuyHistory, 10240)
	secSoaCtx.productCountMgr = NewProductCountMgr()
	secSoaCtx.proxy2SoaRedis = config.RedisClient
	secSoaCtx.soa2ProxyRedis = config.RedisClient
	secSoaCtx.etcdClient = config.EtcdClient
	logs.Debug("init soa config success")
	return
}

type SecKillSoaConf struct {
	UserRequestQueue        string
	UserResponseQueue       string
	WriteGoroutineNum       int
	ReadGoroutineNum        int
	HandleKillGoroutineNum  int
	Read2HandleChanSize     int
	Handle2WriteChanSize    int
	MaxRequestWaitTimeout   int
	SendToWriteChanTimeout  int
	SendToHandleChanTimeout int
	SecProductInfoMap       map[int]*model.ActivityProduct
	TokenPwd                string
}

type SecSoaContext struct {
	proxy2SoaRedis   *redis.Client
	soa2ProxyRedis   *redis.Client
	etcdClient       *etcdClient.Client
	RWSecProductLock sync.RWMutex
	secSoaConf       *SecKillSoaConf
	waitGroup        sync.WaitGroup
	Read2HandleChan  chan *SecSoaRequest
	Handle2WriteChan chan *SecSoaResponse
	HistoryMap       map[int]*UserBuyHistory
	HistoryMapLock   sync.Mutex
	//商品计数
	productCountMgr *ProductCountMgr
}

type SecSoaRequest struct {
	ProductId     int       `json:"productId" form:"productId"`
	UserId        int       `json:"userId"`
	AccessTime    time.Time `json:"accessTime"`
	ClientAddr    string    `json:"clientAddr"`
	ClientReferer string    `json:"clientReferer"`
}

type SecSoaResponse struct {
	ProductId int    `json:"productId"`
	UserId    int    `json:"userId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token"`
	TokenTime int64  `json:"tokenTime"`
}
