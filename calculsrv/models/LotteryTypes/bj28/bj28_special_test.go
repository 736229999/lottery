package bj28

import (
	"calculsrv/models/dbmgr"
	"calculsrv/models/ltryset"
	"fmt"
	"testing"
)

func Test_PaserBetNum(t *testing.T) {
	fmt.Print("---------------------开始测试-------------------------\n")
	//LotterySettings.Init()
	dbmgr.

	// var oddsMap map[string]float64 = make(map[string]float64)

	// oddsMap["1"] = 195

	o := &BJ28{}

	o.gameTag = "bj28"
	o.currentExpect = 855549
	o.openCode = []int{5, 3, 4}
	o.openCodeString = "5,3,4"
	o.Settings = ltryset.LotteriesSettings["bj28"]

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
	orders := dbmgr.GetLotteryOrderRecord(o.gameTag, o.currentExpect)
	o.SettlementOrders(orders, o.openCodeString)

	// fmt.Print("\n")

	// for _, v := range arrayInt {
	// 	fmt.Print(v, "\n")
	// }

	fmt.Print("----------------------测试结束--------------------------\n")
}
