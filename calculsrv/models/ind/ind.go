package ind

import (
	"bytes"
	"calculsrv/models/ctrl"
	"calculsrv/models/encmgr"
	"calculsrv/models/gb"
	"common/httpmgr"
	"encoding/json"
	"errors"
	"time"

	"github.com/astaxie/beego"
)

//独立彩票要冲对应组的独立彩票服务器去获取记录
func GetRecordByNew(gameName string) (gb.LtryRecordByNewest, error) {
	req := ReqLtryNewRecord{}
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

	count := 3 //重试次数,这里暂时写死
	for ; count > 0; count-- {
		resp, err := httpmgr.Post("http://"+ctrl.IndSrv.Ip+":"+ctrl.IndSrv.Port+"/getRecordByNew", body)
		if err != nil {
			beego.Warn("从API服务器获取彩票信息失败, API : ", ctrl.IndSrv.Ip, " Error : ", err, " 剩余重试次数 : ", count, " 彩票名称 : ", gameName)
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

		return ret, nil
	}

	str := "Get lottery record by newest fail ! Game name : " + gameName
	return ret, errors.New(str)
}

//获取独立彩票的历史记录
func GetHist(gameName string) ([]RespGetHist, error) {
	req := ReqGetHist{}
	ret := []RespGetHist{}

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

	count := 3 //重试次数,这里暂时写死
	for ; count > 0; count-- {
		resp, err := httpmgr.Post("http://"+ctrl.IndSrv.Ip+":"+ctrl.IndSrv.Port+"/getHist", body)
		if err != nil {
			beego.Warn("从API服务器获取彩票信息失败, API : ", ctrl.IndSrv.Ip, " Error : ", err, " 重试次数 : ", count, " 彩票名称 : ", gameName)
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

		return ret, nil
	}

	str := "Get lottery record by newest fail ! Game name : " + gameName
	return ret, errors.New(str)
}

//按期数获得某一期某个彩票的历史记录
func GetRecordByExpect(gameName string, expect int) (gb.LotteryRecord, error) {
	req := ReqRecordByExpect{}
	ret := gb.LotteryRecord{}

	req.GameName = gameName
	req.Expect = expect

	data, err := json.Marshal(req)
	if err != nil {
		return ret, err
	}

	cipher, err := encmgr.Instance().AesPrkEnc(data)
	if err != nil {
		return ret, err
	}

	body := bytes.NewBuffer(cipher)

	count := 3 //重试次数,这里暂时写死
	for ; count > 0; count-- {
		resp, err := httpmgr.Post("http://"+ctrl.IndSrv.Ip+":"+ctrl.IndSrv.Port+"/getRecordByExpect", body)
		if err != nil {
			beego.Warn("从API服务器获取彩票信息失败, API : ", ctrl.IndSrv.Ip, " Error : ", err, " 重试次数 : ", count, " 彩票名称 : ", gameName)
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

		return ret, nil
	}

	str := "Get lottery record by expect fail ! Game name : " + gameName
	return ret, errors.New(str)
}

type ReqLtryNewRecord struct {
	GameName string
}

type ReqGetHist struct {
	GameName string
}

type RespGetHist struct {
	GameName           string    // 名称
	Expect             int       // 期次
	OpenCode           string    // 开奖号码
	OpenTime           time.Time // 开奖时间
	OpenTimeStamp      int64     //开奖时间(时间戳)
	RecordingTime      time.Time //记录时间(从api 获取到结果存库的时间,理论上 最多减去间隔访问时间,就是实际获取结果时间,也就是我们这里的开奖时间)
	RecordingTimeStamp int64     //记录时间(时间戳)
}

type ReqRecordByExpect struct {
	GameName string
	Expect   int
}
