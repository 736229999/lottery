// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"gamesrv/controllers"

	"github.com/astaxie/beego"
)

func init() {

	//------------------------------这些是CalculationServer 发来的消息------------------------------
	//Claculation Server 注册消息
	beego.Router("/RegistCalculation", &controllers.RegistCalculation{})
	//更新gamesrv 彩票设置信息
	beego.Router("/UpdateLotterySettings", &controllers.UpdateLotterySettings{})
	//更新gamesrv 彩票信息
	beego.Router("/UpdateLotteryInfo", &controllers.UpdateLotteryInfo{})
	//服务器启动时更新20所有采种20条历史记录
	beego.Router("/UpdateLotteryRecord", &controllers.UpdateLotteryRecord{})
	//更新一个彩种的历史记录,这个消息用于专门更新历史记录,比如计算服掉期补全时使用
	beego.Router("/UpdateRecordForLostExpect", &controllers.UpdateRecordForLostExpect{})

	//---------------------------------客户端消息-------------------------------------
	//检查gamesrv状态
	beego.Router("/CheckGameSrvStatus", &controllers.CheckGameSrvStatus{})
	//得到所有彩票的设置信息
	beego.Router("/GetLotteriesSettings", &controllers.GetLotteriesSettings{})
	//根据不同的标识得到彩种信息
	beego.Router("/GetLotteryInfo", &controllers.GetLotteryInfo{})
	//得到采种历史记录
	beego.Router("/GetLotteryRecord", &controllers.GetLotteryRecord{})
	//下注(一注中含有多个订单)
	beego.Router("/LotteryBetting", &controllers.LotteryBetting{})
	//智能追号
	beego.Router("/IntelligentTrack", &controllers.IntelligentTrack{})

	//获取订单信息(根据标志来使用不同的方式获取订单)
	beego.Router("/GetOrderInfo", &controllers.GetOrderInfo{})
	//获取订单信息
	beego.Router("/GetOrderInfoByGameTag", &controllers.GetOrderInfoByGameTag{})
	//根据日期获取订单
	beego.Router("/GetOrderInfoByGameTag", &controllers.GetOrderInfoByDay{})
	//获取最新中奖
	beego.Router("/GetNewestWinning", &controllers.GetNewestWinning{})

	//在Loging服，注册或登陆成功以后，获取账号信息
	beego.Router("/Hohenheim", &controllers.GetAccountInfo{})

	//修改帐号密码(不知道以前为什么会把修改登陆密码写在这里....)
	//beego.Router("/ModifyPassword", &controllers.ModifyPassword{})

	//修改账号附加信息
	beego.Router("/ModifyAdditionalInfo", &controllers.ModifyAdditionalInfo{})

	//设定，修改资金密码
	beego.Router("/ModifyMoneyPassword", &controllers.ModifyMoneyPassword{})

	//获取该用户充值渠道信息
	beego.Router("/GetRechargeChannels", &controllers.GetRechargeChannels{})

	//用户银行转账充值
	beego.Router("/BankTransferRecharge", &controllers.BankTransferRecharge{})

	//用户申请提款
	beego.Router("/UserRequestBankWithdrawals", &controllers.UserRequestBankWithdrawals{})

	//用户撤销提款 只有在等待审核时可以撤销
	beego.Router("/CancelBankWithdrawals", &controllers.CancelBankWithdrawals{})

	//充值记录查询
	beego.Router("/GetRechargeRecord", &controllers.GetRechargeRecord{})

	//提款记录查询
	beego.Router("/GetDrawingsRecord", &controllers.GetDrawingsRecord{})

	//稽核消息(一个非常蛋碎的消息,目的是获取很多数据,来进行用户真正能提款的金额计算)
	beego.Router("/GetInspectMoney", &controllers.GetInspectMoney{})

	//生成邀请码
	beego.Router("/CreateInviteCode", &controllers.CreateInviteCode{})

	//得到所有邀请码相关信息
	beego.Router("/GetInviteCodeInfo", &controllers.GetInviteCodeInfo{})

	//修改邀请码信息(目前只有停用邀请码,和修改邀请码备注这两个功能)
	beego.Router("/ChangeInviteCodeInfo", &controllers.ChangeInviteCodeInfo{})

	//按账户名获得代理信息(查自己,查下属代理都是这个消息)
	beego.Router("/GetAgentReport", &controllers.GetAgentReport{})

	//获得所有下属代理商,和玩家账号
	beego.Router("/GetLowerAgentAndUser", &controllers.GetLowerAgentAndUser{})

	//得到指定下级的投注明细,充值明细,提款明细
	beego.Router("/GetLowerInfo", &controllers.GetLowerInfo{})

	//获得活动信息
	beego.Router("/GetActivity", &controllers.GetActivity{})

	//获得某个邀请码下属账号个数
	beego.Router("/GetInviteCodeLowerNum", &controllers.GetInviteCodeLowerNum{})

	//获得公告信息
	beego.Router("/GetAnnouncement", &controllers.GetAnnouncement{})

	//获得充值二维码
	beego.Router("/GetQrCode", &controllers.GetQrCode{})

	//二维码转账充值
	beego.Router("/QrCodeTransferRecharge", &controllers.QrCodeTransferRecharge{})

	//---------------------------来自LoginServer的信息---------------------------
	//玩家预登录 这条消息要加ip来源验证 等完成功能后来加
	beego.Router("/PrePlayer", &controllers.PrePlayer{})

	//---------------------------来自后台ManageServer的信息---------------------------
	//改变改用户充值渠道, 公告信息, 充值二维码等改变信息
	beego.Router("/GBRechargeQudiao", &controllers.ChangeRechargeChannels{})

	//更新二维码信息
	beego.Router("/UpdateQrCode", &controllers.UpdateQrCode{})

	//更新活动信息
	beego.Router("/UpdateActivity", &controllers.UpdateActivity{})
}
