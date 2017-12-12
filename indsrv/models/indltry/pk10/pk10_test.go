package pk10

import (
	"indsrv/models/rd"
	"strconv"
	"testing"

	"github.com/astaxie/beego"
)

func Test_PaserBetNum(t *testing.T) {
	beego.Debug("--- 开始测试\n")

	// //1：初始化加密管理类(这个必须是第一初始化的,后续的消息都要依赖)
	// err := encmgr.Instance().Init()
	// if err != nil {
	// 	beego.Error(err)
	// 	return
	// }

	// //2. 初始化控制管理类(这个必须是第二初始化,向控制服发送注册信息,获取密钥会放入 encmgr 中)
	// err = ctrl.Init()
	// if err != nil {
	// 	beego.Error(err)
	// 	return
	// }

	// //3.初始化数据库(链接数据库)(dbmgr 不再是单列类)
	// err = dbmgr.Init()
	// if err != nil {
	// 	beego.Error(err)
	// 	return
	// }

	err := rd.Init()
	if err != nil {
		beego.Error(err)
		return
	}

	var oc []int

	for len(oc) < 10 {
		r := rd.Intn(1, 10)
		exist := false

		for _, v := range oc {
			if v == r {
				exist = true
				break
			}
		}

		if !exist {
			oc = append(oc, r)
		}
	}

	var ocs string
	ocs = strconv.Itoa(oc[0])
	for i := 1; i < 10; i++ {
		ocs = ocs + "," + strconv.Itoa(oc[i])
	}

	beego.Debug("--- 测试结束\n")
}
