package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"seckill/config"
	_ "seckill/router"
	"seckill/service"
	"seckill/service/soa"
)

func main() {
	err := config.InitAll()
	if err != nil {
		panic(fmt.Sprintf("init failed, err:%v", err))
		return
	}
	service.InitKillService()
	err = soa.InitSoaConfig()
	if err != nil {
		panic(fmt.Sprintf("soa config init return, err:%v", err))
		return
	}
	err = soa.Run()
	if err != nil {
		panic(fmt.Sprintf("service run return, err:%v", err))
		return
	}

	beego.Run()
}
