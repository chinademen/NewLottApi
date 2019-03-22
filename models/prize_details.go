package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
)

type tPrizeDetails struct {
	Table
}

var PrizeDetails = &tPrizeDetails{Table: Table{TableName: "prize_details"}}

/*
 * 獲取用戶奖金详情缓存
 */
func (m *tPrizeDetails) GetDetailsNew(sGroupId string) map[string]map[string]interface{} {
	sCacheKey := MakeCacheKeyOfGroupNew(sGroupId)
	sData := redisClient.Redis.StringRead(sCacheKey)

	var aDetails = map[string]map[string]interface{}{}
	var sLastMethodId uint32 = 0
	if len(sData) < 1 {
		mConditions := map[string]string{
			"group_id": sGroupId,
		}

		aPrizeDetails, _ := dao.GetAllPrizeDetails(mConditions, nil, nil, nil, 0, 10000)
		for _, oPrizeDetails := range aPrizeDetails {
			sMethodId := fmt.Sprintf("%d", oPrizeDetails.MethodId)
			sLevel := fmt.Sprintf("%d", oPrizeDetails.Level)
			if sLastMethodId != oPrizeDetails.MethodId {

				aDetails[sMethodId] = map[string]interface{}{
					"method_id":   sMethodId,
					"method_name": oPrizeDetails.MethodName,
					"level":       map[string]float64{},
				}
				sLastMethodId = oPrizeDetails.MethodId
			}
			m := map[string]float64{sLevel: oPrizeDetails.Prize}
			aDetails[sMethodId]["level"] = m
		}

		sResult, _ := json.Marshal(aDetails)
		redisClient.Redis.StringWrite(sCacheKey, string(sResult), 0)
	} else {
		json.Unmarshal([]byte(sData), &aDetails)
	}
	return aDetails
}

/*
 * 根據獎金組和基礎玩法獲取對應的獎金設置
 * @param sGroupId string
 * @param sBasicMethodId string
 */
func (m *tPrizeDetails) GetPrizeSetting(sGroupId, sBasicMethodId string) map[int]float64 {
	mCondition := map[string]string{
		"group_id":  sGroupId,
		"method_id": sBasicMethodId,
	}
	aPrizeDetails, _ := dao.GetAllPrizeDetails(mCondition, nil, nil, nil, 0, 10000)
	var mPrize = map[int]float64{}
	for _, oPrizeDetail := range aPrizeDetails {
		mPrize[int(oPrizeDetail.Level)] = oPrizeDetail.Prize //map[任务条件等级]獎金
	}
	return mPrize
}
