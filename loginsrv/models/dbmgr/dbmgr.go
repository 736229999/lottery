package dbmgr

import (
	"loginsrv/models/ctrl"
	"sync"
	"time"

	"github.com/astaxie/beego"

	mgo "gopkg.in/mgo.v2"
)

//DbMgr 数据库管理员
type DbMgr struct {
	DbServiceMap map[string]*DbService
}

var sInstance *DbMgr
var once sync.Once

// Instance 获取DbMgr单例，Once 线程安全
func Instance() *DbMgr {
	once.Do(func() {
		sInstance = &DbMgr{}
		sInstance.DbServiceMap = make(map[string]*DbService)
	})
	return sInstance
}

func (o *DbMgr) Init() error {
	//判断改服务器是试玩还是正式服;
	// serverType, err := beego.AppConfig.Int("ServerType")
	// if err != nil {
	// 	beego.Emergency(err)
	// 	return
	// }

	if ctrl.SelfSrv.Type == 0 { //试玩服务器
		beego.Info("--- Login Server  : Trial !")
	} else if ctrl.SelfSrv.Type == 1 {
		beego.Info("--- Login Server  : Formal !")
	} else {
		beego.Error("Login Server Type Error !")
	}

	//if serverType == 0 {
	//	beego.Info("------------------- 本服务器是 试玩服务器 ！---------------------")
	//连接试玩服数据库
	var dburl = ctrl.DbSrv.Ip + ":" + ctrl.DbSrv.Port

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{dburl},
		Timeout:  time.Second * 5,
		Database: LoginDbName,
		Username: DbUserName,
		Password: DbPwd,
	}
	err := o.Connect(dialInfo, ctrl.DbSrv.Ip, LoginDbName, AccountInfoCollection)
	if err != nil {
		return err
	}
	// } else {
	// 	beego.Info("------------------- 本服务器是 正式服务器 ！！ ！ ---------------------")
	// 	//连接正式服数据库
	// 	dialInfo_1 := &mgo.DialInfo{
	// 		Addrs:    []string{FormalDbIp},
	// 		Timeout:  time.Second * 5,
	// 		Database: LoginDbName,
	// 		Username: DbUserName,
	// 		Password: DbPwd,
	// 	}
	// 	o.Connect(dialInfo_1, FormalDbIp, LoginDbName, AccountInfoCollection)
	// }

	beego.Info("--- Init DB Mgr Done !  ")
	return nil
}

//Connect 连接数据库
func (o *DbMgr) Connect(dialInfo *mgo.DialInfo, serverName string, dbName string, collectionName string) error {
	//如果没有就创建DbServer对象,注意这里是使用伪构造的方法来创建对象,相当之恶心;
	dbService := newDbService()
	var err error
	//建立连接 Dial
	dbService.Session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}

	//defer dbService.Session.Close()
	//设置连接模式;
	dbService.Session.SetMode(mgo.Strong, true)

	if err != nil {
		defer dbService.Session.Close()
		return err
	}

	dbService.DbMap[dbName] = dbService.Session.DB(dbName)

	dbService.CollectionMap[collectionName] = dbService.DbMap[dbName].C(collectionName)

	o.DbServiceMap[serverName] = dbService

	beego.Info("--- DbMgr Connect DB Done : ", dialInfo.Addrs, " - ", dbName, " - ", collectionName, "!")

	//首先判断之前有没有这个连接在map里面;
	//dbService, ok := dbMgr.DbServiceMap[dialInfo.Addrs]

	// if ok {
	// 	fmt.Println("---------------- Connection has been established : ", dbURL, " ! ----------------")

	// 	//如果已经和数据库服务物理机建立连接,那么看看是不是连接到了一样的数据库
	// 	_, ok := dbService.DbMap[dbName]
	// 	if ok {
	// 		fmt.Println("---------------- Link to the same db server : ", dbName, " ! ----------------")
	// 		//如果链接到想同的数据库,那么再检查是否打开同样的表(Collection)
	// 		_, ok := dbService.CollectionMap[collectionName]
	// 		if ok {
	// 			fmt.Println("---------------- Open the same collection : ", collectionName, " ! ----------------")
	// 			return
	// 		}
	// 		//链接到同一个数据库服务,同一个数据库,但是打开不同的表单;
	// 		dbService.CollectionMap[collectionName] = dbService.DbMap[dbName].C(collectionName)
	// 		if dbService.CollectionMap[collectionName] == nil {
	// 			beego.Emergency("---------------- Error connect  : ", dbURL, " DB Name : ", dbName, "Collection Name : ", collectionName, "! ----------------")
	// 			return
	// 		}

	// 		fmt.Println("---------------- DbMgr Connect DB Done : ", dbURL, " - ", dbName, " - ", collectionName, "! ----------------")
	// 		return
	// 	}
	// 	//同一个数据库服务下,连接不同的数据库
	// 	dbService.DbMap[dbName] = dbService.Session.DB(dbName)
	// 	if dbService.DbMap[dbName] == nil {
	// 		beego.Emergency("---------------- Error connect  : ", dbURL, " - ", dbName, "! ----------------")
	// 		return
	// 	}
	// 	//同一个数据库服务下,打开不同的表单(Collection)
	// 	dbService.CollectionMap[collectionName] = dbService.DbMap[dbName].C(collectionName)
	// 	if dbService.CollectionMap[collectionName] == nil {
	// 		beego.Emergency("---------------- Error connect  : ", dbURL, " DB Name : ", dbName, "Collection Name : ", collectionName, "! ----------------")
	// 		return
	// 	}

	// 	fmt.Println("---------------- DbMgr Connect DB Done : ", dbURL, " - ", dbName, " - ", collectionName, "! ----------------")
	// 	return
	// }
	return nil
}

//在DbService已经存在的情况下,连接不同的数据库
func connectDb(dbName string, collectionName string) {

}

//在DbService已经存在的情况下,连接不同的表
func connectCollection(collectionName string) {

}
