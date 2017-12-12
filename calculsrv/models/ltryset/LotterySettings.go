package ltryset

import (
	"bytes"
	"calculsrv/models/dbmgr"
	"calculsrv/models/gamemgr"
	"calculsrv/models/gb"
	"common/utils"
	"encoding/json"

	"github.com/astaxie/beego"
)

var LotteriesSettings map[string]map[int]gb.LotterySettings = make(map[string]map[int]gb.LotterySettings)

func Init() error {
	data, err := dbmgr.GetLtrySet()
	if err != nil {
		//beego.Debug("------------------------- GetLotterySettings Error : ", err, " -------------------------")
		return err
	}

	//将数据库中得到的数据转换为统一的格式
	for _, v := range data {
		var s gb.LotterySettings
		if _, ok := LotteriesSettings[v["lottery_name"].(string)]; !ok {
			LotteriesSettings[v["lottery_name"].(string)] = make(map[int]gb.LotterySettings)
		}

		s.Name = v["lottery_name"].(string)
		s.Id = v["odds_mode"].(int)
		s.SingleLimit = float64(v["quota_single"].(int))
		s.OrderLimit = float64(v["quota_bet"].(int))
		s.OddsMap = utils.OddsInterface2Map(v["odds_value"])

		LotteriesSettings[v["lottery_name"].(string)][s.Id] = s
	}

	//发送给Game服务器以便用户频繁获取信息
	b, _ := json.Marshal(LotteriesSettings)
	body := bytes.NewBuffer(b)
	gamemgr.Instance().SendMsgToGameServers("/UpdateLotterySettings", body)

	beego.Info("---  Lottery Settings Init Done !")
	return nil
}
