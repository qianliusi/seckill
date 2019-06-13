package soa

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gopkg.in/redis.v4"
	"seckill/model"
	"seckill/service"
	"seckill/utils"
	"time"
)

func Run() (err error) {
	go runProcess()
	return
}
func runProcess() {
	for i := 0; i < secSoaCtx.secSoaConf.ReadGoroutineNum; i++ {
		secSoaCtx.waitGroup.Add(1)
		go HandleReader()
	}
	for i := 0; i < secSoaCtx.secSoaConf.WriteGoroutineNum; i++ {
		secSoaCtx.waitGroup.Add(1)
		go HandleWrite()
	}
	for i := 0; i < secSoaCtx.secSoaConf.HandleKillGoroutineNum; i++ {
		secSoaCtx.waitGroup.Add(1)
		go HandleSecKillReq()
	}
	logs.Debug("all process started")
	secSoaCtx.waitGroup.Wait()
	logs.Debug("wait all goroutine exited")
	return
}

func HandleReader() {
	logs.Debug("read goroutine running")
	for {
		data, err := secSoaCtx.proxy2SoaRedis.BLPop(0, secSoaCtx.secSoaConf.UserRequestQueue).Result()
		if err != nil {
			if err != redis.Nil {
				logs.Error("pop from queue failed,%v", err)
			}
			continue
		}
		req := &SecSoaRequest{}
		err = json.Unmarshal([]byte(data[1]), &req)
		if err != nil {
			logs.Error("json.Unmarshal failed, err:%v", err)
			continue
		}
		logs.Debug("pop from queue, data:%v", req)
		now := time.Now().Unix()
		if now-req.AccessTime.Unix() >= int64(secSoaCtx.secSoaConf.MaxRequestWaitTimeout) {
			logs.Warn("req[%v] is expire", req)
			continue
		}
		timer := time.NewTicker(time.Millisecond * time.Duration(secSoaCtx.secSoaConf.SendToHandleChanTimeout))
		select {
		case secSoaCtx.Read2HandleChan <- req:
		case <-timer.C:
			logs.Warn("send to handle chan timeout, req:%v", req)
			break
		}
	}
}

func HandleWrite() {
	logs.Debug("handle write running")
	for res := range secSoaCtx.Handle2WriteChan {
		err := sendToRedis(res)
		if err != nil {
			logs.Error("send to redis, err:%v, res:%v", err, res)
			continue
		}
	}
}

func sendToRedis(res interface{}) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		logs.Error("marshal failed, err:%v", err)
		return
	}
	_, err = secSoaCtx.soa2ProxyRedis.RPush(secSoaCtx.secSoaConf.UserResponseQueue, string(data)).Result()
	if err != nil {
		logs.Error("rpush failed, err:%v, data:%v", err, res)
	}
	return
}

func HandleSecKillReq() {
	logs.Debug("handle sec kill running")
	for req := range secSoaCtx.Read2HandleChan {
		logs.Debug("begin process request:%v", req)
		res, err := HandleSecKill(req)
		if err != nil {
			logs.Warn("process request %v failed, err:%v", err)
			res = &SecSoaResponse{
				Code: service.ErrServiceBusy,
			}
		}
		timer := time.NewTicker(time.Millisecond * time.Duration(secSoaCtx.secSoaConf.SendToWriteChanTimeout))
		select {
		case secSoaCtx.Handle2WriteChan <- res:
		case <-timer.C:
			logs.Warn("send to response chan timeout, res:%v", res)
			break
		}
	}
	return
}

func HandleSecKill(req *SecSoaRequest) (res *SecSoaResponse, err error) {
	secSoaCtx.RWSecProductLock.RLock()
	defer secSoaCtx.RWSecProductLock.RUnlock()
	res = &SecSoaResponse{}
	res.UserId = req.UserId
	res.ProductId = req.ProductId
	product, ok := secSoaCtx.secSoaConf.SecProductInfoMap[req.ProductId]
	if !ok {
		logs.Error("not found product:%v", req.ProductId)
		res.Code = service.ErrNotFoundProductId
		return
	}

	if product.Status == model.ProductStatusForceSaleOut || product.Status == model.ProductStatusSaleOut {
		res.Code = service.ErrProductSaleOut
		return
	}
	now := time.Now().Unix()
	if product.SecSpeed > 0 {
		if product.RateLimitMgr == nil {
			product.RateLimitMgr = &model.SecCounter{}
		}
		alreadySoldCount := product.RateLimitMgr.Check(now)
		if alreadySoldCount >= product.SecSpeed {
			res.Code = service.ErrRetry
			return
		}
	}
	userHistory, ok := secSoaCtx.HistoryMap[req.UserId]
	if product.BuyLimit > 0 {
		secSoaCtx.HistoryMapLock.Lock()
		if !ok {
			userHistory = &UserBuyHistory{
				history: make(map[int]int, 128),
			}
			secSoaCtx.HistoryMap[req.UserId] = userHistory
		}
		historyCount := userHistory.GetProductBuyCount(req.ProductId)
		secSoaCtx.HistoryMapLock.Unlock()
		if historyCount >= product.BuyLimit {
			res.Code = service.ErrAlreadyBuy
			return
		}
	}
	curSoldCount := secSoaCtx.productCountMgr.Count(req.ProductId)
	logs.Info("curSoldCount:%v", curSoldCount)
	if curSoldCount >= product.Total {
		res.Code = service.ErrSoldout
		product.Status = model.ProductStatusSaleOut
		return
	}
	curRate := utils.RandRangeInt(0, 101)
	logs.Debug("curRate:%v productRate:%v", curRate, product.BuyRate)
	if curRate > product.BuyRate {
		res.Code = service.ErrRetry
		return
	}
	userHistory.Add(req.ProductId, 1)
	secSoaCtx.productCountMgr.Add(req.ProductId, 1)
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId, req.ProductId, now, secSoaCtx.secSoaConf.TokenPwd)
	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = now
	return
}
