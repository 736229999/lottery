package Utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

var increnment int = 0

var mutex sync.Mutex

func GetToken() string {
	token := time.Now().Format("20060102150405") + // 年月日时分秒 +
		fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000) + // 纳秒前三位 +
		fmt.Sprintf("%04d", GetIncrementNumFour()) + // 千位循环
		fmt.Sprintf("%03d", rand.Intn(100))
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(token))
	md5Token := hex.EncodeToString(md5Ctx.Sum(nil))
	return md5Token
}

//得到一个4位循环自增数字
func GetIncrementNumFour() int {
	mutex.Lock()
	increnment = (increnment + 1) % 1000
	defer mutex.Unlock()
	return increnment
}

// //检查服务器状态 返回Game服务器url
func CheckGameServerState(url string) bool {
	to, err := time.ParseDuration("5s")
	if err != nil {
		beego.Debug("------------------------- time.ParseDuration() : ", err, " -------------------------")
		return false
	}

	c := &http.Client{
		Timeout: to}
	_, err_0 := c.PostForm(url, nil)
	if err_0 != nil {
		beego.Debug("------------------------- CheckGameServerState ! Server Id : ", err, " -------------------------")
		return false
	}

	return true
}

// type count struct {
// 	ID             string
// 	SEQUENCE_VALUE int
// }
// type test struct {
// 	ID   int
// 	NAME string
// }

// func main() {
// 	session, err := mgo.Dial("192.168.0.180:27017")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	db := session.DB("LoginServer")
// 	collect := db.C("Test2")
// 	for i := 0; i < 10; i++ {
// 		err = collect.Insert(&test{ID: getNextSequenceValue("productid"), NAME: "1111" + strconv.Itoa(i)})
// 	}

// }

// func getNextSequenceValue(s string) int {
// 	session, err := mgo.Dial("192.168.0.180:27017")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	db := session.DB("LoginServer")
// 	counters := db.C("Counters")
// 	counters.Update(bson.M{"id": s}, bson.M{"$inc": bson.M{"sequence_value": 1}})
// 	con := count{}
// 	err = counters.Find(bson.M{"id": s}).One(&con)
// 	return con.SEQUENCE_VALUE
// }
