package controllers

import (
	"NewLottApi/dao"
	"NewLottApi/lib"
	"NewLottApi/log"
	"NewLottApi/models"
	"common"
	"encoding/json"
	"fmt"
	"lotteryJobs/base/zizhu"
	Ldefined "lotteryJobs/defined"
	Lmodels "lotteryJobs/models"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//GameController 投注类
type GameController struct {
	MainController
}

type Bet struct {
	maxBetGroup int
	minBetGroup int
	betTime     int64
	prizeLimit  int64
	betRecordId int64
	IsTrace     bool
	WinTrueStop bool
	clientIP    string
	proxyIP     string
	BetIssue    map[string]string
	Lottery     *dao.Lotteries
	Series      *dao.Series
	User        *dao.Users
	Account     *dao.Accounts
	Balls       []*models.BetData
}

// @Title GetProjectListServer
// @Description 获取投注记录列表
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /bet_record_list.do [post]
func (c *GameController) GetProjectListServer() {

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
	issue := webLoginMap["issue"]             //期号
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

			sWhere += fmt.Sprintf(" and bought_at >= '%s'", startAt) //sql条件 非必要
		}

		if len(endAt) > 0 {

			//检查结束时间
			if reStatus, reMsg, eTime = models.ChkIsDate(endAt); reStatus != 200 {
				//输出json数据
				c.RenderJson(reStatus, reMsg, d)
			}

			//sql条件 非必要
			sWhere += fmt.Sprintf(" and bought_at <= '%s'", endAt)
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

	if len(issue) > 0 {
		sWhere += fmt.Sprintf(" and issue='%s' ", issue)
	}

	sqlPage := 1 //默认第一页
	if len(page) > 0 {
		pageInt, pageErr := strconv.Atoi(page)
		if pageErr == nil {
			sqlPage = pageInt
		}
	}

	sqlRow := 20 //默认每页取20行
	if len(rows) > 0 {
		rowsInt, rowsErr := strconv.Atoi(rows)
		if rowsErr == nil {
			sqlRow = rowsInt
		}
	}

	iOffset := (sqlPage - 1) * sqlRow

	oProjects := Lmodels.Project.GetAllLists(sWhere, iOffset, sqlRow)
	count := Lmodels.Project.GetCount(sWhere)

	reData := []map[string]interface{}{}
	if len(oProjects) > 0 {
		aSeriesWays := Lmodels.SeriesWay.IdTitleList() //所有转账类型
		for _, oProject := range oProjects {

			oLottery := Lmodels.Lottery.GetInfo(oProject["lottery_id"])
			sLottery := BaseControllerX.GetLang(oLottery["name"])

			traceId := "--"
			if len(oProject["trace_id"]) > 0 {
				traceId = oProject["trace_id"]
			}

			seriesWaysTitle := "--"
			if seriesWaysTitleV, seriesWaysTitleOk := aSeriesWays[oProject["way_id"]]; seriesWaysTitleOk == true {
				seriesWaysTitle = seriesWaysTitleV
			}

			winningNumber := "--"
			if len(oProject["winning_number"]) > 0 {
				winningNumber = oProject["winning_number"]
			}

			betNumber := oProject["display_bet_number"]
			if len(betNumber) > 5 {
				betNumber = Lmodels.Project.GetDisplayBetNumberForListAttribute(betNumber)
			}

			prize := "--"
			if len(oProject["prize"]) > 0 {
				prize = oProject["prize"]
			}

			reMap := map[string]interface{}{}

			reMap["id"] = oProject["id"]
			reMap["user_id"] = oProject["user_id"]
			reMap["username"] = oProject["username"]
			reMap["terminal_id"] = oProject["terminal_id"]
			reMap["serial_number"] = Lmodels.GetSerialNumberShortAttribute(oProject["serial_number"])
			reMap["trace_id"] = traceId
			reMap["prize_group"] = oProject["prize_group"]
			reMap["lottery_id"] = oProject["lottery_id"]
			reMap["lottery"] = sLottery
			reMap["issue"] = oProject["issue"]
			reMap["way_id"] = oProject["way_id"]
			reMap["way"] = seriesWaysTitle //通过way_id 得到series_ways表 name
			reMap["bet_number"] = betNumber
			reMap["multiple"] = oProject["multiple"]
			reMap["coefficient"] = oProject["coefficient"]
			reMap["amount"] = oProject["amount"]
			reMap["winning_number"] = winningNumber
			reMap["prize"] = prize
			reMap["status"] = BaseControllerX.GetLang(Lmodels.Project.GetFormattedStatusAttribute(oProject["status"]))
			reMap["bought_at"] = oProject["bought_at"]
			reMap["is_traced"] = Lmodels.Project.GetTraceIdAttribute(oProject["trace_id"])

			reData = append(reData, reMap)

		}
	}

	d["list"] = reData
	d["count"] = count

	reStatus = 200
	reMsg = "yes"

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)

}

// @Title Bet
// @Description 根据前端提交过来的用户投注数据完成下注功能
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /bet.do [post]
func (this *GameController) Bet() {

	Status := 200
	Msg := "success"
	d := map[string]interface{}{}

	webLoginMap := this.WebLogin()

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			Status = 500
			Msg = "系统错误"

			//回答日志
			sErr := common.InterfaceToString(err)
			if len(sErr) > 1 {
				log.AddLog("", "response", "user_id:"+webLoginMap["user_id"]+"==>data:"+sErr, log.BET)
			}
			this.RenderJson(Status, Msg, err)
		}

	}()

	///////////////////
	//初始化用戶投注數據//
	///////////////////
	sUserId := webLoginMap["user_id"]
	iUserId, _ := strconv.Atoi(sUserId)
	oUser, _ := dao.GetUsersById(int(iUserId))

	oBet := new(Bet)
	if oUser.Id == 0 {
		Status = 4001
		Msg = "用戶不存在"
		this.RenderJson(Status, Msg, d)
	}
	if int(oUser.Blocked) == models.BLOCK_BUY || int(oUser.Blocked) == models.BLOCK_FUND_OPERATE {
		Status = 4002
		Msg = "该用户无法投注"
		this.RenderJson(Status, Msg, d)
	}
	oBet.User = oUser

	//檢查用戶賬戶
	oAccount := models.Accounts.GetAccountInfoByUserId(fmt.Sprintf("%d", oUser.Id))
	if oAccount.Id == 0 {
		Status = 4003
		Msg = "该用户账户错误"
		this.RenderJson(Status, Msg, d)
	}

	oBet.Account = &oAccount

	sIp := webLoginMap["client_ip"]
	oBet.clientIP = sIp
	oBet.proxyIP = sIp

	//判斷是否開啓壓縮
	sOrin := webLoginMap["betdata"]
	sIsCompress := models.SysConfig.GetPrizeByItem("bet_submit_compress")
	if sIsCompress == "1" {
		sOrin = common.DeCodeBetData(webLoginMap["betdata"], 3)
	}

	//請求日志
	log.AddLog("", "request", "user_id:"+sUserId+"==>data:"+sOrin, log.BET)
	if debug {
		fmt.Println("*********************")
		fmt.Println("解壓後字符串", sOrin)
		fmt.Println("*********************")
	}

	////////////////////
	//結構體解析json方法//
	///////////////////
	var oData = new(models.All)
	json.Unmarshal([]byte(sOrin), &oData)

	sGameId := fmt.Sprintf("%d", oData.GameId)                 //彩种id
	sIsTrace := fmt.Sprintf("%d", oData.IsTrace)               //是否追号投注： 1-追号 0-不追号
	aBalls := oData.Balls                                      //玩法等投主内容数组
	sAmount := oData.Amount                                    //投注金额
	sTraceWinStop := fmt.Sprintf("%d", oData.TraceWinStop)     //是否中奖后停止追号1-是 0-否
	sTraceStopValue := fmt.Sprintf("%d", oData.TraceStopValue) //旧版预留
	var mBetIssue = map[string]string{}
	for sIssueCode, iM := range oData.Orders {
		mBetIssue[sIssueCode] = strconv.Itoa(iM)
	}

	oBet.BetIssue = mBetIssue
	oBet.Balls = aBalls

	/////////////////////////////////////
	//解析json方法2(如果結構體解析數據失敗)//
	////////////////////////////////////
	var mRequestData = map[string]interface{}{}
	err := json.Unmarshal([]byte(sOrin), &mRequestData)
	if err != nil {
		Status = 4004
		Msg = "错误请求"
		this.RenderJson(Status, Msg, err)
	}

	ziZhuInfoArr := []map[string]string{}   //自主开奖-选号信息;彩种id,投注号码,玩法,单注金额...
	ziZhuOrdersArr := []map[string]string{} //自主开奖-追号信息;奖期和倍数
	var aBallDatas []*models.BetData
	if len(oBet.BetIssue) < 1 || len(oBet.Balls) < 1 || oData.GameId == 0 || len(sAmount) < 1 || sTraceWinStop == "0" || oData.Multiple < 1 {
		for sKey, value := range mRequestData {
			sValue := common.InterfaceToString(value)
			switch sKey {
			case "amount":
				sAmount = sValue
			case "gameId":
				sGameId = sValue
			case "isTrace":
				sIsTrace = sValue
			case "traceWinStop":
				sTraceWinStop = sValue
			case "sTraceStopValue":
				sTraceStopValue = sValue
			case "orders":
				mOrders := value.(map[string]interface{})
				var mIssues = map[string]string{}
				for sIssueCode, multiple := range mOrders {

					multipleStr := common.InterfaceToString(multiple)
					mIssues[sIssueCode] = multipleStr

					ziZhuOrders := map[string]string{
						"issue":    sIssueCode,  //奖期
						"multiple": multipleStr, //倍数
					}
					ziZhuOrdersArr = append(ziZhuOrdersArr, ziZhuOrders)
				}
				oBet.BetIssue = mIssues
			case "balls":
				aBallData := value.([]interface{})
				for _, ball := range aBallData {

					mBallData := ball.(map[string]interface{})
					oBall := new(models.BetData)

					for sKeyBall, v := range mBallData {
						sV := common.InterfaceToString(v)
						switch sKeyBall {
						case "ball":
							oBall.Ball = sV
						case "moneyunit":
							oBall.Moneyunit = sV
						case "multiple":
							oBall.Multiple = sV
						case "num":
							iNum, _ := strconv.Atoi(sV)
							oBall.Num = iNum
						case "onePrice":
							iOnePrice, _ := strconv.Atoi(sV)
							oBall.OnePrice = iOnePrice
						case "prizeGroup":
							oBall.PrizeGroup = sV
						case "type":
							oBall.Type = sV
						case "viewBalls":
							oBall.ViewBalls = sV
						case "wayId":
							oBall.WayId = sV
						case "jsId":
							iJsId, _ := strconv.Atoi(sV)
							oBall.JsId = iJsId
						case "extra":
							oExt := new(models.Ext)
							mExt := v.(map[string]interface{})
							if len(mExt) > 1 {
								if position, ok := mExt["position"]; ok {
									oExt.Position = common.InterfaceToString(position)
								}
								if seat, ok := mExt["seat"]; ok {
									oExt.Seat = common.InterfaceToString(seat)
								}
							}
							oBall.Extra = oExt
						case "totalMoney":

						case "typeCN":
						case "rebate":
						case "moneyUnitData":
						}
					}
					aBallDatas = append(aBallDatas, oBall)

					ziZhuInfo := map[string]string{
						"balls_betNum":      oBall.Ball,                          //投注号码
						"balls_position":    oBall.Extra.Position,                //万千百十个
						"balls_multiple":    oBall.Multiple,                      //倍数,不追号的时候取这个倍数
						"balls_num":         common.InterfaceToString(oBall.Num), //计算出的注数
						"balls_prizeGroup":  oBall.PrizeGroup,                    //奖金组
						"balls_wayId":       oBall.WayId,                         //玩法id
						"balls_coefficient": oBall.Moneyunit,                     //投注模式; 1=2元, 0.5=1元, 0.1=2角
					}
					ziZhuInfoArr = append(ziZhuInfoArr, ziZhuInfo)

				}
			}
		}
		oBet.Balls = aBallDatas
	}

	if debug {
		fmt.Println("mRequestData-->", mRequestData)
		fmt.Println("###########實例化數據#############")
		fmt.Println("sGameId-->", sGameId)
		fmt.Println("sIsTrace-->", sIsTrace)
		fmt.Println("sAmount-->", sAmount)
		fmt.Println("BetIssue-->", oBet.BetIssue)
		fmt.Println("sTraceWinStop-->", sTraceWinStop)
		fmt.Println("sTraceStopValue-->", sTraceStopValue)
		for _, oBall := range oBet.Balls {
			fmt.Println("==============================")
			fmt.Println("oBall.ball", oBall.Ball)
			fmt.Println("oBall.extra", oBall.Extra)
			fmt.Println("oBall.jsid", oBall.JsId)
			fmt.Println("oBall.Moneyunit", oBall.Moneyunit)
			fmt.Println("oBall.Multiple", oBall.Multiple)
			fmt.Println("oBall.num", oBall.Num)
			fmt.Println("oBall.OnePrice", oBall.OnePrice)
			fmt.Println("oBall.PrizeGroup", oBall.PrizeGroup)
			fmt.Println("oBall.Type", oBall.Type)
			fmt.Println("oBall.ViewBalls", oBall.ViewBalls)
			fmt.Println("oBall.WayId", oBall.WayId)
		}
		fmt.Println("################################")
	}

	///////////
	//數據判斷//
	///////////
	if len(oBet.Balls) < 1 {
		Status = 4005
		Msg = "投注数据错误"
		this.RenderJson(Status, Msg, d)
	}

	//判断彩种是否存在
	oLottery := models.Lotteries.GetInfo(sGameId)
	if len(sGameId) < 1 || oLottery.Id < 1 {
		Status = 4006
		Msg = "彩种信息有误"
		this.RenderJson(Status, Msg, d)
	}
	oBet.Lottery = oLottery

	//如果是測試用戶
	if oUser.IsTester == 1 {
		if int(oLottery.Status) != models.STATUS_AVAILABLE_FOR_TESTER && int(oLottery.Status) != models.STATUS_AVAILABLE {
			Status = 4007
			Msg = "彩种不可用"
			this.RenderJson(Status, Msg, d)
		}
	} else {
		if int(oLottery.Status) == models.STATUS_AVAILABLE_FOR_TESTER {
			Status = 4008
			Msg = "当前彩种针对测试玩家开放"
			this.RenderJson(Status, Msg, d)
		}

		if int(oLottery.Status) == models.STATUS_NOT_AVAILABLE || int(oLottery.Status) == models.STATUS_CLOSED_FOREVER {
			Status = 4009
			Msg = "当前彩种已关闭"
			this.RenderJson(Status, Msg, d)
		}
	}

	//找到系列
	oSeries := models.Series.GetInfo(fmt.Sprintf("%d", oLottery.SeriesId))
	oBet.Series = oSeries

	//判断投注的奖期和号码
	if len(oBet.BetIssue) < 1 {
		Status = 4010
		Msg = "投注信息有误"
		this.RenderJson(Status, Msg, d)
	}

	//判断投注金额
	fAmount, err := strconv.ParseFloat(sAmount, 64)
	if err != nil || fAmount == 0 {
		Status = 4011
		Msg = "投注金额有误"
		this.RenderJson(Status, Msg, err.Error())
	}

	//判斷賬戶餘額
	if oAccount.Available < fAmount {
		Status = 4012
		Msg = "用戶余额不足"
		this.RenderJson(Status, Msg, d)
	}

	//是否追号
	if len(sIsTrace) < 1 || len(sTraceWinStop) < 1 || len(sTraceStopValue) < 1 {
		Status = 4013
		Msg = "是否追号必需"
		this.RenderJson(Status, Msg, err.Error())
	}

	//是否追號
	iIsTrace, _ := strconv.Atoi(sIsTrace)
	if iIsTrace == 1 {
		oBet.IsTrace = true
	} else {
		oBet.IsTrace = false
	}

	//是否追號立停
	if sTraceWinStop == "1" {
		oBet.WinTrueStop = true
	} else {
		oBet.WinTrueStop = false
	}

	//奖金组
	iUserPrizeGroupId, sGroupName := models.UserPrizeSets.GetGroupId(sUserId, sGameId)
	iGroupName, _ := strconv.Atoi(sGroupName)
	if iUserPrizeGroupId < 1 {
		Status = 4014
		Msg = "奖金设置丟失"
		this.RenderJson(Status, Msg, iUserPrizeGroupId)
	}

	//最大獎金組＆最小獎金組
	iMaxBetGroup := models.UserPrizeSets.GetMaxBetGroup(uint(iGroupName), oBet.Series, oBet.Lottery)
	sMinBetGroup := models.SysConfig.GetPrizeByItem("player_min_grize_group")

	if len(sMinBetGroup) < 1 {
		Status = 4015
		Msg = "奖金设置有誤"
		this.RenderJson(Status, Msg, iUserPrizeGroupId)
	}
	iMinBetGroup, _ := strconv.Atoi(sMinBetGroup)
	oBet.minBetGroup = iMinBetGroup

	if iMaxBetGroup >= uint(iGroupName) {
		oBet.maxBetGroup = iGroupName
	} else {
		oBet.maxBetGroup = int(iMaxBetGroup)
	}

	//當前投注時間
	oBet.betTime = time.Now().Unix()
	for sIssue, _ := range oBet.BetIssue {
		oIssue := models.Issues.GetIssue(strconv.Itoa(oLottery.Id), sIssue)
		if oIssue.Id < 1 {
			Status = 4016
			Msg = "獎期錯誤"
			this.RenderJson(Status, Msg, oBet)
		}

		iEndTime, _ := strconv.Atoi(oIssue.EndTime)
		if oBet.betTime >= int64(iEndTime) {
			Status = 4016
			Msg = "超時投注"
			this.RenderJson(Status, Msg, oBet)
		}
	}

	if !models.Issues.CheckBetIssue(oLottery, oBet.betTime) {
		Status = 4016
		Msg = "奖期过期"
		this.RenderJson(Status, Msg, oBet)
	}

	//判断投注最大奖金组
	sPrizeLimit := models.SysConfig.GetPrizeByItem("bet_max_prize")
	iPrizeLimit, _ := strconv.ParseInt(sPrizeLimit, 0, 64)
	oBet.prizeLimit = iPrizeLimit

	//////////////
	//生成游戏记录//
	//////////////
	iTerminalId, _ := strconv.Atoi(webLoginMap["terminal"])
	iMechantId, _ := strconv.Atoi(webLoginMap["merchant_id"])
	oBet.betRecordId = models.BetRecords.CreateRecord(oUser, len(aBalls), iIsTrace, oLottery.Id, iTerminalId, iMechantId, aBalls, "")
	if oBet.betRecordId < 1 {
		Status = 4017
		Msg = "生成游戏记录失敗"
		this.RenderJson(Status, Msg, d)
	}

	///////////
	//生成投注//
	//////////
	iCode, errMsg, aBetResults := this.doBet(oBet)
	sUrl := NetIp
	var bSuccess bool
	if iCode == 200 {

		ziZhuDataArr := []map[string]string{}
		for _, ziZhuInfoRow := range ziZhuInfoArr {

			for _, ziZhuOrdersRow := range ziZhuOrdersArr {

				ziZhuData := ziZhuInfoRow
				ziZhuData["order_issue"] = ziZhuOrdersRow["issue"]       //奖期
				ziZhuData["order_multiple"] = ziZhuOrdersRow["multiple"] //倍数,追号的时候取这个倍数

				//多个玩法同时投注 或 单个玩法投注 的公有数据
				ziZhuData["public_seriesId"] = common.InterfaceToString(oLottery.SeriesId) //彩种系列
				ziZhuData["public_lotteryId"] = sGameId                                    //彩种id
				ziZhuData["public_traceWinStop"] = sTraceWinStop                           //是否中奖后停止追号1-是 0-否

				ziZhuData["test_betMoney"] = fmt.Sprintf("1个奖期投注金额 = 基本金(2) * 注数(%s) * 投注模式(%s) * 选号倍数(%s) * 追号倍数(%s) ",
					ziZhuData["balls_num"], ziZhuData["balls_coefficient"], ziZhuData["balls_multiple"], ziZhuData["order_multiple"])
				ziZhuDataArr = append(ziZhuDataArr, ziZhuData)
			}
		}

		if debug {
			fmt.Println("===ziZhuDataArr===", ziZhuDataArr)
		}

		zz := new(zizhu.ZizhuBase)
		go zz.SplitAlgorithm(&ziZhuDataArr) //自主开奖 拆单计奖

		bSuccess = true
	} else {
		bSuccess = false
	}

	result := map[string]interface{}{
		"isSuccess":     bSuccess,
		"type":          iCode,
		"bet_record_id": oBet.betRecordId,
		"data": map[string]interface{}{
			"tplData": aBetResults,
			"msg":     errMsg,
			"link":    sUrl + "/v1/game/bet.do",
		},
	}

	//回答日志
	b, _ := json.Marshal(result)
	log.AddLog("", "response", "user_id:"+sUserId+"==>data:"+string(b), log.BET)

	///////////////
	//输出json数据//
	//////////////
	this.RenderJson(iCode, errMsg, result)
}

/*
 * 投注數據整合程序
 */
func (this *GameController) doBet(oBet *Bet) (int, string, map[int]map[string]interface{}) {

	oAccount := oBet.Account
	aSeriesWays := map[int]*dao.SeriesWays{}

	webLoginMap := this.WebLogin()

	///////////
	//整合數據//
	//////////
	iCompileResult, aBetNumbers, aSeriesWays := this.compileBetData(oBet, aSeriesWays)
	if debug {
		fmt.Println("xxxxxxxxx整合後返回的數據xxxxxxxxxxx")
		fmt.Println("iCompileResult-->", iCompileResult)
		fmt.Println("aBetNumbers-->", aBetNumbers)
		fmt.Println("aSeriesWays-->", aSeriesWays)
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
	switch iCompileResult {
	case -1:
		return 4018, "奖金组错误", nil
	case -2:
		return 4019, "投注号码错误", nil
	case -3:
		return 4020, "投注金额错误", nil
	}

	/////////////
	//檢查賬戶鎖//
	////////////
	bIsLockAccount := models.Accounts.IsLockAccount(oAccount)
	if bIsLockAccount {
		return iCompileResult, "用戶账户被锁", nil
	}

	///////////////////
	///建立追號或者注單//
	//////////////////
	iSingal, aTrace, aProjects, aSeriesWays := this.CompileTaskAndProjects(oBet, aBetNumbers, aSeriesWays)
	if iSingal != 200 {
		switch iSingal {
		case 1001:
			return 4021, "投注价格错误", nil
		case 1002:
			return 4022, "投注注数错误", nil
		case 1003:
			return 4023, "投注总额过低(必须大于2分钱)", nil
		}
		return iCompileResult, "投注数据错误", nil
	}

	if debug {
		fmt.Println("XXXXXXXXXX當前解析的模型數據XXXXXXXXXX")
		fmt.Println("aTrace-->", aTrace)
		fmt.Println("aProjects-->", aProjects)
		fmt.Println("aSeriesWays-->", aSeriesWays)
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	}

	var aBetResults map[int]map[string]interface{}

	//組裝需要用到的額外數據
	mExtraData := map[string]string{
		"clientIP":      oBet.clientIP,
		"proxyIP":       oBet.proxyIP,
		"merchant_id":   webLoginMap["merchant_id"],
		"bet_record_id": fmt.Sprintf("%d", oBet.betRecordId),
		"terminal_id":   webLoginMap["terminal"],
	}

	if oBet.IsTrace {

		//追號任務
		aBetResults = this.CreateTrace(oBet, aSeriesWays, aTrace, mExtraData)
	} else {

		//注單任務
		aBetResults = this.CreateProjects(oBet, aProjects, aSeriesWays, mExtraData)
	}

	return iCompileResult, "success", aBetResults
}

/**
 * 格式化投注数据
 * @param array aBetNumbers 投注所有數據模型集
 * @param array aSeriesWays 系列投注方式們
 */
func (this *GameController) compileBetData(oBet *Bet, aSeriesWays map[int]*dao.SeriesWays) (int, []map[string]string, map[int]*dao.SeriesWays) {

	oSeries := oBet.Series
	oLottery := oBet.Lottery
	aBalls := oBet.Balls

	var iSeriesId int
	var oPrizeGroup *dao.PrizeGroups
	var aBetNumbers = []map[string]string{}
	var aPrizeGroups = map[string]*dao.PrizeGroups{}
	if oSeries.LinkTo > 0 {
		iSeriesId = int(oSeries.LinkTo)
	} else {
		iSeriesId = int(oLottery.SeriesId)
	}

	for _, oBall := range aBalls {
		sPrizeGroup := oBall.PrizeGroup
		iPrizeGroup, _ := strconv.Atoi(sPrizeGroup)

		//判斷獎金組範圍
		if debug {
			fmt.Println("sPrizeGroup-->", oBall.PrizeGroup)
			fmt.Println("maxBetGroup-->", oBet.maxBetGroup)
			fmt.Println("minBetGroup-->", oBet.minBetGroup)
		}

		if iPrizeGroup > oBet.maxBetGroup || iPrizeGroup < oBet.minBetGroup {
			return -1, nil, aSeriesWays
		}

		oPrizeGroup = models.PrizeGroups.GetPrizeGroupByClassicPrize(sPrizeGroup, fmt.Sprintf("%d", iSeriesId))
		if oPrizeGroup.Id == 0 {
			return -1, nil, aSeriesWays
		}
		aPrizeGroups[oPrizeGroup.Name] = oPrizeGroup

		iWayId, _ := strconv.Atoi(oBall.WayId)
		var oSeriesWay *dao.SeriesWays
		if obj, ok := aSeriesWays[iWayId]; ok {
			oSeriesWay = obj
		} else {
			oSeriesWay, _ = dao.GetSeriesWaysById(iWayId)
			if oSeriesWay.Id == 0 {
				return -1, nil, aSeriesWays
			}
			aSeriesWays[iWayId] = oSeriesWay
		}

		//整合投注數據
		mData := this._CompileBetData(oBall, sPrizeGroup, fmt.Sprintf("%d", oPrizeGroup.Id))
		if debug {
			fmt.Println("0000000本次整合的投注數據00000000")
			fmt.Println("mData-->", mData)
			fmt.Println("0000000000000000000000000000000")
		}

		//經過算法包計算投注號碼和投注注數是否合理
		iSingleCount, _, mData := this.CountBetNumber(oBet, oSeriesWay, mData)
		if debug {
			fmt.Println("11111111111111通過算法包整合的投注數據1111111111111111")
			fmt.Println("iSingleCount-->", iSingleCount)
			fmt.Println("mData-->", mData)
			fmt.Println("111111111111111111111111111111111111111111111111111")
		}

		if fmt.Sprintf("%d", iSingleCount) != mData["single_count"] { //注數錯誤
			return -3, nil, aSeriesWays
		}

		sBetNumber := mData["bet_number"]
		if sBetNumber == "" { //投注碼錯誤
			return -2, nil, aSeriesWays
		}

		//如果改系列玩法規定投注碼數格式需要拆分,找到算法包對應函數進行拆解
		//		if oSeriesWay.NeedSplit > 0 {
		//			var mWayMaps map[string]int
		//			json.Unmarshal([]byte(oSeriesWay.WayMaps), &mWayMaps)
		//			aPositions := common.MapKeys(mWayMaps, nil)
		//			aMethods := []string{"RandomEnum", "RondamMixCombin", "RandomConstituted", "RandomSpan", "RandomSum"}
		//			fMultiple, _ := strconv.ParseFloat(mData["multiple"], 64)
		//			fCoefficient, _ := strconv.ParseFloat(mData["coefficient"], 64)
		//			fPrice, _ := strconv.ParseFloat(mData["price"], 64)

		//			if oSeriesWay.WayFunction == "RandomSeparatedConstituted" {
		//				aBetNumberPos := strings.Split(sBetNumber, "|")

		//				//對於每個位置
		//				for _, sPosition := range aPositions {
		//					aPos := strings.Split(sPosition, "")
		//					var aBetNumberOfPosition []string
		//					iSplitWayId := mWayMaps[sPosition]
		//					for _, sPos := range aPos {
		//						iPos, _ := strconv.Atoi(sPos)
		//						if aBetNumberPos[iPos] == "" {
		//							continue
		//						}
		//						aBetNumberOfPosition = append(aBetNumberOfPosition, aBetNumberPos[iPos])
		//					}
		//					sSplitNumber := strings.Join(aBetNumberOfPosition, "|")
		//					mBetNumbers := this._CompileSplitBetData(fmt.Sprintf("%d", iSplitWayId), sSplitNumber, sPosition, aSeriesWays, fMultiple, fCoefficient, fPrice, mData)
		//					aBetNumbers = append(aBetNumbers, mBetNumbers)
		//				}
		//			} else if common.InArray(oSeriesWay.WayFunction, aMethods) {
		//				aChoosedPositions := strings.Split(mData["position"], "")
		//				for _, sPosition := range aPositions {
		//					aPos := strings.Split(sPosition, "")
		//					iSplitWayId := mWayMaps[sPosition]
		//					for _, sPos := range aPos {
		//						if !common.InArray(sPos, aChoosedPositions) {
		//							continue
		//						}
		//					}
		//					mBetNumbers := this._CompileSplitBetData(fmt.Sprintf("%d", iSplitWayId), sBetNumber, sPosition, aSeriesWays, fMultiple, fCoefficient, fPrice, mData)
		//					aBetNumbers = append(aBetNumbers, mBetNumbers)
		//				}
		//			}
		//		} else {
		aBetNumbers = append(aBetNumbers, mData)
		//		}
	}

	if len(aBetNumbers) < 1 {
		return -2, nil, aSeriesWays
	}

	return 200, aBetNumbers, aSeriesWays
}

/*
 * 組合投注數據
 * @param *models.BetData   oBall    投注數據單個模型
 * @param string sPrizeGroup         獎金組
 * @param string sGroupId            獎金組id
 */
func (this *GameController) _CompileBetData(oBall *models.BetData, sPrizeGroup, sGroupId string) map[string]string {
	fMoneyUnit, _ := strconv.ParseFloat(oBall.Moneyunit, 64)
	fNum, _ := strconv.ParseFloat(oBall.Moneyunit, 64)
	fMultiple, _ := strconv.ParseFloat(oBall.Moneyunit, 64)
	aData := map[string]string{
		"way":            oBall.WayId,                         //投注方式id
		"bet_number":     oBall.Ball,                          //投注數據
		"coefficient":    fmt.Sprintf("%.4f", fMoneyUnit),     //投注模式
		"multiple":       oBall.Multiple,                      //投注倍數
		"single_count":   strconv.Itoa(oBall.Num),             //投注數量
		"amount":         fmt.Sprintf("%f", fMultiple*fNum*2), //投注金額
		"price":          "2",                                 //投注價格(同onePrice)
		"prize_group":    sPrizeGroup,                         //獎金組
		"prize_group_id": sGroupId,                            //獎金組id
	}
	if oBall.Extra != nil {
		aData["position"] = oBall.Extra.Position
	}

	return aData
}

/*
 * 組合拆分後的投注數據
 * @param string     sSplitWayId       拆分的系列玩法id
 * @param string     sSplitNumber      拆分數字
 * @param string     sPosition         拆分位置
 * @param map[int]　　*dao.SeriesWays   系列投注方式們
 * @param float64     fMultiple　　　 　倍數
 * @param float64     fCoefficient     投注金額模式
 * @param float64     fPrice           投注價格
 * @param map[string]string   mData    投注的數據
 */
func (this *GameController) _CompileSplitBetData(oBet *Bet, sSplitWayId, sSplitNumber, sPosition string, aSeriesWays map[int]*dao.SeriesWays, fMultiple, fCoefficient, fPrice float64, mData map[string]string) map[string]string {
	iSplitWayId, _ := strconv.Atoi(sSplitWayId)
	var oSeriesWay *dao.SeriesWays
	if iSplitWayId < len(aSeriesWays) {
		oSeriesWay = aSeriesWays[iSplitWayId]
	} else {
		oSeriesWay, _ = dao.GetSeriesWaysById(iSplitWayId)
	}
	aTmpData := map[string]string{
		"bet_number": sSplitNumber,
		"position":   sPosition,
	}
	iSingleCount, sDisplayNumber, _ := this.CountBetNumber(oBet, oSeriesWay, aTmpData)

	//金額=倍數*投注模式*價格*注數
	fAmount := fMultiple * fCoefficient * fPrice * float64(iSingleCount)

	aRealData := mData
	aRealData["way"] = sSplitWayId
	aRealData["bet_number"] = sSplitNumber
	aRealData["single_count"] = fmt.Sprintf("%d", iSingleCount)
	aRealData["amount"] = fmt.Sprintf("%f", fAmount)
	aRealData["display_bet_number"] = sDisplayNumber
	aRealData["position"] = sPosition
	return aRealData
}

/*
 * 根據玩法包計算投注數據
 * @param *dao.SeriesWays     oSeriesWay  //系列投注方式模型
 * @param map[string]string   mData       //投注數據
 */
func (this *GameController) CountBetNumber(oBet *Bet, oSeriesWay *dao.SeriesWays, mData map[string]string) (int, string, map[string]string) {

	//////////////////////////////////////
	//查詢檢查需要用的基礎玩法基礎投注方式模型//
	//////////////////////////////////////
	var sBasicMethodId string

	//查詢基礎投注方式
	oBasicWay := models.BasicWays.GetInfoById(int(oSeriesWay.BasicWayId))
	if oBasicWay.Id < 1 {
		oBasicWay, _ = dao.GetBasicWaysById(int(oSeriesWay.BasicWayId))
	}
	if oBasicWay.Id < 1 {
		return 0, "找不到投注方式", mData
	}

	aBasicMethodIds := strings.Split(oSeriesWay.BasicMethods, ",")
	if len(aBasicMethodIds) == 1 {
		sBasicMethodId = aBasicMethodIds[0]
	} else {
		sBasicMethodId = "0"
		iMaxDigitalCount := 0
		for _, sId := range aBasicMethodIds {
			oTmpBasicMethod := models.BasicMethods.GetInfo(sId)
			if int(oTmpBasicMethod.DigitalCount) > iMaxDigitalCount {
				iMaxDigitalCount = int(oTmpBasicMethod.DigitalCount)
				sBasicMethodId = sId
			}
		}
	}

	//查詢基礎玩法
	oBasicMethod := models.BasicMethods.GetInfo(sBasicMethodId)
	if oBasicMethod.Id < 1 {
		iBasicMethodId, _ := strconv.Atoi(sBasicMethodId)
		oBasicMethod, _ = dao.GetBasicMethodsById(iBasicMethodId)
		if oBasicMethod.Id < 1 {
			return 0, "找不到玩法", mData
		}
	}

	if len(mData) > 1 {
		sRealBetNumber := models.SeriesWays.CompileBetNumber(mData["bet_number"], fmt.Sprintf("%d", oBet.Series.Id), oBasicWay, oBasicMethod)
		mData["bet_number"] = sRealBetNumber
		if debug {
			fmt.Println("==============CompileBetNumber=============")
			fmt.Println("sRealBetNumber-->", sRealBetNumber)
		}
	}

	if debug {
		fmt.Println("==============before進入算法包=============")
		fmt.Println("mData-->", mData)
	}

	iCount, sDisplayNumber, mData := models.SeriesWays.Count(oBasicWay, oBasicMethod, mData)
	if debug {
		fmt.Println("==============after進入算法包=============")
		fmt.Println("iCount->", iCount)
		fmt.Println("sDisplayNumber->", sDisplayNumber)
		fmt.Println("mData->", mData)
	}
	return iCount, sDisplayNumber, mData
}

/**
 * 生成追号任务数组及注单数组
 *
 * @param array     aBetNumbers　投注的數據們
 * @param array     aSeriesWays  系列投注方式們
 */
func (this *GameController) CompileTaskAndProjects(oBet *Bet, aBetNumbers []map[string]string, aSeriesWays map[int]*dao.SeriesWays) (int, []map[string]map[string]string, []map[string]string, map[int]*dao.SeriesWays) {

	bTrace := oBet.IsTrace
	oLottery := oBet.Lottery
	mBetIssues := oBet.BetIssue //map[獎期]期數
	oUser := oBet.User

	//判斷是否即時彩
	if oLottery.IsInstant > 0 {
		bTrace = false
	}
	var fTotalValidAmount float64 = 0             //追號單產生的總金額
	var aTrace = []map[string]map[string]string{} //追號們
	var aProjects = []map[string]string{}         //注單們
	var mIssueEndTimes = map[string]string{}      //獎期結束時間
	var iSingal int = 200

	//對於每個投注數據
	for _, mBetNumber := range aBetNumbers {

		//獲取系列玩法
		sSeriesWay := mBetNumber["way"]
		iSeriesWay, _ := strconv.Atoi(sSeriesWay)
		var oSeriesWay *dao.SeriesWays
		if iSeriesWay < len(aSeriesWays) {
			oSeriesWay = aSeriesWays[iSeriesWay]
		} else {
			oSeriesWay, _ = dao.GetSeriesWaysById(iSeriesWay)
		}

		//獲取獎金配置,和最大獎金
		mPrizeSettingOfWay, fMaxPrize := this.MakePrizeSettingArray(mBetNumber["prize_group_id"], oSeriesWay)

		//獲取最大投注倍數
		fPrizeLimit, _ := strconv.ParseFloat(fmt.Sprintf("%d", oBet.prizeLimit), 64)
		var fMaxMultiple float64
		if fMaxPrize > 0 {
			//			fMaxMultiple = math.Floor(fPrizeLimit / fMaxPrize)
			//最大投注倍数=奖金/最小投注额
			fMaxMultiple = math.Floor(fPrizeLimit / 0.02)
		}

		//判斷系列價格和投注價格是否相同
		if fmt.Sprintf("%d", oSeriesWay.Price) != mBetNumber["price"] {
			iSingal = 1001
			break
		}

		//投注數據格式化
		fOriginalMultiple, _ := strconv.ParseFloat(mBetNumber["multiple"], 64) //投注數據的投注倍數
		fCoefficient, _ := strconv.ParseFloat(mBetNumber["coefficient"], 64)   //投注數據中的投注模式
		fSingleCount, _ := strconv.ParseFloat(mBetNumber["single_count"], 64)  //投注數據中的單次數量

		//單倍投注金額=數量*系列價格*投注模式（如2角）
		fSingleAmount := fSingleCount * float64(oSeriesWay.Price) * fCoefficient
		fSingleAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", fSingleAmount), 64)

		var fMaxBetMultiple float64
		if debug {
			fmt.Println("===============投注數據格式化================")
			fmt.Println("fMaxMultiple-->", fMaxMultiple)
			fmt.Println("fOriginalMultiple-->", fOriginalMultiple)
			fmt.Println("fCoefficient-->", fCoefficient)
			fmt.Println("fSingleCount-->", fSingleCount)
			fmt.Println("fSingleAmount-->", fSingleAmount)
		}

		//最大注數*倍數*投注模式
		fMaxBetMultiple = float64(common.MapMaxStr(mBetIssues)) * fOriginalMultiple * fCoefficient
		if debug {
			fmt.Println("fMaxBetMultiple-->", fMaxBetMultiple)
		}
		if fMaxBetMultiple > 0 && fMaxBetMultiple > fMaxMultiple {
			iSingal = 1002
			break
		}

		//計算總金額
		var fValidBaseAmount float64

		//追号
		if oBet.IsTrace {
			for _, orderMutiple := range oBet.BetIssue {
				iOrderMutiple, _ := strconv.Atoi(orderMutiple)
				fValidBaseAmount += fSingleAmount * fOriginalMultiple * float64(iOrderMutiple)
			}
		} else {
			fValidBaseAmount = fSingleAmount * fOriginalMultiple
		}

		if debug {
			fmt.Println("fValidBaseAmount-->總金額", fValidBaseAmount)
			fmt.Println("===========================================")
		}
		if fValidBaseAmount < 0.02 {
			iSingal = 1003
			break
		}

		//判断风控
		bIsOverLimit := models.Projects.IsOverMoneyLimit(oUser.MerchantId, oLottery.Id, oSeriesWay.Id, fValidBaseAmount)
		if !bIsOverLimit {
			iSingal = 1004
			break
		}

		//是否追號
		if bTrace {

			//組合單個投注號碼追號
			mTraces, mEnd := this.AddTraceTaskQueue(oBet, fOriginalMultiple, mBetNumber, mPrizeSettingOfWay, mIssueEndTimes)
			mIssueEndTimes = mEnd
			aTrace = append(aTrace, mTraces)
		} else {

			//組合單個注單投注
			mProject, mEnd := this.AddSingleProject(oBet, fSingleCount, fOriginalMultiple, fSingleAmount, mPrizeSettingOfWay, mBetNumber, mIssueEndTimes)
			mIssueEndTimes = mEnd
			aProjects = append(aProjects, mProject)
		}

		fTotalOrderMultiple := common.MapSumF(mBetIssues)
		fTotalValidAmount += fTotalOrderMultiple * fValidBaseAmount
	}

	return iSingal, aTrace, aProjects, aSeriesWays
}

/**
 * 生成奖金设置数组，供投注功能使用
 *
 * @param string iPrizeGroupId       獎金組id
 * @param *dao.SeriesWays oSeriesWay 系列玩法
 */
func (this *GameController) MakePrizeSettingArray(sPrizeGroupId string, oSeriesWay *dao.SeriesWays) (map[string]map[int]float64, float64) {

	var mPrizeSettingOfMethods = map[string]map[int]float64{}
	var fMaxPrize float64 = 0

	//拆分系列玩法的基礎投注方式id
	aMethodIds := strings.Split(oSeriesWay.BasicMethods, ",")
	for _, sMethodId := range aMethodIds {

		//獲取對應玩法和獎金組的獎金詳情
		mPrize := models.PrizeDetails.GetPrizeSetting(sPrizeGroupId, sMethodId)
		mPrizeSettingOfMethods[sMethodId] = mPrize //map[玩法id] -> map[獎金等級]獎金金額

		//令最大獎金爲獎金等級１級的獎金
		if fMaxPrize < mPrize[1] {
			fMaxPrize = mPrize[1]
		}
	}

	return mPrizeSettingOfMethods, fMaxPrize
}

/**
 * 向追号任务数组中增加一个任务
 *
 * @param fOriginalMultiple float64 投注數據中的倍數
 * @param mBetNumber map[string]string 投注數據
 * @param aPrizeSettingOfWay interface{} 玩法獎金設置
 * @param aIssueEndTimes map[string]string 投注獎期結束時間
 */
func (this *GameController) AddTraceTaskQueue(oBet *Bet, fOriginalMultiple float64, mBetNumber map[string]string, aPrizeSettingOfWay interface{}, aIssueEndTimes map[string]string) (map[string]map[string]string, map[string]string) {
	mIssues := oBet.BetIssue
	oLottery := oBet.Lottery

	mEndTimes := map[string]string{}    //結束時間-->map[獎期號]獎期結束時間
	aTraceIssues := map[string]string{} //追號獎期-->map[獎期號]總倍數

	for sIssue, sOrderCount := range mIssues {
		iOrderCount, _ := strconv.Atoi(sOrderCount)           //投注的倍数
		fMultiple := fOriginalMultiple * float64(iOrderCount) //總倍數=注數*倍數
		aTraceIssues[sIssue] = fmt.Sprintf("%f", fMultiple)
		if _, ok := aIssueEndTimes[sIssue]; !ok {
			oIssue := models.Issues.GetIssue(fmt.Sprintf("%d", oLottery.Id), sIssue)
			aIssueEndTimes[sIssue] = oIssue.EndTime
		}
		mEndTimes[sIssue] = aIssueEndTimes[sIssue]
	}
	jsData, _ := json.Marshal(aPrizeSettingOfWay)
	mBetNumber["prize_set"] = string(jsData)
	mTraces := map[string]map[string]string{
		"bet":       mBetNumber,
		"issues":    aTraceIssues,
		"end_times": mEndTimes,
	}
	return mTraces, aIssueEndTimes
}

/**
 * 向注单数组中增加一个注单
 * @param fSingleCount float64    投注注數
 * @param iMultiple    int    投注倍數
 * @param fSingleAmount float64  单倍金额
 * @param aPrizeSettingOfWay interface{}  獎金設置
 * @param mBetNumber map[string]string  投注數據
 * @param aIssueEndTimes map[string]string  獎期結束時間
 */
func (this *GameController) AddSingleProject(oBet *Bet, fSingleCount, fMultiple, fSingleAmount float64, aPrizeSettingOfWay interface{}, mBetNumber, aIssueEndTimes map[string]string) (map[string]string, map[string]string) {
	oLottery := oBet.Lottery
	mBetIssues := oBet.BetIssue
	var mOrderInfo map[string]string

	for sIssue, _ := range mBetIssues {
		oIssue := models.Issues.GetIssue(fmt.Sprintf("%d", oLottery.Id), sIssue)
		if _, ok := aIssueEndTimes[sIssue]; !ok {
			aIssueEndTimes[sIssue] = oIssue.EndTime
		}

		prizeSet, _ := json.Marshal(aPrizeSettingOfWay)
		mOrderInfo = map[string]string{
			"issue":         sIssue,
			"end_time":      oIssue.EndTime,
			"single_count":  fmt.Sprintf("%f", fSingleCount),
			"multiple":      fmt.Sprintf("%f", float64(common.MapSum(mBetIssues))*fMultiple),
			"single_amount": fmt.Sprintf("%f", fSingleAmount),
			"prize_set":     string(prizeSet),
			"prize_group":   mBetNumber["prize_group"],
		}
		for sKey, sValue := range mBetNumber {
			mOrderInfo[sKey] = sValue
		}
		break
	}
	return mOrderInfo, aIssueEndTimes
}

/**
 * 追号任务入库
 * @param aSeriesWays    int                             系列玩法
 * @param aTraces        []map[string]map[string]string  組裝後的追號數據
 * @param mExtraData     map[string]string               公共數據
 */
func (this *GameController) CreateTrace(oBet *Bet, aSeriesWays map[int]*dao.SeriesWays, aTraces []map[string]map[string]string, mExtraData map[string]string) map[int]map[string]interface{} {
	oUser := oBet.User
	oAccount := oBet.Account
	oLottery := oBet.Lottery
	iBetTime := oBet.betTime
	bStopOnPrized := oBet.WinTrueStop

	o := orm.NewOrm()
	sqlErr := o.Begin()
	if sqlErr != nil {
		log.AddLog("", "sqlErr", fmt.Sprintf("user_id:%d;lottery_id:%d==>err:%s", oUser.Id, oLottery.Id, sqlErr.Error()), log.SQL)
	}

	defer func() {
		if err := recover(); err != nil {
			o.Rollback()
		}

	}()

	var aBetResults = map[int]map[string]interface{}{}
	for _, mTrace := range aTraces {
		iWay, _ := strconv.Atoi(mTrace["bet"]["way"])
		bReturn, sErrMsg := models.Traces.CreateTrace(o, oUser, oAccount, aSeriesWays[iWay], oLottery, bStopOnPrized, mExtraData, fmt.Sprintf("%d", iBetTime), mTrace)
		if !bReturn {
			sRes := map[string]interface{}{
				"way":    mTrace["bet"]["way"],
				"ball":   mTrace["bet"]["bet_number"],
				"reason": bReturn,
				"error":  sErrMsg,
			}
			aBetResults[0] = sRes
			sqlErr = o.Rollback()
			if sqlErr != nil {
				log.AddLog("", "sqlErr", fmt.Sprintf("user_id:%d;lottery_id:%d==>err:%s", oUser.Id, oLottery.Id, sqlErr.Error()), log.SQL)
			}
			break
		} else {
			sRes := map[string]interface{}{
				"way":  mTrace["bet"]["way"],
				"ball": mTrace["bet"]["bet_number"],
			}
			aBetResults[1] = sRes
		}
	}
	sqlErr = o.Commit()
	if sqlErr != nil {
		log.AddLog("", "sqlErr", fmt.Sprintf("user_id:%d;lottery_id:%d==>err:%s", oUser.Id, oLottery.Id, sqlErr.Error()), log.SQL)
	}
	return aBetResults
}

/**
 * 注单入库
 * @param aSeriesWays    int                  系列玩法
 * @param aProjects      []map[string]string  組裝好的注單數組
 * @param mExtraData     map[string]string    公共數據
 */
func (this *GameController) CreateProjects(oBet *Bet, aProjects []map[string]string, aSeriesWays map[int]*dao.SeriesWays, mExtraData map[string]string) map[int]map[string]interface{} {
	oUser := oBet.User
	oLottery := oBet.Lottery

	o := orm.NewOrm()
	sqlErr := o.Begin()
	if sqlErr != nil {
		log.AddLog("", "sqlErr", fmt.Sprintf("user_id:%d;lottery_id:%d==>err:%s", oUser.Id, oLottery.Id, sqlErr.Error()), log.SQL)
	}
	defer func() {
		if err := recover(); err != nil {
			o.Rollback()
		}

	}()

	if len(aProjects) < 1 {
		return nil
	}

	var aBetResults = map[int]map[string]interface{}{}
	for _, mProject := range aProjects {

		iWay, _ := strconv.Atoi(mProject["way"]) //系列投注方式id
		oProject := models.Projects.CompileProjectData(oUser, aSeriesWays[iWay], oLottery, fmt.Sprintf("%d", oBet.betTime), mExtraData, mProject)
		bSuccess, _, msg, err := models.Projects.AddProject(o, oUser, oBet.Account, aSeriesWays[iWay], oProject, mExtraData)
		if !bSuccess {
			sRes := map[string]interface{}{
				"way":    mProject["way"],
				"ball":   mProject["bet_number"],
				"reason": msg,
				"error":  err,
			}
			aBetResults[0] = sRes
			sqlErr = o.Rollback()
			if sqlErr != nil {
				log.AddLog("", "sqlErr", fmt.Sprintf("user_id:%d;lottery_id:%d==>err:%s", oUser.Id, oLottery.Id, sqlErr.Error()), log.SQL)
			}
			break
		} else {
			sRes := map[string]interface{}{
				"way":  mProject["way"],
				"ball": mProject["bet_number"],
			}
			aBetResults[1] = sRes
		}
	}
	sqlErr = o.Commit()
	if sqlErr != nil {
		log.AddLog("", "sqlErr", fmt.Sprintf("user_id:%d;lottery_id:%d==>err:%s", oUser.Id, oLottery.Id, sqlErr.Error()), log.SQL)
	}
	return aBetResults
}

/*
 * 獲取正在銷售的獎期生成投注線程id
 */
func (this *GameController) GetOnSaleAndAddThread(oBet *Bet) bool {
	sLotteryId := fmt.Sprintf("%d", oBet.Lottery.Id)
	sResult := Lmodels.Issues.GetOnSaleIssue(sLotteryId)
	if len(sResult) < 1 {
		return false
	}
	sDbThreadId := Lmodels.User.GetDbThreadId()
	Lmodels.BetThread.AddThread(sLotteryId, sResult, sDbThreadId)
	return true
}

/*
 * 獲取正在銷售的獎期
 */
func (this *GameController) GetOnSaleIssue(oBet *Bet) bool {
	sOnSaleIssue, _ := models.Issues.GetOnSaleIssue(oBet.Lottery.Id)
	if len(sOnSaleIssue) > 1 {
		return true
	}
	return false
}

// @Title trend
// @Description 开奖走势
// @Param	params		query 	controllers.ParamsInputType		true		"lottery_id=彩种ＩＤ,num_type=玩法：五星四星前三..."
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /tides.do [post]
func (c *GameController) TrendViewData() {

	var reStatus int = 500
	var reMsg string = "no"
	d := map[string]interface{}{}

	webLoginMap := c.WebLogin()

	//merchantId := webLoginMap["merchant_id"]
	//userId := webLoginMap["user_id"]

	lotteryId := webLoginMap["lottery_id"] //彩种ＩＤ
	numType := webLoginMap["num_type"]     //玩法：五星四星前三
	issuesNum := webLoginMap["issues_num"] //期数
	beginTime := webLoginMap["begin_time"] //开始时间
	endTime := webLoginMap["end_time"]     //结束时间

	if len(lotteryId) == 0 {
		c.RenderJson(501, "彩种id不能为空", d) //输出json数据
	}

	if len(numType) == 0 {
		c.RenderJson(502, "玩法id不能为空", d) //输出json数据
	}

	if len(issuesNum) == 0 { //设置默认值
		issuesNum = "30"
	}

	oLottery := Lmodels.Lottery.GetInfo(lotteryId) //获得彩种信息
	if len(oLottery) == 0 {
		c.RenderJson(503, "彩种id有误", d) //输出json数据
	}

	oSeriesIdentifier := ""
	lotSerId := oLottery["series_id"]
	if (lotSerId == "1") || lotSerId == "2" || lotSerId == "3" {

		oSeriesIdentifier = "UserTrend"
	} else {

		oSeries := models.Series.RGetOneById(lotSerId)
		oSeriesIdentifier = strings.ToLower(oSeries["identifier"])
	}

	result := [][][][]string{}
	switch oSeriesIdentifier {

	case "UserTrend":
		obj := new(lib.UserTrend)

		/*
			isSuccess, data, statistics, omissionBarStatus := obj.GetTrendDataByParams(lotteryId, numType, beginTime, endTime, count)

			result["isSuccess"] = isSuccess
			result["data"] = data
			result["statistics"] = statistics
			result["omissionBarStatus"] = omissionBarStatus
		*/
		reStatus, reMsg, result = obj.GetTrendDataByParams(lotteryId, numType, beginTime, endTime, issuesNum)
		break

	}

	d["data"] = result

	c.RenderJson(reStatus, reMsg, d) //输出json数据
}

// @Title project_detailt
// @Description 获取某个注单详情
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /project_detailt.do [post]
func (c *GameController) GetProjectDetailt() {

	reStatus := 500
	reMsg := "id不能为空"
	d := map[string]interface{}{}

	webLoginMap := c.WebLogin()

	//	merchantId := webLoginMap["merchant_id"]
	userId := webLoginMap["user_id"]

	projects_id := webLoginMap["id"] //注单id
	if len(projects_id) == 0 {
		c.RenderJson(reStatus, reMsg, d) //输出json数据
	}

	oProject := Lmodels.Project.GetInfo(projects_id)
	if len(oProject) == 0 {
		reStatus, reMsg = 501, "id错误"
		c.RenderJson(reStatus, reMsg, d) //输出json数据
	}

	if oProject["user_id"] != userId {
		reStatus, reMsg = 502, "无权查看订单"
		c.RenderJson(reStatus, reMsg, d) //输出json数据
	}

	traceId := "0"
	oTrace := map[string]string{}
	isTraced := 0
	if len(oProject["trace_id"]) > 0 && oProject["trace_id"] != "0" { //是否有追号任务ID
		traceId = oProject["trace_id"]
		oTrace = Lmodels.Trace.GetInfo(traceId)
		isTraced = 1
	}

	prize_set_formatted := Lmodels.Project.GetPrizeSetFormattedAttribute(oProject["prize_set"], oProject["coefficient"])

	oLottery := Lmodels.Lottery.GetInfo(oProject["lottery_id"])
	sLottery := BaseControllerX.GetLang(oLottery["name"])

	iStopOnWin := "0"
	if stopOnWonV, stopOnWonOk := oTrace["stop_on_won"]; stopOnWonOk == true {
		iStopOnWin = stopOnWonV
	}

	//	amountF64, _ := common.Str2Float64(oProject["amount"])
	coefficientF64, _ := common.Str2Float64(oProject["coefficient"])
	/*
		var betCount float64 = 0.00
		if coefficientF64 > 0 {
			betCount = amountF64 / coefficientF64 / 2
		}
	*/

	coefficient := ""
	if _, coefficientOk := oProject["coefficient"]; coefficientOk == true {
		coefficient = models.GetCoefficientText(coefficientF64)
	}

	splittedWinningNumber := Lmodels.Project.GetSplittedWinningNumberAttribute(oProject["winning_number"], oLottery)
	amountFormatted := Lmodels.Project.GetAmountFormattedAttribute(oProject["amount"])

	winningNumber := "--"
	if winningNumberV, winningNumberOK := oProject["winning_number"]; winningNumberOK == true {
		winningNumber = winningNumberV
	}

	prize := "--"
	if prizeV, prizeOk := oProject["prize"]; prizeOk == true {
		prize = prizeV
	}

	formattedStatus := Lmodels.Project.GetFormattedStatusAttribute(oProject["status"])
	formattedStatusFY := BaseControllerX.GetLang(formattedStatus)

	canBeCanceled := 0
	if oProject["status"] == Ldefined.PROJECT_STATUS_NORMAL {
		canBeCanceled = 1
	}

	data := map[string]interface{}{
		"id":                  oProject["id"],
		"user_id":             oProject["user_id"],
		"username":            oProject["username"],
		"terminal_id":         oProject["terminal_id"],        //终端id
		"serial_number":       oProject["serial_number"],      //注单编号
		"trace_id":            traceId,                        //追号任务ID
		"prize_group":         oProject["prize_group"],        //投注时的奖金组
		"prize_set_formatted": prize_set_formatted,            //二星组选 : 一等奖: 4.875元
		"lottery_id":          oProject["lottery_id"],         //彩种id
		"lottery_identifier":  oLottery["identifier"],         //彩种标识符
		"lottery":             sLottery,                       //彩种名翻译
		"issue":               oProject["issue"],              //奖期
		"way_id":              oProject["way_id"],             //投注方式id（Series Way Id）
		"way":                 oProject["title"],              //标题
		"bet_number":          oProject["display_bet_number"], //显示出来的投注号码
		//		"bet_count":               fmt.Sprintf("%.2f", betCount),  //amount / coefficient / 2
		"bet_count":               oProject["single_count"],
		"multiple":                oProject["multiple"],    //倍数
		"coefficient":             oProject["coefficient"], //投注金额模式(1=2元;0.5=1元;0.1=2角0.05=1角;0.01=2分;0.001=2厘)
		"formatted_coefficient":   coefficient,             // 2元;2角;2分;2厘...
		"splitted_winning_number": splittedWinningNumber,   //开奖号码切割成单个数字
		"amount":                  oProject["amount"],      //总金额
		"amount_formatted":        amountFormatted,         //总金额处理50000 -> 5,000.0000
		"winning_number":          winningNumber,           //开奖号码
		"prize":                   prize,                   //奖金
		"status":                  oProject["status"],      //0: 正常；1：已撤销；2：未中奖；3：已中奖；4：已派奖；5：系统撤销（通过redis和帐变表去重复）
		"formatted_status":        formattedStatusFY,       //状态码翻译
		"can_be_canceled":         canBeCanceled,           //状态正常返回1,否则0
		"bought_at":               oProject["bought_at"],
		"is_traced":               isTraced,   //有追号任务ID返回1,否则0
		"stop_on_win":             iStopOnWin, //中奖即停
	}

	reStatus = 200
	reMsg = "yes"
	d = data

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}

// @Title account_change_list
// @Description 账变记录
// @Param	params	query 	controllers.ParamsInputType	true		"created_at_from=开始时间&created_at_to=结束时间&type_id=账变类别ID&page=1&rows=20"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /account_change_list.do [post]
func (c *GameController) GetAccountChangeList() {

	reStatus := 500
	reMsg := "id不能为空"
	d := map[string]interface{}{}

	webLoginMap := c.WebLogin()

	merchantId := webLoginMap["merchant_id"] //登录获取 商户id
	userId := webLoginMap["user_id"]         //登录获取 用户id

	startAt := webLoginMap["created_at_from"] //账变生成开始时间
	endAt := webLoginMap["created_at_to"]     //账变生成结束时间
	typeId := webLoginMap["type_id"]          //账变类别ID

	page := webLoginMap["page"] //第几页
	rows := webLoginMap["rows"] //每页取多少行

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

	if len(typeId) > 0 {
		sWhere += fmt.Sprintf(" and type_id='%s' ", typeId)
	}

	sqlPage := 1 //默认第一页
	if len(page) > 0 {
		pageInt, pageErr := strconv.Atoi(page)
		if pageErr == nil {
			sqlPage = pageInt
		}
	}

	sqlRow := 20 //默认每页取20行
	if len(rows) > 0 {
		rowsInt, rowsErr := strconv.Atoi(rows)
		if rowsErr == nil {
			sqlRow = rowsInt
		}
	}

	iOffset := (sqlPage - 1) * sqlRow
	oTransactions := Lmodels.Transaction.GetAllLists(sWhere, iOffset, sqlRow) //转账记录
	count := Lmodels.Transaction.GetCount(sWhere)

	reData := []map[string]interface{}{}
	if len(oTransactions) > 0 {
		aTransactionTypes := Lmodels.TransactionType.GetAllTransactionTypesArray() //所有转账类型
		aSeriesWays := Lmodels.SeriesWay.IdTitleList()                             //所有转账类型

		for _, oTransaction := range oTransactions {

			reMap := map[string]interface{}{}

			reMap["id"] = oTransaction["id"]
			reMap["username"] = oTransaction["username"]
			reMap["serial_number"] = oTransaction["serial_number"]
			reMap["type_id"] = oTransaction["type_id"]
			reMap["is_income"] = oTransaction["is_income"]
			reMap["trace_id"] = oTransaction["trace_id"]
			reMap["lottery_id"] = oTransaction["lottery_id"]
			reMap["issue"] = oTransaction["issue"]
			reMap["way_id"] = oTransaction["way_id"]
			reMap["coefficient"] = oTransaction["coefficient"]
			reMap["project_id"] = oTransaction["project_id"]
			reMap["amount"] = oTransaction["amount"]
			reMap["available"] = oTransaction["available"]
			reMap["note"] = oTransaction["note"]
			reMap["tag"] = oTransaction["tag"]
			reMap["balance"] = oTransaction["balance"]
			reMap["extra_data"] = oTransaction["extra_data"]

			reMap["ablance"] = oTransaction["available"]     //帐变前可用额度
			reMap["created_at"] = oTransaction["created_at"] //订单创建时间

			oLottery := Lmodels.Lottery.GetInfo(oTransaction["lottery_id"]) //彩种id

			lotName := "--"
			if lotNameV, lotNameOk := oLottery["name"]; lotNameOk == true {
				lotName = BaseControllerX.GetLang(lotNameV)
			}
			reMap["lottery"] = lotName //lotteries表 name字段 并且翻译

			seriesWaysTitle := "--"
			if seriesWaysTitleV, seriesWaysTitleOk := aSeriesWays[oTransaction["way_id"]]; seriesWaysTitleOk == true {
				seriesWaysTitle = seriesWaysTitleV
			}
			reMap["way"] = seriesWaysTitle //通过way_id 得到series_ways表 name

			reMap["formatted_amount"] = Lmodels.Transaction.GetAmountFormattedAttribute(oTransaction["is_income"], oTransaction["amount"])

			tranCnTitle := "--"
			if tranCnTitleV, tranCnTitleOk := aTransactionTypes[oTransaction["type_id"]]; tranCnTitleOk == true {
				tranCnTitle = tranCnTitleV
			}

			//			reMap["description"] = BaseControllerX.GetLang(oTransaction["description"]) //描述（来源于帐变类型）-并且翻译
			reMap["description"] = tranCnTitle
			reMap["type"] = tranCnTitle //通过type_id 得到 transaction_types表 cn_title 字段

			reMap["serial_number"] = Lmodels.GetSerialNumberShortAttribute(oTransaction["serial_number"]) //serial_number取后6位

			reMap["transfer_in"] = "--"
			reMap["transfer_out"] = "--"

			switch oTransaction["type_id"] {

			//1
			case Ldefined.TRAN_TYPE_PLAT_TRANSFER_IN:
				reMap["transfer_in"] = BaseControllerX.GetLang("platform")
				reMap["transfer_out"] = "转入"
				break

			//2
			case Ldefined.TRAN_TYPE_PLAT_TRANSFER_OUT:
				reMap["transfer_in"] = BaseControllerX.GetLang("platform")
				reMap["transfer_out"] = "转出"
				break

			//3
			case Ldefined.TRAN_TYPE_TRANSFER_IN:
				reMap["transfer_in"] = BaseControllerX.GetLang("myself")
				reMap["transfer_out"] = BaseControllerX.GetLang("parent")
				break

			//4
			case Ldefined.TRAN_TYPE_TRANSFER_OUT:
				reMap["transfer_in"] = BaseControllerX.GetLang("children")
				reMap["transfer_out"] = BaseControllerX.GetLang("myself")
				break

				//5
			case Ldefined.TRAN_TYPE_GET_SALARY:
				reMap["transfer_in"] = BaseControllerX.GetLang("myself")
				reMap["transfer_out"] = BaseControllerX.GetLang("parent")
				break

				//6
			case Ldefined.TRAN_TYPE_PAY_SALARY:
				reMap["transfer_in"] = BaseControllerX.GetLang("children")
				reMap["transfer_out"] = BaseControllerX.GetLang("myself")
				break

				//7
			case Ldefined.TRAN_TYPE_DEPOSIT:
				oExtraData := map[string]interface{}{}
				json.Unmarshal([]byte(oTransaction["extra_data"]), &oExtraData)
				reMap["platform"] = oExtraData
				break

				//8
			case Ldefined.TRAN_TYPE_DEPOSIT_BY_ADMIN:
				reMap["platform"] = BaseControllerX.GetLang("deposit-by-admin")
				break

			}

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

// @Title cancel_project
// @Description 用户撤单
// @Param	params		query 	controllers.ParamsInputType		true		"id=注单ID"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /cancel_project.do [post]
func (c *GameController) DropProject() {

	reStatus := 500
	reMsg := "id不能为空"
	d := map[string]interface{}{}

	webLoginMap := c.WebLogin()

	userId := webLoginMap["user_id"] //登录获取 用户id
	projectsId := webLoginMap["id"]  //注单ID

	if len(projectsId) == 0 {
		c.RenderJson(reStatus, reMsg, d) //输出json数据
	}

	reStatus, reMsg = Lmodels.Project.Drop(projectsId, userId, Ldefined.PROJECT_DROP_BY_USER)

	//输出json数据
	c.RenderJson(reStatus, reMsg, d)
}

// @Title CurrentUserInfo
// @Description 用戶及賬戶信息,author(leon)
// @Param	merchant_identity		query 	string	true		"接入方标识代码 如:JMG "
// @Param	params				query 	string	true		"token=xxx&terminal_id=xxx"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /current_user_info.do [post]
func (this *GameController) CurrentUserInfo() {

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

	sMerchantIdentity := webLoginMap["identity"]
	sessionUid := webLoginMap["user_id"]

	iUserId, _ := strconv.Atoi(sessionUid)
	oUser, err := dao.GetUsersById(iUserId)

	if oUser.Id < 1 {
		Status = 331
		Msg = "用户不存在"
		result = err
		this.RenderJson(Status, Msg, result)
	}

	oAccount := models.Accounts.GetAccountInfoByUserId(fmt.Sprintf("%d", iUserId))
	if oAccount.Id == 0 {
		Status = 332
		Msg = "用户账户错误"
		this.RenderJson(Status, Msg, result)
	}
	result = map[string]interface{}{
		"username":          oUser.Username,
		"nickname":          oUser.Nickname,
		"merchant_identity": sMerchantIdentity,
		"prize_group":       oUser.PrizeGroup,
		"blocked":           oUser.Blocked,
		"is_tester":         oUser.IsTester,
		"bet_coefficient":   oUser.BetCoefficient,
		"bet_multiple":      oUser.BetMultiple,
		"available":         oAccount.Available,
		"login_ip":          oUser.LoginIp,
		"signin_at":         oUser.SigninAt.Format(common.DATE_FORMAT_YMDHIS),
	}

	this.RenderJson(Status, Msg, result)
}
