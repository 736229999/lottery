package AccountMgr

import "time"
import "loginsrv/models/DbHandle"

type AccountInfo struct {
	//Inc_ID                string    //自增ID
	Account_Name          string    //用户名
	Password              string    //密码
	Regist_Time           time.Time //注册时间
	Regist_Time_Stamp     int64     //注册时间,时间戳
	Regist_Ip             string    //注册IP
	Last_Login_Time       time.Time //最后登录时间
	Last_Login_Time_Stamp int64     //最后登录时间戳
	Last_Login_Ip         string    //最后登录IP
	Registration_Platform int       //注册平台
	Account_Type          int       //账户类型 o 为试玩用户, 1为正常用户
}

func (o *AccountInfo) Init(accountName string) error {
	err := DbHandle.GetAccountInfo(accountName, o)
	if err != nil {
		return err
	}
	return nil
}

func (o *AccountInfo) VerifyPassword(pw string) bool {
	if o.Password == pw {
		return true
	}
	return false
}

func (o *AccountInfo) ModifyPassword(oldPW string, newPw string) bool {
	if o.Password == oldPW {
		o.Password = newPw
		//存库
		if !DbHandle.UpdateAccountPassword(o.Account_Name, newPw) {
			return false
		}
		return true
	} else {
		return false
	}
}
