package controllers

import (
	"github.com/astaxie/beego"
	"github.com/mojocn/base64Captcha"
	"loginsrv/models/Login"
	"encoding/json"
)

type YzmController struct{
	beego.Controller
}

// base64Captcha create http handler
func (o *YzmController) Post() {
	reqCaptcha := &Login.ReqCaptcha{}
	json.Unmarshal(o.Ctx.Input.RequestBody, reqCaptcha)
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigDigit{
		Height:     45,
		Width:      90,
		MaxSkew:    1.7,  // 图像验证码的最大干扰洗漱.
		DotCount:   800,   // 图像验证码干扰圆点的数量.
		CaptchaLen: 4,
	}

	//创建字符公式验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyC, capC := base64Captcha.GenerateCaptcha("", configC)
	//以base64编码
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)

	Login.VerifyInfo[reqCaptcha.Flag] = idKeyC

	o.Ctx.Output.Body([]byte(base64stringC))
}