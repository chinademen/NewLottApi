package models

import (
	"NewLottApi/dao"
	"encoding/json"
)

type tBetRecords struct {
	Table
}

var BetRecords = &tBetRecords{Table: Table{TableName: "bet_records"}}

/** 保存投注原始记录
 *
 * @param User oUser
 * @param arry aData
 * @param string sCompressedStr
 * @return int64  成功时为记录ID，否则为0
 */
func (m *tBetRecords) CreateRecord(oUser *dao.Users, iLenBalls, iIsTrace, iLotteryId, iTerminalId, iMerchantId int, aData []*BetData, sCompressedStr string) int64 {
	b, _ := json.Marshal(aData)
	oModel := new(dao.BetRecords)
	oModel.UserId = uint64(oUser.Id)
	oModel.MerchantId = uint(iMerchantId)
	oModel.Username = oUser.Username
	oModel.IsTester = oUser.IsTester
	oModel.TerminalId = uint8(iTerminalId)
	oModel.LotteryId = uint(iLotteryId)
	oModel.BetCount = uint(iLenBalls)
	oModel.IsTrace = int8(iIsTrace)
	oModel.BetData = string(b)
	oModel.CompressedData = sCompressedStr
	id, _ := dao.AddBetRecords(oModel)
	return id
}
