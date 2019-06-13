package service

import (
	"seckill/config"
	"seckill/model"
	"time"
)

func AddActivity2SecKill(activityId int) (err error) {
	activity, err := model.ActDao.QueryActivityDetail(activityId)
	if err != nil {
		return
	}
	param := &model.APListParam{}
	param.Start = 0
	param.Limit = 1000
	param.ActivityId = activityId
	list, _, err := model.ActProDao.List(param)
	if err != nil {
		return
	}
	conf := &config.SecKillActProConf{}
	conf.Activity = &activity
	conf.Products = list
	err = config.SycConfToEtcd(config.AppConf.SecKillConf.ActivityProductKey, conf)
	return
}

func UpdateCurrentKill(activityId int) (err error) {
	currentKill := config.SecKillConfig.SecKillActivity
	if currentKill != nil && currentKill.Activity.Id == activityId {
		err = AddActivity2SecKill(activityId)
	}
	return
}

func CheckActivitySecProduct(productId int) (err error) {
	conf := &config.SecKillConfig
	conf.RLock()
	defer conf.RUnlock()
	p, ok := conf.SecKillActivity.ProductsMap[productId]
	if !ok {
		err = NewBusinessError(ErrNotFoundProductId)
		return
	}
	activity := conf.SecKillActivity.Activity
	now := time.Now().Unix()
	if now-activity.StartTime.Time.Unix() < 0 {
		err = NewBusinessError(ErrActivityNotStart)
		return
	}
	if now-activity.EndTime.Time.Unix() > 0 {
		err = NewBusinessError(ErrActivityAlreadyEnd)
		return
	}
	if p.Status == model.ProductStatusForceSaleOut || p.Status == model.ProductStatusSaleOut {
		err = NewBusinessError(ErrProductSaleOut)
		return
	}
	return
}
