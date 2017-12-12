package controllers

import (
	"apisrv/models/dbmgr"
	"apisrv/models/encmgr"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type GetLtryHist struct {
	beego.Controller
}

func (o *GetLtryHist) Post() {
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil {
		beego.Debug(err)
		return
	}

	req := ReqLtryHist{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	//输出请求服务器IP
	beego.Info("Req Srv Ip : ", o.Ctx.Input.IP())

	resp := []RespLtryHist{}
	bsonM := bson.M{"game_name": req.GameName}
	err = dbmgr.Instance().HistColl.Find(bsonM).Sort("-expect").Skip(0).Limit(100).All(&resp) //默认查询100条
	if err != nil {
		beego.Error(err)
		return
	}

	//beego.Debug(resp)

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

type ReqLtryHist struct {
	GameName string
}

type RespLtryHist struct {
	GameName           string    `bson:"game_name"`            // 名称
	Expect             int       `bson:"expect"`               // 期次
	OpenCode           string    `bson:"open_code"`            // 开奖号码
	OpenTime           time.Time `bson:"open_time"`            // 开奖时间
	OpenTimeStamp      int64     `bson:"open_time_stamp"`      //开奖时间(时间戳)
	RecordingTime      time.Time `bson:"recording_time"`       //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	RecordingTimeStamp int64     `bson:"recording_time_stamp"` //记录时间(时间戳)
	//maxNum    int       `bson:"id"`       				   // 最大号码取值
}
