package bj28

import (
	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/Order"
	"calculsrv/models/acmgr"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"
	"common/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//红波 0
var red = [8]int{3, 6, 9, 12, 15, 18, 21, 24}

//蓝波 1
var blue = [8]int{2, 5, 8, 11, 17, 20, 23, 26}

//绿波 2
var green = [8]int{1, 4, 7, 10, 16, 19, 22, 25}

//灰 3
var grey = [4]int{0, 13, 14, 27}

//色彩
var color map[int]int = make(map[int]int)

func (o *BJ28) InitNumProperty() {
	color[0] = 3
	color[1] = 2
	color[2] = 1
	color[3] = 0
	color[4] = 2
	color[5] = 1
	color[6] = 0
	color[7] = 2
	color[8] = 1
	color[9] = 0
	color[10] = 2
	color[11] = 1
	color[12] = 0
	color[13] = 3
	color[14] = 3
	color[15] = 0
	color[16] = 2
	color[17] = 1
	color[18] = 0
	color[19] = 2
	color[20] = 1
	color[21] = 0
	color[22] = 2
	color[23] = 1
	color[24] = 0
	color[25] = 2
	color[26] = 1
	color[27] = 3
}

//分析订单(下注)
func (o *BJ28) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {
	//2分析反水设置是否正确(所有类型都一样，所以放在这里)
	if order.Rebate > accountInfo.Rebate {
		beego.Error("订单反水超过限制")
		return false
	}

	//3分析单注金额有没有超过限制
	if order.SingleBetAmount > o.Settings[order.BetType].SingleLimit {
		beego.Error("单注金额超过限制")
		return false
	}
	//------------------------------------------------------------------------------------------------------------------------------------------
	switch order.BetType {
	case 0: //混合
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_0(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		l := len(array)
		//数组数量不得大于10个元素或小于1个元素
		if l > 10 || l < 1 {
			beego.Debug("失败")
			return false
		}

		//分析订单注数 , 投几个号就是几注
		var singleBetNum = l

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 1: //色波
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_1(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		//数组数量不得大于10个元素或小于1个元素
		if l > 3 || l < 1 {
			beego.Debug("失败")
			return false
		}

		//数组数量不得大于10个元素或小于1个元素
		order.SingleBetNum = l

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 2: //豹子
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_2(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		//数组数量不得大于10个元素或小于1个元素
		if l != 1 {
			beego.Debug("失败")
			return false
		}
		order.SingleBetNum = 1

		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 3: //特码包三
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_3(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		//数组数量不得大于10个元素或小于1个元素
		if l != 3 {
			beego.Debug("失败")
			return false
		}
		order.SingleBetNum = 1

		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 4: //特码
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_3(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		l := len(array)
		//数组数量不得大于10个元素或小于1个元素
		if l > 28 || l < 1 {
			beego.Debug("失败")
			return false
		}

		//分析订单注数 , 投几个号就是几注
		var singleBetNum = l

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	default:
		beego.Debug("失败")
		return false
	}

	//8判断订单总额有没有超过限制
	if order.SingleBetAmount*float64(order.SingleBetNum) > o.Settings[order.BetType].OrderLimit {
		beego.Debug("失败")
		return false
	}

	//9计算订单总额（所有类型都一样）
	order.OrderAmount = order.SingleBetAmount * float64(order.SingleBetNum)
	return true
}

//结算这个彩种当期所有订单
func (o *BJ28) SettlementOrders(orders []gb.Order, openCode string) {
	l := len(orders)
	//没有订单
	if l < 1 {
		return
	}
	//金额流水数组
	BalanceRecourds := []BalanceRecordMgr.BalanceRecord{}
	for i := 0; i < l; i++ {
		if !o.settlementOrder(&orders[i], utils.PaserOpenCodeToArray(openCode)) {
			continue
		}

		//更新订单,上线稳定后,改为批量更新订单,并且钱要用整形,以分为单位
		//结算结果要保留两位小数 4舍5入
		s := fmt.Sprintf("%.2f", orders[i].Settlement)
		orders[i].Settlement, _ = strconv.ParseFloat(s, 64)

		ss := fmt.Sprintf("%.2f", orders[i].RebateAmount)
		orders[i].RebateAmount, _ = strconv.ParseFloat(ss, 64)

		//更新订单,上线稳定后,改为批量更新订单
		orders[i].OpenCode = openCode
		dbmgr.UpdateOrder(&orders[i])

		order := orders[i]
		//获得这条订单的账户信息
		accountInfo := acmgr.AccountInfo{}
		err := accountInfo.Init(order.AccountName)
		if err != nil {
			beego.Emergency(err)
			return
		}

		//生成流水记录
		BalanceRecourd := BalanceRecordMgr.BalanceRecord{}
		BalanceRecourd.Serial_Number = Order.Instance().GetOrderNumber()
		BalanceRecourd.Account_name = accountInfo.Account_Name
		BalanceRecourd.Money_Before = accountInfo.Money
		BalanceRecourd.Money = order.Settlement
		BalanceRecourd.Money_After = accountInfo.Money + order.Settlement //注意 ： 投注只有减钱，而结算只有加钱
		BalanceRecourd.Gap_Money = 0
		BalanceRecourd.Type = 1    //1订单(我这里只有1)
		BalanceRecourd.Subitem = 2 //1投注, 2结算
		BalanceRecourd.Trading_Time = utils.GetNowUTC8Time().Unix()
		BalanceRecourd.Status = 1
		BalanceRecourd.Order_Number = order.OrderNumber
		BalanceRecourds = append(BalanceRecourds, BalanceRecourd)
		//结算只有加钱。。。。
		accountInfo.AddMoney(order.Settlement)
		//更新用户信息
		err = accountInfo.UpdataDb()
		if err != nil {
			beego.Emergency(err)
			return
		}
	}

	//插~!~
	bl := len(BalanceRecourds)
	if bl == 1 {
		dbmgr.InsertBalanceRecord(BalanceRecourds[0])
	} else if bl > 1 {
		dbmgr.BulkInsertBalanceRecord(BalanceRecourds)
	}

	//结算订单完成
}

//结算这个彩种一个订单
func (o *BJ28) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}

	//转下注号码为数组
	switch order.BetType {
	case 0: //0 混合
		o.WinningAndLose_0(order, openCode)
	case 1: //色波
		o.WinningAndLose_1(order, openCode)
	case 2: //豹子
		o.WinningAndLose_2(order, openCode)
	case 3: //特码包三
		o.WinningAndLose_3(order, openCode)
	case 4: //特码
		o.WinningAndLose_4(order, openCode)
	default:
		beego.Debug("失败")
		return false
	}
	return true
}

//混合
func (o *BJ28) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	//判断中几注
	var winningBetNum = 0

	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	sum := openCode[0] + openCode[1] + openCode[2]

	//大小 , 单双, 大单大双小单小双 , 极大极小 一共4个开奖结果
	// 0 大, 1 小, 2 单 , 3 双, 4 大单, 5 大双, 6 小单, 7 小双, 8 极大, 9 极小

	//判断开奖结果能中的号码
	var resultCode [4]int

	//小
	if sum < 14 {
		resultCode[0] = 1
		if sum%2 == 0 { //小单, 小双
			resultCode[2] = 7
		} else {
			resultCode[2] = 6
		}

		if sum < 6 { //极小
			resultCode[3] = 9
		} else {
			resultCode[3] = 10 //如果不是极小 那么赋值10 客户端无法下注10
		}

	} else { //大
		resultCode[0] = 0
		if sum%2 == 0 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}

		if sum > 21 {
			resultCode[3] = 8
		} else {
			resultCode[3] = 10 //如果不是极小 那么赋值10 客户端无法下注10
		}
	}

	//判断单双
	if sum%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
				winningBetNum++
				//计算中了多少钱
				if _, ok := order.Odds["1"]; ok {
					oddsCode := strconv.Itoa(v + 1)
					if v, ok := order.Odds[oddsCode]; ok {
						ret += v * order.SingleBetAmount
					} else {
						beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
						return
					}
				} else {
					beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
				}
			}
		}
	}

	//计算反水 单注金额 * 反水 * 单注数量 -1
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//总结算
	order.Settlement = ret + order.RebateAmount
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//色波
func (o *BJ28) WinningAndLose_1(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_1(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	sum := openCode[0] + openCode[1] + openCode[2]

	c := color[sum]

	for _, v := range betNumbers {
		if c == v {
			winningBetNum++
			oddsCode := strconv.Itoa(v + 1)
			if v, ok := order.Odds[oddsCode]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
	}

	//计算反水 单注金额 * 反水 * 单注数量 -1
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//总结算
	order.Settlement = ret + order.RebateAmount
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//豹子
func (o *BJ28) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	var winningBetNum = 0
	//判断3个号码是否相同
	if openCode[0] == openCode[1] && openCode[0] == openCode[2] {
		winningBetNum = 1
	}

	//如果中奖数组<1 证明没有中奖
	if winningBetNum < 1 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum == 1 {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//中的注数
			winningBetNum = 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	order.WinningBetNum = winningBetNum
	order.Status = 1
}

//特码包三
func (o *BJ28) WinningAndLose_3(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_3(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum int

	sum := openCode[0] + openCode[1] + openCode[2]

	for _, v := range betNumbers {
		if v == sum {
			winningBetNum = 1
			break
		}
	}

	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

		//计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)

		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}

		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//特码
func (o *BJ28) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_3(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	sum := openCode[0] + openCode[1] + openCode[2]

	for _, v := range betNumbers {
		if v == sum {
			winningBetNum = 1
			oddsCode := strconv.Itoa(v + 1)
			if v, ok := order.Odds[oddsCode]; ok {
				ret += v * order.SingleBetAmount * 1
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
			break
		}
	}

	//计算反水 单注金额 * 反水 * 单注数量 -1
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//总结算
	order.Settlement = ret + order.RebateAmount
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *BJ28) PaserNormalBetNum_0(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于10(PK10玩法)
		if i < 0 || i > 9 {
			return false, nil
		}
		arrayInt = append(arrayInt, i)
	}

	//验证3:是否有重复数字(下注数字不能重复)
	if !o.CheckRepeatInt(arrayInt) {
		return false, nil
	}

	return true, arrayInt
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *BJ28) PaserNormalBetNum_1(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于10(PK10玩法)
		if i < 0 || i > 2 {
			return false, nil
		}
		arrayInt = append(arrayInt, i)
	}

	//验证3:是否有重复数字(下注数字不能重复)
	if !o.CheckRepeatInt(arrayInt) {
		return false, nil
	}

	return true, arrayInt
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *BJ28) PaserNormalBetNum_2(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}

		if i != 0 {
			return false, nil
		}
		arrayInt = append(arrayInt, i)
	}

	//验证3:是否有重复数字(下注数字不能重复)
	if !o.CheckRepeatInt(arrayInt) {
		return false, nil
	}

	return true, arrayInt
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *BJ28) PaserNormalBetNum_3(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于10(PK10玩法)
		if i < 0 || i > 27 {
			return false, nil
		}
		arrayInt = append(arrayInt, i)
	}

	//验证3:是否有重复数字(下注数字不能重复)
	if !o.CheckRepeatInt(arrayInt) {
		return false, nil
	}

	return true, arrayInt
}

//检查是否有重复数字符(int数组)
func (o *BJ28) CheckRepeatInt(array []int) bool {
	var newArray []int = []int{}
	for _, i := range array {
		if len(newArray) == 0 {
			newArray = append(newArray, i)
		} else {
			for k, v := range newArray {
				if i == v {
					return false
				}
				if k == len(newArray)-1 {
					newArray = append(newArray, i)
				}
			}
		}
	}
	return true
}
