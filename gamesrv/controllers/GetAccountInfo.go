package controllers

import (
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"encoding/json"

	"github.com/astaxie/beego"
)

type GetAccountInfo struct {
	beego.Controller
}

func (o *GetAccountInfo) Post() {
	req := &GlobalData.ReqGetAccountInfo{}
	resp := &RespGetAccountInfo{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, req)
	if err != nil {
		beego.Emergency("--------------------------------- Marshal Error : ", err, "-----------------------------------")
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil { //1 未找到账号
		resp.Status = 1
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(req.Token)
	if !b { //2 token错误
		resp.Status = 2
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//返回用户信息
	respAccountInfo := RespAccountInfo{}
	respAccountInfo.Status = 0
	respAccountInfo.AccountType = accountInfo.Account_Type
	respAccountInfo.AccountStatus = accountInfo.Account_Status
	respAccountInfo.AccountName = accountInfo.Account_Name
	respAccountInfo.Money = accountInfo.Money
	respAccountInfo.AgentMoney = accountInfo.Agent_Money
	respAccountInfo.Rebate = accountInfo.Rebate
	if accountInfo.Money_Password == "" {
		respAccountInfo.MoneyPassword = false
	} else {
		respAccountInfo.MoneyPassword = true
	}
	respAccountInfo.MobilePhone = accountInfo.Mobile_Phone
	respAccountInfo.QQ = accountInfo.QQ
	respAccountInfo.WeChat = accountInfo.WeChat
	respAccountInfo.WeiBo = accountInfo.WeiBo
	respAccountInfo.Email = accountInfo.Email
	respAccountInfo.Address = accountInfo.Address
	respAccountInfo.BankName = accountInfo.Bank_Name
	respAccountInfo.BankCard = accountInfo.Bank_Card
	respAccountInfo.CardHolder = accountInfo.Card_Holder
	respAccountInfo.BankOfDeposit = accountInfo.Bank_Of_Deposit
	respAccountInfo.Group = accountInfo.Group
	respAccountInfo.IsAgent = accountInfo.Is_Agent
	respAccountInfo.AgentId = accountInfo.Agent_Id
	respAccountInfo.AgentLv = accountInfo.Agent_Lv
	respAccountInfo.LowLvAgentCount = accountInfo.Low_Lv_Agent_Count
	respAccountInfo.LowLvAccountCount = accountInfo.Low_Lv_Account_Count
	respAccountInfo.BelongAgentId = accountInfo.Belong_Agent_Id
	respAccountInfo.BelongAgent = accountInfo.Belong_Agent
	respAccountInfo.InviteCode = accountInfo.Invite_Code
	respAccountInfo.RegistTime = accountInfo.Regist_Time
	respAccountInfo.RegistIp = accountInfo.Regist_Ip
	respAccountInfo.LastLoginTime = accountInfo.Last_Login_Time
	respAccountInfo.LastLoginIp = accountInfo.Last_Login_Ip
	respAccountInfo.TotalBetAmount = accountInfo.Total_Bet_Amount
	respAccountInfo.OnlineDrawingCount = accountInfo.Online_Drawing_Count
	respAccountInfo.BankDrawingCount = accountInfo.Bank_Drawing_Count
	respAccountInfo.RechargeAmount = accountInfo.Recharge_Amount
	respAccountInfo.FavorableAmount = accountInfo.Favorable_Amount
	respAccountInfo.InspectMoney = accountInfo.Inspect_Money
	respAccountInfo.RemainingAgentInviteCodeCount = accountInfo.Remaining_Agent_Invite_Code_Count
	respAccountInfo.RemainingUserInviteCodeCount = accountInfo.Remaining_User_Invite_Code_Count
	respAccountInfo.RemainAgentCount = accountInfo.Remain_Agent_Count

	body, err := json.Marshal(respAccountInfo)
	if err != nil {
		beego.Emergency("--------------------------------- Marshal Error : ", err, "-----------------------------------")
		return
	}
	o.Ctx.Output.Body(body)

}

type RespGetAccountInfo struct {
	Status int `json:"status"`
}

type RespAccountInfo struct {
	Status                        int     `json:"status"`
	AccountType                   int     `json:"accountType"`
	AccountStatus                 int     `json:"accountStatus"` //用户状态 1正常 2冻结
	AccountName                   string  `json:"accountName"`
	Money                         float64 `json:"money"`
	AgentMoney                    float64 `json:"agentMoney"` //代理商余额
	Rebate                        float64 `json:"rebate"`
	MoneyPassword                 bool    `json:"moneyPassword"`                     //是否有资金密码
	MobilePhone                   string  `json:"mobilePhone"`                       //手机号
	QQ                            string  `json:"qq"`                                //QQ号
	WeChat                        string  `json:"weChat"`                            //微信号
	WeiBo                         string  `json:"weiBo"`                             //微博
	Email                         string  `json:"email"`                             //邮箱
	Address                       string  `json:"address"`                           //地址
	BankName                      string  `json:"bankName"`                          //银行地址
	BankCard                      string  `json:"bank_Card"`                         //银行卡号
	CardHolder                    string  `json:"card_Holder"`                       //持卡人
	BankOfDeposit                 string  `json:"bank_Of_Deposit"`                   //开户银行
	Group                         int     `json:"group"`                             //用户组   1没有组
	IsAgent                       int     `json:"isAgent"`                           //是否是代理商
	AgentId                       int     `json:"agentId"`                           //代理商ID
	AgentLv                       int     `json:"agentLv"`                           //代理商层级(最多30级)
	LowLvAgentCount               int     `json:"lowLvAgentCount"`                   //下级代理商数
	LowLvAccountCount             int     `json:"lowLvAccountCount"`                 //下级普通用户数
	BelongAgentId                 int     `json:"belongAgentId"`                     //所属代理商id
	BelongAgent                   string  `json:"belongAgent"`                       //所属代理商账号
	InviteCode                    string  `json:"invite_Code"`                       //邀请码
	RegistTime                    int64   `json:"regist_Time"`                       //用户注册时间
	RegistIp                      string  `json:"regist_Ip"`                         //注册ip
	LastLoginTime                 int64   `json:"last_Login_Time"`                   //上次登录时间
	LastLoginIp                   string  `json:"last_Login_Ip"`                     //上次登录IP
	TotalBetAmount                int     `json:"total_Bet_Amount"`                  //当前投注总额
	OnlineDrawingCount            int     `json:"online_Drawing_Count"`              //在线免费提款次数
	BankDrawingCount              int     `json:"bank_Drawing_Count"`                //银行提款免费次数
	RechargeAmount                int     `json:"recharge_Amount"`                   //总充值金额
	FavorableAmount               int     `json:"favorable_Amount"`                  //总优惠金额
	InspectMoney                  int     `json:"inspect_Money"`                     //稽核金额
	RemainingAgentInviteCodeCount int     `json:"remaining_Agent_Invite_Code_Count"` //剩余可生成代理商邀请码次数 200
	RemainingUserInviteCodeCount  int     `json:"remaining_User_Invite_Code_Count"`  //剩余可生普通用户邀请码次数 200
	RemainAgentCount              int     `json:"remain_Agent_Count"`                //剩余可生成代理商次数 30
}
