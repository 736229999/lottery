package controllers

import (
	"gamesrv/models/LotterySettings"

	"encoding/json"

	"github.com/astaxie/beego"
)

type UpdateLotterySettings struct {
	beego.Controller
}

func (o *UpdateLotterySettings) Post() {
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &LotterySettings.LotteriesSettings)
	if err != nil {
		beego.Debug(err)
		return
	}
}
