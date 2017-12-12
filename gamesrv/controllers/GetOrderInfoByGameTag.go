package controllers

import (
	"encoding/json"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"gamesrv/models/AccountMgr"

	"github.com/astaxie/beego"
)

type GetOrderInfoByGameTag struct {
	beego.Controller
}

func (o *GetOrderInfoByGameTag) Post() {
	req := GlobalData.ReqGetOrderInfoByGameTag{}
	resp := GlobalData.RespGetOrderInfoByGameTag{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil {
		resp.Status = 1
	}
	//token 是否正确
	b := accountInfo.VerifyToken(req.Token)
	if !b {
		resp.Status = 2
	}
	//读取订单
	//通过账号名得到这个帐号的订单(注意,之前要验证 帐号和token是否合法才能进行这个查询, 等上线后这个改为在game服务器建立缓存池,不要每次都从数据库读取)
	//默认读取30条 （查询条件为，账号名，彩种，期数）
	err = dbmgr.Instance().GetOrderByAccountNameAndGameTag(accountInfo.Account_Name, req.GameTag, req.Expect, req.Skip, 30, &(resp.Orders))
	if err != nil {
		resp.Status = 3
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}
