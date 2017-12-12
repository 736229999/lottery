package controllers

import (
	"encoding/json"
	"gamesrv/models/LotteryManager"
	"time"

	"github.com/astaxie/beego"
)

type GetLotteryInfo struct {
	beego.Controller
}

func (o *GetLotteryInfo) Post() {
	req := make(map[string]string)

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug("------------------------- ", err, " -------------------------")
		return
	}

	resp := make(map[string]LotteryManager.LotteryInfo)
	if t, ok := req["Type"]; ok {
		beego.Debug("------------------------- Get Lottery Info : ", t, " -------------------------")
		//Hot类比较特殊所以列出来
		switch t {
		case "All":
			resp = o.getAllLotteriesInfo()
		case "Hot":
			resp = o.getHotLotteriesInfo()
		case "EX5", "K3", "SSC", "PK10", "PL3_0", "HK6_0", "PCDD":
			resp = o.getLotteriesInfo(t)
		default:
			resp = o.getLotteryInfo(t)
		}

		if b, err := json.Marshal(resp); err == nil {
			o.Ctx.Output.Body(b)
		}
	}

	//req := make(map[string]interface{})
	// send := make(map[string]interface{})

	// req := &lottery.ReqLotteryListWithKind{}

	// json.Unmarshal(o.Ctx.Input.RequestBody, req)
	// if o.Ctx.Input.RequestBody == nil {
	// 	return
	// }

	//kind := req["kind"].(string)
	//beego.Debug(lottery.GetLotteryListWithKind(req.Kind))
	// send["lotteryList"] = LotteryManager.Instance().GetLotteryListWithKind(req.Kind)

	// bufres, _ := json.Marshal(send)
	// beego.Debug(string(bufres))
	// o.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	// o.Ctx.Output.Body(bufres)
}

//得到所有彩票信息
func (o *GetLotteryInfo) getAllLotteriesInfo() map[string]LotteryManager.LotteryInfo {
	AllLottery := make(map[string]LotteryManager.LotteryInfo)

	for _, v := range LotteryManager.Instance().GetLtryInfoMap() {
		v.ServerNowTime = time.Now()
		AllLottery[v.GameTag] = v
	}
	return AllLottery
}

//得到热门彩票类型
func (o *GetLotteryInfo) getHotLotteriesInfo() map[string]LotteryManager.LotteryInfo {
	hotLottery := make(map[string]LotteryManager.LotteryInfo)

	for _, v := range LotteryManager.Instance().GetLtryInfoMap() {
		if v.Recommend == 1 {
			v.ServerNowTime = time.Now()
			hotLottery[v.GameTag] = v
		}
	}

	return hotLottery
}

//得到彩票大类
func (o *GetLotteryInfo) getLotteriesInfo(t string) map[string]LotteryManager.LotteryInfo {
	litteries := make(map[string]LotteryManager.LotteryInfo)

	for _, v := range LotteryManager.Instance().GetLtryInfoMap() {
		if v.ParentName == t {
			v.ServerNowTime = time.Now()
			litteries[v.GameTag] = v
		}
	}

	return litteries
}

//得到具体彩票信息
func (o *GetLotteryInfo) getLotteryInfo(t string) map[string]LotteryManager.LotteryInfo {
	litteries := make(map[string]LotteryManager.LotteryInfo)

	info, err := LotteryManager.Instance().GetLtryInfo(t)
	if err != nil {
		beego.Error(err)
		return litteries
	}
	info.ServerNowTime = time.Now()
	litteries[info.GameTag] = info

	// if v, ok := LotteryManager.Instance().LotteriesInfoMap[t]; ok {
	// 	v.ServerNowTime = time.Now()
	// 	litteries[v.GameTag] = v
	// }
	return litteries
}
