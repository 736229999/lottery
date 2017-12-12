// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"ctrlsrv/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//------------------------------ 所有服务器消息 ---------------------------
	//服务器验证(服务器注册)
	beego.Router("/srvRegist", &controllers.SrvRegist{})

	//得到服务器依赖信息(比如:计算服 需要得到 Game 服务器的ip)
	beego.Router("/getDependSrv", &controllers.GetDependSrv{})

	//------------------------------ Api 服务器消息 ---------------------------
	beego.Router("/getLtryInfo", &controllers.GetLtryInfo{})
}
