package controllers

import (
	"calculsrv/models/gb"
	"calculsrv/models/ltrymgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type LotteryBetting struct {
	beego.Controller
}

func (o *LotteryBetting) Post() {
	//解析消息
	reqBettingInfo := gb.MsgBettingInfo{}
	respBettingInfo := gb.MsgRespBetting{}

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &reqBettingInfo); err != nil {
		beego.Debug(err)
		return
	}

	//验证消息是否有这个彩票实体 (自动彩种)
	if lottery, ok := ltrymgr.Instance().LtryifMap[reqBettingInfo.GameTag]; ok {
		//解析订单(分拆订单)
		respBettingInfo.Status = lottery.AnalyticalBetting(reqBettingInfo)
	} else {
		beego.Error("失败,没有这个彩票类型 : ", reqBettingInfo.GameTag)
		return
	}
	// else if lottery, ok := ManualLottery.Instance().LotteriesInterfaceMap[reqBettingInfo.GameTag]; ok { //验证消息是否有这个彩票实体 (手动彩种)
	// 	//解析订单(分拆订单)
	// 	respBettingInfo.Status = lottery.AnalyticalBetting(reqBettingInfo)
	// }
	// else {
	// 	beego.Error("失败,没有这个彩票类型 : ", reqBettingInfo.GameTag)
	// 	return
	// }

	body, err := json.Marshal(respBettingInfo)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
	return
}

//下注失败错误码
//0 下注成功
//1 彩票状态不正确
//2 未在可下注时间
//3 期数错误
//4 订单信息错误
//5 订单为0
