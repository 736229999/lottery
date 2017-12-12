package AnnouncementMgr

import (
	"gamesrv/models/dbmgr"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

//公告管理员
//活动信息管理
type AnnouncementMgr struct {
	Announcements        []Announcement //为分类所有公告
	AllAnnouncements     []Announcement //面向所有用户的公告
	GroupAnnouncements   []Announcement //面向特定用户组的公告
	AccountAnnouncements []Announcement //面向特定用户的公告
}

var sInstance *AnnouncementMgr
var once sync.Once

//单例模式
func Instance() *AnnouncementMgr {
	once.Do(func() {
		sInstance = &AnnouncementMgr{}
		sInstance.init()
	})

	return sInstance
}

//启动时初始化,
func (o *AnnouncementMgr) init() {
	//从数据库获取活动信息
	o.UpdateAnnouncement()
}

//更新活动数据
func (o *AnnouncementMgr) UpdateAnnouncement() {
	//清空之前的数据
	o.Announcements = nil
	o.AllAnnouncements = nil
	o.GroupAnnouncements = nil
	o.AccountAnnouncements = nil

	bm := bson.M{"status": 1}
	dbmgr.Instance().AnnouncementCollection.Find(bm).All(&o.Announcements)

	//分析公告
	for _, v := range o.Announcements {
		if v.Type == 3 { //所有用户
			o.AllAnnouncements = append(o.AllAnnouncements, v)
		} else if v.Type == 2 { //用户组
			o.GroupAnnouncements = append(o.AllAnnouncements, v)
		} else if v.Type == 1 { //指定用户
			o.AccountAnnouncements = append(o.AllAnnouncements, v)
		}
	}
}

type Announcement struct {
	Platform string `bson:"platform"` //平台 : all  ios  android
	Type     int    `bson:"type"`     //3 所有用户 2指定用户
	Title    string `bson:"title"`    //标题
	Content  string `bson:"content"`  //内容
	Time     int    `bson:"time"`     //发布时间
	Status   int    `bson:"status"`   //状态 1可用, 2不可用
	Group    []int  `bson:"group"`    //用户组
	Accounts string `bson:"accounts"` //指定账号,当Type 类型为2的时候 只有指定账号的玩家能看到公告
}
