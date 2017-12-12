package gamemgr

import (
	"bytes"
	"calculsrv/models/ctrl"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

//这里统一管理所有的Game服务器
//以及负责发送消息给Game服务器

//数据库管理员类
type GameSrvMgr struct {
}

var GameSrvMap map[string]GameSrv = make(map[string]GameSrv) //Game Server url map ,k 为 服务器名字

//检查过后可用Game服务器列表(全局方便使用,如果出现异步问题,要放回类中并且加锁)
//var AvailableGameServerInfo map[string]GameServerInfo = make(map[string]GameServerInfo)

var sInstance *GameSrvMgr
var once sync.Once

func Instance() *GameSrvMgr {
	once.Do(func() {
		sInstance = &GameSrvMgr{}
	})
	return sInstance
}

func (o *GameSrvMgr) Init() error {
	//解析数据库中获得的GameServer信息成url
	//o.CheckGameSrv()
	//查看Game服务器状态,直到Game服务器有响应,才算完成初始化,才能通信
	//o.RegistCalculation()

	//发送GameServer需要的数据
	//o.UpdateGsLotteriesInfo()
	//o.UpdateGsLotteriesRecord()

	for _, v := range ctrl.GameSrv {
		var srv GameSrv

		srv.Name = v.Name
		srv.Ip = v.Ip
		srv.Port = v.Port
		srv.Url = "http://" + v.Ip + ":" + v.Port
		srv.Status = 1 //应为 ctrl 只会给状态为正常的服务器,所以这里的 status 肯定是 1 正常

		GameSrvMap[srv.Name] = srv
	}

	//---------------------检查game服务器可用状态 由于现在是根据ctrl 服的配置启动,不依赖之前的服务器启动顺序,有可能calculsrv 启动的时候game并没有启动
	//所以就不用检查服务器状态了,每次要往game发消息的时候,如果game没有响应就抛出错误

	// req := ReqCheckGameSrv{}
	// req.Cipher = conf.AesCipher

	// data, err := json.Marshal(req)

	// cipher, err := encmgr.Instance().AesPrkEnc(data)
	// if err != nil {
	// 	return err
	// }

	// body := bytes.NewBuffer(cipher)

	// //ch := make(chan int)
	// ch := make(chan int, len(GameSrvMap))
	// beego.Debug(GameSrvMap)

	// for _, v := range GameSrvMap {
	// 	go func(g GameSrv, b *bytes.Buffer) {
	// 		resp, err := httpmgr.Post(g.Url, body)
	// 		if err != nil {
	// 			g.Status = 0

	// 			beego.Error(err)
	// 		}

	// 		resp, err = encmgr.Instance().AesPrkDec(resp)
	// 		if err != nil {
	// 			g.Status = 0
	// 			beego.Error(err)
	// 		}

	// 		ret := RespCheckGameSrv{}
	// 		err = json.Unmarshal(resp, &ret)
	// 		if err != nil {
	// 			g.Status = 0
	// 			beego.Error(err)
	// 		}

	// 		if ret.Status != 1 { //1正常
	// 			g.Status = 0
	// 		}

	// 		GameSrvMap[g.Name] = g

	// 		ch <- 1
	// 	}(v, body)
	// }

	// for i := 0; i < len(GameSrvMap); i++ {
	// 	<-ch
	// }

	// beego.Debug(GameSrvMap)

	beego.Info("--- Init Game Mgr Done !")
	return nil
}

//解析从 ctrl 服中获得的GameServer信息 并检查服务器是否可用
// func (o *GameSrvMgr) CheckGameSrv() {

// 	for _, v := range ctrl.GameSrv {
// 		var srv GameSrv

// 		srv.Name = v.Name
// 		srv.Status = 1 //应为 ctrl 只会给状态为正常的服务器,所以这里的 status 肯定是 1 正常
// 		srv.Url = "http://" + v.Ip + ":" + v.Port

// 		GameSrvMap[srv.Name] = srv
// 	}

// 	//beego.Debug(GameSrvMap)
// }

// 发送消息查看Game服务器状态, 如果发现不对的 要修改管理数据库
// func (o *GameServerMgr) RegistCalculation() {
// 	//加密Game服务器注册暗号
// 	//获得单列
// 	//rsa := encrypt.Instance()
// 	cipher, err := encmgr.Instance().RsaEnc(common.CalculRegistGameSign)
// 	//解析公钥
// 	//rsa.DecodePuk(gb.RsaCalToGamePublicKey)
// 	//加密原文
// 	//cipher, err := rsa.RsaEncrypt(common.CalculRegistGameSign)
// 	if err != nil {
// 		beego.Emergency(err)
// 		return
// 	}

// 	data := &gb.ServerRegistInfo{}
// 	data.Id, _ = beego.AppConfig.Int("ServerId")
// 	data.Cipher = cipher
// 	//从配置文件读取服务器坚挺的端口
// 	data.Port, _ = beego.AppConfig.Int("httpport")

// 	b, err := json.Marshal(data)
// 	if err != nil {
// 		beego.Debug("------------------------- ", err, " -------------------------")
// 	}

// 	body := bytes.NewBuffer(b)
// 	//通道用于并发
// 	var ch = make(chan string, len(o.gameServerInfoMap))

// 	for _, v := range o.gameServerInfoMap {
// 		//注意:循环开携程的时候,不能信赖外部变量
// 		go func(info GameServerInfo, body bytes.Buffer) {
// 			resp, err := o.sendMsgToGameServer(info.Url+"/RegistCalculation", &body, CheckGameServerTimeOut)

// 			if err != nil {
// 				beego.Error("------------------------- RegistCalculation() Error ! Server Id : ", info.Id, "Error : ", err, " -------------------------")
// 				//改变服务器状态；这个地方在所有功能完成后来改,这里要把自身从 Available 里面去掉,然后重新检查服务器状态;
// 				ctrl.Instance().ChangeGameServerStatus(info.Id, 2)
// 			} else {
// 				var tempData map[string]interface{}
// 				err := json.Unmarshal(resp, &tempData)
// 				if err != nil {
// 					beego.Error("------------------------- RegistCalculation : ", err, "-------------------------")
// 				} else {
// 					if status, ok := tempData["Status"]; ok {
// 						if status == "OK" {
// 							AvailableGameServerInfo[info.Id] = info
// 						} else {
// 							beego.Error(("------------------------- RegistCalculation : Error ! -------------------------"))
// 						}
// 					}
// 				}
// 			}
// 			ch <- info.Name
// 		}(v, *body)
// 	}

// 	for i := 0; i < len(o.gameServerInfoMap); i++ {
// 		<-ch
// 	}

// 	beego.Info("-------------------------  Regist Calculation Done ! -------------------------")
// }

//发送信息给所有状态正常的服务器 r 代表router 消息名字, 所有发送给GameServer的消息 GameServer统一返回[Status]= OK ,resp为空即代表错误
// func (o *GameSrvMgr) SendMsgToGameServers(r string, body *bytes.Buffer) {
// 	//通道用于并发
// 	//var ch = make(chan int, len(AvailableGameServerInfo))
// 	for _, v := range AvailableGameServerInfo {
// 		go func(g GameServerInfo, body bytes.Buffer) {
// 			_, err := o.sendMsgToGameServer(g.Url+r, &body, DefaultTimeOut)
// 			if err != nil {
// 				beego.Debug("------------------------- Send Msg to Game Server Error ! Server Id : ", g.Name, " -- ", err, " -------------------------")
// 				//改变服务器状态；这个地方在所有功能完成后来改,这里要把自身从 Available 里面去掉,然后重新检查服务器状态;
// 				//ctrl.Instance().ChangeGameServerStatus(g, 2)
// 			}
// 		}(v, *body)
// 	}
// }

//向一个GameServer发送消息(发送标准格式)
// func (o *GameSrvMgr) sendMsgToGameSrv(url string, data *bytes.Buffer, timeOut string) ([]byte, error) {
// 	to, err := time.ParseDuration(timeOut)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c := &http.Client{
// 		Timeout: to}
// 	resp, err := c.Post(url, "application/json;charset=utf-8", data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return body, nil
// }

type GameSrv struct {
	Name   string //服务器名
	Ip     string
	Port   string
	Url    string //服务器Url
	Status int    //服务器状态 0关闭 1为正常 2为维护
}

type ReqCheckGameSrv struct {
	Cipher []byte
}

type RespCheckGameSrv struct {
	Status int //0关闭 1正常, 2维护}
}

//发送信息给所有状态正常的服务器 r 代表router 消息名字, 所有发送给GameServer的消息 GameServer统一返回[Status]= OK ,resp为空即代表错误
func (o *GameSrvMgr) SendMsgToGameServers(r string, body *bytes.Buffer) {
	//通道用于并发
	//var ch = make(chan int, len(AvailableGameServerInfo))
	for _, v := range GameSrvMap {
		go func(g GameSrv, body bytes.Buffer) {
			_, err := o.sendMsgToGameSrv(g.Url+r, &body, DefaultTimeOut)
			if err != nil {
				beego.Error("--- Send Msg to Game Server Error ! Server Id : ", g.Name, " -- ", err)
				//改变服务器状态；这个地方在所有功能完成后来改,这里要把自身从 Available 里面去掉,然后重新检查服务器状态;
				//ctrl.Instance().ChangeGameServerStatus(g, 2)
			}
		}(v, *body)
	}
}

//向一个GameServer发送消息(发送标准格式)
func (o *GameSrvMgr) sendMsgToGameSrv(url string, data *bytes.Buffer, timeOut string) ([]byte, error) {
	to, err := time.ParseDuration(timeOut)
	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Timeout: to}
	resp, err := c.Post(url, "application/json;charset=utf-8", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

const (
	//默认超时时间
	DefaultTimeOut = "10s"

	//检查Game服务器超时时间
	CheckGameServerTimeOut = "5s"

	//给GameServer发送消息超时时间
	SendMsgToGameServerTimeOut = "10s"
)
