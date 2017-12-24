package controllers

import (
	"encoding/json"
	//"strings"
	"time"

	"loginsrv/models/DbHandle"
	"loginsrv/models/Login"
	"loginsrv/models/Utils"
	"loginsrv/models/ctrl"

	"github.com/astaxie/beego"
	"github.com/mojocn/base64Captcha"
)

type RegistController struct {
	beego.Controller
}

func (o *RegistController) Post() {
	beego.Debug("注册验证码------------------------注册验证码")
	reqRegist := &Login.ReqRegist{}
	respLogin := &Login.RespLogin{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, reqRegist)
	if err != nil {
		beego.Error("--- Marshal Error : ", err)
		return
	}
	if ctrl.SelfSrv.Type == 1 {
		if reqRegist.AccountType != 1 {
			return
		}
		//验证密码合法 注意:现在只有正式服务器注册接收密码
		if !Login.VerifyPassword(reqRegist.Password) {
			return
		}
	} else {
		if reqRegist.AccountType != 0 {
			return
		}
	}

	//平台判断 //0苹果，1安卓，2，Wap，3 PC网页端
	if reqRegist.RegistrationPlatform < 0 || reqRegist.RegistrationPlatform > 3 {
		return
	}

	//验证账户名合法
	if !Login.VerifyAccountName(reqRegist.AccountName) {
		return
	}

	//验证码,标识不能为空
	if reqRegist.Captcha == "" || reqRegist.Flag == "" {
		return
	}

	//是否找到验证码
	cpatcha, ok := Login.VerifyInfo[reqRegist.Flag]

	if !ok {
		respLogin.State = 1
		//beego.Debug("验证码超时或无效")
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error("--- Marshal Error : ", err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}
	//一旦找到了验证码就删除掉
	delete(Login.VerifyInfo, reqRegist.Flag)

	//验证码是否输入错误
	verifyResult := base64Captcha.VerifyCaptcha(cpatcha, reqRegist.Captcha) //传入之前生成后存的验证码和前台输入的验证码，这个方法返回bool
	if verifyResult {
		beego.Debug("------------------------------验证码正确-------------------------")
	} else {
		beego.Debug("------------------------------验证码错误-------------------------")
		respLogin.State = 2
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error("--- Marshal Error : ", err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}
	beego.Debug("验证码正确继续运行")
	//if strings.ToLower(reqRegist.Captcha) != strings.ToLower(cpatcha) {
	//	respLogin.State = 2
	//	msg, err := json.Marshal(respLogin)
	//	if err != nil {
	//		beego.Error("--- Marshal Error : ", err)
	//		return
	//	}
	//	o.Ctx.Output.Body(msg)
	//	return
	//}
	//这里要去数据库中查询看是否被注册过了
	accountInfo := DbHandle.FindAccountInfo(reqRegist.AccountName)
	//用户已经注册
	if accountInfo.Password != "" {
		respLogin.State = 3
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}
	//用户可以注册插入数据库
	t := time.Now()
	accountInfo.Account_Name = reqRegist.AccountName
	accountInfo.Password = reqRegist.Password
	accountInfo.Regist_Time = t
	accountInfo.Regist_Time_Stamp = t.Unix()
	if reqRegist.Ip == "" { //如果注册IP为空那么就认为是app的注册,获取消息来源的IP
		accountInfo.Regist_Ip = o.Ctx.Input.IP()
	} else { //如果注册IP有值,那么就认为是来自web或wap的注册
		accountInfo.Regist_Ip = reqRegist.Ip
	}

	accountInfo.Last_Login_Time = t
	accountInfo.Last_Login_Time_Stamp = t.Unix()
	accountInfo.Last_Login_Ip = accountInfo.Regist_Ip
	accountInfo.Registration_Platform = reqRegist.RegistrationPlatform
	accountInfo.Account_Type = reqRegist.AccountType

	//入库
	if !DbHandle.InsertAccount(accountInfo) {
		beego.Error("--- 严重错误,注册入库错误")
		respLogin.State = 4
		msg, err := json.Marshal(respLogin)
		if err != nil {
			beego.Error(err)
			return
		}
		o.Ctx.Output.Body(msg)
		return
	}

	//-------------------------------现在注册成功以后直接登录 给客户端就算是登陆成功了-------------------------

	//根据账号类型查账GameServerIp,如果是试玩用户,直接发 0 号 试玩服务器IP 正式用户就发送从1号开始的服务器,只有当1号服务器出问题 没有响应的时候才会发送2号服务器 依次类推
	gameIp, err := Login.GetGameIp()
	if err != nil {
		beego.Emergency(err)
		respLogin.State = 4
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
	respLogin.AccountType = reqRegist.AccountType

	//发送给Game服务器消息
	reqPrepareLogin := &Login.ReqPrepareLogin{}
	reqPrepareLogin.AccountName = reqRegist.AccountName
	reqPrepareLogin.Token = respLogin.Token
	reqPrepareLogin.Flag = reqRegist.Flag
	reqPrepareLogin.AccountType = reqRegist.AccountType
	reqPrepareLogin.InviteCode = reqRegist.InviteCode
	reqPrepareLogin.RegistTimeStamp = accountInfo.Regist_Time_Stamp
	reqPrepareLogin.RegistIp = accountInfo.Regist_Ip
	reqPrepareLogin.LastLoginTimeStamp = accountInfo.Last_Login_Time_Stamp
	reqPrepareLogin.LastLoginIp = accountInfo.Last_Login_Ip

	//现在假设通信成功那么Game服务器就肯定会是正常的状态,以后来改,这里面要加入Game的失败可能性判断
	//这里要返回 game服务器的错误给客户端
	if !Login.PrepareLogin(reqPrepareLogin, respLogin.GameIp) {
		beego.Error("--- 严重错误，预登录失败")
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
