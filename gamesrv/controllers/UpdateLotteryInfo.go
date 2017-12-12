package controllers

import (
	"encoding/json"
	"gamesrv/models/LotteryManager"
	"gamesrv/models/Utils"

	"github.com/astaxie/beego"
)

type UpdateLotteryInfo struct {
	beego.Controller
}

func (o *UpdateLotteryInfo) Post() {
	var data LotteryManager.LotteryInfo
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &data)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	if Utils.VerifyGameTag(data.GameTag) == false {
		return
	}
	//更新这个彩种的信息(这个操作 没有放入下面的判断中是应为 还未知是否有 更新彩票信息在初始化时比更新历史记录消息来的快的情况)
	//看来消息加密时必须的了

	LotteryManager.Instance().SetLtryInfo(data.GameTag, data)

	//LotteryManager.Instance().LotteriesInfoMap[data.GameTag] = data

	records, err := LotteryManager.Instance().GetLtryRecord(data.GameTag)
	if err != nil {
		beego.Error(err)
		return
	}

	if len(records) > 0 {
		//去头部老记录，增加的新记录在头部
		records = records[1:]
		//头部再添加新的元素
		lr := LotteryManager.LotteryRecord{}
		lr.Expect = data.CurrentExpect
		lr.GameName = data.GameTag
		lr.OpenCode = data.OpenCodeString
		lr.OpenTime = data.CurrentOpenTime
		records = append(records, lr)

		//LotteryManager.Instance().LotteriesRecordMap[data.GameTag] = records
		LotteryManager.Instance().SetLtryRecord(data.GameTag, records)

		send := make(map[string]interface{})
		send["Status"] = "OK"
		bufres, _ := json.Marshal(send)
		o.Ctx.Output.Body(bufres)
	}

	//更新这个彩种的历史记录
	// if records, ok := LotteryManager.Instance().LotteriesRecordMap[data.GameTag]; ok {
	// 	if len(records) > 0 {
	// 		//去头部老记录，增加的新记录在头部
	// 		records = records[1:]
	// 		//头部再添加新的元素
	// 		lr := LotteryManager.LotteryRecord{}
	// 		lr.Expect = data.CurrentExpect
	// 		lr.GameName = data.GameTag
	// 		lr.OpenCode = data.OpenCodeString
	// 		lr.OpenTime = data.CurrentOpenTime
	// 		records = append(records, lr)

	// 		LotteryManager.Instance().LotteriesRecordMap[data.GameTag] = records

	// 		send := make(map[string]interface{})
	// 		send["Status"] = "OK"
	// 		bufres, _ := json.Marshal(send)
	// 		o.Ctx.Output.Body(bufres)
	// 	}
	// }
}
