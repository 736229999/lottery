package Order

import (
	"calculsrv/models/gb"
	"fmt"
	"sync"
	"time"
)

func Instance() *Order {
	once.Do(func() {
		sInstance = &Order{}
		sInstance.init()
	})
	return sInstance
}

type Order struct {
	index int
}

var sInstance *Order

var once sync.Once

/*
初始化订单号码提供器
*/
func (o *Order) init() {
	o.index = 0
}

/*
获取订单号
*/
func (o *Order) GetOrderNumber() string {
	o.index = (o.index + 1) % 1000
	// 本地编号取四位
	orderNum := time.Now().Format("20060102150405") + // 年月日时分秒 +
		fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000) + // 纳秒前三位 +
		fmt.Sprintf("%03d", o.index) + // 千位循环
		gb.MachineCode //计算服id
	return orderNum
}
