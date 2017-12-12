// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"loginsrv/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//获取Game服务器IP
	beego.Router("/bewithyou", &controllers.GetGameIpController{})
	//获取验证码
	beego.Router("/maytheforce", &controllers.CaptchaController{})
	//注册
	beego.Router("/makethman", &controllers.RegistController{})
	//登陆
	beego.Router("/manners", &controllers.LoginController{})
	//修改密码
	beego.Router("/ModifyPassword", &controllers.ModifyPasswordController{})

	//---------------------------------------GameServer------------------------------------
	//beego.Router("/RegistGameServer", &controllers.RegistGameServerController{})
}
