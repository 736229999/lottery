package controllers

import (
	"ctrlsrv/models/encmgr"
	"ctrlsrv/models/srvmgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type GetDependSrv struct {
	beego.Controller
}

func (o *GetDependSrv) Post() {
	//每条来自其他服务器的消息都要验证来源IP是否正确
	if !srvmgr.Instance().VerifySrv(o.Ctx.Input.IP()) {
		beego.Error("Wrongful server access ! IP : ", o.Ctx.Input.IP())
		return
	}

	//根据不同功能的服务器 返回其依赖服务器信息
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		beego.Debug(err)
		return
	}

	req := ReqGetDependSrv{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	resp := []RespGetDependSrv{}
	for _, v := range srvmgr.Instance().SrvInfos {
		if v.Func == req.Func && o.Ctx.Input.IP() == v.Ip { //同一ip上 不会出现重复功能的服务器 所以以这个为条件来判断
			//请求这条信息的服务器信息
			selfInfo := RespGetDependSrv{}
			selfInfo.Name = v.Name
			selfInfo.Ip = v.Ip
			selfInfo.Port = v.Port
			selfInfo.Func = v.Func
			selfInfo.Type = v.Type

			resp = append(resp, selfInfo)
			//请求的服务器的依赖信息
			for _, i := range v.Depend {
				if depend, ok := srvmgr.Instance().SrvInfos[i]; ok {
					if depend.Status == 1 {
						srvInfo := RespGetDependSrv{}
						srvInfo.Name = depend.Name
						srvInfo.Ip = depend.Ip
						srvInfo.Port = depend.Port
						srvInfo.Func = depend.Func
						srvInfo.Type = depend.Type
						resp = append(resp, srvInfo)
					} else {
						beego.Error("Dependency server status is incorrect ! : ", depend)
					}
				} else {
					beego.Error("Failed to find dependency server ! Request server : ", v.Name, " --- Depend server : ", i)
				}
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

			break
		}
	}

}

type ReqGetDependSrv struct {
	Func string
}

type RespGetDependSrv struct {
	Name string
	Ip   string
	Port string
	Func string
	Type int
}
