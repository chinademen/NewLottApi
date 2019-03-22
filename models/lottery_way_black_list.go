package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
)

type tLotteriesWayBlackList struct {
	Table
}

var LotteriesWayBlackList = &tLotteriesWayBlackList{Table: Table{TableName: "lottery_way_black_list"}}

/*
 * 过滤彩种投注方式
 */
func (c *tLotteriesWayBlackList) GetTerminalBlackList(sLotteryId, sTerminalId, sSeriesWayId string) []*dao.LotteryWayBlackList {
	sCacheKey := CompileTerminalBlackListAllCacheKey(sLotteryId, sTerminalId, sSeriesWayId)
	sCahceData := redisClient.Redis.StringRead(sCacheKey)
	var aResult []*dao.LotteryWayBlackList
	if len(sCahceData) < 1 {
		aConditions := map[string]string{
			"lottery_id":  sLotteryId,
			"series_way":  sSeriesWayId,
			"terminal_id": sTerminalId,
		}
		aData, _ := dao.GetAllLotteryWayBlackList(aConditions, nil, nil, nil, 0, 1000)
		b, _ := json.Marshal(aData)
		redisClient.Redis.StringWrite(sCacheKey, string(b), -1)
		json.Unmarshal(b, &aResult)
	} else {
		json.Unmarshal([]byte(sCahceData), &aResult)
	}

	return aResult
}
