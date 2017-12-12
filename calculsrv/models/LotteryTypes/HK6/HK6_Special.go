//这个文件放置彩种特殊的一些功能 比如 HK6 的API获取没有下一期 只有自己计算
package HK6

import (
	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/acmgr"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"
	"common/utils"

	"calculsrv/models/Order"

	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

//红波 0
var red = [17]int{1, 2, 7, 8, 12, 13, 18, 19, 23, 24, 29, 30, 34, 35, 40, 45, 46}

//蓝波 1
var blue = [16]int{3, 4, 9, 10, 14, 15, 20, 25, 26, 31, 36, 37, 41, 42, 47, 48}

//绿波 2
var green = [16]int{5, 6, 11, 16, 17, 21, 22, 27, 28, 32, 33, 38, 39, 43, 44, 49}

//色彩
var color map[int]int = make(map[int]int)

//生肖 0.鼠 1.牛 2.虎 3.兔 4.龙 5.蛇 6.马 7.羊 8.猴 9.鸡 10.狗 11.猪
var zodiac map[int]int = make(map[int]int)

//五行 0.金 1.水 2.木 3.火 4.土
var fiveLine map[int]int = make(map[int]int)

//进行中奖结果计算
// func (o *HK6) CalculateWinningResult(openCode []int) {

// }

func (o *HK6) InitNumProperty() {
	color[1] = 0
	color[2] = 0
	color[3] = 1
	color[4] = 1
	color[5] = 2
	color[6] = 2
	color[7] = 0
	color[8] = 0
	color[9] = 1
	color[10] = 1
	color[11] = 2
	color[12] = 0
	color[13] = 0
	color[14] = 1
	color[15] = 1
	color[16] = 2
	color[17] = 2
	color[18] = 0
	color[19] = 0
	color[20] = 1
	color[21] = 2
	color[22] = 2
	color[23] = 0
	color[24] = 0
	color[25] = 1
	color[26] = 1
	color[27] = 2
	color[28] = 2
	color[29] = 0
	color[30] = 0
	color[31] = 1
	color[32] = 2
	color[33] = 2
	color[34] = 0
	color[35] = 0
	color[36] = 1
	color[37] = 1
	color[38] = 2
	color[39] = 2
	color[40] = 0
	color[41] = 1
	color[42] = 1
	color[43] = 2
	color[44] = 2
	color[45] = 0
	color[46] = 0
	color[47] = 1
	color[48] = 1
	color[49] = 2

	zodiac[1] = 9
	zodiac[2] = 8
	zodiac[3] = 7
	zodiac[4] = 6
	zodiac[5] = 5
	zodiac[6] = 4
	zodiac[7] = 3
	zodiac[8] = 2
	zodiac[9] = 1
	zodiac[10] = 0
	zodiac[11] = 11
	zodiac[12] = 10
	zodiac[13] = 9
	zodiac[14] = 8
	zodiac[15] = 7
	zodiac[16] = 6
	zodiac[17] = 5
	zodiac[18] = 4
	zodiac[19] = 3
	zodiac[20] = 2
	zodiac[21] = 1
	zodiac[22] = 0
	zodiac[23] = 11
	zodiac[24] = 10
	zodiac[25] = 9
	zodiac[26] = 8
	zodiac[27] = 7
	zodiac[28] = 6
	zodiac[29] = 5
	zodiac[30] = 4
	zodiac[31] = 3
	zodiac[32] = 2
	zodiac[33] = 1
	zodiac[34] = 0
	zodiac[35] = 11
	zodiac[36] = 10
	zodiac[37] = 9
	zodiac[38] = 8
	zodiac[39] = 7
	zodiac[40] = 6
	zodiac[41] = 5
	zodiac[42] = 4
	zodiac[43] = 3
	zodiac[44] = 2
	zodiac[45] = 1
	zodiac[46] = 0
	zodiac[47] = 11
	zodiac[48] = 10
	zodiac[49] = 9

	fiveLine[1] = 3
	fiveLine[2] = 3
	fiveLine[3] = 0
	fiveLine[4] = 0
	fiveLine[5] = 1
	fiveLine[6] = 1
	fiveLine[7] = 2
	fiveLine[8] = 2
	fiveLine[9] = 3
	fiveLine[10] = 3
	fiveLine[11] = 4
	fiveLine[12] = 4
	fiveLine[13] = 1
	fiveLine[14] = 1
	fiveLine[15] = 2
	fiveLine[16] = 2
	fiveLine[17] = 0
	fiveLine[18] = 0
	fiveLine[19] = 4
	fiveLine[20] = 4
	fiveLine[21] = 1
	fiveLine[22] = 1
	fiveLine[23] = 3
	fiveLine[24] = 3
	fiveLine[25] = 0
	fiveLine[26] = 0
	fiveLine[27] = 4
	fiveLine[28] = 4
	fiveLine[29] = 2
	fiveLine[30] = 2
	fiveLine[31] = 3
	fiveLine[32] = 3
	fiveLine[33] = 0
	fiveLine[34] = 0
	fiveLine[35] = 1
	fiveLine[36] = 1
	fiveLine[37] = 2
	fiveLine[38] = 2
	fiveLine[39] = 3
	fiveLine[40] = 3
	fiveLine[41] = 4
	fiveLine[42] = 4
	fiveLine[43] = 1
	fiveLine[44] = 1
	fiveLine[45] = 2
	fiveLine[46] = 2
	fiveLine[47] = 0
	fiveLine[48] = 0
	fiveLine[49] = 4
}

//分析订单(下注)
func (o *HK6) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {

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
	case 0, 1, 2, 3, 4, 5, 6: //两面两面正一（多赔率）//两面两面正二（多赔率）//两面两面正三（多赔率）//两面两面正四（多赔率）//两面两面正五（多赔率）//两面两面正六（多赔率） //两面两面特码（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_0(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)
		//最少下一个号,最多下8个号
		if l > 8 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
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

	case 7: //两面两面总和（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_7(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)
		//最少下一个号,最多下8个号
		if l > 4 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
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

	case 8, 9, 10, 12, 13, 14, 15, 16, 17: //正码正码正码 //特码AB特码AB特码A //特码AB特码AB特码B（赔率是根据最高返利率来定）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 49 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败3")
			return false
		}

	case 11: //特码AB特码AB其他（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_11(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)
		//最少下一个号,最多下8个号
		if l > 14 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
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

	case 18, 19, 20, 21, 22, 23: //正码~正码~正码一（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_18(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)
		//最少下一个号,最多下8个号
		if l > 13 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
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

	case 24: //正码过关正码过关正码过关（多赔率中嵌套多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserSixDigitBetNumAllowSpaces(order.BetNums)
		if !ok {
			return false
		}

		var betCount int = 0
		for _, v := range array {
			l := len(v)
			if l != 0 {
				//beego.Debug("失败")
				betCount++
			}
		}

		//最少下两个个号,最多下6个号
		if betCount > 6 || betCount < 2 {
			return false
		}

		//只有一注
		order.SingleBetNum = 1

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 25: //连码连码四全中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 4 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 4)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败3")
			return false
		}
	case 26, 27: //连码连码三全中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 3 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 3)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败3")
			return false
		}

	case 28, 29, 30: //连码连码二全中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 2 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 2)
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败3")
			return false
		}

	case 31: //连肖连肖二肖连（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 2 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 2)
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

	case 32: //连肖连肖三肖连（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 3 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 3)
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
	case 33: //连肖连肖四肖连（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 4 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 4)
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
	case 34: //连肖连肖五肖连（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 5 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 5)
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
	case 35: //连尾连尾二尾碰（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_35(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 2 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 2)
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
	case 36: //连尾连尾三尾碰（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_35(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 3 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 3)
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
	case 37: //连尾连尾四尾碰（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_35(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 4 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 4)
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
	case 38: //连尾连尾五尾碰（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_35(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 5 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 5)
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
	case 39, 61: //自选不中自选不中五不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 5 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 5)
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
	case 40, 62: //自选不中自选不中六不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 6 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 6)
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
	case 41, 63: //自选不中自选不中七不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 7 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 7)
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
	case 42, 64: //自选不中自选不中八不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 11 || l < 8 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 8)
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
	case 43, 65: //自选不中自选不中九不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 12 || l < 9 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 9)
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
	case 44, 66: //自选不中自选不中十不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 13 || l < 10 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 10)
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
	case 45: //自选不中自选不中十一不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 13 || l < 11 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 11)
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
	case 46: //自选不中自选不中十二不中
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_8(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 14 || l < 12 {
			return false
		}

		//6分析订单注数
		singleBetNum := utils.AnalysisCombination(l, 12)
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
	case 47, 48, 49: //生肖生肖十二肖（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 12 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
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

	case 50: //生肖生肖总肖（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_50(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 6 || l < 1 {
			return false
		}

		//6分析订单注数
		singleBetNum := l
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

	case 51: //合肖合肖中（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 11 || l < 1 {
			return false
		}

		//6分析订单注数
		order.SingleBetNum = 1

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 52: //合肖合肖不中（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_31(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 1 {
			return false
		}

		//6分析订单注数
		order.SingleBetNum = 1

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["1"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败")
			return false
		}
	case 53: //色波色波三色波（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_53(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 3 || l < 1 {
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
	case 54, 55: //色波色波半波（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_54(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 12 || l < 1 {
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

	case 56: //色波色波七色波（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_56(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 4 || l < 1 {
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
	case 57: //尾数尾数头尾数（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_57(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 15 || l < 1 {
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
	case 58: //尾数尾数正特尾数（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_58(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 10 || l < 1 {
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
	case 59: //七码五行七码五行七码（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_59(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 16 || l < 1 {
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
	case 60: //七码五行七码五行五行（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_60(order.BetNums)
		if !ok {
			return false
		}

		l := len(array)

		if l > 5 || l < 1 {
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
func (o *HK6) PaserNormalBetNum_0(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于7
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *HK6) PaserNormalBetNum_7(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于7
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *HK6) PaserNormalBetNum_8(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于7
		if i < 1 || i > 49 {
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
func (o *HK6) PaserNormalBetNum_11(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于7
		if i < 0 || i > 13 {
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
func (o *HK6) PaserNormalBetNum_18(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于7
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
func (o *HK6) PaserNormalBetNum_31(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *HK6) PaserNormalBetNum_50(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *HK6) PaserNormalBetNum_53(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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
func (o *HK6) PaserNormalBetNum_54(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *HK6) PaserNormalBetNum_56(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *HK6) PaserNormalBetNum_57(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

		if i < 0 || i > 14 {
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
func (o *HK6) PaserNormalBetNum_58(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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
func (o *HK6) PaserNormalBetNum_59(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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
func (o *HK6) PaserNormalBetNum_60(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

		if i < 0 || i > 4 {
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
func (o *HK6) PaserNormalBetNum_35(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {

		i, err := strconv.Atoi(v)

		if err != nil {
			return false, nil
		}

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
func (o *HK6) PaserNormalBetNum(betNum string) (bool, []int) {
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
func (o *HK6) PaserNormalBetNumFor26(betNum string) (bool, []int) {
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

//解析下注号码,得到注数一维数组(这个由大小单双的判定调用 应为大小单双每个数字和其他的不一样 只能 是 0 或 1)
func (o *HK6) PaserNormalBetNumForBigSmall(betNum string) (bool, []int) {
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
func (o *HK6) PaserTwoDigitBetNum(betNum string) (bool, [][]int) {
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
func (o *HK6) PaserThreeDigitBetNum(betNum string) (bool, [][]int) {
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
func (o *HK6) PaserFiveDigitBetNumAllowSpaces(betNum string) (bool, [][]int) {
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

//解析6位下注号码,得到注数六维数组(用于有6位选择数字的情况,例如:定位胆5)
func (o *HK6) PaserSixDigitBetNumAllowSpaces(betNum string) (bool, [][]int) {
	array := strings.Split(betNum, ";")
	if len(array) != 6 {
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

			if i < 0 || i > 12 {
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
func (o *HK6) CheckRepeatAndEmptyString(array []string) bool {
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
func (o *HK6) CheckRepeatInt(array []int) bool {
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
func (o *HK6) SettlementOrders(orders []gb.Order, openCode string) {
	l := len(orders)
	//没有订单
	if l < 1 {
		return
	}
	//金额流水数组
	BalanceRecourds := []BalanceRecordMgr.BalanceRecord{}
	for i := 0; i < l; i++ {
		//这里有个大bug 如果是已经开奖的订单,这里还是会生成一个结算
		if !o.settlementOrder(&orders[i], utils.PaserOpenCodeToArray(openCode)) {
			continue
		}

		//更新订单,上线稳定后,改为批量更新订单,并且钱要用整形,以分为单位
		//结算结果要保留两位小数 4舍5入
		s := fmt.Sprintf("%.2f", orders[i].Settlement)
		orders[i].Settlement, _ = strconv.ParseFloat(s, 64)

		ss := fmt.Sprintf("%.2f", orders[i].RebateAmount)
		orders[i].RebateAmount, _ = strconv.ParseFloat(ss, 64)

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
func (o *HK6) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}

	//转下注号码为数组
	switch order.BetType {
	case 0: //两面两面正一（多赔率）
		o.WinningAndLose_0(order, openCode)
	case 1: //两面两面正二（多赔率）
		o.WinningAndLose_1(order, openCode)
	case 2: //两面两面正三（多赔率）
		o.WinningAndLose_2(order, openCode)
	case 3: //两面两面正四（多赔率）
		o.WinningAndLose_3(order, openCode)
	case 4: //两面两面正五（多赔率）
		o.WinningAndLose_4(order, openCode)
	case 5: //两面两面正六（多赔率）
		o.WinningAndLose_5(order, openCode)
	case 6: //两面两面特码（多赔率）
		o.WinningAndLose_6(order, openCode)
	case 7: //两面两面总和（多赔率）
		o.WinningAndLose_7(order, openCode)
	case 8: //正码正码正码
		o.WinningAndLose_8(order, openCode)
	case 9: //特码AB特码AB特码A
		o.WinningAndLose_9(order, openCode)
	case 10: //特码AB特码AB特码B（赔率是根据最高返利率来定）
		o.WinningAndLose_10(order, openCode)
	case 11: //特码AB特码AB其他（多赔率）
		o.WinningAndLose_11(order, openCode)
	case 12: //正码特正码特正特一
		o.WinningAndLose_12(order, openCode)
	case 13: //正码特正码特正特二
		o.WinningAndLose_13(order, openCode)
	case 14: //正码特正码特正特三
		o.WinningAndLose_14(order, openCode)
	case 15: //正码特正码特正特四
		o.WinningAndLose_15(order, openCode)
	case 16: //正码特正码特正特五
		o.WinningAndLose_16(order, openCode)
	case 17: //正码特正码特正特六
		o.WinningAndLose_17(order, openCode)
	case 18: //正码~正码~正码一（多赔率）
		o.WinningAndLose_18(order, openCode)
	case 19: //正码~正码~正码二（多赔率）
		o.WinningAndLose_19(order, openCode)
	case 20: //正码~正码~正码三（多赔率）
		o.WinningAndLose_20(order, openCode)
	case 21: //正码~正码~正码四（多赔率）
		o.WinningAndLose_21(order, openCode)
	case 22: //正码~正码~正码五（多赔率）
		o.WinningAndLose_22(order, openCode)
	case 23: //正码~正码~正码六（多赔率）
		o.WinningAndLose_23(order, openCode)
	case 24: //正码过关正码过关正码过关（多赔率中嵌套多赔率）
		o.WinningAndLose_24(order, openCode)
	case 25: //连码连码四全中
		o.WinningAndLose_25(order, openCode)
	case 26: //连码连码三全中
		o.WinningAndLose_26(order, openCode)
	case 27: //连码连码三中二
		o.WinningAndLose_27(order, openCode)
	case 28: //连码连码二全中
		o.WinningAndLose_28(order, openCode)
	case 29: //连码连码二中特
		o.WinningAndLose_29(order, openCode)
	case 30: //连码连码特串
		o.WinningAndLose_30(order, openCode)
	case 31: //连肖连肖二肖连（多赔率）
		o.WinningAndLose_31(order, openCode)
	case 32: //连肖连肖三肖连（多赔率）
		o.WinningAndLose_32(order, openCode)
	case 33: //连肖连肖四肖连（多赔率）
		o.WinningAndLose_33(order, openCode)
	case 34: //连肖连肖五肖连（多赔率）
		o.WinningAndLose_34(order, openCode)
	case 35: //连尾连尾二尾碰（多赔率）
		o.WinningAndLose_35(order, openCode)
	case 36: //连尾连尾三尾碰（多赔率）
		o.WinningAndLose_36(order, openCode)
	case 37: //连尾连尾四尾碰（多赔率）
		o.WinningAndLose_37(order, openCode)
	case 38: //连尾连尾五尾碰（多赔率）
		o.WinningAndLose_38(order, openCode)
	case 39: //自选不中自选不中五不中
		o.WinningAndLose_39(order, openCode)
	case 40: //自选不中自选不中六不中
		o.WinningAndLose_40(order, openCode)
	case 41: //自选不中自选不中七不中
		o.WinningAndLose_41(order, openCode)
	case 42: //自选不中自选不中八不中
		o.WinningAndLose_42(order, openCode)
	case 43: //自选不中自选不中九不中
		o.WinningAndLose_43(order, openCode)
	case 44: //自选不中自选不中十不中
		o.WinningAndLose_44(order, openCode)
	case 45: //自选不中自选不中十一不中
		o.WinningAndLose_45(order, openCode)
	case 46: //自选不中自选不中十二不中
		o.WinningAndLose_46(order, openCode)
	case 47: //生肖生肖十二肖（多赔率）
		o.WinningAndLose_47(order, openCode)
	case 48: //生肖生肖正肖（多赔率）
		o.WinningAndLose_48(order, openCode)
	case 49: //生肖生肖一肖（多赔率）
		o.WinningAndLose_49(order, openCode)
	case 50: //生肖生肖总肖（多赔率）
		o.WinningAndLose_50(order, openCode)
	case 51: //合肖合肖中（多赔率）
		o.WinningAndLose_51(order, openCode)
	case 52: //合肖合肖不中（多赔率）
		o.WinningAndLose_52(order, openCode)
	case 53: //色波色波三色波（多赔率）
		o.WinningAndLose_53(order, openCode)
	case 54: //色波色波半波（多赔率）
		o.WinningAndLose_54(order, openCode)
	case 55: //色波色波半半波（多赔率）
		o.WinningAndLose_55(order, openCode)
	case 56: //色波色波七色波（多赔率）
		o.WinningAndLose_56(order, openCode)
	case 57: //尾数尾数头尾数（多赔率）
		o.WinningAndLose_57(order, openCode)
	case 58: //尾数尾数正特尾数（多赔率）
		o.WinningAndLose_58(order, openCode)
	case 59: //七码五行七码五行七码（多赔率）
		o.WinningAndLose_59(order, openCode)
	case 60: //七码五行七码五行五行（多赔率）
		o.WinningAndLose_60(order, openCode)
	case 61: //中一中一五中一
		o.WinningAndLose_61(order, openCode)
	case 62: //中一中一六中一
		o.WinningAndLose_62(order, openCode)
	case 63: //中一中一七中一
		o.WinningAndLose_63(order, openCode)
	case 64: //中一中一八中一
		o.WinningAndLose_64(order, openCode)
	case 65: //中一中一九中一
		o.WinningAndLose_65(order, openCode)
	case 66: //中一中一十中一
		o.WinningAndLose_66(order, openCode)
	default:
		beego.Error("没有这个下注类型")
		return false
	}

	return true
}

//判断是否中奖 (0 两面两面正一（多赔率）)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[0] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[0] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[0]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[0]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (1  两面两面正二（多赔率）（多赔率）)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_1(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[1] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[1] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[1]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[1]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (2  两面两面正三（多赔率)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[2] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[2] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[2]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[2]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (3  两面两面正四（多赔率)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_3(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[3] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[3] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[3]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[3]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (4  两面两面正五（多赔率)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[4] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[4] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[4]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[4]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (5  两面两面正六（多赔率)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_5(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[5] == 49 { //和
		//算中一注
		winningBetNum = 1

		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[5] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[5]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[5]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (6  两面两面特码（多赔率)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_6(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[6] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [4]int
		if openCode[6] < 25 {
			resultCode[0] = 1
		} else {
			resultCode[0] = 0
		}

		if openCode[6]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[6]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (7  两面两面总和)（多赔率(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_7(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_0(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//先得出7个号码总和
	sum := o.openCode[0] + o.openCode[1] + o.openCode[2] + o.openCode[3] + o.openCode[4] + o.openCode[5] + o.openCode[6]

	//判断开奖结果能中的号码
	var resultCode [2]int

	var winningBetNum = 0

	if sum > 174 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	if sum%2 == 0 {
		resultCode[1] = 3
	} else {
		resultCode[1] = 2
	}

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
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
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (8  正码正码正码)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_8(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode []int
	resultCode = append(resultCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningBetNum = 0

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
				winningBetNum++
			}
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (9  特码AB特码AB特码A)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_9(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[6] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (10  特码AB特码AB特码B（赔率是根据最高返利率来定）)(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_10(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[6] {
			winningBetNum++
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
	} else {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (11  特码AB特码AB其他（多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_11(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_11(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	if openCode[6] == 49 { //和
		//算中一注
		winningBetNum = 1
		//返还这个订单所有金额
		order.Settlement = order.OrderAmount

	} else {
		//判断开奖结果能中的号码
		var resultCode [6]int
		if openCode[6] < 25 {
			resultCode[0] = 1
			if openCode[6]%2 == 0 {
				resultCode[5] = 13
			} else {
				resultCode[5] = 11
			}
		} else {
			resultCode[0] = 0
			if openCode[6]%2 == 0 {
				resultCode[5] = 12
			} else {
				resultCode[5] = 10
			}
		}

		if openCode[6]%2 == 0 {
			resultCode[1] = 3
		} else {
			resultCode[1] = 2
		}

		var t = openCode[6]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[2] = 5
		} else {
			resultCode[2] = 4
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[3] = 7
		} else {
			resultCode[3] = 6
		}
		//判断尾大尾小
		if t1 < 5 {
			resultCode[4] = 9
		} else {
			resultCode[4] = 8
		}

		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}

		//总结算
		order.Settlement = ret
	}

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (12 )(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_12(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[0] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (13 )(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_13(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[1] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (14 )(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_14(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[2] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (15 )(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_15(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[3] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (16 )(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_16(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[4] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (17 )(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_17(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0

	for _, v := range betNumbers {
		if v == o.openCode[5] {
			winningBetNum++
			break
		}
	}

	if winningBetNum == 0 {
		order.Settlement = 0
	} else {
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

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (18  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_18(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_18(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	var rebateCount int
	var colorBet []int

	//判断开奖结果能中的号码
	var resultCode [6]int
	if openCode[0] < 25 {
		resultCode[0] = 1 //小
	} else {
		resultCode[0] = 0 //大
	}

	if openCode[0]%2 == 0 {
		resultCode[1] = 3 //双
	} else {
		resultCode[1] = 2 //单
	}

	var t = openCode[0]
	t0 := t / 10
	t1 := t % 10

	sam := t0 + t1

	//判断和大 和小
	if sam < 7 {
		resultCode[2] = 5 //和小
	} else {
		resultCode[2] = 4 //和大
	}
	//判断和单 和双
	if sam%2 == 0 {
		resultCode[3] = 7 //和双
	} else {
		resultCode[3] = 6 //和单
	}

	//判断尾大尾小
	if t1 < 5 {
		resultCode[4] = 9 //尾小
	} else {
		resultCode[4] = 8 //尾大
	}

	//查看开奖号码色波
	c := color[openCode[0]]
	if c == 0 {
		resultCode[5] = 10
	} else if c == 1 {
		resultCode[5] = 11
	} else {
		resultCode[5] = 12
	}

	if openCode[0] == 49 { //和
		//计算要反几注钱
		for _, v := range betNumbers {
			if v < 10 {
				rebateCount++
			} else {
				colorBet = append(colorBet, v)
			}
		}

		ret += order.SingleBetAmount * float64(rebateCount)

		for _, i := range colorBet {
			if resultCode[5] == i {
				winningBetNum++
				oddsCode := strconv.Itoa(i + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}

	} else {
		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum + rebateCount
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (19  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_19(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_18(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	var rebateCount int
	var colorBet []int

	//判断开奖结果能中的号码
	var resultCode [6]int
	if openCode[1] < 25 {
		resultCode[0] = 1 //小
	} else {
		resultCode[0] = 0 //大
	}

	if openCode[1]%2 == 0 {
		resultCode[1] = 3 //双
	} else {
		resultCode[1] = 2 //单
	}

	var t = openCode[1]
	t0 := t / 10
	t1 := t % 10

	sam := t0 + t1

	//判断和大 和小
	if sam < 7 {
		resultCode[2] = 5 //和小
	} else {
		resultCode[2] = 4 //和大
	}
	//判断和单 和双
	if sam%2 == 0 {
		resultCode[3] = 7 //和双
	} else {
		resultCode[3] = 6 //和单
	}

	//判断尾大尾小
	if t1 < 5 {
		resultCode[4] = 9 //尾小
	} else {
		resultCode[4] = 8 //尾大
	}

	//查看开奖号码色波
	c := color[openCode[1]]
	if c == 0 {
		resultCode[5] = 10
	} else if c == 1 {
		resultCode[5] = 11
	} else {
		resultCode[5] = 12
	}

	if openCode[1] == 49 { //和
		//计算要反几注钱
		for _, v := range betNumbers {
			if v < 10 {
				rebateCount++
			} else {
				colorBet = append(colorBet, v)
			}
		}

		ret += order.SingleBetAmount * float64(rebateCount)

		for _, i := range colorBet {
			if resultCode[5] == i {
				winningBetNum++
				oddsCode := strconv.Itoa(i + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}

	} else {
		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum + rebateCount
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (20  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_20(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_18(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	var rebateCount int
	var colorBet []int

	//判断开奖结果能中的号码
	var resultCode [6]int
	if openCode[2] < 25 {
		resultCode[0] = 1 //小
	} else {
		resultCode[0] = 0 //大
	}

	if openCode[2]%2 == 0 {
		resultCode[1] = 3 //双
	} else {
		resultCode[1] = 2 //单
	}

	var t = openCode[2]
	t0 := t / 10
	t1 := t % 10

	sam := t0 + t1

	//判断和大 和小
	if sam < 7 {
		resultCode[2] = 5 //和小
	} else {
		resultCode[2] = 4 //和大
	}
	//判断和单 和双
	if sam%2 == 0 {
		resultCode[3] = 7 //和双
	} else {
		resultCode[3] = 6 //和单
	}

	//判断尾大尾小
	if t1 < 5 {
		resultCode[4] = 9 //尾小
	} else {
		resultCode[4] = 8 //尾大
	}

	//查看开奖号码色波
	c := color[openCode[2]]
	if c == 0 {
		resultCode[5] = 10
	} else if c == 1 {
		resultCode[5] = 11
	} else {
		resultCode[5] = 12
	}

	if openCode[2] == 49 { //和
		//计算要反几注钱
		for _, v := range betNumbers {
			if v < 10 {
				rebateCount++
			} else {
				colorBet = append(colorBet, v)
			}
		}

		ret += order.SingleBetAmount * float64(rebateCount)

		for _, i := range colorBet {
			if resultCode[5] == i {
				winningBetNum++
				oddsCode := strconv.Itoa(i + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}

	} else {
		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum + rebateCount
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (21  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_21(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_18(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	var rebateCount int
	var colorBet []int

	//判断开奖结果能中的号码
	var resultCode [6]int
	if openCode[3] < 25 {
		resultCode[0] = 1 //小
	} else {
		resultCode[0] = 0 //大
	}

	if openCode[3]%2 == 0 {
		resultCode[1] = 3 //双
	} else {
		resultCode[1] = 2 //单
	}

	var t = openCode[3]
	t0 := t / 10
	t1 := t % 10

	sam := t0 + t1

	//判断和大 和小
	if sam < 7 {
		resultCode[2] = 5 //和小
	} else {
		resultCode[2] = 4 //和大
	}
	//判断和单 和双
	if sam%2 == 0 {
		resultCode[3] = 7 //和双
	} else {
		resultCode[3] = 6 //和单
	}

	//判断尾大尾小
	if t1 < 5 {
		resultCode[4] = 9 //尾小
	} else {
		resultCode[4] = 8 //尾大
	}

	//查看开奖号码色波
	c := color[openCode[3]]
	if c == 0 {
		resultCode[5] = 10
	} else if c == 1 {
		resultCode[5] = 11
	} else {
		resultCode[5] = 12
	}

	if openCode[3] == 49 { //和
		//计算要反几注钱
		for _, v := range betNumbers {
			if v < 10 {
				rebateCount++
			} else {
				colorBet = append(colorBet, v)
			}
		}

		ret += order.SingleBetAmount * float64(rebateCount)

		for _, i := range colorBet {
			if resultCode[5] == i {
				winningBetNum++
				oddsCode := strconv.Itoa(i + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}

	} else {
		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum + rebateCount
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (22  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_22(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_18(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	var rebateCount int
	var colorBet []int

	//判断开奖结果能中的号码
	var resultCode [6]int
	if openCode[4] < 25 {
		resultCode[0] = 1 //小
	} else {
		resultCode[0] = 0 //大
	}

	if openCode[4]%2 == 0 {
		resultCode[1] = 3 //双
	} else {
		resultCode[1] = 2 //单
	}

	var t = openCode[4]
	t0 := t / 10
	t1 := t % 10

	sam := t0 + t1

	//判断和大 和小
	if sam < 7 {
		resultCode[2] = 5 //和小
	} else {
		resultCode[2] = 4 //和大
	}
	//判断和单 和双
	if sam%2 == 0 {
		resultCode[3] = 7 //和双
	} else {
		resultCode[3] = 6 //和单
	}

	//判断尾大尾小
	if t1 < 5 {
		resultCode[4] = 9 //尾小
	} else {
		resultCode[4] = 8 //尾大
	}

	//查看开奖号码色波
	c := color[openCode[4]]
	if c == 0 {
		resultCode[5] = 10
	} else if c == 1 {
		resultCode[5] = 11
	} else {
		resultCode[5] = 12
	}

	if openCode[4] == 49 { //和
		//计算要反几注钱
		for _, v := range betNumbers {
			if v < 10 {
				rebateCount++
			} else {
				colorBet = append(colorBet, v)
			}
		}

		ret += order.SingleBetAmount * float64(rebateCount)

		for _, i := range colorBet {
			if resultCode[5] == i {
				winningBetNum++
				oddsCode := strconv.Itoa(i + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}

	} else {
		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum + rebateCount
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (23  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_23(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_18(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	var rebateCount int
	var colorBet []int

	//判断开奖结果能中的号码
	var resultCode [6]int
	if openCode[5] < 25 {
		resultCode[0] = 1 //小
	} else {
		resultCode[0] = 0 //大
	}

	if openCode[5]%2 == 0 {
		resultCode[1] = 3 //双
	} else {
		resultCode[1] = 2 //单
	}

	var t = openCode[5]
	t0 := t / 10
	t1 := t % 10

	sam := t0 + t1

	//判断和大 和小
	if sam < 7 {
		resultCode[2] = 5 //和小
	} else {
		resultCode[2] = 4 //和大
	}
	//判断和单 和双
	if sam%2 == 0 {
		resultCode[3] = 7 //和双
	} else {
		resultCode[3] = 6 //和单
	}

	//判断尾大尾小
	if t1 < 5 {
		resultCode[4] = 9 //尾小
	} else {
		resultCode[4] = 8 //尾大
	}

	//查看开奖号码色波
	c := color[openCode[5]]
	if c == 0 {
		resultCode[5] = 10
	} else if c == 1 {
		resultCode[5] = 11
	} else {
		resultCode[5] = 12
	}

	if openCode[5] == 49 { //和
		//计算要反几注钱
		for _, v := range betNumbers {
			if v < 10 {
				rebateCount++
			} else {
				colorBet = append(colorBet, v)
			}
		}

		ret += order.SingleBetAmount * float64(rebateCount)

		for _, i := range colorBet {
			if resultCode[5] == i {
				winningBetNum++
				oddsCode := strconv.Itoa(i + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}

	} else {
		for _, v := range betNumbers {
			for _, i := range resultCode {
				if v == i {
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
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum + rebateCount
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (24  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_24(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserSixDigitBetNumAllowSpaces(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	//一共买了几个号
	var betCount int

	//判断开奖结果能中的号码
	var resultCode [6][6]int
	for i := 0; i < 6; i++ {
		if openCode[i] < 25 {
			resultCode[i][0] = 1 //小
		} else {
			resultCode[i][0] = 0 //大
		}

		if openCode[i]%2 == 0 {
			resultCode[i][1] = 3 //双
		} else {
			resultCode[i][1] = 2 //单
		}

		var t = openCode[i]
		t0 := t / 10
		t1 := t % 10

		sam := t0 + t1

		//判断和大 和小
		if sam < 7 {
			resultCode[i][2] = 5 //和小
		} else {
			resultCode[i][2] = 4 //和大
		}
		//判断和单 和双
		if sam%2 == 0 {
			resultCode[i][3] = 7 //和双
		} else {
			resultCode[i][3] = 6 //和单
		}

		//判断尾大尾小
		if t1 < 5 {
			resultCode[i][4] = 9 //尾小
		} else {
			resultCode[i][4] = 8 //尾大
		}

		//查看开奖号码色波
		c := color[openCode[i]]
		if c == 0 {
			resultCode[i][5] = 10
		} else if c == 1 {
			resultCode[i][5] = 11
		} else {
			resultCode[i][5] = 12
		}
	}

	//赔率数组
	var oddsArray []float64
	//判断是否中奖
	l := len(betNumbers)
	for i := 0; i < l; i++ {
		if betNumbers[i] != nil {
			betCount++
			if openCode[i] == 49 { //和 绿
				if betNumbers[i][0] < 10 { //下注号码小于10 开合就返还单注金额 算中一注
					winningBetNum++
					//ret += order.SingleBetAmount
					oddsArray = append(oddsArray, 1)
				} else if betNumbers[i][0] == 11 { //如果下注是绿波 那么算他中奖
					winningBetNum++
					//找出当前这个中奖号的赔率
					oddsCode := strconv.Itoa(13*i + betNumbers[i][0] + 1) //计算赔率号, 位数 * 13 + 买的号码 + 1
					if v, ok := order.Odds[oddsCode]; ok {
						//ret += v * order.SingleBetAmount * 1
						oddsArray = append(oddsArray, v)
					} else {
						beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
						return
					}
				} else { //既没有买10以下 ,又没有买绿波 不中奖
					break
				}
			} else { //如果不是开和
				var isWin = false
				for _, v := range resultCode[i] {
					if v == betNumbers[i][0] {
						isWin = true
						winningBetNum++
						oddsCode := strconv.Itoa(13*i + betNumbers[i][0] + 1) //计算赔率号, 位数 * 13 + 买的号码 + 1
						if v, ok := order.Odds[oddsCode]; ok {
							//ret += v * order.SingleBetAmount * 1
							oddsArray = append(oddsArray, v)
						} else {
							beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
							return
						}
						break
					}
				}
				if isWin == false {
					break
				}
			}
		}
	}

	//下注号码数 和 中奖号码数一样 那么就算中奖
	if betCount == winningBetNum {
		var oddsSum float64 = 1
		for _, v := range oddsArray {
			oddsSum *= v
		}

		ret += order.SingleBetAmount * oddsSum

		//总结算
		order.Settlement = ret
		//中奖注数(这个玩法只中一注)
		order.WinningBetNum = 1
	} else {
		//总结算
		order.Settlement = 0
		//中奖注数
		order.WinningBetNum = 0
	}

	//订单以结算
	order.Status = 1
}

//判断是否中奖 (25  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_25(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 4 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 4)
		//反水
		order.RebateAmount = 0
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

//判断是否中奖 (26  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_26(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 3 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 3)
		//反水
		order.RebateAmount = 0
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

//判断是否中奖 (27  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_27(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = 0
	var winningBetNum_2 = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 2 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//中三个数的组合
		winningBetNum = utils.AnalysisCombination(winningNumbers, 3)
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}

		//中两个数的组合
		winningBetNum_2 := utils.AnalysisCombination(winningNumbers, 2) * (len(betNumbers) - winningNumbers)

		//反水
		order.RebateAmount = 0
		//计算中了多少钱
		if v, ok := order.Odds["2"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum_2)
			//更新order
			order.Settlement += ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	//中的注数
	order.WinningBetNum = winningBetNum + winningBetNum_2
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (28  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_28(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 2 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 2)
		//反水
		order.RebateAmount = 0
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

//判断是否中奖 (29  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_29(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningNumbers int

	//第一种中奖方式 找出前6个号码买中几个
	for _, v := range betNumbers {
		for _, i := range tempOpenCode {
			if v == i {
				winningNumbers++
			}
		}
	}
	//第二种中奖方式 查看有没有买中特码
	isWinSpecial := false
	for _, v := range betNumbers {
		if v == o.openCode[6] {
			isWinSpecial = true
			break
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 1 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//判断正吗
		winningBetNum = utils.AnalysisCombination(winningNumbers, 2)
		//计算中了多少钱
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}

		if isWinSpecial == true {
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * float64(winningNumbers)
				//更新order
				order.Settlement += ret
			} else {
				beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
			}
		}
		//反水
		order.RebateAmount = 0
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (30  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_30(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var tempOpenCode []int
	tempOpenCode = append(tempOpenCode, openCode[0], openCode[1], openCode[2], openCode[3], openCode[4], openCode[5])

	var winningNumbers int

	// 查看有没有买中特码
	isWinSpecial := false
	for _, v := range betNumbers {
		if v == o.openCode[6] {
			isWinSpecial = true
			break
		}
	}

	if isWinSpecial == true {
		for _, v := range betNumbers {
			for _, i := range tempOpenCode {
				if v == i {
					winningNumbers++
				}
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 1 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = winningNumbers
		//反水
		order.RebateAmount = 0
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

//判断是否中奖 (31  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_31(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [7]int
	for i := 0; i < 7; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int

	var recordWinningNum []int
	for _, v := range betNumbers {
		for _, i := range zodiacCode {
			if v == i {
				winningNumbers++
				recordWinningNum = append(recordWinningNum, v)
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 2 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 2)
		//反水
		order.RebateAmount = 0
		//判断中奖号码中有没有鸡
		var isHaveChicken = false
		for _, v := range recordWinningNum {
			if v == 9 {
				//计算中了多少钱(因为中奖号码中有鸡,所以赔率采用10号鸡的赔率)
				if v, ok := order.Odds["10"]; ok {
					ret += v * order.SingleBetAmount * float64(len(recordWinningNum)-1)
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				//计算中了多少钱(非鸡注数)
				if v, ok := order.Odds["1"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 2))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				isHaveChicken = true
				break
			}
		}

		if isHaveChicken == false {

			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (32  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_32(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [7]int
	for i := 0; i < 7; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int

	var recordWinningNum []int
	for _, v := range betNumbers {
		for _, i := range zodiacCode {
			if v == i {
				winningNumbers++
				recordWinningNum = append(recordWinningNum, v)
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 3 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 3)
		//反水
		order.RebateAmount = 0
		//判断中奖号码中有没有鸡
		var isHaveChicken = false
		for _, v := range recordWinningNum {
			if v == 9 {
				//计算中了多少钱(因为中奖号码中有鸡,所以赔率采用10号鸡的赔率)
				if v, ok := order.Odds["10"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 2))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				//计算中了多少钱
				if v, ok := order.Odds["1"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 3))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				isHaveChicken = true
				break
			}
		}

		if isHaveChicken == false {
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (33  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_33(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [7]int
	for i := 0; i < 7; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int

	var recordWinningNum []int
	for _, v := range betNumbers {
		for _, i := range zodiacCode {
			if v == i {
				winningNumbers++
				recordWinningNum = append(recordWinningNum, v)
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 4 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 4)
		//反水
		order.RebateAmount = 0
		//判断中奖号码中有没有鸡
		var isHaveChicken = false
		for _, v := range recordWinningNum {
			if v == 9 {
				//计算中了多少钱(因为中奖号码中有鸡,所以赔率采用10号鸡的赔率)
				if v, ok := order.Odds["10"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 3))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				//计算中了多少钱
				if v, ok := order.Odds["1"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 4))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				isHaveChicken = true
				break
			}
		}

		if isHaveChicken == false {
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}

		}
		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (34  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_34(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [7]int
	for i := 0; i < 7; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int

	var recordWinningNum []int
	for _, v := range betNumbers {
		for _, i := range zodiacCode {
			if v == i {
				winningNumbers++
				recordWinningNum = append(recordWinningNum, v)
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 5 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 5)
		//反水
		order.RebateAmount = 0
		//判断中奖号码中有没有鸡
		var isHaveChicken = false
		for _, v := range recordWinningNum {
			if v == 9 {
				//计算中了多少钱(因为中奖号码中有鸡,所以赔率采用10号鸡的赔率)
				if v, ok := order.Odds["10"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 4))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				//计算中了多少钱
				if v, ok := order.Odds["1"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 5))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				isHaveChicken = true
				break
			}
		}

		if isHaveChicken == false {
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}

		}
		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (35  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_35(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_35(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个开奖号码的尾数
	var tailCode [7]int
	for i := 0; i < 7; i++ {
		tailCode[i] = o.openCode[i] % 10
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tailCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 2 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 2)
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (36  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_36(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_35(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个开奖号码的尾数
	var tailCode [7]int
	for i := 0; i < 7; i++ {
		tailCode[i] = o.openCode[i] % 10
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tailCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 3 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 3)
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (37  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_37(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_35(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个开奖号码的尾数
	var tailCode [7]int
	for i := 0; i < 7; i++ {
		tailCode[i] = o.openCode[i] % 10
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range tailCode {
			if v == i {
				winningNumbers++
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers < 4 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 4)
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (38  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_38(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_35(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个开奖号码的尾数
	var tailCode [7]int
	for i := 0; i < 7; i++ {
		tailCode[i] = o.openCode[i] % 10
	}

	var winningNumbers int

	var recordWinningNum []int
	for _, v := range betNumbers {
		for _, i := range tailCode {
			if v == i {
				winningNumbers++
				recordWinningNum = append(recordWinningNum, v)
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 5 证明没有中奖
	if winningNumbers < 5 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = utils.AnalysisCombination(winningNumbers, 5)
		//反水
		order.RebateAmount = 0
		//判断中奖号码中有没有鸡
		var isHaveChicken = false
		for _, v := range recordWinningNum {
			if v == 0 {
				if v, ok := order.Odds["1"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 4))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				//计算中了多少钱
				if v, ok := order.Odds["2"]; ok {
					ret += v * order.SingleBetAmount * float64(utils.AnalysisCombination(len(recordWinningNum)-1, 5))
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}

				isHaveChicken = true
				break
			}
		}

		if isHaveChicken == false {
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * float64(winningBetNum)
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}

		}
		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (39  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_39(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 5)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (40  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_40(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 6)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (41  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_41(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 7)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (42  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_42(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 8)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (43  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_43(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 9)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (44 （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_44(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 10)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (45 （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_45(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 11)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (46 （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_46(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	//找出未中奖号码数
	resultCode := len(betNumbers) - winningNumbers

	var winningBetNum = utils.AnalysisCombination(resultCode, 12)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (47  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_47(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode int

	zodiacCode = zodiac[o.openCode[6]]

	var winningNumbers int

	for _, v := range betNumbers {
		if v == zodiacCode {
			winningNumbers++
			oddsCode := strconv.Itoa(v + 1)
			if v, ok := order.Odds[oddsCode]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
			break
		}
	}

	//反水
	order.RebateAmount = 0

	//更新order
	order.Settlement = ret

	//中的注数
	order.WinningBetNum = winningNumbers
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (48  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_48(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [6]int
	for i := 0; i < 6; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range zodiacCode {
			if v == i {
				winningNumbers++
				oddsCode := strconv.Itoa(v + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
			}
		}
	}

	//反水
	order.RebateAmount = 0

	//更新order
	order.Settlement = ret

	//中的注数
	order.WinningBetNum = winningNumbers
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (49  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_49(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [7]int
	for i := 0; i < 7; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range zodiacCode {
			if v == i {
				winningNumbers++
				oddsCode := strconv.Itoa(v + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}
	}

	var winningBetNum = 0
	//如果中奖数组< 4 证明没有中奖
	if winningNumbers == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		winningBetNum = winningNumbers
		//反水
		order.RebateAmount = 0

		//更新order
		order.Settlement = ret
	}
	//中的注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (50  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_50(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_50(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//找出7个中奖号码的属性
	var zodiacCode [7]int
	for i := 0; i < 7; i++ {
		zodiacCode[i] = zodiac[o.openCode[i]]
	}

	var winningNumbers int = 0
	//不同生肖的个数
	var zodiacCount = make(map[int]int)
	for _, v := range zodiacCode {
		if _, ok := zodiacCount[v]; !ok {
			zodiacCount[v] = 1
		}
	}

	//可中奖结果
	var resultCode [2]int
	l := len(zodiacCount)
	if l == 2 || l == 3 || l == 4 {
		resultCode[0] = 0
	} else if l == 5 {
		resultCode[0] = 1
	} else if l == 6 {
		resultCode[0] = 2
	} else if l == 7 {
		resultCode[0] = 3
	}

	if l%2 == 1 {
		resultCode[1] = 4
	} else {
		resultCode[1] = 5
	}

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
				winningNumbers++
				//计算中了多少钱
				oddsCode := strconv.Itoa(v + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
			}
		}
	}

	//反水
	order.RebateAmount = 0

	//更新order
	order.Settlement = ret

	//中的注数
	order.WinningBetNum = winningNumbers
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (51  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_51(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	l := len(betNumbers)

	//找出7个中奖号码的属性
	var zodiacCode int

	zodiacCode = zodiac[o.openCode[6]]

	var winningNumbers int

	for _, v := range betNumbers {
		if v == zodiacCode {
			winningNumbers++
			oddsCode := strconv.Itoa(l)
			if v, ok := order.Odds[oddsCode]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
			break
		}
	}

	//反水
	order.RebateAmount = 0

	//更新order
	order.Settlement = ret

	//中的注数
	order.WinningBetNum = winningNumbers
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (52  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_52(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_31(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	l := len(betNumbers)

	//找出7个中奖号码的属性
	var zodiacCode int

	zodiacCode = zodiac[o.openCode[6]]

	var winningNumbers int

	var isWinning = false
	for _, v := range betNumbers {
		if v == zodiacCode {
			isWinning = true
			break
		}
	}

	if isWinning == false {
		winningNumbers++
		oddsCode := strconv.Itoa(l)
		if v, ok := order.Odds[oddsCode]; ok {
			ret += v * order.SingleBetAmount * 1
		} else {
			beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
			return
		}
	}

	//反水
	order.RebateAmount = 0

	//更新order
	order.Settlement = ret

	//中的注数
	order.WinningBetNum = winningNumbers
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (53 （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_53(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_53(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	c := color[o.openCode[6]]

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

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (54  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_54(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_54(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	c := color[o.openCode[6]]

	var resultCode [2]int

	if c == 0 { //红
		if o.openCode[6] > 24 {
			resultCode[0] = 0 //红大
		} else {
			resultCode[0] = 1 //红小
		}

		if o.openCode[6]%2 == 1 {
			resultCode[1] = 2 //红单
		} else {
			resultCode[1] = 3 //红双
		}
	} else if c == 1 { //蓝
		if o.openCode[6] > 24 {
			resultCode[0] = 4 //蓝大
		} else {
			resultCode[0] = 5 //蓝小
		}

		if o.openCode[6]%2 == 1 {
			resultCode[1] = 6 //蓝单
		} else {
			resultCode[1] = 7 //蓝双
		}
	} else { //绿
		if o.openCode[6] > 24 {
			resultCode[0] = 8 //绿大
		} else {
			resultCode[0] = 9 //绿小
		}

		if o.openCode[6]%2 == 1 {
			resultCode[1] = 10 //绿单
		} else {
			resultCode[1] = 11 //绿双
		}
	}

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
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
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (55  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_55(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_54(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	c := color[o.openCode[6]]

	var resultCode int

	if c == 0 { //红
		if o.openCode[6] > 24 {
			if o.openCode[6]%2 == 1 {
				resultCode = 0 //红大单
			} else {
				resultCode = 1 //红大双
			}
		} else {
			if o.openCode[6]%2 == 1 {
				resultCode = 2 //红小单
			} else {
				resultCode = 3 //红小双
			}
		}

	} else if c == 1 { //蓝
		if o.openCode[6] > 24 {
			if o.openCode[6]%2 == 1 {
				resultCode = 4 //蓝大单
			} else {
				resultCode = 5 //蓝大双
			}
		} else {
			if o.openCode[6]%2 == 1 {
				resultCode = 6 //蓝小单
			} else {
				resultCode = 7 //蓝小双
			}
		}
	} else { //绿
		if o.openCode[6] > 24 {
			if o.openCode[6]%2 == 1 {
				resultCode = 8 //绿大单
			} else {
				resultCode = 9 //绿大双
			}
		} else {
			if o.openCode[6]%2 == 1 {
				resultCode = 10 //绿小单
			} else {
				resultCode = 11 //绿小双
			}
		}
	}

	for _, v := range betNumbers {
		if v == resultCode {
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

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (56  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_56(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_56(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	//找出7个号码的颜色
	var colorArray [7]int
	for i := 0; i < 7; i++ {
		colorArray[i] = color[o.openCode[i]]
	}

	//0 1 2 红蓝绿
	var colorCount [3]float64

	for i := 0; i < 7; i++ {
		c := colorArray[i]
		if i == 6 {
			colorCount[c] += 1.5
		} else {
			colorCount[c] += +1
		}
	}

	var result int

	//判断哪个颜色最多
	if colorCount[0] == colorCount[1] { //第一个颜色 等于 第二个颜色
		if colorCount[2] > colorCount[0] {
			result = 2
		} else {
			result = 3
		}
	} else if colorCount[0] == colorCount[2] { //第一个颜色 等于 第三个颜色
		if colorCount[1] > colorCount[0] {
			result = 1
		} else {
			result = 3
		}
	} else if colorCount[1] == colorCount[2] { //第二个颜色 等于 第三个颜色
		if colorCount[0] > colorCount[1] {
			result = 0
		} else {
			result = 3
		}
	} else { //谁大谁赢
		if colorCount[0] > colorCount[1] {
			if colorCount[0] > colorCount[2] {
				result = 0
			} else {
				result = 1
			}
		} else {
			if colorCount[1] > colorCount[2] {
				result = 1
			} else {
				result = 2
			}
		}
	}

	for _, v := range betNumbers {
		if v == result {
			winningBetNum++
			oddsCode := strconv.Itoa(v + 1)
			if v, ok := order.Odds[oddsCode]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
			break
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (57  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_57(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_57(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	var resultCode [2]int
	resultCode[0] = o.openCode[6] / 10
	resultCode[1] = o.openCode[6]%10 + 5

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
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
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (58  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_58(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_58(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	var resultCode [6]int

	for i := 0; i < 6; i++ {
		resultCode[i] = o.openCode[i] % 10
	}

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
				winningBetNum++
				oddsCode := strconv.Itoa(v + 1)
				if v, ok := order.Odds[oddsCode]; ok {
					ret += v * order.SingleBetAmount * 1
				} else {
					beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
					return
				}
				break
			}
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (59  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_59(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_59(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0
	//大
	var big = 0

	//单
	var odd = 0

	for _, v := range o.openCode {
		if v > 24 {
			big++
		}

		if v%2 == 1 {
			odd++
		}
	}

	var resultCode [2]int

	resultCode[0] = big
	resultCode[1] = odd + 8

	for _, v := range betNumbers {
		for _, i := range resultCode {
			if v == i {
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
	}

	//总结算
	order.Settlement = ret
	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (60  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_60(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_60(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断中几注
	var winningBetNum = 0

	var resultCode int = fiveLine[o.openCode[6]]

	for _, v := range betNumbers {
		if v == resultCode {
			winningBetNum++
			oddsCode := strconv.Itoa(v + 1)
			if v, ok := order.Odds[oddsCode]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
			break
		}
	}

	//总结算
	order.Settlement = ret

	//中奖注数
	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (61  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_61(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = winningNumbers * utils.AnalysisCombination(len(betNumbers)-winningNumbers, 4)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (62  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_62(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = winningNumbers * utils.AnalysisCombination(len(betNumbers)-winningNumbers, 5)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (63  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_63(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = winningNumbers * utils.AnalysisCombination(len(betNumbers)-winningNumbers, 6)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (64  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_64(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = winningNumbers * utils.AnalysisCombination(len(betNumbers)-winningNumbers, 7)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (65  （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_65(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = winningNumbers * utils.AnalysisCombination(len(betNumbers)-winningNumbers, 8)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//判断是否中奖 (66 （多赔率）(返回用户输赢情况)(会有多注中奖情况)
func (o *HK6) WinningAndLose_66(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_8(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers int

	for _, v := range betNumbers {
		for _, i := range o.openCode {
			if v == i {
				winningNumbers++
			}
		}
	}

	var winningBetNum = winningNumbers * utils.AnalysisCombination(len(betNumbers)-winningNumbers, 9)
	//如果中奖数组< 4 证明没有中奖
	if winningBetNum == 0 {
		//反水
		order.RebateAmount = 0
		//更新order
		order.Settlement = 0
	} else {
		//反水
		order.RebateAmount = 0

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

//计算下期香港6合彩的开奖日期和时间
//香港6合彩 开奖时间为 二，四，六，晚上9点30分
func (o *HK6) CalculationNextOpenTime(currentOpenTime time.Time) time.Time {
	//计算两天以后的日期
	nextOpenTime := utils.DateAfterTheDay(currentOpenTime, 2)

	//如果是周一再延后一天
	if nextOpenTime.Weekday() == time.Monday {
		nextOpenTime = utils.DateAfterTheDay(nextOpenTime, 1)
	}

	//开奖时间
	d, _ := time.ParseDuration("21h30m")

	return nextOpenTime.Add(d)
}

//计算下期香港六合彩开奖期数
//由于六合彩期数的特殊性,不能一直+1 跨年时 必须从1开始 年数也要变化
func (o *HK6) CalculationNextExpect(nextOpenTime time.Time, currentExpect int) int {
	//如果下一起的开奖时间和现在是同一年,那么期数+1就可以了
	if nextOpenTime.Year() == utils.GetNowUTC8Time().Year() {
		return currentExpect + 1
	} else if nextOpenTime.Year()-utils.GetNowUTC8Time().Year() == 1 {
		//如果 下一起年数 - 现在年数 == 1 说明下一期就过年了 ,下一期期数 = 下一起开奖的年数 * 100 + 1 比如 2018001
		beego.Info("!!!!!!!!!!!!!!!!!!!!!!  跨年啦!~!~!~~ 希望有看到这条消息的那一天啊~!~!~!~!~  !!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
		return nextOpenTime.Year()*1000 + 1

	} else {
		beego.Error("六合彩 期数计算错误 快检查!~!~~ \n")
		//ctrl.Instance().ChangeLotteryStatus(o.gameTag, gb.LotteryStatus_Maintain)
		return 0
	}
	return 0
}
