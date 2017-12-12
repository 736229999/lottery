package main

import (
	"indsrv/models/ctrl"
	"indsrv/models/dbmgr"
	"indsrv/models/encmgr"
	"indsrv/models/ltrymgr"
	"indsrv/models/rd"

	_ "indsrv/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Info("--- Start Independent Lottery Server")

	//1：初始化加密管理类(这个必须是第一初始化的,后续的消息都要依赖)
	err := encmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//2. 初始化控制管理类(这个必须是第二初始化,向控制服发送注册信息,获取密钥会放入 encmgr 中)
	err = ctrl.Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//3.初始化数据库(链接数据库)(dbmgr 不再是单列类)
	err = dbmgr.Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//4.初始化随机数
	err = rd.Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//5.初始化独立彩票管理
	err = ltrymgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	beego.Run()
}
