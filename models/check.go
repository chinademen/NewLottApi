package models

import (
	cm "common"
	"regexp"
	"strings"

	"github.com/astaxie/beego/validation"
)

/**
 * 用户输入验证规则
 */
type regInput struct {
	Account  string `valid:"Required;AlphaNumeric;MinSize(6);MaxSize(12);Match(/^[a-zA-Z].[a-zA-Z0-9]/)"` //必填，长度6-12位,首位必须字母,字母和数字，其他类型不能通过
	Password string `valid:"Required;MinSize(6);MaxSize(16);Match(/^[a-zA-Z0-9]{6,16}$/)"`                //必填，长度6-20位,字母和数字，其他类型不能通过
}

type transferInput struct {
	Account  string `valid:"Required;AlphaNumeric;MinSize(6);MaxSize(12);Match(/^[a-zA-Z].[a-zA-Z0-9]/)"` //必填，长度6-12位,首位必须字母,字母和数字，其他类型不能通过
	Password string `valid:"Required;MinSize(6);MaxSize(16);Match(/^[a-zA-Z0-9]{6,16}$/)"`                //必填，长度6-20位,字母和数字，其他类型不能通过
	OrderId  string `valid:"Required;MinSize(5);MaxSize(45)"`                                             //必填
	Amount   string `valid:"Required;MinSize(1);MaxSize(20);Match(/^[0-9]*$/)"`                           //必填，必须时整数
}

type assignRegInput struct {
	Account  string `valid:"Required;AlphaNumeric;MinSize(6);MaxSize(12);Match(/^[a-zA-Z].[a-zA-Z0-9]/)"` //必填，长度6-12位,首位必须字母,字母和数字，其他类型不能通过
	Password string `valid:"Required;MinSize(6);MaxSize(16);Match(/^[a-zA-Z0-9]{6,16}$/)"`                //必填，长度6-20位,字母和数字，其他类型不能通过
}

/**
* 自定义表单验证提示
 */
var MessageTmpls = map[string]string{
	"Required":     "不能为空",
	"Min":          "最小值为 %d",
	"Max":          "最大值为 %d",
	"Range":        "请输入数字范围%d到%d",
	"MinSize":      "最小长度为%d位",
	"MaxSize":      "最大长度为%d位",
	"Length":       "固定长度为%d位",
	"Alpha":        "必须是26个英文字母",
	"Numeric":      "必须输入一个数字",
	"AlphaNumeric": "必须是字符或者数字",
	"Match":        "必须符合规则 %s",
	"NoMatch":      "必须不符合规则 %s",
	"AlphaDash":    "必须是数字或者字母或者-_",
	"Email":        "必须是email地址",
	"IP":           "必须是ipv4地址",
	"Base64":       "必须是base64位地址",
	"Mobile":       "必须是手机号码",
	"Tel":          "必须是固定电话号码",
	"Phone":        "必须是固定电话或者手机号码",
	"ZipCode":      "必须是邮政编码",
}

/**
* 将提示错误的key替换成中文
 */
func titleInput(key string) string {
	title := key
	switch key {
	case "Account":
		title = "用户名"
		break
	case "Password":
		title = "密码"
		break
	case "Phone":
		title = "手机号码"
		break
	case "Email":
		title = "邮箱"
		break
	default:
		title = key
	}
	return title
}

/**
* 用户名和密码的检测
* @account	string	用户名
* return 	int		状态值
* return	string	错误信息
 */
func ChkUserAndPwd(account, password string) (int, string) {
	status := 401
	msg := "验证失败"
	//第一步判断用户名输入的合法性
	status, msg = ChkUserInput(account, password)
	return status, msg
}

/**
* 验证用户名和密码
 */
func ChkUserInput(account, password string) (int, string) {
	status := 200
	msg := "请求完成"
	//判断输入的数据是否合法
	valid := validation.Validation{}
	validation.MessageTmpls = MessageTmpls
	u := regInput{Account: account, Password: password}
	b, _ := valid.Valid(&u)
	if !b {
		status = 441
		for _, err := range valid.Errors {
			msg = titleInput(err.Field) + err.Message
			break
		}
	}
	return status, msg
}

/**
* 验证转账参数是否合法
 */
func ChkTransferParamInput(accesscode, account, password, amount, url, orderId, isPush string) (int, string) {
	status := 200
	msg := "请求完成"
	if isPush == "true" {
		if len(url) < 6 {
			status = 442
			msg = "url格式不正确"
		} else {
			preUrl := cm.Substr(url, 0, 5)
			if preUrl != "http:" && preUrl != "https" {
				status = 443
				msg = "url格式不正确"
			}
		}
	}
	if status == 200 {
		//判断输入的数据是否合法
		valid := validation.Validation{}
		validation.MessageTmpls = MessageTmpls
		u := transferInput{Account: account, Password: password, OrderId: orderId, Amount: amount}
		b, _ := valid.Valid(&u)
		if !b {
			status = 444
			for _, err := range valid.Errors {
				msg = err.Key + err.Message
			}
		}

	}

	return status, msg
}

/**
* sql注入过滤判断
* @param	str	需要做判断的字符串
* return	bool	如果包含注入关键词输出true，否则false
 */
func ChkDanger(str string) bool {
	str = strings.ToLower(str)
	danger_keys := `(?:')|(?:--)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, _ := regexp.Compile(danger_keys)
	return re.MatchString(str)
}

/**
 * 判断IP白名单和接入商账户是否正确
 * @return int
 */
func CheckIPAndAccess(ip string, sIdentity string) (int, string) {
	if ip == "127.0.0.1" || ip == "::1" {
		return 200, "正常"
	}
	//第一步骤，检查接入商是否存在
	merchantInfo := Merchants.GetInfoForRedis(sIdentity)
	if len(merchantInfo) < 1 {
		return 301, "接入商不存在"
	}
	if merchantInfo["status"] == "0" {
		//接入商被锁定
		return 304, "接入商被锁定"
	}

	chkIP := MerchantIp.CheckIP(merchantInfo["id"], ip)
	if chkIP {
		return 200, "正常"
	} else {
		return 303, "您的IP:" + ip + "不在访问白名单中"
	}
}

/*
 * 检查日期格式是否正确
 */
func ChkIsDate(dateStr string) (int, string, int64) {

	rStatus := 200
	rMsg := "ok"

	if len(dateStr) == 0 {
		return 701, "日期不能为空", 0
	}
	dTime, dErr := cm.FormatStr2Unix(dateStr, cm.DATE_FORMAT_YMDHIS)
	if dErr != nil {
		return 702, "日期格式不对", 0
	}

	if dTime < 946656000 {
		return 703, "日期不能小于2000-01-01 00:00:00", 0
	}

	return rStatus, rMsg, dTime
}
