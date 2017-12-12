package histmgr

import (
	"apisrv/conf"
	"apisrv/models/apimgr"
	"apisrv/models/ctrl"
	"apisrv/models/dbmgr"
	"common/httpmgr"
	"common/utils"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"sync"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type histmgr struct {
}

var sInstance *histmgr
var once sync.Once

func Instance() *histmgr {
	once.Do(func() {
		sInstance = &histmgr{}
		sInstance.init()
	})
	return sInstance
}

func (o *histmgr) init() {
	beego.Info("--- Sttart Complement The Historical Records :")

	for _, v := range ctrl.Instance().Ltrys {

		beego.Info("---", v.Game_name)

		o.ComplementRecord(v.Game_name)
	}

	beego.Info("--- All Lottery Record Complement Completion !")
}

func (o *histmgr) ComplementRecord(gameName string) {

	ltry, ok := ctrl.Instance().Ltrys[gameName]
	if !ok {
		beego.Error("Not have this lottery : ", gameName)
		return
	}

	lastHistRecord := o.GetLtryLastRecord(ltry.Game_name)

	if ltry.Freq == 1 { //1 为高频彩

		now := utils.GetNowUTC8Time()

		recordTime := utils.ConvertToUTC8Time(lastHistRecord.Open_time)

		if utils.DateSub(utils.GetNowUTC8Time(), utils.ConvertToUTC8Time(lastHistRecord.Open_time)) == 0 {
			//数据库中查询到历史记录,如果最后一次记录是同一天, 查看期数 补全
			result, err := o.GetLtryRecordByDate(ltry.Game_name, utils.GetNowUTC8Time())
			if err != nil {
				beego.Error(err)
				return
			}
			//由于服务器启动时间有可能今天并没有记录 比如新疆时时彩 要上午10点才开奖,9点启动服务器就会找不到今天的记录,所以rows == 0 并不代表数据错误,加上code 为空的判断才能说明数据是否正确(还是有点累赘后面来改)
			//这个明早来测试,应为新的获取记录函数里面就已经解析了json格式,不知道再没有数据的时候解析json数据会不会出错
			//按照原来的代码解析会成功 但是没有数据,一定记得明早来测试
			if result.Rows == 0 && result.Code == "" {
				beego.Error("--- 严重问题 : ", ltry.Game_name, " Api获取历史记录失败 ")
				//通过控制中心改变彩票状态
				return
			} else {
				o.insertLtryHistByDay(o.FilterRecord(result, lastHistRecord.Expect))
			}
		} else {
			var d int
			if lastHistRecord.Expect == 0 || utils.DateSub(now, recordTime) >= 4 {
				d = 3
			} else {
				d = utils.DateSub(utils.GetNowUTC8Time(), utils.ConvertToUTC8Time(lastHistRecord.Open_time))
			}
			for i := d; i >= 0; i-- {
				//获取结果往日结果
				result, err := o.GetLtryRecordByDate(ltry.Game_name, utils.DateBeforeTheDay(utils.GetNowUTC8Time(), i))
				if err != nil {
					beego.Error(err)
					return
				}
				//等完成功能来加 获取失败以后将菜种设置为维护
				o.insertLtryHistByDay(result)
			}
		}
	} else if ltry.Freq == 0 { //0 为低频彩
		//彩种是低频彩的情况,由于低频菜api按天查询只会给出从今年1月开始的记录所以没有办法保证可以有100期的历史记录
		result, err := o.GetLtryRecordByDate(ltry.Game_name, utils.GetNowUTC8Time())
		if err != nil {
			beego.Error(err)
			return
		}
		if result.Rows == 0 {
			beego.Error("--- 严重问题 : ", ltry.Game_name, " Api获取历史记录失败 !")
			//通过控制中心改变彩票状态
			return
		} else {
			o.insertLtryHistByDay(o.FilterRecord(result, lastHistRecord.Expect))
		}
	}
}

//从数据库中获取一个彩种最后的开采记录
func (o histmgr) GetLtryLastRecord(gameName string) LtryHistRecord {
	ret := LtryHistRecord{}

	bsonM := bson.M{"game_name": gameName}
	err := dbmgr.Instance().HistColl.Find(bsonM).Sort("-expect").One(&ret)
	if err != nil {
		beego.Debug(err, "Lottery name : ", gameName)
	}

	return ret
}

//按日期获取一个彩种当日所有开采(这里有个问题,这个函数应该是提一个统一调用接口,返回的应该是统一的数据格式,而不管这个数据是从开采网还是其他什么api提供商来得,所以现在先统一格式,等以后加上其他api以后再将其他API获取的数据都以这种统一的格式返回)
//不同的彩种访问间隔时间是不一样的目前开采网按天访问是需要间隔5秒,访问带下一期的是 1秒间隔,函数负责返回api返回的数据结构
func (o *histmgr) GetLtryRecordByDate(gameName string, date time.Time) (LtryRecordDay, error) {
	var ret LtryRecordDay

	d := date.Format(utils.TF_D)

	//循环每一个api提供商 目前优先使用开彩票网
	for i := 0; i < len(apimgr.Instance().ApiDay); i++ {
		//beego.Debug("--- 开始循环 api 提供商 id :", i)
		apiMap := apimgr.Instance().ApiDay[i]

		if apiArray, ok := apiMap[gameName]; ok {
			//循环每一个不同的接口(目前优先使用开彩网的高防主接口,然后是副接口,其次是他妈的 坑爹的 V ! I ! P ! 接口)
			for j := len(apiArray) - 1; j >= 0; j-- {
				count := conf.GetRecordOneDayRetryCount
				for ; count > 0; count-- {
					url := apiArray[j] + d
					resp, err := httpmgr.Get(url)
					if err == nil {
						err := json.Unmarshal(resp, &ret)
						if err == nil {
							return ret, nil
						}
					}
					time.Sleep(conf.GetRecordByNewestSleepTime * time.Second)
				}
			}
		} else {
			return ret, errors.New("There is no API for this lottery : " + gameName)
		}
		beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", gameName)
	}

	return ret, errors.New("--- Failed to obtain data from all API providers !!! Game Name : " + gameName)
}

//插入一天的开奖记录.这里面要转换数据
func (o histmgr) insertLtryHistByDay(data LtryRecordDay) {
	//由于从开采网获取的历史记录 是从晚到早,所以存库时要反过来
	for i := len(data.Data) - 1; i >= 0; i-- {
		ltryRecord := &LtryHistRecord{}
		//游戏标志
		ltryRecord.Game_name = utils.ConvertGameName(data.Code)

		//开奖期数
		var err error
		ltryRecord.Expect, err = strconv.Atoi(data.Data[i].Expect)
		if err != nil {
			beego.Error("--- Insert Lottery History Error ! ")
			return
		}

		//开奖号码
		ltryRecord.Open_code = data.Data[i].Opencode
		//ltryRecord.Open_code = strings.Replace(ltryRecord.Open_code, "+", ",", -1)

		//开奖时间 注意加上时区 prc 北京时间;
		loc, _ := time.LoadLocation("PRC")
		utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.Data[i].OpenTime, loc)
		ltryRecord.Open_time = utcTime
		if err != nil {
			beego.Error("--- Insert Lottery History Error ! ")
			return
		}

		//开奖时间戳
		ltryRecord.Open_time_stamp = data.Data[i].OpenTimeStamp

		if err != nil {
			beego.Error("--- Insert Lottery History Error ! ")
			return
		}

		ltryRecord.Recording_time = utils.GetNowUTC8Time()
		//记录时间,入库时间(时间戳);
		ltryRecord.Recording_time_stamp = time.Now().Unix()
		//插入数据库
		dbmgr.Instance().HistColl.Insert(ltryRecord)
	}
}

//过滤记录 将从api按天获取的历史记录去掉数据库中已经存在的(筛选记录) 参数1 为api获取的历史记录, 参数2 为数据库中查出的最后期数
func (o *histmgr) FilterRecord(data LtryRecordDay, lastExpect int) LtryRecordDay {
	var result []LtryRecord

	for _, v := range data.Data {
		expect, _ := strconv.Atoi(v.Expect)
		//过滤数据只有大于数据库最后一期的数据才储存
		if expect > lastExpect {
			result = append(result, v)
		}
	}
	data.Data = result
	return data
}

//这个是数据库中保存的彩票历史记录结构
type LtryHistRecord struct {
	Game_name            string    // 名称
	Expect               int       // 期次
	Open_code            string    // 开奖号码
	Open_time            time.Time // 开奖时间
	Open_time_stamp      int64     //开奖时间(时间戳)
	Recording_time       time.Time //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	Recording_time_stamp int64     //记录时间(时间戳)
}

//获取的按天查询历史记录(同样这个也是统一格式目前使用的是开采网中格式,等以后加上了其他的api提供商以后,都要将数据全部统一成这个格式)
type LtryRecordDay struct {
	Rows   int
	Code   string
	Remain string
	Data   []LtryRecord
}

//一条记录开彩记录(这个是从api获取的一条彩票开采记录,目前这个结构式开采网的,以后加入了新的api提供商以后,就都需要把记录修改为这种格式)
type LtryRecord struct {
	Expect        string
	Opencode      string
	OpenTime      string
	OpenTimeStamp int64
}
