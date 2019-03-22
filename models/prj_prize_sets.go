package models

import (
	"NewLottApi/dao"
	"common"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type tPrjPrizeSet struct {
	Table
}

var PrjPrizeSet = &tPrjPrizeSet{Table: Table{TableName: "prj_prize_sets"}}

/*
 * 返回最新的中奖记录
 * iCount int 条数
 */
func (m *tPrjPrizeSet) GetPrizeDetails(iCount int) map[string]interface{} {
	o := orm.NewOrm()

	//获取最新派奖信息
	sSql := fmt.Sprintf("select * from %s where status in (0,1,2) and prize >= 0 order by id desc", m.TableName)
	var aDatas []dao.PrjPrizeSets
	o.Raw(sSql).QueryRows(&aDatas)
	var aNewestPrizeSend []map[string]string
	for _, oPrj := range aDatas {

		//查询彩种名
		oLottery, err := dao.GetLotteriesById(int(oPrj.LotteryId))
		if err != nil || oLottery.Id == 0 {
			continue
		}
		sLotName := strings.ToLower(oLottery.Name)
		if sZhName, ok := LotZhName[sLotName]; ok {
			sLotName = sZhName
		}
		sName := oPrj.Username[:2] + "***" + oPrj.Username[len(oPrj.Username)-2:]
		mInfo := map[string]string{
			"username": sName,
			"lottery":  sLotName,
			"prize":    fmt.Sprintf("%.4f", oPrj.Prize),
		}
		aNewestPrizeSend = append(aNewestPrizeSend, mInfo)
	}

	//获取昨日总派奖金额
	sYesterdayBegin, sYesterdayEnd := common.GetYesterdayDateTime()
	var result []orm.Params
	sSql = fmt.Sprintf("select sum(prize) as sum_prize from %s where sent_at between '%s' and '%s'", m.TableName, sYesterdayBegin, sYesterdayEnd)
	o.Raw(sSql).Values(&result)

	mPlatPrize := map[string]interface{}{
		"prize_list": aNewestPrizeSend,
	}

	if result[0]["sum_prize"] == nil {
		mPlatPrize["total_prize"] = "0"
	} else {
		mPlatPrize["total_prize"] = result[0]["sum_prize"]
	}

	return mPlatPrize
}
