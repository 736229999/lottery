package controllers

import (
	"apisrv/models/encmgr"
	"apisrv/models/ltrymgr"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
)

type GetLtryNewestRecord struct {
	beego.Controller
}

func (o *GetLtryNewestRecord) Post() {
	//这里就不验证是否是ctrl 服里面的服务器了 应为有消息加密 在这里加入ip输出 就可以知道有没有其他人在访问我的api服务器
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		//输出请求服务器IP
		beego.Error(err, " --- Req Srv Ip : ", o.Ctx.Input.IP())
		return
	}

	req := ReqLtryNewestRecord{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	resp := RespLtryNewestRecord{}
	if ltry, ok := ltrymgr.Instance().Ltrys[req.GameName]; ok {
		resp.GameName = ltry.GetGameName()
		resp.CurrentExpect = ltry.GetCurrentExpect()
		resp.OpenCode = ltry.GetOpenCode()
		resp.OpenCodeStr = ltry.GetOpenCodeStr()
		resp.CurrentOpenTime = ltry.GetCurrentOpenTime()
		resp.NextExpect = ltry.GetNextExpect()
		resp.NextOpenTime = ltry.GetNextOpenTime()
	} else if req.GameName == "HK6" {
		hk6 := ltrymgr.Instance().HK6
		resp.GameName = hk6.GetGameName()
		resp.CurrentExpect = hk6.GetCurrentExpect()
		resp.OpenCode = hk6.GetOpenCode()
		resp.OpenCodeStr = hk6.GetOpenCodeStr()
		resp.CurrentOpenTime = hk6.GetCurrentOpenTime()
		resp.NextExpect = hk6.GetNextExpect()
		resp.NextOpenTime = hk6.GetNextOpenTime()
	} else {
		beego.Error("Not have this lottery : ", req.GameName)
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

type ReqLtryNewestRecord struct {
	GameName string
}

type RespLtryNewestRecord struct {
	GameName        string
	CurrentExpect   int
	OpenCode        []int
	OpenCodeStr     string
	CurrentOpenTime time.Time

	NextExpect   int
	NextOpenTime time.Time
}
