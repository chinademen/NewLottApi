package models

import (
	_ "NewLottApi/dao"
)

type tApiLogs struct {
	TbName string
}

var ApiLogs = &tApiLogs{TbName: "api_logs"}

func (m *tApiLogs) AddLogs() {

}
