package router

import (
	"github.com/astaxie/beego"
	"seckill/controller/activity"
	"seckill/controller/kill"
	"seckill/controller/login"
	"seckill/controller/product"
	"seckill/controller/sec_conf"
	"seckill/filter"
)

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, filter.UserLoginFilter)
	//登录
	beego.Router("/login", &login.LoginController{}, "*:Login")
	beego.Router("/index", &login.LoginController{}, "*:Index")
	beego.Router("/login.json", &login.LoginController{}, "*:LoginJson")

	//商品管理
	beego.Router("/product/index", &product.ProductController{}, "*:Index")
	beego.Router("/product/list.json", &product.ProductController{}, "*:List")
	beego.Router("/product/add.json", &product.ProductController{}, "*:Add")
	beego.Router("/product/edit.json", &product.ProductController{}, "*:Edit")
	//活动管理
	beego.Router("/activity/index", &activity.ActivityController{}, "*:Index")
	beego.Router("/activity/list.json", &activity.ActivityController{}, "*:List")
	beego.Router("/activity/add.json", &activity.ActivityController{}, "*:Add")
	beego.Router("/activity/edit.json", &activity.ActivityController{}, "*:Edit")
	beego.Router("/activity/detail", &activity.ActivityController{}, "*:Detail")
	beego.Router("/activity/detail.json", &activity.ActivityController{}, "*:DetailJson")
	beego.Router("/activity/add2kill.json", &activity.ActivityController{}, "*:Add2Kill")
	//活动商品管理
	beego.Router("/activity/product/list.json", &activity.ActProController{}, "*:List")
	beego.Router("/activity/product/add.json", &activity.ActProController{}, "*:Add")
	beego.Router("/activity/product/edit.json", &activity.ActProController{}, "*:Edit")

	//秒杀配置
	beego.Router("/sec_conf/index", &sec_conf.SecConfigController{}, "*:Index")
	beego.Router("/sec_conf/detail.json", &sec_conf.SecConfigController{}, "*:DetailJson")
	beego.Router("/sec_conf/user_access_limit/edit.json", &sec_conf.SecConfigController{}, "*:EditUserAccessLimit")
	beego.Router("/sec_conf/refer_white_list/edit.json", &sec_conf.SecConfigController{}, "*:EditReferWhiteList")
	beego.Router("/sec_conf/blacklist/edit.json", &sec_conf.SecConfigController{}, "*:EditBlackList")

	//秒杀
	beego.Router("/seckill/detail", &kill.KillController{}, "*:Detail")
	beego.Router("/seckill/kill.json", &kill.KillController{}, "*:SecKill")
	beego.Router("/seckill/info.json", &kill.KillController{}, "*:SecInfo")
}
