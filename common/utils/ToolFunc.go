//公用工具函数

package utils

import (
	"calculsrv/models/gb"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

//将从API或数据库中获得的开奖号码转换成int数组
func PaserOpenCodeToArray(openCode string) []int {
	code := strings.Split(openCode, ",")
	var arr []int
	for _, v := range code {
		i, _ := strconv.Atoi(v)
		arr = append(arr, i)
	}
	return arr
}

//分析组合(n个数里面 m个数为一组，有几种组合,不重复)(订单注数)
func AnalysisCombination(n int, m int) int {
	//公式为 n!/(m!*(n-m)!)
	return Factorial(n) / (Factorial(m) * Factorial(n-m))
}

//每有一个胆码重复,总注数就要减去 (n - 1)! / ((m - 1)! * (n - m)!)
//例如 : 在SSC 组选60中 , n为6, m为3 ,这个公式就是计算的n为6 和 n为5之间的差值
func CombinationDifference(n int, m int) int {
	return Factorial(n-1) / (Factorial(m-1) * Factorial(n-m))
}

//阶乘函数
func Factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

//转换数组到Interface
func ConvertArrayToInterface(array interface{}) []interface{} {
	if reflect.TypeOf(array).Kind() != reflect.Slice {
		beego.Debug("------------------------- Convert Array To Interface Error ! -------------------------")
		return nil
	}
	arrValue := reflect.ValueOf(array)
	resultArray := make([]interface{}, arrValue.Len())
	for i := 0; i < arrValue.Len(); i++ {
		resultArray[i] = arrValue.Index(i).Interface()
	}
	return resultArray
}

//将从数据库获得的interface数据转换为 map[string]float64格式方便使用
func OddsInterface2Map(src interface{}) map[string]float64 {
	var oddsMap map[string]float64 = make(map[string]float64)
	m := src.(map[string]interface{})
	for k, v := range m {
		oddsMap[k] = v.(float64)
	}
	return oddsMap
}

//获取订单号(注意现在订单号和现金交易流水号用的是同一规则)
var index int = 0

func GetOrderNumber() string {
	index = (index + 1) % 1000
	// 本地编号取四位
	orderNum := time.Now().Format("20060102150405") + // 年月日时分秒 +
		fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000) + // 纳秒前三位 +
		fmt.Sprintf("%03d", index) + // 千位循环
		gb.MachineCode //计算服id
	return orderNum
}
