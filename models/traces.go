package models

import (
	"NewLottApi/dao"
	"common"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type tTraces struct {
	TbName string
}

const (
	ERRNO_TRACE_ERROR_SAVE_ERROR        = -310
	ERRNO_TRACE_ERROR_LOW_BALANCE       = -214
	ERRNO_TRACE_ERROR_DATA_ERROR        = -320
	ERRNO_TRACE_DETAIL_SAVE_FAILED      = -303
	ERRNO_PRJ_GERENATE_FAILED_NO_DETAIL = -341
	STATUS_USER_STOPED                  = 2
	STATUS_RUNNING                      = 0
)

var TracesStatus = map[int]string{
	STATUS_RUNNING:     "进行中",
	1:                  "已完成",
	STATUS_USER_STOPED: "用户终止",
	3:                  "管理员终止",
	5:                  "系统终止",
}

var Traces = &tTraces{TbName: "accounts"}

/*
 * 針對每個追號進行信息存儲
 * @param *dao.Users                   oUser           用戶信息模型
 * @param *dao.Accounts                oAccount        用戶賬戶模型
 * @param *dao.SeriesWays              oSeriesWay      系列投注方式模型
 * @param *dao.Lotteries               oLottery        彩種模型
 * @param bool                         bStopOnPrized   追號贏了是否終止追號
 * @param map[string]string            mExtraData      額外數據
 * @param string                       sBetTime        投注時間
 * @param map[string]map[string]string mTrace          拼裝的追號數據
 */
func (m *tTraces) CreateTrace(o orm.Ormer, oUser *dao.Users, oAccount *dao.Accounts, oSeriesWay *dao.SeriesWays, oLottery *dao.Lotteries, bStopOnPrized bool, mExtraData map[string]string, sBetTime string, mTrace map[string]map[string]string) (bool, string) {

	oTrace, sFirstIssue := m.CompileTraceData(oUser, oAccount, oSeriesWay, oLottery, bStopOnPrized, mExtraData, sBetTime, mTrace)
	iReturn, sErrMsq := m.AddTrace(o, oUser, oAccount, oTrace, oLottery, oSeriesWay, mTrace["end_times"], mTrace["issues"], mExtraData, sBetTime, sFirstIssue)
	if iReturn != 200 {
		return false, sErrMsq
	}
	return true, ""
}

/**
 * 建立追号任务
 * @param *dao.Users                   oUser           用戶信息模型
 * @param *dao.Accounts                oAccount        用戶賬戶模型
 * @param *dao.SeriesWays              oSeriesWay      系列投注方式模型
 * @param *dao.Lotteries               oLottery        彩種模型
 * @param map[string]string            mEndTimes       結束時間map
 * @param map[string]string            mIssues         投注獎期map
 * @param map[string]string            mExt            額外數據map
 * @param map[string]map[string]string mTrace          拼裝的追號數據
 */
func (m *tTraces) AddTrace(o orm.Ormer, oUser *dao.Users, oAccount *dao.Accounts, oTrace *dao.Traces, oLottery *dao.Lotteries, oSeriesWay *dao.SeriesWays, mEndTimes, mIssues, mExt map[string]string, sBetTime, sFirstIssue string) (int, string) {

	//判斷用戶賬戶餘額是否足夠
	if oAccount.Available < oTrace.Amount {
		err := errors.New("餘額不足")
		return ERRNO_TRACE_ERROR_LOW_BALANCE, err.Error()
	}

	//判斷投注模式是否在默認模式裏面
	if !common.InArray(fmt.Sprintf("%.3f", oTrace.Coefficient), common.MapKeysStr(Coefficients)) {
		err := errors.New("模式錯誤")
		return ERRNO_TRACE_ERROR_SAVE_ERROR, err.Error()
	}

	//儲存追號單信息
	iTraceId, err := dao.AddTraces(o, oTrace)
	if iTraceId < 1 {
		return ERRNO_TRACE_ERROR_SAVE_ERROR, err.Error()
	}
	oTrace.Id = int(iTraceId)

	//生成追號任務並且凍結資金賬變
	iReturn, oAccount := TransactionList.AddTraceTransaction(o, oUser, oAccount, oTrace, oSeriesWay, TYPE_FREEZE_FOR_TRACE, oTrace.Amount, mExt)
	if iReturn != ERRNO_CREATE_SUCCESSFUL {
		err := errors.New("系統錯誤")
		return iReturn, err.Error()
	}

	//生成追號細節詳情
	iTraceDetailId, err, oTraceDetail := TracesDetails.AddDetails(o, oTrace, mIssues, mExt, mEndTimes, sFirstIssue)
	if iTraceDetailId < 1 {
		return ERRNO_TRACE_DETAIL_SAVE_FAILED, err.Error()
	}

	//開始關聯並且生成注單
	oTraceDetail.Id = int(iTraceDetailId)
	iErrno, err := m.GenerateProjectOfIssue(o, oTrace, oUser, oAccount, oTraceDetail, oLottery, oSeriesWay, 1, 1, sBetTime, mExt)
	if err != nil {
		return iErrno, "生成注單錯誤"
	}
	return iErrno, ""
}

/*
 * 根據追號信息生成關聯注單
 * @param *dao.Users                   oUser           用戶信息模型
 * @param *dao.Accounts                oAccount        用戶賬戶模型
 * @param *dao.SeriesWays              oSeriesWay      系列投注方式模型
 * @param *dao.Lotteries               oLottery        彩種模型
 * @param *dao.Traces                  oTrace          追號模型
 * @param *dao.TraceDetails            oTraceDetail    追號詳情模型
 * @param int                          iStatus         結束時間map
 * @param int                          iCount          投注獎期map
 * @param string                       sBetTime        投注時間
 */
func (m *tTraces) GenerateProjectOfIssue(o orm.Ormer, oTrace *dao.Traces, oUser *dao.Users, oAccount *dao.Accounts, oTraceDetail *dao.TraceDetails, oLottery *dao.Lotteries, oSeriesWay *dao.SeriesWays, iStatus, iCount int, sBetTime string, mExt map[string]string) (int, error) {
	var iErrno int = 200

	//拼接注單保存對象
	iSin, err := TracesDetails.GenerateProject(o, oTrace, oUser, oAccount, oTraceDetail, oLottery, oSeriesWay, sBetTime, mExt)
	if iSin != 200 {
		return ERRNO_PRJ_GERENATE_FAILED_NO_DETAIL, err
	}

	//更新追號信息
	iErrno, er := m.UpdateFinishedInformation(o, oTrace, 1, oTraceDetail.Amount)
	if iErrno != 200 {
		return iErrno, er
	}

	return iErrno, nil
}

/*
 * 每次生成注單，更新追號信息
 * @param *dao.Traces                  oTrace              追號模型
 * @param int                          iIncrementCount     增加的數量
 * @param float64                      fIncrementAmount    增加的金額
 */
func (m *tTraces) UpdateFinishedInformation(o orm.Ormer, oTrace *dao.Traces, iIncrementCount int, fIncrementAmount float64) (int, error) {
	if iIncrementCount <= 0 || fIncrementAmount <= 0 {
		return ERRNO_TRACE_ERROR_DATA_ERROR, nil
	}

	//完成獎期數=已經完成數+新增數(一般爲1)
	iFinishedIssueCount := int(oTrace.FinishedIssues) + iIncrementCount
	fFinishedAmount := oTrace.FinishedAmount + fIncrementAmount
	oTrace.Status = STATUS_RUNNING
	oTrace.FinishedAmount = fFinishedAmount
	oTrace.FinishedIssues = uint32(iFinishedIssueCount)

	//是否完成所有追號獎期
	bFinished := false
	if uint32(iFinishedIssueCount) == oTrace.TotalIssues {
		bFinished = true
	}

	err := m.SetFinished(o, oTrace, bFinished)
	if err != nil {
		return ERRNO_TRACE_DETAIL_SAVE_FAILED, err
	}

	return 200, nil
}

/*
 * 更新追號信息
 * @param *dao.Traces                  oTrace        追號模型
 * @param bool                         bByCancel     是否完成 true-->取消追號;false-->取消
 */
func (m *tTraces) SetFinished(o orm.Ormer, oTrace *dao.Traces, bFinished bool) error {
	if bFinished {
		oTrace.Status = int8(STATUS_USER_STOPED)
	}

	err := dao.UpdateTracesById(o, oTrace)
	return err
}

/**
 * 获取预约清单
 *
 * @param *dao.Traces                  oTrace          追號模型
 * @param int                          iStatus         結束時間map
 * @param int                          iCount          投注獎期map
 * @param string                       sBeginIssue     開始獎期
 */
func (m *tTraces) getDetails(oTrace *dao.Traces, iStatus, iCount int, sBeginIssue string) []*dao.TraceDetails {
	mCondition := map[string]string{
		"trace_id": fmt.Sprintf("%d", oTrace.Id),
		"status":   fmt.Sprintf("%d", iStatus),
	}
	if sBeginIssue == "" {
		mCondition["issue__lte"] = sBeginIssue
	}
	aOrderBy := []string{"issue"}
	aTraceDatail, _ := dao.GetAllTraceDetails(mCondition, nil, nil, aOrderBy, 0, int64(iCount))
	return aTraceDatail
}

/**
 * 生成追号任务属性数组
 *
 * @param *dao.Users                   oUser           用戶信息模型
 * @param *dao.Accounts                oAccount        用戶賬戶模型
 * @param *dao.SeriesWays              oSeriesWay      系列投注方式模型
 * @param *dao.Lotteries               oLottery        彩種模型
 * @param bool                         bStopOnPrized   結束時間map
 * @param string                       sBetTime        結束時間map
 * @param map[string]string            mExtraData      投注獎期map
 * @param map[string]map[string]string aTrace          追號數據
 */
func (m *tTraces) CompileTraceData(oUser *dao.Users, oAccount *dao.Accounts, oSeriesWay *dao.SeriesWays, oLottery *dao.Lotteries, bStopOnPrized bool, mExtraData map[string]string, sBetTime string, aTrace map[string]map[string]string) (*dao.Traces, string) {
	mBet := aTrace["bet"]
	fSingleCount, _ := strconv.ParseFloat(mBet["single_count"], 64)
	fCoefficient, _ := strconv.ParseFloat(mBet["coefficient"], 64)
	fSingleAmount := fSingleCount * fCoefficient * float64(oSeriesWay.Price)
	iTotalNumberAccount := m.GetTotalNumberCount(oSeriesWay)
	sDisplayBetNumber := ""
	if sValue, ok := aTrace["bet"]["display_bet_number"]; ok {
		sDisplayBetNumber = sValue
	} else {
		sDisplayBetNumber = aTrace["bet"]["bet_number"]
	}

	//將獎期排序後取出第一條
	sFirstIssue := m.GetFirstIssue(aTrace["issues"])
	if len(sBetTime) < 1 {
		sBetTime = fmt.Sprintf("%d", time.Now().Unix())
	}

	var oTrace = new(dao.Traces)

	iTerminalId, _ := strconv.Atoi(mExtraData["terminal_id"])
	iMerchantId, _ := strconv.Atoi(mExtraData["merchant_id"])
	iSingleCount, _ := strconv.Atoi(mBet["single_count"])

	oTrace.UserId = oUser.Id
	oTrace.Username = oUser.Username
	oTrace.TerminalId = uint8(iTerminalId)
	oTrace.AccountId = int(oUser.AccountId)
	oTrace.IsTester = oUser.IsTester
	oTrace.Title = oSeriesWay.Name
	oTrace.WayTotalCount = iTotalNumberAccount
	oTrace.SingleCount = iSingleCount
	oTrace.BetRate = fSingleCount / float64(iTotalNumberAccount)
	oTrace.DisplayBetNumber = sDisplayBetNumber
	oTrace.LotteryId = uint8(oLottery.Id)
	oTrace.StartIssue = sFirstIssue
	oTrace.Coefficient = fCoefficient
	oTrace.MerchantId = iMerchantId
	oTrace.SerialNumber = common.Uniqid(fmt.Sprintf("%d", oUser.Id), true)
	oTrace.WayId = oSeriesWay.Id
	oTrace.Position = oSeriesWay.Position

	oTrace.PrizeGroup = mBet["prize_group"]
	oTrace.PrizeSet = mBet["prize_set"]
	oTrace.TotalIssues = uint32(len(aTrace["issues"]))
	oTrace.Position = mBet["position"]
	oTrace.BetNumber = mBet["bet_number"]
	oTrace.SingleAmount = fSingleAmount
	oTrace.Amount = fSingleAmount * common.MapSumF(aTrace["issues"])
	oTrace.Ip = mExtraData["clientIP"]
	oTrace.ProxyIp = mExtraData["proxyIP"]
	iBetRecordId, _ := strconv.Atoi(mExtraData["bet_record_id"])
	oTrace.BetRecordId = uint64(iBetRecordId)
	oTrace.BoughtAt = common.DateFormat(sBetTime, common.DATE_FORMAT_YMDHIS)

	oTrace.Status = STATUS_RUNNING
	if bStopOnPrized {
		oTrace.StopOnWon = 1
	} else {
		oTrace.StopOnWon = 0
	}

	return oTrace, sFirstIssue
}

/**
 * 获取總共數目
 * @param *dao.SeriesWays        oSeriesWay         系列投注方式模型
 */
func (m *tTraces) GetTotalNumberCount(oSeriesWay *dao.SeriesWays) int {
	aAllCount := strings.Split(oSeriesWay.AllCount, ",")
	if oSeriesWay.BasicWayId == WAY_MULTI_SEQUENCING {
		return common.ArrayMax(aAllCount) * int(oSeriesWay.DigitalCount)
	}
	return common.ArraySum(aAllCount)
}

/*
 * 獲取投注獎期的第一期
 * @param map[string]string        mIssues         獎期集
 */
func (m *tTraces) GetFirstIssue(mIssues map[string]string) string {
	var mIssueKey = map[int]string{}
	var min int = 999999999999999
	for s, _ := range mIssues {
		var iKey int
		if strings.Contains(s, "-") {
			sKey := strings.Replace(s, "-", "", -1)
			iKey, _ = strconv.Atoi(sKey)
			mIssueKey[iKey] = s
		} else {
			iKey, _ = strconv.Atoi(s)
		}

		mIssueKey[iKey] = s
		if iKey < min {
			min = iKey
		}
	}

	return mIssueKey[min]
}
