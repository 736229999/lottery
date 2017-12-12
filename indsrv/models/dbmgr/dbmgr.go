package dbmgr

import (
	"indsrv/models/ctrl"
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

var indSess *mgo.Session

var histDb *mgo.Database

var PK10Coll *mgo.Collection

func Init() error {
	var dburl = ctrl.DbSrv.Ip + ":" + ctrl.DbSrv.Port
	beego.Info("--- Independent DB URL : ", dburl)

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{dburl},
		Timeout:  time.Second * 3,
		Database: dbName,
		Username: userName,
		Password: pwd,
	}

	//连接计算服数据库（授权访问）
	var err error
	indSess, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}

	// 采用 Strong 模式
	indSess.SetMode(mgo.Strong, true)

	histDb = indSess.DB(dbName)

	PK10Coll = histDb.C(pk10Coll)

	beego.Info("--- Init DB Done !")
	return nil
}
