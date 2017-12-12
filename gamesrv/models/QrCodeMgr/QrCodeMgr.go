package QrCodeMgr

import (
	"gamesrv/models/dbmgr"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

//公告管理员
//活动信息管理
type QrCodeMgr struct {
	QrCodes []QrCode
}

var sInstance *QrCodeMgr
var once sync.Once

//单例模式
func Instance() *QrCodeMgr {
	once.Do(func() {
		sInstance = &QrCodeMgr{}
		sInstance.init()
	})

	return sInstance
}

//启动时初始化,
func (o *QrCodeMgr) init() {
	//从数据库获取活动信息
	o.UpdateQrCodeMgr()
}

//更新活动数据
func (o *QrCodeMgr) UpdateQrCodeMgr() {
	//清空之前的切片
	o.QrCodes = nil

	bm := bson.M{"status": 1}
	dbmgr.Instance().QrCodeCollection.Find(bm).All(&o.QrCodes)
}

type QrCode struct {
	QrId       int    `bson:"id"`
	QrPath     string `bson:"qr_path"`
	QrPlatform string `bson:"qr_platform"`
	Title      string `bson:"title"`
	Subtitle   string `bson:"subtitle"`
	Status     int    `bson:"status"`
	Group      []int  `bson:"group"`
}
