package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
	"strings"

	"lotteryJobs/models"
	//	"github.com/astaxie/beego/orm"
)

const (
	RedisKeyUserPrizeSets = "user_prize_sets"
)

type tUserPrizeSets struct {
	Table
}

var UserPrizeSets = &tUserPrizeSets{Table: Table{TableName: "user_prize_sets"}} //彩种玩法分组表

/**
 * 获取用户奖金组ID
 * @param int iUserId
 * @param int iLotteryId
 * @return int
 */
func (m *tUserPrizeSets) GetGroupId(sUserId, sLotteryId string) (int, string) {
	oGroup := m.GetUserPrizeSet(sUserId, sLotteryId)
	return int(oGroup.GroupId), oGroup.PrizeGroup
}

/*
 * 获取 当前用户 存在的玩法
 * sLotteryId string
 * sUserId string
 */
func (m *tUserPrizeSets) GetUserLotteryMethod(sUserId string) map[string]map[string]string {
	aData := models.UserPrizeSet.GetUserLotteryMethod(sUserId)
	var aResult = map[string]map[string]string{}
	for _, mData := range aData {
		aResult[mData["id"]] = mData
	}
	return aResult
}

/**
 * 获取用户奖金组设置
 * @param string iUserId
 * @param string iLotteryId
 * @return dao.UserPrizeSets
 */
func (m *tUserPrizeSets) GetUserPrizeSet(sUserId, sLotteryId string) dao.UserPrizeSets {
	sCacheKey := CompileCacheKeyOfUserLottery(sUserId, sLotteryId)
	sCacheData := redisClient.Redis.StringRead(sCacheKey)

	var aGroups = []dao.UserPrizeSets{}
	if len(sCacheData) < 1 {
		mConditions := map[string]string{
			"user_id":    sUserId,
			"lottery_id": sLotteryId,
		}
		aGroups, _ = dao.GetAllUserPrizeSets(mConditions, nil, nil, nil, 0, 1)
		if len(aGroups) < 1 {
			return dao.UserPrizeSets{}
		}
		jsons, _ := json.Marshal(aGroups) //转换成JSON返回的是byte[]
		sResult := string(jsons)

		//缓存
		redisClient.Redis.StringWrite(sCacheKey, sResult, 0)
	} else {
		json.Unmarshal([]byte(sCacheData), &aGroups)
	}

	if len(aGroups) < 1 {
		return dao.UserPrizeSets{}
	}

	return aGroups[0]
}

/*
 * 獲取最大投注獎金祖
 */
func (m *tUserPrizeSets) GetMaxBetGroup(iGroupFromUser uint, oSeries *dao.Series, oLottery *dao.Lotteries) uint {
	iMaxBetGroup := oLottery.MaxBetGroup
	if iMaxBetGroup <= 0 {
		iMaxBetGroup = uint(oSeries.MaxBetGroup)
	}
	if iMaxBetGroup >= iGroupFromUser {
		iMaxBetGroup = iGroupFromUser
	}

	return iMaxBetGroup
}

/*
 * 获取彩种投注方式的奖金列表
 */
func (m *tUserPrizeSets) GetPrizeSetOfUsers(sUserId, sLotteryId string, oSeriesWay *dao.SeriesWays) interface{} {
	var mData = map[string]interface{}{}
	aMethodIds := strings.Split(oSeriesWay.BasicMethods, ",")
	iGroupId, _ := m.GetGroupId(sUserId, sLotteryId)
	for _, sMethodId := range aMethodIds {
		mData[sMethodId] = PrizeDetails.GetPrizeSetting(fmt.Sprintf("%d", iGroupId), sMethodId)
	}
	return mData
}
