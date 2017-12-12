package PK10

import (
	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/Order"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"
	"common/utils"

	"calculsrv/models/acmgr"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

var sumArray = [17]int{2, 2, 4, 4, 6, 6, 8, 8, 10, 8, 8, 6, 6, 4, 4, 2, 2}

func (o *PK10) AnalyticalOrder(order *gb.Order, accountInfo *acmgr.AccountInfo) bool {
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
	case 0: //冠军冠军冠军
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大10 应为前一玩法最少选择一个数字,最多选择10个数字,选择了几个就是几注
		if l < 1 || l > 10 {
			//beego.Debug("失败")
			return false
		}

		//分析订单注数(前一玩法,选择了几个数就是几注)
		singleBetNum := l
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 1: //冠亚军冠亚军冠亚军
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserTwoDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		//每个数组数量不得大于11个元素或小于1个元素，应为前二直选每一位 最多只能选11个数字 1 - 11 或最少选择1个数字
		for _, v := range array {
			if len(v) > 10 || len(v) < 1 {
				//beego.Debug("失败")
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
		//分析订单注数 (第一组元素个数 * 第二组元素个数) - 两组重复数字次数
		singleBetNum := len(array[0])*len(array[1]) - repateNum
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 2: //冠亚季军冠亚季军冠亚季军
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserThreeDigitBetNum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		//每个数组数量不得大于11个元素或小于1个元素，应为前三组选每一位 最多只能选11个数字 1 - 11 或最少选择1个数字
		for _, v := range array {
			if len(v) > 10 || len(v) < 1 {
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
				for _, j := range array[2] {
					if v == i && v == j {
						repeatT123 += 1
					}
				}
			}
		}

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

		//公式
		l1 := len(array[0])
		l2 := len(array[1])
		l3 := len(array[2])
		//(l1 * l2 * l3) - repeatT12 * l3 -repeatT23 * l1 - repeatT13 * l2 + repeatT123 * 2
		singleBetNum := l1*l2*l3 - repeatT12*l3 - repeatT23*l1 - repeatT13*l2 + repeatT123*2
		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}

	case 3, 4: //定位胆定位胆第一名~第五名,定位胆定位胆第六名~第十名
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserFiveDigitBetNum(order.BetNums)
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
			if l > 10 {
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

	case 19, 20, 21, 22, 23: //两面两面冠军 - 第五名
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_5(order.BetNums)
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
			if l > 1 {
				//beego.Debug("失败")
				return false
			}
			singleBetNum += l
		}

		//5个数组的长度不得小于1 大于50
		if tl < 1 || tl > 3 {
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
			beego.Debug("失败")
			return false
		}

	case 16: //冠亚军和冠亚军和和值
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumForSum(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)
		if l < 1 || l > 17 {
			return false
		}

		var singleBetNum = 0
		for _, v := range array {
			singleBetNum += sumArray[v-3]
		}

		order.SingleBetNum = singleBetNum

		//7计算赔率(注意玩家下注时候的赔率就是结算时的赔率,就算服务器赔率有更改) 赔率计算 赔率 = 赔率 - (赔率 * 返水)
		if odds, ok := order.Odds["1"]; ok {
			order.Odds["1"] = odds - odds*order.Rebate
		} else {
			beego.Debug("失败1")
			return false
		}
	case 17: //冠亚军和冠亚军和和值大小（多赔率）
		//4分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumFor17(order.BetNums)
		if !ok {
			//beego.Debug("失败")
			return false
		}

		l := len(array)

		if l != 1 {
			beego.Debug("失败")
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
	case 18: //冠亚军和冠亚军和和值单双（多赔率）
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNumForBigSmall(order.BetNums)
		if !ok {
			return false
		}
		l := len(array)
		//数组数量不得小于1 大2 应为大小玩法,只能选择大小两个数 单为0, 双为1
		if l < 0 || l > 1 {
			//beego.Debug("失败")
			return false
		}

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

	case 24, 25, 26, 27, 28: //两面两面第六名
		//分析订单下注数字是否正确(如果正确返回解析后的int数组)
		ok, array := o.PaserNormalBetNum_24(order.BetNums)
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
			if l > 1 {
				//beego.Debug("失败")
				return false
			}
			singleBetNum += l
		}

		//5个数组的长度不得小于1 大于50
		if tl < 1 || tl > 2 {
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
func (o *PK10) PaserNormalBetNum(betNum string) (bool, []int) {
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
		if i < 1 || i > 10 {
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
func (o *PK10) PaserNormalBetNumForBigSmall(betNum string) (bool, []int) {
	//分割下注字符
	array := strings.Split(betNum, ",")

	//验证1:将字符数组转换为int数组(如果有不合法的情况直接false)
	var arrayInt []int
	for _, v := range array {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false, nil
		}
		//验证2:每位数字不能小于0和大于1(PK10龙虎玩法)
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)(16)
func (o *PK10) PaserNormalBetNumFor17(betNum string) (bool, []int) {
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

//解析下注号码,得到注数一维数组(没有;号 分割的这种下注数字)(16)
func (o *PK10) PaserNormalBetNumForTwoSum(betNum string) (bool, []int) {
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
		if i < 3 || i > 19 {
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

//解析下注号码,(冠亚军和冠亚军和和值)
func (o *PK10) PaserNormalBetNumForSum(betNum string) (bool, []int) {
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
		if i < 3 || i > 19 {
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
func (o *PK10) PaserTwoDigitBetNum(betNum string) (bool, [][]int) {
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
			if i < 1 || i > 10 {
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
func (o *PK10) PaserThreeDigitBetNum(betNum string) (bool, [][]int) {
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
			//每位数字不能小于1和大于10(PK10玩法)
			if i < 1 || i > 10 {
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
//注意 由于定位胆
func (o *PK10) PaserFiveDigitBetNum(betNum string) (bool, [][]int) {
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
			if i < 1 || i > 10 {
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

func (o *PK10) PaserNormalBetNum_5(betNum string) (bool, [][]int) {
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
			if i < 0 || i > 1 {
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

func (o *PK10) PaserNormalBetNum_24(betNum string) (bool, [][]int) {
	//分割下注位 由于定位胆的特殊性这里要做特殊处理, ;号分割以后会出现空的字段 只要总分个数为5 那么这种空字段的情况是可以出现的
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
			if j == "" {
				continue
			}
			i, err := strconv.Atoi(j)
			if err != nil {
				return false, nil
			}
			//每位数字不能小于1和大于10(PK10玩法)
			if i < 0 || i > 1 {
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

//检查是否有重复数字符(int数组)
func (o *PK10) CheckRepeatInt(array []int) bool {
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
func (o *PK10) SettlementOrders(orders []gb.Order, openCode string) {
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
func (o *PK10) settlementOrder(order *gb.Order, openCode []int) bool {
	if order.Status == 1 {
		beego.Debug("严重错误 : 订单状态错误 ! 正常情况下是不可能出现这个错误的,没有开奖的订单不可能已经结算,所以如果出现这个错误代表服务器被黑")
		return false
	}
	//转下注号码为数组
	switch order.BetType {
	case 0: //冠军冠军冠军
		o.WinningAndLose_0(order, openCode)
	case 1: //冠亚军冠亚军冠亚军
		o.WinningAndLose_1(order, openCode)
	case 2: //冠亚季军冠亚季军冠亚季军
		o.WinningAndLose_2(order, openCode)
	case 3: //定位胆定位胆第一名~第五名
		o.WinningAndLose_3(order, openCode)
	case 4: //定位胆定位胆第六名~第十名
		o.WinningAndLose_4(order, openCode)
	case 16: //冠亚军和冠亚军和和值
		o.WinningAndLose_16(order, openCode)
	case 17: //冠亚军和冠亚军和和值大小（多赔率）
		o.WinningAndLose_17(order, openCode)
	case 18: //冠亚军和冠亚军和和值单双（多赔率）
		o.WinningAndLose_18(order, openCode)
	case 19: //两面两面冠军
		o.WinningAndLose_19(order, openCode)
	case 20: //两面两面亚军
		o.WinningAndLose_20(order, openCode)
	case 21: //两面两面季军
		o.WinningAndLose_21(order, openCode)
	case 22: //两面两面第四名
		o.WinningAndLose_22(order, openCode)
	case 23: //两面两面第五名
		o.WinningAndLose_23(order, openCode)
	case 24: //两面两面第六名
		o.WinningAndLose_24(order, openCode)
	case 25: //两面两面第七名
		o.WinningAndLose_25(order, openCode)
	case 26: //两面两面第八名
		o.WinningAndLose_26(order, openCode)
	case 27: //两面两面第九名
		o.WinningAndLose_27(order, openCode)
	case 28: //两面两面第十名
		o.WinningAndLose_28(order, openCode)

	default:
		beego.Debug("失败")
		return false
	}

	return true
}

//判断是否中奖 (0 冠军冠军冠军)
func (o *PK10) WinningAndLose_0(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningNumbers = 0
	var winningBetNum = 0
	//判断开奖结果第一位有没有在下注号码中
	for _, v := range betNumbers {
		if v == openCode[0] {
			winningNumbers = 1
			break
		}
	}
	//开出的第一位号码没有在下注号码里面
	if winningNumbers == 0 {
		//一注没中 计算反水 单注金额 * 反水 * 单注数量
		ret += (order.SingleBetAmount * order.Rebate) * float64(order.SingleBetNum)
		//反水
		order.RebateAmount = ret
		//更新order
		order.Settlement = ret
	} else if winningNumbers == 1 {
		winningBetNum = 1
		//前一玩法只会有一注中奖 计算反水 单注金额 * 反水 * 单注数量 -1
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

//判断是否中奖 (1冠亚军冠亚军冠亚军)(注意 PK10前二 和 EX5 的前二直选是一模一样的)
func (o *PK10) WinningAndLose_1(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserTwoDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningNumbers = 0
	var winningBetNum = 0
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
	//前两位都相同才中奖
	if flag_1 == 1 && flag_2 == 1 {
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

//判断是否中奖 (2 冠亚季军冠亚季军冠亚季军)(注意 PK10前三 和 EX5 的前三直选是一模一样的)
func (o *PK10) WinningAndLose_2(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserThreeDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0
	var winningNumbers = 0
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
	//前两位都相同才中奖
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

//判断是否中奖 (3 定位胆定位胆第一名~第五名)
func (o *PK10) WinningAndLose_3(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var winningBetNum = 0
	//判断中奖注数(注意:每一位只会有一个数字中奖)
	//应为每一位中奖的赔率不一样,所以要设置标志位看看是哪一位中奖
	//第一名
	for _, v := range betNumbers[0] {
		if v == openCode[0] {
			winningBetNum++
			break
		}
	}

	//第二名
	for _, v := range betNumbers[1] {
		if v == openCode[1] {
			winningBetNum++
			break
		}
	}

	//第三名
	for _, v := range betNumbers[2] {
		if v == openCode[2] {
			winningBetNum++
			break
		}
	}

	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[3] {
			winningBetNum++
			break
		}
	}

	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[4] {
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

//判断是否中奖 (4 定位胆定位胆第六名~第十名)
func (o *PK10) WinningAndLose_4(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserFiveDigitBetNum(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	var winningBetNum = 0
	//判断中奖注数(注意:每一位只会有一个数字中奖)
	//应为每一位中奖的赔率不一样,所以要设置标志位看看是哪一位中奖
	//第一名
	for _, v := range betNumbers[0] {
		if v == openCode[5] {
			winningBetNum++
			break
		}
	}

	//第二名
	for _, v := range betNumbers[1] {
		if v == openCode[6] {
			winningBetNum++
			break
		}
	}

	//第三名
	for _, v := range betNumbers[2] {
		if v == openCode[7] {
			winningBetNum++
			break
		}
	}

	//第四名
	for _, v := range betNumbers[3] {
		if v == openCode[8] {
			winningBetNum++
			break
		}
	}

	//第五名
	for _, v := range betNumbers[4] {
		if v == openCode[9] {
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

//判断是否中奖 (16 冠亚军和冠亚军和和值)(返回用户输赢情况)
func (o *PK10) WinningAndLose_16(order *gb.Order, openCode []int) {
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

//判断是否中奖 (17 冠亚军和冠亚军和和值大小（多赔率）)(返回用户输赢情况)
func (o *PK10) WinningAndLose_17(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumFor17(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	retSum := openCode[0] + openCode[1]

	var winningBetNum = 0

	if retSum == 11 {
		winningBetNum = 1
		ret += order.SingleBetAmount
	} else if retSum > 11 {
		if betNumbers[0] == 0 {
			winningBetNum = 1
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
	} else if retSum < 11 {
		if betNumbers[0] == 1 {
			winningBetNum = 1
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
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

//判断是否中奖 (18  冠亚军和冠亚军和和值单双（多赔率）(返回用户输赢情况)
func (o *PK10) WinningAndLose_18(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNumForBigSmall(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}
	retSum := openCode[0] + openCode[1]

	var winningBetNum = 0

	if retSum == 11 {
		winningBetNum = 1
		ret += order.SingleBetAmount
	} else if retSum%10%2 == 1 {
		if betNumbers[0] == 0 {
			winningBetNum = 1
			//计算中了多少钱
			if v, ok := order.Odds["1"]; ok {
				ret += v * order.SingleBetAmount * 1
			} else {
				beego.Emergency("------------------------- 严重错误 : 订单结算失败 -------------------------")
				return
			}
		}
	} else if retSum%10%2 == 0 {
		if betNumbers[0] == 1 {
			winningBetNum = 1
			//计算中了多少钱
			if v, ok := order.Odds["2"]; ok {
				ret += v * order.SingleBetAmount * 1
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

//判断是否中奖 (19 两面两面冠军)
func (o *PK10) WinningAndLose_19(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_5(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1
	if o.openCode[0] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[0]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	//龙 0 ,虎 1
	if o.openCode[0] > o.openCode[9] {
		resultCode[2] = 0
	} else {
		resultCode[2] = 1
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	if len(betNumbers[2]) != 0 {
		if betNumbers[2][0] == resultCode[2] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (20 两面两面亚军)
func (o *PK10) WinningAndLose_20(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_5(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1
	if o.openCode[1] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[1]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	//龙 0 ,虎 1
	if o.openCode[1] > o.openCode[8] {
		resultCode[2] = 0
	} else {
		resultCode[2] = 1
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	if len(betNumbers[2]) != 0 {
		if betNumbers[2][0] == resultCode[2] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (21 两面两面季军)
func (o *PK10) WinningAndLose_21(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_5(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1
	if o.openCode[2] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[2]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	//龙 0 ,虎 1
	if o.openCode[2] > o.openCode[7] {
		resultCode[2] = 0
	} else {
		resultCode[2] = 1
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	if len(betNumbers[2]) != 0 {
		if betNumbers[2][0] == resultCode[2] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (22 两面两面第四名)
func (o *PK10) WinningAndLose_22(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_5(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1
	if o.openCode[3] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[3]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	//龙 0 ,虎 1
	if o.openCode[3] > o.openCode[6] {
		resultCode[2] = 0
	} else {
		resultCode[2] = 1
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	if len(betNumbers[2]) != 0 {
		if betNumbers[2][0] == resultCode[2] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (23 两面两面第五名)
func (o *PK10) WinningAndLose_23(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_5(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [3]int

	//大 0 , 小 1
	if o.openCode[4] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[4]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	//龙 0 ,虎 1
	if o.openCode[4] > o.openCode[5] {
		resultCode[2] = 0
	} else {
		resultCode[2] = 1
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	if len(betNumbers[2]) != 0 {
		if betNumbers[2][0] == resultCode[2] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (24 两面两面第六名)
func (o *PK10) WinningAndLose_24(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_24(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1
	if o.openCode[5] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[5]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (25 两面两面第七名)
func (o *PK10) WinningAndLose_25(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_24(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1
	if o.openCode[6] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[6]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (26 两面两面第八名)
func (o *PK10) WinningAndLose_26(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_24(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1
	if o.openCode[7] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[7]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (27 两面两面第九名)
func (o *PK10) WinningAndLose_27(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_24(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1
	if o.openCode[8] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[8]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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

//判断是否中奖 (28 两面两面第十名)
func (o *PK10) WinningAndLose_28(order *gb.Order, openCode []int) {
	//最后结果 输赢多少钱
	var ret float64
	ok, betNumbers := o.PaserNormalBetNum_24(order.BetNums)
	if !ok {
		beego.Debug("------------------------- 严重错误 : 订单结算失败 -------------------------")
		return
	}

	var resultCode [2]int

	//大 0 , 小 1
	if o.openCode[9] > 5 {
		resultCode[0] = 0
	} else {
		resultCode[0] = 1
	}

	//单 0 , 双 1
	if o.openCode[9]%2 == 0 {
		resultCode[1] = 1
	} else {
		resultCode[1] = 0
	}

	var winningBetNum = 0

	if len(betNumbers[0]) != 0 {
		if betNumbers[0][0] == resultCode[0] {
			winningBetNum++
		}
	}

	if len(betNumbers[1]) != 0 {
		if betNumbers[1][0] == resultCode[1] {
			winningBetNum++
		}
	}

	//开出的第一位号码没有在下注号码里面
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
