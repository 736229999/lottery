package ActivityMgr

import (
	"gamesrv/models/dbmgr"
	"sync"
)

//活动信息管理
type ActivityMgr struct {
	ActivityInfoArray []ActivityInfo
}

var sInstance *ActivityMgr
var once sync.Once

//单例模式
func Instance() *ActivityMgr {
	once.Do(func() {
		sInstance = &ActivityMgr{}
		sInstance.init()
	})

	return sInstance
}

//启动时初始化,
func (o *ActivityMgr) init() {
	//从数据库获取活动信息
	o.UpdateActivityInfo()
}

//更新活动数据
func (o *ActivityMgr) UpdateActivityInfo() {
	//注意这里我确定是否会产生内存泄露,理论上是应该不会的,等实际情况
	dbmgr.Instance().ActivityCollection.Find(nil).All(&o.ActivityInfoArray)
}

type ActivityInfo struct {
	Image        string `bson:"image_url"`
	ContentUrl   string `bson:"content_url"`
	Sort         int    `bson:"sort"`
	Status       int    `bson:"status"`
	Show         int    `bson:"show"` //1 显示 2不显示
	Title        string `bson:"title"`
	ActivityTime string `bson:"activity_time"`
}
