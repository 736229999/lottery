package conf

//注意 这个文件中得东西以后都要写到加密配置文件中,这样才方便修改

//公钥用于加密消息(所有Ctrl之外的服务器都使用这个加密发送消息到Ctrl)
var RsaPuk = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo+Low3DkBUwY+2Lr/FjT
xB+b2cXbTWH+KQ0aAzwChVVE/FESIfkGJ1QjBq8krzwgpSS/x0LRYF9Bx9u2Q88t
PSy6KSPcbf+SyKtw6eLt250nOSIDUZFkC7MZldMEWKSagIhOyUaBZnNNZ/WABPwy
BGgsPljMMf+GJACHRyvnE0X2iBm5y86RuLHhCBg8hXPYoj9kpUU1hBrwhEvcGw16
qgpSU1mZoUseLZY2uRKDV93URQoZ5tQ11ffYIz2jicGq/y/T63uqJR5OfVaAjqCR
3SC/5TAFCwE8nrNNV5pPdwvLBwquxWI21u4XiP3QIBf4Tdb3mEPvcU/QZmzkTDNR
oQIDAQAB
-----END PUBLIC KEY-----
`)

//RSA 私钥用于解密消息(所有 Ctrl 发送其他服务器的消息都使用这个解密)
var RsaPrk = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA5r1fp4nXNrfp0WRetwjnaFQ/B88Vws6WPo0GgGhZL0SpbYDG
RWSfP/eIlE2N1MvOb6Fe74N3lhCkHlgzQuS+zNLq72AbR1DRXRxYD59jdN2fX/ME
KsHFvvRPSTl201TtjJSor7A4w3YgJn1sSzbbddma1t8BVtqFDbKzKpqLQQscsvbK
lNi2y/BHbpNFQtpSFYctGsYwTrB/mVB/wMwKnoZkPAIiljeOrodT9Rce7Dgsc5WY
JsNXRxDM8jO+g554c8wy3euq3CUFKOSaxA8n6PpGDnIWsYt7O7ch6BhCiLoLQpsw
+GPdLbVxryp4ZOtk+rFr7/5iihNUDR+PnlZwTQIDAQABAoIBABo/SoVFYczgqOFf
2IJbqe8tPipGPUR2uZKN+kJbHGILHnbNYkB8jLz7DHdKRysAsA/0vFbkjpLse67T
+6jNWoL0LzNhrFi1cct0sPO9/tNJbpq8iynN9w+cvPQobELn80m9k17A3gQmCMw4
EjyQFfDW+w0cOwvFOcNwI39iKxsP0DQYHaqeJomkwO37x0M/SgreT1SahR8vo5Hx
no7fSMhZWGxPmOg2zOfH4NX7BAzoS5K/TyWC3PQlg5TVIAZJ/AXAsJ6HR/rdFR+I
9RFmNzk+SstkPWNQVALLjr4HU/YWKtywd1sCeEPCzdYy40zDssn15Ekw6i9dd3MS
fbqLgUUCgYEA/CJYR+y/novdPPFWEvY/IwL+3T8VZPKnchhar7o45VHNSR7GhsXg
OJQ5wfni2DflFKIIhFHdsTs/D9noBo1ZZ1dXYhoNVxMzOOMvDUZ+vDVsUhbeig+r
C+Df1WyIf52zwtsEvCD0LxXVJlnaRwN/Hvw/jd/5Cg7VqMbMWQJvNosCgYEA6kcN
pAssY1uJLt0XGwz01YKJD4JYPhjxx25x3Q7JjRX7Ue5/Cy7RnV8Oz6ASlpSVSGhF
RwmtJPbSZoNx1xf01uFv0tFmzw7tDz9I826MUBxxjZP06LYljIoA9+J5aCUWhvEJ
dRmUcrZrLS7lnuaXUOQm3J4TcjnwewJnYJrUp4cCgYB2JsYF1mypKFOhdlpmglxt
1L6IAULOTmOnNMBybqqw05eGd1SC3YFIIjW6r5Xcyryf4ZpqH07q+Z+AlxWC1IQb
yDMbtqefsVCkjNmEuA81tTcwdKUOP29hHpzlj3mbi9QsMKRUYIDs+6cp1JtUAdB1
PbGuk1FJpE/9SkOCRoDFJwKBgEPDQYLNaI1kkI1pjaFZYgfPte0yemubs3NH0s9p
04pnbUYJgd2uMRsfv5z2Y+oCGBvVbRRbDCXA7qKQKoFYgI0Wr81+nAoP+ymJ2IPw
2cziLUYSIaid5sZ7tEP+0bb540YsuduRBosXkHCFPA12DRZsp4DwiBdmAtTRoS0k
G5ZJAoGBAK1MqVNTqKAmziacvibC6eyf3SUBANE2d9DKAGmK+umerXR9KN1jNc4s
dkoMmo8dGa+rIT6RFeNuCEP99/R4AmmdQfn7Wl2NDXoPoZGGWg3pYpsl3oH3xXNv
nVmMHBlxlD7m0cq4X4PkDBHLS3uUvv0fdvETspRPRF/EMShez0kx
-----END RSA PRIVATE KEY-----
`)

//RSA 暗号
var RsaCipher = []byte("Manners Makyth man ! , 不知礼 无以立也 !")

//AES 暗号
var AesCipher = []byte("May the force be with you !")

//ctrlsrv IP(由于启动的时候必须从控制服务获取信息,所以这个只有再在每个服务器代码中写死,等项目空闲时改为加密配置文件的形式)
//本地测试
//var CtrlsrvIP = "http://192.168.1.182:8877"

//线上
var CtrlsrvIP = "http://47.52.118.161:8877"

//超时时间
var Timeout = "5s"

//本服务器功能
var SrvFunc = "api"

const (
	//按天获取历史记录重试次数
	GetRecordOneDayRetryCount = 3
	//按天获取历史记录每次间隔时间(睡眠时间) 秒
	GetRecordOneDaySleepTime = 5

	//按最新(带下期)历史记录重试次数
	GetRecordByNewestRetryCount = 3
	//按最新(带下期)记录每次间隔时间(睡眠时间) 秒
	GetRecordByNewestSleepTime = 5

	//监听时间间隔(睡眠时间)
	ListeningSleepTime = 5
)

const (
	//--------------------------------------- 高频彩 -----------------------------------
	//江西11选5
	EX5_JiangXi = "EX5_JiangXi"
	//山东11选5
	EX5_ShanDong = "EX5_ShanDong"
	//上海11选5
	EX5_ShangHai = "EX5_ShangHai"
	//北京
	EX5_BeiJing = "EX5_BeiJing"
	//福建
	EX5_FuJian = "EX5_FuJian"
	//黑龙江
	EX5_HeiLongJiang = "EX5_HeiLongJiang"
	//江苏
	EX5_JiangSu = "EX5_JiangSu"

	//广西快3
	K3_GuangXi = "K3_GuangXi"
	//吉林快3
	K3_JiLin = "K3_JiLin"
	//安徽快3
	K3_AnHui = "K3_AnHui"
	//北京
	K3_BeiJing = "K3_BeiJing"
	//福建
	K3_FuJian = "K3_FuJian"
	//河北
	K3_HeBei = "K3_HeBei"
	//上海
	K3_ShangHai = "K3_ShangHai"

	//重庆时时彩
	SSC_ChongQing = "SSC_ChongQing"
	//天津时时彩
	SSC_TianJin = "SSC_TianJin"
	//新疆时时彩
	SSC_XinJiang = "SSC_XinJiang"

	//北京PK10
	PK10_BeiJing = "PK10_BeiJing"

	//-------------------------------------- PC蛋蛋(PCDD) -------------------------------------
	//北京28
	BJ28 = "BJ28"

	//-------------------------------------- 低频彩 -------------------------------------
	//排列3
	PL3 = "PL3"
	//香港彩
	HK6 = "HK6"
)
