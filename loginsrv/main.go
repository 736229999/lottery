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
	"loginsrv/models/DbHandle"
	"loginsrv/models/ctrl"
	"loginsrv/models/dbmgr"
	"loginsrv/models/encmgr"
	_ "loginsrv/routers"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
)


func main() {
	beego.Debug("master的分支")
	beego.Debug("这是一个新功能")
	beego.Debug("bug分支")
	beego.Debug("bra的新功能")
	//设定随机数种子
	rand.Seed(time.Now().Unix())

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

	//3.初始化数据库	 ---验证是试玩服务器还是正式服务器，然后连接数据库
	err = dbmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}
	//初始化数据库处理类；
	DbHandle.Init()
	beego.Info("--- Login server initialization done !")

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

