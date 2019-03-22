// @APIVersion 1.0.0
// @Title 新彩票中心API调试窗口
// @Description 可以直接调用api,所有aes加密数据清到/aes模块进行加密
// @Author roland
package routers

import (
	"NewLottApi/controllers"
	"NewLottApi/models"
	"common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/beego/i18n"
)

var logpath string = ""  //日志保存目录
var clientIP string = "" //请求端ip

func init() {

	ns := beego.NewNamespace("/v1",

		beego.NSBefore(FilterPublic),
		beego.NSNamespace("/aes", //只有商户调用接口会加密。
			beego.NSInclude(
				&controllers.AesController{},
			),
		),
		beego.NSNamespace("/user", //商户对接调用接口统一采用AES加密后传递数据。
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/public",
			beego.NSInclude(
				&controllers.PublicController{},
			),
		),
		beego.NSNamespace("/xml",
			beego.NSInclude(
				&controllers.XMLController{},
			),
		),
		beego.NSNamespace("/game", //平台内部调用的接口无需加密传递数据
			beego.NSInclude(
				&controllers.GameController{},
				&controllers.TraceController{},
			),
		),
		beego.NSNamespace("/admin", //后台调用接口
			beego.NSInclude(
				&controllers.AdminApiController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

/*
 * 钩子公用
 */
var FilterPublic = func(ctx *context.Context) {

	chk := 900
	msg := "请求太频繁!"
	d := map[string]interface{}{}

	//设置日志保存位置
	logpath = controllers.Logpath

	//开启日志
	common.SetLogPath(logpath)

	//將ip質料保存到數據庫中，方便以後使用
	clientIP = "127.0.0.1"

	//读取body所有内容
	io := ctx.Input.RequestBody
	sBody := string(io)

	//进行防注入判断和cc攻击过滤
	ChkDebug := beego.AppConfig.String("chk_debug")
	if ChkDebug == "on" {

		isDanger := models.ChkDanger(sBody)
		if isDanger || models.ChkAccessNumber(clientIP) {

			if isDanger {
				chk = 901
				msg = "可疑的请求，包含非法字符"
			}
			renderJson(ctx, chk, msg, d)

			//包含注入关键词
			common.LogsWithFileName(logpath, "hacker_", "====clientIP==="+clientIP+"\n"+"====sBody==="+sBody+"\n"+"====msg==="+msg) //黑客日志
		}
	}

	setLangVer(ctx) //设置多语言
}

/**
 * 返回json给请求端
 */
func renderJson(ctx *context.Context, status int, message string, d interface{}) {
	//将数据装载到json返回值
	out := controllers.JsonOut{status, message, d}
	ctx.Output.JSON(out, true, true)
}

// setLangVer sets site language version.
func setLangVer(ctx *context.Context) bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	//	lang := this.Input().Get("lang")
	lang := ctx.Request.FormValue("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		//		lang = this.Ctx.GetCookie("lang")
		lang = ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		//		al := this.Ctx.Request.Header.Get("Accept-Language")
		al := ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is 中文.
	if len(lang) == 0 {
		lang = "zh-CN"
		//		lang = "en-US"
		isNeedRedir = false
	}
	lang = "zh-CN" //写死

	curLang := langType{
		Lang: lang, //语言
	}

	// Save language information in cookies.
	if !hasCookie {
		//		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
		ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	for _, v := range langTypes {
		if lang == v.Lang {
			curLang.Name = v.Name //语言说明
		}
	}

	// Set language properties.
	controllers.BaseControllerX.SetLang(curLang.Lang)

	return isNeedRedir
}
