package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

type ChangeInviteCodeInfo struct {
	beego.Controller
}

func (o *ChangeInviteCodeInfo) Post() {
	req := ReqChangeInviteCodeInfo{}
	resp := RespChangeInviteCodeInfo{}

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

	//验证是否是代理商
	if accountInfo.Is_Agent == 0 {
		resp.Status = 3 //不是代理商(只有代理商才能修改邀请码信息)
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//找到要修改的邀请码的信息 (注意查询条件,agentId 和 code 都要)
	inviteCodeInfo := InviteCodeInfo{}

	bm := bson.M{"agent_id": accountInfo.Agent_Id, "code": req.InviteCode}
	err = dbmgr.Instance().InvitationCodeCollection.Find(bm).One(&inviteCodeInfo)
	if err != nil {
		resp.Status = 4 //4 没有找到这个邀请码
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	// if req.ChangeStatus == 0 { //目前只有停用邀请码这一个功能
	// 	inviteCodeInfo.Status = 0
	// }

	if len(req.Remark) > 80 { //当remark 字段不为空,并且小于80个字符时,认为是要修改备注信息
		resp.Status = 4 //不是代理商(只有代理商才能修改邀请码信息)
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	selector := bson.M{"code": req.InviteCode}
	data := bson.M{"$set": bson.M{"remark": req.Remark, "status": req.ChangeStatus}}
	dbmgr.Instance().InvitationCodeCollection.Update(selector, data)

	resp.Status = 0 //0 成功
	bufres, _ := json.Marshal(resp)
	o.Ctx.Output.Body(bufres)
}

type ReqChangeInviteCodeInfo struct {
	AccountName  string `json:"accountName"`
	Token        string `json:"token"`
	Flag         string `json:"flag"`
	InviteCode   string `json:"inviteCode"`   //邀请码
	ChangeStatus int    `json:"changeStatus"` //改变为什么状态  目前只能改为0停用
	Remark       string `json:"remark"`       //备注  80个字符
}

type RespChangeInviteCodeInfo struct {
	Status int `json:"status"`
}

type InviteCodeInfo struct {
	//Type       int     `bson:"type"` //邀请码类型 1,为代理商 2为用户
	//CreateTime int64   `bson:"create_time"`
	//Rebate     float64 `bson:"rebate"`
	Status int    `bson:"status"` //状态1为可用
	Remark string `bson:"remark"` //备注
	//Code   string `bson:"code"`   //邀请码
}
