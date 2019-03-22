package controllers

import (
	"NewLottApi/dao"
	"NewLottApi/models"
	"encoding/json"
	"fmt"
	Ldefined "lotteryJobs/defined"
	Lmodels "lotteryJobs/models"
	"strconv"
)

//TraceController 追号类
type TraceController struct {
	MainController
}

// @Title trace_list
// @Description 追号列表
// @Param	params	query	controllers.ParamsInputType	true		"created_at_from=开始时间&created_at_to=结束时间&page=页&lottery_id=彩种id&status=状态&start_issue=开始奖期"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /trace_list.do [post]
func (c *TraceController) GetTraceListServer() {

	var reStatus int = 500
	var reMsg string = "no"
	var d = map[string]interface{}{}

	webLoginMap := c.WebLogin()
	merchantId := webLoginMap["merchant_id"]
	userId := webLoginMap["user_id"]

	startAt := webLoginMap["created_at_from"] //开始时间
	endAt := webLoginMap["created_at_to"]     //结束时间
	lotteryId := webLoginMap["lottery_id"]    //彩种ID
	status := webLoginMap["status"]           //注单状态
	startIssue := webLoginMap["start_issue"]  //开始奖期
	page := webLoginMap["page"]               //第几页
	rows := webLoginMap["rows"]               //每页取多少行

	sWhere := fmt.Sprintf("merchant_id='%s' and user_id='%s'", merchantId, userId) //sql条件 必要

	//如果有时间搜索
	if len(startAt) > 0 || len(endAt) > 0 {

		var sTime int64
		var eTime int64

		if len(startAt) > 0 {
			//检查开始时间
			if reStatus, reMsg, sTime = models.ChkIsDate(startAt); reStatus != 200 {

				//输出json数据
				c.RenderJson(reStatus, reMsg, d)
			}

			sWhere += fmt.Sprintf(" and created_at >= '%s'", startAt) //sql条件 非必要
		}

		if len(endAt) > 0 {

			//检查结束时间
			if reStatus, reMsg, eTime = models.ChkIsDate(endAt); reStatus != 200 {
				//输出json数据
				c.RenderJson(reStatus, reMsg, d)
			}

			sWhere += fmt.Sprintf(" and created_at <= '%s'", endAt) //sql条件 非必要
		}

		if eTime < sTime {
			//输出json数据
			c.RenderJson(501, "结束时间必须大于开始时间", d)
		}
	}

	if len(lotteryId) > 0 {
		sWhere += fmt.Sprintf(" and lottery_id='%s' ", lotteryId)
	}

	if len(status) > 0 {
		sWhere += fmt.Sprintf(" and status='%s' ", status)
	}

	if len(startIssue) > 0 {
		sWhere += fmt.Sprintf(" and start_issue='%s' ", startIssue)
	}

	oTraces := Lmodels.Trace.GetAllLists(sWhere, page, rows) //转账记录
	count := Lmodels.Trace.GetCount(sWhere)

	reData := []map[string]interface{}{}
	if len(oTraces) > 0 {

		for _, oTrace := range oTraces {

			oLottery := Lmodels.Lottery.GetInfo(oTrace["lottery_id"])
			sLottery := BaseControllerX.GetLang(oLottery["name"])

			prize := "--"
			if prizeV, prizeOk := oTrace["prize"]; prizeOk == true {
				prize = prizeV
			}

			reMap := map[string]interface{}{}

			reMap["id"] = oTrace["id"]
			reMap["user_id"] = oTrace["user_id"]
			reMap["username"] = oTrace["username"]
			reMap["terminal_id"] = oTrace["terminal_id"] //终端id

			reMap["serial_number"] = Lmodels.Trace.GetSerialNumberShortAttribute(oTrace["serial_number"]) //追号编码

			reMap["prize_group"] = oTrace["prize_group"] //奖金组
			reMap["lottery_id"] = oTrace["lottery_id"]
			reMap["lottery"] = sLottery
			reMap["total_issues"] = oTrace["total_issues"]
			reMap["finished_issues"] = oTrace["finished_issues"]
			reMap["canceled_issues"] = oTrace["canceled_issues"]
			reMap["stop_on_won"] = oTrace["stop_on_won"]
			reMap["start_issue"] = oTrace["start_issue"]
			reMap["way_id"] = oTrace["way_id"]
			reMap["way"] = oTrace["title"]
			reMap["bet_number"] = oTrace["display_bet_number"]
			reMap["coefficient"] = oTrace["coefficient"]
			reMap["single_amount"] = oTrace["single_amount"]
			reMap["amount"] = oTrace["amount"]
			reMap["prize"] = prize
			reMap["status"] = Lmodels.Trace.GetFormattedStatusAttribute(oTrace["status"])
			reMap["bought_at"] = oTrace["bought_at"]
			reMap["updated_at"] = oTrace["updated_at"]
			reMap["finished_amount_formatted"] = oTrace["finished_amount"]
			reMap["canceled_amount_formatted"] = oTrace["canceled_amount"]

			reData = append(reData, reMap)
		}
	}

	reStatus = 200
	reMsg = "yes"

	d["count"] = count
	d["list"] = reData

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}

// @Title trace_detailt
// @Description 追号详情
// @Param	params	query	controllers.ParamsInputType	true		"id=追号ID"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /trace_detailt.do [post]
func (c *TraceController) GetTraceDetailServer() {

	var reStatus int = 500
	var reMsg string = "no"
	var d = map[string]interface{}{}

	webLoginMap := c.WebLogin()

	userId := webLoginMap["user_id"]

	id := webLoginMap["id"] //追号ID

	if len(id) == 0 {
		c.RenderJson(501, "id不能为空", d) //输出json数据
	}

	oTrace := Lmodels.Trace.GetInfo(id)
	if len(oTrace) == 0 {
		c.RenderJson(502, "id有误", d) //输出json数据
	}

	if oTrace["user_id"] != userId {
		c.RenderJson(503, "无权查看此订单", d) //输出json数据
	}

	oLottery := Lmodels.Lottery.GetInfo(oTrace["lottery_id"])

	statusRunning := 0
	canBeCanceled := 0
	if oTrace["status"] == Ldefined.TRACE_STATUS_RUNNING {
		statusRunning = 1
		canBeCanceled = 1
	}

	d["id"] = oTrace["id"]
	d["username"] = oTrace["username"]
	d["terminal_id"] = oTrace["terminal_id"]
	d["serial_number"] = oTrace["serial_number"]
	d["prize_group"] = oTrace["prize_group"]
	d["lottery"] = BaseControllerX.GetLang(oLottery["name"])
	d["lottery_identifier"] = oLottery["identifier"]
	d["total_issues"] = oTrace["total_issues"]

	d["finished_issues"] = oTrace["finished_issues"]
	d["canceled_issues"] = oTrace["canceled_issues"]
	d["stop_on_won"] = oTrace["stop_on_won"]
	d["formatted_stop_on_won"] = BaseControllerX.GetLang(Lmodels.Trace.GetFormattedStopOnWonAttribute(oTrace["stop_on_won"]))
	d["start_issue"] = oTrace["start_issue"]
	d["way"] = oTrace["title"]

	d["bet_number"] = oTrace["display_bet_number"]
	d["formatted_coefficient"] = Lmodels.Trace.GetFormattedCoefficientAttribute(oTrace["coefficient"])

	d["single_amount"] = oTrace["single_amount"]
	d["amount_formatted"] = oTrace["amount"]
	d["finished_amount_formatted"] = oTrace["finished_amount"]
	d["canceled_amount_formatted"] = oTrace["canceled_amount"]
	d["status_running"] = statusRunning
	d["prize"] = oTrace["prize"]
	d["formatted_status"] = Lmodels.Trace.GetFormattedStatusAttribute(oTrace["status"])
	d["status"] = oTrace["status"]
	d["bought_at"] = oTrace["bought_at"]
	d["can_be_canceled"] = canBeCanceled

	reStatus = 200
	reMsg = "yes"

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}

// @Title cancel_issue_trace
// @Description 取消某期追号
// @Param	params	query	controllers.ParamsInputType	true		"trace_id=关联追号ID&ids=[1,2,3,4,5]追号详情ID"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_issue_trace.do [post]
func (c *TraceController) CancelTraceReserveServer() {

	var reStatus int = 500
	var reMsg string = "no"
	var d = map[string]interface{}{}

	webLoginMap := c.WebLogin()

	userId := webLoginMap["user_id"] //token 获取的用户id

	traceId := webLoginMap["trace_id"] //追号ID
	ids := webLoginMap["ids"]          //追号详情ID

	if len(traceId) == 0 {
		c.RenderJson(501, "追号id不能为空", d) //输出json数据
	}

	if len(ids) == 0 {
		c.RenderJson(502, "期数id不能为空", d) //输出json数据
	}

	aDetailIds := []int{}
	aDetailIdErr := json.Unmarshal([]byte(ids), &aDetailIds)
	if aDetailIdErr != nil || len(aDetailIds) == 0 {
		c.RenderJson(503, "追号详情id有误", d) //输出json数据
	}

	reStatus, reMsg = Lmodels.Trace.CancelDetail(traceId, userId, Ldefined.TRACE_DETAIL_STATUS_USER_CANCELED, aDetailIds)

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}

// @Title cancel_trace
// @Description 终止追号
// @Param	params	query	controllers.ParamsInputType	true		"id=追号ID"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_trace.do [post]
func (c *TraceController) StopTraceServer() {

	var reStatus int = 500
	var reMsg string = "no"
	var d = map[string]interface{}{}

	webLoginMap := c.WebLogin()

	userId := webLoginMap["user_id"] //token 获取的用户id

	traceId := webLoginMap["id"] //追号ID

	if len(traceId) == 0 {
		c.RenderJson(501, "追号id不能为空", d) //输出json数据
	}

	oTrace := Lmodels.Trace.GetInfo(traceId)
	if len(oTrace) == 0 {
		c.RenderJson(502, "追号id有误", d) //输出json数据
	}

	if oTrace["status"] != Ldefined.TRACE_STATUS_RUNNING {
		c.RenderJson(503, "订单状态不可取消", d) //输出json数据
	}

	if oTrace["user_id"] != userId {
		c.RenderJson(504, "无权操作此订单", d) //输出json数据
	}

	//锁账户
	accountId := oTrace["account_id"]
	bIsLock, _ := Lmodels.Account.Lock(accountId)
	if bIsLock == false {
		c.RenderJson(505, "事务占用", d) //输出json数据
	}

	//开启事务
	db := new(Lmodels.Table)
	Tx := db.BeginTransaction()

	//终止追号
	bIsSuccess, _ := Lmodels.Trace.Terminate(Tx, oTrace, Ldefined.TRACE_STATUS_USER_STOPED)
	if bIsSuccess == true {

		comErr := Tx.Commit() //提交事务
		if comErr == nil {

			reStatus = 200
			reMsg = "终止追号单成功"
		} else {
			reStatus = 506
			reMsg = "提交事务失败"
		}

	} else {
		Tx.Rollback() //回滚事务
		reStatus = 507
		reMsg = "终止追号单失败"
	}
	Lmodels.Account.UnLock(accountId)

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}

// @Title GetTraceProjectDetail
// @Description 获得追号注单详情,author(leon)
// @Param	params	query	controllers.ParamsInputType	true		"id=追号ID"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /trace_project_detail.do [post]
func (this *TraceController) GetTraceProjectDetail() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			this.RenderJson(500, "系统错误", err)
		}

	}()

	Status := 200
	Msg := "success"
	var result interface{}

	webLoginMap := this.WebLogin()

	sUserId := webLoginMap["user_id"]
	iUserId, _ := strconv.Atoi(sUserId)
	oUser, _ := dao.GetUsersById(iUserId) //获取用户数据
	if oUser.Id == 0 {
		Status = 501
		Msg = "用戶信息错误"
		this.RenderJson(Status, Msg, result)
	}

	sTraceId := webLoginMap["trace_id"]
	iTraceId, _ := strconv.Atoi(sTraceId)
	oTrace, err := dao.GetTracesById(iTraceId)
	if oTrace.Id == 0 || err != nil {
		Status = 7001
		Msg = "沒有当前追号信息"
		this.RenderJson(Status, Msg, err)
	}

	if oTrace.UserId != iUserId {
		Status = 7002
		Msg = "没有权限查看订单"
		this.RenderJson(Status, Msg, result)
	}

	limitRow := 5                    //默认显示5条
	limitPage := webLoginMap["page"] //获取分页

	//分页
	limitPageInt := 1
	if len(limitPage) > 0 {

		pageInt, pageErr := strconv.Atoi(limitPage)
		if pageErr == nil {
			limitPageInt = pageInt
		}
	}
	limitOff := (limitPageInt - 1) * limitRow

	sSql := fmt.Sprintf("select * from trace_details where trace_id = %s order by issue asc limit %d, %d", sTraceId, limitOff, limitRow)
	aTraceDetail := Lmodels.TraceDetail.GetSqlQuery(sSql)
	if len(aTraceDetail) < 1 {
		Status = 7003
		Msg = "追号数据有误"
		this.RenderJson(Status, Msg, err)
	}

	//获取总数
	sTotalSql := fmt.Sprintf("select count(0) as total from trace_details where trace_id = %s limit 1", sTraceId)
	qTotal := Lmodels.TraceDetail.GetSqlQuery(sTotalSql)
	total := 0
	if sTotal, totalOk := qTotal[0]["total"]; totalOk == true {

		totalInt, totalErr := strconv.Atoi(sTotal)
		if totalErr == nil {
			total = totalInt
		}
	}

	var aArr []map[string]string
	for _, mTraceDetail := range aTraceDetail {
		iStatus, _ := strconv.Atoi(mTraceDetail["status"])
		if iStatus == models.STATUS_WAITING {
			mTraceDetail["status_waiting"] = "1"
			mTraceDetail["can_be_canceled"] = "1"
		} else {
			mTraceDetail["status_waiting"] = "0"
			mTraceDetail["can_be_canceled"] = "0"
		}
		mTraceDetail["display_status"] = models.TraceDetailStatus[iStatus]
		aArr = append(aArr, mTraceDetail)

	}
	var mResault = map[string]interface{}{
		"list":  aArr,
		"count": total,
	}
	this.RenderJson(Status, Msg, mResault)
}

// @Title SelectLotteries
// @Description
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /get_select_lottery.do [post]
func (c *TraceController) SelectLotteries() {

	var reStatus int = 500
	var reMsg string = "no"
	var d = map[string]interface{}{}

	webLoginMap := c.WebLogin()

	//	userId := webLoginMap["user_id"] //token 获取的用户id
	terminal := webLoginMap["terminal"] //设备类型, 1-PC 2-手机

	lotteryId := webLoginMap["lottery_id"]    //筛选 彩票信息ID
	seriesId := ""                            //前台传的lottery_id,通过lottery_id得到series_id 来获取玩法组信息
	groupId := webLoginMap["group_id"]        //筛选 玩法组ID
	groupIdOk := ""                           //当筛选group_id正确,才会有值
	groupWayId := webLoginMap["group_way_id"] //筛选 玩法ID

	//所有彩票信息,选中的会多出一列selected=selected
	ReLotteryMap := []map[string]string{}

	lotteryMap := Lmodels.Lottery.GetLists()
	for _, lotteryRow := range lotteryMap {

		ReLotteryRow := map[string]string{
			"id":        lotteryRow["id"],
			"name":      lotteryRow["name"],
			"status":    lotteryRow["status"],
			"series_id": lotteryRow["series_id"],
		}

		if len(lotteryId) > 0 && lotteryId == lotteryRow["id"] {
			ReLotteryRow["selected"] = "selected"
			seriesId = lotteryRow["series_id"]
		}

		ReLotteryRow["name_cn"] = BaseControllerX.GetLang(lotteryRow["name"])
		ReLotteryMap = append(ReLotteryMap, ReLotteryRow)
	}

	//玩法组信息,选中的会多出一列selected=selected
	ReGroupMap := []map[string]string{}
	if len(seriesId) > 0 {

		groupMap := Lmodels.WayGroup.RGetListsBySeriesId(seriesId, terminal)
		for _, groupRow := range groupMap {

			ReGroupRow := map[string]string{
				"id":        groupRow["id"],
				"title":     groupRow["title"],
				"series_id": groupRow["series_id"],
			}

			if len(groupId) > 0 && groupId == groupRow["id"] {
				ReGroupRow["selected"] = "selected"
				groupIdOk = groupRow["id"]
			}

			ReGroupMap = append(ReGroupMap, ReGroupRow)
		}
	}

	//玩法组分类信息
	ReGroupClassMap := []map[string]string{}
	if len(groupIdOk) > 0 {
		groupClassMap := Lmodels.WayGroup.RGetListsByParentId(seriesId, terminal, groupIdOk)
		for _, groupClassRow := range groupClassMap {

			ReGroupClassRow := map[string]string{
				"id":        groupClassRow["id"],
				"title":     groupClassRow["title"],
				"series_id": groupClassRow["series_id"],
			}

			if len(groupWayId) > 0 && groupWayId == groupClassRow["id"] {
				ReGroupClassRow["selected"] = "selected"
			}

			ReGroupClassMap = append(ReGroupClassMap, ReGroupClassRow)
		}
	}

	//玩法信息,选中的会多出一列selected=selected
	ReGroupWayMap := []map[string]string{}
	if len(ReGroupClassMap) > 0 {

		//通过玩法分类,获取玩法信息
		for _, row := range ReGroupClassMap {

			groupWayMap := Lmodels.WayGroupWays.RGetListsByGroupid(seriesId, terminal, row["id"])
			for _, wayRow := range groupWayMap {

				reWayRow := map[string]string{
					"id":        wayRow["id"],
					"title":     wayRow["title"],
					"series_id": wayRow["series_id"],
				}

				ReGroupWayMap = append(ReGroupWayMap, reWayRow)
			}
		}
	}

	d["lottery"] = ReLotteryMap //select 游戏名称
	d["group"] = ReGroupMap     //select 玩法群
	d["group_class"] = ReGroupClassMap
	d["group_way"] = ReGroupWayMap //select 玩法

	//	WhereLottery := "" //彩票筛选
	//	WhereGroup := ""   //玩法组筛选
	//	WhereWay := ""     //玩法筛选

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}
