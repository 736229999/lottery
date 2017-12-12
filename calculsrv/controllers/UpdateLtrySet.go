package controllers

import (
	"calculsrv/models/ltrymgr"
	"calculsrv/models/ltryset"
	"encoding/json"

	"github.com/astaxie/beego"
)

type UpdateLtrySet struct {
	beego.Controller
}

func (o *UpdateLtrySet) Post() {
	req := ReqUpdateLtrySet{}
	resp := RespUpdateLtrySet{}

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
		beego.Error(err)
		return
	}

	beego.Warn("收到 UpdateLtrySet 消息 : ", req.GameParentName)

	//重新设置彩票设置(赔率)
	err := ltryset.Init()
	if err != nil {
		beego.Error(err)
		resp.Status = 1
		body, err := json.Marshal(resp)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(body)
		return
	}

	for _, v := range ltrymgr.Instance().LtryifMap {
		if v.GetParentName() == req.GameParentName {
			if !v.UpdateLtrySet() {
				resp.Status = 2
			}
		}
	}

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}
	o.Ctx.Output.Body(body)
}

type ReqUpdateLtrySet struct {
	GameParentName string `json:"gameParentName"`
}

type RespUpdateLtrySet struct {
	Status int `json:"status"`
}
