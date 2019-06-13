package model

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jmoiron/sqlx"
	"time"
	//"errors"
	"fmt"
)

const (
	StatusNormal  = 0
	StatusDisable = 1
	StatusExpire  = 2
)

type Activity struct {
	Id           int            `db:"id" json:"id" form:"id"`
	Name         string         `db:"name" json:"name" form:"name"`
	Status       int            `db:"status" json:"status" form:"form"`
	StartTime    mysql.NullTime `db:"startTime" json:"startTime" form:"startTime"`
	StartDate    time.Time      `json:"startDate" form:"startDate"`
	EndTime      mysql.NullTime `db:"endTime" json:"endTime" form:"endTime"`
	EndDate      time.Time      `json:"endDate" form:"endDate"`
	Total        int            `db:"total" json:"total" form:"total"`
	SecSpeed     int            `db:"secSpeed" json:"secSpeed" form:"secSpeed"`
	BuyLimit     int            `db:"buyLimit" json:"buyLimit" form:"buyLimit"`
	BuyRate      int            `db:"buyRate" json:"buyRate" form:"buyRate"`
	StartTimeStr string         `json:"startTimeStr" form:"startTimeStr"`
	EndTimeStr   string         `json:"endTimeStr" form:"endTimeStr"`
	StatusStr    string         `json:"statusStr" form:"statusStr"`
}

type SecProductInfoConf struct {
	ProductId         int
	StartTime         int64
	EndTime           int64
	Status            int
	Total             int
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	//每秒最多能卖多少个
	SoldMaxLimit int
}

type ActivityDao struct {
}

func NewActivityDao() *ActivityDao {
	return &ActivityDao{}
}

func (p *ActivityDao) QueryActivityDetail(activityId int) (activity Activity, err error) {
	sql := "select id, name,status,startTime,endTime, total,secSpeed,buyLimit,buyRate from activity where id=?"
	err = Db.Get(&activity, sql, activityId)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	fillActivity(&activity)
	return
}

func fillActivity(v *Activity) {
	if v == nil {
		return
	}
	v.StartTimeStr = v.StartTime.Time.Format("2006-01-02 15:04:05")
	v.EndTimeStr = v.EndTime.Time.Format("2006-01-02 15:04:05")
	if v.Status == StatusNormal {
		v.StatusStr = "正常"
		if time.Now().Unix() > v.EndTime.Time.Unix() {
			v.StatusStr = "已结束"
		}
	} else if v.Status == StatusDisable {
		v.StatusStr = "已售罄"
	}
}

func (p *ActivityDao) QueryActivityList(param *PageParam) (list []Activity, total int, err error) {
	sql := "select count(*) from activity"
	err = Db.Get(&total, sql)
	if err != nil {
		logs.Warn("QueryActivityList failed, err:%v sql:%v", err, sql)
		return
	}
	if total == 0 {
		return
	}
	sql = "select id, name,status,startTime,endTime, total,secSpeed,buyLimit,buyRate from activity order by id limit ?,?"
	err = Db.Select(&list, sql, param.Start, param.Limit)
	if err != nil {
		logs.Warn("QueryActivityList failed, err:%v sql:%v", err, sql)
		return
	}
	for i := range list {
		fillActivity(&list[i])
	}
	return
}

func (p *ActivityDao) EditActivity(a *Activity) (err error) {
	sql := "update activity set name=?, startTime=?, endTime=?, total=?, secSpeed=?, buyLimit=?, buyRate=? where id=?"
	_, err = Db.Exec(sql, a.Name, a.StartTime.Time, a.EndTime.Time, a.Total, a.SecSpeed, a.BuyLimit, a.BuyRate, a.Id)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	return
}
func (p *ActivityDao) AddActivity(activity *Activity) (err error) {
	sql := "insert into activity(name, startTime, endTime, total, secSpeed, buyLimit, buyRate)values(?,?,?,?,?,?,?)"
	_, err = Db.Exec(sql, activity.Name, activity.StartTime.Time, activity.EndTime.Time, activity.Total, activity.SecSpeed, activity.BuyLimit, activity.BuyRate)
	if err != nil {
		logs.Warn("exec sql failed, err:%v sql:%v", err, sql)
		return
	}

	/*err = p.SyncToEtcd(activity)
	if err != nil {
		logs.Warn("sync to etcd failed, err:%v data:%v", err, activity)
		return
	}*/
	return
}

func (p *ActivityDao) ProductValid(productId int, total int) (valid bool, err error) {
	sql := "select id, name, total, status from product where id=?"
	var productList []*Product
	err = Db.Select(&productList, sql, productId)
	if err != nil {
		logs.Warn("select product failed, err:%v", err)
		return
	}

	if len(productList) == 0 {
		err = fmt.Errorf("product[%v] is not exists", productId)
		return
	}

	if total > productList[0].Total {
		err = fmt.Errorf("product[%v] 的数量非法", productId)
		return
	}

	valid = true
	return
}

/*func (p *ActivityDao) SyncToEtcd(activity *Activity) (err error) {

	if strings.HasSuffix(EtcdPrefix, "/") == false {
		EtcdPrefix = EtcdPrefix + "/"
	}

	etcdKey  := fmt.Sprintf("%s%s", EtcdPrefix, EtcdProductKey)
	secProductInfoList, err := loadProductFromEtcd(etcdKey)


	var secProductInfo SecProductInfoConf
	secProductInfo.EndTime =  activity.EndTime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.Id
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate

	secProductInfoList = append(secProductInfoList, secProductInfo)

	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		logs.Error("json marshal failed, err:%v", err)
		return
	}

	_, err = EtcdClient.Put(context.Background(), etcdKey, string(data))
	if err != nil {
		logs.Error("put to etcd failed, err:%v, data[%v]", err, string(data))
		return
	}

	logs.Debug("put to etcd succ, data:%v", string(data))
	return
}*/

func loadProductFromEtcd(key string) (secProductInfo []SecProductInfoConf, err error) {

	logs.Debug("start get from etcd succ")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := EtcdClient.Get(ctx, key)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", key, err)
		return
	}

	logs.Debug("get from etcd succ, resp:%v", resp)
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed, err:%v", err)
			return
		}

		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	/*
		updateSecProductInfo(conf, secProductInfo)
		logs.Debug("update product info succ, data:%v", secProductInfo)

		initSecProductWatcher(conf)

		logs.Debug("init etcd watcher succ")
	*/
	return
}
