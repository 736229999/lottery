package controllers

import (
	"gamesrv/models/AnnouncementMgr"
	"gamesrv/models/QrCodeMgr"
	"gamesrv/models/RechargeMgr"

	"github.com/astaxie/beego"
)

type ChangeRechargeChannels struct {
	beego.Controller
}

func (o *ChangeRechargeChannels) Post() {
	beego.Debug("---------------------- 后台 改变充值渠道 , 公告信息--------------------------")
	RechargeMgr.Instance().GetRechargeChannels()
	QrCodeMgr.Instance().UpdateQrCodeMgr()
	AnnouncementMgr.Instance().UpdateAnnouncement()
}
