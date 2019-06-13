package activity

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-sql-driver/mysql"
	"seckill/controller/base"
	"seckill/model"
	"seckill/service"
	"time"
)

var activityDao = model.ActDao

type ActivityController struct {
	base.Controller
}

func (p *ActivityController) Index() {
	p.HomeViewer("activity/index.html", nil)
}

func (p *ActivityController) Detail() {
	var data = make(map[interface{}]interface{})
	data["activityId"] = p.GetString("id")
	p.DetailViewer("activity/detail.html", data)
}

func (p *ActivityController) List() {
	var total int
	var rows interface{}
	var err error
	defer func() {
		p.JsonPage(total, rows, err)
	}()
	param := &model.PageParam{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	param.Parse()
	rows, total, err = activityDao.QueryActivityList(param)
}
func (p *ActivityController) DetailJson() {
	var data interface{}
	var err error
	activityId, err := p.GetInt("id")
	defer func() {
		p.Json(data, err)
	}()
	data, err = activityDao.QueryActivityDetail(activityId)
}

func (p *ActivityController) Add2Kill() {
	var err error
	activityId, err := p.GetInt("id")
	defer func() {
		p.Json(nil, err)
	}()
	err = service.AddActivity2SecKill(activityId)
}

func (p *ActivityController) Add() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	param := &model.Activity{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	if err = checkParam(param); err != nil {
		return
	}
	err = activityDao.AddActivity(param)
}
func checkParam(param *model.Activity) (err error) {
	param.StartTime = mysql.NullTime{Time: param.StartDate}
	param.EndTime = mysql.NullTime{Time: param.EndDate}
	if param.EndTime.Time.Unix() <= param.StartTime.Time.Unix() {
		err = fmt.Errorf("开始时间小于结束时间")
		logs.Error(err)
		return
	}
	now := time.Now().Unix()
	if param.EndTime.Time.Unix() <= now || param.StartTime.Time.Unix() <= now {
		err = fmt.Errorf("开始时间或者结束时间小于现在")
		logs.Error(err)
		return
	}
	if len(param.Name) == 0 {
		err = fmt.Errorf("活动的名字不能为空")
		return
	}
	return
}

func (p *ActivityController) Edit() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	param := &model.Activity{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	if err = checkParam(param); err != nil {
		return
	}
	err = activityDao.EditActivity(param)
	if err != nil {
		return
	}
	err = service.UpdateCurrentKill(param.Id)
}
