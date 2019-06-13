package base

import (
	"github.com/astaxie/beego"
	"seckill/service"
	"seckill/utils"
)

const (
	loginPage = iota
	homePage
	detailPage
)

type Controller struct {
	beego.Controller
}

func (p *Controller) GetLoginUserId() int {
	return utils.VerifyUserLogin(p.Ctx).UserId

}
func (p *Controller) LoginViewer(template string, data map[interface{}]interface{}) {
	p.parseViewer(template, data, loginPage)
}
func (p *Controller) HomeViewer(template string, data map[interface{}]interface{}) {
	p.parseViewer(template, data, homePage)
}
func (p *Controller) DetailViewer(template string, data map[interface{}]interface{}) {
	p.parseViewer(template, data, detailPage)
}
func (p *Controller) parseViewer(template string, data map[interface{}]interface{}, layoutType int) {
	p.Data["base"] = "/"
	p.Data["loginUser"] = "秒杀"
	if data != nil {
		for k, v := range data {
			p.Data[k] = v
		}
	}
	p.TplName = template
	layout := "layout/"
	switch layoutType {
	case loginPage:
		layout += "empty.html"
	case homePage:
		layout += "frame.html"
	case detailPage:
		layout += "detail.html"
	}
	p.Layout = layout
}

func (p *Controller) Json(data interface{}, err error) {
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"
	result["data"] = data
	if err != nil {
		result["code"] = -1
		result["message"] = err.Error()
		if businessError, ok := err.(*service.BusinessError); ok {
			result["code"] = businessError.Code
			result["message"] = businessError.Message
		}
	}
	p.Data["json"] = result
	p.ServeJSON()
}

func (p *Controller) JsonPage(total int, data interface{}, err error) {
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"
	result["rows"] = data
	result["total"] = total
	if err != nil {
		result["code"] = -1
		result["message"] = err.Error()
	}
	p.Data["json"] = result
	p.ServeJSON()
}
