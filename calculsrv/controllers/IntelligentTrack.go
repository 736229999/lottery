package controllers

import (
	"github.com/astaxie/beego"
)

type IntelligentTrack struct {
	beego.Controller
}

func (o *IntelligentTrack) Post() {
	// //解析消息
	// req := gb.MspIntelligentTrackBetting{}
	// resp := gb.RespIntelligentTrackBetting{}

	// if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
	// 	beego.Debug("------------------------- PlaceOrder 消息解析错误 ! -------------------------")
	// 	return
	// }

	//检查用户名
	// accountInfo := AccountMgr.AccountInfo{}
	// err := accountInfo.Init(req.AccountName)
	// if err != nil {
	// 	return
	// }

	// if !accountInfo.VerifyToken(req.Token) {
	// 	resp.Status = 1 //	token 错误
	// 	return
	// }

	// //总计金额
	// var TotalAmount float64 = 0

	// //开始循环订单
	// for _, v := range req.TrackOrders {
	// 	for _, j := range v.Multiple {
	// 		TotalAmount += v.SingleBetAmount * float64(j)
	// 	}
	// }

	// // //验证消息是否有这个彩票实体
	// // if lottery, ok := LotteryManager.Instance().LotteriesInterfaceMap[reqBettingInfo.GameTag]; ok {
	// // 	//解析订单(分拆订单)
	// // 	respBettingInfo.Status = lottery.AnalyticalBetting(reqBettingInfo)
	// // } else {
	// // 	beego.Debug("失败,没有这个彩票类型")
	// // 	return
	// // }

	// body, err := json.Marshal(resp)
	// if err != nil {
	// 	beego.Debug(err)
	// 	return
	// }
	// o.Ctx.Output.Body(body)
	// return
}
