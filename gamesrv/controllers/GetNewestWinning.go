package controllers

import (
	"encoding/json"
	"gamesrv/models/GlobalData"
	"gamesrv/models/dbmgr"

	"github.com/astaxie/beego"
)

type GetNewestWinning struct {
	beego.Controller
}

func (o *GetNewestWinning) Post() {
	respGetNewestWinning := GlobalData.RespGetNewestWinning{}
	//得到最新的中奖记录 默认读取20条
	err := dbmgr.Instance().GetNewestWinning(&(respGetNewestWinning.NewestWinningInfos))
	if err != nil {
		return
	}

	for i := 0; i < len(respGetNewestWinning.NewestWinningInfos); i++ {
		str := []byte(respGetNewestWinning.NewestWinningInfos[i].Account_Name)
		respGetNewestWinning.NewestWinningInfos[i].Account_Name = string(str[0:3])
	}

	for _, v := range respGetNewestWinning.NewestWinningInfos {
		str := []byte(v.Account_Name)
		v.Account_Name = string(str[0:3])
	}

	body, err := json.Marshal(respGetNewestWinning)
	if err != nil {
		beego.Debug(err)
		return
	}
	o.Ctx.Output.Body(body)
}
