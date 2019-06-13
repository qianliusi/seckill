package kill

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"seckill/config"
	"seckill/controller/base"
	"seckill/model"
	"seckill/service"
	"strings"
	"time"
)

type KillController struct {
	base.Controller
}

func (p *KillController) SecKill() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	secRequest := &model.SecRequest{}
	if err = p.ParseForm(secRequest); err != nil {
		return
	}
	secRequest.ResultChan = make(chan *model.SecResult, 1)
	secRequest.AccessTime = time.Now()
	secRequest.UserId = p.GetLoginUserId()
	if len(p.Ctx.Request.RemoteAddr) > 0 {
		secRequest.ClientAddr = strings.Split(p.Ctx.Request.RemoteAddr, ":")[0]
	}
	secRequest.ClientReferer = p.Ctx.Request.Referer()
	secRequest.CloseNotify = p.Ctx.ResponseWriter.CloseNotify()
	reqBytes, err := json.Marshal(secRequest)
	logs.Debug("SecKill request:[%v]", string(reqBytes))
	data, err = service.SecKill(secRequest)
	return
}

func (p *KillController) SecInfo() {
	p.Json(config.SecKillConfig.SecKillActivity, nil)
}

func (p *KillController) Detail() {
	var data = make(map[interface{}]interface{})
	data["activityId"] = p.GetString("id")
	p.DetailViewer("seckill/detail.html", data)
}
