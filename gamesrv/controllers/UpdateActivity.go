package controllers

import (
	"gamesrv/models/ActivityMgr"

	"github.com/astaxie/beego"
)

type UpdateActivity struct {
	beego.Controller
}

func (o *UpdateActivity) Post() {
	ActivityMgr.Instance().UpdateActivityInfo()
}
