package models

import (
	"NewLottApi/dao"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type tTracesDetails struct {
	Table
}

const (
	STATUS_WAITING         = 0
	STATUS_FINISHED        = 1
	STATUS_USER_CANCELED   = 2
	STATUS_ADMIN_CANCELED  = 3
	STATUS_SYSTEM_CANCELED = 4
	STATUS_USER_DROPED     = 5
)

var TraceDetailStatus = map[int]string{
	STATUS_WAITING:         "等待",
	STATUS_FINISHED:        "完成",
	STATUS_USER_CANCELED:   "取消",
	STATUS_ADMIN_CANCELED:  "取消",
	STATUS_SYSTEM_CANCELED: "取消",
	STATUS_USER_DROPED:     "用戶撤銷",
}

var TracesDetails = &tTracesDetails{Table: Table{TableName: "trace_details"}}

/*
 * 添加追號詳情
 * @param *dao.Traces       oTrace     追號模型
 * @param map[string]string mIssueDetails   獎期細節map[獎期]倍數
 * @param map[string]string mExt       額外數據
 * @param map[string]string mEndTimes　結束時間們
 *
 * return int,error,第一期追號細節對象
 */
func (m *tTracesDetails) AddDetails(o orm.Ormer, oTrace *dao.Traces, mIssueDetails, mExt, mEndTimes map[string]string, sFirstIssue string) (int64, error, *dao.TraceDetails) {
	var signal int64
	var e error

	iMerchantId, _ := strconv.Atoi(mExt["merchant_id"])
	var oFistIssueDetail *dao.TraceDetails

	//遍歷獎期生成追號細節
	for sIssue, sMultiple := range mIssueDetails {
		oDetail := new(dao.TraceDetails)
		iEndTime, _ := strconv.Atoi(mEndTimes[sIssue])
		iAcountId := oTrace.AccountId
		oDetail.TraceId = uint64(oTrace.Id)
		oDetail.LotteryId = oTrace.LotteryId
		oDetail.AccountId = int64(iAcountId)
		oDetail.UserId = int64(oTrace.UserId)
		oDetail.Issue = sIssue
		oDetail.EndTime = iEndTime
		oDetail.Multiple = sMultiple
		oDetail.Status = STATUS_WAITING
		oDetail.BoughtAt = oTrace.BoughtAt
		fMultiple, _ := strconv.ParseFloat(sMultiple, 64)
		oDetail.Amount = oTrace.SingleAmount * fMultiple
		oDetail.MerchantId = iMerchantId

		//存儲預約清單
		iDetailId, err := dao.AddTraceDetails(o, oDetail)
		if iDetailId < 1 && err != nil {
			e = err
			break
		}

		if sIssue == sFirstIssue {
			oDetail.Id = int(iDetailId)
			signal = iDetailId
			oFistIssueDetail = oDetail
		}
	}
	return signal, e, oFistIssueDetail
}

/**
 * 完成当期预约的实例化
 *
 * @param *dao.Traces        oTrace
 * @param *dao.Users         oUser
 * @param *dao.Accounts      oAccount
 * @param *dao.TraceDetails  oTraceDetail
 * @param *dao.Lotteries     oLottery
 * @param *dao.SeriesWays    oSeriesWay
 * @param string             sBetTime
 * @param map[string]string  mExt
 */
func (m *tTracesDetails) GenerateProject(o orm.Ormer, oTrace *dao.Traces, oUser *dao.Users, oAccount *dao.Accounts, oTraceDetail *dao.TraceDetails, oLottery *dao.Lotteries, oSeriesWay *dao.SeriesWays, sBetTime string, mExt map[string]string) (int, error) {
	if int(oTraceDetail.Status) != STATUS_WAITING {
		return 700, nil
	}

	//生成注單數據模型
	oProject := Projects.CompileTraceProjectData(oUser, oTraceDetail, oTrace, oSeriesWay, oLottery, sBetTime, mExt)

	//添加注單
	bSuccess, iProjectId, _, err := Projects.AddProject(o, oUser, oAccount, oSeriesWay, oProject, mExt)
	if !bSuccess {
		return 701, err
	}

	//生成注單
	iCode, err := TransactionList.AddProjectTransaction(o, oUser, oAccount, oProject, oSeriesWay, TYPE_UNFREEZE_FOR_BET, oTrace.Amount, mExt)
	if iCode != ERRNO_CREATE_SUCCESSFUL {
		return 702, err
	}

	//更新追號詳情
	oTraceDetail.ProjectId = uint64(iProjectId)
	oTraceDetail.BoughtAt = oTrace.BoughtAt
	oTraceDetail.Status = int8(STATUS_FINISHED)
	err = dao.UpdateTraceDetailsById(o, oTraceDetail)
	if err != nil {
		return 703, err
	}

	return 200, nil
}
