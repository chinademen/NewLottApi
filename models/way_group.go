package models

import (
	"NewLottApi/dao"
	"common"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	RedisKeyWayGroups = "way_groups"
)

type tWayGroups struct {
	Table
	id string
}

type WayGroupData struct {
	Id           int             `json:"id"`
	Pid          uint            `json:"pid"`
	SeriesWayId  uint            `json:"series_way_id"`
	Price        uint16          `json:"price"`
	BetNote      string          `json:"bet_note"`
	BonusNote    string          `json:"bonus_note"`
	BasicMethods string          `json:"basic_methods"`
	NameCn       string          `json:"name_cn"`
	NameEn       string          `json:"name_en"`
	Children     []*WayGroupData `json:"children"`
}

type WayPrizeData struct {
	name          string
	prize         string
	max_multiple  string
	display_prize string
}

var WayGroups = &tWayGroups{Table: Table{TableName: "way_groups"}} //彩种玩法分组表
var NoDyMostPrizeSeriesWayIDs = []int{175, 157, 109, 110}          //不显示最高奖金系列投注方式ids

/*
 * 選號盤三層數組(玩法群－玩法－投注方式)
 */
func (m *tWayGroups) GetWayGroupSettings(iSeriesId, iTerminalId, iLotteryId int, aPrizeDetails map[string]map[string]interface{}, sPrizeLimit, sUserId string) (interface{}, map[int]map[string]interface{}) {
	return m.GetWayGroups(iSeriesId, iTerminalId, iLotteryId, true, aPrizeDetails, sPrizeLimit, sUserId)
}

/**
 * 获取彩种的玩法
 * @param      iSeriesId
 * @param bool bForBet
 * @param int  iTerminalId
 * @param int  iLotteryId
 * @param string  sPrizeLimit
 * @param map[string]map[string]interface{}  aPrizeDetails
 *
 * @return string,map[int]map[string]string
 */
func (m *tWayGroups) GetWayGroups(iSeriesId, iTerminalId, iLotteryId int, bForBet bool, aPrizeDetails map[string]map[string]interface{}, sPrizeLimit, sUserId string) (interface{}, map[int]map[string]interface{}) {
	sTerminalId := fmt.Sprintf("%d", iTerminalId)
	sSeriesId := fmt.Sprintf("%d", iSeriesId)
	sLotteryId := fmt.Sprintf("%d", iLotteryId)

	aPrizeSettings := map[int]map[string]interface{}{} //奖金设置拼装数组

	//緩存key
	sCacheKey := MakeCacheKeyOfLotterySelectionPlate(sLotteryId, sTerminalId, bForBet) //緩存選好盤
	sCacheKeyPrizeSetting := CompileCacheKeyOfUserPrizeSetting(sLotteryId, sUserId)

	//从缓存读取数据
	sResult := redisClient.Redis.StringRead(sCacheKey)
	sPrizeData := redisClient.Redis.StringRead(sCacheKeyPrizeSetting)

	//从数据库读取数据
	aArr1 := []*WayGroupData{}
	if len(sResult) < 1 || len(sPrizeData) < 1 {
		mConditions := map[string]string{
			"series_id":         sSeriesId,
			"terminal_id":       sTerminalId,
			"parent_id__isnull": "true",
			"for_display":       "1",
		}
		aSortBy := []string{"sequence"}
		aOrder := []string{"asc"}
		aMainGroups, _ := dao.GetAllWayGroups(mConditions, nil, aSortBy, aOrder, 0, 10000) //玩法群
		for _, oMainGroup := range aMainGroups {
			aData := new(WayGroupData)
			aData.Id = oMainGroup.Id
			aData.Pid = oMainGroup.ParentId
			aData.NameCn = oMainGroup.Title
			aData.NameEn = oMainGroup.EnTitle

			mCon := map[string]string{
				"series_id":   sSeriesId,
				"terminal_id": sTerminalId,
				"parent_id":   fmt.Sprintf("%d", oMainGroup.Id),
				"for_display": "1",
			}

			aSubGroups, _ := dao.GetAllWayGroups(mCon, nil, aSortBy, aOrder, 0, 10000) //玩法组
			aArr2 := []*WayGroupData{}
			for _, oSubGroup := range aSubGroups {
				aSubData := new(WayGroupData)
				aSubData.Id = oSubGroup.Id
				aSubData.Pid = oSubGroup.ParentId
				aSubData.NameCn = oSubGroup.Title
				aSubData.NameEn = oSubGroup.EnTitle

				aWays := m.GetWays(bForBet, fmt.Sprintf("%d", oSubGroup.Id)) //玩法
				aArr3 := []*WayGroupData{}
				for _, oWay := range aWays {
					oSeriesWay, _ := dao.GetSeriesWaysById(int(oWay.SeriesWayId))

					//组装基础玩法数据
					aWayData := new(WayGroupData)
					aWayData.Id = int(oWay.SeriesWayId)
					aWayData.Pid = oWay.GroupId
					aWayData.SeriesWayId = oWay.SeriesWayId
					aWayData.NameCn = oWay.Title
					aWayData.Price = oSeriesWay.Price
					aWayData.NameEn = oWay.EnTitle
					aWayData.BetNote = oSeriesWay.BetNote
					aWayData.BonusNote = oSeriesWay.BonusNote
					aWayData.BasicMethods = oSeriesWay.BasicMethods

					//组装奖金设置数据
					aBasicMethodIds := strings.Split(oSeriesWay.BasicMethods, ",")

					var aWayPrizes, aWayPrizesMin []string
					var sPrize string
					var iMaxMultiple float64
					for _, sBasicMethodId := range aBasicMethodIds {
						mInterface := aPrizeDetails[sBasicMethodId]

						//防止斷言錯誤
						var mPrize = map[string]interface{}{}
						b, _ := json.Marshal(mInterface["level"])
						json.Unmarshal(b, &mPrize)

						var mPrizeDetail = map[string]float64{}
						for sKey, v := range mPrize {
							i, _ := strconv.ParseFloat(common.InterfaceToString(v), 64)
							mPrizeDetail[sKey] = i
						}

						aWayPrizes = append(aWayPrizes, fmt.Sprintf("%f", common.MapMaxF(mPrizeDetail)))       //最大獎金
						aWayPrizesMin = append(aWayPrizesMin, fmt.Sprintf("%f", common.MapMinF(mPrizeDetail))) //最小獎金
					}
					fPrizeLimit, _ := strconv.ParseFloat(sPrizeLimit, 64)

					fMaxPrize := common.ArrayMaxF(aWayPrizes)
					if len(sPrizeLimit) > 0 {
						fPrize := fMaxPrize
						sPrize = fmt.Sprintf("%f", fPrize)
						if fPrize > 0 {
							iMaxMultiple = math.Floor(fPrizeLimit / fPrize)
						}

					} else {
						sPrize = strings.Join(aWayPrizes, ",")
						iMaxMultiple = 0
					}

					fMinPrize := common.ArrayMinF(aWayPrizesMin)

					//根据奖金是否开启显示
					sDisplayPrize := "0"
					if fMaxPrize == fMinPrize {
						sDisplayPrize = "1"
					}

					//根据投注方式开启是否开启显示
					if common.InArrayInt(oSeriesWay.Id, NoDyMostPrizeSeriesWayIDs) {
						sDisplayPrize = "0"
					}

					aPrizeSettings[oSeriesWay.Id] = map[string]interface{}{
						"name":          aWayData.NameCn,
						"prize":         sPrize,
						"max_multiple":  iMaxMultiple,
						"display_prize": sDisplayPrize,
					}

					aArr3 = append(aArr3, aWayData)
				}

				if len(aArr3) >= 1 {
					aSubData.Children = aArr3
					aArr2 = append(aArr2, aSubData)
				}

			}
			if len(aArr2) >= 1 {
				aData.Children = aArr2
				aArr1 = append(aArr1, aData)
			}

		}

		//緩存選好盤
		jsons, _ := json.Marshal(aArr1) //转换成JSON返回的是byte[]
		sResult = string(jsons)
		redisClient.Redis.StringWrite(sCacheKey, sResult, -1)

		//緩存獎金組
		jsonsForPrize, _ := json.Marshal(aPrizeSettings)
		sPrizeData = string(jsonsForPrize)
		redisClient.Redis.StringWrite(sCacheKeyPrizeSetting, sPrizeData, -1)

	} else {
		json.Unmarshal([]byte(sResult), &aArr1)
		json.Unmarshal([]byte(sPrizeData), &aPrizeSettings)
	}
	return aArr1, aPrizeSettings
}

/*
 * 获得玩法投注方式
 */
func (m *tWayGroups) GetWays(bForBet bool, sWaitGroupId string) []dao.WayGroupWays {
	var sField string
	if bForBet {
		sField = "for_display"
	} else {
		sField = "for_search"
	}
	mConditions := map[string]string{
		"group_id": sWaitGroupId,
		sField:     "1",
	}

	aSortBy := []string{"sequence"}
	aOrderBy := []string{"asc"}

	aData, _ := dao.GetAllWayGroupWays(mConditions, nil, aSortBy, aOrderBy, 0, 10000)
	return aData
}
