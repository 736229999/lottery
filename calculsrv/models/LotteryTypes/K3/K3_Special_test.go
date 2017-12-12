package K3

import (
	"calculsrv/models/LotterySettings"
	"fmt"
	"testing"
)

func Test_PaserBetNum(t *testing.T) {
	fmt.Print("---------------------开始测试-------------------------\n")
	//LotterySettings.Init()
	//DbMgr.Instance()

	// var oddsMap map[string]float64 = make(map[string]float64)

	// oddsMap["1"] = 195

	o := &K3{}
	o.gameTag = "K3_GuangXi"
	o.currentExpect = 20170712074
	o.openCode = []int{4, 5, 6}
	o.openCodeString = "04,05,06"
	o.Settings = LotterySettings.LotteriesSettings["K3_GuangXi"]

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

// func Test_PaserBetNum(t *testing.T) {
// 	o := &EX5{}
// 	if o.PaserBetNum("1,2,3,4,5,6,7,8,9,10,11,12") {
// 		t.Log(o.PaserBetNum("成功"))
// 	} else {
// 		t.Log(o.PaserBetNum("错误"))
// 	}

// }
