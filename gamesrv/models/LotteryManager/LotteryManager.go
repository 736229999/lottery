package LotteryManager

import (
	"errors"
	"gamesrv/models/Utils"
	"sync"

	"github.com/astaxie/beego"
)

var sInstance *LotteryManager
var once sync.Once

type LotteryManager struct {
	lotteriesInfoMap map[string]LotteryInfo //所有彩票类型信息
	infoMapLock      sync.RWMutex

	lotteriesRecordMap map[string][]LotteryRecord //所有彩票历史记录用于保存，计算服每次更新消息就是这个
	recordMapLock      sync.RWMutex
}

/*
初始化实例
*/
func Instance() *LotteryManager {
	once.Do(func() {
		sInstance = &LotteryManager{}
		sInstance.initManager()
	})
	return sInstance
}

/*
初始化操作
*/
func (o *LotteryManager) initManager() {
	o.lotteriesInfoMap = make(map[string]LotteryInfo)
	o.lotteriesRecordMap = make(map[string][]LotteryRecord)
}

//得到所有彩票信息map
func (o *LotteryManager) GetLtryInfoMap() map[string]LotteryInfo {
	return o.lotteriesInfoMap
}

//读取彩票信息map数据(读锁)
func (o *LotteryManager) GetLtryInfo(k string) (LotteryInfo, error) {
	o.infoMapLock.RLock()
	defer o.infoMapLock.RUnlock()
	if ret, ok := o.lotteriesInfoMap[k]; ok {
		return ret, nil
	} else {
		b := LotteryInfo{}
		str := "Not have this lottery type : " + k
		return b, errors.New(str)
	}
}

//设置彩票信息map数据(通用锁)
func (o *LotteryManager) SetLtryInfo(k string, v LotteryInfo) {
	o.infoMapLock.Lock()
	defer o.infoMapLock.Unlock()
	o.lotteriesInfoMap[k] = v
}

//读取彩票历史记录map(读锁)
func (o *LotteryManager) GetLtryRecord(k string) ([]LotteryRecord, error) {
	o.recordMapLock.RLock()
	defer o.recordMapLock.RUnlock()

	if ret, ok := o.lotteriesRecordMap[k]; ok {
		return ret, nil
	} else {
		str := "Not have this lottery history record : " + k
		return nil, errors.New(str)
	}
}

//设置彩票历史记录map(通用锁)
func (o *LotteryManager) SetLtryRecord(k string, v []LotteryRecord) {
	o.recordMapLock.Lock()
	defer o.recordMapLock.Unlock()

	o.lotteriesRecordMap[k] = v
}

//更新一个彩票历史记录
func (o *LotteryManager) UpdateLotteryRecord(record []LotteryRecord) {
	if record != nil && len(record) > 0 {
		//这里验证一下 是不是有这个类型
		if Utils.VerifyGameTag(record[0].GameName) {
			o.recordMapLock.Lock()
			defer o.recordMapLock.Unlock()
			o.lotteriesRecordMap[record[0].GameName] = record
		} else {
			beego.Error("Update lottery record error ! Not have this lottery type  : ", record[0].GameName, " !")
		}
	} else {
		beego.Error("Update lottery record error !")
	}
}

//增加一个采种的一条历史记录,例如掉期补全时
func (o *LotteryManager) UpdateLotteryOneRecord(data LotteryRecord) {

	records, err := o.GetLtryRecord(data.GameName)
	if err != nil {
		beego.Error(err)
		return
	}

	if len(records) > 0 {
		//去头部老记录，增加的新记录在头部
		records = records[1:]
		//头部再添加新的元素
		lr := LotteryRecord{}
		lr.Expect = data.Expect
		lr.GameName = data.GameName
		lr.OpenCode = data.OpenCode
		lr.OpenTime = data.OpenTime
		records = append(records, lr)

		o.SetLtryRecord(data.GameName, records)
	}

	// if records, ok := o.LotteriesRecordMap[data.GameName]; ok {
	// 	if len(records) > 0 {
	// 		//去头部老记录，增加的新记录在头部
	// 		records = records[1:]
	// 		//头部再添加新的元素
	// 		lr := LotteryRecord{}
	// 		lr.Expect = data.Expect
	// 		lr.GameName = data.GameName
	// 		lr.OpenCode = data.OpenCode
	// 		lr.OpenTime = data.OpenTime
	// 		records = append(records, lr)

	// 		o.LotteriesRecordMap[data.GameName] = records
	// 	}
	// }
}
