package ltrymgr

import (
	"calculsrv/models/LotteryTypes/EX5"
	"calculsrv/models/LotteryTypes/HK6"
	"calculsrv/models/LotteryTypes/K3"
	"calculsrv/models/LotteryTypes/PK10"
	"calculsrv/models/LotteryTypes/PL3"
	"calculsrv/models/LotteryTypes/SSC"
	"calculsrv/models/LotteryTypes/bj28"
	"calculsrv/models/apimgr"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"
	"calculsrv/models/ind"
	"common/utils"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

type LtryMgr struct {
	LtryifMap   map[string]Ltryif         //所有彩票对象	k为gameTag
	LtryInfoMap map[string]gb.LotteryInfo //从后台获取的所有彩票信息
}

//-------------------------------------------------- Single mode --------------------------------------------------

var sInstance *LtryMgr
var once sync.Once

// Instance 获取dbmgr单例，Once 线程安全
func Instance() *LtryMgr {
	once.Do(func() {
		sInstance = &LtryMgr{}
		//sInstance.initLotteriesInfo()
	})
	return sInstance
}

//--------------------------------------------------  Method --------------------------------------------------
//初始化所有彩票(注意这个函数在程序启动时会根据数据库的彩票状态调用一次，然后在程序运行过程中会根据管理后台的消息调用)
//新的架构,这个地方彩票依然是从后台管理服获取,应为每一个代理商的彩票设置是不一样的
func (o *LtryMgr) Init() error {
	o.LtryifMap = make(map[string]Ltryif)

	//从管理员后台获取菜种信息
	o.LtryInfoMap = dbmgr.UpdateLotteriesInfo()

	//状态正常的彩票彩初始化
	for _, v := range o.LtryInfoMap {
		if v.Status == gb.LotteryStatus_Nonmal {
			o.initLtry(v)
		}
	}

	beego.Debug("--- Init All Lottery Done !")
	return nil
}

//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//!!!!!!!!!!!!!!!!!!!!!! 新加彩票 要加这里!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//更新彩票类型信息。新加彩票要放这里，并且管理后台发来的开启彩票功能也要调用这个函数
func (o *LtryMgr) initLtry(lif gb.LotteryInfo) {
	switch lif.GameTag {
	//11选5
	case gb.EX5_JiangXi, gb.EX5_ShanDong, gb.EX5_ShangHai, gb.EX5_BeiJing, gb.EX5_FuJian, gb.EX5_HeiLongJiang, gb.EX5_JiangSu:
		ltry, err := EX5.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	//快3
	case gb.K3_GuangXi, gb.K3_JiLin, gb.K3_AnHui, gb.K3_BeiJing, gb.K3_FuJian, gb.K3_HeBei, gb.K3_ShangHai, gb.K3_JiangSu:
		ltry, err := K3.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	//时时彩
	case gb.SSC_ChongQing, gb.SSC_TianJin, gb.SSC_XinJiang:
		ltry, err := SSC.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	// //PK10
	case gb.PK10_BeiJing, gb.PK10_F:
		ltry, err := PK10.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	//排列3
	case gb.PL3:
		ltry, err := PL3.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	//香港6 六合彩 现在加入监听系统应为 api 服务器负责开奖
	case gb.HK6:
		ltry, err := HK6.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	//PCDD 北京28
	case gb.BJ28:
		ltry, err := bj28.Init(lif)
		if err != nil {
			beego.Error(err)
			return
		}
		o.LtryifMap[lif.GameTag] = ltry

	default:
		beego.Error("Init Lottery Error : ", lif.GameTag, " Not Have Type !")
		return
	}

	//每个彩票用一个线程监听开采
	go o.ListeningLottery(lif.GameTag)
}

//注意 !!!!!!!!!!!!!
//监听开采(监听开奖)(注意 : 这里开始就是多携程环境)

func (o *LtryMgr) ListeningLottery(gameTag string) {
	beego.Info("--- Start listening : ", gameTag, " ", "\n")

	ltry := o.LtryifMap[gameTag]
	//写个死循环来开采 呵呵
	for {
		//计算现在距离下期开奖还有多少时间
		nextLotteryOpenTimeDuration := ltry.GetNextReqTime().Sub(utils.GetNowUTC8Time())

		if nextLotteryOpenTimeDuration >= 0 {
			time.Sleep(nextLotteryOpenTimeDuration)
		}

		o.LooplisteningLottery(gameTag)
	}
}

//上面调用 循环监听结果
func (o *LtryMgr) LooplisteningLottery(gameTag string) {
	ltry := o.LtryifMap[gameTag]
	var record gb.LtryRecordByNewest
	var err error
	for {
		//暂时在这里加入判断,等后面功能完成再来思考这里是否要为独立彩票单独列一个类
		if gameTag == "PK10_F" {
			record, err = ind.GetRecordByNew(gameTag)
			if err != nil {
				beego.Error(err)
				return
			}
		} else {
			//通过API得到最新的一期的开彩记录（按最新，带下期）
			record, err = apimgr.Instance().GetRecordByNewest(gameTag)
			if err != nil {
				beego.Error(err)
				return
			}
		}

		nowExpect := record.CurrentExpect
		//最新期数不能小于等于当前期数
		if nowExpect <= ltry.GetCurrentExpect() {
			time.Sleep(2 * time.Second)
		} else {
			ltry.StartLottery(record)
			return
		}
	}
}
