package srvmgr

import (
	"ctrlsrv/models/dbmgr"
	"sync"

	"github.com/astaxie/beego"
)

type SrvMgr struct {
	SrvInfos map[string]srvInfo //使用服务器名称作为key(这个是从数据库中读取的允许注册的服务器信息)
}

var sInstance *SrvMgr
var once sync.Once

func Instance() *SrvMgr {
	once.Do(func() {
		sInstance = &SrvMgr{}
		sInstance.init()
	})
	//beego.Debug("*****SrvMgr  instance Done")
	return sInstance
}

//
func (o *SrvMgr) init() {
	o.SrvInfos = make(map[string]srvInfo)

	ret := []srvInfo{}
	dbmgr.Instance().SrvCollection.Find(nil).All(&ret)

	for _, v := range ret {
		o.SrvInfos[v.Name] = v
	}

	beego.Info("--- Init Srv Mgr Done !")
}

//验证这个服务器是否在服务器map里
func (o SrvMgr) VerifySrv(ip string) bool {
	for _, v := range o.SrvInfos {
		if v.Ip == ip {
			return true
		}
	}

	return false
}

//服务器信息(注意 这里读取数据库的时候没有 加上 bson 尾部,这样代码看起来更优雅,并且可以加上 json 传输字段,但是写得时候一定要注意不能打错字)
type srvInfo struct {
	Company string //公司名称(盘口名称)
	Name    string //服务器名
	Ip      string
	Port    string
	Func    string   //服务器功能 login, game, calcul, api
	Type    int      //服务器类型 0.试玩服 1.正式服 2.测试服
	Status  int      //0.不可用 1.正常(可用)
	Depend  []string //依赖服务器名, 比如 api 服务器必须要有一个 api 数据库服务器,不管数据库是不是和api服务器在同一台物理机
}
