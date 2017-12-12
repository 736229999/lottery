package BalanceRecordMgr

//现金流水(资产记录)
type BalanceRecord struct {
	Serial_Number string  //流水号
	Account_name  string  //账号名
	Money_Before  float64 //改变之前玩家钱钱
	Money         float64 //加减的钱钱
	Money_After   float64 //改变之后玩家的钱钱
	Gap_Money     float64 //优惠金额(这里有很多情况)
	Type          int     //类型 1,订单 2,充值	3,提款 4,活动(对于计算服务器,这里只有1的情况)
	Subitem       int     //对于计算服务器来说 tpye 只会为1 ,subitem 1 表示投注, 2表示结算
	Trading_Time  int64   //交易时间
	Status        int     //状态,1 成功 (对于我服务器来说 这里只有1 成功)
	Order_Number  string  //这条记录所对应的订单号
}

//存入数据库
// func (o BalanceRecord) Save()  {
// 	DbMgr.Instance().InsertOrder(o)
// }
