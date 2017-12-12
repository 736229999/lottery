package apimgr

//API 提供商 开采网

//开采网
const (
	//--------------------------------------- 高频彩 -----------------------------------
	//-------------------------------------- EX5 --------------------------------------
	//江西11选5
	EX5_JiangXiApi_kcw = "&code=jx11x5"
	//山东11选5
	EX5_ShanDongApi_kcw = "&code=sd11x5"
	//上海11选5
	EX5_ShangHaiApi_kcw = "&code=sh11x5"
	//北京11选5
	EX5_BeiJingApi_kcw = "&code=bj11x5"
	//福建11选5
	EX5_FuJianApi_kcw = "&code=fj11x5"
	//黑龙江11选5
	EX5_HeiLongJiangApi_kcw = "&code=hlj11x5"
	//江苏11选5
	EX5_JiangSuApi_kcw = "&code=js11x5"

	//----------------------------------------- K3 --------------------------------------
	//广西快3
	K3_GuangXiApi_kcw = "&code=gxk3"
	//吉林快3
	K3_JiLinApi_kcw = "&code=jlk3"
	//安徽快3
	K3_AnHuiApi_kcw = "&code=ahk3"
	//北京快3
	K3_BeiJingApi_kcw = "&code=bjk3"
	//福建快3
	K3_FuJianApi_kcw = "&code=fjk3"
	//河北快3
	K3_HeBeiApi_kcw = "&code=hebk3"
	//上海快3
	K3_ShangHaiApi_kcw = "&code=shk3"
	//江苏快三
	K3_JiangSu = "&code=jsk3"

	//----------------------------------------- SSC --------------------------------------
	//重庆时时彩
	SSC_ChongQingApi_kcw = "&code=cqssc"
	//天津时时彩
	SSC_TianJinApi_kcw = "&code=tjssc"
	//新疆时时彩
	SSC_XinJiangApi_kcw = "&code=xjssc"
	//内蒙古时时彩
	SSC_NeiMengGuApi_kcw = "&code=nmgssc"
	//云南时时彩
	SSC_YunNanApi_kcw = "&code=ynssc"

	//----------------------------------------- PK10 --------------------------------------
	//北京PK10
	PK10_BeiJingApi_kcw = "&code=bjpk10"

	//----------------------------------------- PC蛋蛋 --------------------------------------
	//北京28(这个菜种有点奇葩,大类名叫PC蛋蛋,彩种叫北京28,开采网api名字却是 bjkl8(北京快乐8))
	BJ28_kcw = "&code=bjkl8"

	//加拿大3.5(这他妈也是个奇葩的名字)
	CAKENOApi_kcw = "&code=cakeno"

	//-------------------------------------- 低频彩 -------------------------------------
	//排列3
	PL3Api_kcw = "&code=pl3"
	//香港六合彩
	HK6Api_kcw = "&code=hk6"
	//大乐透
	DLTApi_kcw = "&code=dlt"
	//福彩3D
	FC3DApi_kcw = "&code=fc3d"
	//七星彩
	QXCApi_kcw = "&code=qxc"
	//双色球
	SSQApi_kcw = "&code=ssq"

	//-------------------------------------- API结构 -------------------------------------
	//API头(高防主接口)(在没有写好Api服务器之前,试玩服务器先暂时使用 这个高防主接口)
	ApiHeadPrimary_kcw = "http://c.apiplus.net"
	//API头(备用接口)
	ApiHeadSecondary_kcw = "http://z.apiplus.net"
	//API头(坑爹VIP接口)
	ApiHeadVip_kcw = "http://101.37.126.3:7825"

	//按天查询
	ByDay_kcw = "/daily.do?token=784F8AEBB9331158"
	//按最新查询
	ByNewest_kcw = "/newly.do?token=784F8AEBB9331158"
	//尾部
	Tail_kcw = "&rows=1&format=json&extend=true"
	//按天查询尾部,使用的时候还要在这个尾部后面加上日期
	Tail_kcw_day = "&format=json&date="

	//"http://c.apiplus.net/daily.do?token=784f8aebb9331158&code=hk6"
)
