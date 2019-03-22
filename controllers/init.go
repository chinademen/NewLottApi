package controllers

import (
	"NewLottApi/models"
	"common"
	"strings"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type i18nLocale struct {
	i18n.Locale
}

type MainController struct {
	beego.Controller
}

/*
* 定义一个json的返回值类型(首字母必须大写)
 */
type JsonOut struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

//接收的 identity 参数的数据类型
type IdentityInputType struct {
	string
}

//接收的 params 参数的数据类型
type ParamsInputType struct {
	string
}

type Params struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	PrizeGroup    string `json:"prize_group"`
	Ip            string `json:"ip"`
	Device        string `json:"device"`
	Amount        string `json:"amount"`
	TypeName      string `json:"type"`
	OrderNumber   string `json:"order_number"`
	NewPassword   string `json:"new_password"`
	NewPrizeGroup string `json:"new_prize_group"`
}

var (
	debug           = false
	Logpath         = ""
	NetIp           = ""
	BaseControllerX = &i18nLocale{}
)

/**
 * 包的入口函数，用于初始化
 */
func init() {

	var controlDebug = beego.AppConfig.String("control_debug")
	Logpath = beego.AppConfig.String("logpath")
	NetIp = beego.AppConfig.String("net")

	//判断是否开启调试
	if controlDebug == "on" {
		debug = true
	}
	common.SetLogPath(Logpath)
}

/*
 * 过滤请求方式
 */
func (this *MainController) Prepare() {
	var d interface{}
	if this.Ctx.Request.Method != "POST" {
		this.RenderJson(501, "请求方式错误", d)
	}
}

/**
 *返回json给请求端
 */
func (c *MainController) RenderJson(status int, message string, d interface{}) {

	res := JsonOut{status, message, d}

	//输出json数据
	c.Data["json"] = &res
	c.ServeJSON()
	c.StopRun()
}

/**
 *返回字符串给请求端
 */
func (c *MainController) RenderString(message string) {

	c.Ctx.WriteString(message)
	c.StopRun()
}

//CheckError check error
func (c *MainController) CheckError(value interface{}, err error, msg string) *JsonOut {
	if err != nil {
		var result JsonOut
		result.Status = 500
		result.Msg = err.Error()
		result.Data = err
		return &result
	}
	if value == nil || value == 0 {
		var result JsonOut
		result.Status = 500
		result.Msg = msg
		return &result
	}
	return nil
}

//设置语言 中文或者英文...
func (i *i18nLocale) SetLang(lang string) {
	i.cSetLang(lang)
}

func (i *i18nLocale) cSetLang(lang string) {
	i.Lang = lang
}

//翻译语言
func (i *i18nLocale) GetLang(translation string) string {
	translationLow := strings.ToLower(translation)
	return i.Tr(translationLow)
}

/*
 *
 */
func (c *MainController) GetLoginInfo() (string, string, string, string) {
	clientIP := c.Ctx.Input.IP()
	userAgent := c.Ctx.Input.UserAgent()            //header中的user-agent
	browser, _ := common.GetBrowserOS(userAgent)    //浏览器类型,操作系统类型
	sVisitType := common.GetWebVisitType(userAgent) //设备类型
	terminal := "2"
	if sVisitType == "pc" {
		terminal = "1"
	}
	return clientIP, userAgent, browser, terminal
}

/*
 * 平台内部接口 钩子
 */
func (c *MainController) WebLogin() map[string]string {

	chk := 200
	msg := "验证通过"
	errData := map[string]interface{}{}
	pMap := map[string]string{}

	//设置param参数
	paramStr := strings.Trim(c.GetString("params"), " ")
	chk, msg, pMap = models.ChkInputAndMap(paramStr, false)
	if chk != 200 {
		c.RenderJson(chk, msg, errData)
	}

	if len(pMap["token"]) < 1 {
		chk = 602
		msg = "token不能为空"
		c.RenderJson(chk, msg, errData) //将数据装载到json返回值
	}

	clientIP, _, browser, terminal := c.GetLoginInfo()
	tokenMap := map[string]string{}

	//比較用戶登錄信息和數據庫的不同
	chk, msg, tokenMap = models.UserToken.CheckToken(pMap["token"], clientIP, terminal, browser)
	if chk != 200 {
		c.RenderJson(chk, msg, errData)
	}

	usMap := models.Users.RGetById(tokenMap["user_id"])
	if len(usMap["id"]) == 0 || len(usMap["username"]) == 0 {
		chk = 604
		msg = "用户信息错误"
		c.RenderJson(chk, msg, errData)
	}

	merchantMap := models.Merchants.RGetById(usMap["merchant_id"])
	if len(merchantMap["id"]) == 0 || len(merchantMap["identity"]) == 0 {
		chk = 605
		msg = "商户信息错误"
		c.RenderJson(chk, msg, errData)
	}

	pMap["identity"] = merchantMap["identity"]
	pMap["merchant_id"] = merchantMap["id"]
	pMap["user_id"] = usMap["id"]
	pMap["username"] = usMap["username"]
	pMap["client_ip"] = clientIP
	pMap["terminal"] = terminal

	return pMap
}

/**
 * 钩子商户对接调用接口
 */
func (c *MainController) FilterUser() map[string]string {

	//将数据装载到json返回值
	chk := 101
	msg := "商户标识不能为空"
	d := map[string]interface{}{}

	//1 先判断商户
	identity := strings.Trim(c.GetString("merchant_identity"), " ")
	if len(identity) == 0 {
		c.RenderJson(chk, msg, d)
	}

	merchantRow := models.Merchants.RGetByIdentity(identity)
	if len(merchantRow) == 0 {
		chk = 120
		msg = "商户不存在"
		c.RenderJson(chk, msg, d)
	}

	if merchantRow["status"] == "0" {
		chk = 102
		msg = "商户未激活"
		c.RenderJson(chk, msg, d)
	}

	if len(merchantRow["safe_key"]) == 0 {

		chk = 106
		msg = "商家没有加密键"
		c.RenderJson(chk, msg, d)
	}

	models.SetPrivateKey(merchantRow["safe_key"]) //3 设置商户的密钥
	merchantId := merchantRow["id"]               //商户id

	//ip
	clientIP := c.Ctx.Input.IP()
	userAgent := c.Ctx.Input.UserAgent()         //header中的user-agent
	browser, _ := common.GetBrowserOS(userAgent) //浏览器类型,操作系统类型

	//2 判断ip是否合法, 只有在merchants_ip设置的ip白名单中, 才可访问api
	// chk, msg = models.IsMerchantsIpWhite(merchantId, clientIP) //此时获取到的clientIP=商户ip
	// if chk != 200 {
	// 	c.RenderJson(chk, msg, d)
	// }

	//3 设置param参数
	paramStr := strings.Trim(c.GetString("params"), " ")
	chk, msg, pMap := models.ChkInputAndMap(paramStr, true)
	if chk != 200 {
		//将数据装载到json返回值
		c.RenderJson(chk, msg, d)
	}

	//保存session
	pMap["identity"] = identity
	pMap["merchant_id"] = merchantId
	pMap["browser"] = browser    //浏览器类型
	pMap["client_ip"] = clientIP //此时获取到的clientIP=商户ip
	return pMap
}

/*
 * 后台调用接口钩子
 */
func (c *MainController) AdminAPI() map[string]string {

	//将数据装载到json返回值
	d := map[string]interface{}{}
	chk := 200
	msg := "验证通过"

	pMap := map[string]string{}

	//设置param参数
	paramStr := strings.Trim(c.GetString("params"), " ")

	chk, msg, pMap = models.ChkInputAndMap(paramStr, false)
	if chk != 200 {
		c.RenderJson(chk, msg, d)
	}

	//验证签名
	sSignSource := ""
	if _, ok := pMap["admin_name"]; ok {
		sSignSource += pMap["admin_name"]
	}
	if _, ok := pMap["project_id"]; ok {
		sSignSource += pMap["project_id"]
	}
	if _, ok := pMap["trace_id"]; ok {
		sSignSource += pMap["trace_id"]
	}

	if _, ok := pMap["lottery_id"]; ok {
		sSignSource += pMap["lottery_id"]
	}

	if _, ok := pMap["issue"]; ok {
		sSignSource += pMap["issue"]
	}

	if _, ok := pMap["new_code"]; ok {
		sSignSource += pMap["new_code"]
	}

	if _, ok := pMap["begin_time"]; ok {
		sSignSource += pMap["begin_time"]
	}

	if _, ok := pMap["code_center_id"]; ok {
		sSignSource += pMap["code_center_id"]
	}
	sSignSource += beego.AppConfig.String("admin_api_key")
	sSign := common.GetMd5(sSignSource)
	if pMap["sign"] != sSign {
		chk = 602
		msg = "签名错误"
		c.RenderJson(chk, msg, d)
	}

	pMap["client_ip"] = c.Ctx.Input.IP()
	pMap["client_ip"] = "45.76.188.143"
	return pMap
}
