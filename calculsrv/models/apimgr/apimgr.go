//添加新彩票时 需要在 Config 文件中添加 Const 字段,然后在 init 函数中 将新彩票的Api放入map中

package apimgr

import (
	"bytes"
	"calculsrv/models/ctrl"
	"calculsrv/models/encmgr"
	"calculsrv/models/gb"
	"common/httpmgr"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

//Api管理员类
type ApiMgr struct {
}

var sInstance *ApiMgr
var once sync.Once

// //单例模式
func Instance() *ApiMgr {
	once.Do(func() {
		sInstance = &ApiMgr{}
		sInstance.init()
	})
	return sInstance
}

//初始化 添加彩票Api 都在这里面
func (o *ApiMgr) init() {

}

//从api服务器获取最新记录(带下期)(这里不会一直请求,报错会返回,需要在外部设计直到请求到正确结果的循环)
func (o *ApiMgr) GetRecordByNewest(gameName string) (gb.LtryRecordByNewest, error) {
	req := ReqLtryNewestRecord{}
	ret := gb.LtryRecordByNewest{}

	req.GameName = gameName

	data, err := json.Marshal(req)
	if err != nil {
		return ret, err
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		return ret, err
	}

	body := bytes.NewBuffer(cipher)

	//在api服务器中循环请求(可能存在多个api服务器的情况)
	for _, v := range ctrl.ApiSrv {
		count := 3 //重试次数,这里暂时写死
		for ; count > 0; count-- {
			resp, err := httpmgr.Post("http://"+v.Ip+":"+v.Port+"/getLtryNewestRecord", body)
			if err != nil {
				beego.Warn("从API服务器获取彩票信息失败, API : ", v.Ip, " Error : ", err, " 重试次数 : ", count, " 彩票名称 : ", gameName)
				//休眠5秒再进行下次请求
				time.Sleep(5 * time.Second)
				continue
			}

			//解码错误, 和 json 解析错误 直接返回
			plaintext, err := encmgr.Instance().AesPrkDec(resp)
			if err != nil {
				return ret, err
			}

			err = json.Unmarshal(plaintext, &ret)
			if err != nil {
				return ret, err
			}

			//beego.Debug(ret)
			return ret, nil
		}
	}

	str := "Get lottery record by newest fail ! Game name : " + gameName
	return ret, errors.New(str)
	//httpmgr.Post(ctrl.Api)
}

type ReqLtryNewestRecord struct {
	GameName string
}

type RespLtryNewestRecord struct {
	GameName        string
	CurrentExpect   int
	OpenCode        []int
	OpenCodeStr     string
	CurrentOpenTime time.Time

	NextExpect   int
	NextOpenTime time.Time
}

//获取一个菜种的历史记录
func (o *ApiMgr) GetLtryHist(gameName string) ([]RespGetLtryHist, error) {
	req := ReqGetLtryHist{}

	req.GameName = gameName

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(cipher)

	//在api服务器中循环请求(可能存在多个api服务器的情况)

	ret := []RespGetLtryHist{}
	for _, v := range ctrl.ApiSrv {
		count := 3 //重试次数,这里暂时写死
		for ; count > 0; count-- {
			resp, err := httpmgr.Post("http://"+v.Ip+":"+v.Port+"/getLtryHist", body)
			if err != nil {
				beego.Warn("从API服务器获取彩票信息历史记录失败, API : ", v.Ip, " Error : ", err, " 重试次数 : ", count)
				continue
			}
			//beego.Debug(resp)
			//解码错误, 和 json 解析错误 直接返回
			plaintext, err := encmgr.Instance().AesPrkDec(resp)
			if err != nil {
				return ret, err
			}

			err = json.Unmarshal(plaintext, &ret)
			if err != nil {
				return ret, err
			}

			//beego.Debug(ret)
			return ret, nil
		}
	}

	return ret, errors.New("Get lottery record by newest fail !")
}

type ReqGetLtryHist struct {
	GameName string
}

//获取某一个菜种某一期的历史记录
func (o *ApiMgr) GetLtryRecordByExpect(gameName string, expect int) (gb.LotteryRecord, error) {
	req := ReqGetLtryRecordByExpect{}
	req.GameName = gameName
	req.Expect = expect

	ret := gb.LotteryRecord{}

	data, err := json.Marshal(req)
	if err != nil {
		return ret, err
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		return ret, err
	}

	body := bytes.NewBuffer(cipher)

	//在api服务器中循环请求(可能存在多个api服务器的情况)
	for _, v := range ctrl.ApiSrv {
		count := 3 //重试次数,这里暂时写死
		for ; count > 0; count-- {
			resp, err := httpmgr.Post("http://"+v.Ip+":"+v.Port+"/getLtryRecordByExpect", body)
			if err != nil {
				beego.Warn("从API服务器获取彩票信息失败, API : ", v.Ip, " Error : ", err, " 重试次数 : ", count)
				continue
			}

			//解码错误, 和 json 解析错误 直接返回
			plaintext, err := encmgr.Instance().AesPrkDec(resp)
			if err != nil {
				return ret, err
			}

			err = json.Unmarshal(plaintext, &ret)
			if err != nil {
				return ret, err
			}

			return ret, nil
		}
	}

	return ret, errors.New("Get lottery record by expect fail !")
}

//按天获取一个菜种的历史记录
func (o *ApiMgr) GetLtryHistByDay(gameName string, date time.Time) (gb.LotteryRecordByDayFromApi, error) {
	req := ReqGetLtryHistByDay{}
	req.GameName = gameName
	req.Date = date

	ret := gb.LotteryRecordByDayFromApi{}

	data, err := json.Marshal(req)
	if err != nil {
		return ret, err
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		return ret, err
	}

	body := bytes.NewBuffer(cipher)

	//在api服务器中循环请求(可能存在多个api服务器的情况)
	for _, v := range ctrl.ApiSrv {
		count := 3 //重试次数,这里暂时写死
		for ; count > 0; count-- {
			resp, err := httpmgr.Post("http://"+v.Ip+":"+v.Port+"/getLtryHistByDay", body)
			if err != nil {
				beego.Warn("从API服务器获取彩票信息失败, API : ", v.Ip, " Error : ", err, " 重试次数 : ", count)
				continue
			}

			//解码错误, 和 json 解析错误 直接返回
			plaintext, err := encmgr.Instance().AesPrkDec(resp)
			if err != nil {
				return ret, err
			}

			err = json.Unmarshal(plaintext, &ret)
			if err != nil {
				return ret, err
			}

			//beego.Debug(ret)
			return ret, nil
		}
	}

	return ret, errors.New("Get lottery record by day fail !")

}

type ReqGetLtryHistByDay struct {
	GameName string
	Date     time.Time
}

type ReqGetLtryRecordByExpect struct {
	GameName string
	Expect   int
}

type RespGetLtryHist struct {
	GameName           string    // 名称
	Expect             int       // 期次
	OpenCode           string    // 开奖号码
	OpenTime           time.Time // 开奖时间
	OpenTimeStamp      int64     //开奖时间(时间戳)
	RecordingTime      time.Time //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	RecordingTimeStamp int64     //记录时间(时间戳)
	//maxNum    int       `bson:"id"`       				   // 最大号码取值
}
