package gb

import (
	"time"
)

//计算服 注册Game结构
//服务器注册信息Game 向 login ,  calculation 向 game 都使用这个结构
type ServerRegistInfo struct {
	Id     int    `json:"id"`     //服务器id
	Port   int    `json:"port"`   //服务器端口(注意:这个服务器端口是指,注册服务器用于接收消息的端口)
	Cipher []byte `json:"cipher"` //rsa密文
}

//--------------------------------------------------------------------------------
//Game服务器信息(从后台服务器获取)
type GameServerInfo struct {
	Id     int    `bson:"id"`          //服务器ID
	Name   string `bson:"server_name"` //服务器名字
	Ip     string `bson:"server_ip"`   //服务器IP
	Port   string `bson:"port"`        //端口
	Type   int    `bson:"type"`        //服务器类型, 0 试玩  ,1 正式
	Status int    `bson:"status"`      //Game服务器状态 1 正常 2错误
}

//彩票信息(从后台服务器获取)
type LotteryInfo struct {
	Id            int    `bson:"id"`             //彩票id
	GameTag       string `bson:"game_tag"`       //游戏标志（游戏名称）
	ParentName    string `bson:"parent_name"`    //父类名字(游戏大类型)
	Frequency     string `bson:"frequency"`      //是否是高频彩 high 高 low 低
	Status        int    `bson:"status"`         //游戏状态 0 关闭，1 正常，2 维护
	Recommend     int    `bson:"recommend"`      //是否推荐 1 推荐，2 不推荐
	RecommendSort int    `bson:"recommend_sort"` //推荐排序(数字越大,越推荐)
	Sort          int    `bson:"sort"`           //彩票排序
}

//一条开奖记录(从计算服务器数据库获取,存数据库也是这个结构)
type LotteryRecord struct {
	GameName           string    `bson:"game_name"`            // 名称
	Expect             int       `bson:"expect"`               // 期次
	OpenCode           string    `bson:"open_code"`            // 开奖号码
	OpenTime           time.Time `bson:"open_time"`            // 开奖时间
	OpenTimeStamp      int64     `bson:"open_time_stamp"`      //开奖时间(时间戳)
	RecordingTime      time.Time `bson:"recording_time"`       //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	RecordingTimeStamp int64     `bson:"recording_time_stamp"` //记录时间(时间戳)
	//maxNum    int       `bson:"id"`       				   // 最大号码取值
}

//所有彩票统一的彩票设置结构
type LotterySettings struct {
	Name        string
	Id          int                //odds mode
	OddsMap     map[string]float64 //赔率数组应为玩法的特殊性,导致有些玩法有多个赔率(比如快3的和值) key 为 赔率id 例如快3 和值的赔率id是 3-18
	SingleLimit float64            //注单限额
	OrderLimit  float64            //订单限额
}

//---------------------------------------------------------------------------------
//从开采网API 获取的按天查询历史记录
type LotteryRecordByDayFromApi struct {
	Rows   int                    `json:"rows"`
	Code   string                 `json:"code"`
	Remain string                 `json:"remain"`
	Data   []LotteryRecordFromApi `json:"data"`
}

//开采网 api 一条记录
type LotteryRecordFromApi struct {
	Expect        string `json:"expect"`
	Opencode      string `json:"opencode"`
	OpenTime      string `json:"opentime"`
	OpenTimeStamp int64  `json:"opentimestamp"`
}

//开采网 按最新获取一掉记录(带下期)
type LotteryRecordByNewestFromApi struct {
	Rows   int                           `json:"rows"`
	Code   string                        `json:"code"`
	Remain string                        `json:"remain"`
	Next   []LotteryRecordNextFromApi    `json:"next"` //下一期信息(下一期期数, 下一期开奖时间)
	Open   []LotteryRecordCurrentFromApi `json:"open"` //最新一期信息 (注意:这里面没有时间戳,存库的时候要自己添加一个)
	Time   string                        `json:"time"` //查询时间
}

//从 api服务器 获取的最新带下一期的记录
type LtryRecordByNewest struct {
	GameName        string
	CurrentExpect   int
	OpenCode        []int
	OpenCodeStr     string
	CurrentOpenTime time.Time

	NextExpect   int
	NextOpenTime time.Time
}

//开采网 下一期开采信息
type LotteryRecordNextFromApi struct {
	Expect   string `json:"expect"`
	OpenTime string `json:"opentime"`
}

//开采网 当前这期信息
type LotteryRecordCurrentFromApi struct {
	Expect   string `json:"expect"`
	Opencode string `json:"opencode"`
	OpenTime string `json:"opentime"`
}

//---------------------------------------------------------------------------
//消息结构 ----------------------------------------------------------------------------------------------
//更新Game服务器彩票信息
type MsgLotteryInfo struct {
	Id            int    `json:"id"`             //彩票ID
	GameTag       string `json:"gameTag"`        //游戏标志（游戏名称）
	ParentName    string `bson:"parent_name"`    //父类名字(游戏大类型)
	Frequency     string `json:"frequency"`      //是否是高频彩 high 高 low 低
	Status        int    `json:"status"`         //游戏状态 0 关闭，1 正常，2 维护
	Recommend     int    `json:"recommend"`      //是否推荐 数字越高越推荐
	RecommendSort int    `bson:"recommend_sort"` //推荐排序(数字越大,越推荐)
	Sort          int    `json:"sort"`

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

	//RemainingOpenTime time.Duration `json:remainingOpenTime` //剩余开奖时间
}

//----------------------------------------------------由GameServer发来的信息--------------------------------------
//订单信息（一次投注）
type MsgBettingInfo struct {
	AccountName string     `json:"accountName"` //用户名
	Token       string     `json:"token"`
	GameTag     string     `json:"gameTag"` //游戏名称
	Expect      int        `json:"expect"`  //对应彩票期数
	Orders      []MsgOrder `json:"orders"`  //订单信息
}

//投注返回
type MsgRespBetting struct {
	Status int `json:"status"` //状态码
}

//单个订单信息
type MsgOrder struct {
	BetType         int     `json:"bettingType"`     //投注类型
	SingleBetAmount float64 `json:"singleBetAmount"` //单注金额
	BetNums         string  `json:"betNums"`         //投注数字
	Rebate          float64 `json:"rebate"`          //反水
}

//智能追号
type MspIntelligentTrackBetting struct {
	AccountName string       `json:"accountName"` //用户名
	Token       string       `json:"token"`       //Token
	TrackOrders []TrackOrder `json:"trackOrders"` //智能追号订单组
}

//追号信息
type TrackOrder struct {
	GameTag         string  `json:"gameTag"`         //游戏名称
	Expect          []int   `json:"expect"`          //要追那些期
	Multiple        []int   `json:"multiple"`        //有多少要追的期数就有多少倍数
	BetType         int     `json:"bettingType"`     //投注类型
	SingleBetAmount float64 `json:"singleBetAmount"` //单注金额
	BetNums         string  `json:"betNums"`         //投注数字
	Rebate          float64 `json:"rebate"`          //反水
	IsContinue      bool    `json:"isContinue"`      //中奖后是否继续
}

type RespIntelligentTrackBetting struct {
	Status int `json:"status"` //状态
}

//----------------------------------------------------
//根据GameServer的信息生成的订单信息用于存库
type Order struct {
	OrderNumber     string             `bson:"order_number"`      //订单号
	OrderType       int                `bson:"order_type"`        //订单类型,0 普通，1追号
	Status          int                `bson:"status"`            //订单状态,0 未结算, 1结算
	AccountName     string             `bson:"account_name"`      //玩家帐号
	GameTag         string             `bson:"game_tag"`          //游戏tag
	Expect          int                `bson:"expect"`            //游戏期数
	BetType         int                `bson:"bet_type"`          //投注类型
	SingleBetAmount float64            `bson:"single_bet_amount"` //单注金额
	SingleBetNum    int                `bson:"single_bet_num"`    //注数
	OrderAmount     float64            `bson:"order_amount"`      //订单金额(单注金额 * 注数)
	BetNums         string             `bson:"bet_nums"`          //投注数字
	Odds            map[string]float64 `bson:"odds"`              //赔率map key 为数据库中读取的赔率id
	Rebate          float64            `bson:"rebate"`            //反水
	OpenCode        string             `bson:"open_code"`         //开奖号码
	RebateAmount    float64            `bson:"rebate_amount"`     //反水金额
	Settlement      float64            `bson:"settlement"`        //输赢总结算
	WinningBetNum   int                `bson:"winning_bet_num"`   //中奖注数
	BettingTime     int64              `bson:"betting_time"`      //下注时间
}

//用户信息
type AccountInfo struct {
	Increment_Code  int
	Account_Id      int
	Account_Type    int
	Account_Status  int //用户状态 1正常 2冻结
	Account_Name    string
	Flag            string
	Token           string
	Money           float64
	Rebate          float64
	Money_Password  string
	Mobile_Phone    int64  //手机号
	QQ              int64  //QQ号
	WeChat          string //微信号
	WeiBo           string //微博
	Email           string //邮箱
	Address         string //地址
	Bank_Card       int64  //银行卡号
	Card_Holder     string //持卡人
	Bank_Name       string //银行名称
	Bank_Of_Deposit string //开户银行
	Group           int    //用户组   1没有组
	Agent           int    //代理商
	Remark          string
	Regist_Time     int64 //用户注册时间
}
