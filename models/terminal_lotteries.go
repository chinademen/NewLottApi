package models

import (
	"lotteryJobs/models"
)

type tTerminalLotteries struct {
	Table
}

var TerminalLotteries = &tTerminalLotteries{Table: Table{TableName: "terminal_lotteries"}}

const (
	Terminal_STATUS_CLOSED                    = 0
	Terminal_STATUS_TESTING                   = 1
	Terminal_STATUS_AVAILABLE_FOR_NORMAL_USER = 2
	Terminal_STATUS_AVAILABLE                 = 3
)

/*
 * 根据给定的状态值返回实际所需要的所有状态值的数组
 * string sId
 */
func (m *tTerminalLotteries) GetLotteryIds(sId string, iStatus int) []string {
	var aIds []string
	aData := models.TerminalLotteries.GetLotteryIds(sId, iStatus)
	for _, mData := range aData {
		aIds = append(aIds, mData["lottery_id"])
	}
	return aIds
}
