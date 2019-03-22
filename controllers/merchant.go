package controllers

import (
	"NewLottApi/models"

	"github.com/astaxie/beego"
)

type MerchantController struct {
	beego.Controller
}

func (this *MerchantController) Post() {
	sIdentity := this.GetString("identity")

	accessid := models.Merchants.GetId(sIdentity)

	//定义需要输出的数据格式
	d := map[string]interface{}{}
	d["accessid"] = accessid

	//将数据装载到json返回值
	res := JsonOut{200, "接入商正常", d}

	if len(accessid) < 1 {
		res = JsonOut{301, "接入商不存在", d}
	}

	//输出json数据
	this.Data["json"] = &res
	this.ServeJSON()
}
