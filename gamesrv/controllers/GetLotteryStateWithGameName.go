package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
)

type GetLotteryStateWithGameName struct {
	beego.Controller
}

func (o *GetLotteryStateWithGameName) Post() {
	req := make(map[string]interface{})
	send := make(map[string]interface{})

	json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if o.Ctx.Input.RequestBody == nil {
		return
	}

	//gameName := req["gameName"].(string)
	//send["lotteryCurrentState"] = lottery.GetLotteryStateWithGameTag(gameName)

	bufres, _ := json.Marshal(send)
	o.Ctx.Output.Body(bufres)
}
