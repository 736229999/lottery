package apimgr

import (
	"apisrv/conf"
	"apisrv/models/ctrl"
	"common/httpmgr"
	"common/utils"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

//api服务器负责维护api接口以及,提供给使用类需要的功能函数,数据结构定义再自身类中

//目前所有的API还是写死在conf里面的,等完成功能,来将所有的 api 还是放在控制服,去控制服中获取
//目前只有开采网一家api,写完功能后来增加 第二家开采api
type ApiMgr struct {
	ApiNew [1]map[string][]string // map[游戏名字][]api   最后一个api数组为,0主接口 1备用接口, 2vip接口(按最新查询) (0,为开采网)
	ApiDay [1]map[string][]string //按天(日期)查询(注意使用的时候要在这个url后面加上日期格式) (0,为开采网)
}

var sInstance *ApiMgr
var once sync.Once

func Instance() *ApiMgr {
	once.Do(func() {
		sInstance = &ApiMgr{}
	})
	return sInstance
}

//加新彩种,加新的api,加新的api提供商 都要再这里加
func (o *ApiMgr) Init() error {

	//开采网的API map
	var kcwApiNew = make(map[string][]string)
	var kcwApiDay = make(map[string][]string)
	//这里可以添加其他API供应商的API map

	o.ApiNew[0] = kcwApiNew
	o.ApiDay[0] = kcwApiDay

	for _, v := range ctrl.Instance().Ltrys {
		switch v.Game_name {
		case "EX5_JiangXi":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_JiangXiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_JiangXiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_JiangXiApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_JiangXiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_JiangXiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_JiangXiApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay
		case "EX5_ShanDong":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_ShanDongApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_ShanDongApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_ShanDongApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_ShanDongApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_ShanDongApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_ShanDongApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "EX5_ShangHai":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_ShangHaiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_ShangHaiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_ShangHaiApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_ShangHaiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_ShangHaiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_ShangHaiApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "EX5_BeiJing":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_BeiJingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_BeiJingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_BeiJingApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_BeiJingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_BeiJingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_BeiJingApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "EX5_FuJian":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_FuJianApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_FuJianApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_FuJianApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_FuJianApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_FuJianApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_FuJianApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "EX5_HeiLongJiang":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_HeiLongJiangApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_HeiLongJiangApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_HeiLongJiangApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_HeiLongJiangApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_HeiLongJiangApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_HeiLongJiangApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "EX5_JiangSu":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+EX5_JiangSuApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+EX5_JiangSuApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+EX5_JiangSuApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+EX5_JiangSuApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+EX5_JiangSuApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+EX5_JiangSuApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_GuangXi":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_GuangXiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_GuangXiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_GuangXiApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_GuangXiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_GuangXiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_GuangXiApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_JiLin":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_JiLinApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_JiLinApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_JiLinApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_JiLinApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_JiLinApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_JiLinApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_AnHui":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_AnHuiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_AnHuiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_AnHuiApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_AnHuiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_AnHuiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_AnHuiApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_BeiJing":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_BeiJingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_BeiJingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_BeiJingApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_BeiJingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_BeiJingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_BeiJingApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_FuJian":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_FuJianApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_FuJianApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_FuJianApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_FuJianApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_FuJianApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_FuJianApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_HeBei":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_HeBeiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_HeBeiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_HeBeiApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_HeBeiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_HeBeiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_HeBeiApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_ShangHai":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_ShangHaiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_ShangHaiApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_ShangHaiApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_ShangHaiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_ShangHaiApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_ShangHaiApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "K3_JiangSu":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+K3_JiangSu+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+K3_JiangSu+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+K3_JiangSu+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+K3_JiangSu+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+K3_JiangSu+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+K3_JiangSu+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "SSC_ChongQing":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+SSC_ChongQingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+SSC_ChongQingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+SSC_ChongQingApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+SSC_ChongQingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+SSC_ChongQingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+SSC_ChongQingApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "SSC_TianJin":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+SSC_TianJinApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+SSC_TianJinApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+SSC_TianJinApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+SSC_TianJinApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+SSC_TianJinApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+SSC_TianJinApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "SSC_XinJiang":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+SSC_XinJiangApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+SSC_XinJiangApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+SSC_XinJiangApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+SSC_XinJiangApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+SSC_XinJiangApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+SSC_XinJiangApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "PK10_BeiJing":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+PK10_BeiJingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+PK10_BeiJingApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+PK10_BeiJingApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+PK10_BeiJingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+PK10_BeiJingApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+PK10_BeiJingApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "PL3":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+PL3Api_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+PL3Api_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+PL3Api_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+PL3Api_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+PL3Api_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+PL3Api_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "HK6":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+HK6Api_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+HK6Api_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+HK6Api_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+HK6Api_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+HK6Api_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+HK6Api_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "BJ28":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+BJ28_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+BJ28_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+BJ28_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+BJ28_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+BJ28_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+BJ28_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "DLT":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+DLTApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+DLTApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+DLTApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+DLTApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+DLTApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+DLTApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "FC3D":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+FC3DApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+FC3DApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+FC3DApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+FC3DApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+FC3DApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+FC3DApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "QXC":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+QXCApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+QXCApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+QXCApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+QXCApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+QXCApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+QXCApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "SSQ":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+SSQApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+SSQApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+SSQApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+SSQApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+SSQApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+SSQApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		case "CAKENO":
			var apiArrayNew []string
			apiArrayNew = append(apiArrayNew, ApiHeadPrimary_kcw+ByNewest_kcw+CAKENOApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadSecondary_kcw+ByNewest_kcw+CAKENOApi_kcw+Tail_kcw)
			apiArrayNew = append(apiArrayNew, ApiHeadVip_kcw+ByNewest_kcw+CAKENOApi_kcw+Tail_kcw)

			kcwApiNew[v.Game_name] = apiArrayNew

			var apiArrayDay []string
			apiArrayDay = append(apiArrayDay, ApiHeadPrimary_kcw+ByDay_kcw+CAKENOApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadSecondary_kcw+ByDay_kcw+CAKENOApi_kcw+Tail_kcw_day)
			apiArrayDay = append(apiArrayDay, ApiHeadVip_kcw+ByDay_kcw+CAKENOApi_kcw+Tail_kcw_day)

			kcwApiDay[v.Game_name] = apiArrayDay

		default:
			return errors.New("There is no lottery type : " + v.Game_name)
		}
	}

	beego.Info("--- Init API Mgr Done !")
	return nil
}

//按日期获取一个彩种当日所有开采(这里有个问题,这个函数应该是提一个统一调用接口,返回的应该是统一的数据格式,而不管这个数据是从开采网还是其他什么api提供商来得,所以现在先统一格式,等以后加上其他api以后再将其他API获取的数据都以这种统一的格式返回)
//不同的彩种访问间隔时间是不一样的目前开采网按天访问是需要间隔5秒,访问带下一期的是 1秒间隔,函数负责返回api返回的数据结构
func (o *ApiMgr) GetLtryRecordByDate(gameName string, date time.Time) (LtryRecordDay, error) {
	var ret LtryRecordDay

	d := date.Format(utils.TF_D)

	//循环每一个api提供商 目前优先使用开彩票网
	for i := 0; i < len(o.ApiDay); i++ {
		//beego.Debug("--- 开始循环 api 提供商 id :", i)
		apiMap := o.ApiDay[i]

		if apiArray, ok := apiMap[gameName]; ok {
			//循环每一个不同的接口(目前优先使用开彩网的高防主接口,然后是副接口,其次是他妈的 坑爹的 V ! I ! P ! 接口)
			for j := len(apiArray) - 1; j >= 0; j-- {
				count := conf.GetRecordOneDayRetryCount
				for ; count > 0; count-- {
					url := apiArray[j] + d
					resp, err := httpmgr.Get(url)
					if err == nil {
						err := json.Unmarshal(resp, &ret)
						if err == nil {
							return ret, nil
						}
					}
					time.Sleep(conf.GetRecordByNewestSleepTime * time.Second)
				}
			}
		} else {
			return ret, errors.New("There is no API for this lottery : " + gameName)
		}
		beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", gameName)
	}

	return ret, errors.New("--- Failed to obtain data from all API providers !!! Game Name : " + gameName)
}

//得到一个彩票最新的开彩记录(循环每一个接口,每个接口重试3次)
func (o ApiMgr) GetNewRecord(gameName string) (LtryRecordNew, error) {
	var ret LtryRecordNew

	for i := 0; i < len(o.ApiNew); i++ {
		apiMap := o.ApiNew[i]

		if apiArray, ok := apiMap[gameName]; ok {
			//循环每一个不同的接口(获取最新记录和获取历史记录不一样,优先使用0号高防接口)
			for j := 0; j < len(apiArray); j++ {
				count := 3 //重试次数,这里暂时写死 以后这些全部要写到控制服数据库中去
				for ; count > 0; count-- {
					url := apiArray[j]
					resp, err := httpmgr.Get(url)
					if err == nil {
						err := json.Unmarshal(resp, &ret)
						if err == nil {
							return ret, nil
						}
					}
					beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", gameName, " API index : ", j)
					time.Sleep(conf.GetRecordByNewestSleepTime * time.Second)
				}
			}
		} else {
			return ret, errors.New("There is no API for this lottery : " + gameName)
		}
		beego.Warn("--- Failed to obtain data from the API provider,  Provider ID : ", i, " Game Name : ", gameName)
	}

	return ret, errors.New("--- Failed to obtain data from all API providers !!! Game Name : " + gameName)
}

//获取的按天查询历史记录(同样这个也是统一格式目前使用的是开采网中格式,等以后加上了其他的api提供商以后,都要将数据全部统一成这个格式)
type LtryRecordDay struct {
	Rows   int
	Code   string
	Remain string
	Data   []LtryRecord
}

//一条记录开彩记录(这个是从api获取的一条彩票开采记录,目前这个结构式开采网的,以后加入了新的api提供商以后,就都需要把记录修改为这种格式)
type LtryRecord struct {
	Expect        string
	Opencode      string
	OpenTime      string
	OpenTimeStamp int64
}

//开采网 按最新获取一条记录(带下期)
type LtryRecordNew struct {
	Rows   int
	Code   string
	Remain string
	Next   []LtryRecordNext    //下一期信息(下一期期数, 下一期开奖时间)
	Open   []LtryRecordCurrent //最新一期信息 (注意:这里面没有时间戳,存库的时候要自己添加一个)
	Time   string              //查询时间
}

//开采网 下一期开采信息
type LtryRecordNext struct {
	Expect   string
	Opentime string
}

//开采网 当前这期信息
type LtryRecordCurrent struct {
	Expect   string
	Opencode string
	Opentime string
}
