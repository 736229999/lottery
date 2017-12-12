package apimgr

const (
	//按天获取历史记录重试次数
	GetRecordOneDayRetryCount = 5
	//按天获取历史记录每次间隔时间(睡眠时间) 秒
	GetRecordOneDaySleepTime = 5

	//按最新(带下期)历史记录重试次数
	GetRecordByNewestRetryCount = 3
	//按最新(带下期)历史记录每次间隔时间(睡眠时间) 秒
	GetRecordByNewestSleepTime = 5

	//监听时间间隔(睡眠时间)
	ListeningSleepTime = 3
)
