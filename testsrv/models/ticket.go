package models

import "sync"
import "github.com/astaxie/beego"

//数据库管理员类
type Ticket struct {
	ticket int
}

var sInstance *Ticket
var once sync.Once

/*
取得数据库连接类实例,单例模式
*/
func Instance() *Ticket {
	once.Do(func() {
		sInstance = &Ticket{}
		sInstance.ticket = 30
	})

	return sInstance
}

func (o *Ticket) Output(v int) {
	beego.Debug(v)
}

func (o Ticket) GetTicket() int {
	return o.ticket
}
