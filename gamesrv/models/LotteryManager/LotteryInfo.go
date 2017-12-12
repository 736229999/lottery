package LotteryManager

import "time"

//彩票信息；
type Lottery struct {
	GameTag    string    // 彩票标识
	NextExcept int       // 下一期期次
	OpenCode   string    // 上一期开奖号码
	OpenTime   time.Time // 下一期开奖时间
}
