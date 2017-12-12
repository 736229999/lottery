package controllers

import (
	"gamesrv/models/LotteryManager"
	"encoding/json"

	"github.com/astaxie/beego"
)

type UpdateRecordForLostExpect struct {
	beego.Controller
}

func (o *UpdateRecordForLostExpect) Post() {
	var data LotteryManager.LotteryRecord

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &data)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	LotteryManager.Instance().UpdateLotteryOneRecord(data)

	//发送ok
	send := make(map[string]interface{})
	send["Status"] = "OK"
	bufres, _ := json.Marshal(send)
	o.Ctx.Output.Body(bufres)
}
