package ctrl

import (
	"apisrv/conf"
	"apisrv/models/encmgr"
	"bytes"
	"common/httpmgr"
	"encoding/json"
	"errors"
	"sync"

	"github.com/astaxie/beego"
)

//所有服务器都需要使用Ctrl类来进行开机初始化,每个服务器 只知道 Ctrlsrv 的ip,其他所有信息都需要统统 Ctrlsrv 来获取.以保证数据安全
type Ctrl struct {
	Ltrys map[string]Ltry //key 为彩票名字
}

//General Manager Server (公司总后台)
var GenMgrSrv Srv //总后台管理 服务器 往后台发消息就要用这个

var sInstance *Ctrl
var once sync.Once

func Instance() *Ctrl {
	once.Do(func() {
		sInstance = &Ctrl{}
	})
	return sInstance
}

func (o *Ctrl) Init() error {
	//组发送结构
	req := ReqSrvRegist{}
	req.Cipher = conf.RsaCipher
	req.SrvFunc = conf.SrvFunc

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
	resp, err := httpmgr.Post(conf.CtrlsrvIP+"/srvRegist", body)
	if err != nil {
		return err
	}

	resp, err = encmgr.Instance().RsaDec(resp)
	if err != nil {
		return err
	}

	//获取返回
	ret := RespSrvRegist{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}

	encmgr.Instance().SetAesPuk(ret.AesPuk)
	encmgr.Instance().SetAesPrk(ret.AesPrk)

	//--------------------------- 发送获取彩票信息消息 ------------------------------

	req1 := ReqGetLtryInfo{}
	req1.Func = conf.SrvFunc

	data1, err := json.Marshal(req1)
	if err != nil {
		return err
	}

	cipher1, err := encmgr.Instance().AesPrkEnc(data1)
	if err != nil {
		return err
	}

	body1 := bytes.NewBuffer(cipher1)

	resp1, err := httpmgr.Post(conf.CtrlsrvIP+"/getLtryInfo", body1)
	if err != nil {
		return err
	}

	resp1, err = encmgr.Instance().AesPrkDec(resp1)
	if err != nil {
		return err
	}

	ret1 := []Ltry{}
	err = json.Unmarshal(resp1, &ret1)
	if err != nil {
		return err
	}

	o.Ltrys = make(map[string]Ltry)

	for _, v := range ret1 {
		o.Ltrys[v.Game_name] = v
	}

	//------------------------------------ 从 ctrl 获取 总管理服IP 六合彩等手动开奖的消息必须来自于 总管理服IP ----------------

	req2 := ReqGetDependSrv{}
	req2.Func = conf.SrvFunc

	data2, err := json.Marshal(req2)

	cipher2, err := encmgr.Instance().AesPrkEnc(data2)
	if err != nil {
		return err
	}

	body2 := bytes.NewBuffer(cipher2)

	resp2, err := httpmgr.Post(conf.CtrlsrvIP+"/getDependSrv", body2)
	if err != nil {
		return err
	}

	resp2, err = encmgr.Instance().AesPrkDec(resp2)
	if err != nil {
		return err
	}

	ret2 := []Srv{}
	err = json.Unmarshal(resp2, &ret2)
	if err != nil {
		return err
	}

	for _, v := range ret2 {

		if v.Func == "genmgr" {
			GenMgrSrv = v
		}
	}

	//检查是不是所需服务器信息都有了
	if GenMgrSrv.Name == "" {
		return errors.New("There is no information about the General Mgr Srv server !")
	}

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

type ReqGetLtryInfo struct {
	Func string
}

type Ltry struct {
	Id           int
	Name         string //彩票中文名
	Game_name    string //彩票英文名
	Freq         int    //频率: 0.低频 1.高频
	Api_code_kcw string //kcw 代表这个是开采网的Api彩票请求代码,这个值是要组合在url中使用的
}

//请求依赖服务器信息
type ReqGetDependSrv struct {
	Func string
}

type Srv struct {
	Name string //服务器名字
	Ip   string
	Port string
	Func string //服务器功能
	Type int    //服务器类型  1 正式 0 测试
}
