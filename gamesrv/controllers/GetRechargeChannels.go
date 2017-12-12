package controllers

import (
	"gamesrv/models/AccountMgr"
	"gamesrv/models/QrCodeMgr"
	"gamesrv/models/RechargeMgr"

	"encoding/json"

	"github.com/astaxie/beego"
)

type GetRechargeChannels struct {
	beego.Controller
}

func (o *GetRechargeChannels) Post() {
	req := ReqGetRechargeChannels{}
	resp := RespGetRechargeChannels{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil {
		beego.Debug(err)
		return
	}

	if !accountInfo.VerifyToken(req.Token) {
		resp.Status = 1 //token 错误
	}

	//获取该帐号的充值权限
	resp.RechargeChannels = RechargeMgr.Instance().GetRechargeChannelsByAccount(accountInfo.Group, req.Devices)

	//返回充值二维码
	//根据不同得条件来获得充值二维码
	for _, v := range QrCodeMgr.Instance().QrCodes {
		for _, i := range v.Group {
			if i == accountInfo.Group {
				resp.QrCodes = append(resp.QrCodes, v)
			}
		}
	}

	body, _ := json.Marshal(resp)
	o.Ctx.Output.Body(body)
}

//客户端请求,获取充值渠道信息
type ReqGetRechargeChannels struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Devices     string "json:devices"
}

//用于回复获取充值渠道信息
type RespGetRechargeChannels struct {
	Status           int                              `json:"status"`           //状态码
	RechargeChannels RechargeMgr.RespRechargeChannels `json:"rechargeChannels"` //充值渠道
	QrCodes          []QrCodeMgr.QrCode               `json:"qrCodes"`          //二维码
}
