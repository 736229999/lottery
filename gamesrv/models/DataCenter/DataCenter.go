package DataCenter

import "sync"

//数据中心文件类,将所有需要客户端平凡请求的数据放入这里,作为缓存的形式纯在
type DataCenter struct {
}

var sInstance *DataCenter
var once sync.Once

//单例模式
func Instance() *DataCenter {
	once.Do(func() {
		sInstance = &DataCenter{}
		sInstance.init()
	})

	return sInstance
}

//启动时初始化,
func (o *DataCenter) init() {

}
