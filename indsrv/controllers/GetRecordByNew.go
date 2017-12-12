package controllers

import (
	"encoding/json"
	"indsrv/models/encmgr"
	"indsrv/models/ltrymgr"
	"time"

	"github.com/astaxie/beego"
)

type GetRecordByNew struct {
	beego.Controller
}

func (o *GetRecordByNew) Post() {
	//这里就不验证是否是ctrl 服里面的服务器了 应为有消息加密 在这里加入ip输出 就可以知道有没有其他人在访问我的api服务器
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		//输出请求服务器IP
		beego.Error(err, " --- Req Srv Ip : ", o.Ctx.Input.IP())
		return
	}

	req := ReqLtryNewRecord{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	resp := RespLtryNewestRecord{}
	if ltry, ok := ltrymgr.Instance().LtryifMap[req.GameName]; ok {
		resp.GameName = ltry.GetGameName()
		resp.CurrentExpect = ltry.GetCurrentExpect()
		resp.OpenCode = ltry.GetOpenCode()
		resp.OpenCodeStr = ltry.GetOpenCodeStr()
		resp.CurrentOpenTime = ltry.GetCurrentOpenTime()
		resp.NextExpect = ltry.GetNextExpect()
		resp.NextOpenTime = ltry.GetNextOpenTime()
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

type ReqLtryNewRecord struct {
	GameName string
}

type RespLtryNewestRecord struct {
	GameName        string
	CurrentExpect   int
	OpenCode        []int
	OpenCodeStr     string
	CurrentOpenTime time.Time
	NextExpect      int
	NextOpenTime    time.Time
}
