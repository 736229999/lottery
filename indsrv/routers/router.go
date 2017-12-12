// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"indsrv/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//获取最新彩票记录
	beego.Router("/getRecordByNew", &controllers.GetRecordByNew{})
	//获得历史记录
	beego.Router("/getHist", &controllers.GetHist{})
	//按期数获取某一期,某一菜种的记录
	beego.Router("/getRecordByExpect", &controllers.GetRecordByExpect{})
}
