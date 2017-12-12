package ltrymgr

import (
	"indsrv/models/indltry/pk10"
	"sync"

	"github.com/astaxie/beego"
)

type LtryMgr struct {
	LtryifMap map[string]Ltryif //所有彩票对象	k为game name
}

//-------------------------------------------------- Single mode --------------------------------------------------

var sInstance *LtryMgr
var once sync.Once

// Instance 获取dbmgr单例，Once 线程安全
func Instance() *LtryMgr {
	once.Do(func() {
		sInstance = &LtryMgr{}
	})
	return sInstance
}

//--------------------------------------------------  Method --------------------------------------------------
//初始化所有彩票,由于独立彩票的特殊性,这里就不读取后台信息了,不管后台是否开放这个彩票,只要indsrv启动,就要持续开采
func (o *LtryMgr) Init() error {
	o.LtryifMap = make(map[string]Ltryif)

	//状态正常的彩票彩初始化
	ltry, err := pk10.Init()
	if err != nil {
		return err
	}
	o.LtryifMap[ltry.GameName] = ltry

	beego.Debug("--- Init All Lottery Done !")
	return nil
}
