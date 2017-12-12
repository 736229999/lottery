package dbmgr

import (
	"calculsrv/models/ctrl"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

//数据库管理员类
//type DbMgr struct {
//计算服数据库Session
var CalculSess *mgo.Session

//计算服DB
var CalculDb *mgo.Database

//用户信息
var AcColl *mgo.Collection

// 开彩历史记录
var HistColl *mgo.Collection

// 号码走势记录
var TrendColl *mgo.Collection

// 用户下注
var BetColl *mgo.Collection

// 订单
var OrderColl *mgo.Collection

//------------------- 后台 ----------------
//管理后台数据库Session
var MgrSess *mgo.Session

//管理后台DB
var MgrDb *mgo.Database

//流水记录	(注意这个是需要写入的)
var BalanceRecordColl *mgo.Collection

//彩票信息表名
var LtryColl *mgo.Collection

//Game服务器表名
var GameSrvColl *mgo.Collection

//彩票设定表名(赔率，限额等)
var LtrySetColl *mgo.Collection

//}

//var sInstance *DbMgr
//var once sync.Once

/*
取得数据库连接类实例,单例模式
*/
// func Instance() *DbMgr {
// 	once.Do(func() {
// 		sInstance = &DbMgr{}
// 	})

// 	return sInstance
// }

func Init() error {
	//查看自身是什么服务器
	if ctrl.SelfSrv.Type == 0 { //试玩服务器
		beego.Info("--- Calculation Server  : Trial !")
	} else if ctrl.SelfSrv.Type == 1 {
		beego.Info("--- Calculation Server  : Formal !")
	} else {
		beego.Error("Calculation Server Type Error !")
	}

	//组依赖数据库url
	var dburl = ctrl.DbSrv.Ip + ":" + ctrl.DbSrv.Port
	beego.Info("--- Calculation DB URL : ", dburl)

	//计算服
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{dburl},
		Timeout:  time.Second * 3,
		Database: CalculationDbName,
		Username: CalculationDbUserName,
		Password: CalculationDbPwd,
	}

	//连接计算服数据库（授权访问）
	var err error
	CalculSess, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}

	// 采用 Strong 模式
	CalculSess.SetMode(mgo.Strong, true)

	CalculDb = CalculSess.DB(CalculationDbName)

	AcColl = CalculDb.C(AccountInfoCollection)
	HistColl = CalculDb.C(HistoryCollection)
	TrendColl = CalculDb.C(TrendCollection)
	BetColl = CalculDb.C(BetCollection)
	OrderColl = CalculDb.C(OrderCollection)

	//----------------------连接管理员后台数据库-------------------------
	//组依赖管理数据库url
	var mgrurl = ctrl.MgrDb.Ip + ":" + ctrl.MgrDb.Port
	beego.Info("--- Mgr DB URL : ", mgrurl)

	dialInfo_1 := &mgo.DialInfo{
		Addrs: []string{mgrurl},
		//注意这里如果不是admin验证，那么一定要记住给出对应数据库名
		Database: ManageDbName,
		Timeout:  time.Second * 3,
		Username: ManageDbUserName,
		Password: ManageDbPwd,
	}

	MgrSess, err = mgo.DialWithInfo(dialInfo_1)
	if err != nil {
		return err
	}
	// 采用 Strong 模式
	MgrSess.SetMode(mgo.Strong, true)

	MgrDb = MgrSess.DB(ManageDbName)
	LtryColl = MgrDb.C(LotteryTypeCollection)
	GameSrvColl = MgrDb.C(GameServerCollection)
	LtrySetColl = MgrDb.C(LotteriesSettingsCollection)
	BalanceRecordColl = MgrDb.C(BalanceRecordCollection)

	beego.Info("--- Init DB Done !")
	return nil
}

//插入数据
func Insert(collection *mgo.Collection, msg interface{}) bool {
	err := collection.Insert(msg)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

//批量插入
func BulkInsert(collection *mgo.Collection, msg *[]interface{}) bool {
	err := collection.Insert(*msg...)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
