package IntelligentTrackBettingMgr

import "sync"

var sInstance *IntelligentTrackBettingMgr
var once sync.Once

func Instance() *IntelligentTrackBettingMgr {
	once.Do(func() {
		sInstance = &IntelligentTrackBettingMgr{}
		sInstance.init()
	})

	return sInstance
}

type IntelligentTrackBettingMgr struct {
	IntelligentTrackInfos map[string][]IntelligentTrackInfo //一个帐号智能追号的所有线
}

func (o *IntelligentTrackBettingMgr) init() {

}

//一条智能追号的信息,一个玩家 ,什么游戏,什么下注方式,多少期,多少钱,多少倍下的什么号,反水多少,,,,,,,追
type IntelligentTrackInfo struct {
	TrackNum        string  `json:"TrackNum`         //追踪号码(追号序列号)
	AccountName     string  `json:"accountName`      //账户名
	GameTag         string  `json:"gameTag"`         //游戏名称
	Expect          []int   `json:"expect"`          //要追那些期
	Multiple        []int   `json:"multiple"`        //有多少要追的期数就有多少倍数
	BetType         int     `json:"bettingType"`     //投注类型
	SingleBetAmount float64 `json:"singleBetAmount"` //单注金额
	BetNums         string  `json:"betNums"`         //投注数字
	Rebate          float64 `json:"rebate"`          //反水
	IsContinue      bool    `json:"isContinue"`      //中奖后是否继续
	OrderNumber     string  `json:"orderNumber"`     //这条线所对应的最新订单号
	Status          int     `json:"status"`          //状态:这条线的状态   0 为还在追号,1为已经中奖停止追号,2为已追完全部订单
	TotalAmount     float64 `json:"totalAmount"`     //状态:这条线的总金额
	ReturnAmount    float64 `json:"returnAmount"`    //返回金额 (注意,只有在状态为1 已中奖停止追号的情况下才有返钱)
}
