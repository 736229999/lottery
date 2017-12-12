package dbmgr

const (
	CalculationDbUserName = "$Hy_Db_Calculation#520017*"
	CalculationDbPwd      = "h657YHG&^*rH&*psz1/F6"

	//计算服DbNmae
	CalculationDbName = "CalculationServer"

	//用户信息表名
	AccountInfoCollection = "account_info"
	//历史记录表名
	HistoryCollection = "history"
	//号码走势表名
	TrendCollection = "trend"
	//下注信息表名;
	BetCollection = "bet"
	//订单信息表名
	OrderCollection = "order"
)

const (
	ManageDbUserName = "ManageServer"
	ManageDbPwd      = "manage123server"

	//管理员后台DbNmae
	ManageDbName = "ManageServer"

	//流水记录表名
	BalanceRecordCollection = "balance_record"
	//彩票管理表名
	LotteryTypeCollection = "lottery_type"
	//Game服务器表名
	GameServerCollection = "server"
	//彩票设置表名
	LotteriesSettingsCollection = "lottery_setting"
)

//正式计算服DB
//const FormalCalculationServer = "47.52.112.4:28018"

//正式管理服DB
//const FormalManageServer = "47.52.65.93:27017"

//试玩服计算服DB
//const TrialCalculationServer = "47.52.61.84:28018"

//试玩服管理后台DB(试玩服 还是要获取正式管理服的数据,但是不写入数据)
//const TrialManagerServer = "47.52.65.93:27017"

//const AnHuiK3Url = "http://f.apiplus.cn/ahk3-20.json"
