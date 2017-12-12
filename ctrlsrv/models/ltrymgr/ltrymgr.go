package ltrymgr

import (
	"ctrlsrv/models/dbmgr"
	"sync"

	"github.com/astaxie/beego"
)

type LtryMgr struct {
	Ltrys map[string]Ltry
}

var sInstance *LtryMgr
var once sync.Once

func Instance() *LtryMgr {
	once.Do(func() {
		sInstance = &LtryMgr{}
		sInstance.init()
	})
	//beego.Debug("*****LtryMgr  instance Done")
	return sInstance
}

func (o *LtryMgr) init() {
	o.Ltrys = make(map[string]Ltry)

	ret := []Ltry{}

	err := dbmgr.Instance().LtryCollection.Find(nil).All(&ret)
	if err != nil {
		beego.Error(err)
		return
	}

	for _, v := range ret {
		o.Ltrys[v.Game_name] = v
	}

	beego.Info("--- Init Ltry Mgr  Done !")
}

type Ltry struct {
	Id           int
	Name         string //彩票中文名
	Game_name    string //彩票英文名
	Parent_name  string //彩票大类名字
	Freq         int    //频率: 0.低频 1.高频
	Status       int    //状态: 0.关闭 1.正常 2.维护
	Recommend    int    //推荐: 0.不推荐 1.推荐
	Sort         int    //排序: 用作推荐排序,数字越小越靠前
	Api_code_kcw string //kcw 代表这个是开采网的Api彩票请求代码,这个值是要组合在url中使用的
}
