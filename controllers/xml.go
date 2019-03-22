package controllers

//import (
//	"NewLottApi/models"
//)

//XMLController xml操作
type XMLController struct {
	MainController
}

// @Title SaveXML
// @Description 保存一个xml文件
// @Param	filename		path 	string	true		"xml的存储名称"
// @Param	body		body 	models.XMLBody	true		"xml content"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /saveXML/:filename [post]
//func (o *XMLController) SaveXML() {
//	dir := "xmls"
//	filename := o.Ctx.Input.Param(":filename")
//	xmlContent := string(o.Ctx.Input.RequestBody)
//	err := models.SaveXML(dir, filename, xmlContent)
//	if errdata := o.CheckError(xmlContent, err, "save xml Error"); errdata != nil {
//		o.Data["json"] = errdata
//	} else {
//		var result JsonOut
//		result.Status = 200
//		result.Msg = "SaveXML Complete"
//		result.Data = nil
//		o.Data["json"] = errdata
//	}
//	o.ServeJSON()
//}
