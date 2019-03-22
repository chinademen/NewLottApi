package controllers

import (
	"NewLottApi/models"
	"strings"
)

//AesController AesController
type AesController struct {
	MainController
}

// @Title AesEncryptString
// @Description 用来产生aes数据
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入商 比如TEST001"
// @Param	params		query 	controllers.ParamsInputType		true		"param参数的明文"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /encrypt.do [post]
func (c *AesController) AesEncryptString() {

	//定义需要输出的数据格式
	var reStatus int = 200
	var reMsg string = "ok"
	d := map[string]interface{}{}

	c.FilterUser()

	paramStr := string(c.GetString("params"))
	aesStr := models.HookAES.AesEncryptString(paramStr)
	aesStr = strings.Replace(aesStr, "+", "%2B", -1)

	d["encode"] = aesStr
	c.RenderJson(reStatus, reMsg, d)
}

//AesDecryptString 用来产生aes数据
func (u *AesController) AesDecryptString() {
	decryptString := string(u.Ctx.Input.RequestBody)
	merchantStr := u.GetString("merchant")
	models.SetPrivateKey(merchantStr)
	aesStr := models.HookAES.AesDecryptString(decryptString)
	u.Data["json"] = aesStr
	u.ServeJSON()

}
