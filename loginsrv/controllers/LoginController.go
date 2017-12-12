package controllers

import (
	"encoding/json"
	"loginsrv/models/Login"
	"loginsrv/models/Utils"
	"loginsrv/models/ctrl"

	"loginsrv/models/DbHandle"

	"time"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (o *LoginController) Post() {
	//试玩服务器不提供登录
	if ctrl.SelfSrv.Type == 0 {
		return
	}

	reqLogin := &Login.ReqLogin{}
	respLogin := &Login.RespLogin{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, reqLogin)
	if err != nil {
		beego.Error(err)
		return
	}

	//验证标识
	if reqLogin.Flag == "" {
		return
	}
	//验证账户名
	if !Login.VerifyAccountName(reqLogin.AccountName) {
		return
	}
	//验证密码
	if !Login.VerifyPassword(reqLogin.Password) {
		return
	}

	//判断在数据库中是否存在该用户信息(密码错误或账户不存在都给一个提示)
	accountInfo := DbHandle.FindAccountInfo(reqLogin.AccountName)
	//密码错误或用户不存在
	if reqLogin.Password != accountInfo.Password {
		respLogin.State = 1
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}
	//更新用户登录信息
	t := time.Now()
	var lastLoginIp string
	if reqLogin.Ip == "" {
		lastLoginIp = o.Ctx.Input.IP()
	} else {
		lastLoginIp = reqLogin.Ip
	}
	if !DbHandle.UpdateAccountInfo(accountInfo.Account_Name, t, t.Unix(), lastLoginIp) {
		beego.Error("--- 更新玩家信息错误 !")
		return
	}

	//根据账号类型查账GameServerIp,如果是试玩用户,直接发 0 号 试玩服务器IP 正式用户就发送从1号开始的服务器,只有当1号服务器出问题 没有响应的时候才会发送2号服务器 依次类推
	gameIp, err := Login.GetGameIp()
	if err != nil {
		beego.Error(err)
		respLogin.State = 5
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}
	//得到token
	respLogin.Token = Utils.GetToken()
	respLogin.State = 0
	respLogin.GameIp = gameIp
	respLogin.AccountType = accountInfo.Account_Type
	//发送给Game服务器消息
	reqPrepareLogin := &Login.ReqPrepareLogin{}
	reqPrepareLogin.AccountName = accountInfo.Account_Name
	reqPrepareLogin.Token = respLogin.Token
	reqPrepareLogin.Flag = reqLogin.Flag
	reqPrepareLogin.AccountType = accountInfo.Account_Type
	reqPrepareLogin.AccountId = 0
	reqPrepareLogin.LastLoginTimeStamp = time.Now().Unix()
	reqPrepareLogin.LastLoginIp = lastLoginIp

	if !Login.PrepareLogin(reqPrepareLogin, respLogin.GameIp) {
		beego.Error("--- 严重错误，预登录失败 !")
		respLogin.State = 5
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}

	msg, err := json.Marshal(respLogin)
	if err != nil {
		beego.Error(err)
		return
	}
	o.Ctx.Output.Body(msg)
}
