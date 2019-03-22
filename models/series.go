package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
	"strconv"
)

type tSeries struct {
	TbName string
	Fields string
}

var Series = &tSeries{TbName: "series",
	Fields: "`id`, `type`, `lotto_type`, `name`, `identifier`, `sort_winning_number`, `buy_length`, `wn_length`, `digital_count`, `classic_amount`, `group_type`, `max_percent_group`, `max_prize_group`, `max_real_group`, `max_bet_group`, `valid_nums`, `offical_prize_rate`, `default_way_id`, `link_to`, `lotteries`, `bonus_enabled`"}

var RSeriesKey string = "series:"                   //redis基本key
var RSeriesRowKey string = RSeriesKey + "row:id_%s" //redis 数据库一行key

var AViewTypes = map[string]string{
	"8":  "ba-xing",  // ['wu-xing', '八星'],
	"5":  "wu-xing",  // ['wu-xing', '五星'],
	"4":  "si-xing",  // ['si-xing', '四星'],
	"3":  "san-xing", // ['san-xing', '三星'],
	"3f": "qian-san", // ['qian-san', '前三'],
	"3e": "hou-san",  // ['hou-san', '后三'],
	"2f": "qian-er",  // ['qian-er', '前二'],
	"2e": "hou-er",   // ['hou-er', '后二'],
	"10": "cmc",
	"20": "keno",
}

/*
 *根据Ｉd获取系列数据
 */
func (m *tSeries) GetInfo(sId string) *dao.Series {

	//從單一緩存key讀取
	sCahceKey := CompileOneSeriesCacheKey(sId)
	mCahceData := redisClient.Redis.HashReadAllMap(sCahceKey)
	iId, _ := strconv.Atoi(sId)
	var oSeries *dao.Series
	if len(mCahceData) < 1 {
		oSeries, _ = dao.GetSeriesById(iId)
		if oSeries.Id != 0 {
			m := map[string]string{}
			b, _ := json.Marshal(oSeries)
			json.Unmarshal(b, &m)
			redisClient.Redis.HashWrite(sCahceKey, m, 60)
		}
	} else {
		cahceData, _ := json.Marshal(mCahceData)
		json.Unmarshal(cahceData, &oSeries)
	}

	if oSeries.Id > 0 {
		return oSeries
	}

	//從所有緩存讀取
	aSeries := m.GetAll()
	for _, oValue := range aSeries {
		if oValue.Id == iId {
			return oValue
		}
	}

	return oSeries
}

/*
 *根据条件获取所有系列
 */
func (m *tSeries) GetAll() []*dao.Series {
	sCacheKey := CompileAllLotteryCacheKey()
	sCacheData := redisClient.Redis.StringRead(sCacheKey)
	var aAll []*dao.Series
	if len(sCacheData) > 0 {
		json.Unmarshal([]byte(sCacheData), &aAll)
	} else {
		aAll, _ = dao.GetAllSeries(nil, nil, nil, nil, 0, 1000)
		sAllbyte, _ := json.Marshal(aAll)
		redisClient.Redis.StringWrite(sCacheKey, string(sAllbyte), -1)
	}

	return aAll
}

/*
 * 獲取真正的系列id
 */
func (m *tSeries) GetRealSeriesId(iSeriesId string) string {
	obj := m.RGetOneById(iSeriesId)
	if len(obj["link_to"]) > 0 {
		return obj["link_to"]
	}

	return iSeriesId
}

/*
 查询资料
*/
func (m *tSeries) GetOne(sWhere, sFields string) map[string]string {

	if len(sFields) == 0 {
		sFields = m.Fields
	}

	//从数据库读取结果
	rMap := GetOne(m.TbName, sWhere, sFields)

	return rMap
}

/*
 查询资料 redis
*/
func (m *tSeries) RGetOneById(iSeriesId string) map[string]string {

	rKey := fmt.Sprintf(RSeriesRowKey, iSeriesId)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	sWhere := fmt.Sprintf("id = '%s'", iSeriesId)
	rMap = m.GetOne(sWhere, "")

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RUsersKeyEX)
	}

	return rMap
}

//获取彩种趋势类型
func GetTrendType(seriesId string) []string {

	arr := []string{}
	if seriesId == "1" {
		arr = []string{"5", "4", "3f", "3e", "2f", "2e"}
	}

	if seriesId == "2" {
		arr = []string{"5"}
	}

	if seriesId == "3" {
		arr = []string{"3"}
	}

	if seriesId == "4" {
		arr = []string{"3"}
	}

	if seriesId == "5" {
		arr = []string{"10"}
	}

	if seriesId == "7" {
		arr = []string{"20"}
	}

	if seriesId == "8" {
		arr = []string{"5"}
	}

	if seriesId == "9" {
		arr = []string{"8"}
	}

	return arr
}
