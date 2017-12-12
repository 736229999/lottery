// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"apisrv/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//获取最新彩票记录
	beego.Router("/getLtryNewestRecord", &controllers.GetLtryNewestRecord{})

	//获取彩票历史记录
	beego.Router("/getLtryHist", &controllers.GetLtryHist{})

	//按天获取彩票历史记录
	beego.Router("/getLtryHistByDay", &controllers.GetLtryHistByDay{})

	//按彩票和期数来获取一期结果
	beego.Router("/getLtryRecordByExpect", &controllers.GetLtryRecordByExpect{})

	//------------------------------------------ 来自总管理后台的消息 ------------------------------------------
	//初始化手动开奖彩票
	beego.Router("/initManLtry", &controllers.InitLtry{})
	//手动开奖
	beego.Router("/startManLtry", &controllers.StartLtry{})
}
