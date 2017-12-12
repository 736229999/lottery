package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"
	"gamesrv/models/GlobalData"

	"github.com/astaxie/beego"
)

type GetRechargeRecord struct {
	beego.Controller
}

func (o *GetRechargeRecord) Post() {
	cReq := ReqGetRechargeRecord{}
	cResp := make(map[string]interface{})
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &cReq)
	if err != nil {
		beego.Debug(err)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(cReq.AccountName)
	if err != nil { //1 未找到账号
		cResp["status"] = 1
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}
	//验证token
	b := accountInfo.VerifyToken(cReq.Token)
	if !b { //9 token错误
		cResp["status"] = 9
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}
	//判断查询类型
	if cReq.Type < 1 || cReq.Type > 3 {
		return
	}

	//查询充值记录
	data := []GlobalData.RechargeRecord{}
	//默认查询30条
	err = dbmgr.Instance().GetRechargeRecord(accountInfo.Account_Name, cReq.Skip, 30, cReq.Type, &data)
	if err != nil {
		beego.Error(err)
		return
	}

	body, err := json.Marshal(data)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}

type ReqGetRechargeRecord struct {
	AccountName string `json:"accountName"` //账号名,提款账号
	Token       string `json:"token"`       //token
	Flag        string `json:"flag"`        //flag
	Skip        int    `json:"skip"`        //跳过数目
	Type        int    `json:"type"`        //查询类型 1.全部 2.成功 3.等待
}
