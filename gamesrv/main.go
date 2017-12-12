/*
*
*          ┌─┐       ┌─┐
*       ┌──┘ ┴───────┘ ┴──┐
*       │                 │
*       │       ───       │
*       │  ─┬┘       └┬─  │
*       │                 │
*       │       ─┴─       │
*       │                 │
*       └───┐         ┌───┘
*           │         │
*           │         │
*           │         │
*           │         └──────────────┐
*           │                        │
*           │                        ├─┐
*           │                        ┌─┘
*           │                        │
*           └─┐  ┐  ┌───────┬──┐  ┌──┘
*             │ ─┤ ─┤       │ ─┤ ─┤
*             └──┴──┘       └──┴──┘
*                 神兽保佑
*                 代码无BUG!
 */
package main

import (
	"gamesrv/models/DataCenter"
	"gamesrv/models/LotteryManager"
	"gamesrv/models/RechargeMgr"
	"gamesrv/models/ctrl"
	"gamesrv/models/dbmgr"
	"gamesrv/models/encmgr"
	_ "gamesrv/routers"

	"github.com/astaxie/beego"
)

//gamesrv功能缺失备忘
//gamesrv 启动要像 计算服务 和LoginServer 发送注册消息

func main() {

	//1.初始化encmgr
	err := encmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//2.初始化 ctrl
	err = ctrl.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//3.初始化数据库管理员;
	err = dbmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//初始化充值渠道管理
	RechargeMgr.Instance()

	// //初始化 彩票管理员；
	LotteryManager.Instance()

	// //初始化 数据中心
	DataCenter.Instance()

	// //向LoginServer 发送Game服务器启动消息
	//Utils.Registgamesrv()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
