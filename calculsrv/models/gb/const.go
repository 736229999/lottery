package gb

//机器编号,同时会存在多个计算服务器
const MachineCode = "1"

//全局使用的const 或者变量等暂时放这里面
const (
	LotteryStatus_Close    = 0
	LotteryStatus_Nonmal   = 1
	LotteryStatus_Maintain = 2
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
	//江苏
	K3_JiangSu = "K3_JiangSu"

	//重庆时时彩
	SSC_ChongQing = "SSC_ChongQing"
	//天津时时彩
	SSC_TianJin = "SSC_TianJin"
	//新疆时时彩
	SSC_XinJiang = "SSC_XinJiang"

	//北京PK10
	PK10_BeiJing = "PK10_BeiJing"

	//PCDD 北京28
	BJ28 = "BJ28"

	//-------------------------------------- 低频彩 -------------------------------------
	//排列3
	PL3 = "PL3"
	//香港彩
	HK6 = "HK6"

	//-------------------------------------- 独立彩票 --------------------------------------
	//急速pk10
	PK10_F = "PK10_F"
)

//默认历史记录条数
const DefaultHistoryNum = 100

//Cal 向 Game 用公钥
var RsaCalToGamePublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDbKQCgjc8S4AnrxW2AmrnZ1lGv
Sf6me64mPsDy0ZOsluFmEOh1ul4GzzuP046gzsF2/VMPMeK7EpOy5nfik3khZ/DR
hy1pl9CI+6hQO++4P9vZgEkorJji05CXy/l8pOb3G7E5zob/YMwZya00yLeC8U7N
DcTGyA9LrUeNPGKsJwIDAQAB
-----END PUBLIC KEY-----
`)
