package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
)

type ModifyAdditionalInfo struct {
	beego.Controller
}

func (o *ModifyAdditionalInfo) Post() {
	req := GlobalData.ReqModifyAdditionalInfo{}
	resp := GlobalData.ResqModifyAdditionalInfo{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	accountInfo := AccountMgr.AccountInfo{}
	//获取账号信息
	err = accountInfo.Init(req.AccountName)
	if err != nil {
		beego.Emergency(err)
		return
	}

	if !accountInfo.VerifyToken(req.Token) {
		return
	}

	if !accountInfo.VerifyFlag(req.Flag) {
		return
	}

	//验证所有附加信息是否合法
	if len(req.Mobile_Phone) > 30 {
		return
	}

	if len(req.QQ) > 20 {
		return
	}

	if len(req.WeChat) > 30 {
		return
	}

	if len(req.WeiBo) > 30 {
		return
	}

	if len(req.Email) > 30 {
		return
	}

	if req.MoneyPassword == accountInfo.Money_Password {
		//开户行地址 这个要200个字节
		if len(req.Address) > 200 {
			return
		}

		if len(req.Bank_Card) > 25 {
			return
		}

		if len(req.Card_Holder) > 20 {
			return
		}

		//银行名称 这个要30字节
		if len(req.Bank_Name) > 30 {
			return
		}

		//开户银行地址30个字节100个中文
		if len(req.Bank_Of_Deposit) > 100 {
			return
		}
		err = dbmgr.Instance().UpdateAccountAdditionInfo(req)
		if err != nil {
			beego.Error(err)
			resp.Status = 1
		} else {
			resp.Status = 0
		}
	} else if req.MoneyPassword == "" {
		err = dbmgr.Instance().UpdateAccountAdditionInfoNotHaveBank(req)
		if err != nil {
			beego.Emergency(err)
			resp.Status = 1
		} else {
			resp.Status = 0
		}
	} else {
		resp.Status = 3
	}

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}
	o.Ctx.Output.Body(body)
}
