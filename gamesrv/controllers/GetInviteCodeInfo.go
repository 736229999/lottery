package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

type GetInviteCodeInfo struct {
	beego.Controller
}

func (o *GetInviteCodeInfo) Post() {
	req := ReqGetInviteCodeInfo{}
	resp := RespStatus{}

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
	bm := bson.M{"agent_id": accountInfo.Agent_Id}
	respInviteCodeDetail := []InviteCodeDetail{}
	err = dbmgr.Instance().InvitationCodeCollection.Find(bm).All(&respInviteCodeDetail)
	if err != nil { //4 查询数据失败
		resp.Status = 4
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//返回客户端
	bufres, _ := json.Marshal(respInviteCodeDetail)
	o.Ctx.Output.Body(bufres)
}

//客户端消息请求结构
type ReqGetInviteCodeInfo struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Flag        string `json:"flag"`
}

//回复客户端状态结构
type RespStatus struct {
	Status int `json:"status"` //状态码
}

//回复客户端详细信息结构
// type RespInviteCodeInfo struct {
// 	InviteCodeInfos []InviteCodeDetail `json:"inviteCodeInfos"`
// }

//邀请码结构
type InviteCodeDetail struct {
	Type       int     `bson:"type"` //邀请码类型 1,为代理商 2为用户
	CreateTime int64   `bson:"create_time"`
	Rebate     float64 `bson:"rebate"`
	Status     int     `bson:"status"` //状态1为可用
	Remark     string  `bson:"remark"` //备注
	Code       string  `bson:"code"`   //邀请码
}
