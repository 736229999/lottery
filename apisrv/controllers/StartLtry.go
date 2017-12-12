package controllers

import (
	"apisrv/models/ltrymgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type StartLtry struct {
	beego.Controller
}

func (o *StartLtry) Post() {
	//验证开奖消息是否来自总管理后台
	//测试的时候暂时注释掉IP验证,上服务器的时候记得取消
	// if ctrl.GenMgrSrv.Ip != o.Ctx.Input.IP() {
	// 	beego.Error("The server may be attacked, and the initial manual lottery message IP source is incorrect")
	// 	return
	// }

	req := ReqStartLtry{}
	resp := RespStartLtry{}

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
		beego.Error(err)
		return
	}

	resp.Status = ltrymgr.Instance().StartManLtry(req.LtryName, req.CurrentExpect, req.CurrentOpenCode, req.CurrentOpenTime, req.NextExpect, req.NextOpenTime)

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}

	o.Ctx.Output.Body(body)
}

type ReqStartLtry struct {
	LtryName        string `json:"lotteryName"`
	CurrentExpect   int    `json:"currentExpect"`
	CurrentOpenCode string `json:"currentOpenCode"` //逗号分隔的开奖结果 如: 1,2,3,4,5,6,7
	CurrentOpenTime string `json:"currentOpenTime"` //2006-01-02 15:04:05 格式
	NextExpect      int    `json:"nextExpect"`
	NextOpenTime    string `json:"nextOpenTime"`
}

type RespStartLtry struct {
	Status int `json:"status"`
}
