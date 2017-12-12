package EX5

import (
	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/Order"
	"calculsrv/models/acmgr"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"
	"common/utils"
	"fmt"
	"sort"
	"strings"

	"strconv"

	"github.com/astaxie/beego"
)

//分析订单(下注)
func (o *EX5) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {

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
	case 0, 9: //任选一~任选四任选任选二 ,前二前二组选
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于2个元素，应为任选2 最多只能选11个数字 1 - 11 或最少选择两个数字
		if l > 11 || l < 2 {
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

	case 1, 11: //任选一~任选四任选任选三, 前三前三组选
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			beego.Debug("失败1")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于3个元素，应为任选3 最多只能选11个数字 1 - 11 或最少选择3个数字
		if l > 11 || l < 3 {
			beego.Debug("失败2")
			return false
		}

		//6分析订单注数 注意 任选3 3个数为一组
		singleBetNum := utils.AnalysisCombination(l, 3)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败3")
			return false
		}
		order.SingleBetNum = singleBetNum
	case 2: //任选一~任选四任选任选四
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于4个元素，应为任选4 最多只能选11个数字 1 - 11 或最少选择4个数字
		if l > 11 || l < 4 {
			//beego.Debug("失败")
			return false
		}
		//6分析订单注数 注意 任选4 4个数为一组
		singleBetNum := utils.AnalysisCombination(l, 4)
		order.SingleBetNum = singleBetNum

		//7计算订单金额
		order.OrderAmount = order.SingleBetAmount * float64(singleBetNum)

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}
	case 3: //任选五~任选八任选任选五
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于5个元素，应为任选5 最多只能选11个数字 1 - 11 或最少选择5个数字
		if l > 11 || l < 5 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 任选5 5个数为一组
		singleBetNum := utils.AnalysisCombination(l, 5)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}
	case 4: //任选五~任选八任选任选六
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于6个元素，应为任选6 最多只能选11个数字 1 - 11 或最少选择6个数字
		if l > 11 || l < 6 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 任选6 6个数为一组
		singleBetNum := utils.AnalysisCombination(l, 6)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}
	case 5: //任选五~任选八任选任选七
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于7个元素，应为任选7 最多只能选11个数字 1 - 11 或最少选择7个数字
		if l > 11 || l < 7 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 任选6 6个数为一组
		singleBetNum := utils.AnalysisCombination(l, 7)
		order.SingleBetNum = singleBetNum
		//判断订单总额有没有超过限制
		if order.SingleBetAmount*float64(singleBetNum) > o.Settings[order.BetType].OrderLimit {
			//beego.Debug("失败")
			return false
		}
		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 6: //任选五~任选八任选任选八
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量不得大于11个元素或小于8个元素，应为任选8 最多只能选11个数字 1 - 11 或最少选择8个数字
		if l > 11 || l < 8 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 任选8 8个数为一组
		singleBetNum := utils.AnalysisCombination(l, 8)
		order.SingleBetNum = singleBetNum
		//判断订单总额有没有超过限制
		if order.SingleBetAmount*float64(singleBetNum) > o.Settings[order.BetType].OrderLimit {
			//beego.Debug("失败")
			return false
		}

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 7: //任选一~任选四任选一
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l > 11 || l < 1 {
			//beego.Debug("失败")
			return false
		}

		//6分析订单注数 注意 前一直选 选了几个数就是几注
		singleBetNum := l
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 8: //前二前二直选
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			beego.Debug("失败")
			return false
		}

		//5每个数组数量不得大于11个元素或小于1个元素，应为前二直选每一位 最多只能选11个数字 1 - 11 或最少选择1个数字
		for _, v := range array {
			if len(v) > 11 || len(v) < 1 {
				beego.Debug("失败")
				return false
			}
		}

		//找出两组重复数对数
		var repateNum int = 0
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repateNum += 1
				}
			}
		}
		//6分析订单注数 (第一组元素个数 * 第二组元素个数) - 两组重复数字次数
		singleBetNum := len(array[0])*len(array[1]) - repateNum
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 10: //前三前三直选
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserThreeDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		//5每个数组数量不得大于11个元素或小于1个元素，应为前三组选每一位 最多只能选11个数字 1 - 11 或最少选择1个数字
		for _, v := range array {
			if len(v) > 11 || len(v) < 1 {
				//beego.Debug("失败")
				return false
			}
		}

		//分别找出12组, 23组, 13组 和 123组的重复数
		var repeatT12 int = 0
		var repeatT23 int = 0
		var repeatT13 int = 0
		var repeatT123 int = 0

		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					repeatT12 += 1
				}
			}
		}

		for _, i := range array[1] {
			for _, j := range array[2] {
				if i == j {
					repeatT23 += 1
				}
			}
		}

		for _, v := range array[0] {
			for _, j := range array[2] {
				if v == j {
					repeatT13 += 1
				}
			}
		}

		for _, v := range array[0] {
			for _, i := range array[1] {
				for _, j := range array[2] {
					if v == i && v == j {
						repeatT123 += 1
					}
				}
			}
		}

		//公式
		l1 := len(array[0])
		l2 := len(array[1])
		l3 := len(array[2])

		//6分析订单注数(l1 * l2 * l3) - repeatT12 * l3 -repeatT23 * l1 - repeatT13 * l2 + repeatT123 * 2
		singleBetNum := l1*l2*l3 - repeatT12*l3 - repeatT23*l1 - repeatT13*l2 + repeatT123*2
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 12, 19: //任选一~任选四胆拖任选二, 前二前二胆拖
		//4 分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l1 := len(array[0])
		l2 := len(array[1])
		//5 任选2胆拖 胆码只能选择一个,拖码必须 >= 1 并且 < 11
		if l1 != 1 {
			return false
		}
		if l2 < 1 || l2 > 10 {
			return false
		}

		//不能有数字重复
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					return false
				}
			}
		}

		//6分析订单注数 (任选2胆拖, 由于 胆码只能选1个 并且拖码不能和胆码重复,那么订单数 == 拖码数)
		singleBetNum := l2
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 13, 20: //任选一~任选四胆拖任选三, 前三前三胆拖
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		//任选3胆拖的玩法 ,也只有两个数组
		l1 := len(array[0])
		l2 := len(array[1])

		//4判断胆码和拖码选择的个数是否合法 任选3胆拖，和前三组选胆拖都只能选择两个胆码
		if l1 == 1 {
			if l2 < 2 || l2 > 10 {
				return false
			}
		} else if l1 == 2 {
			if l2 < 1 || l2 > 9 {
				return false
			}
		} else {
			return false
		}

		//胆码拖码 数字是否有重复
		for _, v := range array[0] {
			for _, j := range array[1] {
				if v == j {
					return false
				}
			}
		}

		var singleBetNum int = 0

		//6分析订单注数  如果 胆码为1 那么拖码的排列组合就等于注数
		if l1 == 1 {
			singleBetNum = utils.AnalysisCombination(l2, 2)
		} else if l1 == 2 { //如果胆码个数为2,那么注数就等于拖码个数
			singleBetNum = l2
		} else {
			return false
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 14: //任选一~任选四胆拖任选四
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		//任选4胆拖的玩法 ,也只有两个数组
		l1 := len(array[0])
		l2 := len(array[1])

		//5判断胆码和拖码选择的个数是否合法
		//任选4胆拖 胆码最少选择一个 最多选择3个数
		if l1 < 1 || l1 > 3 {
			return false
		}
		//根据胆码数量,拖码的最小和最大选择数是变化的
		if l2 < 4-l1 || l2 > 11-l1 {
			return false
		}

		//胆码拖码 数字是否有重复
		for _, v := range array[0] {
			for _, j := range array[1] {
				if v == j {
					return false
				}
			}
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l2, 4-l1)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 15: //任选五~任选八胆拖任选五
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		//任选5胆拖的玩法 ,也只有两个数组
		l1 := len(array[0])
		l2 := len(array[1])

		//5判断胆码和拖码选择的个数是否合法
		//任选5胆拖 胆码最少选择一个 最多选择4个数
		if l1 < 1 || l1 > 4 {
			return false
		}
		//根据胆码数量,拖码的最小和最大选择数是变化的
		if l2 < 5-l1 || l2 > 11-l1 {
			return false
		}

		//胆码拖码 数字是否有重复
		for _, v := range array[0] {
			for _, j := range array[1] {
				if v == j {
					return false
				}
			}
		}

		//6分析订单注数 等于 2组元素个数 的 5-1组元素个数的排列组合
		singleBetNum := utils.AnalysisCombination(l2, 5-l1)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 16: //任选五~任选八胆拖任选六
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		//任选5胆拖的玩法 ,也只有两个数组
		l1 := len(array[0])
		l2 := len(array[1])

		//5判断胆码和拖码选择的个数是否合法
		//任选6胆拖 胆码最少选择一个 最多选择5个数
		if l1 < 1 || l1 > 5 {
			return false
		}
		//根据胆码数量,拖码的最小和最大选择数是变化的
		if l2 < 6-l1 || l2 > 11-l1 {
			return false
		}

		//胆码拖码 数字是否有重复
		for _, v := range array[0] {
			for _, j := range array[1] {
				if v == j {
					return false
				}
			}
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l2, 6-l1)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 17: //任选五~任选八胆拖任选七
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		//任选5胆拖的玩法 ,也只有两个数组
		l1 := len(array[0])
		l2 := len(array[1])

		//5判断胆码和拖码选择的个数是否合法
		//任选7胆拖 胆码最少选择一个 最多选择6个数
		if l1 < 1 || l1 > 6 {
			return false
		}
		//根据胆码数量,拖码的最小和最大选择数是变化的
		if l2 < 7-l1 || l2 > 11-l1 {
			return false
		}

		//胆码拖码 数字是否有重复
		for _, v := range array[0] {
			for _, j := range array[1] {
				if v == j {
					return false
				}
			}
		}

		//6分析订单注数 等于 2组元素个数 的 5-1组元素个数的排列组合
		singleBetNum := utils.AnalysisCombination(l2, 7-l1)
		order.SingleBetNum = singleBetNum

		//判断订单总额有没有超过限制
		if order.SingleBetAmount*float64(singleBetNum) > o.Settings[order.BetType].OrderLimit {
			//beego.Debug("失败")
			return false
		}

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 18: //任选五~任选八胆拖任选八
		//4 分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		//5 任选5胆拖的玩法 ,也只有两个数组
		l1 := len(array[0])
		l2 := len(array[1])

		//判断胆码和拖码选择的个数是否合法
		//任选8胆拖 胆码最少选择一个 最多选择7个数
		if l1 < 1 || l1 > 7 {
			return false
		}
		//根据胆码数量,拖码的最小和最大选择数是变化的
		if l2 < 8-l1 || l2 > 11-l1 {
			return false
		}

		//胆码拖码 数字是否有重复
		for _, v := range array[0] {
			for _, j := range array[1] {
				if v == j {
					return false
				}
			}
		}

		//6分析订单注数 等于 2组元素个数 的 5-1组元素个数的排列组合
		singleBetNum := utils.AnalysisCombination(l2, 8-l1)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败")
			return false
		}

	case 21: //定位胆定位胆定位胆
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
			if l < 0 || l > 11 {
				//beego.Debug("失败")
				return false
			}
			singleBetNum += l
		}

		//5个数组的长度不得小于1 大于50
		if tl < 1 || tl > 55 {
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

	case 22, 23, 24: //总和总和大小, 总和总和尾数大小,总和总和单双,趣味趣味龙虎
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumForBigSmall(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大2 应为大小玩法,只能选择大小两个数 大为1, 小为0
		if l != 1 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(大小玩法单双玩法,注数不是1,就是2)
		singleBetNum := l
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

	// case 26: //趣味趣味猜单双
	// 	//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
	// 	ok, array := o.PaserNormalBetNumFor26(order.BetNums)
	// 	if !ok {
	// 		//beego.Debug("失败")
	// 		return false
	// 	}

	// 	l := len(array)

	// 	if l < 1 || l > 6 {
	// 		beego.Debug("失败")
	// 		return false
	// 	}

	// 	//由于这个玩法的特殊性,每位下注数只能是0,1,2,3,4,5
	// 	if array[0] < 0 || array[0] > 5 {
	// 		return false
	// 	}

	// 	//6分析订单注数
	// 	order.SingleBetNum = l

	// 	//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
	// 	if _, ok := order.Odds["1"]; ok {
	// 		for k, v := range order.Odds {
	// 			order.Odds[k] = v - v*order.Rebate
	// 		}
	// 	} else {
	// 		beego.Debug("失败")
	// 		return false
	// 	}

	// case 27: //趣味趣味猜中位
	// 	//分析订单下注数字是否正确(如果正确返回解析后的int数组)
	// 	ok, array := o.PaserNormalBetNum(order.BetNums)
	// 	if !ok {
	// 		return false
	// 	}
	// 	l := len(array)

	// 	//数组数量不得小于1 大10 应为前一玩法最少选择一个数字,最多选择10个数字,选择了几个就是几注
	// 	if l < 1 || l > 7 {
	// 		//beego.Debug("失败")
	// 		return false
	// 	}
	// 	//由于这个玩法的特殊性,每位下注数只能是 3 - 9
	// 	if array[0] < 3 || array[0] > 9 {
	// 		return false
	// 	}

	// 	//分析订单注数(前一玩法,选择了几个数就是几注)
	// 	singleBetNum := l
	// 	order.SingleBetNum = singleBetNum

	// 	//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
	// 	if odds, ok := order.Odds["1"]; ok {
	// 		order.Odds["1"] = odds - odds*order.Rebate
	// 	} else {
	// 		beego.Debug("失败1")
	// 		return false
	// 	}
	case 28: //两面万位
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor28(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大6
		if l < 1 || l > 6 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(大小玩法单双玩法,注数不是1,就是2)
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

	case 29: //两面千位
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor29(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大6
		if l < 1 || l > 5 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(大小玩法单双玩法,注数不是1,就是2)
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

	case 30: //两面百位
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor30(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大6
		if l < 1 || l > 4 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(大小玩法单双玩法,注数不是1,就是2)
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

	case 31: //两面十位
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor31(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大6
		if l < 1 || l > 3 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(大小玩法单双玩法,注数不是1,就是2)
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
	case 32, 33: //两面个位
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor32(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大6
		if l < 1 || l > 2 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(大小玩法单双玩法,注数不是1,就是2)
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *EX5) PaserNormalBetNum(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		//beego.Debug(v)
		i, err := strconv.Atoi(v)
		//beego.Debug(i)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
		if i < 1 || i > 11 {
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
func (o *EX5) PaserNormalBetNumFor26(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		//beego.Debug(v)
		i, err := strconv.Atoi(v)
		//beego.Debug(i)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
		if i < 0 || i > 5 {
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

func (o *EX5) PaserNormalBetNumFor28(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
		if i < 0 || i > 11 {
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

func (o *EX5) PaserNormalBetNumFor29(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
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

func (o *EX5) PaserNormalBetNumFor30(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
		if i < 0 || i > 7 {
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

func (o *EX5) PaserNormalBetNumFor31(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
		if i < 0 || i > 5 {
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

func (o *EX5) PaserNormalBetNumFor32(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于1和大于11
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

//解析下注号码,得到注数一维数组(这个由大小单双的判定调用 应为大小单双每个数字和其他的不一样 只能 是 0 或 1)
func (o *EX5) PaserNormalBetNumForBigSmall(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于1(PK10的 大小单双玩法)
		if i < 0 || i > 1 {
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
func (o *EX5) PaserTwoDigitBetNum(betNum string) (bool, [][]int) {
	//分割下注位
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
			//每位数字不能小于1和大于11
			if i < 1 || i > 11 {
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

//解析3位下注号码,得到注数二维数组(用于有3位选择数字的情况,例如:前三直选)
func (o *EX5) PaserThreeDigitBetNum(betNum string) (bool, [][]int) {
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
			//每位数字不能小于1和大于11
			if i < 1 || i > 11 {
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
func (o *EX5) PaserFiveDigitBetNumAllowSpaces(betNum string) (bool, [][]int) {
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
			if i < 1 || i > 11 {
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

//检查是否有重复字符和空字符(字符串数组)
func (o *EX5) CheckRepeatAndEmptyString(array []string) bool {
	//排个序先(字符串对比)
	sort.Strings(array)
	arrayLen := len(array)
	for i := 0; i < arrayLen; i++ {
		if (i > 0 && array[i-1] == array[i]) || len(array[i]) == 0 {
			return false
		}
	}
	return true
}

//检查是否有重复数字符(int数组)
func (o *EX5) CheckRepeatInt(array []int) bool {
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

//结算这个彩种当期所有订单
func (o *EX5) SettlementOrders(orders []gb.Order, openCode string) {
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
			beego.Error(err)
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
func (o *EX5) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}

	//转下注号码为数组
	switch order.BetType {
	case 0: //任选一~任选四任选任选二
		o.WinningAndLose_0(order, openCode)
	case 1: //任选一~任选四任选任选三
		o.WinningAndLose_1(order, openCode)
	case 2: //任选一~任选四任选任选四
		o.WinningAndLose_2(order, openCode)
	case 3: //任选五~任选八任选任选五
		o.WinningAndLose_3(order, openCode)
	case 4: //任选五~任选八任选任选六
		o.WinningAndLose_4(order, openCode)
	case 5: //任选五~任选八任选任选七
		o.WinningAndLose_5(order, openCode)
	case 6: //任选五~任选八任选任选八
		o.WinningAndLose_6(order, openCode)
	case 7: //任选一~任选四任选一
		o.WinningAndLose_7(order, openCode)
	case 8: //前二前二直选
		o.WinningAndLose_8(order, openCode)
	case 9: //前二前二组选
		o.WinningAndLose_9(order, openCode)
	case 10: //前三前三直选
		o.WinningAndLose_10(order, openCode)
	case 11: //前三前三组选
		o.WinningAndLose_11(order, openCode)
	case 12: //任选一~任选四胆拖任选二
		o.WinningAndLose_12(order, openCode)
	case 13: //任选一~任选四胆拖任选三
		o.WinningAndLose_13(order, openCode)
	case 14: //任选一~任选四胆拖任选四
		o.WinningAndLose_14(order, openCode)
	case 15: //任选五~任选八胆拖任选五
		o.WinningAndLose_15(order, openCode)
	case 16: //任选五~任选八胆拖任选六
		o.WinningAndLose_16(order, openCode)
	case 17: //任选五~任选八胆拖任选七
		o.WinningAndLose_17(order, openCode)
	case 18: //任选五~任选八胆拖任选八
		o.WinningAndLose_18(order, openCode)
	case 19: //前二前二组选胆拖
		o.WinningAndLose_19(order, openCode)
	case 20: //前三前三组选胆拖
		o.WinningAndLose_20(order, openCode)
	case 21: //定位胆定位胆定位胆
		o.WinningAndLose_21(order, openCode)
	case 22: //总和总和大小
		o.WinningAndLose_22(order, openCode)
	case 23: //总和总和尾数大小
		o.WinningAndLose_23(order, openCode)
	case 24: //总和总和单双
		o.WinningAndLose_24(order, openCode)
	// case 25: //趣味趣味龙虎
	// 	o.WinningAndLose_25(order, openCode)
	// case 26: //趣味趣味猜单双
	// 	o.WinningAndLose_26(order, openCode)
	// case 27: //趣味趣味猜中位
	//o.WinningAndLose_27(order, openCode)
	case 28: //总和总和单双
		o.WinningAndLose_28(order, openCode)
	case 29: //总和总和单双
		o.WinningAndLose_29(order, openCode)
	case 30: //总和总和单双
		o.WinningAndLose_30(order, openCode)
	case 31: //总和总和单双
		o.WinningAndLose_31(order, openCode)
	case 32: //总和总和单双
		o.WinningAndLose_32(order, openCode)
	case 33: //总和总和单双
		o.WinningAndLose_33(order, openCode)
	default:
		beego.Debug("失败")
		return false
	}
	return true
}

//判断是否中奖 (0 任选二)(返回用户输赢情况)(会有多注中奖情况)
func (o *EX5) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}
	l := len(winningNumbers)
	var winningBetNum int = 0
	//如果中奖数组<2 证明没有中奖
	if l < 2 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 2 {
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
	} else {
		winningBetNum = utils.AnalysisCombination(l, 2)
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

//判断是否中奖 (1 任选三)(返回用户输赢情况)(会有多注中奖情况)
func (o *EX5) WinningAndLose_1(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}
	l := len(winningNumbers)
	var winningBetNum int = 0
	//如果中奖数组<3 证明没有中奖
	if l < 3 {
		//中的注数
		winningBetNum = 0
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 3 {
		//中的注数
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		if v, ok := order.Odds["1"]; ok {
			//结算 = 赔率 * 单注金额 * 中奖注数
			ret += v * order.SingleBetAmount * 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	} else {
		winningBetNum = utils.AnalysisCombination(l, 3)
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (2 任选四)(返回用户输赢情况)(会有多注中奖情况)
func (o *EX5) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}
	l := len(winningNumbers)
	var winningBetNum int = 0
	//如果中奖数组< 4 证明没有中奖
	if l < 4 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 4 {
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
		winningBetNum = utils.AnalysisCombination(l, 4)
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//给用户加钱

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

//判断是否中奖 (3 任选五)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_3(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}
	l := len(winningNumbers)
	var winningBetNum int = 0
	//如果中奖数组< 5 证明没有中奖
	if l < 5 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 5 {
		winningBetNum = 1
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * 1
			//给用户加钱

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

//判断是否中奖 (4 任选六)(返回用户输赢情况)(会有多注中奖的情况)
func (o *EX5) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}
	var winningBetNum int = 0
	l := len(winningNumbers)
	//如果中奖数组< 5 证明没有中奖
	if l < 5 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 5 {
		//任选6 选择的数个数 - 5 就是中奖数
		winningBetNum = len(betNumbers) - 5
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (5 任选七)(返回用户输赢情况)(会有多注中奖的情况)
func (o *EX5) WinningAndLose_5(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}
	var winningBetNum int = 0
	l := len(winningNumbers)
	//如果中奖数组< 5 证明没有中奖
	if l < 5 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 5 {
		//任选7 选择的数个数 - 5 剩下的数字进行两两排列组合就是中奖数
		winningBetNum = utils.AnalysisCombination(len(betNumbers)-5, 2)
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (6 任选八)(返回用户输赢情况)(会有多注中奖的情况)
func (o *EX5) WinningAndLose_6(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
			}
		}
	}

	var winningBetNum int = 0
	l := len(winningNumbers)
	//如果中奖数组< 5 证明没有中奖
	if l < 5 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if l == 5 {
		//任选八 选择的数个数 - 5 剩下的数字进行三三排列组合就是中奖注数
		winningBetNum = utils.AnalysisCombination(len(betNumbers)-5, 3)
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (7 任选一~任选四任选一)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_7(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum int = 0

	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningBetNum++
				break
			}
		}
	}

	//如果中奖数组< 5 证明没有中奖
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
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (8 前二直选)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_8(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断有没有 按位开出下注的号码
	var flag_1 = 0
	var flag_2 = 0

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

	var winningBetNum = 0
	//前两位都相同才中奖
	if flag_1 == 1 && flag_2 == 1 {
		winningBetNum = 1
	}

	if winningBetNum < 1 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum == 1 {
		//前二直选只会有一注中奖
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (9 前二组选)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_9(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断有没有 按位开出下注的号码
	var flag_1 = 0
	var flag_2 = 0

	for _, v := range betNumbers {
		if v == openCode[0] {
			flag_1 = 1
			break
		}
	}

	if flag_1 == 1 {
		for _, v := range betNumbers {
			if v == openCode[1] {
				flag_2 = 1
				break
			}
		}
	}
	var winningBetNum = 0
	//开出的前两个数字都在下注数字里面就算中奖不用按位计算
	if flag_1 == 1 && flag_2 == 1 {
		winningBetNum = 1
	}
	//如果中奖数组<2 证明没有中奖
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

	order.Status = 1
}

//判断是否中奖 (10 前三直选)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_10(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断有没有 按位开出下注的号码
	var flag_0 = 0
	var flag_1 = 0

	for _, v := range betNumbers[0] {
		if v == openCode[0] {
			flag_0 = 1
			break
		}
	}

	if flag_0 == 1 {
		for _, v := range betNumbers[1] {
			if v == openCode[1] {
				flag_1 = 1
				break
			}
		}
		if flag_1 == 1 {
			for _, v := range betNumbers[2] {
				if v == openCode[2] {
					winningBetNum = 1
					break
				}
			}
		}
	}

	if winningBetNum < 1 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum == 1 {
		//前三直选只会有一注中奖
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

//判断是否中奖 (11 前三组选)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_11(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断有没有 按位开出下注的号码
	var flag_1 = 0
	var flag_2 = 0
	var flag_3 = 0
	//第一个号码
	for _, v := range betNumbers {
		if v == openCode[0] {
			flag_1 = 1
			break
		}
	}
	//第二个号码
	if flag_1 == 1 {
		for _, v := range betNumbers {
			if v == openCode[1] {
				flag_2 = 1
				break
			}
		}
	}
	//第三个号码
	if flag_2 == 1 {
		for _, v := range betNumbers {
			if v == openCode[2] {
				flag_3 = 1
				break
			}
		}
	}
	//开出的前3个数字都在选择的数组中就算中奖
	if flag_1 == 1 && flag_2 == 1 && flag_3 == 1 {
		winningBetNum = 1
	}

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

//判断是否中奖 (12 任选二胆拖)(返回用户输赢情况)(会有多注中奖情况)
func (o *EX5) WinningAndLose_12(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0
	//判断下注号码中有几个开出
	for _, v := range openCode {
		//胆码有没有中
		if v == betNumbers[0][0] {
			for _, i := range openCode {
				for _, x := range betNumbers[1] {
					if i == x {
						winningBetNum++
					}
				}
			}
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
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		//反水
		order.RebateAmount = ret
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
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

//判断是否中奖 (13 任选三胆拖)(返回用户输赢情况)(会有多注中奖情况)
func (o *EX5) WinningAndLose_13(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumberCount = 0
	//胆码
	l1 := len(betNumbers[0])

	if l1 == 1 {
		for _, v := range openCode {
			//胆码有没有中
			if v == betNumbers[0][0] {
				for _, i := range openCode {
					for _, x := range betNumbers[1] {
						if i == x {
							winningNumberCount++
						}
					}
				}
				break
			}
		}
	} else if l1 == 2 {
		//先判断两个胆码有没有开出
		t := 0
		for _, v := range openCode {
			for _, i := range betNumbers[0] {
				if v == i {
					t++
				}
			}
		}
		//两个胆码中了
		if t == 2 {
			for _, i := range openCode {
				for _, x := range betNumbers[1] {
					if i == x {
						winningNumberCount++
					}
				}
			}
		}
	} else {
		return
	}
	var winningBetNum = 0
	if l1 == 1 {
		if winningNumberCount < 2 {
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		} else if winningNumberCount >= 2 {
			winningBetNum = utils.AnalysisCombination(winningNumberCount, 2)
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		}
	} else if l1 == 2 {
		if winningNumberCount < 1 {
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		} else if winningNumberCount >= 1 {
			winningBetNum = winningNumberCount
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningNumberCount)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningNumberCount)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (14 任选四胆拖)(返回用户输赢情况)(会有多注中奖情况)
func (o *EX5) WinningAndLose_14(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//胆码
	l1 := len(betNumbers[0])
	//先判断胆码有几个中奖
	d := 0
	for _, v := range betNumbers[0] {
		for _, i := range openCode {
			if v == i {
				d++
			}
		}
	}

	var winningBetNum = 0
	//胆码中奖
	if d == l1 {
		t := 0
		for _, v := range betNumbers[1] {
			for _, i := range openCode {
				if v == i {
					t++
				}
			}
		}

		if t < 4-l1 {
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		} else {
			winningBetNum = utils.AnalysisCombination(t, 4-l1)
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
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

//判断是否中奖 (15 任选五胆拖)(返回用户输赢情况)(只会有一注中奖)
func (o *EX5) WinningAndLose_15(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//胆码
	l1 := len(betNumbers[0])
	//先判断胆码有几个中奖
	d := 0
	for _, v := range betNumbers[0] {
		for _, i := range openCode {
			if v == i {
				d++
			}
		}
	}
	var winningBetNum = 0
	//胆码中奖
	if d == l1 {
		t := 0
		for _, v := range betNumbers[1] {
			for _, i := range openCode {
				if v == i {
					t++
				}
			}
		}

		if t < 5-l1 {
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		} else if t == 5-l1 {
			winningBetNum = 1
			//任选5只会有一注中奖
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

//判断是否中奖 (16 任选六胆拖)(返回用户输赢情况)
func (o *EX5) WinningAndLose_16(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//胆码数
	l1 := len(betNumbers[0])
	//拖码数
	l2 := len(betNumbers[1])
	//先判断胆码有几个中奖
	d := 0
	for _, v := range betNumbers[0] {
		for _, i := range openCode {
			if v == i {
				d++
			}
		}
	}
	//拖码中奖数
	t := 0
	for _, v := range betNumbers[1] {
		for _, i := range openCode {
			if v == i {
				t++
			}
		}
	}

	var winningBetNum = 0
	//开出了选择的5个号码的情况
	if d+t == 5 {
		//胆码全中的情况
		if l1 == d {
			winningBetNum = l2 - t
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1-d == 1 { //胆错1
			//中的注数
			winningBetNum = 1
			//达成这个条件中一注
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

//判断是否中奖 (17 任选七胆拖)(返回用户输赢情况)
func (o *EX5) WinningAndLose_17(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//胆码数
	l1 := len(betNumbers[0])
	//拖码数
	l2 := len(betNumbers[1])
	//先判断胆码有几个中奖
	d := 0
	for _, v := range betNumbers[0] {
		for _, i := range openCode {
			if v == i {
				d++
			}
		}
	}
	//拖码中奖数
	t := 0
	for _, v := range betNumbers[1] {
		for _, i := range openCode {
			if v == i {
				t++
			}
		}
	}
	var winningBetNum = 0
	//开出了选择的5个号码的情况（中奖的情况）
	if d+t == 5 {
		if l1+l2 == 7 {
			//只投了一注
			winningBetNum = 1
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1 == d {
			//胆码全中的情况，l2 - t 进行 7 - l1 - t 排列组合
			winningBetNum = utils.AnalysisCombination(l2-t, 7-l1-t)
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1-d == 1 {
			//胆错1
			winningBetNum = l2 - t
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1-d == 2 {
			winningBetNum = 1
			//胆错2
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

//判断是否中奖 (18 任选八胆拖)(返回用户输赢情况)
func (o *EX5) WinningAndLose_18(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//胆码数
	l1 := len(betNumbers[0])
	//拖码数
	l2 := len(betNumbers[1])
	//先判断胆码有几个中奖
	d := 0
	for _, v := range betNumbers[0] {
		for _, i := range openCode {
			if v == i {
				d++
			}
		}
	}
	//拖码中奖数
	t := 0
	for _, v := range betNumbers[1] {
		for _, i := range openCode {
			if v == i {
				t++
			}
		}
	}
	var winningBetNum = 0
	//开出了选择的5个号码的情况（中奖的情况）
	if d+t == 5 {
		if l1+l2 == 8 {
			//只投了一注
			winningBetNum = 1
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * 1
				//给用户加钱

				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1 == d {
			//胆码全中的情况，l2 - t 进行 8 - l1 - t 排列组合
			winningBetNum = utils.AnalysisCombination(l2-t, 8-l1-t)
			beego.Debug("1: ", winningBetNum)
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1-d == 1 {
			//胆错1(任选8 胆全中，和胆错一时一样的公式)
			winningBetNum = utils.AnalysisCombination(l2-t, 8-l1-t)
			beego.Debug("2: ", winningBetNum)
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1-d == 2 {
			//胆错2（相当于 任选7 的胆错1）
			winningBetNum = l2 - t
			beego.Debug("3: ", winningBetNum)
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
				//更新order
				order.Settlement = ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		} else if l1-d == 3 {
			winningBetNum = 1
			//胆错3 （只有一注中）
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

//判断是否中奖 (19 前二组选胆拖)(返回用户输赢情况)
func (o *EX5) WinningAndLose_19(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0
	//判断胆码有没有中奖
	if betNumbers[0][0] == openCode[0] || betNumbers[0][0] == openCode[1] {
		for _, v := range betNumbers[1] {
			if v == openCode[0] || v == openCode[1] {
				winningBetNum = 1
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
		winningBetNum = 1
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

//判断是否中奖 (20 前三组选胆拖)(返回用户输赢情况)
func (o *EX5) WinningAndLose_20(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//胆码数
	l1 := len(betNumbers[0])

	var flag = 0
	var flag_2 = 0

	var winningBetNum = 0
	var tmpOpenCode = openCode[:3]
	if l1 == 1 {
		for _, v := range tmpOpenCode {
			if v == betNumbers[0][0] {
				for _, i := range betNumbers[1] {
					for _, j := range tmpOpenCode {
						if i == j {
							flag_2++
						}
					}
				}
				break
			}
		}

		if flag_2 == 2 {
			winningBetNum = 1
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
			//反水
			order.RebateAmount = ret
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
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
	} else if l1 == 2 {
		for i := 0; i < 3; i++ {
			for _, v := range betNumbers[0] {
				if v == openCode[i] {
					flag++
				}
			}
		}
		//胆码中
		if flag == 2 {
			for i := 0; i < 3; i++ {
				for _, v := range betNumbers[1] {
					if openCode[i] == v {
						flag_2 = 1
						break
					}
				}
				if flag_2 == 1 {
					break
				}
			}
			//拖码也中
			if flag_2 == 1 {
				winningBetNum = 1
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
		} else {
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (21 定位胆定位胆定位胆)(返回用户输赢情况)
func (o *EX5) WinningAndLose_21(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNumAllowSpaces(order.BetNums)
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

	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers > 0 && winningNumbers < 6 {
		//定位胆 1-5 最多只有5注中奖
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

//判断是否中奖 (22 总和总和大小)(返回用户输赢情况)
func (o *EX5) WinningAndLose_22(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//计算5个开奖号码的和值
	retSum := openCode[0] + openCode[1] + openCode[2] + openCode[3] + openCode[4]
	//判断
	var big = 0
	var small = 0
	var sum = 0

	if retSum >= 31 {
		big = 1
	} else if retSum <= 29 {
		small = 1
	} else {
		sum = 1
	}

	var winningBetNum = 0

	if big == 1 {
		if betNumbers[0] == 0 {
			winningBetNum = 1
			//中奖就没有反水,应为只下了一注
			order.RebateAmount = 0
			//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
	} else if small == 1 {
		if betNumbers[0] == 1 {
			winningBetNum = 1
			//中奖就没有反水,应为只下了一注
			order.RebateAmount = 0
			//计算中了多少钱(由于 每一位的中奖赔率不一样,所以这里要每一位的中奖都分开计算)
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * 1
				//更新order
				order.Settlement = ret
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
	} else if sum == 1 {
		winningBetNum = 1
		ret += order.SingleBetAmount
	}
	//计算反水
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//最后结算金额(中奖金额 + 反水)
	order.Settlement = ret + order.RebateAmount

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (23  总和总和尾数大小)(返回用户输赢情况)
func (o *EX5) WinningAndLose_23(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	Sum := openCode[0] + openCode[1] + openCode[2] + openCode[3] + openCode[4]

	ds := (Sum % 10)
	//默认大为 0, 小为 1
	var d = 0
	if ds < 5 {
		d = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == d {
		winningBetNum = 1

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

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (24  总和总和单双)(返回用户输赢情况)
func (o *EX5) WinningAndLose_24(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	Sum := openCode[0] + openCode[1] + openCode[2] + openCode[3] + openCode[4]

	ds := (Sum % 10) % 2
	//默认单为 0, 双为 1
	var d = 0
	if ds == 0 {
		d = 1
	}

	var winningBetNum = 0
	if betNumbers[0] == d {
		winningBetNum = 1

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

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (25 趣味趣味龙虎)(返回用户输赢情况)
// func (o *EX5) WinningAndLose_25(order *gb.Order, openCode []int) {
// 	//最后结果 输赢多少钱
// 	var ret float64
// 	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
// 	if !ok {
// 		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
// 		return
// 	}
// 	//找出前两位
// 	var tmpOpenCode []int
// 	tmpOpenCode = append(tmpOpenCode, openCode[0])
// 	tmpOpenCode = append(tmpOpenCode, openCode[4])

// 	var dragon = 0
// 	var tiger = 0

// 	if tmpOpenCode[0] > tmpOpenCode[1] {
// 		dragon = 1
// 	} else if tmpOpenCode[0] < tmpOpenCode[1] {
// 		tiger = 1
// 	}

// 	var winningBetNum = 0
// 	if betNumbers[0] == 0 {
// 		if dragon == 1 {
// 			winningBetNum = 1
// 		}
// 	} else if betNumbers[0] == 1 {
// 		if tiger == 1 {
// 			winningBetNum = 1
// 		}
// 	}

// 	if winningBetNum == 0 {
// 		//一注没中 计算反水 单注金额 * 反水 * 单注数量
// 		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
// 		//反水
// 		order.RebateAmount = ret
// 		//更新order
// 		order.Settlement = ret
// 	} else if winningBetNum == 1 {
// 		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
// 		order.RebateAmount = ret
// 		//计算中了多少钱
// 		if v, ok := order.Odds["1"]; ok {
// 			ret += v * order.SingleBetAmount * 1
// 			//更新order
// 			order.Settlement = ret
// 		} else {
// 			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
// 			return
// 		}
// 	}

// 	order.WinningBetNum = winningBetNum
// 	//订单以结算
// 	order.Status = 1
// }

// //判断是否中奖 (26  总和总和单双)(返回用户输赢情况)
// func (o *EX5) WinningAndLose_26(order *gb.Order, openCode []int) {
// 	//最后结果 输赢多少钱
// 	var ret float64
// 	ok, betNumbers := o.PaserNormalBetNumFor26(order.BetNums)
// 	if !ok {
// 		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
// 		return
// 	}

// 	//默认单为 0 双为 1
// 	var d = [5]int{0, 0, 0, 0, 0}

// 	if openCode[0]%2 == 0 {
// 		d[0] = 1
// 	}

// 	if openCode[1]%2 == 0 {
// 		d[1] = 1
// 	}

// 	if openCode[2]%2 == 0 {
// 		d[2] = 1
// 	}

// 	if openCode[3]%2 == 0 {
// 		d[3] = 1
// 	}

// 	if openCode[4]%2 == 0 {
// 		d[4] = 1
// 	}

// 	//计算双个数
// 	var countD = 0
// 	for _, v := range d {
// 		if v == 1 {
// 			countD++
// 		}
// 	}

// 	var winningBetNum = 0
// 	for _, v := range betNumbers {
// 		if countD == v {
// 			winningBetNum = 1
// 			break
// 		}
// 	}

// 	if winningBetNum == 0 {
// 		//一注没中 计算反水 单注金额 * 反水 * 单注数量
// 		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
// 		//反水
// 		order.RebateAmount = ret
// 		//更新order
// 		order.Settlement = ret
// 	} else if winningBetNum == 1 {
// 		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
// 		order.RebateAmount = ret
// 		//计算中了多少钱
// 		if v, ok := order.Odds["1"]; ok {
// 			ret += v * order.SingleBetAmount * 1
// 			//更新order
// 			order.Settlement = ret
// 		} else {
// 			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
// 			return
// 		}
// 	}

// 	order.WinningBetNum = winningBetNum
// 	//订单以结算
// 	order.Status = 1
// }

// //判断是否中奖 (27 趣味趣味猜中位)(返回用户输赢情况)
// func (o *EX5) WinningAndLose_27(order *gb.Order, openCode []int) {
// 	//最后结果 输赢多少钱
// 	var ret float64
// 	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
// 	if !ok {
// 		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
// 		return
// 	}

// 	//从小到大排序(注意 openCode 是一个切片 必须是copy操作才行,否者tmpOpenCode的排序,会改变openCode的值)
// 	tmpOpenCode := make([]int, 5)
// 	copy(tmpOpenCode, openCode)

// 	sort.Sort(common.IntSlice(tmpOpenCode))

// 	var winningBetNum = 0
// 	for _, v := range betNumbers {
// 		if v == tmpOpenCode[2] {
// 			winningBetNum = 1
// 			break
// 		}
// 	}

// 	if winningBetNum == 0 {
// 		//一注没中 计算反水 单注金额 * 反水 * 单注数量
// 		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
// 		//反水
// 		order.RebateAmount = ret
// 		//更新order
// 		order.Settlement = ret
// 	} else if winningBetNum == 1 {
// 		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-1)
// 		order.RebateAmount = ret
// 		//计算中了多少钱
// 		if v, ok := order.Odds["1"]; ok {
// 			ret += v * order.SingleBetAmount * 1
// 			//更新order
// 			order.Settlement = ret
// 		} else {
// 			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
// 			return
// 		}
// 	}

// 	order.WinningBetNum = winningBetNum
// 	//订单以结算
// 	order.Status = 1
// }

//判断是否中奖 28 两面万位
func (o *EX5) WinningAndLose_28(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor28(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [6]int

	//大 0 , 小 1 , 单 2, 双 3, 龙万 4, 虎千 5, 龙万 6, 虎百 7, 龙万 8, 虎十 9, 龙万 10, 虎个 11
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
	} else {
		resultCode[2] = 5
	}

	if o.openCode[0] > o.openCode[2] {
		resultCode[3] = 6
	} else {
		resultCode[3] = 7
	}

	if o.openCode[0] > o.openCode[3] {
		resultCode[4] = 8
	} else {
		resultCode[4] = 9
	}

	if o.openCode[0] > o.openCode[4] {
		resultCode[5] = 10
	} else {
		resultCode[5] = 11
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

//判断是否中奖 29 两面千位
func (o *EX5) WinningAndLose_29(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor29(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [5]int

	//大 0 , 小 1 , 单 2, 双 3, 龙千 4, 虎百 5, 龙千 6, 虎十 7, 龙千 8, 虎个 9
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
	} else {
		resultCode[2] = 5
	}

	if o.openCode[1] > o.openCode[3] {
		resultCode[3] = 6
	} else {
		resultCode[3] = 7
	}

	if o.openCode[1] > o.openCode[4] {
		resultCode[4] = 8
	} else {
		resultCode[4] = 9
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

//判断是否中奖 30 两面百位
func (o *EX5) WinningAndLose_30(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor30(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [4]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 龙百 6, 虎个 7
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
	} else {
		resultCode[2] = 5
	}

	if o.openCode[2] > o.openCode[4] {
		resultCode[3] = 6
	} else {
		resultCode[3] = 7
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

//判断是否中奖 31 两面十位
func (o *EX5) WinningAndLose_31(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1 , 单 2, 双 3, 龙十 4, 虎个 5
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
	} else {
		resultCode[2] = 5
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

//判断是否中奖 32 两面个位
func (o *EX5) WinningAndLose_32(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor32(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1 , 单 2, 双 3
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

//判断是否中奖 33 两面总和
func (o *EX5) WinningAndLose_33(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor32(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	sum := o.openCode[0] + o.openCode[1] + o.openCode[2] + o.openCode[3] + o.openCode[4]

	//和大 0 , 和小 1 , 和单 2, 和双 3
	if sum > 30 {
		resultCode[0] = 0
	} else if sum < 30 {
		resultCode[0] = 1
	}

	if sum%2 == 0 {
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
