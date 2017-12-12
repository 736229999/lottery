package main

import (
	"calculsrv/models/ctrl"
	"calculsrv/models/dbmgr"
	"calculsrv/models/encmgr"
	"calculsrv/models/gamemgr"
	"calculsrv/models/ltrymgr"
	"calculsrv/models/ltryset"
	_ "calculsrv/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Info("------------------------- Star Server -------------------------")

	//1：初始化加密管理类(这个必须是第一初始化的,后续的消息都要依赖)
	err := encmgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//2. 初始化控制管理类(这个必须是第二初始化,向控制服发送注册信息,获取密钥会放入 encmgr 中)
	err = ctrl.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//3.初始化数据库(链接数据库)(dbmgr 不再是单列类)
	err = dbmgr.Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//4.初始化Game Server管理
	err = gamemgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	//5.彩票设置初始化
	err = ltryset.Init()
	if err != nil {
		beego.Error(err)
		return
	}

	err = ltrymgr.Instance().Init()
	if err != nil {
		beego.Error(err)
		return
	}

	// //第二补:控制中心实例并初始化(负责服务器状态,彩票状态等,管理员后台的操作)
	// ControlCenter.Instance()

	// //彩票历史管理员实例并初始化(这里会补全在服务器没有启动时中间差的数据)
	// HistoryManager.Instance()

	// //补全未开机时间所有的未开奖订单
	// //OrderCheckSys.Instance()

	// //彩票管理员实例并初始化(这时会通过api去获取一次最新的记录)
	// LotteryManager.Instance()

	beego.Run()
}

//----------------------------这个复杂的超时.....以后来研究下玩玩-----------------------------

//超时机制后面会用到
// c := http.Client{
// 	Transport: &http.Transport{
// 		Dial: func(netw, addr string) (net.Conn, error) {
// 			deadline := time.Now().Add(25 * time.Second)         //返回超时
// 			c, err := net.DialTimeout(netw, addr, time.Second*5) //连接超时
// 			if err != nil {
// 				return nil, err
// 			}
// 			c.SetDeadline(deadline)
// 			return c, nil
// 		},
// 	},
// }
// c.PostForm("http://google.com/", nil)
