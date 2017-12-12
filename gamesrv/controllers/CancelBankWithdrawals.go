package controllers

import (
	"github.com/astaxie/beego"
)

type CancelBankWithdrawals struct {
	beego.Controller
}

//做到一半发现不行,操作人员有可能误操作
func (o *CancelBankWithdrawals) Post() {
	// req := ReqCancelBankWithdrawals{}
	// resp := RespCancelBankWithdrawals{}

	// err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	// if err != nil {
	// 	beego.Debug(err)
	// 	return
	// }

	// //得到账户
	// accountInfo := AccountMgr.AccountInfo{}
	// err = accountInfo.Init(req.AccountName)
	// if err != nil { //1 未找到账号
	// 	resp.Status = 1
	// 	bufres, _ := json.Marshal(resp)
	// 	o.Ctx.Output.Body(bufres)
	// 	return
	// }

	// //验证token
	// b := accountInfo.VerifyToken(req.Token)
	// if !b { //2 token错误
	// 	resp.Status = 2
	// 	bufres, _ := json.Marshal(resp)
	// 	o.Ctx.Output.Body(bufres)
	// 	return
	// }

	//取消提款 只有状态为审核中的可以取消

}

//用户取消提款请求结构结构
type ReqCancelBankWithdrawals struct {
	AccountName string `json:"accountName"`
	Token       string `json:"Token"`
	Flag        string `json:"flag"`
}

//用户取消提款返回结构
type RespCancelBankWithdrawals struct {
	Status int `json:"status"`
}
