package controllers

import (
	"NewLottApi/configs"
	"NewLottApi/dao"
	"NewLottApi/models"
	"common"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//PublicController 公共类
type PublicController struct {
	MainController
	Merchant *dao.Merchants
}

// @Title Token
// @Description 验证token是否正確,阻止非法登陸,author(leon）
// @Param	params	query	string	true		"token=1&ip=2,ip=接口调用方服务器IP"
// @Success	200	{JsonOut}	success!
// @Failure	500	error
// @router /token.do [post]
func (c *PublicController) Token() {

	Status := 200
	Msg := "ok"
	var result interface{}

	webLoginMap := c.WebLogin()
	if len(webLoginMap["ip"]) < 1 {
		//ip不能为空
	}

	if len(webLoginMap["token"]) < 1 {
		//token不能为空
	}

	if webLoginMap["client_ip"] != webLoginMap["ip"] {
		//ip不一致
	}

	c.RenderJson(Status, Msg, result)
}

// @Title LotteryData
// @Description 商户可以投注的彩种数据,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /lottery_data.do [post]
func (this *PublicController) LotteryData() {
	var result interface{}
	//	sMerchantId := webLoginMap["merchant_id"]
	//	if len(sMerchantId) < 1 {
	//		this.RenderJson(3008, "error", result)
	//	}
	//	aCloseIds := models.MerchantsLotteryClose.GetInfoById(sMerchantId)
	//	result = models.Lotteries.GetNeedLotteriesData(aCloseIds)
	this.RenderJson(200, "success", result)
}

// @Title Load
// @Description 根据传递过来的token和游戏标记打开对应的彩票选号盘,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /load-data.do [post]
func (this *PublicController) Load() {
	Status := 200
	Msg := "success"
	var result interface{}

	webLoginMap := this.WebLogin()

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			Status = 500
			Msg = "系统错误"
			this.RenderJson(Status, Msg, err)
		}

	}()

	///////////
	//判斷彩種//
	//////////
	sLotteryIden := webLoginMap["lottery"]                         //彩种标记
	oLottery := models.Lotteries.GetInfoByIdentifier(sLotteryIden) //获取彩种数据
	if oLottery.Id < 1 {
		Status = 3001
		Msg = "彩种不存在"
		this.RenderJson(Status, Msg, result)
	}

	///////////
	//判斷用戶//
	//////////
	sUserId := webLoginMap["user_id"]
	iUserId, _ := strconv.Atoi(sUserId)
	oUser, _ := dao.GetUsersById(iUserId) //获取用户数据
	if oUser.Id == 0 {
		Status = 3002
		Msg = "用户信息错误"
		this.RenderJson(Status, Msg, result)
	}
	if debug {
		fmt.Println("彩種標識:", sLotteryIden)
	}

	///////////
	//判斷系列//
	//////////
	oSeries := models.Series.GetInfo(fmt.Sprintf("%d", oLottery.SeriesId))
	if oSeries.Id == 0 {
		Status = 3003
		Msg = "游戏系列错误"
		this.RenderJson(Status, Msg, result)
	}

	/////////////////////
	//獲取獎金細節和獎金組//
	////////////////////
	aPrizeDetails, sGroupName := models.PrizeGroups.GetPrizeSettingsOfUser(sUserId, fmt.Sprintf("%d", oLottery.Id))
	if len(aPrizeDetails) < 1 {
		Status = 3004
		Msg = "找不到奖金组"
		this.RenderJson(Status, Msg, result)
	}
	if debug {
		fmt.Println("oSeries===>", oSeries)
		fmt.Println("aPrizeDetails==>>", aPrizeDetails)
	}

	//////////////
	//獲取系統配置//
	/////////////
	sPrizeLimit := models.SysConfig.GetPrizeByItem("bet_max_prize")
	sMinPrizeGroup := models.SysConfig.GetPrizeByItem("player_min_grize_group")
	sBetSubmitCompress := models.SysConfig.GetPrizeByItem("bet_submit_compress")
	if len(sPrizeLimit) < 1 || len(sMinPrizeGroup) < 1 || len(sBetSubmitCompress) < 1 {
		Status = 3005
		Msg = "系统配置错误"
		this.RenderJson(Status, Msg, result)
	}
	iMaxPrizeGroup, _ := strconv.Atoi(sGroupName)
	iMinPrizeGroup, _ := strconv.Atoi(sMinPrizeGroup)

	sTerminalId := webLoginMap["terminal"]
	iTerminalId, _ := strconv.Atoi(sTerminalId)

	///////////////////////
	//獲取玩法信息和獎金配置//
	//////////////////////
	sWayGroups, aPrizeSettings := models.WayGroups.GetWayGroupSettings(oSeries.Id, iTerminalId, oLottery.Id, aPrizeDetails, sPrizeLimit, sUserId)
	if debug {
		fmt.Println("aPrizeSettings======>>", aPrizeSettings)
	}

	if sWayGroups == nil {
		Status = 3006
		Msg = "找不到玩法组"
		this.RenderJson(Status, Msg, result)
	}

	//////////////
	//獲取傭金詳情//
	/////////////
	aOptionalPrizeSettings := models.PrizeGroups.GetPrizeCommissions(oSeries, iMaxPrizeGroup, iMinPrizeGroup)

	var aIssues []map[string]interface{}
	var iDefaultMultiple int = 0
	var iMaxBetGroup int = 0
	var fDefaultCoeffcient float64 = 0

	//////////////
	//獲取投注獎期//
	//////////////
	if oLottery.IsInstant < 1 {
		aIssues = models.Issues.GetIssuesForBet(oLottery, 0)
	}

	if oUser.BetCoefficient <= 0 {
		iDefaultMultiple = 1
	}

	if debug {
		fmt.Println("iDefaultMultiple======>>", iDefaultMultiple)
	}

	////////////
	//判斷獎金組//
	////////////
	if oLottery.MaxBetGroup > 0 {
		iMaxBetGroup = int(oLottery.MaxBetGroup)
	} else {
		iMaxBetGroup = oSeries.MaxBetGroup
	}

	iPrizeGroup, _ := strconv.Atoi(oUser.PrizeGroup)
	if iMaxBetGroup >= iPrizeGroup {
		iMaxBetGroup = iPrizeGroup
	}

	if oUser.BetCoefficient <= 0 {
		fDefaultCoeffcient = 1.000
	} else {
		fDefaultCoeffcient = oUser.BetCoefficient
	}
	sDefaultCoeffcient := fmt.Sprintf("%.3f", fDefaultCoeffcient)
	if debug {
		fmt.Println("iMaxBetGroup======>>", iMaxBetGroup)
		fmt.Println("sDefaultCoeffcient======>>", sDefaultCoeffcient)
	}

	aIssuesData := models.Issues.GetIssueListForRefreshNew(fmt.Sprintf("%d", oLottery.Id))
	if debug {
		fmt.Println("aIssuesData======>>", aIssuesData)
	}

	//本機公網ip
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	iNow := time.Now().In(shanghai).Unix()
	sUrl := NetIp

	sLotName := strings.ToLower(oLottery.Name)
	if sZhName, ok := models.LotZhName[sLotName]; ok {
		sLotName = sZhName
	}
	mGameInfo := map[string]interface{}{
		"gameId":                oLottery.Id,                      //彩种id号
		"gameSeriesId":          oLottery.SeriesId,                //系列id
		"gameNameEn":            oLottery.Identifier,              //彩种英文名
		"gameNameCn":            sLotName,                         //彩种中文名
		"wayGroups":             sWayGroups,                       //玩法组
		"defaultMethodId":       oSeries.DefaultWayId,             //选号盘中默认展开玩法
		"prizeSettings":         aPrizeSettings,                   //玩法奖金配置
		"uploadPath":            sUrl + "/bets/upload-bet-number", //注单导入上传和解析结果返回地址
		"jsPath":                "/assets/js-min/",                //个玩法前端js加载路径
		"submitUrl":             sUrl + "/game/bet.do",            //下注提交地址
		"loadDataUrl":           sUrl + "/game/load-data.do",      //数据加载地址, 数据同：/game/load-data.do
		"loadIssueUrl":          sUrl + "/game/load-numbers.do",   //历史奖期数据加载地址（返回 issueHistory 数据）
		"optionalPrizes":        aOptionalPrizeSettings,           //详见奖金原始配置数据
		"currentTime":           iNow,                             //当前时间戳
		"availableCoefficients": models.Coefficients,              //可用投注模式配置[数组]
		"defaultMultiple":       iDefaultMultiple,                 //默认投注倍数
		"defaultCoefficient":    sDefaultCoeffcient,               //默认投注模式(2角)
		"prizeLimit":            sPrizeLimit,                      //系统设定单注最大奖金限制
		"maxPrizeGroup":         iMaxBetGroup,                     //最大奖金组
		"betSubmitCompress":     sBetSubmitCompress,               //是否开启数据压缩

		"traceMaxTimes": len(aIssues), //max追号期数
		"gameNumbers":   aIssues,      //新奖期配置数据[数组]详见新奖期配置说明，期号和截止时间
		"issueHistory":  aIssuesData,  //奖期历史数据[数组]，包含历史期号/开奖号码/官方开奖时间信息
	}

	if len(aIssues) >= 1 {
		mGameInfo["currentNumber"] = aIssues[0]["number"]         //当前期号
		mGameInfo["currentNumberTime"] = aIssues[0]["time_stamp"] //当期截止时间

	}

	if mGameInfo == nil {
		Status = 3007
		Msg = "彩种数据缺失"
		result = mGameInfo
	}

	this.RenderJson(Status, Msg, mGameInfo)
}

// @Title LoadIssues
// @Description 历史奖期数据加载地址(20条),author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /load-numbers.do [post]
func (this *PublicController) LoadIssues() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			this.RenderJson(500, "系统错误", err)
		}

	}()
	Status := 200
	Msg := "success"
	var res interface{}

	webLoginMap := this.WebLogin()

	sLotteryIden := webLoginMap["lottery"] //彩种标记
	if len(sLotteryIden) < 1 {
		this.RenderJson(301, "彩种数据不能为空", res)
	}

	///////////
	//判斷彩種//
	//////////
	oLottery := models.Lotteries.GetInfoByIdentifier(sLotteryIden) //获取彩种数据
	if oLottery.Id < 1 {
		Status = 501
		Msg = "彩种不存在"
		this.RenderJson(Status, Msg, res)
	}

	var aIssues []map[string]interface{}
	var issueHistoryData interface{}

	//////////////
	//獲取投注獎期//
	//////////////
	if oLottery.IsInstant < 1 {
		aIssues = models.Issues.GetIssuesForBet(oLottery, 0)
	}

	aIssuesData := models.Issues.GetIssueListForRefreshNew(fmt.Sprintf("%d", oLottery.Id))
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	iNow := time.Now().In(shanghai).Unix()
	if len(aIssues) >= 1 {
		issueHistoryData = map[string]interface{}{
			"currentNumber":     aIssues[0]["number"],
			"currentNumberTime": aIssues[0]["time_stamp"],
			"currentTime":       iNow,
			"issueHistory":      aIssuesData,
		}
	} else {
		issueHistoryData = aIssuesData
	}

	res = issueHistoryData
	this.RenderJson(200, "success", res)
}

// @Title NoticeList
// @Description 前端请求加载公告标题和内容数据(20条),author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /notice_list.do [post]
func (this *PublicController) NoticeList() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			this.RenderJson(500, "系统错误", err)
		}

	}()

	Status := 200
	Msg := "success"
	var result interface{}

	result = models.Bulletin.GetInfo()
	this.RenderJson(Status, Msg, result)
}

// @Title PlatPrizeData,author(leon）
// @Description 平台最新中奖信息和昨日总派奖金额,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /plat_prize_data.do [post]
func (this *PublicController) PlatPrizeData() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			this.RenderJson(500, "系统错误", err)
		}

	}()

	result := models.PrjPrizeSet.GetPrizeDetails(10)
	this.RenderJson(200, "success", result)
}

// @Title GetGameMenu,author(leon）
// @Description 以分组形式返回游戏列表,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /get_game_menu.do [post]
func (this *PublicController) GetGameMenu() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			this.RenderJson(500, "系统错误", err)
		}

	}()

	Status := 200
	Msg := "success"
	var result interface{}
	var aLotteryGroups = map[string]map[string]interface{}{}

	webLoginMap := this.WebLogin()

	sUserId := webLoginMap["user_id"]
	sTerminalId := webLoginMap["terminal"]
	if len(sTerminalId) < 1 || len(sUserId) < 1 {
		Status = 3010
		Msg = "参数有误"
		this.RenderJson(Status, Msg, result)

	}
	iUserId, _ := strconv.Atoi(sUserId)
	oUser, _ := dao.GetUsersById(iUserId) //获取用户数据
	if oUser.Id == 0 {
		Status = 3011
		Msg = "用戶信息错误"
		this.RenderJson(Status, Msg, result)
	}

	//優先從緩存裏讀取
	sCacheKey := models.CompileUserLotteryMethodCacheKey()
	sCacheData := redisClient.Redis.HashReadField(sCacheKey, sUserId)
	if len(sCacheData) < 1 {

		//////////////
		//獲取配置文件//
		/////////////
		aLotteries := models.Lotteries.GetAll()
		mGameSettings := configs.Setting
		aSpecialGroups := mGameSettings["settings"].(map[string][]int)
		aHotIds := aSpecialGroups["hot"]
		aNewIds := aSpecialGroups["new"]
		a24hIds := aSpecialGroups["24h"]
		aNewWinIds := aSpecialGroups["new_win"]

		mLotteryGroupNames := mGameSettings["games"].(map[string]string)
		aGroups := mGameSettings["groups"].(map[string][]int)

		iStatus := models.STATUS_AVAILABLE
		if oUser.IsTester == 1 {
			iStatus = models.STATUS_AVAILABLE_FOR_TESTER
		}

		//////////////////
		//處理用戶玩法數據//
		/////////////////

		//用戶玩法
		aUserLotteriesMethods := models.UserPrizeSets.GetUserLotteryMethod(sUserId)

		//從所有彩種過濾出符合彩種的數據
		var aAvailableLotteryIds []int
		var aAvailaBleLotteryObj = map[int]*dao.Lotteries{}
		for _, oLottery := range aLotteries {
			if oLottery.Status == uint8(iStatus) {

				//判斷是否在用戶玩法內
				if _, ok := aUserLotteriesMethods[fmt.Sprintf("%d", oLottery.Id)]; !ok {
					continue
				}

				aAvailableLotteryIds = append(aAvailableLotteryIds, oLottery.Id)
				aAvailaBleLotteryObj[oLottery.Id] = oLottery
			}
		}

		///////////
		//組裝數據//
		//////////
		for sKey, aIds := range aGroups {
			iCount := len(aIds)
			aLotteryGroups[sKey] = map[string]interface{}{
				"name":      mLotteryGroupNames[sKey],
				"count":     iCount,
				"lotteries": nil,
			}
			var mLotteriesData = map[int]interface{}{}
			for _, iId := range aIds {
				if !common.InArrayInt(iId, aAvailableLotteryIds) {
					iCount--
					aLotteryGroups[sKey]["count"] = iCount
					continue
				}

				//判斷是否在可用彩種對象裏面
				oLot := aAvailaBleLotteryObj[iId]
				if oLot.Id == 0 {
					iCount--
					continue
				}

				//判斷第三方平臺（表已經刪除）
				iHot := 0
				iNew := 0
				i24H := 0
				iNewWin := 0
				iPLatId := 0
				if common.InArrayInt(oLot.Id, aHotIds) {
					iHot = 1
				}
				if common.InArrayInt(oLot.Id, aNewIds) {
					iNew = 1
				}
				if common.InArrayInt(oLot.Id, a24hIds) {
					i24H = 1
				}
				if common.InArrayInt(oLot.Id, aNewWinIds) {
					iNewWin = 1
				}

				sLotName := strings.ToLower(oLot.Name)
				if sZhName, ok := models.LotZhName[sLotName]; ok {
					sLotName = sZhName
				}
				mData := map[string]interface{}{
					"name":       sLotName,
					"identifier": oLot.Identifier,
					"hot":        iHot,
					"new":        iNew,
					"24h":        i24H,
					"new_win":    iNewWin,
					"is_third":   iPLatId,
				}

				mLotteriesData[oLot.Id] = mData
			}
			aLotteryGroups[sKey]["lotteries"] = mLotteriesData

		}

		//緩存
		b, _ := json.Marshal(aLotteryGroups)
		mCache := map[string]string{
			sUserId: string(b),
		}
		redisClient.Redis.HashWrite(sCacheKey, mCache, -1)
	} else {
		json.Unmarshal([]byte(sCacheData), &aLotteryGroups)
	}

	///////////
	//返回數據//
	//////////
	this.RenderJson(Status, Msg, aLotteryGroups)
}

// @Title GetBetRecord
// @Description 以分頁形式返回遊戲記錄,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数(token=)"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /get_bet_record.do [post]
func (this *PublicController) GetBetRecord() {

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

	//获取注单状态
	aProjectStatus := models.ProjectsStatus

	//獲取追號狀態
	aTracesStatus := models.TracesStatus

	//獲取系列
	sUserId := webLoginMap["user_id"]
	iUserId, _ := strconv.Atoi(sUserId)
	oUser, _ := dao.GetUsersById(iUserId) //获取用户数据
	if oUser.Id == 0 {
		Status = 3009
		Msg = "用戶信息错误"
		this.RenderJson(Status, Msg, result)
	}

	//判斷用戶狀態
	iStatus := models.STATUS_AVAILABLE
	if oUser.IsTester == 1 {
		iStatus = models.STATUS_AVAILABLE_FOR_TESTER
	}

	//所有系列所有彩種
	var aAllSeriesSlice = []map[string]interface{}{}
	var aAllLotteriesSlice = []map[string]interface{}{}
	aAllSeries := models.Series.GetAll()
	aAllLotteries := models.Lotteries.GetAll()

	//循環系列
	for _, oSeries := range aAllSeries {
		mSeriesData := map[string]interface{}{}
		sSeriesName := oSeries.Name
		if sSName, ok := models.SeriesName[oSeries.Name]; ok {
			sSeriesName = sSName
		}

		//拼接系列數據
		mSeriesData["id"] = oSeries.Id
		mSeriesData["game_type"] = "number"
		mSeriesData["name"] = sSeriesName
		mSeriesData["identifier"] = oSeries.Identifier
		aAllSeriesSlice = append(aAllSeriesSlice, mSeriesData)

	}

	//循環彩種
	for _, oLotttery := range aAllLotteries {
		if oLotttery.Status == uint8(iStatus) {
			sLotName := strings.ToLower(oLotttery.Name)
			if sZhName, ok := models.LotZhName[sLotName]; ok {
				sLotName = sZhName
			}
			mLottery := map[string]interface{}{
				"id":         oLotttery.Id,
				"series_id":  oLotttery.SeriesId,
				"game_type":  "number",
				"name":       sLotName,
				"identifier": oLotttery.Identifier,
			}
			aAllLotteriesSlice = append(aAllLotteriesSlice, mLottery)
		}
	}

	//transactions_type
	aField := []string{"id", "cn_title"}
	var mTransactionTypes = map[int]interface{}{}
	aTransactionTypes, _ := dao.GetAllTransactionTypes(nil, aField, nil, nil, 0, 1000)
	for _, oTransactionTypes := range aTransactionTypes {
		mTransactionTypes[oTransactionTypes.Id] = oTransactionTypes.CnTitle
	}

	result = map[string]interface{}{
		"series":            aAllSeriesSlice,
		"lotteries":         aAllLotteriesSlice,
		"project_status":    aProjectStatus,
		"traces_status":     aTracesStatus,
		"transaction_types": mTransactionTypes,
		"game_type":         "number",
	}

	this.RenderJson(Status, Msg, result)
}

// @Title PrintProjects
// @Description 打印注单,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数(token=)"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /print_projects.do [post]
func (this *PublicController) PrintProjects() {

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
	if oUser.Id < 1 {
		Status = 3010
		Msg = "用戶信息错误"
		this.RenderJson(Status, Msg, result)
	}

	sBetRecordId := webLoginMap["bet_record_id"]
	if len(sBetRecordId) < 1 {
		Status = 3011
		Msg = "参数错误"
		this.RenderJson(Status, Msg, result)
	}
	iBetRecordId, _ := strconv.Atoi(sBetRecordId)
	oBetRecord, err := dao.GetBetRecordsById(iBetRecordId)
	if oBetRecord.Id < 1 {
		Status = 3012
		Msg = "没有该记录"
		this.RenderJson(Status, Msg, err)
	}

	//获取对应注单详情
	mWhere := map[string]string{
		"bet_record_id": sBetRecordId,
	}
	aProjects, _ := dao.GetAllProjects(mWhere, nil, nil, nil, 0, 100)
	var aData = []map[string]interface{}{}
	for _, mProject := range aProjects {
		if mProject.UserId != int64(iUserId) {
			Status = 3013
			Msg = "没有该权限"
			this.RenderJson(Status, Msg, result)
		}

		oLottery := models.Lotteries.GetInfo(strconv.Itoa(int(mProject.LotteryId)))
		if oLottery.Id < 1 {
			Status = 3014
			Msg = "查询不到该彩种"
			this.RenderJson(Status, Msg, result)
		}

		iTrace := 0
		if mProject.TraceId > 0 {
			iTrace = 1
		}

		var fBetCount float64 = 0
		if mProject.Coefficient > 0 {
			fBetCount = mProject.Amount / mProject.Coefficient / 2
		}
		sCoefficient := fmt.Sprintf("%.3f", mProject.Coefficient)
		sCoefficient = models.Coefficients[sCoefficient]
		mData := map[string]interface{}{
			"id":                 mProject.Id,
			"lottery":            oLottery.Name,
			"serial_number":      mProject.SerialNumber,
			"amount":             fmt.Sprintf("%.3f", mProject.Amount),
			"way":                mProject.Title,
			"issue":              mProject.Issue,
			"bought_at":          mProject.BoughtAt,
			"coefficient":        sCoefficient,
			"bet_count":          fBetCount,
			"multiple":           mProject.Multiple,
			"is_trace":           iTrace,
			"display_bet_number": mProject.DisplayBetNumber,
			"bet_record_id":      sBetRecordId,
		}
		aData = append(aData, mData)
	}

	this.RenderJson(Status, Msg, aData)
}

// @Title Test
// @Description 測試接口,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数(token=)"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /test.do [post]
func (this *PublicController) Test() {
	//	webLoginMap := this.WebLogin()
	oLottery := models.Lotteries.GetInfo("1")
	aIssues := models.Issues.GetIssuesForBet(oLottery, 0)

	this.RenderJson(200, "test", aIssues)
}

// @Title GetDayProjects
// @Description 獲取介入商某天的所有注單，默認昨天,,author(leon）
// @Param	params	query 	controllers.ParamsInputType		true		"param参数(token=)"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /day.do [post]
func (this *PublicController) GetDayProjects() {
	Status := 200
	Msg := "success"
	var result interface{}

	// webLoginMap := this.WebLogin()
	//1.判斷接入商
	//2.判斷日期
	//3.

	this.RenderJson(Status, Msg, result)
}
