package pk10

import (
	"common/utils"
	"indsrv/models/dbmgr"
	"indsrv/models/rd"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//注意 由于急速fpk10是一分钟一期,那么原6位数的期数只能支撑不到两年时间,所以急速fpk10期数为7位
type PK10 struct {
	GameName        string
	CurrentExpect   int
	OpenCode        []int
	OpenCodeStr     string
	CurrentOpenTime time.Time

	NextExpect   int
	NextOpenTime time.Time

	Interval int //开奖间隔 为分钟
}

//所有的这些设置等项目完成之后 所有的参数都要改到我自己的服务器后台去
func Init() (*PK10, error) {
	ltry := &PK10{}
	ltry.GameName = "PK10_F"
	ltry.Interval = 1

	record := &Record{}
	err := dbmgr.PK10Coll.Find(bson.M{"game_name": ltry.GameName}).Sort("-expect").One(&record)

	if err != nil {
		if err.Error() == "not found" {
			beego.Error("--- No records found, add a record !")
			record.Game_name = ltry.GameName
			record.Expect = 1888888
			record.Open_code, _ = ltry.RandomOpenCode()
			record.Open_time = utils.GetIntegerMin(time.Now())
			record.Open_time_stamp = record.Open_time.UnixNano()
			record.Recording_time = time.Now()
			record.Recording_time_stamp = record.Recording_time.UnixNano()

			dbmgr.PK10Coll.Insert(record)
		} else {
			beego.Error("--- Init PK10 failed !")
			return nil, err
		}
	}

	ltry.CurrentExpect = record.Expect
	ltry.OpenCode = utils.PaserOpenCodeToArray(record.Open_code)
	ltry.OpenCodeStr = record.Open_code
	ltry.CurrentOpenTime = record.Open_time
	ltry.NextExpect = ltry.CurrentExpect + 1
	ltry.NextOpenTime = utils.TimeAfterMin(ltry.CurrentOpenTime, ltry.Interval)

	//itv = interval
	itv := time.Now().Sub(ltry.CurrentOpenTime)

	itvi := int(itv.Minutes())

	if itvi >= ltry.Interval {
		ltry.ComplementExpect(itvi / ltry.Interval)
	}

	go ltry.TimerLtry()

	beego.Info("--- Init FPK10 Done !")
	return ltry, nil
}

func (o *PK10) ComplementExpect(itv int) {
	beego.Info("--- PK10 补全 : ", itv, "期 !")

	record := &Record{}
	record.Game_name = o.GameName

	for i := 0; i < itv; i++ {
		o.CurrentExpect++
		record.Expect = o.CurrentExpect
		record.Open_code, _ = o.RandomOpenCode()
		o.CurrentOpenTime = utils.TimeAfterMin(o.CurrentOpenTime, o.Interval)
		record.Open_time = o.CurrentOpenTime
		record.Open_time_stamp = record.Open_time.UnixNano()
		record.Recording_time = time.Now()
		record.Recording_time_stamp = record.Recording_time.UnixNano()

		dbmgr.PK10Coll.Insert(record)
	}

	o.CurrentExpect = record.Expect
	o.OpenCode = utils.PaserOpenCodeToArray(record.Open_code)
	o.OpenCodeStr = record.Open_code
	o.CurrentOpenTime = record.Open_time
	o.NextExpect = o.CurrentExpect + 1
	o.NextOpenTime = utils.TimeAfterMin(o.CurrentOpenTime, o.Interval)
}

//数学知识不足,暂时现以这种排重的方式完成,后期补完数学知识后再来重写
func (o *PK10) RandomOpenCode() (string, []int) {
	var oc []int

	for len(oc) < 10 {
		r := rd.Intn(1, 10)
		exist := false

		for _, v := range oc {
			if v == r {
				exist = true
				break
			}
		}

		if !exist {
			oc = append(oc, r)
		}
	}

	var ocs string
	ocs = strconv.Itoa(oc[0])
	for i := 1; i < 10; i++ {
		ocs = ocs + "," + strconv.Itoa(oc[i])
	}
	return ocs, oc
}

func (o *PK10) TimerLtry() {
	beego.Info("--- PK10 开始定时开采 !")

	for {
		if time.Now().After(o.NextOpenTime) {

			record := &Record{}
			record.Game_name = o.GameName
			o.CurrentExpect++
			record.Expect = o.CurrentExpect
			record.Open_code, _ = o.RandomOpenCode()
			record.Open_time = utils.GetIntegerMin(time.Now())
			record.Open_time_stamp = record.Open_time.UnixNano()
			record.Recording_time = time.Now()
			record.Recording_time_stamp = record.Recording_time.UnixNano()

			o.CurrentExpect = record.Expect
			o.OpenCode = utils.PaserOpenCodeToArray(record.Open_code)
			o.OpenCodeStr = record.Open_code
			o.CurrentOpenTime = record.Open_time
			o.NextExpect = o.CurrentExpect + 1
			o.NextOpenTime = utils.TimeAfterMin(o.CurrentOpenTime, o.Interval)

			dbmgr.PK10Coll.Insert(record)

			beego.Info("--- 开采啦 ! Game Name : ", o.GameName, "  Expect : ", o.CurrentExpect, " OpenCode : ", o.OpenCodeStr)

			time.Sleep(o.NextOpenTime.Sub(time.Now()))
		}
		d, _ := time.ParseDuration("200ms")
		time.Sleep(d)
	}
}

//------------------------- Interface ---------------------
func (o *PK10) GetGameName() string {
	return o.GameName
}

func (o *PK10) GetCurrentExpect() int {
	return o.CurrentExpect
}

func (o *PK10) GetOpenCode() []int {
	return o.OpenCode
}

func (o *PK10) GetOpenCodeStr() string {
	return o.OpenCodeStr
}

func (o *PK10) GetCurrentOpenTime() time.Time {
	return o.CurrentOpenTime
}

func (o *PK10) GetNextExpect() int {
	return o.NextExpect
}

func (o *PK10) GetNextOpenTime() time.Time {
	return o.NextOpenTime
}

//这个是数据库中保存的彩票历史记录结构
type Record struct {
	Game_name            string    // 名称
	Expect               int       // 期次
	Open_code            string    // 开奖号码
	Open_time            time.Time // 开奖时间
	Open_time_stamp      int64     //开奖时间(时间戳)
	Recording_time       time.Time //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	Recording_time_stamp int64     //记录时间(时间戳)
}
