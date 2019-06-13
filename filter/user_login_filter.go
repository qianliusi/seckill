package filter

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"seckill/config"
	"seckill/utils"
)

var UserLoginFilter = func(ctx *context.Context) {
	if ctx.Request.RequestURI != "/login" && ctx.Request.RequestURI != "/login.json" {
		loginTicket := utils.VerifyUserLogin(ctx)
		if loginTicket == nil || loginTicket.UserId == 0 {
			// 验证登录失败
			logs.Info("verify login failed,ticket is empty")
			ctx.Redirect(302, config.AppConf.AppPath+"/login")
		}
	}
}
