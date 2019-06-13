package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"seckill/utils"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
)

type APListParam struct {
	PageParam
	ActivityId int `form:"activityId"`
}
type ActivityProduct struct {
	Id           int         `db:"id" json:"id" form:"id"`
	ActivityId   int         `db:"activityId" json:"activityId" form:"activityId"`
	ProductId    int         `db:"productId" json:"productId" form:"productId"`
	Activity     *Activity   `json:"activity"`
	Product      *Product    `json:"product"`
	Status       int         `db:"status" json:"status" form:"status"`
	Total        int         `db:"total" json:"total" form:"total"`
	SecSpeed     int         `db:"secSpeed" json:"secSpeed" form:"secSpeed"`
	BuyLimit     int         `db:"buyLimit" json:"buyLimit" form:"buyLimit"`
	BuyRate      int         `db:"buyRate" json:"buyRate" form:"buyRate"`
	StatusStr    string      `json:"statusStr" form:"statusStr"`
	RateLimitMgr TimeCounter `json:"-"`
}

type ActivityProductDao struct {
}

func NewActivityProductDao() *ActivityProductDao {
	return &ActivityProductDao{}
}

func (p *ActivityProductDao) Edit(ap *ActivityProduct) (err error) {
	sql := "update activity_product set status=?, total=?, secSpeed=?, buyLimit=?,buyRate=? where id=?"
	_, err = Db.Exec(sql, ap.Status, ap.Total, ap.SecSpeed, ap.BuyLimit, ap.BuyRate, ap.Id)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	return
}

func (p *ActivityProductDao) Add(ap *ActivityProduct) (err error) {
	sql := "insert into activity_product(activityId,productId,status,total,secSpeed,buyLimit,buyRate)values(?,?,?,?,?,?,?)"
	_, err = Db.Exec(sql, ap.ActivityId, ap.ProductId, ap.Status, ap.Total, ap.SecSpeed, ap.BuyLimit, ap.BuyRate)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	return
}

func (p *ActivityProductDao) Detail(apId int) (ap ActivityProduct, err error) {
	sql := "select id, activityId,productId,status,total,secSpeed,buyLimit,buyRate from activity_product where id=?"
	err = Db.Get(&ap, sql, apId)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	return
}

func fillAp(list []ActivityProduct) (err error) {
	if len(list) == 0 {
		return
	}
	var productIds = utils.NewSet()
	for _, v := range list {
		productIds.Add(v.ProductId)
	}
	param := &ProductListParam{}
	param.Start = 0
	param.Limit = 100
	param.Ids = productIds.Array()
	lists, _, err := ProDao.QueryProductList(param)
	if err != nil {
		return
	}
	pMap := make(map[int]*Product)
	for i, v := range lists {
		pMap[v.Id] = &lists[i]
	}
	for i, v := range list {
		fillActPro(&list[i])
		if product, ok := pMap[v.ProductId]; ok {
			list[i].Product = product
		}
	}
	return
}

func fillActPro(v *ActivityProduct) {
	if v == nil {
		return
	}
	if v.Status == StatusNormal {
		v.StatusStr = "正常"
	} else if v.Status == StatusDisable {
		v.StatusStr = "已售罄"
	}
}

func (p *ActivityProductDao) List(param *APListParam) (list []ActivityProduct, total int, err error) {
	sql := "select count(*) from activity_product where activityId=?"
	err = Db.Get(&total, sql, param.ActivityId)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	if total == 0 {
		return
	}
	sql = "select id, activityId,productId,status,total,secSpeed,buyLimit,buyRate from activity_product where activityId=? limit ?,?"
	err = Db.Select(&list, sql, param.ActivityId, param.Start, param.Limit)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	err = fillAp(list)
	return
}
