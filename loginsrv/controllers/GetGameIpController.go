package controllers

import (
	"encoding/json"

	"loginsrv/models/Login"

	"github.com/astaxie/beego"
)

type GetGameIpController struct {
	beego.Controller
}

//在用户既没有登录,也没有注册,客户端没有任何服务器信息的情况下才会请回请求这条消息(也就是说首次使用客户端的时候才会请求这条信息)
func (o *GetGameIpController) Post() {
	ip, err := Login.GetGameIp()
	if err != nil {
		beego.Error(err)
		return
	}

	send := make(map[string]interface{})
	send["ip"] = ip
	bufres, _ := json.Marshal(send)
	o.Ctx.Output.Body(bufres)
}
