package dbmgr

import (
	"errors"
	"gamesrv/models/GlobalData"
	"sync"
	"time"

	"gamesrv/models/Order"
	"gamesrv/models/ctrl"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//数据库管理员类
type DbMgr struct {
	calculationSession    *mgo.Session //计算服数据库Session
	calculationDb         *mgo.Database
	AccountInfoCollection *mgo.Collection //用户信息表名
	orderCollection       *mgo.Collection // 订单
	IncrementIdCollection *mgo.Collection //自增数表名

	ManageDbSession            *mgo.Session    //管理员数据库Session
	ManageDb                   *mgo.Database   //管理员数据库名
	TransferBankCardCollection *mgo.Collection //线上转账表名
	OnlinePaymentCollection    *mgo.Collection //在线支付表名
	RechargeOrderCollection    *mgo.Collection //充值订单表名
	DrawingsOrderCollection    *mgo.Collection //提款订单表名
	PlayerGroupCollection      *mgo.Collection //用户组表名
	PayTypeCollection          *mgo.Collection //支付类型表名
	InvitationCodeCollection   *mgo.Collection //邀请码表
	PlatformConfCollection     *mgo.Collection //平台相关信息表
	AgentCountHour             *mgo.Collection //代理商统计包边(按小时)
	AgentCountDay              *mgo.Collection //代理商统计包边(按天)
	AgentCountMonth            *mgo.Collection //代理商统计包边(按月)
	ActivityCollection         *mgo.Collection //活动信息
	AnnouncementCollection     *mgo.Collection //公告信息
	QrCodeCollection           *mgo.Collection //充值二维码信息
}

var sInstance *DbMgr
var once sync.Once

/*
取得数据库连接类实例,单例模式
*/
func Instance() *DbMgr {
	once.Do(func() {
		sInstance = &DbMgr{}
	})

	return sInstance
}

//每个彩票最后一次开奖期数	k 为GameTag
//var LotteriesLastRecordtMap map[string]DbsLotteryRecord = make(map[string]DbsLotteryRecord)

/*
初始化工作
链接数据库
*/
func (o *DbMgr) Init() error {
	//判断本服务器是试玩还是正式服
	if ctrl.SelfSrv.Type == 0 { //试玩服务器
		beego.Info("--- Game Server  : Trial !")
		//-------------------------------------- 计算服数据库 ------------------------------------------------
		var trialDbUrl = ctrl.DbSrv.Ip + ":" + ctrl.DbSrv.Port

		beego.Info("--- Calculation DB URL : ", trialDbUrl)

		dialInfo := &mgo.DialInfo{
			Addrs:    []string{trialDbUrl},
			Timeout:  time.Second * 3,
			Username: TrialDbUserName,
			Password: TrialDbPwd,
		}
		//连接计算服数据库
		var err error
		o.calculationSession, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}
		// 采用 Strong 模式
		o.calculationSession.SetMode(mgo.Strong, true)
		o.calculationDb = o.calculationSession.DB(CalculationDbName)

		o.AccountInfoCollection = o.calculationDb.C(accountInfoCollection)

		o.orderCollection = o.calculationDb.C(OrderCollection)

		//试玩服务器不用连接正式的管理服，所以管理服的信息就写在试玩服里面就可以了
		o.ManageDb = o.calculationSession.DB(ManageDbName)
		//试玩服没有这些表,这里只是为了统一
		o.TransferBankCardCollection = o.ManageDb.C(TransferBankCardCollection)
		o.OnlinePaymentCollection = o.ManageDb.C(OnlinePaymentCollection)
		o.RechargeOrderCollection = o.ManageDb.C(RechargeOrderCollection)
		o.DrawingsOrderCollection = o.ManageDb.C(DrawingsOrderCollection)
	} else if ctrl.SelfSrv.Type == 1 {
		beego.Info("--- Game Server  : Formal !")
		//-------------------------------------- 计算服数据库 ------------------------------------------------
		var dburl = ctrl.DbSrv.Ip + ":" + ctrl.DbSrv.Port

		beego.Info("--- Calculation DB URL : ", dburl)

		dialInfo := &mgo.DialInfo{
			Addrs:    []string{dburl},
			Timeout:  time.Second * 3,
			Database: CalculationDbName,
			Username: CalculationDbUserName,
			Password: CalculationDbPwd,
		}
		//连接计算服数据库
		var err error
		o.calculationSession, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}
		// 采用 Strong 模式
		o.calculationSession.SetMode(mgo.Strong, true)
		o.calculationDb = o.calculationSession.DB(CalculationDbName)

		o.AccountInfoCollection = o.calculationDb.C(accountInfoCollection)

		o.orderCollection = o.calculationDb.C(OrderCollection)

		o.IncrementIdCollection = o.calculationDb.C(IncrementIdCollection)

		//-------------------------------------- 管理服数据库 ------------------------------------------------
		var mgrdbUrl = ctrl.MgrDb.Ip + ":" + ctrl.MgrDb.Port

		beego.Info("--- Mgr DB URL : ", mgrdbUrl)

		dialInfo_1 := &mgo.DialInfo{
			Addrs:    []string{mgrdbUrl},
			Timeout:  time.Second * 3,
			Database: ManageDbName,
			Username: ManageDbUserName,
			Password: ManageDbPwd,
		}

		o.ManageDbSession, err = mgo.DialWithInfo(dialInfo_1)
		if err != nil {
			return err
		}

		o.ManageDbSession.SetMode(mgo.Strong, true)

		o.ManageDb = o.ManageDbSession.DB(ManageDbName)

		o.TransferBankCardCollection = o.ManageDb.C(TransferBankCardCollection)

		o.OnlinePaymentCollection = o.ManageDb.C(OnlinePaymentCollection)

		o.RechargeOrderCollection = o.ManageDb.C(RechargeOrderCollection)

		o.DrawingsOrderCollection = o.ManageDb.C(DrawingsOrderCollection)

		o.PlayerGroupCollection = o.ManageDb.C(PlayerGroupCollection)

		o.PayTypeCollection = o.ManageDb.C(PayTypeCollection)

		o.InvitationCodeCollection = o.ManageDb.C(InvitationCode)

		o.PlatformConfCollection = o.ManageDb.C(PlatformConf)

		o.AgentCountHour = o.ManageDb.C(AgentCountHour)

		o.AgentCountDay = o.ManageDb.C(AgentCountDay)

		o.AgentCountMonth = o.ManageDb.C(AgentCountMonth)

		o.ActivityCollection = o.ManageDb.C(Activity)

		o.AnnouncementCollection = o.ManageDb.C(Announcement)

		o.QrCodeCollection = o.ManageDb.C(QrCode)
	} else {
		beego.Error("Login Server Type Error !")
	}

	beego.Info("------------------------- Init DB Mgr Done ! ------------------------- ")
	return nil
}

//得到一个用户信息(AccountMgr专用,以后所有的AccountInfo的地方都要使用这个新的AccountMgr里面的AccountInfo类)
func (o DbMgr) InitAccountInfo(accountName string, accountInfo interface{}) error {
	err := o.AccountInfoCollection.Find(bson.M{"account_name": accountName}).One(accountInfo)
	if err != nil {
		return err
	}
	return nil
}

//通过账号ID来得到一个用户信息(同上)
func (o DbMgr) InitAccountInfoById(accountID int, accountInfo interface{}) error {
	err := o.AccountInfoCollection.Find(bson.M{"account_id": accountID}).One(accountInfo)
	if err != nil {
		return err
	}
	return nil
}

//得到一个代理商信息
func (o DbMgr) InitAgentInfo(agentId int, accountInfo interface{}) error {
	err := o.AccountInfoCollection.Find(bson.M{"agent_id": agentId}).One(accountInfo)
	if err != nil {
		return err
	}
	return nil
}

//插入数据
func (o DbMgr) Insert(collection *mgo.Collection, msg interface{}) bool {
	err := collection.Insert(msg)
	if err != nil {
		beego.Emergency("------------------------- Insert Error ", err, " ------------------------- ")
		return false
	}
	return true
}

//批量插入
func (o DbMgr) bulkInsert(collection *mgo.Collection, msg *[]interface{}) bool {
	err := collection.Insert(*msg...)
	if err != nil {
		beego.Emergency("------------------------- Bulk Insert Error : ", err, " ------------------------- ")
		return false
	}
	return true
}

func (o DbMgr) FindAccount(account string) *GlobalData.AccountInfo {
	//查询单条数据
	result := &GlobalData.AccountInfo{}
	err := o.AccountInfoCollection.Find(bson.M{"account_name": account}).One(&result)
	if err != nil {
	}
	return result
}

// func (o DbMgr) UpdateAccountInfo(account string, flag string, token string, lastLoginTime time.Time, lastLoginTimeStamp int64, lastLoginIp string) bool {
// 	//更新数据
// 	selector := bson.M{"account_name": account}
// 	data := bson.M{"$set": bson.M{"flag": flag, "token": token}}

// 	err := o.accountInfoCollection.Update(selector, data)
// 	if err != nil {
// 		beego.Error(err)
// 		return false
// 	}
// 	return true
// }

// fun

func (o DbMgr) InsertAccountInfo(msg interface{}) bool {
	err := o.AccountInfoCollection.Insert(msg)
	if err != nil {
		beego.Error("------------------------- Bulk Insert Error : ", err, " ------------------------- ")
		return false
	}
	return true
}

//通过账号名得到这个帐号的订单(注意,之前要验证 帐号和token是否合法才能进行这个查询, 等上线后这个改为在game服务器建立缓存池,不要每次都从数据库读取)
func (o DbMgr) GetOrderByAccountName(accountName string, skip int, limit int, searchType int, ret *([]Order.Order)) error {
	var bm bson.M
	switch searchType {
	case 1: //全部-全部
		bm = bson.M{"account_name": accountName}
	case 2: //全部中奖
		bm = bson.M{"account_name": accountName, "winning_bet_num": bson.M{"$ne": 0}}
	case 3: //全部未开奖
		bm = bson.M{"account_name": accountName, "status": 0}
	case 4: //普通-全部
		bm = bson.M{"account_name": accountName, "order_type": 0}
	case 5: //普通-中奖
		bm = bson.M{"account_name": accountName, "order_type": 0, "winning_bet_num": bson.M{"$ne": 0}}
	case 6: //普通-待开奖
		bm = bson.M{"account_name": accountName, "order_type": 0, "status": 0}
	case 7: //追号-全部
		bm = bson.M{"account_name": accountName, "order_type": 1}
	case 8: //追号-中奖
		bm = bson.M{"account_name": accountName, "order_type": 1, "winning_bet_num": bson.M{"$ne": 0}}
	case 9: //追号-待开奖
		bm = bson.M{"account_name": accountName, "order_type": 1, "status": 0}
	default:
		return errors.New("订单查询失败, 没有这个查询类型")
	}

	err := o.orderCollection.Find(bm).Sort("-betting_time").Skip(skip).Limit(limit).All(ret)
	if err != nil {
		return err
	}
	return nil
}

//通过账号名得到这个帐号的订单(注意,之前要验证 帐号和token是否合法才能进行这个查询, 等上线后这个改为在game服务器建立缓存池,不要每次都从数据库读取)
func (o DbMgr) GetOrderByAccountNameAndGameTag(accountName string, gameTag string, expect int, skip int, limit int, ret *([]Order.Order)) error {
	bsonM := bson.M{"account_name": accountName, "game_tag": gameTag, "expect": expect}
	//按时间倒序；
	err := o.orderCollection.Find(bsonM).Sort("-betting_time").Skip(skip).Limit(limit).All(ret)
	if err != nil {
		return err
	}
	return nil
}

//得到最新的20条中奖记录 查询条件为 WinningBetNumber > 0
func (o DbMgr) GetNewestWinning(ret *([]GlobalData.NewestWinningInfo)) error {
	bsonM := bson.M{"winning_bet_num": bson.M{"$gt": 0}}
	err := o.orderCollection.Find(bsonM).Sort("-_id").Limit(20).All(ret)
	if err != nil {
		return err
	}
	return nil
}

//更新 用户附加信息(包括银行等信息)
func (o DbMgr) UpdateAccountAdditionInfo(addInfo GlobalData.ReqModifyAdditionalInfo) error {
	selector := bson.M{"account_name": addInfo.AccountName}
	data := bson.M{"$set": bson.M{"mobile_phone": addInfo.Mobile_Phone, "qq": addInfo.QQ, "wechat": addInfo.WeChat, "weibo": addInfo.WeiBo, "email": addInfo.Email, "address": addInfo.Address, "bank_name": addInfo.Bank_Name, "bank_card": addInfo.Bank_Card, "card_holder": addInfo.Card_Holder, "bank_of_deposit": addInfo.Bank_Of_Deposit}}
	err := o.AccountInfoCollection.Update(selector, data)
	if err != nil {
		return err
	}
	return nil
}

//更新用户附加信息（不包括银行等信息）
func (o DbMgr) UpdateAccountAdditionInfoNotHaveBank(addInfo GlobalData.ReqModifyAdditionalInfo) error {
	selector := bson.M{"account_name": addInfo.AccountName}
	data := bson.M{"$set": bson.M{"mobile_phone": addInfo.Mobile_Phone, "qq": addInfo.QQ, "wechat": addInfo.WeChat, "weibo": addInfo.WeiBo, "email": addInfo.Email}}
	err := o.AccountInfoCollection.Update(selector, data)
	if err != nil {
		return err
	}
	return nil
}

//更新资金密码(以后有空来将这里的资金密码改为 int)
func (o DbMgr) UpdateMoneyPassword(accountName string, newMoneyPassword string) error {
	selector := bson.M{"account_name": accountName}
	data := bson.M{"$set": bson.M{"money_password": newMoneyPassword}}
	err := o.AccountInfoCollection.Update(selector, data)
	if err != nil {
		return err
	}
	return nil
}

//获得查询记录
func (o DbMgr) GetRechargeRecord(accountName string, skip int, limit int, searchType int, ret *([]GlobalData.RechargeRecord)) error {
	var bm bson.M
	switch searchType {
	case 1: //全部
		bm = bson.M{"account_name": accountName}
	case 2: //成功
		bm = bson.M{"account_name": accountName, "status": 1}
	case 3: //等待审核
		bm = bson.M{"account_name": accountName, "status": 2}
	default:
		errors.New("获取充值记录失败,没有对应的查询类型")
	}

	err := o.RechargeOrderCollection.Find(bm).Sort("-recharge_time").Skip(skip).Limit(limit).All(ret)
	if err != nil {
		return err
	}
	return nil
}

//获得提款记录
func (o DbMgr) GetDrawingsRecord(accountName string, skip int, limit int, searchType int, ret *([]GlobalData.DrawingsRecord)) error {
	var bm bson.M
	switch searchType {
	case 1: //全部
		bm = bson.M{"account_name": accountName}
	case 2: //成功
		bm = bson.M{"account_name": accountName, "status": 1}
	case 3: //等待审核
		bm = bson.M{"account_name": accountName, "$or": []bson.M{bson.M{"status": 2}, bson.M{"status": 3}}}
	default:
		errors.New("获取充值记录失败,没有对应的查询类型")
	}

	err := o.DrawingsOrderCollection.Find(bm).Sort("-drawings_time").Skip(skip).Limit(limit).All(ret)
	if err != nil {
		return err
	}
	return nil
}

//根据用户组查询支付类型(指定查询PayType)
func (o DbMgr) GetPayTypeByGroupId(groupId int) (int, error) {
	//还是要定义接收结构体
	type Ret struct {
		PayType int `bson:"pay_type"`
	}
	//接收变量
	var ret Ret
	//查询条件
	bm := bson.M{"id": groupId}
	//指定查询字段,1代表查询,0代表剔除这个字段
	sm := bson.M{"pay_type": 1}
	err := o.PlayerGroupCollection.Find(bm).Select(sm).One(&ret)
	if err != nil {
		return 0, err
	}
	return ret.PayType, nil
}

//根据pay_type id来查找线下银行稽核数据
func (o DbMgr) GetInspectInfoBank(payTypeId int) (GlobalData.InspectMoneyBank, error) {
	bm := bson.M{"id": payTypeId}
	var sm bson.M
	sm = bson.M{"drawings.bank": 1}

	var ret GlobalData.InspectMoneyBank
	err := o.PayTypeCollection.Find(bm).Select(sm).One(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

//根据pay_type id来查找线上稽核数据
func (o DbMgr) GetInspectInfoOnline(payTypeId int) (GlobalData.InspectMoneyOnline, error) {
	bm := bson.M{"id": payTypeId}
	var sm bson.M
	sm = bson.M{"drawings.online": 1}
	var ret GlobalData.InspectMoneyOnline
	err := o.PayTypeCollection.Find(bm).Select(sm).One(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

//得到账户自增值
func (o DbMgr) GetAccountIncrmentId() (int, error) {
	bm := bson.M{"name": "account_inc"}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	var ret map[string]interface{}
	_, err := o.IncrementIdCollection.Find(bm).Apply(change, &ret)
	if err != nil {
		return 0, err
	}
	return ret["id"].(int), nil
}

//根据邀请码,查找信息
func (o DbMgr) GetInviteCodeRelatedInfo(inviteCode string) (GlobalData.InviteCodeRelatedInfo, error) {
	bm := bson.M{"code": inviteCode}

	ret := GlobalData.InviteCodeRelatedInfo{}
	err := o.InvitationCodeCollection.Find(bm).One(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

// //更新代理商信息
// func (o DbMgr) UpdateAgentInfo(agentInfo interface{}) bool {
// 	//更新数据
// 	selector := bson.M{"account_name": agentName}
// 	//data := bson.M{"$set": bson.M{"flag": flag, "token": token}}
// 	err := o.accountInfoCollection.Update(selector, agentInfo) //放入整个结构体
// 	if err != nil {
// 		beego.Error(err)
// 		return false
// 	}
// 	return true
// }

//得到代理商邀请码自增ID
func (o DbMgr) GetAgentInviteCodeIncrmentId() (int, error) {
	bm := bson.M{"name": "code_agent"}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	var ret map[string]interface{}
	_, err := o.IncrementIdCollection.Find(bm).Apply(change, &ret)
	if err != nil {
		return 0, err
	}
	return ret["id"].(int), nil
}

//得到用户邀请码自增ID
func (o DbMgr) GetUserInviteCodeIncrmentId() (int, error) {
	bm := bson.M{"name": "code_user"}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	var ret map[string]interface{}
	_, err := o.IncrementIdCollection.Find(bm).Apply(change, &ret)
	if err != nil {
		return 0, err
	}
	return ret["id"].(int), nil
}
