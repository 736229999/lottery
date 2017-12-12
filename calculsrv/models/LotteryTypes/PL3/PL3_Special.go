package PL3

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

//和值数组
var sumArray = [28]int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 63, 69, 73, 75, 75, 73, 69, 63, 55, 45, 36, 28, 21, 15, 10, 6, 3, 1}

//二和值
var twoSumArray = [19]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

//分析订单(下注)
func (o *PL3) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {

	//2分析反水设置是否正确(所有类型都一样，所以放在这里)
	if order.Rebate > accountInfo.Rebate {
		beego.Debug("失败0")
		return false
	}

	//3分析单注金额有没有超过限制
	if order.SingleBetAmount > o.Settings[order.BetType].SingleLimit {
		beego.Debug("失败00")
		return false
	}

	switch order.BetType {
	case 0: //三星三星直选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserThreeDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		var singleBetNum = 1
		for _, v := range array {
			l := len(v)
			//数组数量不得大于10个元素或小于1个元素，应为三码直选每位只能是 0 - 9 这10个数字
			if l > 10 || l < 1 {
				//beego.Debug("失败")
				return false
			}
			//订单注数就是 3个数组相乘
			singleBetNum *= l
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 1: //三星三星和值
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumForThreeSum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		if l < 1 || l > 28 {
			return false
		}

		var singleBetNum = 0
		for _, v := range array {
			singleBetNum += sumArray[v]
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 2: //三星三星组三
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//数组数量不得大于10个元素或小于2个元素，应为组3复式 最多只能选10个数字 0 - 9 或最少选择2个数字
		if l > 10 || l < 2 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数 选择号码数 * （选择号码数 - 1）
		singleBetNum := l * (l - 1)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 3: //三星三星组六
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于3 大于10 应为组六玩法 最少选择3个 最多选择10个号码
		if l < 3 || l > 10 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(三不同号,就是所选数字个数的三三排列组合)
		singleBetNum := utils.AnalysisCombination(l, 3)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 4, 7: //前二前二直选, 后二后二直选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		var singleBetNum = 1
		for _, v := range array {
			l := len(v)
			//数组数量不得大于10个元素或小于1个元素
			if l > 10 || l < 1 {
				//beego.Debug("失败")
				return false
			}
			//分析订单注数 , 每位选择的个数相乘就是订单数
			singleBetNum *= l
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 5, 8: //前二前二和值, 后二后二和值
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumForTwoSum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		if l < 1 || l > 19 {
			return false
		}

		var singleBetNum = 0
		for _, v := range array {
			singleBetNum += twoSumArray[v]
		}

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
	case 6, 9: //前二前二组选 , 后二后二组选
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)

		if l < 2 || l > 10 {
			beego.Debug("失败")
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 2)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 10: //定位胆定位胆定位胆
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserThreeDigitBetNumAllowSpaces(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		//每个数组数量不得大于10个元素或小于1个元素，应为前二直选每一位 最多只能选11个数字 1 - 11 或最少选择1个数字
		//同时分析订单数,定位胆的特殊性选择了几个号码就是几注
		var singleBetNum = 0
		var tl = 0
		for _, v := range array {
			l := len(v)
			tl += l
			if l < 0 || l > 10 {
				//beego.Debug("失败")
				return false
			}
			singleBetNum += l
		}

		//5个数组的长度不得小于1 大于50
		if tl < 1 || tl > 30 {
			beego.Debug("失败")
			return false
		}

		order.SingleBetNum = singleBetNum
		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败1")
			return false
		}

	case 11: //不定胆不定胆一码
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)

		if l < 1 || l > 10 {
			beego.Debug("失败")
			return false
		}

		//6分析订单注数
		order.SingleBetNum = l

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 12: //不定胆不定胆二码
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)

		if l < 1 || l > 10 {
			beego.Debug("失败")
			return false
		}

		//6分析订单注数
		order.SingleBetNum = utils.AnalysisCombination(l, 2)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 13, 14: //前二直选大小单双, 后二直选大小单双
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNumForBigSmall(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 1 || l1 > 4 {
			return false
		}

		if l2 < 1 || l2 > 4 {
			return false
		}

		//l1 * l2 - 重复数
		order.SingleBetNum = l1 * l2

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 15: //两面百位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor15(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l < 1 || l > 4 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 前一直选 选了几个数就是几注
		order.SingleBetNum = l

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败1")
			return false
		}

	case 16: //两面十位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor16(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l < 1 || l > 3 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 前一直选 选了几个数就是几注
		order.SingleBetNum = l

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败1")
			return false
		}

	case 17, 18: //两面百位,两面和值
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor17(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l < 1 || l > 2 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 前一直选 选了几个数就是几注
		order.SingleBetNum = l

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败1")
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
func (o *PL3) SettlementOrders(orders []gb.Order, openCode string) {
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
func (o *PL3) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}

	//转下注号码为数组
	switch order.BetType {
	case 0: //三星三星直选
		o.WinningAndLose_0(order, openCode)
	case 1: //三星三星和值
		o.WinningAndLose_1(order, openCode)
	case 2: //三星三星组三
		o.WinningAndLose_2(order, openCode)
	case 3: //三星三星组六
		o.WinningAndLose_3(order, openCode)
	case 4: //前二前二直选
		o.WinningAndLose_4(order, openCode)
	case 5: //前二前二和值
		o.WinningAndLose_5(order, openCode)
	case 6: //前二前二组选
		o.WinningAndLose_6(order, openCode)
	case 7: //后二后二直选
		o.WinningAndLose_7(order, openCode)
	case 8: //后二后二和值
		o.WinningAndLose_8(order, openCode)
	case 9: //后二后二组选
		o.WinningAndLose_9(order, openCode)
	case 10: //定位胆定位胆定位胆
		o.WinningAndLose_10(order, openCode)
	case 11: //不定胆不定胆一码
		o.WinningAndLose_11(order, openCode)
	case 12: //不定胆不定胆二码
		o.WinningAndLose_12(order, openCode)
	case 13: //大小单双大小单双百十
		o.WinningAndLose_13(order, openCode)
	case 14: //大小单双大小单双十个
		o.WinningAndLose_14(order, openCode)
	case 15: //大小单双大小单双十个
		o.WinningAndLose_15(order, openCode)
	case 16: //大小单双大小单双十个
		o.WinningAndLose_16(order, openCode)
	case 17: //大小单双大小单双十个
		o.WinningAndLose_17(order, openCode)
	case 18: //大小单双大小单双十个
		o.WinningAndLose_18(order, openCode)
	default:
		beego.Debug("失败")
		return false
	}
	return true
}

//判断是否中奖 (0 直选)(返回用户输赢情况)
func (o *PL3) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers = 0
	var winningBetNum = 0
	//判断有没有 按位开出下注的号码
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0

	for _, v := range betNumbers[0] {
		if v == openCode[0] {
			flag_1 = 1
			break
		}
	}

	if flag_1 == 1 {
		for _, v := range betNumbers[1] {
			if v == openCode[1] {
				flag_2 = 1
				break
			}
		}
	}

	if flag_2 == 1 {
		for _, v := range betNumbers[2] {
			if v == openCode[2] {
				flag_3 = 1
				break
			}
		}
	}
	//三位都相同才中奖
	if flag_1 == 1 && flag_2 == 1 && flag_3 == 1 {
		winningNumbers = 1
	}

	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		//前三直选只会有一注中奖
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (1 三星三星和值)(返回用户输赢情况)
func (o *PL3) WinningAndLose_1(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForThreeSum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算前三个开奖号码的和值
	sum := openCode[0] + openCode[1] + openCode[2]
	//循环下注号码是否中奖
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == sum {
			winningBetNum = 1
			break
		}
	}

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum == 1 {
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (2三星三星组三)(返回用户输赢情况)
func (o *PL3) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断下注号码中有几个开出(先看有没有开出两个相同的号码)

	var winningBetNum = 0
	var sameNumber = 0

	if openCode[0] == openCode[1] && openCode[0] == openCode[2] {
		winningBetNum = 0
	} else if openCode[0] == openCode[1] {
		for _, v := range betNumbers {
			if openCode[0] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers {
				if v == openCode[2] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[0] == openCode[2] {
		for _, v := range betNumbers {
			if openCode[0] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers {
				if v == openCode[1] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[1] == openCode[2] {
		for _, v := range betNumbers {
			if openCode[1] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers {
				if v == openCode[0] {
					winningBetNum = 1
					break
				}
			}
		}
	}

	//如果中奖数组<2 证明没有中奖
	if winningBetNum == 0 {
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (3 三星三星组六(返回用户输赢情况)
func (o *PL3) WinningAndLose_3(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers = 0
	var winningBetNum = 0
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers++
				break
			}
		}

	}

	//如果中奖数组<1 证明没有中奖
	if winningNumbers < 3 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 3 {
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (4 前二前二直选)(返回用户输赢情况)
func (o *PL3) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_0 = 0

	var winningNumbers = 0
	var winningBetNum = 0
	//第一名
	for _, v := range betNumbers[0] {
		if v == openCode[0] {
			flag_0 = 1
			break
		}
	}

	if flag_0 == 1 {
		//第二名
		for _, v := range betNumbers[1] {
			if v == openCode[1] {
				winningNumbers = 1
				break
			}
		}
	}

	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (5 前二前二和值)(返回用户输赢情况)
func (o *PL3) WinningAndLose_5(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForTwoSum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算前二个开奖号码的和值
	sum := openCode[0] + openCode[1]
	//循环下注号码是否中奖
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == sum {
			winningBetNum = 1
			break
		}
	}

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum == 1 {
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (6 前二前二组选)(返回用户输赢情况)
func (o *PL3) WinningAndLose_6(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	var winningNumbers = 0
	//判断下注号码中有几个开出
	var tmpOpenCode []int = openCode[:2]
	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}

	}

	//如果中奖数组<1 证明没有中奖
	if winningNumbers < 2 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 2 {
		//三不同号只有一注中奖
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (7 后二后二直选)(返回用户输赢情况)
func (o *PL3) WinningAndLose_7(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_3 = 0

	var winningNumbers = 0
	var winningBetNum = 0
	//第二名
	for _, v := range betNumbers[0] {
		if v == openCode[1] {
			flag_3 = 1
			break
		}
	}

	if flag_3 == 1 {
		//第三名
		for _, v := range betNumbers[1] {
			if v == openCode[2] {
				winningNumbers = 1
				break
			}
		}
	}

	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (8 后二后二和值)(返回用户输赢情况)
func (o *PL3) WinningAndLose_8(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForTwoSum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算后二个开奖号码的和值
	sum := openCode[1] + openCode[2]
	//循环下注号码是否中奖
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == sum {
			winningBetNum = 1
			break
		}
	}

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum == 1 {
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (9 后二后二组选)(返回用户输赢情况)
func (o *PL3) WinningAndLose_9(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	var winningNumbers = 0
	//判断下注号码中有几个开出(注意采用这种方式才不会修改到原切片)
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1], openCode[2])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}

	}

	//如果中奖数组<1 证明没有中奖
	if winningNumbers < 2 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 2 {
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
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (10 定位胆定位胆定位胆)(返回用户输赢情况)
func (o *PL3) WinningAndLose_10(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNumAllowSpaces(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers = 0
	var winningBetNum = 0
	//判断中奖注数(注意:每一位只会有一个数字中奖)
	//第一名
	for _, v := range betNumbers[0] {
		if v == openCode[0] {
			winningNumbers++
			break
		}
	}
	//第二名
	for _, v := range betNumbers[1] {
		if v == openCode[1] {
			winningNumbers++
			break
		}
	}
	//第三名
	for _, v := range betNumbers[2] {
		if v == openCode[2] {
			winningNumbers++
			break
		}
	}

	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers > 0 && winningNumbers < 4 {
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningNumbers)
		order.RebateAmount = ret

		//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
		winningBetNum = winningNumbers
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningNumbers)
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}

	}
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (11 不定胆不定胆一码)(返回用户输赢情况)
func (o *PL3) WinningAndLose_11(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var winningBetNum = 0

	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningBetNum++
				break
			}
		}
	}

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (12 不定胆不定胆二码)(返回用户输赢情况)
func (o *PL3) WinningAndLose_12(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var winningNumbers = 0
	var winningBetNum = 0
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	if winningNumbers < 2 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 2 {
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
	} else if winningNumbers == 3 {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-3)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 3
			//中的注数
			winningBetNum = 3
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (13 大小单双大小单双百十)(返回用户输赢情况)
func (o *PL3) WinningAndLose_13(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNumForBigSmall(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0

	//0 1 2 3 大小单双
	//第1位 默认大
	var B1 = 0
	//第2位 默认大
	var B2 = 0
	//第1位 默认单
	var D1 = 2
	//第2位 默认单
	var D2 = 2

	if openCode[0] < 5 {
		B1 = 1
	}
	if openCode[0]%2 == 0 {
		D1 = 3
	}

	if openCode[1] < 5 {
		B2 = 1
	}
	if openCode[1]%2 == 0 {
		D2 = 3
	}

	var l1Winning = 0
	var l2Winning = 0
	for _, v := range betNumbers[0] {
		if B1 == v {
			l1Winning++
		}
		if D1 == v {
			l1Winning++
		}
	}

	for _, v := range betNumbers[1] {
		if B2 == v {
			l2Winning++
		}
		if D2 == v {
			l2Winning++
		}
	}

	winningBetNum = l1Winning * l2Winning

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 {
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (14 大小单双大小单双十个)(返回用户输赢情况)
func (o *PL3) WinningAndLose_14(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNumForBigSmall(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0

	//0 1 2 3 大小单双
	//第1位 默认大
	var B1 = 0
	//第2位 默认大
	var B2 = 0
	//第1位 默认单
	var D1 = 2
	//第2位 默认单
	var D2 = 2

	if openCode[1] < 5 {
		B1 = 1
	}
	if openCode[1]%2 == 0 {
		D1 = 3
	}

	if openCode[2] < 5 {
		B2 = 1
	}
	if openCode[2]%2 == 0 {
		D2 = 3
	}

	var l1Winning = 0
	var l2Winning = 0
	for _, v := range betNumbers[0] {
		if B1 == v {
			l1Winning++
		}
		if D1 == v {
			l1Winning++
		}
	}

	for _, v := range betNumbers[1] {
		if B2 == v {
			l2Winning++
		}
		if D2 == v {
			l2Winning++
		}
	}

	winningBetNum = l1Winning * l2Winning

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 {
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//更新order
			order.Settlement = ret
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 15 两面百位
func (o *PL3) WinningAndLose_15(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor15(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [4]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6, 龙百 7, 虎个 8, 百个和 9
	if o.openCode[0] > 4 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if o.openCode[0]%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	if o.openCode[0] > o.openCode[1] {
		resultCode[2] = 4
	} else if o.openCode[0] < o.openCode[1] {
		resultCode[2] = 5
	} else {
		resultCode[2] = 6
	}

	if o.openCode[0] > o.openCode[2] {
		resultCode[3] = 7
	} else if o.openCode[0] < o.openCode[2] {
		resultCode[3] = 8
	} else {
		resultCode[3] = 9
	}

	var winningBetNum = 0

	for _, v := range resultCode {
		for _, i := range betNumbers {
			if v == i {
				winningBetNum++
				//计算中了多少钱
				if v, ok := order.Odds[strconv.Itoa(i+1)]; ok {
					ret += v * order.SingleBetAmount
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
			}
		}
	}

	//开出的第一位号码没有在下注号码里面

	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)

	order.Settlement = order.RebateAmount + ret

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 16 两面十位
func (o *PL3) WinningAndLose_16(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor16(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6, 龙百 7, 虎个 8, 百个和 9
	if o.openCode[1] > 4 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if o.openCode[1]%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	if o.openCode[1] > o.openCode[2] {
		resultCode[2] = 4
	} else if o.openCode[1] < o.openCode[2] {
		resultCode[2] = 5
	} else {
		resultCode[2] = 6
	}

	var winningBetNum = 0

	for _, v := range resultCode {
		for _, i := range betNumbers {
			if v == i {
				winningBetNum++
				//计算中了多少钱
				if v, ok := order.Odds[strconv.Itoa(i+1)]; ok {
					ret += v * order.SingleBetAmount
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
			}
		}
	}

	//开出的第一位号码没有在下注号码里面

	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)

	order.Settlement = order.RebateAmount + ret

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 17 两面个位
func (o *PL3) WinningAndLose_17(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor17(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1 , 单 2, 双 3,
	if o.openCode[2] > 4 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if o.openCode[2]%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	var winningBetNum = 0

	for _, v := range resultCode {
		for _, i := range betNumbers {
			if v == i {
				winningBetNum++
				//计算中了多少钱
				if v, ok := order.Odds[strconv.Itoa(i+1)]; ok {
					ret += v * order.SingleBetAmount
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
			}
		}
	}

	//开出的第一位号码没有在下注号码里面

	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)

	order.Settlement = order.RebateAmount + ret

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 18 两面总和
func (o *PL3) WinningAndLose_18(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor17(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	var winningBetNum = 0

	sum := o.openCode[0] + o.openCode[1] + o.openCode[2]
	//大 0 , 小 1 , 单 2, 双 3,
	if sum > 13 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if sum%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	for _, v := range resultCode {
		for _, i := range betNumbers {
			if v == i {
				winningBetNum++
				//计算中了多少钱
				if v, ok := order.Odds[strconv.Itoa(i+1)]; ok {
					ret += v * order.SingleBetAmount
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
			}
		}
	}

	//开出的第一位号码没有在下注号码里面

	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)

	order.Settlement = order.RebateAmount + ret

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *PL3) PaserNormalBetNum(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(PL3玩法)
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)(1)
func (o *PL3) PaserNormalBetNumForThreeSum(betNum string) (bool, []int) {
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

//解析下注号码,得到注数二维数组(用于有两位选择数字的情况,列入前二直选)
func (o *PL3) PaserTwoDigitBetNum(betNum string) (bool, [][]int) {
	//分割下注位(这里有个问题 如果 分割不成功会返回什么?等到测试来验证)
	array := strings.Split(betNum, ";")
	if len(array) != 2 {
		return false, nil
	}

	//分割位数
	var bInt [][]int
	for _, v := range array {
		b := strings.Split(v, ",")
		var b2Int []int
		for _, j := range b {
			i, err := strconv.Atoi(j)
			if err != nil {
				return false, nil
			}
			//每位数字不能小于1和大于9
			if i < 0 || i > 9 {
				return false, nil
			}
			b2Int = append(b2Int, i)
		}
		bInt = append(bInt, b2Int)
	}

	for _, v := range bInt {
		if !o.CheckRepeatInt(v) {
			return false, nil
		}
	}
	return true, bInt
}

//解析下注号码,得到注数二维数组(13,14)
func (o *PL3) PaserTwoDigitBetNumForBigSmall(betNum string) (bool, [][]int) {
	//分割下注位(这里有个问题 如果 分割不成功会返回什么?等到测试来验证)
	array := strings.Split(betNum, ";")
	if len(array) != 2 {
		return false, nil
	}

	//分割位数
	var bInt [][]int
	for _, v := range array {
		b := strings.Split(v, ",")
		var b2Int []int
		for _, j := range b {
			i, err := strconv.Atoi(j)
			if err != nil {
				return false, nil
			}
			//每位数字不能小于1和大于3(pk10 大小单双玩法)
			if i < 0 || i > 3 {
				return false, nil
			}
			b2Int = append(b2Int, i)
		}
		bInt = append(bInt, b2Int)
	}

	for _, v := range bInt {
		if !o.CheckRepeatInt(v) {
			return false, nil
		}
	}
	return true, bInt
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)(5, 8)
func (o *PL3) PaserNormalBetNumForTwoSum(betNum string) (bool, []int) {
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
		if i < 0 || i > 18 {
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

//解析3位下注号码,得到注数二维数组(用于有3位选择数字的情况,例如:前三直选)
func (o *PL3) PaserThreeDigitBetNum(betNum string) (bool, [][]int) {
	//分割下注位
	array := strings.Split(betNum, ";")
	if len(array) != 3 {
		return false, nil
	}
	//分割位数
	var bInt [][]int
	for _, v := range array {
		b := strings.Split(v, ",")
		var b2Int []int
		for _, j := range b {
			i, err := strconv.Atoi(j)
			if err != nil {
				return false, nil
			}
			//每位数字不能小于0和大于9(PL3玩法)
			if i < 0 || i > 9 {
				return false, nil
			}
			b2Int = append(b2Int, i)
		}
		bInt = append(bInt, b2Int)
	}

	for _, v := range bInt {
		if !o.CheckRepeatInt(v) {
			return false, nil
		}
	}
	return true, bInt
}

//解析3位下注号码,得到注数三维数组(用于有3位选择数字的情况,例如:定位胆)(允许空位数)
func (o *PL3) PaserThreeDigitBetNumAllowSpaces(betNum string) (bool, [][]int) {
	//分割下注位 由于定位胆的特殊性这里要做特殊处理, ;号分割以后会出现空的字段 只要总分个数为5 那么这种空字段的情况是可以出现的
	array := strings.Split(betNum, ";")
	if len(array) != 3 {
		return false, nil
	}
	//分割位数
	var bInt [][]int
	for _, v := range array {
		b := strings.Split(v, ",")
		var b2Int []int
		for _, j := range b {
			if j == "" {
				continue
			}
			i, err := strconv.Atoi(j)
			if err != nil {
				return false, nil
			}
			//每位数字不能小于1和大于10(PK10玩法)
			if i < 0 || i > 9 {
				return false, nil
			}
			b2Int = append(b2Int, i)
		}
		bInt = append(bInt, b2Int)
	}

	for _, v := range bInt {
		if !o.CheckRepeatInt(v) {
			return false, nil
		}
	}
	return true, bInt
}

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *PL3) PaserNormalBetNumFor15(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(PL3玩法)
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
func (o *PL3) PaserNormalBetNumFor16(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(PL3玩法)
		if i < 0 || i > 6 {
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
func (o *PL3) PaserNormalBetNumFor17(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(PL3玩法)
		if i < 0 || i > 3 {
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
func (o *PL3) CheckRepeatInt(array []int) bool {
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
