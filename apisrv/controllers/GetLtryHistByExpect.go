package controllers

import (
	"apisrv/models/apimgr"
	"apisrv/models/encmgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type GetLtryHistByExpect struct {
	beego.Controller
}

func (o *GetLtryHistByExpect) Post() {
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		beego.Debug(err)
		return
	}

	req := ReqLtryHistByExpect{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	//输出请求服务器IP
	beego.Info("Req Srv Ip : ", o.Ctx.Input.IP())

	resp, err := apimgr.Instance().GetLtryRecordByDate(req.GameName, req.Date)
	if err != nil {
		beego.Error(err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		beego.Debug(err)
		return
	}

	o.Ctx.Output.Body(cipher)
}

type ReqLtryHistByExpect struct {
	GameName string
	Expect   int
}
