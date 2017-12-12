package BetManger

import "sync"

//数据库管理员类
type BetManger struct {
}

var sInstance *BetManger
var once sync.Once

/*
取得数据库连接类实例,单例模式
*/
func Instance() *BetManger {
	once.Do(func() {
		sInstance = &BetManger{}
		sInstance.init()
	})

	return sInstance
}

func (o *BetManger) init() {

}
