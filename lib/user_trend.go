package lib

import (
	"common"
	"fmt"
	//	"fmt"
	Lmodels "lotteryJobs/models"
	//	"math"
	"strings"
)

type UserTrend struct{}

//不同统计类型号码截取的起始下标
var PositionsMap = map[string]int{
	"5":  0,
	"4":  1,
	"3":  0,
	"3f": 0,
	"3e": 2,
	"2f": 0,
	"2e": 3,
}

//星对应长度
var NumLenMap = map[string]int{
	"5":  5,
	"4":  4,
	"3":  3,
	"3f": 3,
	"3e": 3,
	"2f": 2,
	"2e": 2,
}

/**
 * [getTrendDataByParams 根据查询参数获取奖期开奖号码, 并生成分析后的走势数据]
 * @param  [integer] $iLotteryId [彩种id]
 * @param  [integer] $iNumType   [位数]
 * @param  [integer] $iBeginTime [起始时间秒数]
 * @param  [integer] $iEndTime   [结束时间秒数]
 * @param  [integer] $iCount     [记录条数]
 * @return [Array]               [返回分析数据]
 */
func (ut *UserTrend) GetTrendDataByParams(lotteryId, numType, beginTime, endTime, count string) (int, string, [][][][]string) {

	reStatus := 901
	reMsg := "没有数据"
	re := [][][][]string{}

	IssuesList := Lmodels.Issues.GetIssuesByParams(lotteryId, beginTime, endTime, count)
	if len(IssuesList) == 0 {
		return reStatus, reMsg, re
	}

	if _, ok1 := NumLenMap[numType]; ok1 == false {

		reStatus = 902
		reMsg = "玩法错误"
		return reStatus, reMsg, re
	}

	oLottery := Lmodels.Lottery.GetInfo(lotteryId) //获取彩种信息
	seriesId := oLottery["series_id"]              //彩种系列
	validNums := oLottery["valid_nums"]            //可选的号码范围

	re = ut.GenerateTrendData(IssuesList, numType, seriesId, validNums) //5星,4星,3星所有玩法的公用数据
	switch numType {

	case "5", "4":

	}

	reStatus = 200
	reMsg = "ok"

	return reStatus, reMsg, re
}

/*
获取开奖走势 所有玩法
IssuesList	奖期开奖号码数据
numType		玩法;5星,前三,后三...

*/
func (ut *UserTrend) GenerateTrendData(IssuesList []map[string]string, numType, seriesId, validNums string) [][][][]string {

	data := [][][][]string{}

	fenGe := "" //开奖号码分割符
	if seriesId == "2" {
		fenGe = " "
	}

	xing := NumLenMap[numType]
	forLen := 3 + xing    //期号+开奖号码+万+千+百+十+个+号码分布
	numFengbu := 2 + xing //号码分布下标

	validNumsArr := strings.Split(validNums, ",") //切割有效号码

	biaoJi := map[int]map[string]int{}
	fenBu := map[string]int{}

	for _, issuesRow := range IssuesList { //奖期信息

		wnNumber := issuesRow["wn_number"]                                //开奖号码
		wnNumber1 := common.Substr(wnNumber, PositionsMap[numType], xing) //根据玩法截取需要的开奖号码
		wnNumberArr := strings.Split(wnNumber1, fenGe)

		data1 := [][][]string{}

		for i := 0; i < forLen; i++ { //号码分布&期号&开奖号码&万千百十个

			data2 := [][]string{}

			if i == 0 || i == 1 || i == numFengbu { //期号&开奖号码&号码分布

				if i == numFengbu { //号码分布

					for _, sValidNum2 := range validNumsArr { //有效号码0-9

						wnNumCishu := common.InArrayNumStr(sValidNum2, wnNumberArr) //值在开奖号码出现的次数
						if wnNumCishu > 0 {                                         //如果是开奖号码

							fenBu[sValidNum2] = 0

						} else {

							if _, ok5 := fenBu[sValidNum2]; ok5 {
								fenBu[sValidNum2] = fenBu[sValidNum2] + 1
							} else {
								fenBu[sValidNum2] = 1
							}
						}

						data3 := []string{
							fmt.Sprintf("%d", fenBu[sValidNum2]), //号码连续未出现次数
							sValidNum2,
							fmt.Sprintf("%d", wnNumCishu),
						}

						data2 = append(data2, data3)

					}

				} else { //期号&开奖号码
					//				data2 = [][]string{}

					data3 := []string{issuesRow["issue"], issuesRow["wn_number"]}
					data2 = append(data2, data3)

					//				data1 = append(data1, data2)
				}

			} else { //万千百十个

				key := i - 2 //万千百十个在走势中的起始位置
				sWnNumber := wnNumberArr[key]
				for _, sValidNum2 := range validNumsArr { //有效号码0-9

					isWnNum := "1"
					if sWnNumber == sValidNum2 { //如果是开奖号码

						if _, ok3 := biaoJi[key]; ok3 == true {

							biaoJi[key][sValidNum2] = 0

						} else {

							biaoJi[key] = map[string]int{
								sValidNum2: 0,
							}
						}

						isWnNum = "0"

					} else {

						if _, ok3 := biaoJi[key]; ok3 == true {

							biaoJi[key][sValidNum2] = biaoJi[key][sValidNum2] + 1

						} else {

							biaoJi[key] = map[string]int{
								sValidNum2: 1,
							}
						}

					}

					data3 := []string{
						fmt.Sprintf("%d", biaoJi[key][sValidNum2]),
						sWnNumber,
						"1",
						isWnNum,
					}

					data2 = append(data2, data3)
				}

				//					if len(data1) < (forLen - 1) {
				//					if (len(data1) < numFengbu) || (len(data1) == 6) {

				//					data1 = append(data1, data2)
				//					}

			}

			data1 = append(data1, data2)
		}

		data = append(data, data1)
	}

	return data
}
