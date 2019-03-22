package models

import (
	"NewLottApi/dao"
	"common"
	"common/ext/redisClient"
	"fmt"
	"strings"
	"time"
)

var RMerchantsDomainsKey string = "merchants_domains:"
var RMerchantsDomainsOneKey string = RMerchantsDomainsKey + "one:id_%s:domains_%s" //%s=商户id, %s=域名
var RMerchantsDomainsKeyEX int = 60 * 60 * 24

/**
 * 检测域名是否合法域名(不包含www.)
 * @params map  请求的参数和host
 * @return
 * int 状态
 * string 信息
 */
func CheckHost(merchantId, host string) (int, string) {
	status := 110
	msg := "域名不能为空"
	if len(host) > 0 {
		//获取顶级域名
		domain := common.GetTopDomain(host)
		if debug {
			fmt.Println("CheckHost->", domain)
		}

		status = 111
		msg = "域名不合法"
		domainChk := MmerchantsDomainsByR(merchantId, domain)
		if len(domainChk) > 0 || strings.Count(host, "127.0.0.1") > 0 {
			status = 200
			msg = "域名合法"
		}
	}
	return status, msg
}

/**
 * 读取域名	reids缓存
 * @domain	需要查询的domain
 */
func MmerchantsDomainsByR(merchantId, domain string) string {

	rKey := fmt.Sprintf(RMerchantsDomainsOneKey, merchantId, domain)

	//优先redis
	rStr := redisClient.Redis.StringRead(rKey)
	if len(rStr) > 0 {
		return rStr
	}

	mWhere := map[string]string{"domain": domain}
	sqlMaps, sqlErr := dao.GetAllMerchantsDomains(mWhere, nil, nil, nil, 0, 1)
	if sqlErr == nil {

		if len(sqlMaps) > 0 {

			rStr = sqlMaps[0].Domain
			redisClient.Redis.StringWrite(rKey, rStr, RMerchantsDomainsKeyEX)
		}
	}

	return rStr
}

/**
 * 对过于频发的请求进行锁定,防止cc攻击
 * 设置每个ip一分钟不能超过200次，请求否则，锁定10分钟
 */
func ChkAccessNumber(sIp string) bool {
	res := false
	rKey := fmt.Sprintf("router_ip:%s_%s", sIp, time.Now().Format("200601021504"))
	num := redisClient.Redis.IntRead(rKey)
	if num >= 2000 {
		res = true
	} else {
		num++
		redisClient.Redis.IntReplaceWrite(rKey, num, 60)
	}
	return res
}

/**
* 解密出真实参数到公用变量paramsMap
* @param params	string	传递的参数值
* @param bNeedDes	bool	是否需要解密数据 true:需要  false:不需要
* @return
	int		状态码   完成=200    输入错误=400
	string	文字描述
	map	解密后的参数结果集
*/
func ChkInputAndMap(params string, bNeedDes bool) (int, string, map[string]string) {
	status := 200
	msg := "请求完成"
	res := map[string]string{}

	if len(params) < 1 {
		status = 400
		msg = "params参数不得为空"
	} else {
		paramsMapStr := params

		if bNeedDes {
			//数据解密
			sEncodeParam := strings.Replace(params, "%2B", "+", -1)
			common.LogsWithFileName("", "api_request", "params->"+params+"\r\nen_param->"+sEncodeParam)
			paramsMapStr = HookAES.AesDecryptString(sEncodeParam)

		}
		common.LogsWithFileName("", "api_request", paramsMapStr+"\r\n==================================\r\n\r\n")

		if debug {
			fmt.Println("paramsMapStr->", paramsMapStr)
		}

		//防注入分析
		if ChkDanger(paramsMapStr) == true {
			status = 500
			msg = "非法操作，包含入侵字符串"
		} else {
			//将解密后的数据切割，并解析到公用参数map里面去
			mp := strings.Split(paramsMapStr, "&")
			for i := 0; i < len(mp); i++ {
				if len(mp[i]) > 0 {
					tm := strings.Split(mp[i], "=")
					if len(tm) >= 2 {
						res[tm[0]] = tm[1]
					}
				}
			}
		}

	}

	return status, msg, res
}

/*
 * 是否商户ip白名单
 */
func IsMerchantsIpWhite(merchantId, clientIP string) (int, string) {

	if len(merchantId) == 0 {
		return 401, "商户id不能为空"
	}

	if len(clientIP) == 0 {
		return 402, "ip不能为空"
	}
	getId := MerchantIp.RGetField(merchantId, clientIP, "id")

	if len(getId) > 0 {
		return 200, "ok"
	}

	sErr := fmt.Sprintf("ip非法, id=%s", clientIP)
	return 403, sErr
}
