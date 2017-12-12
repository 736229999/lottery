package Utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

//得到该服务器是否正式服
func IsFormalServer() bool {
	//判断改服务器是试玩还是正式服;
	serverType, err := beego.AppConfig.Int("ServerType")
	if err != nil {
		beego.Emergency(err)
	}

	if serverType == 0 {
		return false
	} else {
		return true
	}
}

//加彩种要加这里
func VerifyGameTag(gameTag string) bool {
	switch gameTag {
	case "EX5_JiangXi":
		return true
	case "EX5_ShanDong":
		return true
	case "EX5_ShangHai":
		return true
	case "EX5_BeiJing":
		return true
	case "EX5_FuJian":
		return true
	case "EX5_HeiLongJiang":
		return true
	case "EX5_JiangSu":
		return true

	case "K3_GuangXi":
		return true
	case "K3_JiLin":
		return true
	case "K3_AnHui":
		return true
	case "K3_BeiJing":
		return true
	case "K3_FuJian":
		return true
	case "K3_HeBei":
		return true
	case "K3_ShangHai":
		return true
	case "K3_JiangSu":
		return true

	case "SSC_ChongQing":
		return true
	case "SSC_TianJin":
		return true
	case "SSC_XinJiang":
		return true
	case "SSC_NeiMengGu":
		return true
	case "SSC_YunNan":
		return true

	case "PK10_BeiJing":
		return true

	case "PK10_F":
		return true

	case "PL3":
		return true

	case "HK6":
		return true

	case "BJ28":
		return true

	default:
		return false
	}
}

//gamesrv启动完毕,像LoginServer发送消息
// func Registgamesrv() {
// 	var loginServerIp []string

// 	if IsFormalServer() {
// 		//正式服向两个正式Login注册
// 		loginServerIp = append(loginServerIp, GlobalData.FormalLoginServerIp_0)
// 		loginServerIp = append(loginServerIp, GlobalData.FormalLoginServerIp_1)
// 	} else {
// 		//试玩服向一个试玩服Login注册
// 		loginServerIp = append(loginServerIp, GlobalData.TrialLoginServerIp)
// 	}

// 	//加密Game服务器注册暗号
// 	//获得单列
// 	rsa := encrypt.Instance()
// 	//解析密钥
// 	rsa.DecodePuk(GlobalData.RsaGameToLoginPublicKey)
// 	//加密原文
// 	cipher, err := rsa.RsaEncrypt(common.GameRegistLoginSign)
// 	if err != nil {
// 		beego.Emergency(err)
// 		return
// 	}

// 	info := &GlobalData.ServerRegistInfo{}
// 	info.Id, _ = beego.AppConfig.Int("ServerId")
// 	info.Cipher = cipher
// 	//从配置文件读取服务器坚挺的端口
// 	serverType, err := beego.AppConfig.Int("httpport")
// 	if err != nil {
// 		beego.Emergency(err)
// 	}
// 	info.Port = serverType

// 	msg, err := json.Marshal(info)
// 	if err != nil {
// 		return
// 	}

// 	body := bytes.NewBuffer(msg)

// 	//通道用于并发
// 	var ch = make(chan int, len(loginServerIp))

// 	for _, v := range loginServerIp {
// 		go func(url string, body bytes.Buffer) {
// 			//超时时间设定为5秒
// 			sendMsgToLoginServer("http://"+url+"/Registgamesrv", &body, "5s")
// 			ch <- 1
// 		}(v, *body)
// 	}

// 	for i := 0; i < len(loginServerIp); i++ {
// 		<-ch
// 	}

// 	beego.Info("-------------------------  gamesrv Regist To Login Server OK ! -------------------------")
// }

//向一个LoginServer发送消息(发送标准格式)
func sendMsgToLoginServer(url string, data *bytes.Buffer, timeOut string) []byte {
	to, err := time.ParseDuration(timeOut)
	if err != nil {
		beego.Error("------------------------- time.ParseDuration() : ", err, " -------------------------")
		return nil
	}

	c := &http.Client{
		Timeout: to}
	resp, err := c.Post(url, "application/json;charset=utf-8", data)
	if err != nil {
		beego.Error("------------------------- sendMsgToLoginServer() Post Error ! Server  : ", err, " -------------------------")
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("------------------------- sendMsgToLoginServer() ioutil.ReadAll(resp.Body) Error ! :", err, "-------------------------")
		return nil
	}
	return body
}

//验证
