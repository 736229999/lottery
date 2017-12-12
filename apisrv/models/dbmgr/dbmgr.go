package dbmgr

import (
	"apisrv/conf"
	"apisrv/models/encmgr"
	"bytes"
	"common/httpmgr"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

type DbMgr struct {
	apiSess *mgo.Session  //Api服数据库session
	apiDb   *mgo.Database //Api数据库

	HistColl *mgo.Collection //历史记录信息表名
}

var sInstance *DbMgr
var once sync.Once

func Instance() *DbMgr {
	once.Do(func() {
		sInstance = &DbMgr{}
	})

	return sInstance
}

func (o *DbMgr) Init() error {
	//去控制服获取api信息
	req := ReqGetDependSrv{}
	req.Func = conf.SrvFunc

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(cipher)

	resp, err := httpmgr.Post(conf.CtrlsrvIP+"/getDependSrv", body)
	if err != nil {
		return err
	}

	resp, err = encmgr.Instance().AesPrkDec(resp)
	if err != nil {
		return err
	}

	ret := []RespGetDependSrv{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}

	for _, v := range ret {
		if v.Func == "db" {
			dial := &mgo.DialInfo{
				Addrs:    []string{v.Ip + ":" + v.Port},
				Timeout:  time.Second * 3,
				Database: ApiDbName,
				Username: apisrvDbUserName,
				Password: apisrvDbPwd,
			}

			//连接计算服数据库（授权访问）
			o.apiSess, err = mgo.DialWithInfo(dial)
			if err != nil {
				return err
			}

			// 采用 Strong 模式
			o.apiSess.SetMode(mgo.Strong, true)

			o.apiDb = o.apiSess.DB(ApiDbName)

			o.HistColl = o.apiDb.C(histColl)

			beego.Info("--- Init Db Mgr Done !")

			return nil
		}
	}
	return errors.New("There is no database server on the depend server .")
}

type ReqGetDependSrv struct {
	Func string
}

type RespGetDependSrv struct {
	Name string
	Ip   string
	Port string
	Func string
}
