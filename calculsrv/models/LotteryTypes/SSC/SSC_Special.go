package SSC

import (
	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/Order"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"

	"calculsrv/models/acmgr"
	"common/utils"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//和值数组
var sumArray = [28]int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 63, 69, 73, 75, 75, 73, 69, 63, 55, 45, 36, 28, 21, 15, 10, 6, 3, 1}

//二和值
var twoSumArray = [19]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

//三号跨度
var threeSpan = [10]int{10, 54, 96, 126, 144, 150, 144, 126, 96, 54}

//两号跨度
var twoSpan = [10]int{10, 18, 16, 14, 12, 10, 8, 6, 4, 2}

//分析订单(下注)
func (o *SSC) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {
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
	//------------------------------------------------------------------------------------------------------------------------------------------
	switch order.BetType {
	case 0: //五星直选复试
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNum(order.BetNums)
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

	case 1: //五星直选组合
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNum(order.BetNums)
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

		//应为组合选择5个数一组就算5注
		singleBetNum *= 5

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

	case 2: //五星组选组选120
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		if l < 5 || l > 10 {
			return false
		}

		//分析订单注数 , 选择的数字个数进行5个数的排列组合
		order.SingleBetNum = utils.AnalysisCombination(l, 5)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 3: //五星组选组选60
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 1 || l1 > 10 {
			return false
		}

		if l2 < 3 || l2 > 10 {
			return false
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum++
				}
			}
		}

		//l2 进行3的组合 * l1 - 重复数*组合差值(组合差值 = l2的3组合数 - (l2 - 1)的3组合数

		//6分析订单注数
		order.SingleBetNum = utils.AnalysisCombination(l2, 3)*l1 - utils.CombinationDifference(l2, 3)*repateNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 4: //五星组选组选30
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 2 || l1 > 10 {
			return false
		}

		if l2 < 1 || l2 > 10 {
			return false
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum++
				}
			}
		}

		//6分析订单注数
		order.SingleBetNum = utils.AnalysisCombination(l1, 2)*l2 - utils.CombinationDifference(l1, 2)*repateNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 5: //五星组选组选20
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 1 || l1 > 10 {
			return false
		}

		if l2 < 2 || l2 > 10 {
			return false
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum++
				}
			}
		}

		//l2 进行2的组合 * l1 - 重复数*组合差值(组合差值 = l2的2组合数 - (l2 - 1)的2组合数
		order.SingleBetNum = utils.AnalysisCombination(l2, 2)*l1 - utils.CombinationDifference(l2, 2)*repateNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 6, 7: //五星组选组选10,五星组选组选5
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 1 || l1 > 10 {
			return false
		}

		if l2 < 1 || l2 > 10 {
			return false
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum++
				}
			}
		}

		//l1 * l2 - 重复数
		order.SingleBetNum = l1*l2 - repateNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 8, 14: //前四直选复试,后四直选复试
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFourDigitBetNum(order.BetNums)
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

	case 9, 15: //前四直选组合,后四直选组合
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFourDigitBetNum(order.BetNums)
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

		//应为组合选择4个数一组就算4注
		singleBetNum *= 4

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

	case 10, 16: //前四组选组选24,后四组选组选24
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		if l < 4 || l > 10 {
			return false
		}

		//分析订单注数 , 选择的数字个数进行5个数的排列组合
		order.SingleBetNum = utils.AnalysisCombination(l, 4)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 11, 17: //前四组选组选12,后四组选组选12
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 1 || l1 > 10 {
			return false
		}

		if l2 < 2 || l2 > 10 {
			return false
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum++
				}
			}
		}

		//l2 进行2的组合 * l1 - 重复数*组合差值(组合差值 = l2的2组合数 - (l2 - 1)的2组合数
		//6分析订单注数
		order.SingleBetNum = utils.AnalysisCombination(l2, 2)*l1 - utils.CombinationDifference(l2, 2)*repateNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 12, 18: //前四组选组选6,后四组选组选6
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		if l < 2 || l > 10 {
			return false
		}

		//分析订单注数 , 选择的数字个数进行2个数的排列组合
		order.SingleBetNum = utils.AnalysisCombination(l, 2)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 13, 19: //前四组选组选4,后四组选组选4
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 1 || l1 > 10 {
			return false
		}

		if l2 < 1 || l2 > 10 {
			return false
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum++
				}
			}
		}

		//l1 * l2 - 重复数
		order.SingleBetNum = l1*l2 - repateNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 20, 24, 28: //前三直选复试,中三直选复试,后三直选复试
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserThreeDigitBetNum(order.BetNums)
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

	case 21, 25, 29: //前三直选和值,中三直选和值,后三直选和值
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

	case 22, 26, 30: //前三组选组三,中三组选组三,后三组选组三
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

	case 23, 27, 31: //前三组选组六,中三组选组六,后三组选组六
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

	case 32, 36: //前二直选复式, 后二直选复式
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

	case 33, 37: //前二直选和值, 后二直选和值
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
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 34, 38: //前二直选大小单双, 后二直选大小单双
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
	case 35, 39: //前二组选复式 , 后二组选复式
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

	case 40: //定位胆定位胆定位胆
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
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
		if tl < 1 || tl > 50 {
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

	case 41, 42, 43, 58, 59, 60, 61: //三星不定胆一码前三, 三星不定胆一码中三, 三星不定胆一码后三
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
	case 44, 45, 46: //三星不定胆二码前三, 三星不定胆二码中三, 三星不定胆二码后三
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
	case 47: //任选任选二复选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		l := len(array)
		var singleBetNum = 0
		var tl = 0

		for i := 0; i < l; i++ {
			for j := i + 1; j < l; j++ {
				singleBetNum += len(array[i]) * len(array[j])
			}
			tl += len(array[i])
		}

		if singleBetNum < 1 || singleBetNum > 1000 {
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

	case 48: //任选任选二组选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 2 || l1 > 5 {
			return false
		}

		if l2 < 2 || l2 > 10 {
			return false
		}

		//由于这个玩法的特殊性,第一个数组 万千百十个由 0.1.2.3.4代替,所以这里加判断
		for _, v := range array[0] {
			if v < 0 || v > 4 {
				return false
			}
		}

		order.SingleBetNum = utils.AnalysisCombination(l1, 2) * utils.AnalysisCombination(l2, 2)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 49: //任选任选三复选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		var arrayLen []int

		//判断每一位的个数不能超过10个数
		for _, v := range array {
			l := len(v)
			if l > 10 {
				return false
			}
			if l > 0 {
				arrayLen = append(arrayLen, l)
			}

		}

		//有几位中选择了号码
		l := len(arrayLen)
		if l < 3 || l > 5 {
			return false
		}

		var singleBetNum = 0
		singleBetNum += arrayLen[0] * arrayLen[1] * arrayLen[2]
		if l > 3 {
			singleBetNum += arrayLen[0]*arrayLen[1]*arrayLen[3] + arrayLen[0]*arrayLen[2]*arrayLen[3] + arrayLen[1]*arrayLen[2]*arrayLen[3]
		}
		if l > 4 {
			singleBetNum += arrayLen[0]*arrayLen[1]*arrayLen[4] + arrayLen[0]*arrayLen[2]*arrayLen[4] + arrayLen[0]*arrayLen[3]*arrayLen[4] + arrayLen[1]*arrayLen[2]*arrayLen[4] + arrayLen[1]*arrayLen[3]*arrayLen[4] + arrayLen[2]*arrayLen[3]*arrayLen[4]
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

	case 50: //任选任选三组三
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 3 || l1 > 5 {
			return false
		}

		if l2 < 2 || l2 > 10 {
			return false
		}

		//由于这个玩法的特殊性,第一个数组 万千百十个由 0.1.2.3.4代替,所以这里加判断
		for _, v := range array[0] {
			if v < 0 || v > 4 {
				return false
			}
		}

		order.SingleBetNum = utils.AnalysisCombination(l1, 3) * utils.AnalysisCombination(l2, 2) * 2

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}

	case 51: //任选任选三组六
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l1 := len(array[0])
		l2 := len(array[1])

		if l1 < 3 || l1 > 5 {
			return false
		}

		if l2 < 3 || l2 > 10 {
			return false
		}

		//由于这个玩法的特殊性,第一个数组 万千百十个由 0.1.2.3.4代替,所以这里加判断
		for _, v := range array[0] {
			if v < 0 || v > 4 {
				return false
			}
		}

		order.SingleBetNum = utils.AnalysisCombination(l1, 3) * utils.AnalysisCombination(l2, 3)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 52: //任选任选四复选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		var arrayLen []int

		//判断每一位的个数不能超过10个数
		for _, v := range array {
			l := len(v)
			if l > 10 {
				return false
			}
			if l > 0 {
				arrayLen = append(arrayLen, l)
			}

		}

		//有几位中选择了号码
		l := len(arrayLen)
		if l < 4 || l > 5 {
			return false
		}

		var singleBetNum = 0

		singleBetNum += arrayLen[0] * arrayLen[1] * arrayLen[2] * arrayLen[3]

		if l > 4 {
			singleBetNum += arrayLen[0]*arrayLen[1]*arrayLen[2]*arrayLen[4] + arrayLen[0]*arrayLen[1]*arrayLen[3]*arrayLen[4] + arrayLen[0]*arrayLen[2]*arrayLen[3]*arrayLen[4] + arrayLen[1]*arrayLen[2]*arrayLen[3]*arrayLen[4]
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

	case 53, 54, 55: //跨度三星跨度前三, 跨度三星跨度中三, 跨度三星跨度后三
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

		var singleBetNum = 0
		for _, v := range array {
			singleBetNum += threeSpan[v]
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 56, 57: //跨度二星跨度前二, 跨度二星跨度后二
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

		var singleBetNum = 0
		for _, v := range array {
			singleBetNum += twoSpan[v]
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 62, 63, 64, 65, 66, 67, 68, 69, 70, 71: //龙虎龙虎万千, 龙虎龙虎万百, 龙虎龙虎万十, 龙虎龙虎万个, 龙虎龙虎千百,  龙虎龙虎千十,  龙虎龙虎千个, 龙虎龙虎百十, 龙虎龙虎百个, 龙虎龙虎十个
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)

		if l != 1 {
			beego.Debug("失败")
			return false
		}

		//由于这个玩法的特殊性,每位下注数只能是0,1,2
		if array[0] < 0 || array[0] > 2 {
			return false
		}

		//6分析订单注数
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

	case 72: //两面万位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor72(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l < 1 || l > 6 {
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

	case 73: //两面千位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor73(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l < 1 || l > 5 {
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

	case 74: //两面百位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor74(order.BetNums)
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

	case 75: //两面十位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor75(order.BetNums)
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

	case 76, 77: //两面十位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor76(order.BetNums)
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

	case 78: //牛牛
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor78(order.BetNums)
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

//解析5位下注号码,得到注数五维数组(用于有5位选择数字的情况)
func (o *SSC) PaserFiveDigitBetNum(betNum string) (bool, [][]int) {
	//分割下注位 由于定位胆的特殊性这里要做特殊处理, ;号分割以后会出现空的字段 只要总分个数为5 那么这种空字段的情况是可以出现的
	array := strings.Split(betNum, ";")
	if len(array) != 5 {
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
			//每位数字不能小于1和大于10(SSC玩法)
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

//解析5位下注号码,得到注数五维数组(用于有5位选择数字的情况,例如:定位胆5)
func (o *SSC) PaserFiveDigitBetNumAllowSpaces(betNum string) (bool, [][]int) {
	//分割下注位 由于定位胆的特殊性这里要做特殊处理, ;号分割以后会出现空的字段 只要总分个数为5 那么这种空字段的情况是可以出现的
	array := strings.Split(betNum, ";")
	if len(array) != 5 {
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

//解析4位下注号码,得到注数四维数组(用于有4位选择数字的情况)
func (o *SSC) PaserFourDigitBetNum(betNum string) (bool, [][]int) {
	//分割下注位 由于定位胆的特殊性这里要做特殊处理, ;号分割以后会出现空的字段 只要总分个数为5 那么这种空字段的情况是可以出现的
	array := strings.Split(betNum, ";")
	if len(array) != 4 {
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
			//每位数字不能小于1和大于10(SSC玩法)
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

//解析3位下注号码,得到注数二维数组(用于有3位选择数字的情况,例如:前三)
func (o *SSC) PaserThreeDigitBetNum(betNum string) (bool, [][]int) {
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
			//每位数字不能小于1和大于10(SSC玩法)
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

//解析下注号码,得到注数二维数组(用于有两位选择数字的情况,列入前二直选)
func (o *SSC) PaserTwoDigitBetNum(betNum string) (bool, [][]int) {
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
			//每位数字不能小于1和大于10(pk10玩法)
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

//解析下注号码,得到注数二维数组(用于有两位选择数字的情况,列入前二直选)
func (o *SSC) PaserTwoDigitBetNumForBigSmall(betNum string) (bool, [][]int) {
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *SSC) PaserNormalBetNum(betNum string) (bool, []int) {
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)(33, 37)
func (o *SSC) PaserNormalBetNumForTwoSum(betNum string) (bool, []int) {
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)(21, 25, 29)
func (o *SSC) PaserNormalBetNumForThreeSum(betNum string) (bool, []int) {
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *SSC) PaserNormalBetNumFor72(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(快3玩法)
		if i < 0 || i > 15 {
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
func (o *SSC) PaserNormalBetNumFor73(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(快3玩法)
		if i < 0 || i > 12 {
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
func (o *SSC) PaserNormalBetNumFor74(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(快3玩法)
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
func (o *SSC) PaserNormalBetNumFor75(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(快3玩法)
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
func (o *SSC) PaserNormalBetNumFor76(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(快3玩法)
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

func (o *SSC) PaserNormalBetNumFor78(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于6(快3玩法)
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
func (o *SSC) CheckRepeatInt(array []int) bool {
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

// //结算这个彩种当期所有订单-----------------------
// func (o *SSC) settlementOrders() {
// 	//从数据库中获取这一期这个采种所有的订单
// 	orders := dbmgr.GetLotteryOrderRecord(o.gameTag, o.currentExpect)
// 	l := len(orders)
// 	//没有订单
// 	if l < 1 {
// 		return
// 	}
// 	//金额流水数组
// 	BalanceRecourds := []BalanceRecordMgr.BalanceRecord{}
// 	for i := 0; i < l; i++ {

// 		o.settlementOrder(&orders[i])

// 		order := orders[i]
// 		//获得这条订单的账户信息
// 		accountInfo := acmgr.AccountInfo{}
// 		err := accountInfo.Init(order.AccountName)
// 		if err != nil {
// 			beego.Emergency(err)
// 			return
// 		}
// 		//生成流水记录
// 		BalanceRecourd := BalanceRecordMgr.BalanceRecord{}
// 		BalanceRecourd.Serial_Number = Order.Instance().GetOrderNumber()
// 		BalanceRecourd.Account_name = accountInfo.Account_Name
// 		BalanceRecourd.Money_Before = accountInfo.Money
// 		BalanceRecourd.Money = order.Settlement
// 		BalanceRecourd.Money_After = accountInfo.Money + order.Settlement //注意 ： 投注只有减钱，而结算只有加钱
// 		BalanceRecourd.Gap_Money = 0
// 		BalanceRecourd.Type = 1    //1订单(我这里只有1)
// 		BalanceRecourd.Subitem = 2 //1投注, 2结算
// 		BalanceRecourd.Trading_Time = utils.GetNowUTC8Time().Unix()
// 		BalanceRecourd.Status = 1
// 		BalanceRecourd.Order_Number = order.OrderNumber
// 		BalanceRecourds = append(BalanceRecourds, BalanceRecourd)
// 		//结算只有加钱。。。。
// 		accountInfo.AddMoney(order.Settlement)
// 		//更新用户信息
// 		err = accountInfo.UpdataDb()
// 		if err != nil {
// 			beego.Emergency(err)
// 			return
// 		}
// 	}
// 	//插~!~
// 	bl := len(BalanceRecourds)
// 	if bl == 1 {
// 		dbmgr.InsertBalanceRecord(BalanceRecourds[0])
// 	} else if bl > 1 {
// 		dbmgr.BulkInsertBalanceRecord(BalanceRecourds)
// 	}
// }

//根据指定期数信息结算,这个函数上线稳定后必须和上面的函数合一,现在不敢轻易修改原流程
//结算这个彩种当期所有订单
func (o *SSC) SettlementOrders(orders []gb.Order, openCode string) {
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
// func (o *SSC) settlementOrder(order *gb.Order) {
// 	if order.Status == 1 {
// 		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
// 		return
// 	}
// 	//转下注号码为数组
// 	switch order.BetType {
// 	case 0: //五星直选复式
// 		o.WinningAndLose_0(order)
// 	case 1: //五星直选组合
// 		o.WinningAndLose_1(order)
// 	case 2: //五星组选组选120
// 		o.WinningAndLose_2(order)
// 	case 3: //五星组选组选60
// 		o.WinningAndLose_3(order)
// 	case 4: //五星组选组选30
// 		o.WinningAndLose_4(order)
// 	case 5: //五星组选组选20
// 		o.WinningAndLose_5(order)
// 	case 6: //五星组选组选10
// 		o.WinningAndLose_6(order)
// 	case 7: //五星组选组选5
// 		o.WinningAndLose_7(order)
// 	case 8: //前四直选复式
// 		o.WinningAndLose_8(order)
// 	case 9: //前四直选组合
// 		o.WinningAndLose_9(order)
// 	case 10: //前四组选组选24
// 		o.WinningAndLose_10(order)
// 	case 11: //前四组选组选12
// 		o.WinningAndLose_11(order)
// 	case 12: //前四组选组选6
// 		o.WinningAndLose_12(order)
// 	case 13: //前四组选组选4
// 		o.WinningAndLose_13(order)
// 	case 14: //后四直选复式
// 		o.WinningAndLose_14(order)
// 	case 15: //后四直选组合
// 		o.WinningAndLose_15(order)
// 	case 16: //后四组选组选24
// 		o.WinningAndLose_16(order)
// 	case 17: //后四组选组选12
// 		o.WinningAndLose_17(order)
// 	case 18: //后四组选组选6
// 		o.WinningAndLose_18(order)
// 	case 19: //后四组选组选4
// 		o.WinningAndLose_19(order)
// 	case 20: //前三直选复式
// 		o.WinningAndLose_20(order)
// 	case 21: //前三直选和值
// 		o.WinningAndLose_21(order)
// 	case 22: //前三组选组三
// 		o.WinningAndLose_22(order)
// 	case 23: //前三组选组六
// 		o.WinningAndLose_23(order)
// 	case 24: //中三直选复式
// 		o.WinningAndLose_24(order)
// 	case 25: //中三直选和值
// 		o.WinningAndLose_25(order)
// 	case 26: //中三组选组三
// 		o.WinningAndLose_26(order)
// 	case 27: //中三组选组六
// 		o.WinningAndLose_27(order)
// 	case 28: //后三直选复式
// 		o.WinningAndLose_28(order)
// 	case 29: //后三直选和值
// 		o.WinningAndLose_29(order)
// 	case 30: //后三组选组三
// 		o.WinningAndLose_30(order)
// 	case 31: //后三组选组六
// 		o.WinningAndLose_31(order)
// 	case 32: //前二直选复式
// 		o.WinningAndLose_32(order)
// 	case 33: //前二直选和值
// 		o.WinningAndLose_33(order)
// 	case 34: //前二直选大小单双
// 		o.WinningAndLose_34(order)
// 	case 35: //前二组选复式
// 		o.WinningAndLose_35(order)
// 	case 36: //后二直选复式
// 		o.WinningAndLose_36(order)
// 	case 37: //后二直选和值
// 		o.WinningAndLose_37(order)
// 	case 38: //后二直选大小单双
// 		o.WinningAndLose_38(order)
// 	case 39: //后二组选复式
// 		o.WinningAndLose_39(order)
// 	case 40: //定位胆定位胆定位胆
// 		o.WinningAndLose_40(order)
// 	case 41: //三星不定胆一码前三
// 		o.WinningAndLose_41(order)
// 	case 42: //三星不定胆一码中三
// 		o.WinningAndLose_42(order)
// 	case 43: //三星不定胆一码后三
// 		o.WinningAndLose_43(order)
// 	case 44: //三星不定胆二码前三
// 		o.WinningAndLose_44(order)
// 	case 45: //三星不定胆二码中三
// 		o.WinningAndLose_45(order)
// 	case 46: //三星不定胆二码后三
// 		o.WinningAndLose_46(order)
// 	case 47: //任选任选二复选
// 		o.WinningAndLose_47(order)
// 	case 48: //任选任选二组选
// 		o.WinningAndLose_48(order)
// 	case 49: //任选任选三复选
// 		o.WinningAndLose_49(order)
// 	case 50: //任选任选三组三
// 		o.WinningAndLose_50(order)
// 	case 51: //任选任选三组六
// 		o.WinningAndLose_51(order)
// 	case 52: //任选任选四复选
// 		o.WinningAndLose_52(order)
// 	case 53: //跨度三星跨度前三
// 		o.WinningAndLose_53(order)
// 	case 54: //跨度三星跨度中三
// 		o.WinningAndLose_54(order)
// 	case 55: //跨度三星跨度后三
// 		o.WinningAndLose_55(order)
// 	case 56: //跨度二星跨度前二
// 		o.WinningAndLose_56(order)
// 	case 57: //跨度二星跨度后二
// 		o.WinningAndLose_57(order)
// 	case 58: //趣味趣味一帆风顺
// 		o.WinningAndLose_58(order)
// 	case 59: //趣味趣味好事成双
// 		o.WinningAndLose_59(order)
// 	case 60: //趣味趣味三星报喜
// 		o.WinningAndLose_60(order)
// 	case 61: //趣味趣味四季发财
// 		o.WinningAndLose_61(order)
// 	case 62: //龙虎龙虎万千
// 		o.WinningAndLose_62(order)
// 	case 63: //龙虎龙虎万百
// 		o.WinningAndLose_63(order)
// 	case 64: //龙虎龙虎万十
// 		o.WinningAndLose_64(order)
// 	case 65: //龙虎龙虎万个
// 		o.WinningAndLose_65(order)
// 	case 66: //龙虎龙虎千百
// 		o.WinningAndLose_66(order)
// 	case 67: //龙虎龙虎千十
// 		o.WinningAndLose_67(order)
// 	case 68: //龙虎龙虎千个
// 		o.WinningAndLose_68(order)
// 	case 69: //龙虎龙虎百十
// 		o.WinningAndLose_69(order)
// 	case 70: //龙虎龙虎百个
// 		o.WinningAndLose_70(order)
// 	case 71: //龙虎龙虎十个
// 		o.WinningAndLose_71(order)
// 	default:
// 		return
// 		beego.Debug("失败")
// 	}
// 	//更新订单,上线稳定后,改为批量更新订单,放在外面的函数更新
// 	order.OpenCode = o.openCodeString
// 	dbmgr.UpdateOrder(order)
// }

//指定一个开奖记录,结算一个订单
//结算这个彩种一个订单
func (o *SSC) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}
	//转下注号码为数组
	switch order.BetType {
	case 0: //五星直选复式
		o.WinningAndLose_0(order, openCode)
	case 1: //五星直选组合
		o.WinningAndLose_1(order, openCode)
	case 2: //五星组选组选120
		o.WinningAndLose_2(order, openCode)
	case 3: //五星组选组选60
		o.WinningAndLose_3(order, openCode)
	case 4: //五星组选组选30
		o.WinningAndLose_4(order, openCode)
	case 5: //五星组选组选20
		o.WinningAndLose_5(order, openCode)
	case 6: //五星组选组选10
		o.WinningAndLose_6(order, openCode)
	case 7: //五星组选组选5
		o.WinningAndLose_7(order, openCode)
	case 8: //前四直选复式
		o.WinningAndLose_8(order, openCode)
	case 9: //前四直选组合
		o.WinningAndLose_9(order, openCode)
	case 10: //前四组选组选24
		o.WinningAndLose_10(order, openCode)
	case 11: //前四组选组选12
		o.WinningAndLose_11(order, openCode)
	case 12: //前四组选组选6
		o.WinningAndLose_12(order, openCode)
	case 13: //前四组选组选4
		o.WinningAndLose_13(order, openCode)
	case 14: //后四直选复式
		o.WinningAndLose_14(order, openCode)
	case 15: //后四直选组合
		o.WinningAndLose_15(order, openCode)
	case 16: //后四组选组选24
		o.WinningAndLose_16(order, openCode)
	case 17: //后四组选组选12
		o.WinningAndLose_17(order, openCode)
	case 18: //后四组选组选6
		o.WinningAndLose_18(order, openCode)
	case 19: //后四组选组选4
		o.WinningAndLose_19(order, openCode)
	case 20: //前三直选复式
		o.WinningAndLose_20(order, openCode)
	case 21: //前三直选和值
		o.WinningAndLose_21(order, openCode)
	case 22: //前三组选组三
		o.WinningAndLose_22(order, openCode)
	case 23: //前三组选组六
		o.WinningAndLose_23(order, openCode)
	case 24: //中三直选复式
		o.WinningAndLose_24(order, openCode)
	case 25: //中三直选和值
		o.WinningAndLose_25(order, openCode)
	case 26: //中三组选组三
		o.WinningAndLose_26(order, openCode)
	case 27: //中三组选组六
		o.WinningAndLose_27(order, openCode)
	case 28: //后三直选复式
		o.WinningAndLose_28(order, openCode)
	case 29: //后三直选和值
		o.WinningAndLose_29(order, openCode)
	case 30: //后三组选组三
		o.WinningAndLose_30(order, openCode)
	case 31: //后三组选组六
		o.WinningAndLose_31(order, openCode)
	case 32: //前二直选复式
		o.WinningAndLose_32(order, openCode)
	case 33: //前二直选和值
		o.WinningAndLose_33(order, openCode)
	case 34: //前二直选大小单双
		o.WinningAndLose_34(order, openCode)
	case 35: //前二组选复式
		o.WinningAndLose_35(order, openCode)
	case 36: //后二直选复式
		o.WinningAndLose_36(order, openCode)
	case 37: //后二直选和值
		o.WinningAndLose_37(order, openCode)
	case 38: //后二直选大小单双
		o.WinningAndLose_38(order, openCode)
	case 39: //后二组选复式
		o.WinningAndLose_39(order, openCode)
	case 40: //定位胆定位胆定位胆
		o.WinningAndLose_40(order, openCode)
	case 41: //三星不定胆一码前三
		o.WinningAndLose_41(order, openCode)
	case 42: //三星不定胆一码中三
		o.WinningAndLose_42(order, openCode)
	case 43: //三星不定胆一码后三
		o.WinningAndLose_43(order, openCode)
	case 44: //三星不定胆二码前三
		o.WinningAndLose_44(order, openCode)
	case 45: //三星不定胆二码中三
		o.WinningAndLose_45(order, openCode)
	case 46: //三星不定胆二码后三
		o.WinningAndLose_46(order, openCode)
	case 47: //任选任选二复选
		o.WinningAndLose_47(order, openCode)
	case 48: //任选任选二组选
		o.WinningAndLose_48(order, openCode)
	case 49: //任选任选三复选
		o.WinningAndLose_49(order, openCode)
	case 50: //任选任选三组三
		o.WinningAndLose_50(order, openCode)
	case 51: //任选任选三组六
		o.WinningAndLose_51(order, openCode)
	case 52: //任选任选四复选
		o.WinningAndLose_52(order, openCode)
	case 53: //跨度三星跨度前三
		o.WinningAndLose_53(order, openCode)
	case 54: //跨度三星跨度中三
		o.WinningAndLose_54(order, openCode)
	case 55: //跨度三星跨度后三
		o.WinningAndLose_55(order, openCode)
	case 56: //跨度二星跨度前二
		o.WinningAndLose_56(order, openCode)
	case 57: //跨度二星跨度后二
		o.WinningAndLose_57(order, openCode)
	case 58: //趣味趣味一帆风顺
		o.WinningAndLose_58(order, openCode)
	case 59: //趣味趣味好事成双
		o.WinningAndLose_59(order, openCode)
	case 60: //趣味趣味三星报喜
		o.WinningAndLose_60(order, openCode)
	case 61: //趣味趣味四季发财
		o.WinningAndLose_61(order, openCode)
	case 62: //龙虎龙虎万千
		o.WinningAndLose_62(order, openCode)
	case 63: //龙虎龙虎万百
		o.WinningAndLose_63(order, openCode)
	case 64: //龙虎龙虎万十
		o.WinningAndLose_64(order, openCode)
	case 65: //龙虎龙虎万个
		o.WinningAndLose_65(order, openCode)
	case 66: //龙虎龙虎千百
		o.WinningAndLose_66(order, openCode)
	case 67: //龙虎龙虎千十
		o.WinningAndLose_67(order, openCode)
	case 68: //龙虎龙虎千个
		o.WinningAndLose_68(order, openCode)
	case 69: //龙虎龙虎百十
		o.WinningAndLose_69(order, openCode)
	case 70: //龙虎龙虎百个
		o.WinningAndLose_70(order, openCode)
	case 71: //龙虎龙虎十个
		o.WinningAndLose_71(order, openCode)
	case 72: //龙虎龙虎十个
		o.WinningAndLose_72(order, openCode)
	case 73: //龙虎龙虎十个
		o.WinningAndLose_73(order, openCode)
	case 74: //龙虎龙虎十个
		o.WinningAndLose_74(order, openCode)
	case 75: //龙虎龙虎十个
		o.WinningAndLose_75(order, openCode)
	case 76: //龙虎龙虎十个
		o.WinningAndLose_76(order, openCode)
	case 77: //龙虎龙虎十个
		o.WinningAndLose_77(order, openCode)
	case 78: //龙虎龙虎十个
		o.WinningAndLose_78(order, openCode)
	default:
		beego.Debug("失败")
		return false
	}

	return true
}

//判断是否中奖 (0 五星直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_0 = 0
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0
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
				flag_1 = 1
				break
			}
		}
		if flag_1 == 1 {
			//第三名
			for _, v := range betNumbers[2] {
				if v == openCode[2] {
					flag_2 = 1
					break
				}
			}
			if flag_2 == 1 {
				//第四名
				for _, v := range betNumbers[3] {
					if v == openCode[3] {
						flag_3 = 1
						break
					}
				}
				if flag_3 == 1 {
					//第五名
					for _, v := range betNumbers[4] {
						if v == openCode[4] {
							winningNumbers = 1
							break
						}
					}
				}
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

//判断是否中奖 (1 五星直选组合)(返回用户输赢情况)
func (o *SSC) WinningAndLose_1(order *gb.Order, openCode []int) {
	ok, betNumbers := o.PaserFiveDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0
	var flag_4 = 0

	//每位个数
	l1 := len(betNumbers[0])
	l2 := len(betNumbers[1])
	l3 := len(betNumbers[2])
	l4 := len(betNumbers[3])

	//中奖注数
	var winningBetNum5 = 0
	var winningBetNum4 = 0
	var winningBetNum3 = 0
	var winningBetNum2 = 0
	var winningBetNum1 = 0

	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[4] {
			flag_4 = 1
			winningBetNum5 = l1 * l2 * l3 * l4
			break
		}
	}
	if flag_4 == 1 {
		//第四名
		for _, v := range betNumbers[3] {
			if v == openCode[3] {
				flag_3 = 1
				winningBetNum4 = l1 * l2 * l3
				break
			}
		}
		if flag_3 == 1 {
			//第三名
			for _, v := range betNumbers[2] {
				if v == openCode[2] {
					flag_2 = 1
					winningBetNum3 = l1 * l2
					break
				}
			}
			if flag_2 == 1 {
				//第二名
				for _, v := range betNumbers[1] {
					if v == openCode[1] {
						flag_1 = 1
						winningBetNum2 = l1
						break
					}
				}
				if flag_1 == 1 {
					//第二名
					for _, v := range betNumbers[0] {
						if v == openCode[0] {
							winningBetNum1 = 1
							break
						}
					}
				}
			}
		}
	}

	//最后结果 输赢多少钱
	var ret float64
	//总中奖注数
	var winningBetNum = 0
	//第5位中奖
	if winningBetNum5 != 0 {
		//中的注数
		winningBetNum += winningBetNum5
		//计算中了多少钱
		if v, ok := order.Odds["5"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum5)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	//第4位中奖
	if winningBetNum4 != 0 {
		//中的注数
		winningBetNum += winningBetNum4
		//计算中了多少钱
		if v, ok := order.Odds["4"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum4)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第3位中奖
	if winningBetNum3 != 0 {
		//中的注数
		winningBetNum += winningBetNum3
		//计算中了多少钱
		if v, ok := order.Odds["3"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum3)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第2位中奖
	if winningBetNum2 != 0 {
		//中的注数
		winningBetNum += winningBetNum2
		//计算中了多少钱
		if v, ok := order.Odds["2"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum2)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第1位中奖
	if winningBetNum1 != 0 {
		//中的注数
		winningBetNum += winningBetNum1
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	//计算反水
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//最后结算金额(中奖金额 + 反水)
	order.Settlement = ret + order.RebateAmount
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (2 五星组选组选120)(返回用户输赢情况)
func (o *SSC) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0

	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 5 证明没有中奖
	if winningNumbers < 5 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 5 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (3 五星组选组选60)(返回用户输赢情况)
func (o *SSC) WinningAndLose_3(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//分析开奖结果,找出能中奖的号码(号码和号码个数)
	var NumAndNumOf = make(map[int]int)
	for _, v := range openCode {
		if m, ok := NumAndNumOf[v]; ok {
			NumAndNumOf[v] = m + 1
		} else {
			NumAndNumOf[v] = 1
		}
	}

	mapl := len(NumAndNumOf)

	var winningNumbers = 0

	//只有开出个数字的情况下才能中奖
	if mapl == 4 {
		//找出重号位和单号位
		var sameNum = 0
		var singleNum []int
		for k, v := range NumAndNumOf {
			if v == 2 {
				sameNum = k
			} else if v == 1 {
				singleNum = append(singleNum, k)
			}
		}
		//开始分别匹配 重号位,和单号位,是否中奖
		for _, v := range betNumbers[0] {
			if v == sameNum {
				for _, i := range betNumbers[1] {
					for _, j := range singleNum {
						if i == j {
							winningNumbers++
							break
						}
					}
				}
				break
			}
		}

	}

	var winningBetNum = 0
	//如果中奖数组< 3 证明没有中奖
	if winningNumbers < 3 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 3 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (4 五星组选组选30)(返回用户输赢情况)
func (o *SSC) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//分析开奖结果,找出能中奖的号码(号码和号码个数)
	var NumAndNumOf = make(map[int]int)
	for _, v := range openCode {
		if m, ok := NumAndNumOf[v]; ok {
			NumAndNumOf[v] = m + 1
		} else {
			NumAndNumOf[v] = 1
		}
	}

	mapl := len(NumAndNumOf)

	var winningNumbers = 0

	//只有开出3个数字的情况下才能中奖
	if mapl == 3 {
		//找出重号位和单号位
		var sameNum []int
		var singleNum int = 0
		for k, v := range NumAndNumOf {
			if v == 2 {
				sameNum = append(sameNum, k)
			} else if v == 1 {
				singleNum = k
			}
		}
		//开始分别匹配 重号位,和单号位,是否中奖(这个比较特殊,先匹配单号位,再来计算同号位的出现次数)
		for _, v := range betNumbers[1] {
			if v == singleNum {
				for _, i := range betNumbers[0] {
					for _, j := range sameNum {
						if i == j {
							winningNumbers++
							break
						}
					}
				}
				break
			}
		}

	}

	var winningBetNum = 0
	//如果中奖数组< 1 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 2 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (5 五星组选组选20)(返回用户输赢情况)
func (o *SSC) WinningAndLose_5(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//分析开奖结果,找出能中奖的号码(号码和号码个数)
	var NumAndNumOf = make(map[int]int)
	for _, v := range openCode {
		if m, ok := NumAndNumOf[v]; ok {
			NumAndNumOf[v] = m + 1
		} else {
			NumAndNumOf[v] = 1
		}
	}

	mapl := len(NumAndNumOf)

	var winningNumbers = 0

	//只有开出个数字的情况下才能中奖
	if mapl == 3 {
		//找出重号位和单号位
		var sameNum = 0
		var singleNum []int
		for k, v := range NumAndNumOf {
			if v == 3 {
				sameNum = k
			} else if v == 1 {
				singleNum = append(singleNum, k)
			}
		}

		//开始分别匹配 重号位,和单号位,是否中奖
		for _, v := range betNumbers[0] {
			if v == sameNum {
				for _, i := range betNumbers[1] {
					for _, j := range singleNum {
						if i == j {
							winningNumbers++
							break
						}
					}
				}
				break
			}
		}
	}
	var winningBetNum = 0
	//如果中奖数组< 2 证明没有中奖
	if winningNumbers < 2 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 2 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (6 五星组选组选10)(返回用户输赢情况)
func (o *SSC) WinningAndLose_6(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//分析开奖结果,找出能中奖的号码(号码和号码个数)
	var NumAndNumOf = make(map[int]int)
	for _, v := range openCode {
		if m, ok := NumAndNumOf[v]; ok {
			NumAndNumOf[v] = m + 1
		} else {
			NumAndNumOf[v] = 1
		}
	}

	mapl := len(NumAndNumOf)

	var winningNumbers = 0

	//只有开出2个数字的情况下才能中奖
	if mapl == 2 {
		//找出重号位和双号位
		var sameNum = 0
		var singleNum = 0
		for k, v := range NumAndNumOf {
			if v == 3 {
				sameNum = k
			} else if v == 2 {
				singleNum = k
			}
		}
		//开始分别匹配 重号位,和单号位,是否中奖
		for _, v := range betNumbers[0] {
			if v == sameNum {
				for _, i := range betNumbers[1] {
					if i == singleNum {
						winningNumbers++
						break
					}
				}
				break
			}
		}

	}

	var winningBetNum = 0
	//如果中奖数组< 1 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (7 五星组选组选5)(返回用户输赢情况)
func (o *SSC) WinningAndLose_7(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//分析开奖结果,找出能中奖的号码(号码和号码个数)
	var NumAndNumOf = make(map[int]int)
	for _, v := range openCode {
		if m, ok := NumAndNumOf[v]; ok {
			NumAndNumOf[v] = m + 1
		} else {
			NumAndNumOf[v] = 1
		}
	}

	mapl := len(NumAndNumOf)

	var winningNumbers = 0

	//只有开出2个数字的情况下才能中奖
	if mapl == 2 {
		//找出重号位和双号位
		var sameNum = 0
		var singleNum = 0
		for k, v := range NumAndNumOf {
			if v == 4 {
				sameNum = k
			} else if v == 1 {
				singleNum = k
			}
		}
		//开始分别匹配 重号位,和单号位,是否中奖
		for _, v := range betNumbers[0] {
			if v == sameNum {
				for _, i := range betNumbers[1] {
					if i == singleNum {
						winningNumbers++
						break
					}
				}
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 1 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (8 前四直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_8(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFourDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_0 = 0
	var flag_1 = 0
	var flag_2 = 0

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
				flag_1 = 1
				break
			}
		}
		if flag_1 == 1 {
			//第三名
			for _, v := range betNumbers[2] {
				if v == openCode[2] {
					flag_2 = 1
					break
				}
			}
			if flag_2 == 1 {
				//第四名
				for _, v := range betNumbers[3] {
					if v == openCode[3] {
						winningNumbers = 1
						break
					}
				}
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

//判断是否中奖 (9 前四直选组合)(返回用户输赢情况)
func (o *SSC) WinningAndLose_9(order *gb.Order, openCode []int) {
	ok, betNumbers := o.PaserFourDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0

	//每位个数
	l1 := len(betNumbers[0])
	l2 := len(betNumbers[1])
	l3 := len(betNumbers[2])

	//中奖注数
	var winningBetNum4 = 0
	var winningBetNum3 = 0
	var winningBetNum2 = 0
	var winningBetNum1 = 0

	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[3] {
			flag_3 = 1
			winningBetNum4 = l1 * l2 * l3
			break
		}
	}
	if flag_3 == 1 {
		//第三名
		for _, v := range betNumbers[2] {
			if v == openCode[2] {
				flag_2 = 1
				winningBetNum3 = l1 * l2
				break
			}
		}
		if flag_2 == 1 {
			//第二名
			for _, v := range betNumbers[1] {
				if v == openCode[1] {
					flag_1 = 1
					winningBetNum2 = l1
					break
				}
			}
			if flag_1 == 1 {
				//第一名
				for _, v := range betNumbers[0] {
					if v == openCode[0] {
						winningBetNum1 = 1
						break
					}
				}
			}
		}
	}

	//最后结果 输赢多少钱
	var ret float64
	//总中奖注数
	var winningBetNum = 0
	//第4位中奖
	if winningBetNum4 != 0 {
		//中的注数
		winningBetNum += winningBetNum4
		//计算中了多少钱
		if v, ok := order.Odds["4"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum4)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第3位中奖
	if winningBetNum3 != 0 {
		//中的注数
		winningBetNum += winningBetNum3
		//计算中了多少钱
		if v, ok := order.Odds["3"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum3)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第2位中奖
	if winningBetNum2 != 0 {
		//中的注数
		winningBetNum += winningBetNum2
		//计算中了多少钱
		if v, ok := order.Odds["2"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum2)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第1位中奖
	if winningBetNum1 != 0 {
		//中的注数
		winningBetNum += winningBetNum1
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	//计算反水
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//最后结算金额(中奖金额 + 反水)
	order.Settlement = ret + order.RebateAmount
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (10 前四组选组选24)(返回用户输赢情况)
func (o *SSC) WinningAndLose_10(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0

	//截取前4个开奖号码 注意(截取使用的是下标,但是 [:4] 不包含下标4的数字)
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3])
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 4 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 4 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (11 前四组选组选12)(返回用户输赢情况)
func (o *SSC) WinningAndLose_11(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//截取前4个开奖号码 注意(截取使用的是下标,但是 [:4] 不包含下标4的数字)
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3])

	//是否有同号出现
	var isSameNum = false
	//同号点数(由于时时彩开奖号码为0 - 9)所以这里给10作为初始化
	var sameNum int = 10
	//不同号数组 len == 2的时候才有中奖可能
	var differentNum []int
	//是否有中奖可能,也就是说有没有开出可中奖的结果
	var isHaveChanceToWin bool = false

	//-----------------------判断是否有中奖可能 ,并分析出可中奖号码 --------------------
	//1,找出是否有同号位
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			if tempOpenCode[i] == tempOpenCode[j] {
				sameNum = tempOpenCode[i]
				isSameNum = true
				break
			}
		}
		if isSameNum == true {
			for _, i := range tempOpenCode {
				if i != sameNum {
					differentNum = append(differentNum, i)
				}
			}
			//判断不同号位的个数 和 是否相同 得出是否可能中奖的结果
			if len(differentNum) == 2 && differentNum[0] != differentNum[1] {
				isHaveChanceToWin = true
			}
			break
		}
	}

	//------------------------- 判断是否中奖 ------------------------
	var diffNumWinNum = 0
	//判断是否中奖
	if isHaveChanceToWin == true {
		for _, v := range betNumbers[0] {
			if v == sameNum { //同号位中奖
				for _, i := range betNumbers[1] {
					for _, j := range differentNum {
						if i == j {
							diffNumWinNum++
						}
					}
				}
				break
			}
		}
	}

	//中奖注数
	var winningBetNum = 0
	//只有不同号位也买中两个号码才算中奖,其他所有情况都不中奖
	if diffNumWinNum == 2 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	} else {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	}

	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (12 前四组选组选6)(返回用户输赢情况)
func (o *SSC) WinningAndLose_12(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//两个同号位是否中奖
	var twoSameNum = 0
	//截取前4个开奖号码 注意(截取使用的是下标,但是 [:4] 不包含下标4的数字)
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3])
	//先判断两个同号为是否中奖
	for _, v := range betNumbers {
		var sameNum = 0
		for _, i := range tempOpenCode {
			if v == i {
				sameNum++
			}
		}
		//开出两个同号位号码
		if sameNum == 2 {
			twoSameNum++
		}
	}

	var winningNumbers = 0
	if twoSameNum == 2 {
		winningNumbers = 1
	}

	var winningBetNum = 0
	//如果中奖数组< 1 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (13 前四组选组选4)(返回用户输赢情况)
func (o *SSC) WinningAndLose_13(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//同号位是否中奖
	var isSameNum = false
	//截取前4个开奖号码 注意(截取使用的是下标,但是 [:4] 不包含下标4的数字)
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3])
	//先判断同号位是否开出
	for _, v := range betNumbers[0] {
		var sameNum = 0
		for _, i := range tempOpenCode {
			if v == i {
				sameNum++
			}
		}
		//开出三个同号位号码
		if sameNum == 3 {
			isSameNum = true
			break
		}
	}

	var winningNumbers = 0
	//如果同号位中奖,判断另外1个不同号位是否中奖
	if isSameNum {
		//判断下注号码中有几个开出
		for _, v := range betNumbers[1] {
			for _, i := range tempOpenCode {
				if v == i {
					winningNumbers++
					break
				}
			}
			if winningNumbers == 1 {
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 2 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (14 后四直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_14(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFourDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0

	var winningNumbers = 0
	var winningBetNum = 0
	//第二名
	for _, v := range betNumbers[0] {
		if v == openCode[1] {
			flag_1 = 1
			break
		}
	}

	if flag_1 == 1 {
		//第三名
		for _, v := range betNumbers[1] {
			if v == openCode[2] {
				flag_2 = 1
				break
			}
		}
		if flag_2 == 1 {
			//第四名
			for _, v := range betNumbers[2] {
				if v == openCode[3] {
					flag_3 = 1
					break
				}
			}
			if flag_3 == 1 {
				//第五名
				for _, v := range betNumbers[3] {
					if v == openCode[4] {
						winningNumbers = 1
						break
					}
				}
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

//判断是否中奖 (15 后四直选组合)(返回用户输赢情况)
func (o *SSC) WinningAndLose_15(order *gb.Order, openCode []int) {
	ok, betNumbers := o.PaserFourDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1], openCode[2], openCode[3], openCode[4])
	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0

	//每位个数
	l1 := len(betNumbers[0])
	l2 := len(betNumbers[1])
	l3 := len(betNumbers[2])

	//中奖注数
	var winningBetNum4 = 0
	var winningBetNum3 = 0
	var winningBetNum2 = 0
	var winningBetNum1 = 0

	//第四名
	for _, v := range betNumbers[3] {
		if v == tmpOpenCode[3] {
			flag_3 = 1
			winningBetNum4 = l1 * l2 * l3
			break
		}
	}
	if flag_3 == 1 {
		//第三名
		for _, v := range betNumbers[2] {
			if v == tmpOpenCode[2] {
				flag_2 = 1
				winningBetNum3 = l1 * l2
				break
			}
		}
		if flag_2 == 1 {
			//第二名
			for _, v := range betNumbers[1] {
				if v == tmpOpenCode[1] {
					flag_1 = 1
					winningBetNum2 = l1
					break
				}
			}
			if flag_1 == 1 {
				//第一名
				for _, v := range betNumbers[0] {
					if v == tmpOpenCode[0] {
						winningBetNum1 = 1
						break
					}
				}
			}
		}
	}

	//最后结果 输赢多少钱
	var ret float64
	//总中奖注数
	var winningBetNum = 0
	//第4位中奖
	if winningBetNum4 != 0 {
		//中的注数
		winningBetNum += winningBetNum4
		//计算中了多少钱
		if v, ok := order.Odds["4"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum4)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第3位中奖
	if winningBetNum3 != 0 {
		//中的注数
		winningBetNum += winningBetNum3
		//计算中了多少钱
		if v, ok := order.Odds["3"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum3)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第2位中奖
	if winningBetNum2 != 0 {
		//中的注数
		winningBetNum += winningBetNum2
		//计算中了多少钱
		if v, ok := order.Odds["2"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum2)
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}
	//第1位中奖
	if winningBetNum1 != 0 {
		//中的注数
		winningBetNum += winningBetNum1
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	//计算反水
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//最后结算金额(中奖金额 + 反水)
	order.Settlement = ret + order.RebateAmount
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (16 后四组选组选24)(返回用户输赢情况)
func (o *SSC) WinningAndLose_16(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0

	//截取后4个开奖号码 注意(截取使用的是下标, [1:] )
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[1], openCode[2], openCode[3], openCode[4])
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 4 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 4 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (17 后四组选组选12)(返回用户输赢情况)
func (o *SSC) WinningAndLose_17(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//截取前4个开奖号码 注意(截取使用的是下标,但是 [:4] 不包含下标4的数字)
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[1], openCode[2], openCode[3], openCode[4])

	//是否有同号出现
	var isSameNum = false
	//同号点数(由于时时彩开奖号码为0 - 9)所以这里给10作为初始化
	var sameNum int = 10
	//不同号数组 len == 2的时候才有中奖可能
	var differentNum []int
	//是否有中奖可能,也就是说有没有开出可中奖的结果
	var isHaveChanceToWin bool = false

	//-----------------------判断是否有中奖可能 ,并分析出可中奖号码 --------------------
	//1,找出是否有同号位
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			if tempOpenCode[i] == tempOpenCode[j] {
				sameNum = tempOpenCode[i]
				isSameNum = true
				break
			}
		}
		if isSameNum == true {
			for _, i := range tempOpenCode {
				if i != sameNum {
					differentNum = append(differentNum, i)
				}
			}
			//判断不同号位的个数 和 是否相同 得出是否可能中奖的结果
			if len(differentNum) == 2 && differentNum[0] != differentNum[1] {
				isHaveChanceToWin = true
			}
			break
		}
	}

	//------------------------- 判断是否中奖 ------------------------
	var diffNumWinNum = 0
	//判断是否中奖
	if isHaveChanceToWin == true {
		for _, v := range betNumbers[0] {
			if v == sameNum { //同号位中奖
				for _, i := range betNumbers[1] {
					for _, j := range differentNum {
						if i == j {
							diffNumWinNum++
						}
					}
				}
				break
			}
		}
	}

	//中奖注数
	var winningBetNum = 0
	//只有不同号位也买中两个号码才算中奖,其他所有情况都不中奖
	if diffNumWinNum == 2 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	} else {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	}

	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (18 前四组选组选6)(返回用户输赢情况)
func (o *SSC) WinningAndLose_18(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//两个同号位是否中奖
	var twoSameNum = 0
	//截取前4个开奖号码 注意(截取使用的是下标, [1:] )
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[1], openCode[2], openCode[3], openCode[4])
	//先判断两个同号为是否中奖
	for _, v := range betNumbers {
		var sameNum = 0
		for _, i := range tempOpenCode {
			if v == i {
				sameNum++
			}
		}
		//开出两个同号位号码
		if sameNum == 2 {
			twoSameNum++
		}
	}

	var winningNumbers = 0
	if twoSameNum == 2 {
		winningNumbers = 1
	}

	var winningBetNum = 0
	//如果中奖数组< 1 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (19 前四组选组选4)(返回用户输赢情况)
func (o *SSC) WinningAndLose_19(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//同号位是否中奖
	var isSameNum = false
	//截取前4个开奖号码 注意(截取使用的是下标, [1:] )
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[1], openCode[2], openCode[3], openCode[4])
	//先判断同号位是否开出
	for _, v := range betNumbers[0] {
		var sameNum = 0
		for _, i := range tempOpenCode {
			if v == i {
				sameNum++
			}
		}
		//开出三个同号位号码
		if sameNum == 3 {
			isSameNum = true
			break
		}
	}

	var winningNumbers = 0
	//如果同号位中奖,判断另外1个不同号位是否中奖
	if isSameNum {
		//判断下注号码中有几个开出
		for _, v := range betNumbers[1] {
			for _, i := range tempOpenCode {
				if v == i {
					winningNumbers++
					break
				}
			}
			if winningNumbers == 1 {
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 2 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (20 前三直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_20(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_0 = 0
	var flag_1 = 0

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
				flag_1 = 1
				break
			}
		}
		if flag_1 == 1 {
			//第三名
			for _, v := range betNumbers[2] {
				if v == openCode[2] {
					winningNumbers = 1
					break
				}
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

//判断是否中奖 (21 前三直选和值)(返回用户输赢情况)
func (o *SSC) WinningAndLose_21(order *gb.Order, openCode []int) {
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

//判断是否中奖 (22 前三组选组三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_22(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断下注号码中有几个开出(先看有没有开出两个相同的号码)
	var sameNumber = 0
	var winningBetNum = 0
	//前三位 开奖结果一样 不算中奖
	if openCode[0] == openCode[1] && openCode[0] == openCode[2] {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//更新order
		order.Settlement = ret
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

//判断是否中奖 (23 前三组选组六)(返回用户输赢情况)
func (o *SSC) WinningAndLose_23(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0], openCode[1], openCode[2])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
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

//判断是否中奖 (24 中三直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_24(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_1 = 0
	var flag_2 = 0

	var winningNumbers = 0
	var winningBetNum = 0
	//第二名
	for _, v := range betNumbers[0] {
		if v == openCode[1] {
			flag_1 = 1
			break
		}
	}

	if flag_1 == 1 {
		//第三名
		for _, v := range betNumbers[1] {
			if v == openCode[2] {
				flag_2 = 1
				break
			}
		}
		if flag_2 == 1 {
			//第四名
			for _, v := range betNumbers[2] {
				if v == openCode[3] {
					winningNumbers = 1
					break
				}
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

//判断是否中奖 (25 中三直选和值)(返回用户输赢情况)
func (o *SSC) WinningAndLose_25(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForThreeSum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算中三个开奖号码的和值
	sum := openCode[1] + openCode[2] + openCode[3]
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

//判断是否中奖 (26 中三组选组三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_26(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断下注号码中有几个开出(先看有没有开出两个相同的号码)
	var sameNumber = 0
	var winningBetNum = 0
	//中三位 开奖结果一样 不算中奖
	if openCode[1] == openCode[2] && openCode[1] == openCode[3] {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//更新order
		order.Settlement = ret
	} else if openCode[1] == openCode[2] {
		for _, v := range betNumbers {
			if openCode[1] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers {
				if v == openCode[3] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[1] == openCode[3] {
		for _, v := range betNumbers {
			if openCode[1] == v {
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
	} else if openCode[2] == openCode[3] {
		for _, v := range betNumbers {
			if openCode[2] == v {
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

//判断是否中奖 (27 中三组选组六)(返回用户输赢情况)
func (o *SSC) WinningAndLose_27(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1], openCode[2], openCode[3])
	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
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

//判断是否中奖 (28 后三直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_28(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var flag_2 = 0
	var flag_3 = 0

	var winningNumbers = 0
	var winningBetNum = 0
	//第三名
	for _, v := range betNumbers[0] {
		if v == openCode[2] {
			flag_2 = 1
			break
		}
	}

	if flag_2 == 1 {
		//第四名
		for _, v := range betNumbers[1] {
			if v == openCode[3] {
				flag_3 = 1
				break
			}
		}
		if flag_3 == 1 {
			//第五名
			for _, v := range betNumbers[2] {
				if v == openCode[4] {
					winningNumbers = 1
					break
				}
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

//判断是否中奖 (29 后三直选和值)(返回用户输赢情况)
func (o *SSC) WinningAndLose_29(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForThreeSum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算中三个开奖号码的和值
	sum := openCode[2] + openCode[3] + openCode[4]
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

//判断是否中奖 (30 后三组选组三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_30(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断下注号码中有几个开出(先看有没有开出两个相同的号码)
	var sameNumber = 0
	var winningBetNum = 0
	//中三位 开奖结果一样 不算中奖
	if openCode[2] == openCode[3] && openCode[2] == openCode[4] {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//更新order
		order.Settlement = ret
	} else if openCode[2] == openCode[3] {
		for _, v := range betNumbers {
			if openCode[2] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers {
				if v == openCode[4] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[2] == openCode[4] {
		for _, v := range betNumbers {
			if openCode[2] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers {
				if v == openCode[3] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[3] == openCode[4] {
		for _, v := range betNumbers {
			if openCode[3] == v {
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

//判断是否中奖 (31 后三组选组六)(返回用户输赢情况)
func (o *SSC) WinningAndLose_31(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[2], openCode[3], openCode[4])
	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
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

//判断是否中奖 (32 前二直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_32(order *gb.Order, openCode []int) {
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

//判断是否中奖 (33 前二直选和值)(返回用户输赢情况)
func (o *SSC) WinningAndLose_33(order *gb.Order, openCode []int) {
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

//判断是否中奖 (34 前二直选大小单双)(返回用户输赢情况)
func (o *SSC) WinningAndLose_34(order *gb.Order, openCode []int) {
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

//判断是否中奖 (35 前二组选复试)(返回用户输赢情况)
func (o *SSC) WinningAndLose_35(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0], openCode[1])
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

//判断是否中奖 (36 后二直选复式)(返回用户输赢情况)
func (o *SSC) WinningAndLose_36(order *gb.Order, openCode []int) {
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
	//第四名
	for _, v := range betNumbers[0] {
		if v == openCode[3] {
			flag_3 = 1
			break
		}
	}

	if flag_3 == 1 {
		//第五名
		for _, v := range betNumbers[1] {
			if v == openCode[4] {
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

//判断是否中奖 (37 后二直选和值)(返回用户输赢情况)
func (o *SSC) WinningAndLose_37(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForTwoSum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算后二个开奖号码的和值
	sum := openCode[3] + openCode[4]
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

//判断是否中奖 (38 后二直选大小单双)(返回用户输赢情况)
func (o *SSC) WinningAndLose_38(order *gb.Order, openCode []int) {
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

	if openCode[3] < 5 {
		B1 = 1
	}
	if openCode[3]%2 == 0 {
		D1 = 3
	}

	if openCode[4] < 5 {
		B2 = 1
	}
	if openCode[4]%2 == 0 {
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

//判断是否中奖 (39 后二组选复试)(返回用户输赢情况)
func (o *SSC) WinningAndLose_39(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[3], openCode[4])
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

//判断是否中奖 (40 定位胆定位胆定位胆)(返回用户输赢情况)
func (o *SSC) WinningAndLose_40(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers = 0

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
	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[3] {
			winningNumbers++
			break
		}
	}
	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[4] {
			winningNumbers++
			break
		}
	}
	var winningBetNum = 0
	//中几个号码就是几注
	winningBetNum = winningNumbers
	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 6 {

		//定位胆 1-5 最多只有5注中奖
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret

		//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
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

//判断是否中奖 (41 三星不定胆一码前三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_41(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var winningBetNum = 0
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0], openCode[1], openCode[2])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
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

//判断是否中奖 (42 三星不定胆一码中三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_42(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var winningBetNum = 0
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1], openCode[2], openCode[3])
	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
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

//判断是否中奖 (43 三星不定胆一码后三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_43(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中奖注数(注意:每一位只会有一个数字中奖)
	var winningBetNum = 0
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[2], openCode[3], openCode[4])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
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

//判断是否中奖 (44 三星不定胆二码前三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_44(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0], openCode[1], openCode[2])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	winningBetNum = utils.AnalysisCombination(winningNumbers, 2)

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

//判断是否中奖 (45 三星不定胆二码中三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_45(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1], openCode[2], openCode[3])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	winningBetNum = utils.AnalysisCombination(winningNumbers, 2)

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

//判断是否中奖 (46 三星不定胆二码后三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_46(order *gb.Order, openCode []int) {
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
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[2], openCode[3], openCode[4])

	for _, v := range betNumbers {
		for _, i := range tmpOpenCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	winningBetNum = utils.AnalysisCombination(winningNumbers, 2)

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

//判断是否中奖 (47 任选任选二复选)(返回用户输赢情况)
func (o *SSC) WinningAndLose_47(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0
	var winningBetNum = 0
	//判断中奖注数先判断 有几位中奖,然后用中奖位数进行2组合(注意:每一位只会有一个数字中奖)
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
	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[3] {
			winningNumbers++
			break
		}
	}
	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[4] {
			winningNumbers++
			break
		}
	}

	winningBetNum = utils.AnalysisCombination(winningNumbers, 2)

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 11 {
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret

		//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
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

//判断是否中奖 (48 任选任选二组选)(返回用户输赢情况)
func (o *SSC) WinningAndLose_48(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0

	//判断中奖注数先判断 有几位中奖,然后用中奖位数进行2组合(注意:每一位只会有一个数字中奖)
	//外层循环为每位下注数,再判断有下注位数上的下注号码有没有中奖
	for _, v := range betNumbers[1] {
		for _, i := range betNumbers[0] {
			if v == openCode[i] {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	winningBetNum = utils.AnalysisCombination(winningNumbers, 2)

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 10 {
		//任选二组选  最多只有10注中奖
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret

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

//判断是否中奖 (49 任选任选三复选)(返回用户输赢情况)
func (o *SSC) WinningAndLose_49(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0
	var winningBetNum = 0
	//判断中奖注数先判断 有几位中奖,然后用中奖位数进行3组合(注意:每一位只会有一个数字中奖)
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
	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[3] {
			winningNumbers++
			break
		}
	}
	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[4] {
			winningNumbers++
			break
		}
	}

	winningBetNum = utils.AnalysisCombination(winningNumbers, 3)

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 11 {
		//任选二复选  最多只有5注中奖
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret

		//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
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

//判断是否中奖 (50 任选任选三组三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_50(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//首先看看在哪几位有下注,生成新的开奖号码组
	var tmpOpenCode []int
	for _, v := range betNumbers[0] {
		tmpOpenCode = append(tmpOpenCode, openCode[v])
	}

	var recordNum = make(map[int]int)
	//用下注号码匹配开奖号码,看看买中几个位, 组成中奖号码组,并且记录号码出现次数
	var winningBetNumbers []int
	for _, v := range tmpOpenCode {
		for _, i := range betNumbers[1] {
			if v == i {
				winningBetNumbers = append(winningBetNumbers, v)
				//记录出现号码,和出现次数
				if m, ok := recordNum[v]; ok {
					recordNum[v] = m + 1
				} else {
					recordNum[v] = 1
				}
				break
			}
		}
	}

	//总共买中位数
	l := len(winningBetNumbers)
	//判断重复数组数

	var winningBetNum = 0
	//判断买中几位
	if l == 3 { //如果买中3位,直接判断有没有重号出现有就中一注
		for _, v := range recordNum {
			if v == 2 { //如果有一个号重复两次就中一注
				winningBetNum += 1
			}
		}
	} else if l == 4 { //买中4位和5位都
		for _, v := range recordNum {
			if v == 4 { //如果4个号重复,直接不中奖
				break
			} else if v == 3 { //如果有3个重就中4注
				winningBetNum += 3
			} else if v == 2 { //两个重号
				winningBetNum += 2
			}
		}
	} else if l == 5 {
		for _, v := range recordNum {
			if v == 5 { //如果5个号重复,直接不中奖
				break
			} else if v == 4 { //4个号重复
				winningBetNum += 6
			} else if v == 3 { //3个号重复
				winningBetNum += 6
			} else if v == 2 {
				winningBetNum += 3
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
	} else if winningBetNum > 0 {
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

//判断是否中奖 (51 任选任选三组六)(返回用户输赢情况)
func (o *SSC) WinningAndLose_51(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0

	//获取买中号码,重复数以及重复数个数
	var sameNum = make(map[int]int)

	for _, v := range betNumbers[0] {
		for _, i := range betNumbers[1] {
			if i == openCode[v] {
				winningNumbers++
				if m, ok := sameNum[i]; ok {
					sameNum[i] = m + 1
				} else {
					sameNum[i] = 1
				}
				break
			}
		}
	}

	var winningBetNum = 0
	mapl := len(sameNum)
	if winningNumbers == 5 {
		if mapl == 5 { //abcde
			winningBetNum = 10
		} else if mapl == 4 { //aabcd
			winningBetNum = 7
		} else if mapl == 3 {
			for _, v := range sameNum {
				if v == 3 { //aaabc
					winningBetNum = 3
				} else { //aabbc
					winningBetNum = 4
				}
			}
		}
	} else if winningNumbers == 4 {
		if mapl == 4 { //abcd
			winningBetNum = 4
		} else if mapl == 3 { //aabc
			winningBetNum = 2
		}
	} else if winningNumbers == 3 {
		if mapl == 3 { //abcd
			winningBetNum = 1
		}
	}

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 11 {
		//任选三组六  最多只有10注中奖
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret

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

//判断是否中奖 (52 任选任选四复选)(返回用户输赢情况)
func (o *SSC) WinningAndLose_52(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0
	var winningBetNum = 0
	//判断中奖注数先判断 有几位中奖,然后用中奖位数进行4组合(注意:每一位只会有一个数字中奖)
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
	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[3] {
			winningNumbers++
			break
		}
	}
	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[4] {
			winningNumbers++
			break
		}
	}

	winningBetNum = utils.AnalysisCombination(winningNumbers, 4)

	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 6 {
		//任选四复选  最多只有5注中奖
		//先计算未中奖的反水
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret

		//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
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

//判断是否中奖 (53 跨度三星跨度前三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_53(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出前三位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0], openCode[1], openCode[2])

	//从小到大排序
	sort.Sort(utils.IntSlice(tmpOpenCode))

	//计算跨度
	span := tmpOpenCode[2] - tmpOpenCode[0]
	//循环下注号码是下中这个跨度
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == span {
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

//判断是否中奖 (54 跨度三星跨度中三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_54(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出前三位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1], openCode[2], openCode[3])
	//从小到大排序
	sort.Sort(utils.IntSlice(tmpOpenCode))

	//计算跨度
	span := tmpOpenCode[2] - tmpOpenCode[0]

	//循环下注号码是下中这个跨度
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == span {
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

//判断是否中奖 (55 跨度三星跨度后三)(返回用户输赢情况)
func (o *SSC) WinningAndLose_55(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出后三位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[2], openCode[3], openCode[4])
	//从小到大排序
	sort.Sort(utils.IntSlice(tmpOpenCode))
	//计算跨度
	span := tmpOpenCode[2] - tmpOpenCode[0]

	//循环下注号码是下中这个跨度
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == span {
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

//判断是否中奖 (56 跨度二星跨度前二)(返回用户输赢情况)
func (o *SSC) WinningAndLose_56(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//beego.Debug("SSC 56 原开奖号码:  ", openCode)

	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0], openCode[1])
	//从小到大排序
	sort.Sort(utils.IntSlice(tmpOpenCode))

	//计算跨度
	span := tmpOpenCode[1] - tmpOpenCode[0]
	//循环下注号码是下中这个跨度
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == span {
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

//判断是否中奖 (57 跨度二星跨度后二)(返回用户输赢情况)
func (o *SSC) WinningAndLose_57(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出后两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[3], openCode[4])
	//从小到大排序
	sort.Sort(utils.IntSlice(tmpOpenCode))

	//计算跨度
	span := tmpOpenCode[1] - tmpOpenCode[0]
	//循环下注号码是下中这个跨度
	var winningBetNum = 0
	for _, v := range betNumbers {
		if v == span {
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

//判断是否中奖 (58 趣味趣味一帆风顺)(返回用户输赢情况)
func (o *SSC) WinningAndLose_58(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningBetNum++
				break
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
	} else {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -中奖注数
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

//判断是否中奖 (59 趣味趣味好事成双)(返回用户输赢情况)
func (o *SSC) WinningAndLose_59(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var recordNum = make(map[int]int)
	//用下注号码匹配开奖号码,看看买中几个位, 组成中奖号码组,并且记录号码出现次数
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				//记录出现号码,和出现次数
				if m, ok := recordNum[i]; ok {
					recordNum[i] = m + 1
				} else {
					recordNum[i] = 1
				}
			}
		}
	}

	//判买出了几个重复2次以上的数
	var winningBetNum = 0
	for _, v := range recordNum {
		if v >= 2 {
			winningBetNum++
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
	} else {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -中奖注数
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

//判断是否中奖 (60 趣味趣味三星报喜)(返回用户输赢情况)
func (o *SSC) WinningAndLose_60(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var recordNum = make(map[int]int)
	//用下注号码匹配开奖号码,看看买中几个位, 组成中奖号码组,并且记录号码出现次数
	for _, v := range openCode {
		for _, i := range betNumbers {
			if v == i {
				//记录出现号码,和出现次数
				if m, ok := recordNum[v]; ok {
					recordNum[v] = m + 1
				} else {
					recordNum[v] = 1
				}
				break
			}
		}
	}
	//判买出了几个重复3次以上的数
	var winningBetNum = 0
	for _, v := range recordNum {
		if v >= 3 {
			winningBetNum++
			break
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
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -中奖注数
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

//判断是否中奖 (61 趣味趣味四季发财)(返回用户输赢情况)
func (o *SSC) WinningAndLose_61(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var recordNum = make(map[int]int)
	//用下注号码匹配开奖号码,看看买中几个位, 组成中奖号码组,并且记录号码出现次数
	for _, v := range openCode {
		for _, i := range betNumbers {
			if v == i {
				//记录出现号码,和出现次数
				if m, ok := recordNum[v]; ok {
					recordNum[v] = m + 1
				} else {
					recordNum[v] = 1
				}
				break
			}
		}
	}
	//判买出了几个重复3次以上的数
	var winningBetNum = 0
	for _, v := range recordNum {
		if v >= 4 {
			winningBetNum++
			break
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
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -中奖注数
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

//判断是否中奖 (62 龙虎龙虎万千)(返回用户输赢情况)
func (o *SSC) WinningAndLose_62(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0])
	tmpOpenCode = append(tmpOpenCode, openCode[1])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (63 龙虎龙虎万百)(返回用户输赢情况)
func (o *SSC) WinningAndLose_63(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0])
	tmpOpenCode = append(tmpOpenCode, openCode[2])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (64 龙虎龙虎万十)(返回用户输赢情况)
func (o *SSC) WinningAndLose_64(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Error("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0])
	tmpOpenCode = append(tmpOpenCode, openCode[3])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (65 龙虎龙虎万个)(返回用户输赢情况)
func (o *SSC) WinningAndLose_65(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[0])
	tmpOpenCode = append(tmpOpenCode, openCode[4])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (66 龙虎龙虎千百)(返回用户输赢情况)
func (o *SSC) WinningAndLose_66(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1])
	tmpOpenCode = append(tmpOpenCode, openCode[2])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (67 龙虎龙虎千十)(返回用户输赢情况)
func (o *SSC) WinningAndLose_67(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1])
	tmpOpenCode = append(tmpOpenCode, openCode[3])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (68 龙虎龙虎千个)(返回用户输赢情况)
func (o *SSC) WinningAndLose_68(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[1])
	tmpOpenCode = append(tmpOpenCode, openCode[4])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (69 龙虎龙虎百十)(返回用户输赢情况)
func (o *SSC) WinningAndLose_69(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[2])
	tmpOpenCode = append(tmpOpenCode, openCode[3])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (70 龙虎龙虎百个)(返回用户输赢情况)
func (o *SSC) WinningAndLose_70(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[2])
	tmpOpenCode = append(tmpOpenCode, openCode[4])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (71 龙虎龙虎十个)(返回用户输赢情况)
func (o *SSC) WinningAndLose_71(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//找出前两位
	var tmpOpenCode []int
	tmpOpenCode = append(tmpOpenCode, openCode[3])
	tmpOpenCode = append(tmpOpenCode, openCode[4])

	var dragon = 0
	var tiger = 0
	var sum = 0

	if tmpOpenCode[0] > tmpOpenCode[1] {
		dragon = 1
	} else if tmpOpenCode[0] < tmpOpenCode[1] {
		tiger = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == 0 {
		if dragon == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 1 {
		if tiger == 1 {
			winningBetNum = 1
		}
	} else if betNumbers[0] == 2 {
		if sum == 1 {
			winningBetNum = 1
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
		if dragon == 1 {
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
		} else if tiger == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		} else if sum == 1 {
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["3"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}

	}

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 72 两面万位
func (o *SSC) WinningAndLose_72(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor72(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [6]int

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

	if o.openCode[0] > o.openCode[3] {
		resultCode[4] = 10
	} else if o.openCode[0] < o.openCode[3] {
		resultCode[4] = 11
	} else {
		resultCode[4] = 12
	}

	if o.openCode[0] > o.openCode[4] {
		resultCode[5] = 13
	} else if o.openCode[0] < o.openCode[4] {
		resultCode[5] = 14
	} else {
		resultCode[5] = 15
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

//判断是否中奖 73 两面千位
func (o *SSC) WinningAndLose_73(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor73(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [5]int

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

	if o.openCode[1] > o.openCode[3] {
		resultCode[3] = 7
	} else if o.openCode[1] < o.openCode[3] {
		resultCode[3] = 8
	} else {
		resultCode[3] = 9
	}

	if o.openCode[1] > o.openCode[4] {
		resultCode[4] = 10
	} else if o.openCode[1] < o.openCode[4] {
		resultCode[4] = 11
	} else {
		resultCode[4] = 12
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

//判断是否中奖 74 两面百位
func (o *SSC) WinningAndLose_74(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor74(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [4]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6, 龙百 7, 虎个 8, 百个和 9
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

	if o.openCode[2] > o.openCode[3] {
		resultCode[2] = 4
	} else if o.openCode[2] < o.openCode[3] {
		resultCode[2] = 5
	} else {
		resultCode[2] = 6
	}

	if o.openCode[2] > o.openCode[4] {
		resultCode[3] = 7
	} else if o.openCode[2] < o.openCode[4] {
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

//判断是否中奖 75 两面十位
func (o *SSC) WinningAndLose_75(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor75(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6, 龙百 7, 虎个 8, 百个和 9
	if o.openCode[3] > 4 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if o.openCode[3]%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	if o.openCode[3] > o.openCode[4] {
		resultCode[2] = 4
	} else if o.openCode[3] < o.openCode[4] {
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

//判断是否中奖 76 两面个位
func (o *SSC) WinningAndLose_76(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor76(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6, 龙百 7, 虎个 8, 百个和 9
	if o.openCode[4] > 4 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if o.openCode[4]%2 == 0 {
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

//判断是否中奖 77 两面总和
func (o *SSC) WinningAndLose_77(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor76(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	var winningBetNum = 0

	sum := o.openCode[0] + o.openCode[1] + o.openCode[2] + o.openCode[3] + o.openCode[4]
	//大 0 , 小 1 , 单 2, 双 3,
	if sum > 22 {
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

//判断是否中奖 78 牛牛
func (o *SSC) WinningAndLose_78(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor76(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	var winningBetNum = 0

	sum := o.openCode[0] + o.openCode[1] + o.openCode[2] + o.openCode[3] + o.openCode[4]
	//大 0 , 小 1 , 单 2, 双 3,
	if sum > 22 {
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
