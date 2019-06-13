package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"seckill/config"
	"seckill/model"
	"strings"
	"time"
)

var secConf = &config.SecKillConfig

func SecKill(req *model.SecRequest) (result *model.SecResult, err error) {
	err = userCheck(req)
	if err != nil {
		err = NewBusinessError(ErrInvalidRequest)
		logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return
	}
	err = antiCheat(req)
	if err != nil {
		err = NewBusinessError(ErrServiceBusy)
		logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return
	}
	err = CheckActivitySecProduct(req.ProductId)
	if err != nil {
		logs.Warn("userId[%d] ActivitySecProduct failed, req[%v]", req.UserId, req)
		return
	}
	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)
	secConf.UserConnMap[userKey] = req.ResultChan
	secConf.SecReqChan <- req
	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
		secConf.UserConnMapLock.Lock()
		delete(secConf.UserConnMap, userKey)
		secConf.UserConnMapLock.Unlock()
	}()
	select {
	case <-ticker.C:
		err = NewBusinessError(ErrProcessTimeout)
		return
	case <-req.CloseNotify:
		err = NewBusinessError(ErrClientClosed)
		return
	case result = <-req.ResultChan:
		if result.Code != 0 {
			err = NewBusinessError(result.Code)
		}
		return
	}
}

func userCheck(req *model.SecRequest) (err error) {
	found := false
	for _, refer := range strings.Split(secConf.ReferWhiteList, ",") {
		if refer == req.ClientReferer {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("invalid request")
		logs.Warn("user[%d] is reject by refer, req[%v]", req.UserId, req)
		return
	}
	return
}

func InitKillService() {
	secConf.RateLimitMgr = &model.RateLimitMgr{
		UserLimitMap: make(map[int]*model.SecMinCounter, 10000),
		IpLimitMap:   make(map[string]*model.SecMinCounter, 10000),
	}
	secConf.SecReqChan = make(chan *model.SecRequest, 10000)
	secConf.UserConnMap = make(map[string]chan *model.SecResult, 10000)
	initRedisProcessFunc()
}

func initRedisProcessFunc() {
	for i := 0; i < config.AppConf.SecKillConf.UserRequestWriteGoRoutineNum; i++ {
		go WriteHandle()
	}

	for i := 0; i < config.AppConf.SecKillConf.UserResponseReadGoRoutineNum; i++ {
		go ReadHandle()
	}
}
