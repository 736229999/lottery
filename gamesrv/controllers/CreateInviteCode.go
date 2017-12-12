package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type CreateInviteCode struct {
	beego.Controller
}

func (o *CreateInviteCode) Post() {
	req := ReqCreateInviteCode{}
	resp := RespCreateInviteCode{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug(err)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}

	err = accountInfo.Init(req.AccountName)
	if err != nil { //1 未找到账号
		beego.Debug(err)
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

	//是代理商可以创建邀请码
	if accountInfo.Is_Agent == 1 {
		if req.InviteCodeType == 1 { //创建代理商邀请码
			if accountInfo.Remaining_Agent_Invite_Code_Count < 1 { //是否还有剩余代理商邀请码生成次数
				resp.Status = 5 //5 剩余代理商邀请码生成次数已经用完
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			//代理商自身可生成代理商邀请码数量-1
			accountInfo.Remaining_Agent_Invite_Code_Count--
			if !accountInfo.Update() {
				resp.Status = 10 //10 更新代理商数据失败
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}

			//开始生成验证码步骤
			inviteCode := InviteCode{}
			//1,去管理数据库找到自增值 然后用 7个9 -自增值 代理商前面+1
			inc, err1 := dbmgr.Instance().GetAgentInviteCodeIncrmentId()
			if err1 != nil {
				beego.Error(err1)
				resp.Status = 6 //6 自增ID获取失败
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			codeInt := 9999999 - inc
			codeStr := "1" + strconv.Itoa(codeInt)
			inviteCode.Code = codeStr
			inviteCode.Type = req.InviteCodeType
			inviteCode.Agent_Id = accountInfo.Agent_Id
			inviteCode.Create_Time = time.Now().Unix()
			if len(req.Remark) > 80 {
				resp.Status = 7 //7备注超过字符限制
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			inviteCode.Remark = req.Remark
			inviteCode.Status = 1
			if req.Rebate > accountInfo.Rebate {
				resp.Status = 8 //8 反水设置超过自身,失败,失败中的失败...
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}

			inviteCode.Rebate = req.Rebate //注意这里给的是百分比而不是小数
			//存邀请码到管理数据库
			err2 := dbmgr.Instance().InvitationCodeCollection.Insert(inviteCode)
			if err2 != nil {
				beego.Error(err2)
				resp.Status = 9 //9 插入新邀请码数据失败
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}

			//返回成功
			resp.Status = 0 //0 成功
			resp.Code = inviteCode.Code
			resp.CreateTime = inviteCode.Create_Time
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		} else if req.InviteCodeType == 2 { //创建普通用户邀请码
			if accountInfo.Remaining_User_Invite_Code_Count < 1 { //是否还有剩余代理商邀请码生成次数
				resp.Status = 5 //5 剩余代理商邀请码生成次数已经用完
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			//代理商自身可生成代理商邀请码数量-1
			accountInfo.Remaining_User_Invite_Code_Count--
			if !accountInfo.Update() {
				resp.Status = 10 //10 更新代理商数据失败
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			//开始生成验证码步骤
			inviteCode := InviteCode{}
			//1,去管理数据库找到自增值 然后用 7个9 -自增值 代理商前面+1
			inc, err1 := dbmgr.Instance().GetUserInviteCodeIncrmentId()
			if err1 != nil {
				beego.Error(err1)
				resp.Status = 6 //6 自增ID获取失败
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			codeInt := 9999999 - inc
			codeStr := "0" + strconv.Itoa(codeInt) //普通用户是0开头
			inviteCode.Code = codeStr
			inviteCode.Type = req.InviteCodeType
			inviteCode.Agent_Id = accountInfo.Agent_Id
			inviteCode.Create_Time = time.Now().Unix()
			if len(req.Remark) > 80 {
				resp.Status = 7 //7备注超过字符限制
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}
			inviteCode.Remark = req.Remark
			inviteCode.Status = 1
			if req.Rebate > accountInfo.Rebate {
				resp.Status = 8 //8 反水设置超过自身,失败,失败中的失败...
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}

			inviteCode.Rebate = req.Rebate //注意这里给的是百分比而不是小数
			//存邀请码到管理数据库
			err2 := dbmgr.Instance().InvitationCodeCollection.Insert(inviteCode)
			if err2 != nil {
				beego.Error(err2)
				resp.Status = 9 //9 插入新邀请码数据失败
				bufres, _ := json.Marshal(resp)
				o.Ctx.Output.Body(bufres)
				return
			}

			//返回成功
			resp.Status = 0 //0 成功
			resp.Code = inviteCode.Code
			resp.CreateTime = inviteCode.Create_Time
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return

		} else {
			resp.Status = 4 //4创建邀请码类型错误
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}

	} else {
		resp.Status = 3 //3不是代理商不能创建邀请码
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}
}

type ReqCreateInviteCode struct {
	AccountName    string  `json:"accountName"`
	Token          string  `json:"token"`
	Flag           string  `json:"flag"`
	InviteCodeType int     `json:"inviteCodeType"` //创建邀请码类型 1 代理商 2,普通用户
	Rebate         float64 `json:"rebate"`         //反税率 注意:这里是给的是浮点数百分比  0.13
	Remark         string  `json:"remark"`         //备注 40个中文字符
}

type RespCreateInviteCode struct {
	Status     int    `json:"status"`
	Code       string `json:"code"`
	CreateTime int64  `json:"createTime"`
}

//管理数据库邀请码结构
type InviteCode struct {
	Type        int     `bson:"type"` //邀请码类型 1,为代理商 2为用户
	Agent_Id    int     `bson:"agent_id"`
	Create_Time int64   `bson:"create_time"`
	Rebate      float64 `bson:"rebate"`
	Status      int     `bson:"status"` //状态1为可用
	Remark      string  `bson:"remark"` //备注
	Code        string  `bson:"code"`   //邀请码
}
