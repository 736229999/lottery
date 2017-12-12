package dbmgr

import (
	"calculsrv/models/BalanceRecordMgr"
	"calculsrv/models/ctrl"
	"calculsrv/models/gb"
	"common/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//------------------------------------------------ GET ---------------------------------------
//得到一个用户信息
func GetAccountInfo(accountName string) (*gb.AccountInfo, error) {
	result := &gb.AccountInfo{}
	err := AcColl.Find(bson.M{"account_name": accountName}).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//得到一个用户信息(AccountMgr专用,以后所有的AccountInfo的地方都要使用这个新的AccountMgr里面的AccountInfo类)
func InitAccountInfo(accountName string, accountInfo interface{}) error {
	err := AcColl.Find(bson.M{"account_name": accountName}).One(accountInfo)
	if err != nil {
		return err
	}
	return nil
}

//得到所有彩票最后一期记录信息
func GetLotteriesLastRecordByGameTag(lotteries map[string]gb.LotteryInfo) map[string]gb.LotteryRecord {
	var m = make(map[string]gb.LotteryRecord)
	for _, v := range lotteries {
		m[v.GameTag] = GetLotteryLastRecordByGameTag(v.GameTag)
	}

	return m
}

//根据彩票tag得到最后一期信息记录
func GetLotteryLastRecordByGameTag(gameTag string) gb.LotteryRecord {
	var result gb.LotteryRecord
	bsonM := bson.M{"game_tag": gameTag}
	err := HistColl.Find(bsonM).Sort("-expect").One(&result)
	if err != nil {
		beego.Error(err)
	}
	return result
}

//得到指定彩票倒数多少期的历史记录	参数2为 倒数期数
func GetLotteryHistoryByGameTagAndCount(gameTag string, count int) []gb.LotteryRecord {
	var result []gb.LotteryRecord
	bsonM := bson.M{"game_tag": gameTag}
	err := HistColl.Find(bsonM).Sort("-expect").Skip(0).Limit(count).All(&result)
	if err != nil {
		beego.Debug(err)
		return nil
	}
	return result
}

//得到一个彩种开彩这期的所有订单记录
func GetLotteryOrderRecord(gameTag string, expect int) []gb.Order {
	var record []gb.Order
	bsonM := bson.M{"game_tag": gameTag, "expect": expect}
	err := OrderColl.Find(bsonM).All(&record)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return record
}

//得到一个彩种所有未开采记录
func GetLotteryUnsettledOrderRecord(gameTag string) []gb.Order {
	var record []gb.Order
	//查询条件为,指定彩票,和开奖状态为0的彩票(0就是未开采)
	bsonM := bson.M{"game_tag": gameTag, "status": 0}
	err := OrderColl.Find(bsonM).All(&record)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return record
}

//通过账号名得到这个帐号的订单(注意,之前要验证 帐号和token是否合法才能进行这个查询, 等上线后这个改为在game服务器建立缓存池,不要每次都从数据库读取)
func GetOrderByAccountName(accountName string, skip int, limit int, ret interface{}) error {
	bsonM := bson.M{"account_name": accountName}
	err := OrderColl.Find(bsonM).Sort("-expect").Skip(skip).Limit(limit).All(&ret)
	if err != nil {
		return err
	}
	return nil
}

//得到指定彩种一期的历史记录
// func GetOneLotteryRecordByExpect(gameTag string, expect int) (gb.LotteryRecord, error) {
// 	ret := gb.LotteryRecord{}
// 	bsonM := bson.M{"game_tag": gameTag, "expect": expect}
// 	err := HistColl.Find(bsonM).One(&ret)
// 	if err != nil {
// 		return ret, err
// 	}
// 	return ret, nil
// }

//------------------------------------------------ Insert ---------------------------------------
//插入一天的开奖记录.这里面要转换数据
// func InsertLotteryHistoryByDay(data gb.LotteryRecordByDayFromApi) {
// 	//由于从开采网获取的历史记录 是从晚到早,所以存库时要反过来
// 	for i := len(data.Data) - 1; i >= 0; i-- {
// 		LotteryRecord := &gb.LotteryRecord{}
// 		//游戏标志
// 		LotteryRecord.GameName = o.convertGameTag(data.Code)

// 		//开奖期数
// 		var err error
// 		LotteryRecord.Expect, err = strconv.Atoi(data.Data[i].Expect)
// 		if err != nil {
// 			beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 			return
// 		}

// 		//开奖号码
// 		LotteryRecord.OpenCode = data.Data[i].Opencode

// 		//开奖时间 注意加上时区 prc 北京时间;
// 		loc, _ := time.LoadLocation("PRC")
// 		utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.Data[i].OpenTime, loc)
// 		LotteryRecord.OpenTime = utcTime
// 		if err != nil {
// 			beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 			return
// 		}

// 		//开奖时间戳
// 		LotteryRecord.OpenTimeStamp = data.Data[i].OpenTimeStamp

// 		if err != nil {
// 			beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 			return
// 		}

// 		LotteryRecord.RecordingTime = utils.GetNowUTC8Time()
// 		//记录时间,入库时间(时间戳);
// 		LotteryRecord.RecordingTimeStamp = time.Now().Unix()
// 		//插入数据库
// 		Insert(HistColl, LotteryRecord)
// 	}
// }

//插入最新一期开奖数据
// func InsertLotteryHistoryByNewest(data gb.LotteryRecordByNewestFromApi) {
// 	LotteryRecord := &gb.LotteryRecord{}
// 	//游戏标志
// 	LotteryRecord.GameName = o.convertGameTag(data.Code)

// 	//开奖期数
// 	var err error
// 	LotteryRecord.Expect, err = strconv.Atoi(data.Open[0].Expect)
// 	if err != nil {
// 		beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 		return
// 	}

// 	//开奖号码
// 	LotteryRecord.OpenCode = data.Open[0].Opencode

// 	//开奖时间 注意加上时区 prc 北京时间;
// 	loc, _ := time.LoadLocation("PRC")
// 	utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.Open[0].OpenTime, loc)
// 	LotteryRecord.OpenTime = utcTime
// 	if err != nil {
// 		beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 		return
// 	}

// 	//开奖时间戳
// 	LotteryRecord.OpenTimeStamp = utcTime.Unix()

// 	if err != nil {
// 		beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 		return
// 	}
// 	//记录时间,入库时间;
// 	LotteryRecord.RecordingTime = utils.GetNowUTC8Time()
// 	//记录时间,入库时间(时间戳);
// 	LotteryRecord.RecordingTimeStamp = utils.GetNowUTC8Time().Unix()
// 	//插入数据库
// 	Insert(o.HistColl, LotteryRecord)
// }

// //插入一期开奖数据,这个函数用于掉期补全
// func InsertLotteryHistoryOneRecord(gameTag string, data gb.LotteryRecordFromApi) {
// 	LotteryRecord := &gb.LotteryRecord{}
// 	//游戏标志
// 	LotteryRecord.GameName = gameTag

// 	//开奖期数
// 	var err error
// 	LotteryRecord.Expect, err = strconv.Atoi(data.Expect)
// 	if err != nil {
// 		beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 		return
// 	}

// 	//开奖号码
// 	LotteryRecord.OpenCode = data.Opencode

// 	//开奖时间 注意加上时区 prc 北京时间;
// 	loc, _ := time.LoadLocation("PRC")
// 	utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.OpenTime, loc)
// 	LotteryRecord.OpenTime = utcTime
// 	if err != nil {
// 		beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 		return
// 	}

// 	//开奖时间戳
// 	LotteryRecord.OpenTimeStamp = utcTime.Unix()

// 	if err != nil {
// 		beego.Emergency("------------------------- Insert Lottery History Error ! ------------------------- ")
// 		return
// 	}
// 	//记录时间,入库时间;
// 	LotteryRecord.RecordingTime = utils.GetNowUTC8Time()
// 	//记录时间,入库时间(时间戳);
// 	LotteryRecord.RecordingTimeStamp = utils.GetNowUTC8Time().Unix()
// 	//插入数据库
// 	Insert(o.HistColl, LotteryRecord)
// }

//插入一期手动开奖的历史记录
// func (o DbMgr) InsertLotteryHistoryOneRecordByManual(gameTag string, expect int, openCode string, openTime string) {
// 	LotteryRecord := &gb.LotteryRecord{}
// 	LotteryRecord.Expect = expect
// 	LotteryRecord.GameName = gameTag
// 	LotteryRecord.OpenCode = openCode
// 	loc, _ := time.LoadLocation("PRC")
// 	utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", openTime, loc)
// 	LotteryRecord.OpenTime = utcTime
// 	LotteryRecord.OpenTimeStamp = utcTime.Unix()
// 	LotteryRecord.RecordingTime = utils.GetNowUTC8Time()
// 	LotteryRecord.RecordingTimeStamp = utils.GetNowUTC8Time().Unix()

// 	o.insert(o.HistColl, LotteryRecord)
// }

//将api获得的gameTag转换成我看得顺眼的形式,不服吗?我就是看小写开头不顺眼
// func (o DbMgr) convertGameTag(gameTag string) string {
// 	switch gameTag {
// 	case "jx11x5":
// 		return gb.EX5_JiangXi
// 	case "sd11x5":
// 		return gb.EX5_ShanDong
// 	case "sh11x5":
// 		return gb.EX5_ShangHai
// 	case "bj11x5":
// 		return gb.EX5_BeiJing
// 	case "fj11x5":
// 		return gb.EX5_FuJian
// 	case "hlj11x5":
// 		return gb.EX5_HeiLongJiang
// 	case "js11x5":
// 		return gb.EX5_JiangSu

// 	case "gxk3":
// 		return gb.K3_GuangXi
// 	case "jlk3":
// 		return gb.K3_JiLin
// 	case "ahk3":
// 		return gb.K3_AnHui
// 	case "bjk3":
// 		return gb.K3_BeiJing
// 	case "fjk3":
// 		return gb.K3_FuJian
// 	case "hebk3":
// 		return gb.K3_HeBei
// 	case "shk3":
// 		return gb.K3_ShangHai

// 	case "cqssc":
// 		return gb.SSC_ChongQing
// 	case "tjssc":
// 		return gb.SSC_TianJin
// 	case "xjssc":
// 		return gb.SSC_XinJiang
// 	//case "nmgssc":
// 	//	return gb.SSC_NeiMengGu
// 	// case "ynssc":
// 	// 	return gb.SSC_YunNan

// 	case "bjpk10":
// 		return gb.PK10_BeiJing

// 	case "pl3":
// 		return gb.PL3

// 	case "hk6":
// 		return gb.HK6

// 	default:
// 		beego.Emergency("-------------------- DbMgr Convert Game Tag Error ! --------------------")
// 		return ""
// 	}
// }

//插入一个订单
func InsertOrder(order gb.Order) {
	Insert(OrderColl, &order)
}

//批量插入订单
func BulkInsertOrder(order []gb.Order) {
	data := utils.ConvertArrayToInterface(order)
	if data == nil {
		beego.Debug("------------------------- Convert Array To Interface Error !  -------------------------\n")
		return
	}

	BulkInsert(OrderColl, &data)
}

//插入一个流水
func InsertBalanceRecord(balanceRecord BalanceRecordMgr.BalanceRecord) {
	//正式服才插入资金流水数据到管理数据库
	if ctrl.SelfSrv.Type == 1 {
		Insert(BalanceRecordColl, &balanceRecord)
	}
}

//批量插入流水
func BulkInsertBalanceRecord(balanceRecords []BalanceRecordMgr.BalanceRecord) {
	data := utils.ConvertArrayToInterface(balanceRecords)
	if data == nil {
		beego.Debug("------------------------- Convert Array To Interface Error !  -------------------------\n")
		return
	}
	//正式服才插入资金流水数据到管理数据库
	if ctrl.SelfSrv.Type == 1 {
		BulkInsert(BalanceRecordColl, &data)
	}
}

//---------------------------------------------- Updata ---------------------------------------

//更新订单(由于每个order 结果 都不一样,目前还不知有什么办法可以批量更新所有订单,同时结果不一样的情况,暂时只有一条一条的更新)
func UpdateOrder(order *gb.Order) {
	selector := bson.M{"order_number": order.OrderNumber}
	data := bson.M{"$set": bson.M{"status": order.Status, "winning_bet_num": order.WinningBetNum, "settlement": order.Settlement, "open_code": order.OpenCode, "rebate_amount": order.RebateAmount}}
	OrderColl.Update(selector, data)
}

//更新账户信息
//这里出现个大问题,出现了一个用户没有更新钱的情况,我猜测可能是应为我把这个函数写成了一个单列类的成员喊函数
//所以在多线程调用的时候,这里就出现了售票员问题,前一个用户的钱还没有更新,后一个用户的数据来了,就变成了更新后一个用户的数据
//这样前一个用户的数据就没有更新到
//现在我再调用这个函数的地方加锁,
//如果以后不出现这个问题,那么就证明猜测是正确的,那么就要把这些数据库函数都写成独立的,而不是成员函数
//实际上这个地方完全没有必要使用单列类
//现在倭修改了整个dbmgr现在不再是单列类,如果再出现这个问题那么就加锁,如果再次出现就要怀疑这个数据库命令了
func UpdateAccountInfoMoney(accountName string, money float64, totalBetAmount float64, betAmountImmdiate float64) error {
	selector := bson.M{"account_name": accountName}
	data := bson.M{"$set": bson.M{"money": money, "total_bet_amount": totalBetAmount, "bet_amount_immediate": betAmountImmdiate}}
	err := AcColl.Update(selector, data)
	if err != nil {
		return err
	}
	return nil
}
