package GlobalData

import "time"

//RSA 私钥 (注意 `` 不是 "")
var RsaGameToLoginPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDsQx8pbtf2qsj0a7Y8qCHJ6uYiNmoPA2ZXRSKOw0mZqkIJTJ3M
Y3c74/XhMBWc1bsNMBvfKH+w+BCpSRTbrvXQpF9B5Ks/vRcu192AOCnlFay3FJ6r
Dk+Zt/GyE1q75+mQIthvbiJY6IEA6kZ1isHw+2nj27M0slwlWmPIoD8xnwIDAQAB
AoGAYiwB5tWIJ5cPqYCYWNwELkLNHao+p75h6CFyGqBLFO8KQZL0ftlV0i0HHms8
j86/ytsSuce6R27MfYtAf4hT23IEXl1z7T32pF15mOD+gjIq/1bitXhvg5qw9lTd
I68EQjdj7BzJrB7lZ0eGtVB+NDF0we5NRtimGkyH228HlgECQQD7Bmwbpc4wt4Pl
2APvsUD25gHaBXbZOv7tbUIWcZEdVhO2uDGDe03d3tQVygTMfoIAtQNJUsydF6QQ
/7SBZGi3AkEA8PHMvlINgc7zQSgrtqc9FQaeWUNadRvck0YB1QTpswmMGyWfrWF4
KNoR72261Ao7sTjBjuyxIJO+erSfybGGWQJAXHSPGNfGy7arw9n+CShV7xnkELL2
hSfvgO8+6hlGA3ISrLfGtNCTa2aI4sWXYuIta6k/3p+6cRml+gTULNwYnQJBAN8s
KKfUWpx0ss4URkEdsK8r/TnS8BNu5iUngATuUwS5gCOY+sjecizsqBYsfhNvExE4
79fRltME5jbD3Qk1vKkCQH9XzQrOK0Bl0+PRS8pSqLObNtK1bn9eAsiAFpYWv9Nk
HLbL18iTr511FgLcuhk3DRKBt4XOsgDt5NVSY/R2R04=
-----END RSA PRIVATE KEY-----
`)

//User 用户注册信息
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
