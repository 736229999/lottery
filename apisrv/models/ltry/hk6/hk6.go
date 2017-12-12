package hk6

import (
	"apisrv/models/apimgr"
	"apisrv/models/dbmgr"
	"apisrv/models/histmgr"
	"common/utils"
	"strconv"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

type HK6 struct {
	id        int    //彩票ID
	gameName  string //游戏名字
	frequency int    //彩票频率

	currentExpect   int       //当前这期期数(最近已开奖期数)
	openCode        []int     //当前这期开奖号码(最近已开奖号码)
	openCodeStr     string    //当前这期开奖号码String形式(最近已开奖号码)
	currentOpenTime time.Time //当期开奖时间(最近已开出期数的时间)
	nextExpect      int       //下期开奖期数
	nextOpenTime    time.Time //下期开奖时间

	status int //彩票状态 0.关闭, 1正常, 2维护
	//历史记录

	lock sync.RWMutex //为了避免可能出现的同步读写,出现结构体数据更新一半的情况,尝试加入通用锁
}

func Init(id int, gameName string, freq int) *HK6 {
	o := &HK6{}

	o.id = id
	o.gameName = gameName
	o.frequency = freq
	o.status = 1

	o.UpdataInfo()
	return o
}

//初始化手动彩票 返回错误码,
func (o *HK6) InitLtry(currentExpect int, nextExpect int, nextOpenTime string) int {
	if o.nextExpect != 0 {
		beego.Warn("HK6 has been initialized !")
		return 99
	}
	o.lock.Lock()
	defer o.lock.Unlock()

	//开始验证各项传输数据是否正确
	ret := histmgr.LtryHistRecord{}
	err := dbmgr.Instance().HistColl.Find(bson.M{"game_name": o.gameName}).Sort("-expect").One(&ret)
	if err != nil {
		beego.Error(err)
		return 100
	}

	//1.验证当前期数 和 数据库中最新的期数是否一致
	if ret.Expect != currentExpect {
		beego.Error("Initialize lottery error : The current expect is incorrect !")
		beego.Error("DB expect : ", ret.Expect)
		return 1
	}

	//3.验证开奖时间日期格式是否正确
	loc, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Error(err)
		return 2
	}

	//5.验证下期开奖时间格式是否合法
	t, err := time.ParseInLocation("2006-01-02 15:04:05", nextOpenTime, loc)
	if err != nil {
		beego.Error("Initialize lottery error : The next open time is incorrect !")
		beego.Error(err)
		return 3
	}

	if o.currentOpenTime.After(t) {
		beego.Error("Initialize lottery error : The next open time is incorrect !")
		return 4
	}

	o.nextOpenTime = t

	//6.验证下期期数是否正确 六合彩的期数是一直累加 但是过年的时候就会有变化
	if nextExpect != o.CalNextExpect(o.nextOpenTime, o.currentExpect) {
		beego.Error("Initialize lottery error : The next expect is incorrect !")
		return 5
	}
	o.nextExpect = nextExpect

	beego.Info("--- Initialize HK6 Success !")
	return 0
}

//手动开奖 返回错误码
func (o *HK6) StartLtry(expect int, openCode string, openTime string, nextExpect int, nextOpenTime string) int {
	if o.nextExpect == 0 {
		beego.Warn(o.gameName, " is Not initialized !")
		return 1
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	//1.验证当前期数是否正确
	if expect != o.nextExpect {
		beego.Error("Start lottery error : The expect is incorrect !")
		return 2
	}

	//2.验证开奖号码是否合法
	oc := utils.PaserOpenCodeToArray(openCode)
	if len(oc) != 7 {
		beego.Error("Start lottery error : The open code count is incorrect !")
		return 3
	}

	//3.验证开奖号码每位号码是否合法 六合彩的开奖号码 1 - 49
	for _, v := range oc {
		if v < 1 || v > 49 {
			beego.Error("Start lottery error : The open code Numeric range is incorrect !")
			return 4
		}
	}

	//4.验证开奖时间格式是否正确
	loc, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Error(err)
		return 5
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", openTime, loc)
	if err != nil {
		beego.Error("Start lottery error : The open time format is incorrect !")
		beego.Error(err)
		return 6
	}

	if o.currentOpenTime.After(t) {
		beego.Error("Start lottery error : The open time is incorrect !")
		return 7
	}

	//5.验证下期开奖时间格式是否正确
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", nextOpenTime, loc)
	if err != nil {
		beego.Error("Start lottery error : The next open time format is incorrect !")
		beego.Error(err)
		return 8
	}

	if t.After(t2) {
		beego.Error("Start lottery error : The next open time is incorrect !")
		return 9
	}

	//6.验证下期期数是否正确
	if nextExpect != o.CalNextExpect(t2, expect) {
		beego.Error("Start lottery error : The next expect  is incorrect ! : ", nextExpect)
		return 10
	}

	//验证全部通过更新数据
	o.currentExpect = expect
	o.openCode = oc
	o.openCodeStr = openCode
	o.currentOpenTime = t
	o.nextExpect = nextExpect
	o.nextOpenTime = t2

	histRecord := &histmgr.LtryHistRecord{}
	histRecord.Expect = o.currentExpect
	histRecord.Game_name = o.gameName
	histRecord.Open_code = o.openCodeStr
	histRecord.Open_time = o.currentOpenTime
	histRecord.Open_time_stamp = histRecord.Open_time.Unix()
	histRecord.Recording_time = time.Now()
	histRecord.Recording_time_stamp = histRecord.Recording_time.Unix()

	//将最新记录插入数据库
	err = dbmgr.Instance().HistColl.Insert(histRecord)
	if err != nil {
		beego.Error(err)
		return 11
	}

	beego.Info("--- Start HK6 Success !")

	return 0
}

func (o *HK6) UpdataInfo() {
	//验证API返回的结果是否正确
	newRecord, err := apimgr.Instance().GetNewRecord(o.gameName)

	if err != nil {
		beego.Error("Lottery New Record Error : ", err)
		o.status = 2
		//这里后面来加,如果通过api 获取记录失败,这里要改变通知控制服务器
		return
	}

	if !o.verifyRecord(newRecord) {
		beego.Error("Lottery New Record Error : ", newRecord)
		o.status = 2
		return
	}

	o.currentExpect, _ = strconv.Atoi(newRecord.Open[0].Expect)
	o.openCode = utils.PaserOpenCodeToArray(newRecord.Open[0].Opencode)
	o.openCodeStr = newRecord.Open[0].Opencode

	loc, _ := time.LoadLocation("PRC")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", newRecord.Open[0].Opentime, loc)
	o.currentOpenTime = t

	o.status = 1
}

//验证API给的数据是否正确(由于现在不知道api给的数据会出现什么错误,所以只能进行基本的正确性判断)
//以后遇到问题就在这里面添加错误判断
func (o HK6) verifyRecord(newRecord apimgr.LtryRecordNew) bool {
	if len(newRecord.Open) < 1 {
		return false
	}

	if newRecord.Rows != 1 {
		return false
	}

	return true
}

//得到游戏名字
func (o HK6) GetGameName() string {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.gameName
}

//得到当前期数
func (o HK6) GetCurrentExpect() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.currentExpect
}

//得到开奖号码数组
func (o HK6) GetOpenCode() []int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.openCode
}

//得到开奖号码 string
func (o HK6) GetOpenCodeStr() string {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.openCodeStr
}

//得到当前期开奖时间
func (o HK6) GetCurrentOpenTime() time.Time {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.currentOpenTime
}

//得到下期期数
func (o HK6) GetNextExpect() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.nextExpect
}

//得到下期开奖时间
func (o HK6) GetNextOpenTime() time.Time {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.nextOpenTime
}

//游戏频率(高频还是低频)
func (o HK6) GetFreq() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.frequency
}

func (o HK6) GetStatus() int {
	return o.status
}

//得到一个彩票最新的开彩记录
// func (o HK6) GetNewRecord() (LtryRecordNew, error) {
// 	var ret LtryRecordNew

// 	for i := 0; i < len(apimgr.Instance().ApiNew); i++ {
// 		apiMap := apimgr.Instance().ApiNew[i]

// 		if apiArray, ok := apiMap[o.gameName]; ok {
// 			//循环每一个不同的接口(获取最新记录和获取历史记录不一样,优先使用0号高防接口)
// 			for j := 0; j < len(apiArray); j++ {
// 				count := 3 //重试次数,这里暂时写死 以后这些全部要写到控制服数据库中去
// 				for ; count > 0; count-- {
// 					url := apiArray[j]
// 					resp, err := httpmgr.Get(url)
// 					if err == nil {
// 						err := json.Unmarshal(resp, &ret)
// 						if err == nil {
// 							return ret, nil
// 						}
// 					}
// 					time.Sleep(conf.GetRecordByNewestSleepTime * time.Second)
// 				}
// 			}
// 		} else {
// 			return ret, errors.New("There is no API for this lottery : " + o.gameName)
// 		}
// 		beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", o.gameName)
// 	}

// 	return ret, errors.New("--- Failed to obtain data from all API providers !!! Game Name : " + o.gameName)
// }

//计算下期期数(参数1 下期开奖时间, 参数2 当前期数)
func (o HK6) CalNextExpect(nextOpenTime time.Time, currentExpect int) int {
	//查看下期开奖时间和今天是不是跨年了
	if nextOpenTime.Year() > time.Now().Year() {
		beego.Warn("------------- 跨年了!!!!!! 新年快乐!!!!!!! --------------")
		return nextOpenTime.Year()*100 + 1
	}
	currentExpect++
	return currentExpect
}

//开采网 按最新获取一条记录(带下期)
// type LtryRecordNew struct {
// 	Rows   int
// 	Code   string
// 	Remain string
// 	Next   []LtryRecordNext    //下一期信息(下一期期数, 下一期开奖时间)
// 	Open   []LtryRecordCurrent //最新一期信息 (注意:这里面没有时间戳,存库的时候要自己添加一个)
// 	Time   string              //查询时间
// }

// //开采网 下一期开采信息
// type LtryRecordNext struct {
// 	Expect   string
// 	Opentime string
// }

// //开采网 当前这期信息
// type LtryRecordCurrent struct {
// 	Expect   string
// 	Opencode string
// 	Opentime string
// }
