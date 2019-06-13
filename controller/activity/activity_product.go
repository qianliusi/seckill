package activity

import (
	"seckill/controller/base"
	"seckill/model"
	"seckill/service"
)

var actProDao = model.ActProDao

type ActProController struct {
	base.Controller
}

func (p *ActProController) List() {
	var total int
	var rows interface{}
	var err error
	defer func() {
		p.JsonPage(total, rows, err)
	}()
	param := &model.APListParam{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	param.Parse()
	rows, total, err = actProDao.List(param)
}

func (p *ActProController) Add() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	param := &model.ActivityProduct{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	err = actProDao.Add(param)
	if err != nil {
		return
	}
	err = service.UpdateCurrentKill(param.ActivityId)
}

func (p *ActProController) Edit() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	param := &model.ActivityProduct{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	err = actProDao.Edit(param)
	err = service.UpdateCurrentKill(param.ActivityId)

}
