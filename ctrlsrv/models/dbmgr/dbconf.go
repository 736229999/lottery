package dbmgr

const (
	//数据库root权限用户名
	//ctrlsrvDbUserName = "e0zs+>v2$-Mt2t:A#HUi"
	ctrlsrvDbUserName = "xq001"
	//数据库root权限密码
	//ctrlsrvDbPwd = "n.6mLarq5!31XOxXP7/D"
	ctrlsrvDbPwd = "123456"

	//控制服 db ip 注意,由于控制服只有一个所以再代码里面写死,而且应该不存在频繁变更的情况
	ctrlsrvDbIP = "192.168.1.151:27017"

	//控制服配置数据库名
	confDbName = "conf"

	//服务器配置表名
	srvColl = "srv"

	//加密信息表名
	encryptColl = "encrypt"

	//彩种信息表名
	ltryColl = "lottery"
)
