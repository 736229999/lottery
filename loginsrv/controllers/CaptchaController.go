package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"loginsrv/models/Login"
	"time"

	"github.com/astaxie/beego"
	"github.com/lifei6671/gocaptcha"
)

//验证码模块；
type CaptchaController struct {
	beego.Controller
}

func (o *CaptchaController) Post() {
	reqCaptcha := &Login.ReqCaptcha{}
	//respCaptcha := &Login.RespCommon{}

	json.Unmarshal(o.Ctx.Input.RequestBody, reqCaptcha)

	//代表机器的标志位不能为空
	if reqCaptcha.Flag == "" {
		return
	}

	//将验证码保持到输出流，可以是文件或HTTP流等
	captchaBuff := bytes.NewBuffer( nil)

	//获取验证码；
	font := Login.CreateCaptcha(captchaBuff, gocaptcha.ImageFormatJpeg) // ---生成验证码放入文本流中，返回生成的验证码字符
	if font == "" {
		beego.Emergency("-------------------------- 生成验证码错误 !--------------------------")
		return
	}
	//o.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	//o.Ctx.WriteString(base64.StdEncoding.EncodeToString(captchaBuff.Bytes()))

	body, err := json.Marshal(base64.StdEncoding.EncodeToString(captchaBuff.Bytes()))
	if err != nil {
		return
	}
	o.Ctx.Output.Body(body) //向客户端返回

	Login.VerifyInfo[reqCaptcha.Flag] = font

	//验证码计时，超过一分钟验证码失效
	timer := time.NewTimer(time.Second * 60)

	//这个地方暂时这样做，分协程来处理倒计时是不可取的应该由一个协程来总负责，先完成功能再说了(正确的方式是写计时器,依次执行函数回调)
	go func(flag string, font string) {
		<-timer.C
		if f, ok := Login.VerifyInfo[flag]; ok {
			if f == font {
				Login.Destory(reqCaptcha.Flag)
			}
		}
		//beego.Debug("删除map元素 flag ：", reqCaptcha.Flag)
	}(reqCaptcha.Flag, font)
}

//Get 仅仅用户web端测试,上线记得删除
func (o *CaptchaController) Get() {
	//通过http原生方式将验证码传回去；
	Login.CreateCaptcha(o.Ctx.ResponseWriter, gocaptcha.ImageFormatJpeg)
}
