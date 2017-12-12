package controllers

import (
	"github.com/astaxie/beego"
)

type GetOrderInfoByDay struct {
	beego.Controller
}

func (o *GetOrderInfoByDay) Post() {
	//Req := ReqGetOrderByDay{}
	//Resp := GlobalData.RespGetOrderInfo{}
}

type ReqGetOrderByDay struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Type        int    `json:"type"` //0. 普通订单 1.追号订单 2.合买订单
	Time        string `json:"time"` //时间格式
}
