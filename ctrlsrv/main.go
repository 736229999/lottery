package main

import (
	"ctrlsrv/models/dbmgr"
	"ctrlsrv/models/encmgr"
	"ctrlsrv/models/ltrymgr"
	"ctrlsrv/models/srvmgr"
	_ "ctrlsrv/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Debug("******************************************************************************")
	beego.Debug("--------------------------------------------------------------------------------")
	dbmgr.Instance()// ---连接数据库

	srvmgr.Instance()// ---服务器管理

	encmgr.Instance()// ---加解密

	ltrymgr.Instance()// ---彩票

	beego.Run()
}
