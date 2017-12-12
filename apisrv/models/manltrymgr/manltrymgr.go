package manltrymgr

import (
	"apisrv/models/ltrymgr"
	"sync"

	"github.com/astaxie/beego"
)

//手动开奖彩票管理
//手动开奖系统
//所有需要手动开奖的彩种 全部要放在这里面
//需要手动开奖的彩种,或者是api数据不完善的彩种,所有的启动都需要手动来实现
//也就是说,这些特殊的彩种,要再后台控制开启,关闭,后台控制开奖
type ManlLtryMgr struct {
	//LotteriesInterfaceMap map[string]ManualLotteryInterface
	ltryifMap map[string]ManlLtryif //所有彩票对象	k为Lotter Name (彩票名字)
	lock      sync.RWMutex
}

var sInstance *ManlLtryMgr
var once sync.Once

// Instance 获取dbmgr单例，Once 线程安全
func Instance() *ManlLtryMgr {
	once.Do(func() {
		sInstance = &ManlLtryMgr{}
		sInstance.init()
	})
	return sInstance
}

//初始化
func (o *ManlLtryMgr) init() {
	o.ltryifMap = make(map[string]ManlLtryif)
}

//启动彩种(初始化彩种)注意:这里启动的是手动开奖系统的彩票,目前只有六合彩
//之所以不能在初始化的时候启动,是应为必须手动输入一些参数才能正确启动
//参数为: 1.彩票名字 2.当前期数(已经开奖的这期) 3.当前这期开奖结果 4.当前这期开奖时间 5.下期期数 6.下期开奖时间
func (o *ManlLtryMgr) InitLottery(ltryName string, currentExpect int, currentOpenCode string, currentOpenTime string, nextExpect int, nextOpenTime string) int {
	switch ltryName {
	case "HK6": //香港六合彩(手动)
		hk6 := ltrymgr.Instance().HK6
		if hk6 != nil {
			ec := hk6.InitLtry(currentExpect, currentOpenCode, nextExpect, nextOpenTime)
			if ec != 0 {
				return ec
			}
			o.lock.Lock()
			defer o.lock.Unlock()
			o.ltryifMap[ltryName] = hk6
		} else {
			beego.Error("------------------------- HK6 Already Initialized ! -------------------------")
		}

	default:
		beego.Error("------------------------- Init Manual Lottery Error : ", ltryName, " Not Have Lottery Type ! -------------------------")
		return
	}
}

//手动开奖
func (o ManlLtryMgr) ManualLottery(ltryName string, expect int, openCode string, openTime string, nextExpect int, nextOpenTime string) {
	if ltry, ok := o.ltryifMap[ltryName]; ok {
		ltry.StartLtry(expect, openCode, openTime, nextExpect, nextOpenTime)
	} else {
		beego.Error("------------------------- Manual Lottery Error : ", ltryName, " Not Have Lottery Type ! -------------------------")
	}
}
