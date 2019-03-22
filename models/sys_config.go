package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
)

type tSysConfig struct {
	Table
}

var SysConfig = &tSysConfig{Table: Table{TableName: "sysconfig"}}

/*
 * 獲取系統配置信息
 * @param string sField 字段名字
 * return string
 */
func (m *tSysConfig) GetPrizeByItem(sField string) string {
	sCacheKey := CompileSysConfigsCacheKey(sField)
	sCacheData := redisClient.Redis.StringRead(sCacheKey)
	var oResult *dao.SysConfigs
	if len(sCacheData) >= 1 {
		json.Unmarshal([]byte(sCacheData), &oResult)
	} else {
		mCondition := map[string]string{
			"item": sField,
		}
		aData, _ := dao.GetAllSysConfigs(mCondition, nil, nil, nil, 0, 1)
		if len(aData) < 1 {
			return ""
		}
		b, _ := json.Marshal(aData[0])
		redisClient.Redis.StringWrite(sCacheKey, string(b), -1)
		json.Unmarshal(b, &oResult)
	}

	return oResult.Value
}
