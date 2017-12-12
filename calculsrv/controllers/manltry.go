package controllers

import (
	"calculsrv/models/dbmgr"
	"calculsrv/models/ltrymgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type Manltry struct {
	beego.Controller
}

func (o *Manltry) Post() {
	req := ReqManltry{}
	resp := RespManltry{}

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
		beego.Error(err)
		return
	}

	beego.Info("开始手动开奖 !")

	for _, v := range ltrymgr.Instance().LtryifMap {
		if req.GameName == v.GetGameName() {
			orders := dbmgr.GetLotteryOrderRecord(req.GameName, req.Expect)
			v.SettlementOrders(orders, req.Opencode)
			resp.Status = 0
			body, err := json.Marshal(resp)
			if err != nil {
				beego.Debug(err)
				return
			}
			beego.Info("手动开奖成功 !")
			o.Ctx.Output.Body(body)
			return
		}
	}

	beego.Info("手动开奖失败 !")
	resp.Status = 1
	body, err := json.Marshal(resp)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}

type ReqManltry struct {
	GameName string
	Expect   int
	Opencode string
}

type RespManltry struct {
	Status int // 0成功 1失败没有找到要开奖的菜种
}
