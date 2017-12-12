package controllers

import (
	"apisrv/models/ltrymgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

type InitLtry struct {
	beego.Controller
}

func (o *InitLtry) Post() {
	//验证开奖消息是否来自总管理后台
	//测试的时候暂时注释掉IP验证,上服务器的时候记得取消
	//现在六合彩总后台没有写好暂时取消验证
	// if ctrl.GenMgrSrv.Ip != o.Ctx.Input.IP() {
	// 	beego.Error("The server may be attacked, and the initial manual lottery message IP source is incorrect")
	// 	return
	// }

	//初始化(启动)手动开奖的彩票(注意这里暂时没有加密,等完善稳定后来加密)
	req := ReqInitLtry{}
	resp := RespInitLtry{}

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
		beego.Error(err)
		return
	}

	resp.Status = ltrymgr.Instance().InitManLtry(req.LtryName, req.CurrentExpect, req.NextExpect, req.NextOpenTime)

	body, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}

	o.Ctx.Output.Body(body)
}

type ReqInitLtry struct {
	LtryName      string `json:"lotteryName"`
	CurrentExpect int    `json:"currentExpect"`
	NextExpect    int    `json:"nextExpect"`
	NextOpenTime  string `json:"nextOpenTime"`
}

type RespInitLtry struct {
	Status int `json:"status"` //0. 成功 1.
}
