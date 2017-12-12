package ctrl

import (
	"bytes"
	"common/httpmgr"
	"encoding/json"
	"errors"
	"loginsrv/conf"
	"loginsrv/models/encmgr"
	"sync"

	"github.com/astaxie/beego"
)

type Ctrl struct {
}

var GameSrv []Srv //game 服务器
var DbSrv Srv     //数据库 服务器
var SelfSrv Srv   //本服务器信息(列入,本服务器是计算服,那么是正式 还是试玩,还是测试,都要从 ctrlsrv 获取)

var sInstance *Ctrl
var once sync.Once

func Instance() *Ctrl {
	once.Do(func() {
		sInstance = &Ctrl{}
	})

	return sInstance
}

func (o *Ctrl) Init() error {
	//去Ctrl 服注册,获取密钥信息
	req := ReqSrvRegist{} // ---声明空的ReqSrvRegist结构体，然后赋值
	req.Cipher = conf.RsaCipher
	req.SrvFunc = conf.SrvFunc

	data, err := json.Marshal(req)//---解析封装的注册服务器的密钥参数，序列化json
	beego.Debug(data,"****************************************************************************")
	if err != nil {
		return err
	}

	cipher, err := encmgr.Instance().RsaEnc(data)//---传入原密钥 返回加密后的参数信息 []byte类型
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(cipher) // ---使用加密后的参数信息 []byte类型

	//发送注册消息   ---自定义一个http的消息体用于服务器注册  返回服务器密钥
	resp, err := httpmgr.Post(conf.CtrlsrvIP+"/srvRegist", body)
	if err != nil {
		return err
	}

	resp, err = encmgr.Instance().RsaDec(resp) //---解密服务器返回的密钥
	if err != nil {
		beego.Debug("----------------------------",err)
		return err
	}

	ret := RespSrvRegist{}
	err = json.Unmarshal(resp, &ret) //---反序列化json
	if err != nil {
		return err
	}

	encmgr.Instance().SetAesPuk(ret.AesPuk)
	encmgr.Instance().SetAesPrk(ret.AesPrk)

	//------------------------------------ 从 ctrl 获取 game 服务器, 主数据库 ,等信息 ----------------

	req1 := ReqGetDependSrv{}
	req1.Func = conf.SrvFunc

	data1, err := json.Marshal(req1)

	cipher1, err := encmgr.Instance().AesPrkEnc(data1)
	if err != nil {
		return err
	}

	body1 := bytes.NewBuffer(cipher1)

	resp1, err := httpmgr.Post(conf.CtrlsrvIP+"/getDependSrv", body1)
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
		if v.Func == "game" {
			GameSrv = append(GameSrv, v)
		} else if v.Func == "db" {
			DbSrv = v
		} else if v.Func == conf.SrvFunc {
			SelfSrv = v
		}
	}

	//检查是不是所需服务器信息都有了

	if len(GameSrv) == 0 {
		return errors.New("There is no information about the Game server !")
	}

	if DbSrv.Name == "" {
		return errors.New("There is no information about the Db server !")
	}

	if SelfSrv.Name == "" {
		return errors.New("There is no information about the Self server !")
	}

	beego.Info("--- Server Name : ", SelfSrv.Name)

	beego.Info("--- Init Ctrl Done !")
	return nil
}

//---去ctrl服注册需要的参数，这里封装成struct，参数为 暗号和服务器功能
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
	Type int //0试玩, 1正式
}
