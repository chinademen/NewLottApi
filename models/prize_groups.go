package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

const (
	PrizeGroupsRedisKey = "prize_groups_redis_%v"
)

var PrizeGroupsList []*dao.PrizeGroups

type tPrizeGroups struct {
	TbName string
	Fields string
}

var PrizeGroups = &tPrizeGroups{TbName: "prize_groups", //彩种玩法分组表
	Fields: "`id`, `series_id`, `type`, `name`, `classic_prize`, `water`, `created_at`, `updated_at`"}

var RPrizeGroupsKey string = "prize_groups:"                                 //redis基本key
var RPrizeGroupsRowKey string = RPrizeGroupsKey + "row:name_%s"              //redis 数据库一行key
var RPrizeGroupsRowKey1 string = RPrizeGroupsKey + "row:name_%s-seriesId_%s" //redis 数据库一行key
var RPrizeGroupsKeyEX int = 60 * 60 * 24

//缓存奖金组数据
func init() {
	query := make(map[string]string)
	var err error
	PrizeGroupsList, err = dao.GetAllPrizeGroups(query, nil, nil, nil, 0, 100000)
	if err != nil {
		logs.Error("PrizeGroups init Error:%s", err)
	}
}

//GetPrizeGroupsByID GetPrizeGroupsByID
func GetPrizeGroupsByID(groupid int) *dao.PrizeGroups {
	for _, prizegroup := range PrizeGroupsList {
		if prizegroup.Id == groupid {
			return prizegroup
		}
	}
	return nil
}

//GetPrizeGroupsByName GetPrizeGroupsByName
func GetPrizeGroupsByName(name string) *dao.PrizeGroups {
	for _, prizegroup := range PrizeGroupsList {
		if prizegroup.Name == name {
			return prizegroup
		}
	}
	return nil
}

/**
 * 获得用户奖金设置
 * @param string sUserId
 * @param string sLotteryId
 * @return map[string]map[string]interface{}
 */
func (m *tPrizeGroups) GetPrizeSettingsOfUser(sUserId, sLotteryId string) (map[string]map[string]interface{}, string) {
	sGroupId, sGroupName := UserPrizeSets.GetGroupId(sUserId, sLotteryId)
	if sGroupId < 1 {
		return nil, ""
	}
	return m.GetPrizeDetails(fmt.Sprintf("%d", sGroupId)), sGroupName
}

/**
 * 获得奖金设置详情数组
 * @param string sGroupId
 * @return map[string]map[string]interface{}
 */
func (m *tPrizeGroups) GetPrizeDetails(sGroupId string) map[string]map[string]interface{} {
	return PrizeDetails.GetDetailsNew(sGroupId)
}

/*
 * 獲取傭金詳情
 * @param *dao.Series oSeries
 * @param int iMaxGroup
 * @param int iMinGroup
 * @return []map[string]interface{}
 */
func (m *tPrizeGroups) GetPrizeCommissions(oSeries *dao.Series, iMaxGroup, iMinGroup int) []map[string]interface{} {
	iSeriesId := oSeries.Id
	if oSeries.LinkTo > 0 {
		iSeriesId = int(oSeries.LinkTo)
	}
	aPrizeGroups := m.GetPrizeGroups(iSeriesId, iMaxGroup, iMinGroup)
	var aSettings []map[string]interface{}
	for _, oGroups := range aPrizeGroups {
		rate := (float64(iMaxGroup) - float64(oGroups.ClassicPrize)) / 2000
		mData := map[string]interface{}{
			"prize_group": oGroups.Name,
			"rate":        rate,
		}
		aSettings = append(aSettings, mData)
	}
	return aSettings
}

/*
 * 獲取獎金組
 * @param int iSeriesId
 * @param int iMaxGroup
 * @param int iMinGroup
 * @return []dao.PrizeGroups
 */
func (m *tPrizeGroups) GetPrizeGroups(iSeriesId, iMaxGroup, iMinGroup int) []*dao.PrizeGroups {
	aGroups := []*dao.PrizeGroups{}

	//cahce
	sCacheKey := CompileCacheKeyOfOpSetting(iSeriesId)
	sCacheData := redisClient.Redis.StringRead(sCacheKey)
	if len(sCacheData) < 1 {
		aAllPrizeGroups := m.GetAllPrizeGroups(iSeriesId)
		for _, oPrizeGroup := range aAllPrizeGroups {
			if oPrizeGroup.ClassicPrize < uint32(iMinGroup) || oPrizeGroup.ClassicPrize > uint32(iMaxGroup) {
				continue
			}
			aGroups = append(aGroups, oPrizeGroup)
		}

		b, _ := json.Marshal(aGroups)
		redisClient.Redis.StringWrite(sCacheKey, string(b), -1)
	} else {
		json.Unmarshal([]byte(sCacheData), &aGroups)
	}

	return aGroups
}

/*
 * 獲取所有獎金組
 */
func (m *tPrizeGroups) GetAllPrizeGroups(iSeriesId int) []*dao.PrizeGroups {
	mConditions := map[string]string{
		"series_id": fmt.Sprintf("%d", iSeriesId),
	}
	aSortBy := []string{"name"}
	aOrderBy := []string{"asc"}
	aPrizeGroups, _ := dao.GetAllPrizeGroups(mConditions, nil, aSortBy, aOrderBy, 0, 1000)

	return aPrizeGroups
}

/*
 查询资料
*/
func (m *tPrizeGroups) GetOne(sWhere, sFields string) map[string]string {

	if len(sFields) == 0 {
		sFields = m.Fields
	}

	//从数据库读取结果
	rMap := GetOne(m.TbName, sWhere, sFields)

	return rMap
}

/*
 查询资料	redis
name 	奖金组name
*/
func (m *tPrizeGroups) RGetOneByName(name string) map[string]string {

	rKey := fmt.Sprintf(RPrizeGroupsRowKey, name)

	sWhere := fmt.Sprintf("name = '%s'", name)

	rMap := m.RGetOneByKeyWhere(rKey, sWhere)

	return rMap
}

/**
 * 根据奖金值获取奖金组详情]
 * @param  [string]  iClassicPrize [经典奖金值]
 * @param  [string]  iLotteryType  [彩种类型]
 * @return [dao.PrizeGroups]                 [奖金组详情]
 */
func (m *tPrizeGroups) GetPrizeGroupByClassicPrize(sClassicPrize, sSeriesId string) *dao.PrizeGroups {
	sCacheKey := CompilePrizeGroupCacheKey(sClassicPrize, sSeriesId)
	sData := redisClient.Redis.StringRead(sCacheKey)
	var sResult *dao.PrizeGroups
	if len(sData) < 1 {
		mCondition := map[string]string{
			"series_id":     sSeriesId,
			"classic_prize": sClassicPrize,
		}
		aGroup, _ := dao.GetAllPrizeGroups(mCondition, nil, nil, nil, 0, 1)
		if len(aGroup) < 1 {
			return sResult
		}
		sResult = aGroup[0]
		b, _ := json.Marshal(sResult)
		redisClient.Redis.StringWrite(sCacheKey, string(b), -1)

	} else {
		json.Unmarshal([]byte(sData), &sResult)
	}

	return sResult
}

/**
 * 获取用户奖金组ID
 * @param int $iSeriesId
 * @param string $sGroupName        用以保存奖金组名称
 * @return int | false
 */
func (m *tPrizeGroups) GetGroupId(iSeriesId, sGroupName string) string {
	oGroup := m.getObjectByName(iSeriesId, sGroupName)

	gId := ""
	if len(oGroup["id"]) > 0 {
		gId = oGroup["id"]
	}
	return gId
}

func (m *tPrizeGroups) getObjectByName(iSeriesId, sGroupName string) map[string]string {
	iSeriesId = Series.GetRealSeriesId(iSeriesId)

	rKey := fmt.Sprintf(RPrizeGroupsRowKey1, sGroupName, iSeriesId)
	sWhere := fmt.Sprintf("series_id = '%s' and name='%s'", iSeriesId, sGroupName)

	oGroup := m.RGetOneByKeyWhere(rKey, sWhere)
	return oGroup
}

/*
 查询资料	redis
name 	奖金组name
*/
func (m *tPrizeGroups) RGetOneByKeyWhere(rKey, sWhere string) map[string]string {

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetOne(sWhere, "")

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RPrizeGroupsKeyEX)
	}

	return rMap
}
