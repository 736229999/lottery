package controllers

import (
	"encoding/json"
	"gamesrv/models/LotteryManager"

	"github.com/astaxie/beego"
)

type GetLotteryRecord struct {
	beego.Controller
}

func (o *GetLotteryRecord) Post() {
	req := make(map[string]string)

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	if gameTag, ok := req["gameTag"]; ok {
		//beego.Debug("--- Get Record : ", gameTag)

		records, err := LotteryManager.Instance().GetLtryRecord(gameTag)
		if err != nil {
			beego.Error(err)
			return
		}

		if b, err := json.Marshal(records); err == nil {
			o.Ctx.Output.Body(b)
		}

		// if records, ok := LotteryManager.Instance().LotteriesRecordMap[gameTag]; ok {
		// 	//beego.Debug(records)
		// 	if b, err := json.Marshal(records); err == nil {
		// 		o.Ctx.Output.Body(b)
		// 	}
		// } else {
		// 	beego.Error("Not have this lottery : ", gameTag)
		// }
	}
}
