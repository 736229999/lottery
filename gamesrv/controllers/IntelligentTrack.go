package controllers

import (
	"github.com/astaxie/beego"
)

type IntelligentTrack struct {
	beego.Controller
}

func (o *IntelligentTrack) Post() {

	// reqBody := bytes.NewBuffer(o.Ctx.Input.RequestBody)

	// to, err := time.ParseDuration("5s")
	// if err != nil {
	// 	beego.Debug("------------------------- time.ParseDuration() : ", err, " -------------------------")
	// }

	// c := &http.Client{
	// 	Timeout: to}
	// resp, err2 := c.Post("http://"+GlobalData.CalculationServerIP+"/IntelligentTrack", "application/json;charset=utf-8", reqBody)
	// if err2 != nil {
	// 	beego.Debug(err)
	// }

	// defer resp.Body.Close()

	// respBody, err3 := ioutil.ReadAll(resp.Body)
	// if err3 != nil {
	// 	beego.Debug(err)
	// }
	// o.Ctx.Output.Body(respBody)
}
