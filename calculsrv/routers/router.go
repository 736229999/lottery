package routers

import (
	"calculsrv/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//彩票投注
	beego.Router("/LotteryBetting", &controllers.LotteryBetting{})
	//修改彩票信息(排序,推荐等)
	beego.Router("/UpdateLtryInfo", &controllers.UpdateLtryInfo{})
	//修改彩票设置(赔率)
	beego.Router("/UpdateLtrySet", &controllers.UpdateLtrySet{})
	//手动开奖(Manual lottery 预防掉期补全失败需要手动开奖的情况)
	beego.Router("/ManLtry", &controllers.Manltry{})
}
