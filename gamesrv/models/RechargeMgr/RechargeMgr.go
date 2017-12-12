package RechargeMgr

import (
	"gamesrv/models/dbmgr"
	"sync"

	"github.com/astaxie/beego"
)

var sInstance *RechargeMgr
var once sync.Once

//初始化实例
func Instance() *RechargeMgr {
	once.Do(func() {
		sInstance = &RechargeMgr{}
		sInstance.GetRechargeChannels()
	})
	return sInstance
}

type RechargeMgr struct {
	RechargeChannels RechargeChannel
}

//得到用户可充值渠道
func (o *RechargeMgr) GetRechargeChannels() {
	err := dbmgr.Instance().TransferBankCardCollection.Find(nil).All(&(o.RechargeChannels.TransferBanks))
	if err != nil {
		beego.Error(err)
		return
	}

	err = dbmgr.Instance().OnlinePaymentCollection.Find(nil).All(&(o.RechargeChannels.OnlinePayments))
	if err != nil {
		beego.Error(err)
		return
	}

}

//得到该用户可充值渠道
func (o *RechargeMgr) GetRechargeChannelsByAccount(groupId int, device string) RespRechargeChannels {
	respRechargeChannels := RespRechargeChannels{}

	//第一步根据用户组筛选线下支付渠道
	for _, v := range o.RechargeChannels.TransferBanks {
		if v.Status == 1 { //状态1,表示可用
			for _, i := range v.AccountGroup {
				if groupId == i { //是否符合用户组
					respTransferBank := RespTransferBank{} //如果符合用户组删选要返回给客户端的信息
					respTransferBank.Id = v.Id
					respTransferBank.BankName = v.BankName
					respTransferBank.BankCard = v.BankCard
					respTransferBank.CardHolder = v.CardHolder
					respTransferBank.BankOfDeposit = v.BankOfDeposit

					respRechargeChannels.TransferBanks = append(respRechargeChannels.TransferBanks, respTransferBank)
					break
				}
			}
		}
	}

	//第二部根据用户组,请求平台,筛选支付渠道
	for _, v := range o.RechargeChannels.OnlinePayments {
		if v.Status == 1 { //状态1,表示可用
			for _, i := range v.AccountGroup {
				if groupId == i { //是否符合用户组
					for _, j := range v.Devices {
						if device == j { //是否符合平台
							respOnlinePayment := RespOnlinePayment{}
							respOnlinePayment.Id = v.Id
							respOnlinePayment.CallbackAddress = v.CallbackAddress
							respOnlinePayment.Code = v.Code
							respOnlinePayment.MchId = v.MchId
							respOnlinePayment.PaymentAddress = v.PaymentAddress
							respOnlinePayment.PaymentPlatform = v.PaymentPlatform
							respOnlinePayment.PaymentType = v.PaymentType
							respOnlinePayment.Remark = v.Remark
							respOnlinePayment.OpenBrowser = v.OpenBrowser

							respRechargeChannels.OnlinePayments = append(respRechargeChannels.OnlinePayments, respOnlinePayment)
							break
						}
					}
				}
			}
		}
	}

	return respRechargeChannels
}

//-------返回结构

//所有转账渠道(注意 :这是返回用结构)
type RespRechargeChannels struct {
	TransferBanks  []RespTransferBank  `json:"transferBanks"`
	OnlinePayments []RespOnlinePayment `json:"onlinePayments"`
}

//银行转账渠道(注意 :这是返回用结构)
type RespTransferBank struct {
	Id            int    `bson:"id"`              //id
	BankName      string `bson:"bank_name"`       //银行名称
	BankCard      string `bson:"bank_card"`       //银行卡号卡号
	CardHolder    string `bson:"card_holder"`     //持卡人
	BankOfDeposit string `bson:"bank_of_deposit"` //开户银行
}

//在线充值渠道(注意 :这是返回用结构)
type RespOnlinePayment struct {
	Id              int    `bson:"id"` //id
	PaymentAddress  string `bson:"payment_address"`
	CallbackAddress string `bson:"callback_address"`
	MchId           string `bson:"mchID"`
	PaymentPlatform string `bson:"payment_platform"`
	PaymentType     string `bson:"payment_type"`
	Remark          string `bson:"remark"`
	Code            string `bson:"code"`
	OpenBrowser     int    `bson:"open_browser"`
}

//--------本类使用-------

//所有充值渠道
type RechargeChannel struct {
	TransferBanks  []TransferBank  `json:"transferBanks"`
	OnlinePayments []OnlinePayment `json:"onlinePayments"`
}

//线下转账渠道
type TransferBank struct {
	Status        int    `bson:"status"`          //状态
	Id            int    `bson:"id"`              //id
	BankName      string `bson:"bank_name"`       //银行名称
	BankCard      string `bson:"bank_card"`       //银行卡号卡号
	CardHolder    string `bson:"card_holder"`     //持卡人
	BankOfDeposit string `bson:"bank_of_deposit"` //开户银行
	AccountGroup  []int  `bson:"account_group"`   //用户组
}

//线上支付渠道
type OnlinePayment struct {
	Status          int      `bson:"status"` //状态
	Id              int      `bson:"id"`     //id
	PaymentAddress  string   `bson:"payment_address"`
	CallbackAddress string   `bson:"callback_address"`
	MchId           string   `bson:"mchID"`
	PaymentPlatform string   `bson:"payment_platform"`
	PaymentType     string   `bson:"payment_type"`
	Remark          string   `bson:"remark"`
	Code            string   `bson:"code"`
	Devices         []string `bson:"devices"`
	AccountGroup    []int    `bson:"account_group"` //用户组
	OpenBrowser     int      `bson:"open_browser"`  //是否浏览器打开
}
