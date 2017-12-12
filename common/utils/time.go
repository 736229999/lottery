package utils

import (
	"time"

	"github.com/astaxie/beego"
)

const (
	TF_D   = "20060102"
	TF_DWH = "2006-01-02" //Time Format Date With HyPhen 时间格式带连接符的日期
)

//得到现在北京时间
func GetNowUTC8Time() time.Time {
	nowTime := time.Now()
	// 设置时区 PRC 代表中国
	location, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Error(err)
	}

	return nowTime.In(location)
}

//转换为北京时间
func ConvertToUTC8Time(t time.Time) time.Time {
	location, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Error(err)
	}

	return t.In(location)
}

//计算日期差(返回相差天数 now - t)
func DateSub(now time.Time, t time.Time) int {
	// 如果TimeA和TimeB是同一年份，考虑两时间之间可能间隔较近，这样做是为了增强效率
	if t.Year() == now.Year() {
		return now.YearDay() - t.YearDay()
	}
	location, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Emergency(err)
	}

	dateA := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
	dateB := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	return int((dateB.Sub(dateA).Hours()) / 24)
}

//计算给定日期N天前的日期 只算天数
func DateBeforeTheDay(today time.Time, before int) time.Time {
	location, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Emergency(err)
	}

	date := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, location)
	day, _ := time.ParseDuration("-24h")

	result := date.Add(day * time.Duration(before))
	return result
}

//计算给定日期N天后的日期，只算天数
func DateAfterTheDay(today time.Time, after int) time.Time {
	location, err := time.LoadLocation("PRC")
	if err != nil {
		beego.Emergency(err)
	}

	date := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, location)
	day, _ := time.ParseDuration("24h")

	result := date.Add(day * time.Duration(after))
	return result
}

//取整分数(去掉秒数的时间)
func GetIntegerMin(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

//计算几分钟之后的时间
func TimeAfterMin(t time.Time, after int) time.Time {
	m, _ := time.ParseDuration("1m")
	return t.Add(m * time.Duration(after))
}
