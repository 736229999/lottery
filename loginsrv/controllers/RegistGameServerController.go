package controllers

import (
	"github.com/astaxie/beego"
)

type RegistGameServerController struct {
	beego.Controller
}

func (o *RegistGameServerController) Post() {
	// data := &Login.ServerRegistInfo{}
	// err := json.Unmarshal(o.Ctx.Input.RequestBody, &data)
	// if err != nil {
	// 	beego.Debug("------------------------- ", err, "-------------------------")
	// 	return
	// }
	// //解密
	// rsa := encrypt.Instance()
	// rsa.DecodePrk(GlobalData.RsaGameToLoginPrivateKey)
	// origin, err := rsa.RsaDecrypt(data.Cipher)

	// if origin == common.GameRegistLoginSign {
	// 	Login.GameServerIp[data.Id] = o.Ctx.Input.IP() + ":" + strconv.Itoa(data.Port)
	// 	beego.Info("------------------------- Game Server Regist , ID : ", data.Id, " IP:", Login.GameServerIp[data.Id], "-------------------------")
	// } else {
	// 	beego.Error("------------------------- Game Server Regist Cipher Verify Fail -------------------------")
	// }
}
