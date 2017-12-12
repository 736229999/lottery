package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"

	"github.com/astaxie/beego"
)

type ModifyPassword struct {
	beego.Controller
}

func (o *ModifyPassword) Post() {
	req := GlobalData.ReqModifyPassword{}
	resp := GlobalData.RespModifyPassword{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}
	//检查账户
	accountInfo := AccountMgr.AccountInfo{}
	accountInfo.Init(req.Account_Name)

	if !accountInfo.VerifyToken(req.Token) {
		resp.Status = 1
	}

	if !accountInfo.VerifyFlag(req.Flag) {
		resp.Status = 2
	}

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Emergency(err)
		return
	}
	o.Ctx.Output.Body(body)
}
