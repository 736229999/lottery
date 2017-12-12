package ltrymgr

import (
	"apisrv/conf"
	"apisrv/models/apimgr"
	"apisrv/models/ctrl"
	"apisrv/models/ltry/hk6"
	"common/utils"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

//由于 目前所有自动开奖的彩票都是相同的类,所有没有采用接口的方式,而是同类实体
type LtryMgr struct {
	Ltrys map[string]*Ltry
	HK6   *hk6.HK6
}

var sInstance *LtryMgr
var once sync.Once

func Instance() *LtryMgr {
	once.Do(func() {
		sInstance = &LtryMgr{}
	})

	return sInstance
}

//不管是 手动开奖的彩票 ,还是统一类型的彩票都要再这里执行初始化,应为有部分数据是放在ctrl数据库中得,需要初始化才有这些信息
func (o *LtryMgr) Init() error {
	o.Ltrys = make(map[string]*Ltry)

	//由于并行化写入map会引发资源共享的问题,所以这里暂时采用单线程,以后来优化为多携程处理
	for _, v := range ctrl.Instance().Ltrys {
		switch v.Game_name {

		case "EX5_JiangXi", "EX5_ShanDong", "EX5_ShangHai", "EX5_BeiJing", "EX5_FuJian", "EX5_HeiLongJiang", "EX5_JiangSu":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "K3_GuangXi", "K3_JiLin", "K3_AnHui", "K3_BeiJing", "K3_FuJian", "K3_HeBei", "K3_ShangHai", "K3_JiangSu":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "SSC_ChongQing", "SSC_TianJin", "SSC_XinJiang":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "PK10_BeiJing":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "PL3":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "HK6": //由于六合彩是手动开奖的彩票 所以不加入自动开奖开票管理类
			o.HK6 = hk6.Init(v.Id, v.Game_name, v.Freq)

		case "BJ28":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "DLT":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "FC3D":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "QXC":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "SSQ":
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		case "CAKENO":
			//加拿大3.5 现在api 没有下一期 ,开奖时间间隔也不稳定,先不做 ctrl 服数据库中 状态为0
			o.Ltrys[v.Game_name] = Init(v.Id, v.Game_name, v.Freq)

		default:
			beego.Error("--- Lottery init error not have lottery : " + v.Game_name)
		}
	}

	beego.Info("--- Init Ltry Mgr Done !")

	for _, v := range o.Ltrys {
		go o.ListeningLtry(v.GetGameName())
	}

	return nil
}

//注意 !!!!!!!!!!!!!
//监听开采(监听开奖)(注意 : 这里开始就是多携程环境)
func (o LtryMgr) ListeningLtry(gameName string) {
	beego.Info("--- Start listening : ", gameName, " !\n")

	ltry := o.Ltrys[gameName]
	//写个死循环来开采 呵呵
	for {
		//计算现在距离下期开奖还有多少时间
		nextLtryOpenTimeDuration := ltry.GetNextOpenTime().Sub(utils.GetNowUTC8Time())
		//beego.Debug("--- 当前彩票 ", ltry.GetGameName(), "!")
		//beego.Debug("--- 下期期数 ", ltry.GetNextExpect(), "!")
		//beego.Debug("--- 下期开奖时间 ", ltry.GetNextOpenTime(), "!\n")
		//距离时间大于0 就是还没到开奖时间， 小于0证明已经错过了开奖时间，说明在初始化彩票到 检查监听这段时间已经开除新的一期
		if nextLtryOpenTimeDuration > 0 {
			//beego.Debug("------------------------- 计算下期时间正确 ! 休眠 : ", nextLtryOpenTimeDuration)
			time.Sleep(nextLtryOpenTimeDuration)

		} else if nextLtryOpenTimeDuration <= 0 {
			newestRecord, err := apimgr.Instance().GetNewRecord(gameName)
			if err != nil {
				//等写完功能这里来将彩票状态修改为 2,这部分现在新的控制服还没有做
				beego.Error("--- 严重问题 : ", gameName, " Api获取最新记录失败 ")
			} else {
				nowExpect, err := strconv.Atoi(newestRecord.Open[0].Expect)
				if err != nil {
					beego.Error("--- 发现期数错误,从API获取了正确的数据结构,但是期数有问题 : ")
					beego.Error(err)
					beego.Error(newestRecord)
				} else {
					//最新期数不能小于等于当前期数
					if nowExpect <= ltry.GetCurrentExpect() {
						//beego.Debug("-- 已到开采请求时间,但是未获取到最新开采记录, 休眠 : ", conf.ListeningSleepTime)

						time.Sleep(conf.ListeningSleepTime * time.Second)

					} else /* if nowExpect-lottery.GetCurrentExpect() == 1*/ { //!!!!!!!!!!!!!!! 注意 由于跨天过后 期数的加减就不是==1了 这里暂时先这样，等功能写完来加正确的期数判断
						//记录入库
						//DbMgr.Instance().InsertLotteryHistoryByNewest(newestRecord)
						ltry.StartLtry(newestRecord)
					}
				}
			}
		}
	}
}

//启动彩种(初始化彩种)注意:这里启动的是手动开奖系统的彩票,目前只有六合彩
//之所以不能在初始化的时候启动,是应为必须手动输入一些参数才能正确启动
//参数为: 1.彩票名字 2.当前期数(已经开奖的这期) 3.当前这期开奖结果 4.当前这期开奖时间 5.下期期数 6.下期开奖时间
func (o *LtryMgr) InitManLtry(ltryName string, currentExpect int, nextExpect int, nextOpenTime string) int {
	switch ltryName {
	case "HK6": //香港六合彩(手动)
		if o.HK6 != nil {
			ec := o.HK6.InitLtry(currentExpect, nextExpect, nextOpenTime)
			if ec != 0 {
				return ec
			}
		} else {
			beego.Error("HK6 is nil !")
			return 88
		}

	default:
		beego.Error("Init Manual Lottery Error : ", ltryName, " Not Have Lottery Type !")
		return 99
	}
	return 0
}

//手动彩票开采
func (o *LtryMgr) StartManLtry(ltryName string, expect int, openCode string, openTime string, nextExpect int, nextOpenTime string) int {
	switch ltryName {
	case "HK6": //香港六合彩(手动)
		if o.HK6 != nil {
			ec := o.HK6.StartLtry(expect, openCode, openTime, nextExpect, nextOpenTime)
			if ec != 0 {
				return ec
			}
		} else {
			beego.Error("HK6 is nil !")
			return 88
		}

	default:
		beego.Error("Start Manual Lottery Error : ", ltryName, " Not Have Lottery Type !")
		return 99
	}
	return 0
}

type ReqGetLtryInfo struct {
	Func string
}

type LtryDb struct {
	Id           int
	Name         string //彩票中文名
	Game_name    string //彩票英文名
	Freq         int    //频率: 0.低频 1.高频
	Api_code_kcw string //kcw 代表这个是开采网的Api彩票请求代码,这个值是要组合在url中使用的
}

// //得到一个彩票最新的开彩记录(循环每一个接口,每个接口重试3次)
// func GetNewRecord(gameName string) (LtryRecordNew, error) {
// 	var ret LtryRecordNew

// 	for i := 0; i < len(apimgr.Instance().ApiNew); i++ {
// 		apiMap := apimgr.Instance().ApiNew[i]

// 		if apiArray, ok := apiMap[gameName]; ok {
// 			//循环每一个不同的接口(获取最新记录和获取历史记录不一样,优先使用0号高防接口)
// 			for j := 0; j < len(apiArray); j++ {
// 				count := 3 //重试次数,这里暂时写死 以后这些全部要写到控制服数据库中去
// 				for ; count > 0; count-- {
// 					url := apiArray[j]
// 					resp, err := httpmgr.Get(url)
// 					if err == nil {
// 						err := json.Unmarshal(resp, &ret)
// 						if err == nil {
// 							return ret, nil
// 						}
// 					}
// 					beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", gameName, " API index : ", j)
// 					time.Sleep(conf.GetRecordByNewestSleepTime * time.Second)
// 				}
// 			}
// 		} else {
// 			return ret, errors.New("There is no API for this lottery : " + gameName)
// 		}
// 		beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", gameName)
// 	}

// 	return ret, errors.New("--- Failed to obtain data from all API providers !!! Game Name : " + gameName)
// }

//启动彩种(初始化彩种)注意:这里启动的是手动开奖系统的彩票,目前只有六合彩
//之所以不能在初始化的时候启动,是应为必须手动输入一些参数才能正确启动
//参数为: 1.彩票名字 2.当前期数(已经开奖的这期) 3.当前这期开奖结果 4.当前这期开奖时间 5.下期期数 6.下期开奖时间
// func (o *ManlLtryMgr) InitLottery(ltryName string, currentExpect int, currentOpenCode string, currentOpenTime string, nextExpect int, nextOpenTime string) int {
// 	switch ltryName {
// 	case "HK6": //香港六合彩(手动)
// 		hk6 := ltrymgr.Instance().HK6
// 		if hk6 != nil {
// 			ec := hk6.InitLtry(currentExpect, currentOpenCode, nextExpect, nextOpenTime)
// 			if ec != 0 {
// 				return ec
// 			}
// 			o.lock.Lock()
// 			defer o.lock.Unlock()
// 			o.ltryifMap[ltryName] = hk6
// 		} else {
// 			beego.Error("------------------------- HK6 Already Initialized ! -------------------------")
// 		}

// 	default:
// 		beego.Error("------------------------- Init Manual Lottery Error : ", ltryName, " Not Have Lottery Type ! -------------------------")
// 		return
// 	}
// }

// //开采网 按最新获取一条记录(带下期)
// type LtryRecordNew struct {
// 	Rows   int
// 	Code   string
// 	Remain string
// 	Next   []LtryRecordNext    //下一期信息(下一期期数, 下一期开奖时间)
// 	Open   []LtryRecordCurrent //最新一期信息 (注意:这里面没有时间戳,存库的时候要自己添加一个)
// 	Time   string              //查询时间
// }

// //开采网 下一期开采信息
// type LtryRecordNext struct {
// 	Expect   string
// 	Opentime string
// }

// //开采网 当前这期信息
// type LtryRecordCurrent struct {
// 	Expect   string
// 	Opencode string
// 	Opentime string
//}
