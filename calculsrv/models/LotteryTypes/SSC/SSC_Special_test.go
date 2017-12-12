package SSC

import (
	"calculsrv/models/LotterySettings"
	"fmt"
	"testing"
)

func Test_PaserBetNum(t *testing.T) {
	fmt.Print("---------------------开始测试-------------------------\n")
	//LotterySettings.Init()
	DbMgr.Instance()

	// var oddsMap map[string]float64 = make(map[string]float64)

	// oddsMap["1"] = 195

	o := &SSC{}

	o.gameTag = "SSC_ChongQing"
	o.currentExpect = 20170820066
	o.openCode = []int{6, 0, 7, 8, 5}
	o.openCodeString = "6,0,7,8,5"
	o.Settings = LotterySettings.LotteriesSettings["SSC_ChongQing"]

	// order := gb.Order{}
	// order.AccountName = "try0056"
	// order.BetNums = "1"
	// order.BettingTime = time.Now().Unix()
	// order.BetType = 57
	// order.Expect = 20170711096
	// order.GameTag = "SSC_ChongQing"
	// order.OrderType = 0
	// order.Rebate = 0

	//accountInfo := AccountMgr.AccountInfo{}

	// if !o.AnalyticalOrder(&order, &accountInfo) {
	// 	beego.Debug("投注失败")
	// }
	orders := DbMgr.Instance().GetLotteryOrderRecord(o.gameTag, o.currentExpect)
	o.SettlementOrders(orders, o.openCodeString)

	// fmt.Print("\n")

	// for _, v := range arrayInt {
	// 	fmt.Print(v, "\n")
	// }

	fmt.Print("----------------------测试结束--------------------------\n")
}
