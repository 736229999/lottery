package controllers

import (
	"encoding/json"
	"indsrv/models/dbmgr"
	"indsrv/models/encmgr"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type GetHist struct {
	beego.Controller
}

func (o *GetHist) Post() {
	//这里就不验证是否是ctrl 服里面的服务器了 应为有消息加密 在这里加入ip输出 就可以知道有没有其他人在访问我的api服务器
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		//输出请求服务器IP
		beego.Error(err, " --- Req Srv Ip : ", o.Ctx.Input.IP())
		return
	}

	req := ReqGetHist{}

	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	resp := []RespGetHist{}
	bsonM := bson.M{"game_name": req.GameName}
	err = dbmgr.PK10Coll.Find(bsonM).Sort("-expect").Skip(0).Limit(100).All(&resp) //默认查询100条
	if err != nil {
		beego.Error(err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		beego.Error(err)
		return
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		beego.Debug(err)
		return
	}

	o.Ctx.Output.Body(cipher)
}

type ReqGetHist struct {
	GameName string
}

type RespGetHist struct {
	GameName           string    `bson:"game_name"`            // 名称
	Expect             int       `bson:"expect"`               // 期次
	OpenCode           string    `bson:"open_code"`            // 开奖号码
	OpenTime           time.Time `bson:"open_time"`            // 开奖时间
	OpenTimeStamp      int64     `bson:"open_time_stamp"`      //开奖时间(时间戳)
	RecordingTime      time.Time `bson:"recording_time"`       //记录时间实际获取结果时间
	RecordingTimeStamp int64     `bson:"recording_time_stamp"` //记录时间(时间戳)
}
