package DbHandle

import (
	"loginsrv/models/GlobalData"
	"loginsrv/models/ctrl"
	"loginsrv/models/dbmgr"
	"time"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var userInfoCollection *mgo.Collection

func Init() {
	//判断改服务器是试玩还是正式服;
	// serverType, err := beego.AppConfig.Int("ServerType")
	// if err != nil {
	// 	beego.Emergency(err)
	// 	return
	// }

	//if serverType == 0 {
	userInfoCollection = dbmgr.Instance().DbServiceMap[ctrl.DbSrv.Ip].CollectionMap[dbmgr.AccountInfoCollection]
	// } else {
	// 	userInfoCollection = dbmgr.Instance().DbServiceMap[ctrl.DbSrv.Ip].CollectionMap[dbmgr.AccountInfoCollection]
	// }
}

//DBHandleFind 数据库查询数据
// func FindAccount(account string) string {
// 	result := GlobalData.AccountInfo{}
// 	err := userInfoCollection.Find(bson.M{"account_name": account}).One(&result)
// 	if err != nil {
// 		return ""
// 	}
// 	//beego.Debug("-------------" + result.Password + "------------------")
// 	return result.Password
// }

//得到账号信息（新）
func GetAccountInfo(accountName string, ret interface{}) error {
	err := userInfoCollection.Find(bson.M{"account_name": accountName}).One(ret)
	if err != nil {
		return err
	}
	return nil
}

func FindAccountInfo(account string) GlobalData.AccountInfo {
	result := GlobalData.AccountInfo{}
	err := userInfoCollection.Find(bson.M{"account_name": account}).One(&result)
	if err != nil {
		return result
	}
	//beego.Debug("-------------" + result.Password + "------------------")
	return result
}

//DBHandleInsert 数据库插入数据
func InsertAccount(account GlobalData.AccountInfo) bool {
	err := userInfoCollection.Insert(account)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

//用户登录,更新最后登录信息
func UpdateAccountInfo(account string, lastLoginTime time.Time, lastLoginTimeStamp int64, lastLoginIp string) bool {
	selector := bson.M{"account_name": account}
	data := bson.M{"$set": bson.M{"last_login_ip": lastLoginIp, "last_login_time": lastLoginTime, "last_login_time_stamp": lastLoginTimeStamp}}
	err := userInfoCollection.Update(selector, data)
	if err != nil {
		beego.Debug(err)
		return false
	}
	return true
}

//更新密码
func UpdateAccountPassword(accountName string, newPassword string) bool {
	selector := bson.M{"account_name": accountName}
	data := bson.M{"$set": bson.M{"password": newPassword}}
	err := userInfoCollection.Update(selector, data)
	if err != nil {
		beego.Debug(err)
		return false
	}
	return true
}
