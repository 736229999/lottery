package controllers

import (
	"encoding/json"
	"loginsrv/models/AccountMgr"
	"loginsrv/models/ctrl"

	"github.com/astaxie/beego"
)

type ModifyPasswordController struct {
	beego.Controller
}

func (o *ModifyPasswordController) Post() {
	//试玩服不提供修改密码
	// if !common.IsFormalServer() {
	// 	return
	// }

	if ctrl.SelfSrv.Type == 0 {
		return
	}

	req := make(map[string]interface{})
	resp := make(map[string]interface{})
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, "-------------------------")
		return
	}
	accountName := req["accountName"].(string)
	oldPassword := req["oldPassword"].(string)
	newPassword := req["newPassword"].(string)

	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(accountName)
	if err != nil {
		resp["status"] = 2 //初始化账号出错(查询账号出错)
		body, err := json.Marshal(resp)
		if err != nil {
			beego.Emergency(err)
			return
		}
		o.Ctx.Output.Body(body)
		return
	}

	if !accountInfo.ModifyPassword(oldPassword, newPassword) {
		resp["status"] = 1 //1失败，
		body, err := json.Marshal(resp)
		if err != nil {
			beego.Emergency(err)
			return
		}
		o.Ctx.Output.Body(body)
		return
	}

	resp["status"] = 0 //成功

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Emergency(err)
		return
	}
	o.Ctx.Output.Body(body)
}
