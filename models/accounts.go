package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"fmt"
)

type tAccounts struct {
	TbName string
	Fields string
}

var Accounts = &tAccounts{TbName: "accounts",
	Fields: "`id`, `merchant_id`, `user_id`, `username`, `is_tester`, `balance`, `frozen`, `available`, `status`, `locked`, `created_at`, `updated_at`, `backup_made_at`"}

var RAccountsKey string = "accounts:"                //redis基本key
var RAccountsOneKey string = RAccountsKey + "one:%s" //redis 字符串key
var RAccountsRowKey string = RAccountsKey + "row:%s" //redis 数据库一行key
var RAccountsKeyEX int = 3600

/*
 * 查询资料 accountsId
 */
func (m *tAccounts) GetById(accountsId string) map[string]string {

	//从数据库读取结果
	sWhere := fmt.Sprintf("id = '%s'", accountsId)
	rMap := GetOne(m.TbName, sWhere, m.Fields)

	return rMap
}

/*
 * 查询资料 accountsId
 */
func (m *tAccounts) RGetById(accountsId string) map[string]string {

	rKey := fmt.Sprintf(RAccountsRowKey, accountsId)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetById(accountsId)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RAccountsKeyEX)
	}

	return rMap
}

/**
 * 将插入数据还原成sql语句
 * @param		d		插入资料
 * @return				sql
 */
func (m *tAccounts) GetAddSql(d map[string]string) string {

	sql := GetInsertSql(m.TbName, d)
	return sql
}

/**
 * 将插入数据还原成sql语句
 * @param		d		插入资料
 * @return				sql
 */
func (m *tAccounts) GetAddOnlySql(d map[string]string) string {

	sql := GetInsertTrueSql(m.TbName, d)
	return sql
}

/**
 * 插入数据
 * @param		d		插入资料
 */
func (m *tAccounts) DbInsert(d map[string]string) (int, string) {

	sqlOk, sqlId := Insert(d, m.TbName)
	return sqlOk, sqlId
}

/*
 * 判斷賬戶是否被鎖
 * @param iAccountId int
 */
func (m *tAccounts) IsLockAccount(oAccount *dao.Accounts) bool {
	if oAccount.Locked == 0 {
		return false
	}
	return true
}

/*
 * 獲取用戶的賬戶
 */
func (m *tAccounts) GetAccountInfoByUserId(sUserId string) dao.Accounts {
	mConditions := map[string]string{
		"user_id": sUserId,
	}
	aAccounts, _ := dao.GetAllAccounts(mConditions, nil, nil, nil, 0, 1)
	if len(aAccounts) < 1 {
		return dao.Accounts{}
	}
	return aAccounts[0]
}
