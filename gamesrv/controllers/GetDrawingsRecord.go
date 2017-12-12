package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
)

type GetDrawingsRecord struct {
	beego.Controller
}

func (o *GetDrawingsRecord) Post() {
	cReq := ReqGetDrawingsRecord{}
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

	//查询充值记录
	data := []GlobalData.DrawingsRecord{}
	//默认查询30条
	err = dbmgr.Instance().GetDrawingsRecord(accountInfo.Account_Name, cReq.Skip, 30, cReq.Type, &data)
	if err != nil {
		beego.Error(err)
		return
	}

	//计算其他
	// for _, v := range data {
	// 	v.ActualAmount = v.Money - v.CommissionMoney
	// 	v.AccountBalance = v.MoneyBefore - v.Money
	// }
	body, err := json.Marshal(data)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}

type ReqGetDrawingsRecord struct {
	AccountName string `json:"accountName"` //账号名,提款账号
	Token       string `json:"token"`       //token
	Flag        string `json:"flag"`        //flag
	Skip        int    `json:"skip"`        //跳过数目
	Type        int    `json:"type"`        //查询类型 1.全部 2.成功 3.等待
}

// 提款
// 1.提款时间 drawings_time
// 2.账号名 account_name
// 3.流水号 serial_number
// 4.提款状态 status
// 5.提款金额 money
// 6.手续费 commission_money
// 7.出款金额
// 8.提款后账余额
