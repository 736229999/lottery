package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
)

type GetLowerInfo struct {
	beego.Controller
}

func (o *GetLowerInfo) Post() {
	req := ReqLowerInfo{}
	resp := RespLowerInfoStatus{}

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

	//验证是否代理商
	if accountInfo.Is_Agent != 1 {
		resp.Status = 3 //用户不是代理商
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//如果查询对象不是自己
	if req.SearchAccountName != req.AccountName {
		//找到查询对象
		searchAccount := AccountMgr.AccountInfo{}
		err = searchAccount.Init(req.SearchAccountName)
		if err != nil {
			resp.Status = 4 //没有找到要查询的下级账号
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}

		//判断要查找的账号是否是当前账号的下属账号
		if searchAccount.Belong_Agent_Id != accountInfo.Account_Id {
			resp.Status = 5 //查找的这个账号不是当前账号的下属
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}
	}

	//返回要查找的信息(按查找类型)
	switch req.SearchItem {
	case 1: //查询投注明细
		respOrders := GlobalData.RespGetOrderInfo{}
		err = dbmgr.Instance().GetOrderByAccountName(req.SearchAccountName, req.Skip, 30, req.SearchType, &(respOrders.Orders))
		if err != nil {
			resp.Status = 6 //查询失败
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}
		//返回查询结果
		bufres, _ := json.Marshal(respOrders)
		o.Ctx.Output.Body(bufres)
		return

	case 2: //查询充值明细
		//查询充值记录
		respRecharge := []GlobalData.RechargeRecord{}
		//默认查询30条
		err = dbmgr.Instance().GetRechargeRecord(req.SearchAccountName, req.Skip, 30, req.SearchType, &respRecharge)
		if err != nil {
			resp.Status = 6 //查询失败
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}
		//返回查询结果
		bufres, _ := json.Marshal(respRecharge)
		o.Ctx.Output.Body(bufres)
		return

	case 3: //查询提款明细
		//查询充值记录
		respDrawings := []GlobalData.DrawingsRecord{}
		//默认查询30条
		err = dbmgr.Instance().GetDrawingsRecord(req.SearchAccountName, req.Skip, 30, req.SearchType, &respDrawings)
		if err != nil {
			resp.Status = 6 //查询失败
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}
		//计算其他
		// for _, v := range respDrawings {
		// 	v.ActualAmount = v.Money - v.CommissionMoney
		// 	v.AccountBalance = v.MoneyBefore - v.Money
		// }

		//返回查询结果
		bufres, _ := json.Marshal(respDrawings)
		o.Ctx.Output.Body(bufres)
		return

	default:
		resp.Status = 7 //查询类型错误
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}
}

//客户端请求结构
//说明:1.如果是查投注明细
//这里要和GetOrderInfo一样给出查询下注类型 默认还是给30条
//1,全部-全部 2,全部-中奖 3,全部-待开奖 4,普通-全部 5,普通-中奖 6,普通待开奖, 7追号-全部, 8追号-中奖, 9追号-待开奖
//说明:2.如果是查充值记录
//查询类型 1.全部 2.成功 3.等待
//说明:3.如果是查提款记录
//查询类型 1.全部 2.成功 3.等待
type ReqLowerInfo struct {
	AccountName       string `json:"accountName"`
	Token             string `json:"token"`
	Flag              string `json:"flag"`
	SearchAccountName string `json:"searchAccountName"`
	//SearchAccountID   int    `json:"searchNameID"` //要查询的账号的ID (注意是ID)
	SearchItem int `json:"searchItem"` //查询项目 1.投注明细, 2.充值明细 3.提款明细
	Skip       int `json:"skip"`       //跳过条目数
	SearchType int `json:"searchType"` //根据查询类型不一样来变化传入值

}

//错误验证码返回
type RespLowerInfoStatus struct {
	Status int `json:"status"`
}
