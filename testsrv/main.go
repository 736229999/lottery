package main

import (
	"testsrv/models"
	_ "testsrv/routers"

	"github.com/astaxie/beego"
)

var ticket = 20

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Debug("--- test srv")

	//var lock sync.Mutex

	go func(id int) {
		for i := 0; i < 10; i++ {
			beego.Debug("--Thread one call :")
			models.Instance().Output(i)
		}
	}(1)

	go func(id int) {
		for i := 10; i < 20; i++ {
			beego.Debug("--Thread two call :")
			models.Instance().Output(i)
		}
	}(2)

	go func(id int) {
		for i := 20; i < 30; i++ {
			beego.Debug("--Thread three call :")
			models.Instance().Output(i)
		}
	}(3)

	// go func() {
	// 	for {
	// 		if models.Instance().GetTicket() > 0 {
	// 			beego.Debug("--Thread three call :", models.Instance().SellTicket())
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		if models.Instance().GetTicket() > 0 {
	// 			beego.Debug("--Thread four call :", models.Instance().SellTicket())
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }()

	beego.Run()
}

func SellTicket() {
	for {
		// if ticket > 0 {
		// 	ticket--
		// 	beego.Debug(ticket)
		// } else {
		// 	break
		// }
		//	lock.Lock()
		ticket--
		beego.Debug(ticket)
		if ticket == 0 {
			break
		}
		//	lock.Unlock()
	}
}
