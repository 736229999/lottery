package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

type GetLowerAgentAndUser struct {
	beego.Controller
}

func (o *GetLowerAgentAndUser) Post() {
	req := ReqLowerAgentAndUser{}
	resp := RespLowerAgentAndUserStatus{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug(err)
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

	//验证是否代理商
	if accountInfo.Is_Agent != 1 {
		resp.Status = 3 //用户不是代理商
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//查询这个账号所有的下级
	bm := bson.M{"belong_agent_id": req.SearchIdAccountId}
	ret := []RespLowerAgentAndUser{}
	dbmgr.Instance().AccountInfoCollection.Find(bm).All(&ret)

	bufres, _ := json.Marshal(ret)
	o.Ctx.Output.Body(bufres)
	return
}

//客户端请求结构
type ReqLowerAgentAndUser struct {
	AccountName       string `json:"accountName"`
	SearchIdAccountId int    `json:"searchIdAccountId"`
	Token             string `json:"token"`
	Flag              string `json:"flag"`
}

//错误验证码返回
type RespLowerAgentAndUserStatus struct {
	Status int `json:"status"`
}

//下级信息
type RespLowerAgentAndUser struct {
	AccountID          int     `bson:"account_id"`
	AccountName        string  `bson:"account_name"`
	Money              float64 `bson:"money"`           //钱
	Rebate             float64 `bson:"rebate"`          //反水
	IsAgent            int     `bson:"is_agent"`        //是否是代理商
	AgentLv            int     `bson:"agent_lv"`        //代理商层级(最多30级)
	LastLoginTime      int64   `bson:"last_login_time"` //上次登录时间
	RegistTime         int64   `bson:"regist_time"`
	InspectMoney       int     `bson:"inspect_money"`        //稽核金额
	TotalBetAmount     int     `bson:"total_bet_amount"`     //当前投注总额
	CardHolder         string  `bson:"card_holder"`          //持卡人
	BetAmountImmediate int     `bson:"bet_amount_immediate"` //及时有效投注
}
