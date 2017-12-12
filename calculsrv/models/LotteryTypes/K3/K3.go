package K3

import (
	"bytes"
	"calculsrv/models/apimgr"
	"calculsrv/models/ctrl"
	"calculsrv/models/dbmgr"
	"common/utils"
	"encoding/json"
	"strconv"
	"time"

	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/LotteryTypes/CommonFunc"
	"calculsrv/models/Order"
	"calculsrv/models/acmgr"
	"calculsrv/models/gamemgr"
	"calculsrv/models/gb"
	"calculsrv/models/ltryset"

	"github.com/astaxie/beego"
)

//11选5，所有同类型彩票公用一个类
type K3 struct {
	id            int    //彩票ID
	gameTag       string //游戏名字
	parentName    string //游戏大类名字
	frequency     string //彩票频率
	status        int    //彩票状态  0 关闭，1 正常，2 维护
	recommend     int    //是否推荐 1推荐 2 不推荐
	recommendSort int    //推荐顺序,越大月推荐
	sort          int

	currentExpect   int       //当前这期期数(最近已开奖期数)
	openCode        []int     //当前这期开奖号码(最近已开奖号码)
	openCodeString  string    //当前这期开奖号码String形式(最近已开奖号码)
	currentOpenTime time.Time //当期开奖时间(最近已开出期数的时间)

	nextExpect         int       //下期期数
	nextRequestTime    time.Time //下期请求时间
	nextOpenTime       time.Time //下期开彩时间(官方开奖整点)
	nextClosingBetTime time.Time //下期截至下注时间（下期开彩时间 - 1m）

	afterNextExpect         int       //下下期期数
	afterNextClosingBetTime time.Time //下下期间下-注截至时间

	Settings map[int]gb.LotterySettings //所有游戏玩法设置map(赔率,限额) key 为 odds_mode(玩法id)

	newestRecord gb.LtryRecordByNewest //从Api 获得的最新期数数据
}

//初始化彩票
func Init(lif gb.LotteryInfo) (*K3, error) {
	ltry := &K3{}
	ltry.gameTag = lif.GameTag
	//通过API得到最新的一期的开彩记录（按最新，带下期）
	//开始新的架构:这里是从api服务器获取信息了(赶进度, 这里先就这样写,这里正确的情况下应该是 直到请求到正确的结果为止,外部应该是一个线程初始化一个菜种,现在不是,完成功能来改)
	var err error
	ltry.newestRecord, err = apimgr.Instance().GetRecordByNewest(lif.GameTag)
	if err != nil {
		return nil, err
	}

	//更新彩票在GameSrv的历史记录

	// ltry.newestRecord = apimgr.Instance().GetRecordByNewest(ltry.gameTag)
	// if lottery.newestRecord.Rows == 0 {
	// 	beego.Emergency("------------------------- 严重问题 : ", gameTag, " Api获取最新记录失败 将彩票状态设置为 2 维护! 请立即检查 -------------------------")
	// 	//将彩票状态设为维护
	// 	//ctrl.Instance().ChangeLtryStatus(gameTag, gb.LotteryStatus_Maintain)
	// 	return nil
	// }

	//从管理中心获得的彩票状态
	ltry.id = lif.Id
	ltry.parentName = lif.ParentName
	ltry.frequency = lif.Frequency
	ltry.status = lif.Status
	ltry.recommend = lif.Recommend
	ltry.recommendSort = lif.RecommendSort
	ltry.sort = lif.Sort

	//将彩票历史发送给GameServer
	//注意 发送的顺序不能变，必须先发送历史，这样在GameServer才能正确更新历史记录
	//这个消息 要改到LotteryHistory里面去
	ltry.updateGsLtryRecord()

	//从管理数据库获取彩票玩法相关信息
	ltry.UpdateSettings()

	//更新彩票信息,并发送给GameServer
	ltry.updataInfo(ltry.newestRecord)

	//发送信息给 GameServer
	ltry.updataGsLotteryInfo()

	//补开采(停机时间开奖的)
	ltry.CheckUnsettledOrderRecord()

	return ltry, nil
}

//更新玩法相关信息(赔率,限额等)
func (o *K3) UpdateSettings() {
	if s, ok := ltryset.LotteriesSettings[o.parentName]; ok {
		o.Settings = s
	} else {
		beego.Debug("------------------------------ 赔率获取错误！ -----------------------------")
	}
}

//更新彩票信息
func (o *K3) updataInfo(newestRecord gb.LtryRecordByNewest) {
	//保存API获取的数据结构
	o.newestRecord = newestRecord
	//当前期数
	o.currentExpect = newestRecord.CurrentExpect
	//当前开彩号码
	//beego.Debug(newestRecord.Open[0].Opencode)
	o.openCode = newestRecord.OpenCode
	//当期开采号码String形式
	o.openCodeString = newestRecord.OpenCodeStr
	//当期开奖时间
	//loc, _ := time.LoadLocation("PRC")
	//t, _ := time.ParseInLocation("2006-01-02 15:04:05", newestRecord.Open[0].OpenTime, loc)
	o.currentOpenTime = newestRecord.CurrentOpenTime
	//beego.Debug("当前彩票")
	//beego.Debug(o.gameTag)
	//beego.Debug("当前开期数")
	//beego.Debug(o.currentExpect)
	//下期期数
	o.nextExpect = newestRecord.NextExpect
	//beego.Debug("下期期数")
	//beego.Debug(o.nextExpect)
	//下期请求时间
	o.nextRequestTime = newestRecord.NextOpenTime
	//beego.Debug("下期请求时间")
	//beego.Debug(o.nextRequestTime)
	//下期开彩时间
	o.nextOpenTime = CommonFunc.GetNextOpenTime(o.gameTag, o.nextExpect, o.currentOpenTime, o.nextRequestTime)
	//beego.Debug("下期开彩时间")
	//beego.Debug(o.nextOpenTime)
	//下期截至下注时间
	d, _ := time.ParseDuration("-90s")
	o.nextClosingBetTime = o.nextOpenTime.Add(d)
	//beego.Debug("下期截至下注时间")
	//beego.Debug(o.nextClosingBetTime)
	//下下期期数
	o.afterNextExpect, _ = CommonFunc.GetAfterNextExpect(o.gameTag, o.nextExpect)
	//beego.Debug("下下期期数")
	//beego.Debug(o.afterNextExpect)
	//下下期截至下注时间
	o.afterNextClosingBetTime = CommonFunc.GetAfterNextClosingBetTime(o.gameTag, o.nextExpect, o.nextClosingBetTime)
	//beego.Debug("下下期截至下注时间")
	//beego.Debug(o.afterNextClosingBetTime)
}

//开采
func (o *K3) StartLottery(newestRecord gb.LtryRecordByNewest) {
	beego.Debug("------------------------- 开彩啦 : ", o.gameTag, " -------------------------\n")
	//检查是否掉期
	o.CheckLostExpect(o.gameTag, newestRecord, o.nextExpect)

	//更新类信息
	o.updataInfo(newestRecord)

	//结算订单
	//从数据库中获取这一期这个采种所有的订单
	orders := dbmgr.GetLotteryOrderRecord(o.gameTag, o.currentExpect)
	o.SettlementOrders(orders, o.openCodeString)

	//发送给GS
	o.updataGsLotteryInfo()
}

//检查遗漏期数并且补全开彩和历史记录, 下次更新时考虑将这个功能单开携程来处理,直到正确的补完所有彩票种(现在只查看两天,而且不会重复查看)
func (o *K3) CheckLostExpect(gameTag string, newestRecord gb.LtryRecordByNewest, nextExpect int) {
	newestExpect := newestRecord.CurrentExpect
	//如果 最新获取的期数大于下棋期数,开始进入补全程序
	if newestExpect > nextExpect {
		beego.Error("发现掉期 开始进补全 !\n")
		beego.Error("", gameTag, " 最新期数为 : ", newestExpect, " 当前下期期数为 : ", nextExpect, "\n")

		//记录掉的期数
		var lostExpect []int
		//计算要补的期数(计算期数差),并得到要补的期数切片
		for {
			//得到最新获得期数的上一期期数
			lastExpect, _ := CommonFunc.GetLastExpect(gameTag, newestExpect)
			//如果最新一期的上一期不等于下期期数,证明掉期不止一期,将这期期数保存下来
			if lastExpect != nextExpect {
				lostExpect = append(lostExpect, lastExpect)
				newestExpect = lastExpect
			} else {
				//只掉一期,记录下来,然后跳出for循环
				lostExpect = append(lostExpect, lastExpect)
				break
			}
		}
		//开始获取掉期期数的信息,先按天请求当天的开采列表
		expectRecord, err := apimgr.Instance().GetLtryHistByDay(gameTag, time.Now())
		if err != nil {
			beego.Error(err)
			return
		}

		//查询掉的期数是否在记录中
		var lostExpectInfo []gb.LotteryRecordFromApi

		for _, v := range lostExpect {
			for _, i := range expectRecord.Data {
				tmpExpect, _ := strconv.Atoi(i.Expect)
				if tmpExpect == v { //招到了掉的期信息,保存掉的信息
					lostExpectInfo = append(lostExpectInfo, i)
				}
			}
		}
		if len(lostExpect) != len(lostExpectInfo) {
			//查询完查看是否所有掉期都找到了信息,如果有没找到的再去前一天的日期找
			expectRecord2, err := apimgr.Instance().GetLtryHistByDay(gameTag, utils.DateBeforeTheDay(time.Now(), 1))
			if err != nil {
				beego.Error(err)
				return
			}

			for _, v := range lostExpect {
				for _, i := range expectRecord2.Data {
					tmpExpect, _ := strconv.Atoi(i.Expect)
					if tmpExpect == v { //招到了掉的期信息,保存掉的信息
						lostExpectInfo = append(lostExpectInfo, i)
					}
				}
			}
		}

		//经过两天查找 还是没有找到数据的话,报出错误,并补全找到的数据
		if len(lostExpect) != len(lostExpectInfo) {
			beego.Error("------------------------- ", gameTag, " 掉期补全失败 请手动检查补全 !!-------------------------\n")
		}

		//开始补全找到的数据
		for _, v := range lostExpectInfo {
			//补全开奖
			//从数据库中获取这一期这个采种所有的订单
			tmpExpect, _ := strconv.Atoi(v.Expect)
			orders := dbmgr.GetLotteryOrderRecord(gameTag, tmpExpect)
			//补开采
			o.SettlementOrders(orders, v.Expect)
			//记录入库,新架构这里不需要入库,如果不出意外api服务器会自动补全的
			//dbmgr.InsertLotteryHistoryOneRecord(gameTag, v)
			//组数据发送给GS
			record := gb.LotteryRecord{}
			record.Expect = tmpExpect
			record.GameName = gameTag
			record.OpenCode = v.Opencode

			//开奖时间 注意加上时区 prc 北京时间;
			loc, _ := time.LoadLocation("PRC")
			utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.OpenTime, loc)
			record.OpenTime = utcTime

			b, err := json.Marshal(record)
			if err != nil {
				beego.Debug("-------------------------  json.Marshal() Error : ", err, " -------------------------")
			}
			body := bytes.NewBuffer(b)

			gamemgr.Instance().SendMsgToGameServers("/UpdateRecordForLostExpect", body)
		}
	}
}

//--------------------------------------- 发送到GS的消息 -------------------------------------------------------------------
//新架构 历史记录从 apisrv 去获取
func (o *K3) updateGsLtryRecord() error {

	hist, err := apimgr.Instance().GetLtryHist(o.gameTag)
	if err != nil {
		return err
	}

	//beego.Debug(hist)

	data, err := json.Marshal(hist)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(data)

	gamemgr.Instance().SendMsgToGameServers("/UpdateLotteryRecord", body)

	return nil
}

//更新GameServer的彩票信息
func (o *K3) updataGsLotteryInfo() {
	//组消息
	msg := &gb.MsgLotteryInfo{}
	msg.Id = o.id
	msg.GameTag = o.gameTag
	msg.ParentName = o.parentName
	msg.Frequency = o.frequency
	msg.Status = o.status
	msg.Recommend = o.recommend
	msg.RecommendSort = o.recommendSort
	msg.Sort = o.sort

	msg.CurrentExpect = o.currentExpect
	msg.OpenCode = o.openCode
	msg.OpenCodeString = o.openCodeString
	msg.CurrentOpenTime = o.currentOpenTime

	msg.NextExpect = o.nextExpect
	msg.NextRequestTime = o.nextRequestTime
	msg.NextOpenTime = o.nextOpenTime
	msg.NextClosingBetTime = o.nextClosingBetTime

	msg.AfterNextExpect = o.afterNextExpect
	msg.AfterNextClosingBetTime = o.afterNextClosingBetTime

	msg.Settings = o.Settings

	b, err := json.Marshal(msg)
	if err != nil {
		beego.Debug("-------------------------  json.Marshal() Error : ", err, " -------------------------")
	}
	body := bytes.NewBuffer(b)
	gamemgr.Instance().SendMsgToGameServers("/UpdateLotteryInfo", body)
}

//解析投注
func (o *K3) AnalyticalBetting(bettingInfo gb.MsgBettingInfo) int {
	//验证用户名和Token(用户管理模块，由于现在暂时没有绑定账号的功能，所以只用验证token就行，flag 变了 token肯定变) //这里要请重新登陆
	//这里应该获取账号信息以便后续验证
	accountInfo := &acmgr.AccountInfo{}
	//是否找到帐号
	err := accountInfo.Init(bettingInfo.AccountName)
	if err != nil {
		return 6
	}
	//验证帐号的token是否正确
	if !accountInfo.VerifyToken(bettingInfo.Token) {
		return 6
	}

	//验证彩票状态;
	if o.status != 1 {
		//beego.Debug("失败")
		return 1
	}

	//验证期数(现在可以投注下期期数和下下期) 如果当前时间超过了下期截至投注时间，那么可以下注下下棋
	t := utils.GetNowUTC8Time()
	if bettingInfo.Expect == o.nextExpect { //如果是下一期 是否超过下注时间
		if t.After(o.nextClosingBetTime) {
			return 2
		}
	} else if bettingInfo.Expect == o.afterNextExpect { //如果是下下期 是否超过下注时间
		if t.After(o.afterNextClosingBetTime) {
			return 2
		}
	} else {
		return 3
	}

	//验证有无订单
	if len(bettingInfo.Orders) < 1 {
		return 5
	}

	//订单数组
	orders := []gb.Order{}
	//金额流水数组
	BalanceRecourds := []BalanceRecordMgr.BalanceRecord{}

	//循环解析订单
	for _, v := range bettingInfo.Orders {
		//单注金额数不能小于等于0
		if v.SingleBetAmount <= 0 {
			return 4
		}
		order := gb.Order{}
		order.OrderType = 0 //订单类型 0 为普通用户自己下注的订单，1为智能追号
		order.AccountName = bettingInfo.AccountName
		order.GameTag = o.gameTag
		order.Expect = bettingInfo.Expect
		order.BetType = v.BetType
		order.SingleBetAmount = v.SingleBetAmount
		order.BetNums = v.BetNums
		//将现在的赔率赋值给订单(注意这里的大坑,map直接赋值是引用形式,所以必须手动copy)
		order.Odds = make(map[string]float64)
		t := o.Settings[v.BetType].OddsMap
		for k, v := range t {
			order.Odds[k] = v
		}
		//这里要验证反水是否超过自身
		if v.Rebate/100 > accountInfo.Rebate {
			return 4
		}
		order.Rebate = v.Rebate / 100
		order.Status = 0 //未结算账单
		order.OrderNumber = Order.Instance().GetOrderNumber()
		order.BettingTime = utils.GetNowUTC8Time().Unix()
		if !o.AnalyticalOrder(&order, accountInfo) {
			return 4
		}

		orders = append(orders, order)

		//流水
		BalanceRecourd := BalanceRecordMgr.BalanceRecord{}
		BalanceRecourd.Serial_Number = Order.Instance().GetOrderNumber()
		BalanceRecourd.Account_name = accountInfo.Account_Name
		BalanceRecourd.Money_Before = accountInfo.Money
		BalanceRecourd.Money = order.OrderAmount
		BalanceRecourd.Money_After = accountInfo.Money - order.OrderAmount
		BalanceRecourd.Gap_Money = 0
		BalanceRecourd.Type = 1    //1订单(我这里只有1)
		BalanceRecourd.Subitem = 1 //1投注, 2结算
		BalanceRecourd.Trading_Time = utils.GetNowUTC8Time().Unix()
		BalanceRecourd.Status = 1
		BalanceRecourd.Order_Number = order.OrderNumber
		BalanceRecourds = append(BalanceRecourds, BalanceRecourd)

		//判断玩家身上有没有这么多钱钱,并且扣钱
		if !accountInfo.DeductMoney(order.OrderAmount) {
			beego.Debug("失败")
			return 7
		}
		//记录总投注金额,
		accountInfo.Total_Bet_Amount += order.OrderAmount
		//记录及时有效投注
		accountInfo.Bet_Amount_Immediate += order.OrderAmount
	}

	//全部检查通过,开始生成订单号
	l := len(orders)
	if l == 1 {
		dbmgr.InsertOrder(orders[0])
	} else if l > 1 {
		dbmgr.BulkInsertOrder(orders)
	}

	//是正式服才插入数据
	if ctrl.SelfSrv.Type == 1 {
		//插~!~
		bl := len(BalanceRecourds)
		if bl == 1 {
			dbmgr.InsertBalanceRecord(BalanceRecourds[0])
		} else if bl > 1 {
			dbmgr.BulkInsertBalanceRecord(BalanceRecourds)
		}
	}

	//帐号信息存库
	err = accountInfo.UpdataDb()
	if err != nil {
		beego.Emergency(err)
	}
	return 0
}

//检查未结算彩票订单记录, 补全在服务器停机期间未开采订单
func (o *K3) CheckUnsettledOrderRecord() {
	//1. 找出这个菜种所有未结算订单
	unsettleOrders := dbmgr.GetLotteryUnsettledOrderRecord(o.gameTag)
	if len(unsettleOrders) < 1 {
		return
	}

	//2, 将未结算订单按照期数分批 ,key 为 期数,value 为 同一期的订单切片
	unsettleOrdersMap := make(map[int][]gb.Order)

	for _, v := range unsettleOrders {
		if _, ok := unsettleOrdersMap[v.Expect]; ok {
			unsettleOrdersMap[v.Expect] = append(unsettleOrdersMap[v.Expect], v)
		} else {
			unsettleOrdersMap[v.Expect] = []gb.Order{v}
		}
	}

	//3, 判断订单期数
	for expect, orders := range unsettleOrdersMap {
		//对比期数 == 的情况
		if o.currentExpect == expect {
			o.SettlementOrders(orders, o.openCodeString)
		} else if o.currentExpect > expect { //如果当前期数大于,订单的期数,那么要去数据库中招到对应期数的历史记录来进行结算
			lotteryRecord, err := apimgr.Instance().GetLtryRecordByExpect(o.gameTag, expect)
			if err != nil {
				beego.Error("----------------------- 严重错误 :  根据彩种,和期数获取开采记录失败 请检查 ", err, " 彩种 : ", o.gameTag, " 期数 : ", expect, " -----------------------")
				return
			}
			//3.结算同一期的这些订单
			o.SettlementOrders(orders, lotteryRecord.OpenCode)
		}
	}
	beego.Info("------------------------- ", o.gameTag, " 订单补开奖完成 !  -------------------------")
}

//得到当前期数
func (o *K3) GetCurrentExpect() int {
	return o.currentExpect
}

//得到下棋开彩时间
func (o *K3) GetNextOpenTime() time.Time {
	return o.nextOpenTime
}

//得到下期请求时间
func (o *K3) GetNextReqTime() time.Time {
	return o.nextRequestTime
}

//得到父类名字
func (o *K3) GetParentName() string {
	return o.parentName
}

//更新彩票信息(这个函数用于后台有彩票数据更新的情况,这时要更新彩票信息,并通知Gamesrv)
func (o *K3) UpdateLtryInfo(li gb.LotteryInfo) bool {
	o.recommend = li.Recommend
	o.recommendSort = li.RecommendSort
	o.sort = li.Sort

	o.updataGsLotteryInfo()
	return true
}

//更新彩票设置(赔率)
func (o *K3) UpdateLtrySet() bool {
	o.UpdateSettings()

	o.updataGsLotteryInfo()
	return true
}

//得到彩票名字
func (o K3) GetGameName() string {
	return o.gameTag
}
