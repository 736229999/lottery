package controllers

import (
	"apisrv/models/apimgr"
	"apisrv/models/encmgr"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
)

type GetLtryHistByDay struct {
	beego.Controller
}

func (o *GetLtryHistByDay) Post() {
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		beego.Debug(err)
		return
	}

	req := ReqLtryHistByDay{}
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

type ReqLtryHistByDay struct {
	GameName string
	Date     time.Time
}

type RespLtryHistByDay struct {
	Rows   int                 `json:"rows"`
	Code   string              `json:"code"`
	Remain string              `json:"remain"`
	Data   []LtryRecordFromApi `json:"data"`
}

//开采网 api 一条记录
type LtryRecordFromApi struct {
	Expect        string `json:"expect"`
	Opencode      string `json:"opencode"`
	OpenTime      string `json:"opentime"`
	OpenTimeStamp int64  `json:"opentimestamp"`
}
