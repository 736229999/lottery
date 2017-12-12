package controllers

import (
	"bytes"
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/ctrl"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

//银行转账充值
type UserRequestBankWithdrawals struct {
	beego.Controller
}

func (o *UserRequestBankWithdrawals) Post() {
	//请求结构,和回复结构
	cReq := ReqBankWithdrawal{}
	cResq := make(map[string]interface{})
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
		cResq["status"] = 1
		bufres, _ := json.Marshal(cResq)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(cReq.Token)
	if !b { //9 token错误
		cResq["status"] = 9
		bufres, _ := json.Marshal(cResq)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证资金密码
	b = accountInfo.VerifyMoneyPassword(cReq.MPW)
	if !b {
		cResq["status"] = 11
		bufres, _ := json.Marshal(cResq)
		o.Ctx.Output.Body(bufres)
		return
	}

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

	c := &http.Client{
		Timeout: to}
	respM, err2 := c.Post("http://"+ctrl.MgrSrv.Ip+"/index.php/Inform/saveDrawingsBank", "application/json;charset=utf-8", reqMbody)
	if err2 != nil {
		beego.Error(err2)
		cResq["status"] = 10 //超时
		bufres, _ := json.Marshal(cResq)
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

//客户端请求提款消息
type ReqBankWithdrawal struct {
	AccountName     string `json:"accountName"`     //账号名,提款账号
	Token           string `json:"token"`           //token
	Flag            string `json:"flag"`            //flag
	MPW             string `json:"mpw"`             //资金密码
	WithdrawalMongy string `json:"withdrawalMoney"` //提款金额
}
