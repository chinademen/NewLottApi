package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"fmt"
)

type tMerchantIp struct {
	TbName string
	Fields string
}

var MerchantIp = &tMerchantIp{TbName: "merchants_ip",
	Fields: "`id`, `username`, `password`, `fund_password`, `account_id`, `merchant_id`, `prize_group`, `blocked`, `realname`, `nickname`, `email`, `mobile`, `is_tester`, `bet_multiple`, `bet_coefficient`, `login_ip`, `register_ip`, `token`, `signin_at`, `activated_at`, `register_at`, `deleted_at`, `created_at`, `updated_at`"}

var RMerchantIpKeyEX int = 3600

/**
 * 判断ip是否存在
 * @param		sMerchantId		接入商id
 * @param		ip				访问的接口ip
 * return		判断结果，如果正常返回true
 */
func (c *tMerchantIp) CheckIP(sMerchantId, ip string) bool {
	if len(sMerchantId) < 1 {
		return false
	}
	var queryMap = map[string]string{}
	queryMap["merchant_id"] = sMerchantId
	queryMap["ip"] = ip
	row, _ := dao.GetAllMerchantsIp(queryMap, nil, nil, nil, 0, 100)
	if len(row) > 0 {
		return true
	}

	return false
}

/**
* 从数据库中读取ip白名单列表
* @param		sMerchantId			接入商id
* return		[]string
 */
func (c *tMerchantIp) GetIP(sMerchantId string) []string {
	sRes := []string{}
	queryMap := make(map[string]string)
	queryMap["merchant_id"] = sMerchantId
	row, _ := dao.GetAllMerchantsIp(queryMap, nil, nil, nil, 0, 100)
	if len(row) > 0 {
		for _, v := range row {
			sRes = append(sRes, v.Ip) //v["ip"])
		}
	}
	return sRes
}

/*
 查询资料
merchantId	商户id
ip			用户名
field		查询的字段
*/
func (m *tMerchantIp) GetField(merchantId, ip, field string) string {

	//从数据库读取结果
	sWhere := fmt.Sprintf("merchant_id = '%s' AND ip = '%s'", merchantId, ip)
	rMap := GetOne(m.TbName, sWhere, field)

	rStr := ""
	if len(rMap) > 0 && len(rMap[field]) > 0 {
		rStr = rMap[field]
	}

	return rStr
}

/*
 * 查询资料      redis
 * merchantId	商户id
 * ip			用户名
 */
func (m *tMerchantIp) RGetField(merchantId, ip, field string) string {

	rKey := CopmileMerchantIpOneKey(merchantId, ip)

	//优先读redis缓存
	rString := redisClient.Redis.StringRead(rKey)
	if len(rString) > 0 {
		return rString
	}

	rString = m.GetField(merchantId, ip, field)

	//将结果写入redis，缓存1小时
	if len(rString) > 0 {
		redisClient.Redis.StringReWrite(rKey, rString, RMerchantIpKeyEX)
	}

	return rString
}
