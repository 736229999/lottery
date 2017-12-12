package utils

//得到该服务器是否正式服
// func IsFormalServer() bool {
// 	//判断改服务器是试玩还是正式服;
// 	serverType, err := beego.AppConfig.Int("ServerType")
// 	if err != nil {
// 		beego.Emergency(err)
// 	}

// 	if serverType == 0 {
// 		return false
// 	} else {
// 		return true
// 	}
// }

//用于Int排序(从小到大)
type IntSlice []int

func (s IntSlice) Len() int           { return len(s) }
func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
