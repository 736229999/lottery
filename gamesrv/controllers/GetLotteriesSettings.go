package controllers

import (
	"gamesrv/models/LotterySettings"

	"encoding/json"

	"github.com/astaxie/beego"
)

type GetLotteriesSettings struct {
	beego.Controller
}

func (o *GetLotteriesSettings) Post() {

	body, err := json.Marshal(LotterySettings.LotteriesSettings)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}
