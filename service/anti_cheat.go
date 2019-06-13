package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"seckill/config"
	"seckill/model"
	"strconv"
)

var conf = &config.SecKillConfig

func idBlack(req *model.SecRequest) (err error) {
	_, ok := conf.IdBlack[strconv.Itoa(req.UserId)]
	if ok {
		err = fmt.Errorf("invalid request")
	}
	return
}

func ipBlack(req *model.SecRequest) (err error) {
	_, ok := conf.IpBlack[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
	}
	return
}

func rateLimit(req *model.SecRequest) (err error) {
	err = fmt.Errorf("invalid request")
	conf.RateLimitMgr.Lock()
	defer conf.RateLimitMgr.Unlock()
	//userId访问频率控制
	limit, ok := conf.RateLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		limit = model.NewSecMinCounter()
		conf.RateLimitMgr.UserLimitMap[req.UserId] = limit
	}
	secIdCount := limit.SecCounter.Count(req.AccessTime.Unix())
	if secIdCount > conf.AccessLimitConf.UserSecAccessLimit {
		logs.Error("AccessLimitConf.UserSecAccessLimit checked, userId[%v]", req.UserId)
		return
	}
	minIdCount := limit.MinCounter.Count(req.AccessTime.Unix())
	if minIdCount > conf.AccessLimitConf.UserMinAccessLimit {
		logs.Error("AccessLimitConf.UserMinAccessLimit checked, userId[%v]", req.UserId)
		return
	}
	//ip频率控制
	limit, ok = conf.RateLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		limit = model.NewSecMinCounter()
		conf.RateLimitMgr.IpLimitMap[req.ClientAddr] = limit
	}
	secIpCount := limit.SecCounter.Count(req.AccessTime.Unix())
	minIpCount := limit.MinCounter.Count(req.AccessTime.Unix())
	if secIpCount > conf.AccessLimitConf.IPSecAccessLimit {
		logs.Error("AccessLimitConf.IPSecAccessLimit checked, userId[%v],ip[%v]", req.UserId, req.ClientAddr)
		return
	}
	if minIpCount > conf.AccessLimitConf.IPMinAccessLimit {
		logs.Error("AccessLimitConf.IPMinAccessLimit checked, userId[%v],ip[%v]", req.UserId, req.ClientAddr)
		return
	}
	err = nil
	return
}

func antiCheat(req *model.SecRequest) (err error) {
	err = idBlack(req)
	if err != nil {
		logs.Error("idBlack checked,userId[%v]", req.UserId)
		return
	}
	err = ipBlack(req)
	if err != nil {
		logs.Error("ipBlack checked,userId[%v]", req.UserId)
		return
	}
	err = rateLimit(req)
	if err != nil {
		logs.Error("rateLimit checked,userId[%v],ip[%v]", req.UserId, req.ClientAddr)
		return
	}
	return
}
