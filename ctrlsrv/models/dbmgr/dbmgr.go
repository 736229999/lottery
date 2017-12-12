package dbmgr

import (
	"sync"
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

type DbMgr struct {
	ctrlSess *mgo.Session  //控制服数据库session
	confDb   *mgo.Database //配置数据库

	SrvCollection     *mgo.Collection //服务器信息表名
	EncryptCollection *mgo.Collection //加密信息表名
	LtryCollection    *mgo.Collection //彩种信息表名
}

var sInstance *DbMgr
var once sync.Once

// ---返回数据库连接·？ 这个函数返回值类型为*DbMgr
func Instance() *DbMgr {
	once.Do(func() {
		sInstance = &DbMgr{}
		sInstance.init()
	})
	//beego.Debug("*****DbMgr  instance Done")
	return sInstance
}

// ---初始化获取配置文件中的数据库连接参数，这个方法接收者是*DbMgr类型（该init方法属于DbMgr类型对象中的方法）
func (o *DbMgr) init() {
	dial := &mgo.DialInfo{
		Addrs:    []string{ctrlsrvDbIP},
		Timeout:  time.Second * 3,
		Database: confDbName,
		Username: ctrlsrvDbUserName,
		Password: ctrlsrvDbPwd,
	}

	//连接计算服数据库（授权访问）
	var err error
	o.ctrlSess, err = mgo.DialWithInfo(dial)// ---连接数据库返回*mgo.Session类型的对象
	if err != nil {
		beego.Error(err)
		return
	}

	// 采用 Strong 模式
	o.ctrlSess.SetMode(mgo.Strong, true)

	o.confDb = o.ctrlSess.DB(confDbName)

	o.SrvCollection = o.confDb.C(srvColl)

	o.EncryptCollection = o.confDb.C(encryptColl)

	o.LtryCollection = o.confDb.C(ltryColl)

	beego.Info("--- Init Db Mgr Done !")
}
