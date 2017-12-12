package HistoryManager

// //彩票历史；
// type HistoryMgr struct {
// 	//如果 每次历史查询会损耗很多效率,再放在内存中来操作
// 	//historyList map[string]list.List //彩票实时历史记录,和数据库是同步的,k 为 gameTag ,v 为队列(注意这里list 当作队列使用,list 的每一个元素就是数据库的一条记录)

// }

// var sInstance *HistoryMgr
// var once sync.Once

// func Instance() *HistoryMgr {
// 	once.Do(func() {
// 		sInstance = &HistoryMgr{}
// 		sInstance.init()
// 	})

// 	return sInstance
// }

// func (o *HistoryMgr) init() {
// 	for _, v := range ctrl.LotteriesInfoMap {
// 		beego.Info("------------------------- 开始补全 ", v.GameTag, " ! -------------------------")
// 		o.MakeUpRecord(v.GameTag)
// 	}
// 	beego.Info("------------------------- 所有彩票历史记录补全完成 ! -------------------------\n")
// }

// //启动时 补全彩种历史记录 (函数有效率问题 现在懒得改,反正也就启动执行一次)
// func (o *HistoryMgr) MakeUpRecord(gameTag string) {
// 	//得到彩票信息
// 	lotteryInfo := ctrl.LotteriesInfoMap[gameTag]
// 	//得到最后一期的记录
// 	lastRecord := dbmgr.GetLotteryLastRecordByGameTag(lotteryInfo.GameTag)

// 	//首先判断彩种是高频彩还是低频彩(高频彩要补足3天的记录，低频彩票要补足 100期)
// 	if lotteryInfo.Frequency == "high" {
// 		//现在北京时间
// 		now := Utils.GetNowUTC8Time()
// 		//最后一期的开采时间(转换为北京时间)
// 		recordTime := Utils.ConvertToUTC8Time(lastRecord.OpenTime)

// 		//查看最后一期的日期,如果没有历史记录或是历史记录距离现在超过3天,或者天就要补助3天加上今天启动时的开奖记录
// 		if lastRecord.Expect == 0 || Utils.DateSub(now, recordTime) >= 4 {
// 			beego.Info("------------------------- 数据库中没有 ", lotteryInfo.GameTag, "的历史记录, 该彩种为高频菜,开始补全 3 天的历史记录. -------------------------")
// 			//计算时间差 然后获取差值 存库(注意存库时 时间是从早到晚的顺序,最新的记录在数据库的最后才是正确的)
// 			for i := 3; i >= 0; i-- {
// 				//获取结果往日结果
// 				result := LotteryApiMgr.Instance().GetRecordByDate(lotteryInfo.GameTag, Utils.DateBeforeTheDay(Utils.GetNowUTC8Time(), i))
// 				//如果没有数据返回立马将彩票状态变更为维护
// 				if result.Rows == 0 && result.Code == "" {
// 					beego.Emergency("------------------------- 严重问题 : ", lotteryInfo.GameTag, " Api获取历史记录失败 将彩票状态设置为 2 维护! 请立即检查 -------------------------")
// 					ctrl.Instance().ChangeLotteryStatus(lotteryInfo.GameTag, gb.LotteryStatus_Maintain)
// 					return
// 				} else {
// 					//如果数据正确存库
// 					dbmgr.InsertLotteryHistoryByDay(result)
// 					//放入内存中
// 					//beego.Info("-------------------- 数据正确 存入数据库 ! -------------------------")
// 				}
// 			}
// 		} else if Utils.DateSub(Utils.GetNowUTC8Time(), Utils.ConvertToUTC8Time(lastRecord.OpenTime)) == 0 {
// 			//数据库中查询到历史记录,如果最后一次记录是同一天, 查看期数 补全
// 			result := LotteryApiMgr.Instance().GetRecordByDate(lotteryInfo.GameTag, Utils.GetNowUTC8Time())
// 			//由于服务器启动时间有可能今天并没有记录 比如新疆时时彩 要上午10点才开奖,9点启动服务器就会找不到今天的记录,所以rows == 0 并不代表数据错误,加上code 为空的判断才能说明数据是否正确(还是有点累赘后面来改)
// 			if result.Rows == 0 && result.Code == "" {
// 				beego.Emergency("------------------------- 严重问题 : ", lotteryInfo.GameTag, " Api获取历史记录失败 将彩票状态设置为 2 维护! 请立即检查 -------------------------")
// 				//通过控制中心改变彩票状态
// 				ctrl.Instance().ChangeLotteryStatus(lotteryInfo.GameTag, gb.LotteryStatus_Maintain)
// 				return
// 			} else {
// 				//如果数据正确存库
// 				dbmgr.InsertLotteryHistoryByDay(o.FilterRecord(result, lastRecord.Expect))
// 				//beego.Info("-------------------- 数据正确 存入数据库 ! -------------------------")
// 			}
// 		} else {
// 			//数据库中查询到历史记录,但没有超过3天,然后依次补全到今天
// 			date := Utils.DateSub(Utils.GetNowUTC8Time(), Utils.ConvertToUTC8Time(lastRecord.OpenTime))
// 			//获取结果往日结果
// 			for i := date; i >= 0; i-- {
// 				result := LotteryApiMgr.Instance().GetRecordByDate(lotteryInfo.GameTag, Utils.DateBeforeTheDay(Utils.GetNowUTC8Time(), i))
// 				if result.Rows == 0 && result.Code == "" {
// 					beego.Emergency("------------------------- 严重问题 : ", lotteryInfo.GameTag, " Api获取历史记录失败 将彩票状态设置为 2 维护! 请立即检查 -------------------------")
// 					//通过控制中心改变彩票状态
// 					ctrl.Instance().ChangeLotteryStatus(lotteryInfo.GameTag, gb.LotteryStatus_Maintain)
// 					return
// 				} else {
// 					//如果数据正确存库
// 					dbmgr.InsertLotteryHistoryByDay(o.FilterRecord(result, lastRecord.Expect))
// 					//beego.Info("-------------------- 数据正确 存入数据库 ! -------------------------")
// 				}
// 			}
// 		}
// 	} else if lotteryInfo.Frequency == "low" {
// 		//彩种是低频彩的情况,由于低频菜api按天查询只会给出从今年1月开始的记录所以没有办法保证可以有100期的历史记录
// 		result := LotteryApiMgr.Instance().GetRecordByDate(lotteryInfo.GameTag, Utils.GetNowUTC8Time())
// 		if result.Rows == 0 {
// 			beego.Emergency("------------------------- 严重问题 : ", lotteryInfo.GameTag, " Api获取历史记录失败 将彩票状态设置为 2 维护! 请立即检查 -------------------------")
// 			//通过控制中心改变彩票状态
// 			ctrl.Instance().ChangeLotteryStatus(lotteryInfo.GameTag, gb.LotteryStatus_Maintain)
// 			return
// 		} else {
// 			//如果数据正确存库
// 			dbmgr.InsertLotteryHistoryByDay(o.FilterRecord(result, lastRecord.Expect))
// 			//beego.Info("-------------------- 数据正确 存入数据库 ! -------------------------")
// 		}
// 	}
// }

// //将从api按天获取的历史记录去掉数据库中已经存在的(筛选记录) 参数1 为api获取的历史记录, 参数2 为数据库中查出的最后期数
// func (o *HistoryMgr) FilterRecord(data gb.LotteryRecordByDayFromApi, lastExpect int) gb.LotteryRecordByDayFromApi {
// 	var result []gb.LotteryRecordFromApi

// 	for _, v := range data.Data {

// 		expect, _ := strconv.Atoi(v.Expect)
// 		//过滤数据只有大于数据库最后一期的数据才储存
// 		if expect > lastExpect {
// 			result = append(result, v)
// 		}
// 	}
// 	data.Data = result
// 	return data
// }
