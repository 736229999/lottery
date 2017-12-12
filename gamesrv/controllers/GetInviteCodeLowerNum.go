package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type GetInviteCodeLowerNum struct {
	beego.Controller
}

func (o *GetInviteCodeLowerNum) Post() {
	req := ReqGetInviteCodeLowerNum{}
	resp := RespGetInviteCodeLowerNum{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug(err)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil { //1 未找到账号
		resp.Status = 1
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(req.Token)
	if !b { //2 token错误
		resp.Status = 2
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}
	//是否是代理商,不是代理商不能查这些信息
	if accountInfo.Is_Agent != 1 {
		resp.Status = 3 //该账号不是代理商
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//在管理数据库查找所有这个账号开出的邀请码
	bm := bson.M{"invite_code": req.InviteCode}
	resp.Count, err = dbmgr.Instance().AccountInfoCollection.Find(bm).Count()
	if err != nil { //4 查询数据失败
		resp.Status = 4
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//返回客户端
	bufres, _ := json.Marshal(resp)
	o.Ctx.Output.Body(bufres)
}

type ReqGetInviteCodeLowerNum struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Flag        string `json:"flag"`
	InviteCode  string `json:"inviteCode"`
}

type RespGetInviteCodeLowerNum struct {
	Status int `json:"status"`
	Count  int `json:"count"`
}
