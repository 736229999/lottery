package controllers

import (
	"encoding/json"
	"gamesrv/models/LotteryManager"

	"github.com/astaxie/beego"
)

type UpdateLotteryRecord struct {
	beego.Controller
}

func (o *UpdateLotteryRecord) Post() {
	var data []LotteryManager.LotteryRecord

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &data)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	//如果数据长度大于100 说明有人伪造消息
	if len(data) > 100 {
		return
	}

	//交换数据顺序，有 新---老 变成 老----新
	var newData []LotteryManager.LotteryRecord
	for i := len(data) - 1; i >= 0; i-- {
		newData = append(newData, data[i])
	}

	LotteryManager.Instance().UpdateLotteryRecord(newData)

	//发送ok
	send := make(map[string]interface{})
	send["Status"] = "OK"
	bufres, _ := json.Marshal(send)
	o.Ctx.Output.Body(bufres)
}
