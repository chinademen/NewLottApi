package models

import (
	"NewLottApi/dao"
	"common"
	"encoding/json"
	"fmt"
	"lotteryJobs/models"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type tProjects struct {
	TbName string
	Fields string
}

var ProjectsStatus = map[int]string{
	0: "待开奖",
	1: "已撤销",
	2: "未中奖",
	3: "已中奖",
	5: "系统撤销",
}

var Projects = &tProjects{TbName: "projects",
	Fields: "`id`, `merchant_id`, `terminal_id`, `serial_number`, `trace_id`, `user_id`, `username`, `is_tester`, `account_id`, `prize_group`, `lottery_id`, `issue`, `end_time`, `way_id`, `title`, `position`, `bet_number`, `way_total_count`, `single_count`, `bet_rate`, `display_bet_number`, `multiple`, `single_amount`, `amount`, `winning_number`, `prize`, `prize_sale_rate`, `status`, `status_prize`, `status_sync`, `prize_set`, `single_won_count`, `won_count`, `won_data`, `ip`, `proxy_ip`, `bet_record_id`, `canceled_by`, `bought_at`, `canceled_at`, `prize_sent_at`, `bought_time`, `bet_commit_time`, `coefficient`, `created_at`, `updated_at`"}

/**
 * 组合追號注单数据(追號)
 * @param array     aOrder
 * @param SeriesWay oSeriesWay
 * @param Lottery   oLottery
 * @param array     aExtraData
 */
func (m *tProjects) CompileTraceProjectData(oUser *dao.Users, oTraceDetail *dao.TraceDetails, oTrace *dao.Traces, oSeariesWay *dao.SeriesWays, oLottery *dao.Lotteries, sBetTime string, mExt map[string]string) *dao.Projects {

	var oNewProject = new(dao.Projects)
	iBetTime, _ := strconv.Atoi(sBetTime)
	fMultiple, _ := strconv.ParseFloat(oTraceDetail.Multiple, 64)
	iMerchantId, _ := strconv.Atoi(mExt["merchant_id"])

	oNewProject.TraceId = int64(oTrace.Id)
	oNewProject.UserId = int64(oUser.Id)
	oNewProject.Username = oUser.Username
	oNewProject.AccountId = int64(oTrace.AccountId)
	oNewProject.Issue = oTrace.StartIssue
	oNewProject.EndTime = oTraceDetail.EndTime
	oNewProject.Position = oTrace.Position
	oNewProject.Title = oSeariesWay.Name
	oNewProject.SingleCount = oTrace.SingleCount
	oNewProject.BetNumber = oTrace.BetNumber
	oNewProject.DisplayBetNumber = oTrace.DisplayBetNumber
	oNewProject.LotteryId = uint8(oLottery.Id)
	oNewProject.WayId = oSeariesWay.Id
	oNewProject.Coefficient = oTrace.Coefficient
	oNewProject.SingleAmount = oTrace.SingleAmount
	oNewProject.Amount = oTraceDetail.Amount
	oNewProject.Multiple = int(fMultiple)
	oNewProject.SerialNumber = common.Uniqid(fmt.Sprintf("%d", oUser.Id), true)
	oNewProject.MerchantId = iMerchantId
	oNewProject.Status = 0
	oNewProject.StatusPrize = 0

	oNewProject.BoughtAt = common.DateFormat(sBetTime, common.DATE_FORMAT_YMDHIS)
	oNewProject.BoughtTime = iBetTime
	oNewProject.PrizeSet = oTrace.PrizeSet
	oNewProject.PrizeGroup = oTrace.PrizeGroup
	oNewProject.Ip = oTrace.Ip
	oNewProject.ProxyIp = oTrace.ProxyIp
	oNewProject.TerminalId = int(oTrace.TerminalId)
	oNewProject.IsTester = oUser.IsTester
	oNewProject.BetRecordId = oTrace.BetRecordId
	oNewProject.WayTotalCount = uint64(oTrace.WayTotalCount)
	oNewProject.BetRate = float32(oTrace.BetRate)

	return oNewProject
}

/**
 * 直接组合注单数据(無追號)
 * @param array     aOrder
 * @param SeriesWay oSeriesWay
 * @param Lottery   oLottery
 * @param array     aExtraData
 */
func (m *tProjects) CompileProjectData(oUser *dao.Users, oSeariesWay *dao.SeriesWays, oLottery *dao.Lotteries, sBetTime string, mExt, mProject map[string]string) *dao.Projects {

	//實例化對象
	var oNewProject = new(dao.Projects)

	iBetTime, _ := strconv.Atoi(sBetTime)
	iEndTime, _ := strconv.Atoi(mProject["end_time"])
	fMultiple, _ := strconv.ParseFloat(mProject["multiple"], 64)
	iTerminalId, _ := strconv.Atoi(mExt["terminal_id"])
	iMerchantId, _ := strconv.Atoi(mExt["merchant_id"])
	iBetRecordId, _ := strconv.Atoi(mExt["bet_record_id"])
	iSingleCount, _ := strconv.Atoi(mProject["single_count"])
	fSingelAmount, _ := strconv.ParseFloat(mProject["single_amount"], 64)
	fCoefficient, _ := strconv.ParseFloat(mProject["coefficient"], 64)
	iWayTotalCount := SeriesWays.GetTotalNumberCount(oSeariesWay)

	oNewProject.TraceId = 0
	oNewProject.UserId = int64(oUser.Id)                                          //必要
	oNewProject.Username = oUser.Username                                         //必要
	oNewProject.AccountId = int64(oUser.AccountId)                                //必要
	oNewProject.LotteryId = uint8(oLottery.Id)                                    //必要
	oNewProject.Issue = mProject["issue"]                                         //必要
	oNewProject.TerminalId = iTerminalId                                          //必要
	oNewProject.MerchantId = iMerchantId                                          //必要
	oNewProject.Multiple = int(fMultiple)                                         //必要
	oNewProject.Ip = mExt["clientIp"]                                             //必要
	oNewProject.ProxyIp = mExt["proxyIp"]                                         //必要
	oNewProject.Status = 0                                                        //必要
	oNewProject.StatusPrize = 0                                                   //必要
	oNewProject.BoughtAt = common.DateFormat(sBetTime, common.DATE_FORMAT_YMDHIS) //必要
	oNewProject.SingleAmount = fSingelAmount                                      //必要
	oNewProject.Amount = fSingelAmount * fMultiple                                //必要
	oNewProject.WayId = oSeariesWay.Id                                            //必要

	oNewProject.EndTime = iEndTime
	oNewProject.Title = oSeariesWay.Name
	oNewProject.WayTotalCount = uint64(iWayTotalCount)
	oNewProject.SingleCount = iSingleCount
	oNewProject.Position = mProject["position"]

	oNewProject.BetNumber = mProject["bet_number"]
	oNewProject.PrizeGroup = mProject["prize_group"]
	oNewProject.Coefficient = fCoefficient
	oNewProject.BoughtTime = iBetTime
	oNewProject.BetRecordId = uint64(iBetRecordId)
	oNewProject.SerialNumber = common.Uniqid(fmt.Sprintf("%d", oUser.Id), true)
	oNewProject.IsTester = oUser.IsTester

	if sDisplayBetNumber, ok := mProject["display_bet_number"]; ok {
		oNewProject.DisplayBetNumber = sDisplayBetNumber
	} else {
		oNewProject.DisplayBetNumber = mProject["bet_number"]
	}

	if iWayTotalCount == 0 {
		oNewProject.BetRate = 0
	} else {
		oNewProject.BetRate = float32(iSingleCount / iWayTotalCount)
	}

	if sValue, ok := mProject["prize_set"]; !ok {
		b, _ := json.Marshal(UserPrizeSets.GetPrizeSetOfUsers(fmt.Sprintf("%d", oUser.Id), fmt.Sprintf("%d", oLottery.Id), oSeariesWay))
		oNewProject.PrizeSet = string(b)
	} else {
		oNewProject.PrizeSet = sValue
	}
	return oNewProject
}

/*
 * 添加追號注單(追號)
 * @param *dao.Users          oUser
 * @param *dao.Accounts       oAccount
 * @param *dao.Projects       oProject
 * @param *dao.SeriesWays     oSeariesWay
 * @param map[string]string   mExt
 */
func (m *tProjects) AddProject(o orm.Ormer, oUser *dao.Users, oAccount *dao.Accounts, oSeariesWay *dao.SeriesWays, oProject *dao.Projects, mExt map[string]string) (bool, int, string, error) {

	//生成注單
	iProjectId, err := dao.AddProjects(o, oProject)
	if iProjectId < 1 {
		return false, oProject.Id, "保存注單失敗", err
	}
	oProject.Id = int(iProjectId)

	//生成注單賬變
	iReturn, err := TransactionList.AddProjectTransaction(o, oUser, oAccount, oProject, oSeariesWay, TYPE_BET, oProject.Amount, mExt)
	if iReturn != ERRNO_CREATE_SUCCESSFUL {
		return false, oProject.Id, "保存賬變失敗", err
	}

	return true, oProject.Id, "success", nil
}

func (m *tProjects) GetList(sWhere, sField, sOrder string, offset, limit int) []map[string]string {

	if len(sField) == 0 {
		sField = m.Fields
	}

	if offset != 0 {
		offset = (offset - 1) * limit
	}

	return GetList(m.TbName, sWhere, sField, sOrder, offset, limit)
}

func (m *tProjects) GetCount(sWhere string) int {
	return GetCount(m.TbName, sWhere)
}

/*
 * 是否超出风控
 * @param sMechantId      int
 * @param sLotteryId      int
 * @param sSeriesWayId    int
 */
func (m *tProjects) IsOverMoneyLimit(iMechantId, iLotteryId, iSeriesWayId int, fNowAmount float64) bool {
	mLimitMoney := models.SeriesWayLimit.GetInfo(iMechantId, iLotteryId, iSeriesWayId)

	//如果有配置的话
	if len(mLimitMoney) > 0 {
		fSumMoney := models.Project.GetSumMoney(iMechantId, iLotteryId, iSeriesWayId)
		fLimitMoney, _ := strconv.ParseFloat(mLimitMoney["prize"], 64)
		if fSumMoney+fNowAmount > fLimitMoney {
			return false
		}
	}
	return true
}
