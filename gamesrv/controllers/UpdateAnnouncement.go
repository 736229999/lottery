package controllers

import (
	"gamesrv/models/AnnouncementMgr"

	"github.com/astaxie/beego"
)

type UpdateAnnouncement struct {
	beego.Controller
}

func (o *UpdateAnnouncement) Post() {
	AnnouncementMgr.Instance().UpdateAnnouncement()
	beego.Debug("收到来自后台消息,更新 公告信息")
}
