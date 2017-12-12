package ltrymgr

import (
	"calculsrv/models/gb"
	"time"
)

/*
彩票接口
*/
type Ltryif interface {
	// //得到游戏标识
	GetGameName() string

	//得到当前期数
	GetCurrentExpect() int

	// //得到当前这期开采号码
	// GetOpenCode() []int

	// //得到下期期数
	// GetNextExpect() int

	//得到下棋开彩时间
	GetNextOpenTime() time.Time

	//得到下期请求时间
	GetNextReqTime() time.Time

	// //得到从api获得的开彩记录数据结构
	// GetNewestRecord() gb.LotteryRecordByNewestFromApi

	// 开彩
	StartLottery(newestRecord gb.LtryRecordByNewest)

	//解析投注(一注中有多个订单)
	AnalyticalBetting(bettingInfo gb.MsgBettingInfo) int

	//更新彩票信息(排序,推荐等)
	UpdateLtryInfo(gb.LotteryInfo) bool

	//更新彩票设置(赔率)
	UpdateLtrySet() bool

	//得到彩票名字
	GetParentName() string

	//结算订单(之所以开放这个接口,是应为目前阶段自动补全机制还不完善,后期完善了以后 要取消这个危险的接口)
	SettlementOrders(orders []gb.Order, openCode string)
}
