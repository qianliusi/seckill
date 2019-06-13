package login

import (
	"errors"
	"seckill/controller/base"
	"seckill/utils"
	"strings"
)

type LoginController struct {
	base.Controller
}

var loginUsers = map[string]string{
	"admin": "admin123",
	"qls":   "qls123",
}
var loginUserIds = map[string]int{
	"admin": 1,
	"qls":   2,
}

func (p *LoginController) Login() {
	p.LoginViewer("login.html", nil)
}
func (p *LoginController) Index() {
	p.HomeViewer("index.html", nil)
}

func (p *LoginController) Logout() {
	p.Ctx.SetCookie(utils.LoginTicketKey, "", 0, "/", strings.Split(p.Ctx.Request.Host, ":")[0])
	p.LoginViewer("login.html", nil)
}


func (p *LoginController) LoginJson() {
	userName := p.GetString("username")
	password := p.GetString("password")
	var data interface{}
	var errLogin error
	defer func() {
		p.Json(data, errLogin)
	}()
	v, ok := loginUsers[userName]
	if !ok || v != password {
		errLogin = errors.New("用户名或密码不正确")
		return
	}
	ticket, errLogin := utils.GetUserLoginTicket(loginUserIds[userName], userName)
	p.Ctx.SetCookie(utils.LoginTicketKey, ticket, 60*60*12, "/", strings.Split(p.Ctx.Request.Host, ":")[0])
	data = &utils.LoginTicket{RealName: userName}
}
