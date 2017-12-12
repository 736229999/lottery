package controllers

import (
	"encoding/json"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"gamesrv/models/AccountMgr"

	"github.com/astaxie/beego"
)

type GetOrderInfo struct {
	beego.Controller
}

func (o *GetOrderInfo) Post() {
	cReq := ReqGetOrderInfo{}
	cResp := GlobalData.RespGetOrderInfo{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &cReq)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(cReq.AccountName)
	if err != nil {
		cResp.Status = 1
	}
	//token 是否正确
	b := accountInfo.VerifyToken(cReq.Token)
	if !b {
		cResp.Status = 2
	}

	//读取订单
	//通过账号名得到这个帐号的订单(注意,之前要验证 帐号和token是否合法才能进行这个查询, 等上线后这个改为在game服务器建立缓存池,不要每次都从数据库读取)
	//默认读取30条 （查询条件为，账号名，彩种，期数,查询类型）
	err = dbmgr.Instance().GetOrderByAccountName(accountInfo.Account_Name, cReq.Skip, 30, cReq.SearchType, &(cResp.Orders))
	if err != nil {
		cResp.Status = 3
		return
	}

	body, err := json.Marshal(cResp)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}

//获得订单信息(查询条件更多,针对具体彩票,具体期数)
type ReqGetOrderInfo struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	SearchType  int    `json:"searchType"` //1,全部-全部 2,全部-中奖 3,全部-待开奖 4,普通-全部 5,普通-中奖 6,普通待开奖, 7追号-全部, 8追号-中奖, 9追号-待开奖
	Skip        int    `json:"skip"`       //跳过多少条(从多少条开始查询)
}
