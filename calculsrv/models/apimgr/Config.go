package apimgr

//API 提供商

//开采网
const (
	//--------------------------------------- 高频彩 -----------------------------------
	//-------------------------------------- EX5 --------------------------------------
	//江西11选5
	EX5_JiangXiApi = "&code=jx11x5"
	//山东11选5
	EX5_ShanDongApi = "&code=sd11x5"
	//上海11选5
	EX5_ShangHaiApi = "&code=sh11x5"
	//北京11选5
	EX5_BeiJingApi = "&code=bj11x5"
	//福建11选5
	EX5_FuJianApi = "&code=fj11x5"
	//黑龙江11选5
	EX5_HeiLongJiangApi = "&code=hlj11x5"
	//江苏11选5
	EX5_JiangSuApi = "&code=js11x5"

	//----------------------------------------- K3 --------------------------------------
	//广西快3
	K3_GuangXiApi = "&code=gxk3"
	//吉林快3
	K3_JiLinApi = "&code=jlk3"
	//安徽快3
	K3_AnHuiApi = "&code=ahk3"
	//北京快3
	K3_BeiJingApi = "&code=bjk3"
	//福建
	K3_FuJianApi = "&code=fjk3"
	//河北
	K3_HeBeiApi = "&code=hebk3"
	//上海
	K3_ShangHaiApi = "&code=shk3"

	//----------------------------------------- SSC --------------------------------------
	//重庆时时彩
	SSC_ChongQingApi = "&code=cqssc"
	//天津时时彩
	SSC_TianJinApi = "&code=tjssc"
	//新疆时时彩
	SSC_XinJiangApi = "&code=xjssc"
	//内蒙古
	SSC_NeiMengGuApi = "&code=nmgssc"
	//云南
	SSC_YunNanApi = "&code=ynssc"

	//----------------------------------------- PK10 --------------------------------------
	//北京PK10
	PK10_BeiJing = "&code=bjpk10"

	//-------------------------------------- 低频彩 -------------------------------------
	//排列3
	PL3Api = "&code=pl3"
	//香港彩
	HK6Api = "&code=hk6"

	//-------------------------------------- API结构 -------------------------------------
	//API头(在没有写好Api服务器之前,试玩服务器先暂时使用 这个高防主接口)
	ApiHead = "http://c.apiplus.net"
	//按天查询
	ByDay = "/daily.do?token=784F8AEBB9331158"
	//按最新查询
	ByNewest = "/newly.do?token=784F8AEBB9331158"
	//尾部
	Tail = "&rows=1&format=json&extend=true"

	//"http://c.apiplus.net/daily.do?token=784f8aebb9331158&code=hk6"
)
