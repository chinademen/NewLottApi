package models

import (
	"NewLottApi/dao"

	"common"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

const (
	STATUS_NOT_AVAILABLE        = 0
	STATUS_AVAILABLE_FOR_TESTER = 1
	STATUS_AVAILABLE            = 3
	STATUS_TESTING              = 4
	STATUS_CLOSED_FOREVER       = 8
)

type tLotteries struct {
	TbName string
	Fields string
}

var Lotteries = &tLotteries{TbName: "lotteries",
	Fields: "`id`, `series_id`, `name`, `type`, `lotto_type`, `is_self`, `is_instant`, `high_frequency`, `sort_winning_number`, `valid_nums`, `buy_length`, `wn_length`, `identifier`, `days`, `issue_over_midnight`, `issue_format`, `bet_template`, `begin_time`, `end_time`, `sequence`, `status`, `need_draw`, `daily_issue_count`, `trace_issue_count`, `max_bet_group`, `series_ways`, `created_at`, `updated_at`"}

//系列名語言包
var SeriesName = map[string]string{}

//彩種名語言包
var LotZhName = map[string]string{}

func init() {
	LotZhName = map[string]string{

		"115mmc":        "11选5秒秒彩",
		"115sfc":        "11选5三分彩",
		"3d":            "3D",
		"agother":       "AG(仅用于统计)",
		"ah115":         "安徽11选5",
		"ahk3":          "安徽快三",
		"all-lotteries": "全部彩种",
		"bj115":         "北京11选5",
		"bj5fc":         "北京五分彩",
		"bjkl8":         "北京快乐8",
		"bjlgjt":        "百家乐",
		"bjlqjt":        "百家乐旗舰厅",
		"bjpk10":        "北京PK10",
		"cq11-5":        "重庆11选5",
		"cqklsf":        "重庆快乐十分",
		"cqssc":         "重庆时时彩",
		"dzyx":          "老虎机",
		"dzyxdt":        "电游独立大厅",
		"fj115":         "福建11选5",
		"fjk3":          "福建快三",
		"gabjl":         "GA百家乐",
		"gabrnn":        "GA百人牛牛",
		"gahjk":         "GA二十一点",
		"gajstb":        "GA江苏骰宝",
		"galh":          "GA龙虎",
		"gaqwnn":        "GA趣味牛牛",
		"gasb":          "GA骰宝",
		"gasgj":         "GA水果机",
		"gasgsgj":       "GA三国水果机",
		"gaszpk":        "GA赢三张",
		"gaxszc":        "GA西施早餐",
		"gaxydzp":       "GA幸运大转盘",
		"gbgjt":         "骰宝国际厅",
		"gbqjt":         "骰宝",
		"gd115":         "广东11选5",
		"gdklsf":        "广东快乐十分",
		"gs115":         "甘肃11选5",
		"gsk3":          "甘肃快三",
		"gxk3":          "广西快三",
		"gz115":         "贵州11选5",
		"hb115":         "河北11选5",
		"hebk3":         "河北快三",
		"hgklc":         "韩国快乐彩",
		"hgklssc":       "韩国1.5分彩",
		"hljssc":        "黑龙江时时彩",
		"hnk3":          "河南快三",
		"jczq":          "竞彩足球",
		"jdly":          "机动乐园",
		"jlk3":          "吉林快三",
		"js115":         "江苏11选5",
		"jsk3":          "江苏快三",
		"jx115":         "江西11选5",
		"jxssc":         "江西时时彩",
		"lhgjt":         "龙虎国际厅",
		"lhqjt":         "龙虎",
		"ln115":         "辽宁11选5",
		"lnkl12":        "辽宁快乐十二",
		"lpgjt":         "轮盘",
		"lpqjt":         "轮盘旗舰厅",
		"mdbmc":         "缅甸百秒彩",
		"nmg115":        "内蒙古11选5",
		"pg115":         "苹果11选5",
		"pg3d":          "苹果极速3D",
		"pg3fc":         "苹果三分彩",
		"pg5fc":         "苹果五分彩",
		"pgffc":         "苹果分分彩",
		"pgk3ffc":       "苹果快三分分彩",
		"pgkeno":        "苹果快乐8分分彩",
		"pgmmc":         "苹果秒秒彩",
		"pgpk10":        "苹果极速PK10",
		"pl5":           "排列三/五",
		"rbws":          "日本武士",
		"saxklsf":       "陕西快乐十分",
		"sckl12":        "四川快乐十二",
		"sd115":         "山东11选5",
		"sglb":          "水果拉霸",
		"sh115":         "上海11选5",
		"ssl":           "上海时时乐",
		"sx115":         "山西11选5",
		"sxklsf":        "山西快乐十分",
		"tj115":         "天津11选5",
		"tjssc":         "天津时时彩",
		"twbingo":       "台湾宾果",
		"xjssc":         "新疆时时彩",
		"xryd":          "夏日营地",
		"xylhj":         "幸运老虎机",
		"ynssc":         "云南时时彩",
		"zj115":         "浙江11选5",
		"zjkl12":        "浙江快乐十二",
		"zrbyw":         "捕鱼王",
		"zrbzt":         "真人包桌(VIP包桌)",
		"zrdss":         "赌神赛",
		"zrgjt":         "真人国际厅",
		"zrjmt":         "竞咪(竞咪厅)",
		"zrqjt":         "真人旗舰厅",
		"zryxdt":        "真人游戏大厅",
		"zrzbt":         "真人直播厅",
	}

	SeriesName = map[string]string{
		"SSC":  "时时彩",
		"11-5": "十一选5",
		"3D":   "3D",
		"快三":   "快三",
		"PK10": "PK10",
		"KENO": "基诺",
		"KL12": "快乐十二",
		"KLSF": "快乐十分",
	}
}

/*
 * 根据Ｉd获取彩种数据
 */
func (m *tLotteries) GetInfo(sLotteryId string) *dao.Lotteries {

	//從單一緩存key讀取
	sCahceKey := CompileOneLotteryCacheKey(sLotteryId)
	mCahceData := redisClient.Redis.HashReadAllMap(sCahceKey)
	iId, _ := strconv.Atoi(sLotteryId)
	var oLottery *dao.Lotteries
	if len(mCahceData) < 1 {
		oLottery, _ = dao.GetLotteriesById(iId)
		if oLottery.Id != 0 {
			m := map[string]string{}
			b, _ := json.Marshal(oLottery)
			json.Unmarshal(b, &m)
			redisClient.Redis.HashWrite(sCahceKey, m, 60)
		}
	} else {
		cahceData, _ := json.Marshal(mCahceData)
		json.Unmarshal(cahceData, &oLottery)
	}

	if oLottery.Id > 0 {
		return oLottery
	}

	//從所有緩存讀取
	aLotteries := m.GetAll()
	for _, oValue := range aLotteries {
		if oValue.Id == iId {
			return oValue
		}
	}

	return oLottery
}

/*
 * 根据sIdentifier获取彩种数据
 */
func (m *tLotteries) GetInfoByIdentifier(sIdentifier string) *dao.Lotteries {
	aLotteries := m.GetAll()
	for _, oValue := range aLotteries {
		if oValue.Identifier == sIdentifier {
			return oValue
		}
	}
	var empty *dao.Lotteries
	return empty

}

/*
 * 获取接入商不关闭的彩种列表
 * params []string 對接入商關閉的id
 */
func (c *tLotteries) GetNeedLotteriesData(aCloseIds []string) []orm.ParamsList {
	aAll := c.GetAll()
	var aAllIds []string
	for _, oLottery := range aAll {
		aAllIds = append(aAllIds, fmt.Sprintf("%d", oLottery.Id))
	}
	aId := common.ArrayDiff(aAllIds, aCloseIds)

	sField := "id,name,identifier,series_id,is_instant,status"
	sSql := fmt.Sprintf("select %s from %s where id in (%s) and status = 3", sField, c.TbName, strings.Join(aId, ","))

	o := orm.NewOrm()
	var aNeed []orm.ParamsList
	o.Raw(sSql).ValuesList(&aNeed)
	return aNeed
}

/*
 * 获取所有彩种信息
 */
func (c *tLotteries) GetAll() []*dao.Lotteries {

	sCacheKey := CompileAllLotteryCacheKey()
	iCount := redisClient.Redis.LenList(sCacheKey)
	var aAll []*dao.Lotteries
	if iCount > 0 { //緩存讀取
		aList := redisClient.Redis.LRANGE(sCacheKey, "0", "-1")
		for _, sJson := range aList {
			obj := new(dao.Lotteries)
			json.Unmarshal([]byte(sJson), &obj)
			aAll = append(aAll, obj)
		}
	} else { //數據庫讀取
		aAll, _ = dao.GetAllLotteries(nil, nil, nil, nil, 0, 1000)
		for _, oLottery := range aAll {

			lottery, _ := json.Marshal(oLottery)
			redisClient.Redis.RPushList(sCacheKey, string(lottery))
			redisClient.Redis.KeyExpire(sCacheKey, 3600, 1)
		}
	}
	return aAll
}

func (m *tLotteries) GetList(sWhere, sField, sOrder string, offset, limit int) []map[string]string {

	if len(sField) == 0 {
		sField = m.Fields
	}

	rMap := GetList(m.TbName, sWhere, sField, sOrder, offset, limit)

	return rMap
}

/*
得到 type=1 or type =2 的彩种信息
*/
func GetGroupPrizeLottery() []map[string]string {

	sWhere := "type = 1 OR type = 2"
	oLotteries := Lotteries.GetList(sWhere, "", "", 0, 1000)

	return oLotteries
}
