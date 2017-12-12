package controllers

import (
	"gamesrv/models/AccountMgr"
	"gamesrv/models/QrCodeMgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type GetQrCode struct {
	beego.Controller
}

func (o *GetQrCode) Post() {
	req := ReqGetAnnouncement{}
	resp := []QrCodeMgr.QrCode{}

	respStatus := make(map[string]interface{})

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug(err)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil { //1 未找到账号
		respStatus["status"] = 1
		bufres, _ := json.Marshal(respStatus)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(req.Token)
	if !b { //2 token错误
		respStatus["status"] = 2
		bufres, _ := json.Marshal(respStatus)
		o.Ctx.Output.Body(bufres)
		return
	}

	//根据不同得条件来获得充值二维码
	for _, v := range QrCodeMgr.Instance().QrCodes {
		for _, i := range v.Group {
			if i == accountInfo.Group {
				resp = append(resp, v)
			}
		}
	}

	bufres, _ := json.Marshal(resp)
	o.Ctx.Output.Body(bufres)
}

type ReqGetQrCode struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Flag        string `json:"flag"`
}
