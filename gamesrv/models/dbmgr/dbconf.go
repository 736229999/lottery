package dbmgr

const (
	//-----------------------------------------------正式服信息-------------------------------------------------------
	//-----------------------------------------------Combodia 计算服 ----------------------------------------------------
	//DbUrl_CalculationServer = "47.52.112.4:28018"
	CalculationDbUserName = "$Hy_Db_Calculation#520017*"
	CalculationDbPwd      = "h657YHG&^*rH&*psz1/F6"
	//计算服DbNmae
	CalculationDbName = "CalculationServer"

	//账户信息表名
	accountInfoCollection = "account_info"
	//订单信息表名
	OrderCollection = "order"
	//自增表名(所有需要自增数的都放在这里面)
	IncrementIdCollection = "inc_id"

	//----------------------------------------------- 管理服务 -----------------------------------------------------
	//DbUrl_ManageServer = "47.52.65.93:27017"
	ManageDbUserName = "ManageServer"
	ManageDbPwd      = "manage123server"

	//管理数据库名
	ManageDbName = "ManageServer"
	//线下银行转账表名
	TransferBankCardCollection = "transfer_bank_card"
	//在线支付渠道表名
	OnlinePaymentCollection = "online_payment"
	//充值订单表名
	RechargeOrderCollection = "recharge_order"
	//提款订单表名
	DrawingsOrderCollection = "drawings_order"
	//用户组表名
	PlayerGroupCollection = "player_group"
	//支付类型表名
	PayTypeCollection = "pay_type"
	//邀请码表
	InvitationCode = "invitation_code"
	//平台相关信息表(注意这个表以后就会作为配置文件一样的存在,非常重要,以后所用与平台相关的配置都要放这里面)
	PlatformConf = "platform_conf"
	//代理商统计报表(按小时)
	AgentCountHour = "agent_count_hour"
	//代理商统计报表(按天)
	AgentCountDay = "agent_count_day"
	//代理商统计报表(按月)
	AgentCountMonth = "agent_count_month"
	//活动信息
	Activity = "activity"
	//公告信息
	Announcement = "announcement"
	//充值二维码
	QrCode = "qr_code"
	//-------------------------------------------------- 试玩服信息 ------------------------------------------------------
	DbUrl_TrialServer = "47.52.61.84:28018"
	TrialDbUserName   = "$Hy_Db_Manager#170623*"
	TrialDbPwd        = "23Dcypsz1/2jss1/2j#f80"
)
