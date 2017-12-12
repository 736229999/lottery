package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/AnnouncementMgr"

	"github.com/astaxie/beego"
)

type GetAnnouncement struct {
	beego.Controller
}

func (o *GetAnnouncement) Post() {
	req := ReqGetAnnouncement{}
	resp := []AnnouncementMgr.Announcement{}

	respStatus := make(map[string]interface{})

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug(err)
		return
	}

	//用户名为空代表未登录情况下的公告获取
	if req.AccountName == "" {
		for _, v := range AnnouncementMgr.Instance().AllAnnouncements {
			resp = append(resp, v)
		}

		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil { //1 未找到账号
		respStatus["status"] = 1
		bufres, _ := json.Marshal(respStatus)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(req.Token)
	if !b { //2 token错误
		respStatus["status"] = 2
		bufres, _ := json.Marshal(respStatus)
		o.Ctx.Output.Body(bufres)
		return
	}

	//根据不同得条件来获得信息
	//首先查看 面向所有用户的公告中,哪些对应平台的公告
	for _, v := range AnnouncementMgr.Instance().AllAnnouncements {
		resp = append(resp, v)
	}

	//查看面向用户组公告
	for _, v := range AnnouncementMgr.Instance().GroupAnnouncements {
		if v.Platform == req.Platform {
			for _, i := range v.Group {
				if i == accountInfo.Group {
					resp = append(resp, v)
				}
			}
		}
	}

	//查看面向特定玩家公告
	for _, v := range AnnouncementMgr.Instance().AccountAnnouncements {
		if v.Platform == req.Platform {
			if req.AccountName == v.Accounts {
				resp = append(resp, v)
			}
		}
	}

	bufres, _ := json.Marshal(resp)
	o.Ctx.Output.Body(bufres)
}

type ReqGetAnnouncement struct {
	AccountName string `json:"accountName"`
	Token       string `json:"token"`
	Flag        string `json:"flag"`
	Platform    string `json:"platform"` //all , ios ,android
}
