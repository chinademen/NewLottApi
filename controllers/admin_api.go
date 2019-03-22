package controllers

import (
	"fmt"
	lotteryJobModels "lotteryJobs/models"
	lotteryJobThread "lotteryJobs/thread"
)

//AdminController 后台调用
type AdminApiController struct {
	MainController
}

// @Title Calculate Prize
// @Description 计奖
// @Param	params		formData 	string	true		"参数 lottery_id=xxx&issue=xxx&new_code=xxx&admin_name=xxx&sign=xxx"
// lottery_id		string			"彩种ID"
// issue		string			"奖期号码"
// admin_name		string			"后台管理员名"
// sign			string			"签名"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /calculate_prize.do [post]
func (c *AdminApiController) CalculatePrize() {
	var reStatus int = 200
	var reMsg string = "计奖任务成功提交"
	d := map[string]interface{}{}

	adminAPIMap := c.AdminAPI()

	//1.验证必传参数
	if _, ok := adminAPIMap["lottery_id"]; !ok {
		reStatus = 403
		reMsg = "缺少彩种参数"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["issue"]; !ok {
		reStatus = 403
		reMsg = "缺少奖期参数"
		c.RenderJson(reStatus, reMsg, d)
	}

	if _, ok := adminAPIMap["admin_name"]; !ok {
		reStatus = 403
		reMsg = "缺少admin_name"
		c.RenderJson(reStatus, reMsg, d)
	}

	sLotteryId := adminAPIMap["lottery_id"]
	sIssue := adminAPIMap["issue"]

	mLottery := lotteryJobModels.Lottery.GetInfo(sLotteryId)
	if len(mLottery) == 0 {
		reStatus = 404
		reMsg = "彩种不存在"
		c.RenderJson(reStatus, reMsg, d)
	}
	mIssues := lotteryJobModels.Issues.GetInfo(sLotteryId, sIssue)
	if len(mIssues) == 0 {
		reStatus = 404
		reMsg = "彩种对应的奖期数据不存在"
		c.RenderJson(reStatus, reMsg, d)
	}

	go lotteryJobThread.CalculatePrize(sLotteryId, sIssue)

	c.RenderJson(reStatus, reMsg, d)

}

// @Title Cancel Prize
// @Description 取消原来号码重计奖
// @Param	params		formData 	string	true		"参数 lottery_id=xxx&issue=xxx&new_code=xxx&admin_name=xxx&sign=xx"
// lottery_id		string			"彩种ID"
// issue		string			"奖期号码"
// new_code		string			"新的开奖号码"
// code_center_id	string			"开奖中心编号"
// admin_name		string			"后台管理员名"
// sign			string			"签名"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_prize.do [post]
func (c *AdminApiController) CancelPrize() {
	var reStatus int = 200
	var reMsg string = "重开奖任务成功提交"
	d := map[string]interface{}{}

	adminAPIMap := c.AdminAPI()

	//1.验证必传参数
	if _, ok := adminAPIMap["lottery_id"]; !ok {
		reStatus = 403
		reMsg = "缺少彩种参数"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["issue"]; !ok {
		reStatus = 403
		reMsg = "缺少奖期参数"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["admin_name"]; !ok {
		reStatus = 403
		reMsg = "缺少admin_name"
		c.RenderJson(reStatus, reMsg, d)
	}

	if _, ok := adminAPIMap["new_code"]; !ok {
		reStatus = 403
		reMsg = "缺少新开奖号码"
		c.RenderJson(reStatus, reMsg, d)
	}

	sLotteryId := adminAPIMap["lottery_id"]
	sIssue := adminAPIMap["issue"]
	sNewCode := adminAPIMap["new_code"]
	sAdminName := adminAPIMap["admin_name"]

	mLottery := lotteryJobModels.Lottery.GetInfo(sLotteryId)
	if len(mLottery) == 0 {
		reStatus = 404
		reMsg = "彩种不存在"
		c.RenderJson(reStatus, reMsg, d)
	}
	mIssues := lotteryJobModels.Issues.GetInfo(sLotteryId, sIssue)
	if len(mIssues) == 0 {
		reStatus = 404
		reMsg = "彩种对应的奖期数据不存在"
		c.RenderJson(reStatus, reMsg, d)
	}

	go lotteryJobThread.CancelPrize(sLotteryId, sIssue, sNewCode, sAdminName)

	c.RenderJson(reStatus, reMsg, d)

}

// @Title Cancel Issue
// @Description 撤销指定奖期内的所有正常注单
// @Param	params		formData 	string	true		"参数 lottery_id=xxx&issue=xxx&begin_time=xxx&admin_name=xxx&sign=xx"
// lottery_id		string			"彩种ID"
// issue		string			"奖期号码"
// begin_time		string			"最早开奖时间"
// admin_name		string			"后台管理员名"
// sign			string			"签名"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_issue.do [post]
func (c *AdminApiController) CancelIssue() {
	var reStatus int = 200
	var reMsg string = "重开奖任务成功提交"
	d := map[string]interface{}{}

	adminAPIMap := c.AdminAPI()

	//1.验证必传参数
	if _, ok := adminAPIMap["lottery_id"]; !ok {
		reStatus = 403
		reMsg = "缺少彩种参数"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["issue"]; !ok {
		reStatus = 403
		reMsg = "缺少奖期参数"
		c.RenderJson(reStatus, reMsg, d)
	}

	if _, ok := adminAPIMap["admin_name"]; !ok {
		reStatus = 403
		reMsg = "缺少admin_name"
		c.RenderJson(reStatus, reMsg, d)
	}

	sLotteryId := adminAPIMap["lottery_id"]
	sIssue := adminAPIMap["issue"]
	sBeginTime := adminAPIMap["begin_time"]
	sAdminName := adminAPIMap["admin_name"]

	mLottery := lotteryJobModels.Lottery.GetInfo(sLotteryId)
	if len(mLottery) == 0 {
		reStatus = 404
		reMsg = "彩种不存在"
		c.RenderJson(reStatus, reMsg, d)
	}
	mIssues := lotteryJobModels.Issues.GetInfo(sLotteryId, sIssue)
	if len(mIssues) == 0 {
		reStatus = 404
		reMsg = "彩种对应的奖期数据不存在"
		c.RenderJson(reStatus, reMsg, d)
	}

	go lotteryJobThread.CancelIssue(sLotteryId, sIssue, sBeginTime, sAdminName)

	c.RenderJson(reStatus, reMsg, d)

}

// @Title Cancel Project
// @Description 撤销指定注单
// @Param	params		formData 	string	true		"参数 project_id=xxx&admin_name=xxx&sign=xx"
// project_id		string			"注单ID"
// admin_name		string			"后台管理员名"
// sign			string			"签名"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_project.do [post]
func (c *AdminApiController) CancelProject() {
	var reStatus int = 200
	var reMsg string = "注单撤单任务成功提交"
	d := map[string]interface{}{}

	adminAPIMap := c.AdminAPI()

	//1.验证必传参数
	if _, ok := adminAPIMap["project_id"]; !ok {
		reStatus = 403
		reMsg = "缺少注单id"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["admin_name"]; !ok {
		reStatus = 403
		reMsg = "缺少admin_name"
		c.RenderJson(reStatus, reMsg, d)
	}

	sProjectId := adminAPIMap["project_id"]
	sAdminName := adminAPIMap["admin_name"]
	mProject := lotteryJobModels.Project.GetInfo(sProjectId)
	if len(mProject) == 0 {
		reStatus = 404
		reMsg = "注单不存在"
		c.RenderJson(reStatus, reMsg, d)
	}

	go lotteryJobThread.StopProject([]string{sProjectId}, sAdminName)

	c.RenderJson(reStatus, reMsg, d)

}

// @Title Cancel Trace
// @Description 终止指定追号
// @Param	params		formData 	string	true		"参数 project_id=xxx&admin_name=xxx&sign=xx"
// trace_id		string			"追号ID"
// admin_name		string			"后台管理员名"
// sign			string			"签名"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_trace.do [post]
func (c *AdminApiController) CancelTrace() {
	var reStatus int = 200
	var reMsg string = "撤销追号任务成功提交"
	d := map[string]interface{}{}

	adminAPIMap := c.AdminAPI()

	//1.验证必传参数
	if _, ok := adminAPIMap["trace_id"]; !ok {
		reStatus = 403
		reMsg = "缺少追号id"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["admin_name"]; !ok {
		reStatus = 403
		reMsg = "缺少admin_name"
		c.RenderJson(reStatus, reMsg, d)
	}

	sTraceId := adminAPIMap["trace_id"]
	sAdminName := adminAPIMap["admin_name"]
	mTraceId := lotteryJobModels.Trace.GetInfo(sTraceId)
	if len(mTraceId) == 0 {
		reStatus = 404
		reMsg = "追号任务不存在"
		c.RenderJson(reStatus, reMsg, d)
	}

	go lotteryJobThread.StopTrace([]string{sTraceId}, sAdminName)

	c.RenderJson(reStatus, reMsg, d)

}

// @Title Make Issue Cache
// @Description 生成奖期缓存
// @Param	params		formData 	string	true		"参数 lottery_id=xxx&admin_name=xxx&sign=xx"
// lottery_id		string			"彩种ID"
// admin_name		string			"后台管理员名"
// sign			string			"签名"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /make_issue_cache.do [post]
func (c *AdminApiController) MakeIssueCache() {
	var reStatus int = 200
	var reMsg string = "奖期缓存重建任务成功提交"
	d := map[string]interface{}{}

	adminAPIMap := c.AdminAPI()

	//1.验证必传参数
	if _, ok := adminAPIMap["lottery_id"]; !ok {
		reStatus = 403
		reMsg = "缺少彩种id"
		c.RenderJson(reStatus, reMsg, d)
	}
	if _, ok := adminAPIMap["admin_name"]; !ok {
		reStatus = 403
		reMsg = "缺少admin_name"
		c.RenderJson(reStatus, reMsg, d)
	}

	sAdminName := adminAPIMap["admin_name"]
	if debug {
		fmt.Println("====sAdminName===", sAdminName)
	}

	go lotteryJobThread.MakeIssueListCache(0)

	c.RenderJson(reStatus, reMsg, d)

}
