package controllers

import (
	"ctrlsrv/models/encmgr"
	"ctrlsrv/models/ltrymgr"
	"ctrlsrv/models/srvmgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

//只提供给Api服务器获取彩票信息,应为其他每一个代理商的彩票信息设定都可能是不一样的
type GetLtryInfo struct {
	beego.Controller
}

func (o *GetLtryInfo) Post() {
	//每条来自其他服务器的消息都要验证来源IP是否正确
	if !srvmgr.Instance().VerifySrv(o.Ctx.Input.IP()) {
		beego.Error("Wrongful server access ! IP : ", o.Ctx.Input.IP())
		return
	}

	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		beego.Debug(err)
		return
	}

	req := ReqGetLtryInfo{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	var resp interface{}
	switch req.Func {
	case "api":
		resp = o.getLtryInfo()
	default:
		beego.Error("There is no server for this function !")
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		beego.Debug(data)
		return
	}

	body, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		beego.Error(err)
		return
	}

	o.Ctx.Output.Body(body)
}

func (o GetLtryInfo) getLtryInfo() []RespLtry {
	resp := []RespLtry{}
	for _, v := range ltrymgr.Instance().Ltrys {
		//状态为正常的彩种才返回
		if v.Status == 1 {
			l := RespLtry{}
			l.Id = v.Id
			l.Name = v.Name
			l.Game_name = v.Game_name
			l.Freq = v.Freq
			l.Api_code_kcw = v.Api_code_kcw
			resp = append(resp, l)
		}
	}

	return resp
}

type ReqGetLtryInfo struct {
	Func string
}

type RespLtry struct {
	Id           int
	Name         string //彩票中文名
	Game_name    string //彩票英文名
	Freq         int    //频率: 0.低频 1.高频
	Api_code_kcw string //kcw 代表这个是开采网的Api彩票请求代码,这个值是要组合在url中使用的
}
