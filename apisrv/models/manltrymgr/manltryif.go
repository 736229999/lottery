package manltrymgr

//手动开奖的彩票接口

type ManlLtryif interface {
	//得到游戏标识
	GetGameName() string

	//得到当前期数
	//GetCurrentExpect() int

	//得到当前这期开采号码
	//GetOpenCode() []int

	//得到下期期数
	GetNextExpect() int

	//得到下棋开彩时间
	//GetNextOpenTime() time.Time

	// 开彩
	StartLtry(expect int, openCode string, openTime string, nextExpect int, nextOpenTime string)

	//解析投注(一注中有多个订单)
	//AnalyticalBetting(bettingInfo gb.MsgBettingInfo) int
}
