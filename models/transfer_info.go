package models

import (
	"NewLottApi/dao"
)

type tTransferInfo struct {
	TbName string
}

var TransferInfo = &tTransferInfo{TbName: "transfer_info"}

func (m *tTransferInfo) GetInfo(sUserId, sOrderNumber string) *dao.TransferInfo {
	mCondition := map[string]string{
		"order_number": sOrderNumber,
		"user_id":      sUserId,
	}
	aTransferInfos, _ := dao.GetAllTransferInfo(mCondition, nil, nil, nil, 0, 1)
	if len(aTransferInfos) > 0 {
		return aTransferInfos[0]
	}
	var empty *dao.TransferInfo
	return empty
}
