package controllers

import (
	"calculsrv/models/dbmgr"
	"calculsrv/models/ltrymgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type UpdateLtryInfo struct {
	beego.Controller
}

func (o *UpdateLtryInfo) Post() {
	req := ReqUpdateLtryInfo{}
	resp := RespUpdateLtryInfo{}

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
		beego.Error(err)
		return
	}

	limap := dbmgr.UpdateLotteriesInfo()
	if limap == nil {
		resp.Status = 1
		body, err := json.Marshal(resp)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(body)
		return
	}

	for k, v := range ltrymgr.Instance().LtryifMap {
		if li, ok := limap[k]; ok {
			if !v.UpdateLtryInfo(li) {
				beego.Error("Up date lottery error : ", k)
				resp.Status = 2
			}
		} else {
			beego.Warn("Not have this lottery info : ", k)
		}
	}

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}
	o.Ctx.Output.Body(body)
}

type ReqUpdateLtryInfo struct {
	UpdateType string `json:"updateType"` // 值为 : recommend , sort
}

type RespUpdateLtryInfo struct {
	Status int `json:"status"` // 错误码 1:从管理数据库获取彩票信息失败. 2:有彩票更新不成功
}
