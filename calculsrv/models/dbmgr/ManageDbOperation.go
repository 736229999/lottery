//管理后台数据库操作
package dbmgr

import (
	"calculsrv/models/gb"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//从管理员数据库查询当前彩票状态(如果 从后台开关彩票,都要调用这个函数,并且重新通知)
func UpdateLotteriesInfo() map[string]gb.LotteryInfo {
	//查询条件为什么是大类型的彩票信息，不要为我为什么，呵呵，想吐槽
	var result []gb.LotteryInfo

	bsonM := bson.M{"parent_id": bson.M{"$ne": 0}}
	err := LtryColl.Find(bsonM).All(&result)
	if err != nil {
		beego.Emergency("------------------------- Emergency error :", err, "  -------------------------")
		return nil
	}

	var m = make(map[string]gb.LotteryInfo)
	//beego.Debug("从管理后台获得所有菜种: ")
	for _, v := range result {
		m[v.GameTag] = v
		//beego.Debug(v)
	}

	return m
}

//从管理员数据库获取彩票设置
func GetLtrySet() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := LtrySetColl.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
