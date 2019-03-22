package main

import (
	"NewLottApi/cron"
	_ "NewLottApi/routers"
	"fmt"
	"net"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	var bStatus bool = false

	//配置网卡地址
	sMac := beego.AppConfig.String("mac")

	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces { //获取本机MAC地址
		sMyMacStr := fmt.Sprintf("%s", inter.HardwareAddr)
		if sMyMacStr == sMac {
			bStatus = true
			break
		}
	}

	//如果没有绑定网卡，直接退出
	if !bStatus {
		fmt.Println("Illegal access")
		//		os.Exit(1)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	} else {
		//生产环境下使用grace进行热升级
		beego.BConfig.Listen.Graceful = true
	}

	///////////
	//定时任务//
	///////////
	go cron.StartTasks()

	beego.Run()
}
