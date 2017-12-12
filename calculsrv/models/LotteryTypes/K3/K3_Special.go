package K3

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

//分析订单(下注)
func (o *K3) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {
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
	case 0: //二不同号二不同号标准 (最多中3注)
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//5数组数量不得小于2 大于6 应为二不同号玩法 最少选择2个 最多选择6个号码
		if l > 6 || l < 2 {
			//beego.Debug("失败")
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
	case 1: //二同号二同号单选(目前 同号位只能选择一个,也就是说 有几个单号就有几注)
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			return false
		}
		l1 := len(array[0])
		l2 := len(array[1])

		//二同号按照趣彩来,同号位只能选择一个
		if l1 != 1 {
			//beego.Debug("失败")
			return false
		}

		if l2 < 1 || l2 > 5 {
			return false
		}

		//找出两组重复数对数(不能有重复数字)
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					return false
				}
			}
		}

		//singleBetNum := l1*l2 - repateNum  这个算法是下注不排重,计算排重的情况下用的
		//6分析订单注数(目前在同号位只能选择一个数字的情况下,l2选择了几个数 就是几注)
		order.SingleBetNum = l2

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 2: //二同号二同号复选
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l1 := len(array)
		//5二同号复选
		if l1 < 1 || l1 > 6 {
			return false
		}

		//6分析订单注数(二同号复选注数 ,就是l1的个数)
		singleBetNum := l1
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 3: //三不同号三不同号标准
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于3 大于6 应为三不同号玩法 最少选择3个 最多选择6个号码
		if l < 3 || l > 6 {
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

	case 4: //三同号三同号单选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大于6 应为三不同号玩法 最少选择1个 最多选择6个号码
		if l < 1 || l > 6 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(三同号,数组个数就是下注个数,最少1注,最多6注)
		singleBetNum := l
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 5: //三同号三同号通选
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量必须为6 应为三同号通选 就是所有都选
		if l != 6 {
			//beego.Debug("失败")
			return false
		}
		//分析订单注数(三同号通选,只算1注)
		singleBetNum := 1
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 6: //连号连号三连号
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量必须等于4 应为三连号玩法 客户端传来的是1234  数字1 代表123  2代表234 类推
		if l != 4 {
			return false
		}

		//分析订单注数(三连号只有一注)
		singleBetNum := 1
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 7: //和值和值和值（多赔率）
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumSpecial(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大于16 应为和值最少选择一个数,最多选择16个数
		if l < 1 || l > 16 {
			return false
		}

		//分析订单注数(和值,选择数就是注数)
		singleBetNum := l
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if _, ok := order.Odds["3"]; ok {
			for k, v := range order.Odds {
				order.Odds[k] = v - v*order.Rebate
			}
		} else {
			beego.Debug("失败1")
			return false
		}

	case 8: //二不同号二不同号胆拖
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			return false
		}
		l1 := len(array[0])
		l2 := len(array[1])
		//二不同号胆拖胆码只能选择一个
		if l1 != 1 {
			return false
		}
		//第二个数组最少选择一个,最多选择5个
		if l2 < 1 || l2 > 5 {
			return false
		}

		//第二组和第一组不能有数字重复
		for _, v := range array[1] {
			if v == array[0][0] {
				return false
			}
		}
		//分析订单注数(目前在胆码只能选择一个的情况下,l2选择了几个数 就是几注)
		singleBetNum := l2
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 9: //三不同号三不同号胆拖
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			return false
		}
		l1 := len(array[0])
		l2 := len(array[1])
		//三不同号胆拖胆码最少选择一个,最多选择2个
		if l1 < 1 || l1 > 2 {
			return false
		}
		//第二个数组拖码的最少选择数是根据第一个来变化的
		if l2 < 3-l1 || l2 > 6-l1 {
			return false
		}

		//第二组和第一组不能有数字重复
		for _, v := range array[0] {
			for _, i := range array[1] {
				if v == i {
					return false
				}
			}
		}

		//分析订单注数(3不同胆拖,如果l1 为1 那么l2 两两排列组合为注数, 如果l1为2,那么l2就是注数)
		var singleBetNum = 0
		if l1 == 1 {
			singleBetNum = utils.AnalysisCombination(l2, 2)
		} else if l1 == 2 {
			singleBetNum = l2
		} else {
			return false
		}
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 10: //连号连号二连号
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumSpecial10(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大于16 应为和值最少选择一个数,最多选择16个数
		if l < 1 || l > 15 {
			return false
		}

		//分析订单注数(和值,选择数就是注数)
		singleBetNum := l
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 11, 12: //和值和值大小, 和值和值单双
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
	case 13: //独胆独胆独胆
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}
		l := len(array)
		//5数组数量必须是 大于1 或小于11 应为前一直选 可以选择1 - 11 个号码
		if l > 6 || l < 1 {
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
	case 14: //两面百位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor14(order.BetNums)
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

	case 15: //两面十位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor15(order.BetNums)
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

	case 16, 17: //两面个位
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor16(order.BetNums)
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)
func (o *K3) PaserNormalBetNum(betNum string) (bool, []int) {
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
		if i < 1 || i > 6 {
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
func (o *K3) PaserNormalBetNumForBigSmall(betNum string) (bool, []int) {
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

func (o *K3) PaserNormalBetNumSpecial(betNum string) (bool, []int) {
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
		if i < 3 || i > 18 {
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

func (o *K3) PaserNormalBetNumSpecial10(betNum string) (bool, []int) {
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
		if i < 12 || i > 56 {
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
func (o *K3) PaserTwoDigitBetNum(betNum string) (bool, [][]int) {
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
			//每位数字不能小于1和大于6
			if i < 1 || i > 6 {
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
func (o *K3) PaserNormalBetNumFor14(betNum string) (bool, []int) {
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
func (o *K3) PaserNormalBetNumFor15(betNum string) (bool, []int) {
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
func (o *K3) PaserNormalBetNumFor16(betNum string) (bool, []int) {
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
func (o *K3) CheckRepeatInt(array []int) bool {
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
func (o *K3) SettlementOrders(orders []gb.Order, openCode string) {
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
func (o *K3) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}

	//转下注号码为数组
	switch order.BetType {
	case 0: //二不同号二不同号标准
		o.WinningAndLose_0(order, openCode)
	case 1: //二同号二同号单选
		o.WinningAndLose_1(order, openCode)
	case 2: //二同号二同号复选
		o.WinningAndLose_2(order, openCode)
	case 3: //三不同号三不同号标准
		o.WinningAndLose_3(order, openCode)
	case 4: //三同号三同号单选
		o.WinningAndLose_4(order, openCode)
	case 5: //三同号三同号通选
		o.WinningAndLose_5(order, openCode)
	case 6: //连号连号三连号
		o.WinningAndLose_6(order, openCode)
	case 7: //和值和值和值（多赔率）
		o.WinningAndLose_7(order, openCode)
	case 8: //二不同号二不同号胆拖
		o.WinningAndLose_8(order, openCode)
	case 9: //三不同号三不同号胆拖
		o.WinningAndLose_9(order, openCode)
	case 10: //连号连号二连号
		o.WinningAndLose_10(order, openCode)
	case 11: //和值和值大小
		o.WinningAndLose_11(order, openCode)
	case 12: //和值和值单双
		o.WinningAndLose_12(order, openCode)
	case 13: //独胆独胆独胆
		o.WinningAndLose_13(order, openCode)
	case 14:
		o.WinningAndLose_14(order, openCode)
	case 15:
		o.WinningAndLose_15(order, openCode)
	case 16:
		o.WinningAndLose_16(order, openCode)
	case 17:
		o.WinningAndLose_17(order, openCode)

	default:
		beego.Debug("失败")
		return false
	}

	return true
}

//判断是否中奖 (0 二不同号二不同号标准)(返回用户输赢情况)
func (o *K3) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers []int
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		for _, i := range openCode {
			if v == i {
				winningNumbers = append(winningNumbers, v)
				break
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
	} else if l == 3 {
		//中了3注
		winningBetNum = 3
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -3
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

//判断是否中奖 (1 二同号二同号单选)(返回用户输赢情况)
func (o *K3) WinningAndLose_1(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	//判断下注号码中有几个开出(先看有没有开出两个相同的号码)

	var winningBetNum = 0
	if openCode[0] == openCode[1] {
		var sameNumber = 0
		for _, v := range betNumbers[0] {
			if openCode[0] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers[1] {
				if v == openCode[2] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[0] == openCode[2] {
		var sameNumber = 0
		for _, v := range betNumbers[0] {
			if openCode[0] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers[1] {
				if v == openCode[1] {
					winningBetNum = 1
					break
				}
			}
		}
	} else if openCode[1] == openCode[2] {
		var sameNumber = 0
		for _, v := range betNumbers[0] {
			if openCode[1] == v {
				sameNumber = 1
				break
			}
		}
		if sameNumber == 1 {
			for _, v := range betNumbers[1] {
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
	order.Status = 1
}

//判断是否中奖 (2 二同号二同号复选)(返回用户输赢情况)
func (o *K3) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0

	if openCode[0] == openCode[1] {
		for _, v := range betNumbers {
			if v == openCode[0] {
				winningBetNum = 1
				break
			}
		}

	} else if openCode[0] == openCode[2] {
		for _, v := range betNumbers {
			if v == openCode[0] {
				winningBetNum = 1
				break
			}
		}
	} else if openCode[1] == openCode[2] {
		for _, v := range betNumbers {
			if v == openCode[1] {
				winningBetNum = 1
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
	} else if winningBetNum == 1 {
		//中了一注 计算反水 单注金额 * 反水 * 单注数量 -1
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
	order.Status = 1
}

//判断是否中奖 (3 三不同号三不同号标准)(返回用户输赢情况)
func (o *K3) WinningAndLose_3(order *gb.Order, openCode []int) {
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

//判断是否中奖 (4 三同号三同号单选)(返回用户输赢情况)
func (o *K3) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断下注号码中有几个开出
	if openCode[0] == openCode[1] && openCode[0] == openCode[2] {
		for _, v := range betNumbers {
			if openCode[0] == v {
				winningBetNum = 1
				break
			}
		}
	}

	//如果中奖数组<1 证明没有中奖
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
	order.Status = 1
}

//判断是否中奖 (5 三同号三同号通选)(返回用户输赢情况)
func (o *K3) WinningAndLose_5(order *gb.Order, openCode []int) {
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

//判断是否中奖 (6 连号连号三连号)(返回用户输赢情况)
func (o *K3) WinningAndLose_6(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	var winningBetNum = 0
	//判断是否开出连号
	if openCode[0] == (openCode[1]-1) && openCode[0] == (openCode[2]-2) {
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

//判断是否中奖 (7 和值和值和值（多赔率）)(返回用户输赢情况)
func (o *K3) WinningAndLose_7(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	//得到和值
	sum := openCode[0] + openCode[1] + openCode[2]
	ok, betNumbers := o.PaserNormalBetNumSpecial(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断下注号码中有几个开出
	for _, v := range betNumbers {
		if v == sum {
			winningBetNum++
		}
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
		//和值只有一注中奖
		//计算中了多少钱.判断该用那个赔率(注意这里的odds key 的转换)
		if v, ok := order.Odds[strconv.Itoa(sum)]; ok {
			ret += v * order.SingleBetAmount * 1
			//中的注数
			winningBetNum = 1
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}
	order.WinningBetNum = winningBetNum
	order.Status = 1
}

//判断是否中奖 (8 二不同号二不同号胆拖)(返回用户输赢情况)
func (o *K3) WinningAndLose_8(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	var winningNumbers = 0

	//剔除开奖号码中相同的数

	for _, v := range openCode {
		if betNumbers[0][0] == v { //胆码中
			for _, i := range openCode {
				for _, j := range betNumbers[1] {
					if i == j {
						winningNumbers++
					}
				}
			}
			break
		}
	}

	if openCode[0] == openCode[1] && openCode[0] == openCode[2] { //3个号相同直接不中奖
		winningBetNum = 0
	} else if (openCode[0] == openCode[1] || openCode[1] == openCode[2] || openCode[0] == openCode[2]) && winningNumbers == 2 {
		winningNumbers--
	}

	//如果中奖数组<1 证明没有中奖
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else {
		winningBetNum = winningNumbers
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
		order.RebateAmount = ret
		//计算中了多少钱.判断该用那个赔率
		if v, ok := order.Odds["1"]; ok {
			ret += v * order.SingleBetAmount * float64(winningBetNum)
			//更新order
			order.Settlement = ret
		} else {
			beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		}
	}

	order.WinningBetNum = winningBetNum
	order.Status = 1
}

//判断是否中奖 (9 三不同号三不同号胆拖)(返回用户输赢情况)
func (o *K3) WinningAndLose_9(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//胆码数
	l1 := len(betNumbers[0])

	var winningNumber = 0
	var flag0 = 0
	var flag1 = 0
	var flag2 = 0

	var winningBetNum = 0

	if l1 == 1 {
		for _, v := range openCode {
			if betNumbers[0][0] == v {
				flag0 = 1
				break
			}
		}

		if flag0 == 1 {
			for _, v := range betNumbers[1] {
				for _, i := range openCode {
					if v == i {
						winningNumber++
						break
					}
				}
			}
		}

		//如果中奖数组<1 证明没有中奖
		if winningNumber != 2 {
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		} else if winningNumber == 2 {
			winningBetNum = 1
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
	} else if l1 == 2 {
		for _, v := range openCode {
			if v == betNumbers[0][0] {
				flag0 = 1
				break
			}
		}

		if flag0 == 1 {
			for _, v := range openCode {
				if v == betNumbers[0][1] {
					flag1 = 1
					break
				}
			}

			if flag1 == 1 {
				for _, v := range openCode {
					for _, j := range betNumbers[1] {
						if v == j {
							flag2 = 1
							break

						}
					}
					if flag2 == 1 {
						break
					}
				}
			}
		}

		if flag0 == 1 && flag1 == 1 && flag2 == 1 {
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
			//一注没中 计算反水 单注金额 * 反水 * 单注数量
			ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
			//反水
			order.RebateAmount = ret
			//更新order
			order.Settlement = ret
		}
	}
	order.WinningBetNum = winningBetNum
	order.Status = 1
}

//判断是否中奖 (10 连号连号二连号)(返回用户输赢情况)
func (o *K3) WinningAndLose_10(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumSpecial10(order.BetNums)
	if !ok {
		//如果出现这个错误 说明问题严重要关闭彩票,后来补上
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	//将下注号码拆成3个数组

	var abc []int
	abc = append(abc, openCode[0]*10+openCode[1])
	abc = append(abc, openCode[0]*10+openCode[2])
	abc = append(abc, openCode[1]*10+openCode[2])

	var winningBetNum = 0
	for _, v := range betNumbers {
		for _, i := range abc {
			if v == i {
				winningBetNum++
				break
			}
		}
	}
	//如果中奖数组<1 证明没有中奖
	if winningBetNum == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningBetNum > 0 && winningBetNum < 4 {
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
	order.Status = 1
}

//判断是否中奖 (11 和值和值大小)(返回用户输赢情况)
func (o *K3) WinningAndLose_11(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断
	var big = 0
	var small = 0
	//三同号直接不中奖
	if !(openCode[0] == openCode[1] && openCode[0] == openCode[2]) {
		//计算5个开奖号码的和值
		retSum := openCode[0] + openCode[1] + openCode[2]

		if retSum > 10 {
			big = 1
		} else if retSum < 11 {
			small = 1
		}
	}

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
	}
	//计算反水
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//最后结算金额(中奖金额 + 反水)
	order.Settlement = ret + order.RebateAmount

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (12  和值和值单双)(返回用户输赢情况)
func (o *K3) WinningAndLose_12(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断
	var d = 0
	var s = 0
	//日吗 单双三同号又不算
	//三同号直接不中奖
	//if !(openCode[0] == openCode[1] && openCode[0] == openCode[2]) {
	//计算5个开奖号码的和值
	retSum := openCode[0] + openCode[1] + openCode[2]

	retSum = retSum % 2

	if retSum == 1 {
		d = 1
	} else if retSum == 0 {
		s = 1
	}
	//}

	if d == 1 {
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
	} else if s == 1 {
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
	}
	//计算反水
	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)
	//最后结算金额(中奖金额 + 反水)
	order.Settlement = ret + order.RebateAmount

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}

//判断是否中奖 (13 独胆独胆独胆)(返回用户输赢情况)(只会有一注中奖)
func (o *K3) WinningAndLose_13(order *gb.Order, openCode []int) {
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

//判断是否中奖 14 两面百位
func (o *K3) WinningAndLose_14(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor14(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [4]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6, 龙百 7, 虎个 8, 百个和 9
	if o.openCode[0] > 3 {
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

//判断是否中奖 15 两面十位
func (o *K3) WinningAndLose_15(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor15(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4, 虎十 5, 百十和 6
	if o.openCode[1] > 3 {
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

//判断是否中奖 16 两面个位
func (o *K3) WinningAndLose_16(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor16(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1 , 单 2, 双 3, 龙百 4
	if o.openCode[2] > 3 {
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

//判断是否中奖 17 两面和值
func (o *K3) WinningAndLose_17(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor16(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	var winningBetNum = 0
	if o.openCode[0] != o.openCode[1] || o.openCode[0] != o.openCode[2] { //开豹子直接不中奖
		sum := o.openCode[0] + o.openCode[1] + o.openCode[2]
		//大 0 , 小 1 , 单 2, 双 3,
		if sum > 10 {
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
	}

	//开出的第一位号码没有在下注号码里面

	order.RebateAmount = (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum-winningBetNum)

	order.Settlement = order.RebateAmount + ret

	order.WinningBetNum = winningBetNum
	//订单以结算
	order.Status = 1
}
