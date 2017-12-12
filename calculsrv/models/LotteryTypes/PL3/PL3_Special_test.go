package PL3

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

	o := &PL3{}
	o.gameTag = "PL3"
	o.currentExpect = 2017187
	o.openCode = []int{0, 0, 0}
	o.openCodeString = "0,0,0"
	o.Settings = LotterySettings.LotteriesSettings["PL3"]

	// order := gb.Order{}
	// order.AccountName = "Test001"
	// order.BetNums = "03,04,01,11,10"
	// order.BettingTime = time.Now().Unix()
	// order.BetType = 1
	// order.Expect = 2017060581
	// order.GameTag = "EX5_JiangXi"
	// order.OrderType = 0
	// order.Rebate = 0

	// accountInfo := AccountMgr.AccountInfo{}

	// if !o.AnalyticalOrder(&order, &accountInfo) {
	// 	beego.Debug("投注失败")
	// }

	o.SettlementOrders()

	// fmt.Print("\n")

	// for _, v := range arrayInt {
	// 	fmt.Print(v, "\n")
	// }

	fmt.Print("----------------------测试结束--------------------------\n")
}
