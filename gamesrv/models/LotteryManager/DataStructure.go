package LotteryManager

import (
	"time"
)

//彩票信息,由计算服发送而来
type LotteryInfo struct {
	Id            int    `json:"id"`             //彩票ID
	GameTag       string `json:"gameTag"`        //游戏标志（游戏名称）
	ParentName    string `bson:"parent_name"`    //父类名字(游戏大类型)
	Frequency     string `json:"frequency"`      //是否是高频彩 high 高 low 低
	Status        int    `json:"status"`         //游戏状态 0 关闭，1 正常，2 维护
	Recommend     int    `json:"recommend"`      //是否推荐 数字越高越推荐
	RecommendSort int    `bson:"recommend_sort"` //推荐排序(数字越大,越推荐)
	Sort          int    `json:"sort"`           //彩票排序

	CurrentExpect      int       `json:"currentExpect"`      //当前期数
	OpenCode           []int     `json:"openCode"`           //当期开奖code
	OpenCodeString     string    `json:"openCodeString"`     //当期开奖code string形式
	CurrentOpenTime    time.Time `json:"currentOpenTime"`    //当期开奖时间
	NextExpect         int       `json:"nextExpect"`         //下期期数
	NextRequestTime    time.Time `json:"nextRequestTime"`    //下期请求时间
	NextOpenTime       time.Time `json:"nextOpenTime"`       //下期开彩时间(官方开奖整点)
	NextClosingBetTime time.Time `json:"nextClosingBetTime"` //下期截至下注时间（下期开彩时间 - 1m）

	AfterNextExpect         int       `json:"afterNextExpect"`         //下下期期数
	AfterNextClosingBetTime time.Time `json:"afterNextClosingBetTime"` //下下期截至下注时间

	Settings map[int]LotterySettings `json:"settings"` //游戏设置(赔率，限额等信息)

	ServerNowTime time.Time //服务器现在时间

	//RemainingOpenTime time.Duration `json:remainingOpenTime` //剩余开奖时间
}

type LotteryRecord struct {
	GameName           string    `bson:"game_name"`            // 名称
	Expect             int       `bson:"expect"`               // 期次
	OpenCode           string    `bson:"open_code"`            // 开奖号码
	OpenTime           time.Time `bson:"open_time"`            // 开奖时间
	OpenTimeStamp      int64     `bson:"open_Time_Stamp"`      //开奖时间戳
	RecordingTimeStamp int64     `bson:"recording_Time_stamp"` //记录时间
}

//彩票设定信息(赔率, 限额等,这是一个玩法的信息)(从管理员后台数据库获取)
type LotterySettings struct {
	Id          int             `bson:"odds_mode"`    //玩法ID
	OddsMap     map[int]float64 `bson:"odds_value"`   //游戏赔率
	SingleLimit int             `bson:"quota_single"` //注单限额
	OrderLimit  int             `bson:"quota_bet"`    //订单限额
}
