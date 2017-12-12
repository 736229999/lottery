package controllers

import (
	"github.com/astaxie/beego"
)

type RegistCalculation struct {
	beego.Controller
}

func (o *RegistCalculation) Post() {
	// data := &GlobalData.ServerRegistInfo{}
	// err := json.Unmarshal(o.Ctx.Input.RequestBody, &data)
	// if err != nil {
	// 	beego.Debug("------------------------- ", err, "-------------------------")
	// 	return
	// }
	// //解密
	// rsa := encrypt.Instance()
	// rsa.DecodePrk(GlobalData.RsaCalToGamePrivateKey)
	// origin, err := rsa.RsaDecrypt(data.Cipher)

	// if origin == common.CalculRegistGameSign {
	// 	GlobalData.CalculationServerIp = o.Ctx.Input.IP() + ":" + strconv.Itoa(data.Port)
	// 	beego.Info("--------- Calculation Server Regist , ID : ", data.Id, " IP: ", GlobalData.CalculationServerIp, "-----------")
	// 	//回CalculationServer 消息
	// 	send := make(map[string]interface{})
	// 	send["Status"] = "OK"
	// 	bufres, _ := json.Marshal(send)
	// 	o.Ctx.Output.Body(bufres)
	// } else {
	// 	beego.Error("-------- Calculation Server Regist Cipher Verify Fail ---------")
	// }
}
