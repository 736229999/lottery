package AccountMgr

import (
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//帐号信息类
type AccountInfo struct {
	Account_Id                        int
	Account_Type                      int    //用户 0 试玩,1正式, 2测试
	Account_Status                    int    //用户状态 1正常 2冻结
	Account_Name                      string //
	Flag                              string
	Token                             string
	Money                             float64
	Agent_Money                       float64 //代理商余额
	Rebate                            float64
	Money_Password                    string //资金密码
	Mobile_Phone                      string //手机号
	QQ                                string //QQ号
	WeChat                            string //微信号
	WeiBo                             string //微博
	Email                             string //邮箱
	Address                           string //地址
	Bank_Name                         string //银行名称
	Bank_Card                         string //银行卡号
	Card_Holder                       string //持卡人
	Bank_Of_Deposit                   string //开户银行
	Group                             int    //用户组   1没有组
	Is_Agent                          int    //是否是代理商
	Agent_Id                          int    //代理商ID
	Agent_Lv                          int    //代理商层级(最多30级)
	Low_Lv_Agent_Count                int    //下级代理商数
	Low_Lv_Account_Count              int    //下级普通用户数
	Belong_Agent_Id                   int    //所属代理商id
	Belong_Agent                      string //所属代理商账号
	Remark                            string
	Invite_Code                       string //邀请码
	Regist_Time                       int64  //用户注册时间
	Regist_Ip                         string //注册ip
	Last_Login_Time                   int64  //上次登录时间
	Last_Login_Ip                     string //上次登录IP
	Total_Bet_Amount                  int    //当前投注总额(每笔充值之间的投注总额)
	Bet_Amount_Immediate              int    //当前投注总额
	Between_Bet_Amount                int    //每笔充值之间的投注总额
	Online_Drawing_Count              int    //在线免费提款次数
	Bank_Drawing_Count                int    //银行提款免费次数
	Recharge_Amount                   int    //总充值金额
	Favorable_Amount                  int    //总优惠金额
	Inspect_Money                     int    //核查金额
	Remaining_Agent_Invite_Code_Count int    //剩余可生成代理商邀请码次数 200
	Remaining_User_Invite_Code_Count  int    //剩余可生普通用户邀请码次数 200
	Remain_Agent_Count                int    //剩余可生成代理商次数 30
}

//初始化(从数据库获取用户信息)
func (o *AccountInfo) Init(accountName string) error {
	err := dbmgr.Instance().InitAccountInfo(accountName, o)
	if err != nil {
		return err
	}
	return nil
}

//通过accountID初始化一个账号
func (o *AccountInfo) InitByID(accountId int) error {
	err := dbmgr.Instance().InitAccountInfoById(accountId, o)
	if err != nil {
		return err
	}
	return nil
}

//通过AgentId初始化一个代理商账号(从数据库获取用户信息)
func (o *AccountInfo) InitAgent(agentId int) error {
	err := dbmgr.Instance().InitAgentInfo(agentId, o)
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

//验证帐号的flag 是否相同
func (o AccountInfo) VerifyFlag(flag string) bool {
	if o.Flag == flag {
		return true
	}
	return false
}

//验证资金密码是否相同
func (o AccountInfo) VerifyMoneyPassword(moneyPassword string) bool {
	if o.Money_Password == moneyPassword {
		return true
	}
	return false
}

//更新资金密码
func (o AccountInfo) UpdateMoneyPassword() bool {
	err := dbmgr.Instance().UpdateMoneyPassword(o.Account_Name, o.Money_Password)
	if err != nil {
		beego.Emergency(err)
		return false
	}
	return true
}

//更新数据库中用户信息
func (o AccountInfo) Update() bool {
	//更新数据
	selector := bson.M{"account_name": o.Account_Name}

	err := dbmgr.Instance().AccountInfoCollection.Update(selector, o)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
