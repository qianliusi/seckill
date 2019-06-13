package sec_conf

import (
	"errors"
	"seckill/config"
	"seckill/controller/base"
	"seckill/model"
)

type SecConfigController struct {
	base.Controller
}

func (p *SecConfigController) Index() {
	p.DetailViewer("sec_config/detail.html", nil)
}

func (p *SecConfigController) DetailJson() {
	p.Json(config.SecKillConfig, nil)
}

func (p *SecConfigController) EditUserAccessLimit() {
	var err error
	defer func() {
		p.Json(nil, err)
	}()
	param := &config.AccessLimitConf{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	err = config.SycConfToEtcd(config.AppConf.SecKillConf.AccessLimitKey, param)
}

func (p *SecConfigController) EditReferWhiteList() {
	var err error
	defer func() {
		p.Json(nil, err)
	}()
	param := &config.SecKillConf{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	if param.ReferWhiteList == "" {
		err = errors.New("不能为空")
		return
	}
	err = config.SycConfToEtcd(config.AppConf.SecKillConf.ReferWhiteListKey, param.ReferWhiteList)
}

func (p *SecConfigController) EditBlackList() {
	var err error
	defer func() {
		p.Json(nil, err)
	}()
	blackType := p.GetString("type")
	if blackType == "" {
		err = errors.New("type不能为空")
		return
	}
	param := &model.DataParam{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	if param.Data == "" {
		err = errors.New("不能为空")
		return
	}
	keyQueue := config.AppConf.SecKillConf.UserIpBlackQueueKey
	key := config.AppConf.SecKillConf.UserIpBlackKey
	if blackType == "id" {
		keyQueue = config.AppConf.SecKillConf.UserIdBlackQueueKey
		key = config.AppConf.SecKillConf.UserIdBlackKey
	}
	err = config.RedisPush(keyQueue, param.Data)
	err = config.RedisHset(key, param.Data, "")
}
