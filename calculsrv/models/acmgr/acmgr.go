package acmgr

import (
	"calculsrv/models/dbmgr"
	"calculsrv/models/gb"
	"sync"

	"github.com/astaxie/beego"
)

//帐号信息类
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
	Bank_Of_Deposit      string //开户银行
	Group                int    //用户组   1没有组
	Agent                int    //代理商
	Remark               string
	Regist_Time          int64   //用户注册时间
	Total_Bet_Amount     float64 //总下注金额
	Bet_Amount_Immediate float64 //及时有效投注

	lock sync.Mutex
}

//初始化(从数据库获取用户信息)
func (o *AccountInfo) Init(accountName string) error {
	err := dbmgr.InitAccountInfo(accountName, o)
	if err != nil {
		return err
	}
	return nil
}

//验证帐号的token 是否相同
func (o AccountInfo) VerifyToken(token string) bool {
	if o.Token == token {
		return true
	}
	return false
}

//验证帐号上有没有这么多钱钱,随便就减了 哈哈哈！~
func (o *AccountInfo) DeductMoney(money float64) bool {
	o.lock.Lock()
	defer o.lock.Unlock()
	//这里注意(下注金额必须是大于0,并且玩家身上的金额要大于下注金额)
	if money <= 0 {
		beego.Error("--------------------- Emergency ! ! ! 服务器消息被破解 ,收到负数的下注信息 !!! ---------------------")
		return false
	}

	if o.Money >= money {
		o.Money -= money
		return true
	}
	return false
}

//加钱钱
func (o *AccountInfo) AddMoney(money float64) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Money += money
}

//更新Db中的数据 需要更新的内容为钱钱
func (o AccountInfo) UpdataDb() error {
	o.lock.Lock()
	defer o.lock.Unlock()
	err := dbmgr.UpdateAccountInfoMoney(o.Account_Name, o.Money, o.Total_Bet_Amount, o.Bet_Amount_Immediate)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------全局功能-----------------------------------------------

//用户缓存池 等功能玩成后来写 目前还没有想出 当用户重新登录后怎么通知缓存池的方法,只能先查询数据库
var AccountPool map[string]gb.AccountInfo = make(map[string]gb.AccountInfo)
