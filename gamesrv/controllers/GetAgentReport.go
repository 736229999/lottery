package controllers

import (
	"encoding/json"
	"gamesrv/models/AccountMgr"
	"gamesrv/models/dbmgr"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

type GetAgentReport struct {
	beego.Controller
}

func (o *GetAgentReport) Post() {
	req := ReqGetAgentReport{}
	resp := RespGetAgentReport{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Debug(err)
		return
	}

	//得到账户
	accountInfo := AccountMgr.AccountInfo{}
	err = accountInfo.Init(req.AccountName)
	if err != nil { //1 未找到账号
		resp.Status = 1
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证token
	b := accountInfo.VerifyToken(req.Token)
	if !b { //2 token错误
		resp.Status = 2
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//验证是否代理商
	if accountInfo.Is_Agent != 1 {
		resp.Status = 3 //用户不是代理商
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//得到要查询的账号的信息
	searchAccount := AccountMgr.AccountInfo{}
	err = searchAccount.Init(req.SearchAccount)
	if err != nil {
		resp.Status = 4 //4 未找到要查询账号的信息
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	//看看是不是查自己的信息
	if searchAccount.Account_Id != accountInfo.Account_Id {
		//判断这个要查询的账号的是不是下属
		if searchAccount.Belong_Agent_Id != accountInfo.Agent_Id {
			resp.Status = 5 //5 要查找的账号不属于这个代理商
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}
		//判断这个下属账号是不是代理商
		if searchAccount.Is_Agent != 1 {
			resp.Status = 6 //6 这个查找的账号不是下属代理商
			bufres, _ := json.Marshal(resp)
			o.Ctx.Output.Body(bufres)
			return
		}
	}
	//得到要验证的账户信息
	//根据查询条件返回要查询的信息
	ret := AgentReport{}
	switch req.SearchType {
	case 0: //今天
		bm := bson.M{"agent_id": searchAccount.Account_Id}
		dbmgr.Instance().AgentCountHour.Find(bm).Sort("-count_time").One(&ret)
	case 1: //昨天
		bm := bson.M{"agent_id": searchAccount.Account_Id}
		dbmgr.Instance().AgentCountDay.Find(bm).Sort("-count_time").One(&ret)
	case 2: //本月
		bm := bson.M{"agent_id": searchAccount.Account_Id}
		dbmgr.Instance().AgentCountMonth.Find(bm).Sort("-count_time").One(&ret)
	case 3: //上月
		bm := bson.M{"agent_id": searchAccount.Account_Id}
		dbmgr.Instance().AgentCountMonth.Find(bm).Sort("-count_time").Skip(1).One(&ret)
	default:
		resp.Status = 7 //7 查询条件不正确
		bufres, _ := json.Marshal(resp)
		o.Ctx.Output.Body(bufres)
		return
	}

	bufres, _ := json.Marshal(ret)
	o.Ctx.Output.Body(bufres)
	return
}

type ReqGetAgentReport struct {
	AccountName   string `json:"accountName"`
	Token         string `json:"token"`
	Flag          string `json:"flag"`
	SearchAccount string `json:"SearchAccount"` //查询账号
	SearchType    int    `json:"searchType"`    //查询类型 0.今天  1.昨天. 2.本月 3 上月
}

type RespGetAgentReport struct {
	Status int `json:"status"` //状态码
}

//代理商报表信息
type AgentReport struct {
	Rebate                    float64 `bson:"rebate"`                      //反水
	AgentId                   int     `bson:"agent_id"`                    //代理商id
	AgentLv                   int     `bson:"agent_lv"`                    //代理商级别
	BelongAgentId             int     `bson:"belong_agent_id"`             //所属代理商id
	BelongAgent               string  `bson:"belong_agent"`                //所属代理商账号名
	RegistCount               int     `bson:"regist_count"`                //所有下级代理商,下级用户数
	BetAmount                 int     `bson:"bet_amount"`                  //所有下级用户投注金额(不包括自己)
	BetRebateAmount           int     `bson:"bet_rebate_amount"`           //所有下级用户反水金额
	BetWinningAmount          int     `bson:"bet_winning_amount"`          //所有下级赢的金额
	RechargeAmount            int     `bson:"recharge_amount"`             //所有下级充值金额
	DrawingsAmount            int     `bson:"drawings_amount"`             //所有下级提款金额
	RechargeFavorable         int     `bson:"recharge_favorable"`          //所有下级充值优惠金额
	RechargeFirstPeopleCount  int     `bson:"recharge_first_people_count"` //首冲金额人数
	DrawingAdministrativeFees int     `bson:"drawing_administrative_fees"` //行政费率
	DrawingsFee               int     `bson:"drawings_fee"`                //手续费
	BetPeopleCount            int     `bson:"bet_people_count"`            //总下注人数
	GetRebate                 float64 `bson:"get_rebate"`                  //代理商获得返利
	WinningRebate             float64 `bson:"winning_rebate"`              //中奖赔付,代理商需要承担的金额
	AgentEarnings             float64 `bson:"agent_earnings"`              //代理盈亏
	AgentMoney                float64 `bson:"agent_money"`                 //代理商钱
	ActivityAmountGiven       int     `bson:"activity_amount_given"`       //活动赠送金额
	CountTime                 int64   `bson:"count_time"`                  //统计时间
	Status                    int     `bson:"status"`                      //状态 : 1.已处理 2.待处理 3.待结算
}
