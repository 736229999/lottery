package controllers

import (
	"bytes"
	"gamesrv/models/ctrl"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

type LotteryBetting struct {
	beego.Controller
}

func (o *LotteryBetting) Post() {
	//req := MsgBettingInfo{}
	//send := make(map[string]interface{})

	// if err := json.Unmarshal(o.Ctx.Input.RequestBody, &req); err != nil {
	// 	beego.Debug("------------------------- ", err, " -------------------------")
	// 	return
	// }

	// beego.Debug(req.GameTag)
	// beego.Debug(req.Expect)
	// beego.Debug(req.Orders)
	// beego.Debug(req.UserAccount)
	// beego.Debug(req.Orders[0].BetNums)

	// req := &msgTest{}
	// req.test = time.Now()

	// b, err := json.Marshal(req)
	// if err != nil {
	// 	beego.Debug("------------------------- ", err, " -------------------------")
	// }

	reqBody := bytes.NewBuffer(o.Ctx.Input.RequestBody)

	to, err := time.ParseDuration("5s")
	if err != nil {
		beego.Debug("------------------------- time.ParseDuration() : ", err, " -------------------------")
	}

	c := &http.Client{
		Timeout: to}
	resp, err2 := c.Post("http://"+ctrl.CalculSrv.Ip+":"+ctrl.CalculSrv.Port+"/LotteryBetting", "application/json;charset=utf-8", reqBody)
	if err2 != nil {
		beego.Error(err2)
	}

	defer resp.Body.Close()

	respBody, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		beego.Error(err)
	}
	o.Ctx.Output.Body(respBody)
}

type msgTest struct {
	test time.Time
}

//订单信息（一次投注）
type MsgBettingInfo struct {
	UserAccount string     `json:"userAccount"` //用户名
	Token       string     `json:"Token"`
	GameTag     string     `json:"gameTag"` //游戏名称
	Expect      int        `json:"expect"`  //对应彩票期数
	Orders      []MsgOrder `json:"infos"`   //订单信息
}

//单个订单信息
type MsgOrder struct {
	BetType         int     `json:"bettingType"`     //投注类型
	SingleBetAmount float64 `json:"singleBetAmount"` //单注金额
	BetNums         string  `json:"betNums"`         //投注数字
	Odds            float64 `json:"odds"`            //赔率
	Rebate          float64 `json:"rebate"`          //反水
}
