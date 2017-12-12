package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

//得到上一期期数
func GetLastExpect(gameName string, currentExpect int) (int, error) {
	switch gameName {
	//------------------------------------------------------------ EX5 ------------------------------------------------------------
	case "EX5_JiangXi": //江西11选5 为84期
		if currentExpect%100 == 1 {
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
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
			t := DateBeforeTheDay(time.Now(), 1)
			i, err := strconv.Atoi(t.Format(TF_D))
			if err != nil {
				beego.Debug(err)
				return 0, err
			}
			return i*1000 + 71, nil
		} else {
			currentExpect--
			return currentExpect, nil
		}
	case "BJ28": //北京28 和北京pk10一样 是一直累加的
		currentExpect--
		return currentExpect, nil
	case "PK10_BeiJing": //北京PK10 注意北京pk10 期数是6位数的自增数字,那么当数字达到999999的时候 要根据官方的情况来修改代码
		currentExpect--
		return currentExpect, nil
	case "PL3": //注意!!!!!!!!!!!!! 排列3 情况特殊,如果是跨年,那么应该是上一年的天数作为期数尾部,但是又有几天不开,所以这个等过年时来处理
		currentExpect--
		return currentExpect, nil
	}

	return 0, errors.New("计算上期彩票开奖期数错误,未找到彩票类型 : " + gameName)
}

//将api获得的gameTag转换成我看得顺眼的形式,不服吗?我就是看小写开头不顺眼
func ConvertGameName(gameName string) string {
	switch gameName {
	case "jx11x5":
		return "EX5_JiangXi"
	case "sd11x5":
		return "EX5_ShanDong"
	case "sh11x5":
		return "EX5_ShangHai"
	case "bj11x5":
		return "EX5_BeiJing"
	case "fj11x5":
		return "EX5_FuJian"
	case "hlj11x5":
		return "EX5_HeiLongJiang"
	case "js11x5":
		return "EX5_JiangSu"

	case "gxk3":
		return "K3_GuangXi"
	case "jlk3":
		return "K3_JiLin"
	case "ahk3":
		return "K3_AnHui"
	case "bjk3":
		return "K3_BeiJing"
	case "fjk3":
		return "K3_FuJian"
	case "hebk3":
		return "K3_HeBei"
	case "shk3":
		return "K3_ShangHai"
	case "jsk3":
		return "K3_JiangSu"

	case "cqssc":
		return "SSC_ChongQing"
	case "tjssc":
		return "SSC_TianJin"
	case "xjssc":
		return "SSC_XinJiang"

	case "bjpk10":
		return "PK10_BeiJing"

	case "pl3":
		return "PL3"

	case "hk6":
		return "HK6"
	case "bjkl8":
		return "BJ28"
	case "dlt":
		return "DLT"
	case "fc3d":
		return "FC3D"
	case "qxc":
		return "QXC"
	case "ssq":
		return "SSQ"
	case "cakeno":
		return "CAKENO"

	default:
		beego.Error("Convert game name error , There is no lottery type!")
		return ""
	}
}
