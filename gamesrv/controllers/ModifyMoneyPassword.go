package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"regexp"

	"github.com/astaxie/beego"
)

type ModifyMoneyPassword struct {
	beego.Controller
}

func (o *ModifyMoneyPassword) Post() {
	req := GlobalData.ReqModifyMoneyPassword{}
	resp := GlobalData.RespModifyMoneyPassword{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	//检查账户
	accountInfo := AccountMgr.AccountInfo{}
	accountInfo.Init(req.Account_Name)

	if !accountInfo.VerifyToken(req.Token) {
		beego.Debug(1)
		resp.Status = 1
	}

	if !accountInfo.VerifyFlag(req.Flag) {
		beego.Debug(2)
		resp.Status = 2
	}

	//如果newPassword 没有字段,就认为这时第一次设定资金密码
	if req.NewMoneyPassword == "" {
		if !VerifMoneyPassword(req.OldMoneyPassword) {
			beego.Debug(req.OldMoneyPassword)
			return
		}
		//改变资金密码
		accountInfo.Money_Password = req.OldMoneyPassword
		if !accountInfo.UpdateMoneyPassword() {
			beego.Debug(3)
			resp.Status = 3
		}
	} else {
		//验证老密码是否相等
		if req.OldMoneyPassword != accountInfo.Money_Password {
			beego.Debug(4)
			resp.Status = 4
		} else {
			if !VerifMoneyPassword(req.NewMoneyPassword) {
				return
			}
			accountInfo.Money_Password = req.NewMoneyPassword
			if !accountInfo.UpdateMoneyPassword() {
				resp.Status = 3
			}
		}
	}

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Emergency(err)
		return
	}
	o.Ctx.Output.Body(body)
}

//验证密码合法
func VerifMoneyPassword(password string) bool {
	if password == "" {
		return false
	}

	l := len(password)
	if l != 6 {
		return false
	}

	match, err := regexp.MatchString("^[0-9]*$", password)
	if match == false || err != nil {
		return false
	}
	return true
}
