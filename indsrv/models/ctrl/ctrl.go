package ctrl

import (
	"bytes"
	"common/httpmgr"
	"encoding/json"
	"errors"
	"indsrv/models/encmgr"
	"indsrv/models/scy"

	"github.com/astaxie/beego"
)

var DbSrv Srv   //数据库 服务器
var MgrSrv Srv  //后台管理 服务器 只能接收来自这个服务器的信息
var SelfSrv Srv //本服务器信息(列入,本服务器是计算服,那么是正式 还是试玩,还是测试,都要从 ctrlsrv 获取)

func Init() error {
	//去Ctrl 服注册,获取密钥信息
	req := ReqSrvRegist{}
	req.Cipher = scy.RsaCipher
	req.SrvFunc = scy.SrvFunc

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	cipher, err := encmgr.Instance().RsaEnc(data)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(cipher)

	//发送注册消息
	resp, err := httpmgr.Post(scy.CtrlsrvIP+"/srvRegist", body)
	if err != nil {
		return err
	}

	resp, err = encmgr.Instance().RsaDec(resp)
	if err != nil {
		return err
	}

	ret := RespSrvRegist{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}

	encmgr.Instance().SetAesPuk(ret.AesPuk)
	encmgr.Instance().SetAesPrk(ret.AesPrk)

	//------------------------------------ 从 ctrl 获取 数据库信息,管理服信息 ----------------

	req1 := ReqGetDependSrv{}
	req1.Func = scy.SrvFunc

	data1, err := json.Marshal(req1)

	cipher1, err := encmgr.Instance().AesPrkEnc(data1)
	if err != nil {
		return err
	}

	body1 := bytes.NewBuffer(cipher1)

	resp1, err := httpmgr.Post(scy.CtrlsrvIP+"/getDependSrv", body1)
	if err != nil {
		return err
	}

	resp1, err = encmgr.Instance().AesPrkDec(resp1)
	if err != nil {
		return err
	}

	ret1 := []Srv{}
	err = json.Unmarshal(resp1, &ret1)
	if err != nil {
		return err
	}

	for _, v := range ret1 {
		if v.Func == "db" {
			DbSrv = v
		} else if v.Func == "mgr" {
			MgrSrv = v
		} else if v.Func == scy.SrvFunc {
			SelfSrv = v
		}
	}

	//检查是不是所需服务器信息都有了

	if DbSrv.Name == "" {
		return errors.New("There is no information about the Db server !")
	}

	if MgrSrv.Name == "" {
		return errors.New("There is no information about the Mgr server !")
	}

	if SelfSrv.Name == "" {
		return errors.New("There is no information about the Self server !")
	}

	//------------------------------------ 从 管理服数据库 获取 彩票信息 应为每一个代理商的信息都是不一样的 ----------------

	beego.Info("--- Server Name : ", SelfSrv.Name)

	beego.Info("--- Init Ctrl Done !")
	return nil
}

type ReqSrvRegist struct {
	Cipher  []byte
	SrvFunc string
}

type RespSrvRegist struct {
	AesPuk string //aes 公钥
	AesPrk string //aes 私钥
}

//请求依赖服务器信息
type ReqGetDependSrv struct {
	Func string
}

type Srv struct {
	Name string
	Ip   string
	Port string
	Func string
	Type int
}

type ReqGetLtryInfo struct {
	Func string
}
