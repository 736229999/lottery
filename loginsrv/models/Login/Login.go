package Login

import (
	"bytes"
	"encoding/json"
	"loginsrv/models/Utils"
	"loginsrv/models/ctrl"
	"net/http"
	"regexp"

	"errors"

	"github.com/astaxie/beego"
)

//验证账户名合法
func VerifyAccountName(accountName string) bool {
	if accountName == "" {
		return false
	}

	l1 := len(accountName)
	if l1 < 6 && l1 > 20 {
		return false
	}
	match, err := regexp.MatchString("^[0-9a-z]*$", accountName)
	if match == false || err != nil {
		return false
	}
	return true
}

//验证密码合法
func VerifyPassword(password string) bool {
	if password == "" {
		return false
	}

	l2 := len(password)
	if l2 < 6 && l2 > 20 {
		return false
	}
	match, err := regexp.MatchString("^[0-9a-z]*$", password)
	if match == false || err != nil {
		return false
	}
	return true
}

//预登陆流程
func PrepareLogin(prepareLogin *ReqPrepareLogin, gameIp string) bool {
	//没有找到可用GameServer
	if gameIp == "" {
		beego.Emergency("--------------------------------- 严重错误, 没有可用的 Game Server !-----------------------------------")
		return false
	}

	prepareMsg, err := json.Marshal(prepareLogin)
	if err != nil {
		beego.Emergency("--------------------------------- Marshal Error : ", err, "-----------------------------------")
		return false
	}

	body := bytes.NewBuffer(prepareMsg)
	//这里采用阻塞，确保 服务器收到了以后再通知客户端可以登陆
	resp, err1 := http.Post("http://"+gameIp+"/PrePlayer", "application/json;charset=utf-8", body)
	if err1 != nil {
		beego.Emergency("----------------------------- 严重错误,发送PrePlayer失败 Server : ", gameIp, " --------------------------------")
		return false
	}
	resp.Body.Close()

	return true
}

//得到服务器ip（最新修改为正式服和试玩服务器完全分离，所以返回gameIp就从0开始就行了）
func GetGameIp() (string, error) {

	for _, v := range ctrl.GameSrv {
		if Utils.CheckGameServerState("http://" + v.Ip + ":" + v.Port + "/CheckGameSrvStatus") {
			//gameIp = ip
			return v.Ip + ":" + v.Port, nil
		} else {
			beego.Error("--- 出现Game服务器无法通信的情况,请检查服务器 Server IP : ", v.Ip)
		}
	}

	return "", errors.New("--- 严重错误, 未能找到可用的 Game Server ！请立即检查 Game Server服务器 ")
}
