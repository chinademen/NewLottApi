package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
)

type tBasicWays struct {
	Table
}

const (
	WAY_MULTI_SEQUENCING = 7
)

var BasicWays = &tBasicWays{Table: Table{TableName: "basic_ways"}}

/*
 * 獲取基礎投注方式數據
 */
func (m *tBasicWays) GetInfoById(iId int) *dao.BasicWays {

	sCahceKey := CompileBasicWaysCacheKey()
	sCahceData := redisClient.Redis.StringRead(sCahceKey)
	var aResult []*dao.BasicWays
	if len(sCahceData) < 1 {
		aData, _ := dao.GetAllBasicWays(nil, nil, nil, nil, 0, 10000)
		b, _ := json.Marshal(aData)
		aResult = aData
		redisClient.Redis.StringWrite(sCahceKey, string(b), -1)
		json.Unmarshal(b, &aResult)
	} else {
		json.Unmarshal([]byte(sCahceData), &aResult)
	}

	for _, oValue := range aResult {
		if oValue.Id == iId {
			return oValue
		}
	}
	var empty *dao.BasicWays
	return empty
}
