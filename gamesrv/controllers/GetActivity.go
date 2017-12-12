package controllers

import (
	"gamesrv/models/ActivityMgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type GetActivity struct {
	beego.Controller
}

func (o *GetActivity) Post() {
	body, err := json.Marshal(ActivityMgr.Instance().ActivityInfoArray)
	if err != nil {
		beego.Error(err)
		return
	}
	o.Ctx.Output.Body(body)
}
