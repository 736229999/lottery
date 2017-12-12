package main

import (
	"apisrv/models/apimgr"
	"apisrv/models/ctrl"
	"apisrv/models/dbmgr"
	"apisrv/models/encmgr"
	"apisrv/models/histmgr"
	"apisrv/models/ltrymgr"

	_ "apisrv/routers"

	"github.com/astaxie/beego"
)

func main() {

	//1. 初始化加密管理类(这个必须是第一初始化的,后续的消息都要依赖)
	err := encmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//2. 初始化控制管理类(这个必须是第二初始化,向控制服发送注册信息,获取密钥会放入 encmgr 中)
	err = ctrl.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//3. 初始化数据库类(从 ctrl 服获取这个服务器要依赖的数据库信息,并且连接数据库)
	err = dbmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//4.初始化api管理类(现目前 所有的 API 都是写死在本类对应的配置文件中,等功能做完一定要写入数据库,不然每次修改都要改动源码)
	err = apimgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//5.初始化彩票历史记录类(补全所有的历史信息,由于历史记录有可能有部分彩票无法正确的完成初始化,所以这个初始化就不 return)
	histmgr.Instance()

	//6. 初始化彩票管理类(从控制服获取彩票信息,初始化所有彩票类,开始监听)
	err = ltrymgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	beego.Run()
}
