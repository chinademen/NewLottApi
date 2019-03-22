package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"encoding/json"
)

type tTransactionTypes struct {
	TbName string
}

var TransactionTypes = &tTransactionTypes{TbName: "transaction_types"}

const (
	TYPE_FREEZE_FOR_TRACE = 5
	TYPE_UNFREEZE_FOR_BET = 6
	TYPE_BET              = 7
)

/*
 * 根據id獲取賬變類型
 */
func (m *tTransactionTypes) GetInfo(iId int) *dao.TransactionTypes {
	aTransactionTypes := m.GetAll()
	for _, oValue := range aTransactionTypes {
		if oValue.Id == iId {
			return oValue
		}
	}
	var empty *dao.TransactionTypes
	return empty
}

/*
 * 獲取所有賬變類型
 */
func (m *tTransactionTypes) GetAll() []*dao.TransactionTypes {
	sCacheKey := CompileTransactionsTypeCacheKey()
	sCacheData := redisClient.Redis.StringRead(sCacheKey)
	var aAll []*dao.TransactionTypes
	if len(sCacheData) > 0 {
		json.Unmarshal([]byte(sCacheData), &aAll)
	} else {
		aAll, _ = dao.GetAllTransactionTypes(nil, nil, nil, nil, 0, 1000)
		sAllbyte, _ := json.Marshal(aAll)
		redisClient.Redis.StringWrite(sCacheKey, string(sAllbyte), -1)
	}

	return aAll
}
