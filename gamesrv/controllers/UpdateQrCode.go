package controllers

import (
	"gamesrv/models/QrCodeMgr"

	"github.com/astaxie/beego"
)

type UpdateQrCode struct {
	beego.Controller
}

func (o *UpdateQrCode) Post() {
	QrCodeMgr.Instance().UpdateQrCodeMgr()
	beego.Debug("收到来自后台消息,更新 充值二维码信息")
}
