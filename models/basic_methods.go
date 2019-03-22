package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
	"strconv"
)

type tBasicMethods struct {
	Table
}

var BasicMethods = &tBasicMethods{Table: Table{TableName: "basic_methods"}}

/*
 * 根据Ｉd获取玩法数据
 */
func (m *tBasicMethods) GetInfo(sId string) *dao.BasicMethods {

	//從單一緩存key讀取
	sCahceKey := CompileOneBasicMethodCacheKey(sId)
	mCahceData := redisClient.Redis.HashReadAllMap(sCahceKey)
	iId, _ := strconv.Atoi(sId)
	var oBasicMethod *dao.BasicMethods
	if len(mCahceData) < 1 {
		oBasicMethod, _ = dao.GetBasicMethodsById(iId)
		if oBasicMethod.Id != 0 {
			m := map[string]string{}
			b, _ := json.Marshal(oBasicMethod)
			json.Unmarshal(b, &m)
			redisClient.Redis.HashWrite(sCahceKey, m, 60)
		}
	} else {
		cahceData, _ := json.Marshal(mCahceData)
		json.Unmarshal(cahceData, &oBasicMethod)
	}

	if oBasicMethod.Id > 0 {
		return oBasicMethod
	}

	//從所有緩存讀取
	aBasicMethods := m.GetAll()
	for _, oValue := range aBasicMethods {
		if oValue.Id == iId {
			return oValue
		}
	}

	return oBasicMethod
}

/*
 * 获取所有玩法数据
 */
func (m *tBasicMethods) GetAll() []*dao.BasicMethods {
	sCahceKey := CompileBasicMethodCacheKey()
	sCahceData := redisClient.Redis.StringRead(sCahceKey)
	var aResult []*dao.BasicMethods
	if len(sCahceData) < 1 {
		aData, _ := dao.GetAllBasicMethods(nil, nil, nil, nil, 0, 10000)
		b, _ := json.Marshal(aData)
		redisClient.Redis.StringWrite(sCahceKey, string(b), -1)
		json.Unmarshal(b, &aResult)
	} else {
		json.Unmarshal([]byte(sCahceData), &aResult)
	}

	return aResult
}
