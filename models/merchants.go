package models

import (
	"NewLottApi/dao"
	"common"
	"common/ext/redisClient"
	"fmt"
	"strconv"
)

type tMerchants struct {
	TbName string
	Fields string
}

var Merchants = &tMerchants{TbName: "merchants", //商户表
	Fields: "`id`, `identity`, `name`, `wallet_id`, `safe_key`, `post_url`, `status`, `is_tester`, `template`, `remark`, `created_at`, `updated_at`"}

var RMerchantKey string = "merchants:"
var RMerchantRowKey string = RMerchantKey + "row:identity_%s" //%s=商户唯一标识
var RMerchantRowKeyId string = RMerchantKey + "row:id_%s"     //%s=id
var RMerchantKeyEX int = 60 * 60 * 24

/**
 * 通过截获用户密钥去设置通用加密key
 */
func SetPrivateKey(privateKey string) string {
	HookAES = common.SetAES(privateKey, "", "pkcs5") //设置密钥
	return privateKey
}

//merchantToMap MerchantToMap
func merchantToMap(merchant *dao.Merchants) map[string]string {
	result := make(map[string]string)
	result["created_at"] = merchant.CreatedAt.String()
	result["id"] = strconv.Itoa(merchant.Id)
	result["identity"] = merchant.Identity
	result["isTester"] = strconv.Itoa(int(merchant.IsTester))
	result["name"] = merchant.Name
	result["postUrl"] = merchant.PostUrl
	result["remark"] = merchant.Remark
	result["safeKey"] = merchant.SafeKey
	result["status"] = strconv.Itoa(int(merchant.Status))
	result["template"] = strconv.Itoa(int(merchant.Template))
	result["updated_at"] = merchant.UpdatedAt.String()
	return result
}

/**
*通过code得到接入商的id，redis缓存24小时
* @param		identity	编码
* return		map[string]string
 */
func (m *tMerchants) GetInfoForRedis(identity string) map[string]string {

	merchantMap := map[string]string{}
	if len(identity) < 1 {
		return merchantMap
	}

	rKey := fmt.Sprintf(RMerchantRowKey, identity)

	//优先redis
	mRes := redisClient.Redis.HashReadAllMap(rKey)
	if len(mRes) > 0 {
		return mRes
	}

	merchantData, err := m.GetInfo(identity)
	if err != nil {
		return merchantMap
	}

	if merchantData != nil {

		if len(merchantToMap(merchantData)) > 0 {

			merchantMap = merchantToMap(merchantData)
			//写入redis
			redisClient.Redis.HashWrite(rKey, merchantMap, RMerchantKeyEX)
		}
	}
	return merchantMap
}

/**
* 读取 merchants 某一个指定的字段
* @param identity  		string	接入商的英文编码
* @param field			string	接入商属性字段
* @return				string	返回一个字符串
 */
func (m *tMerchants) MerchantsGetField(identity, field string) string {
	res := ""
	MerchantsRow := m.RGetByIdentity(identity)
	if len(MerchantsRow) > 0 {
		r, ok := MerchantsRow[field]
		if ok {
			res = r
		}
	}
	return res
}

/**
* 从数据库中读取接入商资料
* @param		merchantId
* return		map[string]string
 */
func (m *tMerchants) GetInfoById(merchantId string) map[string]string {
	sWhere := fmt.Sprintf("id='%s'", merchantId)

	row := GetOne(m.TbName, sWhere, m.Fields)
	return row
}

/**
* 从数据库中读取接入商资料
* @param		identity		编码
* return		map[string]string
 */
func (c *tMerchants) GetInfo(identity string) (*dao.Merchants, error) {
	queryMap := make(map[string]string)
	queryMap["identity"] = identity
	merchantList, err := dao.GetAllMerchants(queryMap, nil, nil, nil, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(merchantList) == 0 {
		return nil, nil
	}
	return &(merchantList[0]), nil
}

/**
* 实时读取接入商的id
* @param		identity		编码
* return		string
 */
func (c *tMerchants) GetId(identity string) string {
	id := ""
	sWhere := fmt.Sprintf("identity='%s'", identity)
	row := GetOne(c.TbName, sWhere, "id")

	if len(row) > 0 {
		id = row["id"]
	}
	return id
}

/**
 *@param string
 *@param string
 *@param string
 *@param string
 *@return
 */
func (c *tMerchants) Gets(where, field string, offset, limit int) []map[string]string {
	sWhere := fmt.Sprintf("status='%s'", "1")

	if len(where) > 0 {
		sWhere = sWhere + " and " + where
	}
	sOrder := " id  desc "

	sField := c.Fields

	if len(field) > 0 {
		sField = field
	}
	return GetList(c.TbName, sWhere, sField, sOrder, offset, limit)
}

/*
 查询资料
identity		商户唯一标识
*/
func (m *tMerchants) GetByIdentity(identity string) map[string]string {

	//从数据库读取结果
	sWhere := fmt.Sprintf("identity = '%s'", identity)
	rMap := GetOne(m.TbName, sWhere, m.Fields)

	return rMap
}

/*
 * 查询资料 redis
 * identity		商户唯一标识
 */
func (m *tMerchants) RGetByIdentity(identity string) map[string]string {

	rKey := fmt.Sprintf(RMerchantRowKey, identity)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetByIdentity(identity)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RMerchantKeyEX)
	}

	return rMap
}

/**
* 根据id获取某条记录
 */
func (m *tMerchants) GetRow(id string) map[string]string {
	sWhere := fmt.Sprintf("id='%s'", id)
	return GetOne(m.TbName, sWhere, m.Fields)
}

/*
 查询一行 redis
*/
func (m *tMerchants) RGetById(id string) map[string]string {

	rKey := fmt.Sprintf(RMerchantRowKeyId, id)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetRow(id)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RMerchantKeyEX)
	}

	return rMap
}
