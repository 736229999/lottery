package controllers

import (
	"encoding/json"
	"indsrv/models/dbmgr"
	"indsrv/models/encmgr"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type GetRecordByExpect struct {
	beego.Controller
}

func (o *GetRecordByExpect) Post() {
	plaintext, err := encmgr.Instance().AesPrkDec(o.Ctx.Input.RequestBody)
	if err != nil { //解密错误输出请求服务器ip
		beego.Error(err, " --- Req Srv Ip : ", o.Ctx.Input.IP())
		return
	}

	req := ReqRecordByExpect{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	resp := RespRecord{}

	err = dbmgr.PK10Coll.Find(bson.M{"game_name": req.GameName, "expect": req.Expect}).One(&resp)
	if err != nil {
		beego.Error("Record search error , Game : ", req.GameName, " Expect : ", req.Expect)
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

type ReqRecordByExpect struct {
	GameName string
	Expect   int
}

//一条开奖记录(从计算服务器数据库获取,存数据库也是这个结构)
type RespRecord struct {
	GameName           string    `bson:"game_name"`            // 名称
	Expect             int       `bson:"expect"`               // 期次
	OpenCode           string    `bson:"open_code"`            // 开奖号码
	OpenTime           time.Time `bson:"open_time"`            // 开奖时间
	OpenTimeStamp      int64     `bson:"open_time_stamp"`      //开奖时间(时间戳)
	RecordingTime      time.Time `bson:"recording_time"`       //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	RecordingTimeStamp int64     `bson:"recording_time_stamp"` //记录时间(时间戳)
}
