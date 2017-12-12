package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"gamesrv/models/ctrl"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
)

type PrePlayer struct {
	beego.Controller
}
type count struct {
	ID             string
	SEQUENCE_VALUE int
}

func (o *PrePlayer) Post() {
	reqPrepareLogin := &GlobalData.ReqPrepareLogin{}
	respPrepareLogin := &GlobalData.ResqPrepareLogin{} //目前默认消息到这里只要顺利传达就一定会成功注册或登录

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &reqPrepareLogin)
	if err != nil {
		return
	}

	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(reqPrepareLogin.AccountName)

	if err != nil { //没有这个账号的数据
		//获取平台相关数据
		platformConf := GlobalData.PlatformConf{}

		if ctrl.SelfSrv.Type == 1 {
			err = dbmgr.Instance().PlatformConfCollection.Find(nil).One(&platformConf)
			if err != nil {
				beego.Error(err)
				respPrepareLogin.Status = 13 //获取平台数据失败
				bufres, _ := json.Marshal(respPrepareLogin)
				o.Ctx.Output.Body(bufres)
				return
			}
		}

		accountInfo.Account_Type = reqPrepareLogin.AccountType
		accountInfo.Account_Status = 1 //用户状态 1正常 2冻结
		accountInfo.Account_Name = reqPrepareLogin.AccountName
		accountInfo.Flag = reqPrepareLogin.Flag //写完以后来探究这里保存flag的意义
		accountInfo.Token = reqPrepareLogin.Token
		accountInfo.Regist_Time = reqPrepareLogin.RegistTimeStamp
		accountInfo.Regist_Ip = reqPrepareLogin.RegistIp
		accountInfo.Last_Login_Time = reqPrepareLogin.LastLoginTimeStamp
		accountInfo.Last_Login_Ip = reqPrepareLogin.LastLoginIp
		accountInfo.Group = 1 //默认用户组就是1

		var isFormal bool

		if ctrl.SelfSrv.Type == 1 {
			isFormal = true
		} else {
			isFormal = false
		}

		if accountInfo.Account_Type == 0 && isFormal == false {
			accountInfo.Money = 2000  //试玩用户直接给2000
			accountInfo.Rebate = 0.13 //默认返水13%
			accountInfo.Group = 1
		} else if accountInfo.Account_Type == 1 && isFormal == true { //正式账号才用查找这些信息

			accountInfo.Account_Id, err = dbmgr.Instance().GetAccountIncrmentId() //账号ID 自增ID
			if err != nil {
				beego.Error(err)
				respPrepareLogin.Status = 11 //获取自增ID错误
				bufres, _ := json.Marshal(respPrepareLogin)
				o.Ctx.Output.Body(bufres)
				return
			}
			//查找邀请码相关信息
			inviteCodeInfo, err1 := dbmgr.Instance().GetInviteCodeRelatedInfo(reqPrepareLogin.InviteCode)
			//查找开出这个邀请码代理商的信息
			agentInfo := AccountMgr.AccountInfo{}
			err = agentInfo.InitAgent(inviteCodeInfo.AgentId)
			//达成4个条件才能生成代理商账号
			if err1 == nil && err == nil && inviteCodeInfo.Status == 1 && agentInfo.Remain_Agent_Count > 1 && inviteCodeInfo.Type == 1 {
				//如果上级是测试账号,这个账号也就是测试账号
				if agentInfo.Account_Type == 2 {
					accountInfo.Account_Type = 2
				}

				//招到邀请码信息, 招到生成这个邀请码的代理商信息, 邀请码为可用状态, 查看有没有超过代理商生成数量限制, 邀请码类型为代理商类型才能生成代理商
				accountInfo.Rebate = inviteCodeInfo.Rebate

				accountInfo.Belong_Agent_Id = inviteCodeInfo.AgentId //所属代理商ID 这里就是用户自增号

				accountInfo.Belong_Agent = agentInfo.Account_Name

				accountInfo.Is_Agent = 1 //是代理商

				accountInfo.Agent_Id = accountInfo.Account_Id //如果这个邀请码是开启代理商邀请码,那么这个值就是之身自增号

				accountInfo.Agent_Lv = agentInfo.Agent_Lv + 1 //代理商层级+1

				accountInfo.Remaining_Agent_Invite_Code_Count = platformConf.RemainingAgentInviteCodeCount

				accountInfo.Remaining_User_Invite_Code_Count = platformConf.RemainingUserInviteCodeCount

				accountInfo.Remain_Agent_Count = platformConf.RemainAgentCount //允许生成30个下级代理商

				accountInfo.Invite_Code = reqPrepareLogin.InviteCode

				agentInfo.Low_Lv_Agent_Count++ //上级的下级代理商数+1
				agentInfo.Remain_Agent_Count-- //上级的剩余可生成代理商数-1
				//更新上级代理商信息
				if !agentInfo.Update() {
					beego.Error(err)
					respPrepareLogin.Status = 12 //更新代理商信息错误
					bufres, _ := json.Marshal(respPrepareLogin)
					o.Ctx.Output.Body(bufres)
					return
				}
			} else if err1 == nil && err == nil && inviteCodeInfo.Status == 1 && inviteCodeInfo.Type == 2 {
				//如果上级是测试账号,这个账号也就是测试账号
				if agentInfo.Account_Type == 2 {
					accountInfo.Account_Type = 2
				}
				//生成用户
				accountInfo.Rebate = inviteCodeInfo.Rebate

				accountInfo.Belong_Agent_Id = inviteCodeInfo.AgentId //所属代理商ID 这里就是用户自增号

				accountInfo.Belong_Agent = agentInfo.Account_Name

				accountInfo.Is_Agent = 0 //是普通用户

				accountInfo.Agent_Id = 0

				accountInfo.Agent_Lv = 0

				accountInfo.Remaining_Agent_Invite_Code_Count = 0

				accountInfo.Remaining_User_Invite_Code_Count = 0

				accountInfo.Remain_Agent_Count = 0 //允许生成30个下级代理商

				accountInfo.Invite_Code = reqPrepareLogin.InviteCode

				agentInfo.Low_Lv_Account_Count++ //上级的下级代理商数+1

				//更新上级代理商信息
				if !agentInfo.Update() {
					beego.Error(err)
					respPrepareLogin.Status = 12 //更新代理商信息错误
					bufres, _ := json.Marshal(respPrepareLogin)
					o.Ctx.Output.Body(bufres)
					return
				}
			} else { //只要有错就是公司自己的自来用户
				accountInfo.Money = 0                          //呵呵正式账号没有钱
				accountInfo.Rebate = platformConf.Rebate / 100 //默认为我们自己的用户13%反水 日吗不要问我这里为什么要除100 我不想吐槽姐姐
				accountInfo.Is_Agent = 0                       //只是普通用户
				accountInfo.Agent_Id = 0                       //普通用户没有代理商ID
				accountInfo.Belong_Agent = platformConf.Name
			}

			//获取线上线下免费提款次数
			payType, err2 := dbmgr.Instance().GetPayTypeByGroupId(accountInfo.Group)
			if err2 != nil {
				beego.Error(err2)
				return
			}
			bankInfo, err3 := dbmgr.Instance().GetInspectInfoBank(payType)
			if err3 != nil {
				beego.Error(err2)
				return
			}
			onlineInfo, err4 := dbmgr.Instance().GetInspectInfoOnline(payType)
			if err4 != nil {
				beego.Error(err2)
				return
			}

			accountInfo.Bank_Drawing_Count = bankInfo.Drawings.InspectDetail.CommissionCount     //当用户组银行提款免费次数
			accountInfo.Online_Drawing_Count = onlineInfo.Drawings.InspectDetail.CommissionCount //当用户组在线免费提款次数
		}

		if !dbmgr.Instance().InsertAccountInfo(accountInfo) {
			beego.Error("----------------------------------------- Prepare Error ! -----------------------------------------------")
		}
		return
	}

	//有这个玩家的账号就更新信息
	if reqPrepareLogin.AccountName == accountInfo.Account_Name {
		//更新数据
		accountInfo.Token = reqPrepareLogin.Token
		accountInfo.Flag = reqPrepareLogin.Flag
		accountInfo.Last_Login_Time = reqPrepareLogin.LastLoginTimeStamp
		accountInfo.Last_Login_Ip = reqPrepareLogin.LastLoginIp
		if accountInfo.Update() {
			return
		}
	}
}
