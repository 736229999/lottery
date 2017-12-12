package ltrymgr

import (
	"apisrv/models/apimgr"
	"apisrv/models/dbmgr"
	"apisrv/models/histmgr"
	"strings"
	"sync"

	"common/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type Ltry struct {
	id        int    //彩票ID
	gameName  string //游戏名字
	frequency int    //彩票频率

	currentExpect   int       //当前这期期数(最近已开奖期数)
	openCode        []int     //当前这期开奖号码(最近已开奖号码)
	openCodeString  string    //当前这期开奖号码String形式(最近已开奖号码)
	currentOpenTime time.Time //当期开奖时间(最近已开出期数的时间)

	nextExpect   int       //下期期数
	nextOpenTime time.Time //下期开彩时间

	status int //彩票状态 0.关闭, 1正常, 2维护

	//等完成了功能来这里将历史记录放入类中

	lock sync.RWMutex //为了避免可能出现的同步读写,出现结构体数据更新一半的情况,尝试加入通用锁
}

func Init(id int, gameName string, freq int) *Ltry {
	o := &Ltry{}

	o.id = id
	o.gameName = gameName
	o.frequency = freq
	o.status = 1

	newRecord, err := apimgr.Instance().GetNewRecord(o.gameName)
	if err != nil {
		beego.Error("Lottery New Record Error : ", err)
		o.status = 2
		//这里后面来加,如果通过api 获取记录失败,这里要改变通知控制服务器
		return o
	}

	o.UpdataInfo(newRecord)
	return o
}

func (o *Ltry) UpdataInfo(newRecord apimgr.LtryRecordNew) {

	//验证API返回的结果是否正确

	if !o.verifyRecord(newRecord) {
		beego.Error("Lottery New Record Error : ", newRecord)
		o.status = 2
		return
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	o.currentExpect, _ = strconv.Atoi(newRecord.Open[0].Expect)
	tmp := strings.Replace(newRecord.Open[0].Opencode, "+", ",", -1)
	o.openCode = utils.PaserOpenCodeToArray(tmp)
	o.openCodeString = newRecord.Open[0].Opencode

	loc, _ := time.LoadLocation("PRC")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", newRecord.Open[0].Opentime, loc)
	o.currentOpenTime = t

	o.nextExpect, _ = strconv.Atoi(newRecord.Next[0].Expect)
	o.nextOpenTime, _ = time.ParseInLocation("2006-01-02 15:04:05", newRecord.Next[0].Opentime, loc)

	o.status = 1
}

//验证API给的数据是否正确(由于现在不知道api给的数据会出现什么错误,所以只能进行基本的正确性判断)
//以后遇到问题就在这里面添加错误判断
func (o Ltry) verifyRecord(newRecord apimgr.LtryRecordNew) bool {
	if len(newRecord.Open) < 1 || len(newRecord.Next) < 1 {
		return false
	}

	if newRecord.Rows != 1 {
		return false
	}

	return true
}

//得到游戏名字
func (o Ltry) GetGameName() string {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.gameName
}

//游戏频率(高频还是低频)
func (o Ltry) GetFreq() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.frequency
}

func (o Ltry) GetStatus() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.status
}

func (o Ltry) GetCurrentExpect() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.currentExpect
}

func (o Ltry) GetOpenCode() []int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.openCode
}

func (o Ltry) GetOpenCodeStr() string {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.openCodeString
}

func (o Ltry) GetCurrentOpenTime() time.Time {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.currentOpenTime
}

func (o Ltry) GetNextExpect() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.nextExpect
}

func (o Ltry) GetNextOpenTime() time.Time {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.nextOpenTime
}

func (o *Ltry) StartLtry(ltryRecord apimgr.LtryRecordNew) {
	//检查是否掉期
	o.CheckLostExpect(ltryRecord)

	o.UpdataInfo(ltryRecord)

	//将最新记录插入数据库
	o.InsterHistRecordByNew(ltryRecord)

	beego.Info("--- 开采啦 ! : ", o.gameName, "  -- expect : ", ltryRecord.Open[0].Expect, "  -- open code : ", o.openCode)
	return
}

//使用到这个函数的情况是,向API 获取数据成功,但是掉期,在现在的结构中,如果api访问有问题,会去其他的备用接口访问
func (o *Ltry) CheckLostExpect(ltryRecord apimgr.LtryRecordNew) {
	newestExpect, _ := strconv.Atoi(ltryRecord.Open[0].Expect)
	//如果 最新获取的期数大于下棋期数,开始进入补全程序
	if newestExpect > o.nextExpect {
		beego.Error("--- 发现掉期 开始进补全 ! \n")
		beego.Error("--- ", o.gameName, " 最新期数为 : ", newestExpect, " 当前下期期数为 : ", o.nextExpect, "\n")

		//记录掉得期数
		var lostExpect []int
		//计算要补的期数(计算期数差),并得到要补的期数切片
		for {
			//得到最新获得期数的上一期期数
			lastExpect, _ := utils.GetLastExpect(o.gameName, newestExpect)
			//如果最新一期的上一期不等于下期期数,证明掉期不止一期,将这期期数保存下来
			if lastExpect != o.nextExpect {
				lostExpect = append(lostExpect, lastExpect)
				newestExpect = lastExpect
			} else {
				//只掉一期,记录下来,然后跳出for循环
				lostExpect = append(lostExpect, lastExpect)
				break
			}
		}

		//开始补期
		//获取今天的历史记录
		expectRecord, err := apimgr.Instance().GetLtryRecordByDate(o.gameName, time.Now())
		if err != nil {
			beego.Error(err)
			return
		}
		//查询掉的期数是否在记录中
		var lostExpectInfo []apimgr.LtryRecord

		//查看掉得期数有没有在当天的记录里
		for _, v := range lostExpect {
			for _, i := range expectRecord.Data {
				tmpExpect, _ := strconv.Atoi(i.Expect)
				if tmpExpect == v { //招到了掉的期信息,保存掉的信息
					lostExpectInfo = append(lostExpectInfo, i)
				}
			}
		}

		//查看有没有找到掉期的信息, 如果在上面当天的记录里面没有找到,那么就再去前一天的日期里去查找
		if len(lostExpect) != len(lostExpectInfo) {
			expectRecord2, err := apimgr.Instance().GetLtryRecordByDate(o.gameName, utils.DateBeforeTheDay(time.Now(), 1))
			if err != nil {
				beego.Error(err)
				return
			}

			for _, v := range lostExpect {
				for _, i := range expectRecord2.Data {
					tmpExpect, _ := strconv.Atoi(i.Expect)
					if tmpExpect == v {
						lostExpectInfo = append(lostExpectInfo, i)
					}
				}
			}
		}

		//经过两天查找 还是没有找到数据的话,报出错误,并补全找到的数据
		if len(lostExpect) != len(lostExpectInfo) {
			beego.Error(o.gameName, " 掉期补全失败 请手动检查补全 !! \n")
		}

		//将补全的数据插入数据库以便计算服发现掉期来请求结果
		for _, v := range lostExpectInfo {
			o.InsterHistRecordByApi(v)
		}

		beego.Error("--- 掉期补全完成 !")
	}
}

//向数据库插入 最新的一条记录历史记录 参数 LtryRecordNew
func (o Ltry) InsterHistRecordByNew(ltryRecord apimgr.LtryRecordNew) {
	histRecord := &histmgr.LtryHistRecord{}
	//游戏标志
	histRecord.Game_name = utils.ConvertGameName(ltryRecord.Code)

	//开奖期数
	var err error
	histRecord.Expect, err = strconv.Atoi(ltryRecord.Open[0].Expect)
	if err != nil {
		beego.Error("--- Insert Lottery History Error ! ")
		return
	}

	//开奖号码
	histRecord.Open_code = ltryRecord.Open[0].Opencode

	//开奖时间 注意加上时区 prc 北京时间;
	loc, _ := time.LoadLocation("PRC")
	utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", ltryRecord.Open[0].Opentime, loc)
	histRecord.Open_time = utcTime
	if err != nil {
		beego.Error("--- Insert Lottery History Error ! ")
		return
	}

	//开奖时间戳
	histRecord.Open_time_stamp = histRecord.Open_time.Unix()

	if err != nil {
		beego.Error("--- Insert Lottery History Error ! ")
		return
	}

	histRecord.Recording_time = utils.GetNowUTC8Time()
	//记录时间,入库时间(时间戳);
	histRecord.Recording_time_stamp = time.Now().Unix()

	//插入数据库
	err = dbmgr.Instance().HistColl.Insert(histRecord)
	if err != nil {
		beego.Debug(err)
		return
	}
}

//向数据库插入 最新的一条记录历史记录 参数 apimgr.LtryRecord
func (o Ltry) InsterHistRecordByApi(ltryRecord apimgr.LtryRecord) {
	histRecord := &histmgr.LtryHistRecord{}
	//游戏标志
	histRecord.Game_name = o.gameName
	//开奖期数
	var err error
	histRecord.Expect, err = strconv.Atoi(ltryRecord.Expect)
	if err != nil {
		beego.Error("--- Insert Lottery History Error ! ")
		return
	}

	//开奖号码
	histRecord.Open_code = ltryRecord.Opencode

	//开奖时间 注意加上时区 prc 北京时间;
	loc, _ := time.LoadLocation("PRC")
	utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", ltryRecord.OpenTime, loc)
	histRecord.Open_time = utcTime
	if err != nil {
		beego.Error("--- Insert Lottery History Error ! ")
		return
	}

	//开奖时间戳
	histRecord.Open_time_stamp = histRecord.Open_time_stamp

	if err != nil {
		beego.Error("--- Insert Lottery History Error ! ")
		return
	}

	histRecord.Recording_time = utils.GetNowUTC8Time()
	//记录时间,入库时间(时间戳);
	histRecord.Recording_time_stamp = time.Now().Unix()

	//插入数据库
	err = dbmgr.Instance().HistColl.Insert(histRecord)
	if err != nil {
		beego.Debug(err)
		return
	}
}
