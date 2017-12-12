package GlobalData

import (
	"gamesrv/models/Order"
)

//Game 向 Login 用公钥
var RsaGameToLoginPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDsQx8pbtf2qsj0a7Y8qCHJ6uYi
NmoPA2ZXRSKOw0mZqkIJTJ3MY3c74/XhMBWc1bsNMBvfKH+w+BCpSRTbrvXQpF9B
5Ks/vRcu192AOCnlFay3FJ6rDk+Zt/GyE1q75+mQIthvbiJY6IEA6kZ1isHw+2nj
27M0slwlWmPIoD8xnwIDAQAB
-----END PUBLIC KEY-----
`)

//Calculation 向 Game 用私钥
var RsaCalToGamePrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDbKQCgjc8S4AnrxW2AmrnZ1lGvSf6me64mPsDy0ZOsluFmEOh1
ul4GzzuP046gzsF2/VMPMeK7EpOy5nfik3khZ/DRhy1pl9CI+6hQO++4P9vZgEko
rJji05CXy/l8pOb3G7E5zob/YMwZya00yLeC8U7NDcTGyA9LrUeNPGKsJwIDAQAB
AoGAOFZm+d7aX2DGTBx5VLjxk6T7ZJMh6gwbLWuuT/099/zlPkaUa0cYSqnVBtj8
biwYIY1nX52USkCxRSjnopAEb+8ukx8vygrjsPRelNFEQHccvGFOHxge/uP15GNe
r6tgRvd8oMnAW7NZKU0sbzwNv5E8DdzC5CTeYqtCzQvuZoECQQDu8D64FYWKP3w5
m80bvwoVQp4vL7Hi6STMd2JHgSiw7Rc+4aMtdqEXOapzZ2mn3TfpT0XXDQBW82vS
uL5E0SZnAkEA6s82rofRabOTmkoBgn1FKM1rykLtO9WYJVNnIIcNRwVPwE9k7w+G
k0my7+OIWf0FUG0SWOlj4w1yFhJh4c00QQJAN1zBG4QZEgRNs0dvSduE6ZIq9sH1
VJ8yoJGU2v6JQB1fQnmjhngzMR9yaBTw/S0btFpi00Y26u6x7/xQUb+oRwJBAM7o
5ZCKEhiIq2psCESKSqUXzwIeU5pEL5vEkL1zBFou7gXScGjQT4/+g3UsFPznHwWt
91rt3p03Pe8BJ8un3YECQQCFeHMjtU3Xns59K+DmzIeUSLQAP2Uf216CUpTCh9EG
e+TG01nrowxe44n3arW2DYT8uZ9xUuJ8zGoYP8yEz9nF
-----END RSA PRIVATE KEY-----
`)

//计算服务器Ip(注意:目前同一时间只会有一台计算服连接到Game 服务器,应为现在还没有自己写API服务器,API有IP访问限制,而且还没有时间来解决,多计算服同时向Game发消息的取舍情况,等后期有时间来做,目前先上线!~~!~)
var CalculationServerIp string

const (
	//预计会启动3个LoginServer (Login服务器是固定的)
	TrialLoginServerIp    = "47.52.61.84:8767"  //试玩服
	FormalLoginServerIp_0 = "47.52.89.242:8767" //正式0
	FormalLoginServerIp_1 = "47.52.99.194:8767" //正式1
)

//服务器注册信息Game 向 login ,  calculation 向 game 都使用这个结构
type ServerRegistInfo struct {
	Id     int    `json:"id"`     //服务器id
	Port   int    `json:"port"`   //服务器端口(注意:这个服务器端口是指,注册服务器用于接收消息的端口)
	Cipher []byte `json:"cipher"` //rsa密文
}

//-------------------------------------------------------------------------------
//预登录请求
type ReqPrepareLogin struct {
	AccountName        string `json:"accountName"`
	Token              string `json:"token"`
	Flag               string `json:"flag"`
	AccountId          int    `json:"accountId"`
	AccountType        int    `json:"accountType"`
	InviteCode         string `json:"inviteCode"`
	RegistTimeStamp    int64  `json:"registTime"`
	RegistIp           string `json:"registIp"`
	LastLoginTimeStamp int64  `json:"lastLoginTimeStamp"`
	LastLoginIp        string `json:"lastLoginIp"`
}

//返回客户端预登陆请求
type ResqPrepareLogin struct {
	Status int `json:"status"`
}

//获取账户信息请求
type ReqGetAccountInfo struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Flag        string `json:"flag"`
}

//返回获得订单信息
type RespGetOrderInfo struct {
	Status int           `json:"status"`
	Orders []Order.Order `json:"orders"`
}

//获得订单信息(查询条件更多,针对具体彩票,具体期数)
type ReqGetOrderInfoByGameTag struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	GameTag     string `json:"gameTag"`
	Expect      int    `json:"expect"`
	Skip        int    `json:"skip"`
}

//返回获得订单信息
type RespGetOrderInfoByGameTag struct {
	Status int           `json:"status"`
	Orders []Order.Order `json:"orders"`
}

//返回 最新中奖信息
type RespGetNewestWinning struct {
	NewestWinningInfos []NewestWinningInfo `json:"newestWinningInfos"`
}

type NewestWinningInfo struct {
	Account_Name string  `json:"accountName"`
	Game_Tag     string  `json:"gameTag"`
	Settlement   float64 `json:"settlement"`
}

type ReqModifyPassword struct {
	Account_Name string `json:"accountName"`
	Token        string `json:"token"`
	Flag         string `json:"flag"`
	OldPassword  string `json:"oldPassword"`
	NewPassword  string `json:"newPassword"`
}

type RespModifyPassword struct {
	Status int `json:"status"`
}

type ReqModifyMoneyPassword struct {
	Account_Name     string `json:"accountName"`
	Token            string `json:"token"`
	Flag             string `json:"flag"`
	OldMoneyPassword string `json:"oldMoneyPassword"`
	NewMoneyPassword string `json:"newMoneyPassword"`
}

type RespModifyMoneyPassword struct {
	Status int `json:"status"`
}

type ReqModifyAdditionalInfo struct {
	AccountName     string `json:"accountName"`
	Token           string `json:"token"`
	Flag            string `json:"flag"`
	MoneyPassword   string `json:"moneyPassword"`
	Mobile_Phone    string `json:"mobilePhone"`   //手机号
	QQ              string `json:"qq"`            //QQ号
	WeChat          string `json:"weChat"`        //微信号
	WeiBo           string `json:"weiBo"`         //微博
	Email           string `json:"email"`         //邮箱
	Address         string `json:"address"`       //地址
	Bank_Name       string `json:"bankName"`      //银行名称
	Bank_Card       string `json:"bankCard"`      //银行卡号
	Card_Holder     string `json:"cardHolder"`    //持卡人
	Bank_Of_Deposit string `json:"bankOfDeposit"` //开户银行 （限定len不能大于200）
}

type ResqModifyAdditionalInfo struct {
	Status int `json:"status"`
}

//-----------------------------------------------------
//用户信息
type AccountInfo struct {
	Increment_Code       int
	Account_Id           int
	Account_Type         int
	Account_Status       int //用户状态 1正常 2冻结
	Account_Name         string
	Flag                 string
	Token                string
	Money                float64
	Rebate               float64
	Money_Password       string
	Mobile_Phone         string //手机号
	QQ                   string //QQ号
	WeChat               string //微信号
	WeiBo                string //微博
	Email                string //邮箱
	Address              string //地址
	Bank_Card            string //银行卡号
	Card_Holder          string //持卡人
	Bank_Name            string //银行名称
	Bank_Of_Deposit      string //开户银行
	Group                int    //用户组   1没有组
	Agent                int    //代理商
	Remark               string
	InviteCode           string //邀请码
	Regist_Time          int64  //用户注册时间
	Regist_Ip            string //注册ip
	Last_Login_Time      int64  //上次登录时间
	Last_Login_Ip        string //上次登录IP
	Total_Bet_Amount     int    //当前投注总额
	Online_Drawing_Count int    //在线免费提款次数
	Bank_Drawing_Count   int    //银行提款免费次数
	Recharge_Amount      int    //总充值金额
	Favorable_Amount     int    //总优惠金额
	Inspect_Money        int    //稽核金额
}

//用户线下转账充值
type ReqBankTransfer struct {
	AccountName   string `json:"accountName"`
	Flag          string `json:"flag"`
	MPW           string `json:mpw` //资金密码
	Token         string `json:"token"`
	RemitAmount   string `json:"remitAmount"`   //存款金额
	BankId        string `json:"bankId"`        //存款银行ID
	Remitter      string `json:"Remitter"`      //存款人姓名
	RemitBankName string `json:"remitBankName"` //存款银行
	RemitBankCard string `json:"remitBankCard"` //存款卡号
	RemitType     string `json:"remitType"`     //存款方式
	RemitTime     string `json:"remitTime"`     //存款时间
}

//充值记录
type RechargeRecord struct {
	RechargeTime int            `bson:"recharge_time"` //1.充值时间
	AccountName  string         `bson:"account_name"`  //2.账号名
	OrderNumber  string         `bson:"order_number"`  //3.充值订单号
	SerialNumber string         `bson:"serial_number"` //4.序列号
	Money        int            `bson:"money"`         //5.充值金额
	Favorable    int            `bson:"favorable"`     //6.优惠金额
	RechargeInfo RechargeDetail `bson:"recharge_info"` //7.充值详情
	Status       int            `bson:"status"`        //8.订单状态
}

//充值详情
type RechargeDetail struct {
	ThirdPlatform   string `bson:"third_platform"`    //第三方平台
	ThirdMchid      string `bson:"third_mchid"`       //第三方mchid
	ThirdType       string `bson:"third_type"`        //第三方类型
	OnlinePaymentId string `bson:"online_payment_id"` //第三方付款id
	RechargeTime    string `bson:"recharge_time"`     //充值时间
}

//提款记录
type DrawingsRecord struct {
	DrawingsTime    int     `bson:"drawings_time"`    //1提款时间
	AccountName     string  `bson:"account_name"`     //2.账号名
	SerialNumber    string  `bson:"serial_number"`    //4.序列号
	Status          int     `bson:"status"`           //8.提款状态
	Money           float64 `bson:"money"`            //5.提款金额
	CommissionMoney float64 `bson:"commission_money"` //6.手续费
	MoneyBefore     float64 `bson:"money_before"`     //提款之前金额
	//ActualAmount    float64 `bson:"ActualAmount"`     //实际到账金额  = 提款金额 - 手续费
	//AccountBalance  float64 `bson:"AccountBalance"`   //提款后账户余额 = 提款之前金额 - 提款金额
}

type InspectMoneyBank struct {
	Drawings DrawingsBank `bson:"drawings"`
}

type InspectMoneyOnline struct {
	Drawings DrawingsOnline `bson:"drawings"`
}

type DrawingsBank struct {
	InspectDetail InspectDetailStr `bson:"bank"`
}

type DrawingsOnline struct {
	InspectDetail InspectDetailStr `bson:"online"`
}

type InspectDetailStr struct {
	CommissionPercent       float64 `bson:"commissionPercent"`       //5.当前用户组手续费费率
	CommissionStatus        int     `bson:"commissionStatus"`        //7.当前用户组是否开启免费
	Minimum                 int     `bson:"minimum"`                 //9.当前用户组提款最低值
	Maximum                 int     `bson:"maximum"`                 //10.当前用户组提款最高值
	CommissionMinimum       int     `bson:"commissionMinimum"`       //11.手续费最低值
	CommissionMaximum       int     `bson:"commissionMaximum"`       //12.手续费最高值
	NormalityInspectBroaden int     `bson:"normalityInspectBroaden"` //13.常态性核查放宽额度
	NormalityInspectRate    int     `bson:"normalityInspectRate"`    //14.行政费率
	CommissionCount         int     `bson:"commissionCount"`         //15.免费提款次数
}

//邀请码信息(用于用户注册)
type InviteCodeRelatedInfo struct {
	Type    int     `bson:"type"`     //邀请码可生成用户类型,0为普通用户,1为代理商
	Status  int     `bson:"status"`   //验证码状态 0为不可用 1为可用
	AgentId int     `bson:"agent_id"` //如果是代理商,这个值就是用户表中的Account_Id
	Rebate  float64 `bson:"rebate"`   //反水
	Code    string  `bson:"code"`     //邀请码
}

//平台相关数据信息
type PlatformConf struct {
	Id                            int     `bson:"id"`
	Name                          string  `bson:"name"`
	Status                        int     `bson:"status"`                            //平台状态 1.开启  0.未开启
	Rebate                        float64 `bson:"rebate"`                            //平台反水
	RemainAgentCount              int     `bson:"remain_agent_count"`                //可生成代理商次数
	RemainingAgentInviteCodeCount int     `bson:"remaining_agent_invite_code_count"` //可生成代理商邀请码次数
	RemainingUserInviteCodeCount  int     `bson:"remaining_user_invite_code_count"`  //可生成普通用于邀请码次数
}

//所有彩票统一的彩票设置结构
type LotterySettings struct {
	Name        string
	Id          int                //odds mode
	OddsMap     map[string]float64 //赔率数组应为玩法的特殊性,导致有些玩法有多个赔率(比如快3的和值) key 为 赔率id 例如快3 和值的赔率id是 3-18
	SingleLimit float64            //注单限额
	OrderLimit  float64            //订单限额
}
