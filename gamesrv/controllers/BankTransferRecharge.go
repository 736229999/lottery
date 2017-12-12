package controllers

import (
	"bytes"
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"gamesrv/models/ctrl"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

//银行转账充值
type BankTransferRecharge struct {
	beego.Controller
}

func (o *BankTransferRecharge) Post() {
	//请求结构,和回复结构
	cReq := GlobalData.ReqBankTransfer{}
	cResp := make(map[string]interface{})
	//解析请求结构
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &cReq)
	if err != nil {
		beego.Error("----------- 客户端 银行转账消息 Json 解析错误 ! :", err, " -----------")
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(cReq.AccountName)
	if err != nil { //1 未找到账号
		cResp["status"] = 1
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}
	//验证token
	b := accountInfo.VerifyToken(cReq.Token)
	if !b { //9 token错误
		cResp["status"] = 9
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}
	//验证资金密码
	// b = accountInfo.VerifyMoneyPassword(cReq.MPW)
	// if !b { //9 token错误
	// 	cResp["status"] = 11
	// 	bufres, _ := json.Marshal(cResp)
	// 	o.Ctx.Output.Body(bufres)
	// 	return
	// }

	reqMbodyByte, err := json.Marshal(cReq)
	if err != nil {
		beego.Error(err)
		return
	}

	reqMbody := bytes.NewBuffer(reqMbodyByte)

	to, err := time.ParseDuration("30s")
	if err != nil {
		beego.Error("------------------------- time.ParseDuration() : ", err, " -------------------------")
		return
	}
	//http://cp08com.com/index.php/Inform/
	c := &http.Client{
		Timeout: to}
	respM, err2 := c.Post("http://"+ctrl.MgrSrv.Ip+"/index.php/Inform/saveRechargeBank", "application/json;charset=utf-8", reqMbody)
	if err2 != nil {
		beego.Error(err2)
		cResp["status"] = 10 //超时
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}

	defer respM.Body.Close()

	respMbody, err3 := ioutil.ReadAll(respM.Body)
	if err3 != nil {
		beego.Error(err)
		return
	}

	o.Ctx.Output.Body(respMbody)
}
