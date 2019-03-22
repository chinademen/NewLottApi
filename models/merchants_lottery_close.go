package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
	"strconv"
	"strings"
)

type tMerchantsLotteryClose struct {
	Table
}

var MerchantsLotteryClose = &tMerchantsLotteryClose{Table: Table{TableName: "merchants_lottery_close"}}

/*
 * 返回当前接入商的彩种关闭id切片
 * string sId
 */
func (m *tMerchantsLotteryClose) GetInfoById(sId string) []string {
	oData := m.GetInfo(sId)

	//如果为空，说明该接入商所有彩种都有效
	if oData.Id < 1 {
		return []string{}
	}
	return strings.Split(oData.LotteryIds, ",")

}

/*
 * 獲取接入商關閉的彩種
 * string sId
 */
func (m *tMerchantsLotteryClose) GetInfo(sId string) *dao.MerchantsLotteryClose {
	iId, _ := strconv.Atoi(sId)
	sCacheKey := CompileMerchantsLotteryCloseCacheKey(sId)
	sCahceData := redisClient.Redis.StringRead(sCacheKey)
	var aResult []*dao.MerchantsLotteryClose
	if len(sCahceData) < 1 {
		aData, _ := dao.GetAllMerchantsLotteryClose(nil, nil, nil, nil, 0, 1)
		b, _ := json.Marshal(aData)
		redisClient.Redis.StringWrite(sCacheKey, string(b), -1)
		json.Unmarshal(b, &aResult)
	} else {
		json.Unmarshal([]byte(sCahceData), &aResult)
	}

	for _, oValue := range aResult {
		if oValue.Id == iId {
			return oValue
		}
	}
	var empty *dao.MerchantsLotteryClose
	return empty
}
