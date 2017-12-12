package Order

var OrderMapByAccountName map[string][]Order = make(map[string][]Order)

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
