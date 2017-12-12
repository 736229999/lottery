package CommonFunc

import (
	"common/utils"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

//得到上一期期数
func GetLastExpect(gameTag string, currentExpect int) (int, error) {
	switch gameTag {
	//------------------------------------------------------------ EX5 ------------------------------------------------------------
	case "EX5_JiangXi": //江西11选5 为84期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 84, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "EX5_ShanDong": //山东11选5 为87期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 87, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "EX5_ShangHai": //上海11选5 为90期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 90, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "EX5_BeiJing": //北京11选5 为85期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 85, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "EX5_FuJian": //福建11选5 为90期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 90, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "EX5_HeiLongJiang": //黑龙江11选5 为88期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 88, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "EX5_JiangSu": //江苏11选5 为82期
		if currentExpect%100 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 80, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	//------------------------------------------------------------ K3 ------------------------------------------------------------
	case "K3_GuangXi": //广西快3为78期 注意:广西快3的期数格式为  20170531087
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 78, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "K3_JiLin": //吉林快3 为87期 ,注意:吉林快3的期数格式为  20170531087
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 87, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "K3_AnHui": //安徽快3 为80期 ,注意:安徽快3的期数格式为  20170531087
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 80, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "K3_BeiJing": //北京快3 为97期 ,注意:北京快3的期数格式为 082479 (从第一期开始累计)
		currentExpect--
		return currentExpect, nil
	case "K3_FuJian": //福建快3 为78期 ,注意:福建快3的期数格式为  20170531087
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 78, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "K3_HeBei": //福建快3 为81期 ,注意:福建快3的期数格式为  20170531087
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 81, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "K3_ShangHai": //上海快3 为82期 ,注意:福建快3的期数格式为  20170531087
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 82, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}

	//------------------------------------------------------------ SSC ------------------------------------------------------------
	case "SSC_ChongQing": //重庆时时彩 为120期 格式为2017053120
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 120, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "SSC_TianJin": //天津时时彩 为84期 格式为2017053120
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 84, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "SSC_XinJiang": //新疆时时彩 为84期 格式为2017053120
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 96, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "SSC_YunNan": //云南时时彩 为71期 格式为2017053120
		if currentExpect%1000 == 1 {
			t := utils.DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 71, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}

	case "PK10_BeiJing": //北京PK10 注意北京pk10 期数是6位数的自增数字,那么当数字达到999999的时候 要根据官方的情况来修改代码
		currentExpect--
		return currentExpect, nil

	case "BJ28": //北京PK10 注意北京pk10 期数是6位数的自增数字,那么当数字达到999999的时候 要根据官方的情况来修改代码
		currentExpect--
		return currentExpect, nil

	case "PL3": //注意!!!!!!!!!!!!! 排列3 情况特殊,如果是跨年,那么应该是上一年的天数作为期数尾部,但是又有几天不开,所以这个等过年时来处理
		currentExpect--
		return currentExpect, nil
	}

	return 0, errors.New("计算上期彩票开奖期数错误,未找到彩票类型 : " + gameTag)
}

//得到下期期数(计算下下期期数) 注意 由于六合彩的特殊性,六合彩的计算放在它的类中
func GetAfterNextExpect(gameTag string, nextExpect int) (int, error) {
	switch gameTag {
	//------------------------------------------------------------ EX5 ------------------------------------------------------------
	case "EX5_JiangXi": //江西11选5 为84期
		if nextExpect%100 == 84 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "EX5_ShanDong": //山东11选5 为87期
		if nextExpect%100 == 87 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "EX5_ShangHai": //上海11选5 为90期
		if nextExpect%100 == 90 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "EX5_BeiJing": //北京11选5 为85期
		if nextExpect%100 == 85 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "EX5_FuJian": //福建11选5 为90期
		if nextExpect%100 == 90 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "EX5_HeiLongJiang": //黑龙江11选5 为88期
		if nextExpect%100 == 88 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "EX5_JiangSu": //江苏11选5 为82期
		if nextExpect%100 == 82 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*100 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
		//------------------------------------------------------------ K3 ------------------------------------------------------------
	case "K3_GuangXi": //广西快3为78期 注意:广西快3的期数格式为  20170531087
		if nextExpect%1000 == 78 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "K3_JiLin": //吉林快3 为87期 ,注意:吉林快3的期数格式为  20170531087
		if nextExpect%1000 == 87 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "K3_AnHui": //安徽快3 为80期 ,注意:安徽快3的期数格式为  20170531087
		if nextExpect%1000 == 80 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "K3_BeiJing": //北京快3 为97期 ,注意:北京快3的期数格式为 082479 (从第一期开始累计)
		nextExpect++
		return nextExpect, nil
	case "K3_FuJian": //福建快3 为78期 ,注意:福建快3的期数格式为  20170531087
		if nextExpect%1000 == 78 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "K3_HeBei": //福建快3 为81期 ,注意:福建快3的期数格式为  20170531087
		if nextExpect%1000 == 81 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "K3_ShangHai": //上海快3 为82期 ,注意:福建快3的期数格式为  20170531087
		if nextExpect%1000 == 82 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}

	case "K3_JiangSu": //江苏快3 为82期 ,注意:江苏快3的期数格式为  20170531087
		if nextExpect%1000 == 82 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}

		//------------------------------------------------------------ SSC ------------------------------------------------------------
	case "SSC_ChongQing": //重庆时时彩 为120期 格式为2017053120
		if nextExpect%1000 == 120 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "SSC_TianJin": //天津时时彩 为84期 格式为2017053120
		if nextExpect%1000 == 84 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "SSC_XinJiang": //新疆时时彩 为84期 格式为2017053120
		if nextExpect%1000 == 96 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	case "SSC_YunNan": //云南时时彩 为71期 格式为2017053120
		if nextExpect%1000 == 71 {
			t := utils.DateAfterTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(utils.TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}

	case "PK10_BeiJing": //北京PK10 注意北京pk10 期数是6位数的自增数字,那么当数字达到999999的时候 要根据官方的情况来修改代码
		nextExpect++
		return nextExpect, nil

	case "PK10_F": //自己的急速PK10
		nextExpect++
		return nextExpect, nil

	case "BJ28": //pc蛋蛋 北京28, 期数是6位增数,还有 3个月 期数即将达到999999 要注意这时候怎么办
		nextExpect++
		return nextExpect, nil

	case "PL3":
		//得到明天是多少年
		t := utils.DateAfterTheDay(time.Now(), 1)
		//如果跨年了
		if t.Year() > time.Now().Year() {
			beego.Info("!!!!!!!!!!!!!!!!!!!!!!  跨年啦!~!~!~~ 希望有看到这条消息的那一天啊~!~!~!~!~  !!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
			return t.Year()*1000 + 1, nil
		} else {
			nextExpect++
			return nextExpect, nil
		}
	}

	return 0, errors.New("计算下下期彩票开奖期数错误,为找到彩票类型 : " + gameTag)
}

//得到下期开奖时间(参数为彩票名， 下期期数， 当期开奖时间，api 下期开奖时间）
func GetNextOpenTime(gameTag string, nextExpect int, currentOpentime time.Time, nexRequestTime time.Time) time.Time {
	switch gameTag {
	//------------------------------------------------------------ EX5 ------------------------------------------------------------
	case "EX5_JiangXi": //江西11选5 , 开彩是整点 是每10分钟开一期 ，开奖时间为 10分整点 如 22：00：00， 22：10：00 ，22：20：00
		t, err := getIntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "EX5_ShanDong": //山东11选5 ，开彩时间是 5分， 间隔还是10分
		t, err := getIntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 5, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "EX5_ShangHai":
		t, err := getIntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "EX5_BeiJing":
		t, err := getIntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "EX5_FuJian":
		t, err := getIntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "EX5_HeiLongJiang":
		t, err := getIntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 5, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "EX5_JiangSu": //江苏这个奇葩 是7分开奖
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t

		//------------------------------------------------------------ K3 ------------------------------------------------------------
	case "K3_GuangXi":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t

	case "K3_JiLin": //吉林快3 9分钟一次
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 9)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "K3_AnHui":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "K3_BeiJing":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "K3_FuJian":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "K3_HeBei":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "K3_ShangHai":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "K3_JiangSu":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t

		//------------------------------------------------------------ SSC ------------------------------------------------------------
	case "SSC_ChongQing": //重庆时时彩 1-23期 5分一期, 24-96 10分一期, 96 - 120 5分一期
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "SSC_TianJin":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "SSC_XinJiang":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "SSC_YunNan":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 10)
		if err != nil {
			beego.Debug(err)
		}
		return t

		//------------------------------------------------------------ 其他 ------------------------------------------------------------
	case "PK10_BeiJing": //北京PK10是 5分钟一期
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 5)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "PK10_F": //急速PK10 1分钟一期
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 1)
		if err != nil {
			beego.Debug(err)
		}
		return t
	case "PL3":
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 5)
		if err != nil {
			beego.Debug(err)
		}
		return t

	case "BJ28": //现在所有的彩票下期开采时间 紧紧是去掉了秒数
		t, err := getK3IntegerOpenTime(gameTag, currentOpentime, nexRequestTime, 5)
		if err != nil {
			beego.Debug(err)
		}
		return t
	}
	t := time.Time{}
	beego.Error("计算下期彩票开彩时间错误,未找到彩票类型 : " + gameTag)
	return t
}

//获取下下期截至下注时间
func GetAfterNextClosingBetTime(gameTag string, nextExpect int, nextClosingBetTime time.Time) time.Time {
	switch gameTag {
	//---------------------------------- EX5 ----------------------------------
	case "EX5_JiangXi":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "EX5_ShanDong":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "EX5_ShangHai":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "EX5_BeiJing":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "EX5_FuJian":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "EX5_HeiLongJiang":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "EX5_JiangSu":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	//---------------------------------- K3 ----------------------------------
	case "K3_GuangXi":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "K3_JiLin": //注意吉林快3的开奖间隔为9分钟
		d, _ := time.ParseDuration("9m")
		return nextClosingBetTime.Add(d)
	case "K3_AnHui":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "K3_BeiJing":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "K3_FuJian":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "K3_HeBei":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "K3_ShangHai":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "K3_JiangSu":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)

	//---------------------------------- SSC ----------------------------------
	case "SSC_ChongQing": //这个比较特殊
		var d time.Duration
		e := nextExpect % 100
		if e == 23 { //如果下一期是23期 那么下下期24期就要等485分钟才开奖
			d, _ = time.ParseDuration("485m")
		} else if e >= 24 && e < 96 {
			d, _ = time.ParseDuration("10m")
		} else if (e >= 96 && e <= 120) || (e < 23) {
			d, _ = time.ParseDuration("5m")
		}
		return nextClosingBetTime.Add(d)
	case "SSC_TianJin":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "SSC_XinJiang":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)
	case "SSC_YunNan":
		d, _ := time.ParseDuration("10m")
		return nextClosingBetTime.Add(d)

	//---------------------------------- Other ----------------------------------
	case "PK10_BeiJing":
		d, _ := time.ParseDuration("5m")
		return nextClosingBetTime.Add(d)

	case "PK10_F":
		d, _ := time.ParseDuration("1m")
		return nextClosingBetTime.Add(d)

	case "BJ28":
		d, _ := time.ParseDuration("5m")
		return nextClosingBetTime.Add(d)

	case "PL3":
		d, _ := time.ParseDuration("20h")
		return nextClosingBetTime.Add(d)
	}

	beego.Error("严重错误,没有找到彩票类型")
	return nextClosingBetTime
}

//获取整分整点
func getIntegerOpenTime(gameTag string, currentOpentime time.Time, nexRequestTime time.Time, integerNum int, interval int) (time.Time, error) {
	//取整分钟数
	nm := nexRequestTime.Minute() - (nexRequestTime.Minute() % integerNum)
	//计算下一期的整分钟开奖时间
	nt := time.Date(nexRequestTime.Year(), nexRequestTime.Month(), nexRequestTime.Day(), nexRequestTime.Hour(), nm, 0, 0, nexRequestTime.Location())
	//计算下一期的整分开奖时间是否正确（将当前这期的开奖时间取整+10分钟后的结果应该与下期取整后的时间一样）
	// cm := currentOpentime.Minute() - (currentOpentime.Minute() % integerNum)
	// ct := time.Date(currentOpentime.Year(), currentOpentime.Month(), currentOpentime.Day(), currentOpentime.Hour(), cm+interval, 0, 0, nexRequestTime.Location())

	//if ct == nt {
	return nt, nil
	// } else {
	// 	t := time.Time{}
	// 	ControlCenter.Instance().ChangeLotteryStatus(gameTag, 2)
	// 	return t, errors.New("计算整点开彩时间错误,这是一个严重错误，这说明上一期彩票可能在官方整点之前就开出，或者是API出现问题 上一期彩票的开彩时间错误，与下一期彩票开彩时间间隔大于正常彩票开彩间隔,设置彩票为维护状态")
	// }
}

//获取K3开奖整分数
func getK3IntegerOpenTime(gameTag string, currentOpentime time.Time, nexRequestTime time.Time, interval int) (time.Time, error) {
	//取整分钟数(K3 去掉秒数就OK)
	nm := nexRequestTime.Minute()
	//计算下一期的整分钟开奖时间
	nt := time.Date(nexRequestTime.Year(), nexRequestTime.Month(), nexRequestTime.Day(), nexRequestTime.Hour(), nm, 0, 0, nexRequestTime.Location())
	//计算下一期的整分开奖时间是否正确（将当前这期的开奖时间取整+10分钟后的结果应该与下期取整后的时间一样）
	//cm := currentOpentime.Minute()
	//ct := time.Date(currentOpentime.Year(), currentOpentime.Month(), currentOpentime.Day(), currentOpentime.Hour(), cm+interval, 0, 0, nexRequestTime.Location())

	//if ct == nt {
	return nt, nil
	//}
	// else {
	// 	t := time.Time{}
	// 	ControlCenter.Instance().ChangeLotteryStatus(gameTag, 2)
	// 	return t, errors.New("计算整点开彩时间错误,这是一个严重错误，这说明上一期彩票可能在官方整点之前就开出，或者是API出现问题 上一期彩票的开彩时间错误，与下一期彩票开彩时间间隔大于正常彩票开彩间隔,设置彩票为维护状态")
	// }
}
