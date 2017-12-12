package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
)

type GetInspectMoney struct {
	beego.Controller
}

func (o *GetInspectMoney) Post() {
	cReq := ReqInspectMoney{}
	cResp := RespInspectMoney{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &cReq)
	if err != nil {
		beego.Debug(err)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(cReq.AccountName)
	if err != nil { //1 未找到账号
		cResp.Status = 1
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(cReq.Token)
	if !b { //9 token错误
		cResp.Status = 9
		bufres, _ := json.Marshal(cResp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//组稽核消息所需返回数据
	//这几个数据是用户信息身上的
	cResp.RechargeAmount = accountInfo.Recharge_Amount
	cResp.FavorableAmount = accountInfo.Favorable_Amount

	if cReq.Type == 0 { //线下银行提款稽核金额
		cResp.InspectMoney = accountInfo.Inspect_Money
		cResp.CommissionCount = accountInfo.Bank_Drawing_Count
	} else if cReq.Type == 1 { //线上提款稽核金额
		cResp.InspectMoney = accountInfo.Inspect_Money
		cResp.CommissionCount = accountInfo.Online_Drawing_Count
	}
	cResp.BetAmountImmediate = accountInfo.Bet_Amount_Immediate

	//首先根据用户组,去查询Pay_typ(支付类型)
	var payType int
	payType, err = dbmgr.Instance().GetPayTypeByGroupId(accountInfo.Group)
	if err != nil {
		beego.Error(err)
		return
	}

	//再根据pay_type id,和线上或线下类型 查找剩下的稽核数据
	if cReq.Type == 0 { //线下银行
		var ret GlobalData.InspectMoneyBank
		ret, err = dbmgr.Instance().GetInspectInfoBank(payType)
		if err != nil {
			beego.Error(err)
			return
		}

		cResp.CommissionPercent = ret.Drawings.InspectDetail.CommissionPercent
		cResp.CommissionStatus = ret.Drawings.InspectDetail.CommissionStatus
		cResp.Minimum = ret.Drawings.InspectDetail.Minimum
		cResp.Maximum = ret.Drawings.InspectDetail.Maximum
		cResp.CommissionMinimum = ret.Drawings.InspectDetail.CommissionMinimum
		cResp.CommissionMaximum = ret.Drawings.InspectDetail.CommissionMaximum
		cResp.NormalityInspectBroaden = ret.Drawings.InspectDetail.NormalityInspectBroaden
		cResp.NormalityInspectRate = ret.Drawings.InspectDetail.NormalityInspectRate
	} else if cReq.Type == 1 { //线上
		var ret GlobalData.InspectMoneyOnline
		ret, err = dbmgr.Instance().GetInspectInfoOnline(payType)
		if err != nil {
			beego.Error(err)
			return
		}

		cResp.CommissionPercent = ret.Drawings.InspectDetail.CommissionPercent
		cResp.CommissionStatus = ret.Drawings.InspectDetail.CommissionStatus
		cResp.Minimum = ret.Drawings.InspectDetail.Minimum
		cResp.Maximum = ret.Drawings.InspectDetail.Maximum
		cResp.CommissionMinimum = ret.Drawings.InspectDetail.CommissionMinimum
		cResp.CommissionMaximum = ret.Drawings.InspectDetail.CommissionMaximum
		cResp.NormalityInspectBroaden = ret.Drawings.InspectDetail.NormalityInspectBroaden
		cResp.NormalityInspectRate = ret.Drawings.InspectDetail.NormalityInspectRate
	} else {
		return
	}

	body, err := json.Marshal(cResp)
	if err != nil {
		beego.Error(err)
		return
	}
	o.Ctx.Output.Body(body)
}

//客户端请求稽核消息
type ReqInspectMoney struct {
	AccountName string `json:"accountName"` //账号名,提款账号
	Token       string `json:"token"`       //token
	Flag        string `json:"flag"`        //flag
	Type        int    `json:"type"`        //稽核类型,是在线稽核,还是线下银行稽核 0.线下银行 1.在线
}

//回复客户端稽核消息
type RespInspectMoney struct {
	Status          int `json:"status"`          // 0.状态
	RechargeAmount  int `json:"rechargeAmount"`  //1.总充值金额
	FavorableAmount int `json:"favorableAmount"` //2.总优惠金额
	InspectMoney    int `json:"inspectMoney"`    //3.核查金额
	//TotalBetAmount          int     `json:"totalBetAmount"`          //4.当前投注总额
	BetAmountImmediate      int     `json:"betAmountImmediate"`      //4.当前及时投注
	CommissionPercent       float64 `json:"commissionPercent"`       //5.当前用户组手续费费率
	CommissionCount         int     `json:"commissionCount"`         //6.当前用户组费免费次数
	CommissionStatus        int     `json:"commissionStatus"`        //7.当前用户组是否开启免费
	NormalityInspectRate    int     `json:"normalityInspectRate"`    //8.行政费率
	Minimum                 int     `json:"minimum"`                 //9.当前用户组提款最低值
	Maximum                 int     `json:"maximum"`                 //10.当前用户组提款最高值
	CommissionMinimum       int     `json:"commissionMinimum"`       //11.手续费最低值
	CommissionMaximum       int     `json:"commissionMaximum"`       //12.手续费最高值
	NormalityInspectBroaden int     `json:"normalityInspectBroaden"` //13.常态性核查放宽额度
}
