package controllers

import (
	"ctrlsrv/models/encmgr"
	"ctrlsrv/models/srvmgr"
	"encoding/json"

	"github.com/astaxie/beego"
)

//现在由于是短连接只能验证 这个来源服务器是否再允许的列表中,等完成功能后,将服务器之间改为长连接,这样可以准确的记录已注册服务器的个数,和达到服务器注册功能
type SrvRegist struct {
	beego.Controller
}

func (o *SrvRegist) Post() {
	//每条来自其他服务器的消息都要验证来源IP是否正确
	if !srvmgr.Instance().VerifySrv(o.Ctx.Input.IP()) {
		beego.Error("Wrongful server access ! IP : ", o.Ctx.Input.IP())
		return
	}

	plaintext, err := encmgr.Instance().RsaDec(o.Ctx.Input.RequestBody)
	if err != nil {
		beego.Error(err)
		return
	}

	req := ReqSrvRegist{}
	err = json.Unmarshal(plaintext, &req)
	if err != nil {
		beego.Error(err)
		return
	}

	//验证暗号
	if string(req.Cipher) != encmgr.Instance().GetRsaCipher() {
		beego.Error("RSA Cipher wrong !")
		return
	}

	//验证是什么服务器
	srv := srvmgr.Instance().SrvInfos
	for _, v := range srv {
		if v.Ip == o.Ctx.Input.IP() && v.Func == req.SrvFunc {
			if v.Status == 1 {
				resp := RespSrvRegist{}
				resp.AesPuk = encmgr.Instance().GetAesPuk()
				resp.AesPrk = encmgr.Instance().GetAesPrk()

				body, err := json.Marshal(resp)
				if err != nil {
					beego.Error(err)
					return
				}

				body, err = encmgr.Instance().RsaEnc(body)
				if err != nil {
					beego.Error(err)
					return
				}

				o.Ctx.Output.Body(body)
				beego.Info("--- Server Register Successful : ", v.Ip, " -- Server Func : ", req.SrvFunc)
				break
			}
			beego.Error("Server Status is 0 :", v)
		}
	}
}

type ReqSrvRegist struct {
	Cipher  []byte
	SrvFunc string
}

type RespSrvRegist struct {
	AesPuk string
	AesPrk string
}
