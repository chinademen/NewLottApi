package models

import (
	"NewLottApi/dao"
	"common"
	"common/ext/redisClient"
	"encoding/json"
	"fmt"
	"lotteryJobs/models"
	"strconv"
	"strings"
	"time"
)

type tIssues struct {
	Table
}

var Issues = &tIssues{Table: Table{TableName: "issues"}}

/*
 *根据Ｉd获取系列数据
 */
func (c *tIssues) GetInfo(sId string) *dao.Series {
	iId, _ := strconv.Atoi(sId)
	oSeries, _ := dao.GetSeriesById(iId)
	return oSeries
}

/**
 * 返回指定游戏和奖期号的奖期对象
 * @param string sLotteryId
 * @param string sIssue
 * @return oIssue
 */
func (m *tIssues) GetIssue(sLotteryId string, sIssue string) dao.Issues {
	mCondition := map[string]string{
		"lottery_id": sLotteryId,
		"issue":      sIssue,
	}
	aIssues, _ := dao.GetAllIssues(mCondition, nil, nil, nil, 0, 1)
	if len(aIssues) < 1 {
		return dao.Issues{}
	}
	return aIssues[0]
}

/*
 * 根据条件获取所有系列
 */
func (m *tIssues) GetAll(aConditions map[string]string, aFields []string, aSort []string, aOrder []string, iOffset, iLimit int64) []*dao.Series {
	oSeries, _ := dao.GetAllSeries(aConditions, aFields, aSort, aOrder, iOffset, iLimit)
	return oSeries
}

/*
 * 根據彩種獲取投注的獎期
 */
func (m *tIssues) GetIssuesForBet(oLottery *dao.Lotteries, iCount uint) []map[string]interface{} {
	if iCount == 0 {
		iCount = oLottery.DailyIssueCount * 7
	}

	aIssues := m.GetIssueArrayForBetNew(oLottery.Id, time.Now().Unix(), time.Now().AddDate(0, 0, 7).Unix(), iCount, false)
	return aIssues
}

/*
 * 根據彩種獲取投注的獎期數組
 */
func (m *tIssues) GetIssueArrayForBet(iLotteryId int, iBeginTime, iEndTime int64, iCount uint, bStop bool) []map[string]string {

	aIssues := m._GetIssueArrayForBet(iLotteryId, iBeginTime, iEndTime, iCount, false)

	iNow := time.Now().Unix()
	var iEnd int
	var aIssuesForBet = []map[string]string{}
	if len(aIssues) > 0 {
		iEnd, _ = strconv.Atoi(aIssues[0]["end_time"])
	}
	if iNow < int64(iEnd) {
		aIssues = []map[string]string{}
	}
	if len(aIssues) < int(iCount) && !bStop {
		aIssuesForBet = m.GetIssueArrayForBet(iLotteryId, iBeginTime, iEndTime, iCount, true)
	} else {
		m := map[string]int{
			"0": int(iCount),
			"1": len(aIssues),
		}
		iTrueCount := common.MapMin(m)
		for i := 0; i < iTrueCount; i++ {
			mData := map[string]string{
				"number": aIssues[i]["issue"],
				"time":   aIssues[i]["end_time2"],
			}
			aIssuesForBet = append(aIssuesForBet, mData)
		}
	}

	return aIssuesForBet
}

func (m *tIssues) GetIssueArrayForBetNew(iLotteryId int, iBeginTime, iEndTime int64, iCount uint, bStop bool) []map[string]interface{} {

	var aIssuesForBet = []map[string]interface{}{}
	iNow := time.Now().Unix()

	////////////
	//從緩存讀取//
	////////////
	sCacheKey := CompileIssuesLotteryCacheKey(iLotteryId)

	aIssueCache := redisClient.Redis.LRANGE(sCacheKey, "0", "-1") //獲取當天所有緩存獎期

	for _, sIssueCache := range aIssueCache {
		mCache := map[string]interface{}{}
		err := json.Unmarshal([]byte(sIssueCache), &mCache)
		iEnd, _ := strconv.Atoi(mCache["end_time"].(string))
		if err == nil && int64(iEnd) > iNow {
			mNewData := map[string]interface{}{
				"number":     mCache["issue"],
				"time":       mCache["end_time2"],
				"time_stamp": mCache["end_time"],
			}
			aIssuesForBet = append(aIssuesForBet, mNewData)
		}
	}

	if len(aIssuesForBet) > 1 {
		return aIssuesForBet
	}

	//////////////
	//從數據庫讀取//
	//////////////

	aIssuesData := m._GetIssueArrayForBet(iLotteryId, iBeginTime, iEndTime, iCount, false)
	var iEnd int
	if len(aIssuesData) > 0 {
		iEnd, _ = strconv.Atoi(aIssuesData[0]["end_time"])
	}
	if iNow > int64(iEnd) {
		aIssuesData = []map[string]string{}
	}

	for i := 0; i < len(aIssuesData); i++ {
		mData := map[string]interface{}{
			"number":     aIssuesData[i]["issue"],
			"time":       aIssuesData[i]["end_time2"],
			"time_stamp": aIssuesData[i]["end_time"],
		}
		aIssuesForBet = append(aIssuesForBet, mData)
	}

	return aIssuesForBet
}

/**
 * 获取指定游戏的奖期数组
 * @param int $iLotteryId
 * @param int $iCount
 * @param int $iBeginTime
 * @return Collection
 */
func (m *tIssues) _GetIssueArrayForBet(iLotteryId int, iBeginTime, iEndTime int64, iCount uint, bOrderDesc bool) []map[string]string {

	sWhere := fmt.Sprintf(" lottery_id = '%d' and end_time between '%d' and '%d'", iLotteryId, iBeginTime, iEndTime)
	orderBy := ""
	if bOrderDesc {
		orderBy += " offical_time desc"
	} else {
		orderBy += " offical_time asc"
	}

	aIssues := models.Issues.GetList(sWhere, "*", orderBy, 0, int(iCount))
	return aIssues
}

type IssuesInfo struct {
	Issue       string `json:"issue"`
	WnNumber    string `json:"wn_number"`
	OfficalTime string `json:"offical_time"`
}

type IssuesHistory struct {
	Issues       []IssuesInfo `json:"issues"`
	LastNumber   IssuesInfo   `json:"last_number"`
	CurrentIssue string       `json:"current_issue"`
}

/*
 * 獲取12條之前的獎期新
 */
func (m *tIssues) GetIssueListForRefreshNew(sLotteryId string) map[string]interface{} {
	var IssuesData = map[string]interface{}{}
	iNow := time.Now().Unix()
	sSql := fmt.Sprintf("select id,issue,end_time from %s where lottery_id = %s and wn_number = '' and end_time >= %d and status = 1 limit 1", m.TableName, sLotteryId, iNow)
	aRecent := models.Issues.GoSql(sSql)

	sNowIssue := aRecent[0]["issue"]
	if len(aRecent) >= 1 {
		IssuesData["current_issue"] = sNowIssue
	} else {
		IssuesData["current_issue"] = ""
	}

	sSql = fmt.Sprintf("select %s from %s where lottery_id = %s and id < %s order by offical_time desc limit 0,%d", "id,issue,wn_number,offical_time", m.TableName, sLotteryId, aRecent[0]["id"], 12)
	aHistoryIssues := models.Issues.GoSql(sSql)

	IssuesData["issues"] = aHistoryIssues
	return IssuesData
}

/*
 * 獲取12條之前的獎期
 */
func (m *tIssues) GetIssueListForRefresh(iLotteryId int) map[string]interface{} {
	sOnSaleIssue, sOnSaleEndTime := m.GetOnSaleIssue(iLotteryId)
	mLastNumber := m.GetLastWnNumber(iLotteryId)
	aHistoryWnNumbers := m.GetRecentIssues(iLotteryId, 12, sOnSaleEndTime)
	return map[string]interface{}{
		"issues":        aHistoryWnNumbers,
		"last_number":   mLastNumber,
		"current_issue": sOnSaleIssue,
	}
}

/*
 * 获取正在销售中的奖期
 * @param int $iLotteryId
 *
 * return string,string
 */
func (m *tIssues) GetOnSaleIssue(iLotteryId int) (string, string) {
	sIssueInfo := m.GetOnSaleIssueInfo(iLotteryId)
	if len(sIssueInfo) >= 1 {
		aIssueInfo := strings.Split(sIssueInfo, ",")
		return aIssueInfo[0], aIssueInfo[1]
	} else {
		return "", ""
	}
}

/*
 * 獲取最新的獎期
 */
func (m *tIssues) GetLastWnNumber(iLotteryId int) map[string]string {
	sCacheKey := CompileLastWnNumberCacheKey(fmt.Sprintf("%d", iLotteryId))
	mInfo := redisClient.Redis.HashReadAllMap(sCacheKey)
	if len(mInfo) < 1 {
		aIssues := m.GetRecentIssuesFromDb(iLotteryId, 1, 0)
		if len(aIssues) >= 1 {
			mInfo["issue"] = aIssues[0].Issue
			mInfo["wn_number"] = aIssues[0].WnNumber
			mInfo["offical_time"] = aIssues[0].OfficalTime
		}
		redisClient.Redis.HashWrite(sCacheKey, mInfo, 60)
	} else {
		mInfo["issue"] = ""
		mInfo["wn_number"] = ""
		mInfo["offical_time"] = ""
	}
	return mInfo
}

/*
 * 從數據庫獲取最新的獎期
 */
func (m *tIssues) GetRecentIssuesFromDb(iLotteryId, iCount, iSkipCount int) []dao.Issues {
	sCacheKey := CompileRecentIssuesCacheKey(fmt.Sprintf("%d", iLotteryId))
	sCacheData := redisClient.Redis.StringRead(sCacheKey)

	var aResult []dao.Issues
	if len(sCacheData) < 1 {
		iStart := time.Now().AddDate(0, 0, -12).Unix()
		iEnd := time.Now().Unix()
		mCondtions := map[string]string{
			"end_time2__qte": fmt.Sprintf("%d", iStart),
			"end_time2__lte": fmt.Sprintf("%d", iEnd),
			"end_time__lte":  fmt.Sprintf("%d", iEnd),
			"lottery_id":     fmt.Sprintf("%d", iLotteryId),
		}
		aFields := []string{"issue", "wn_number", "offical_time"}
		aOrderBy := []string{"-issue"}
		aResult, _ = dao.GetAllIssues(mCondtions, aFields, aOrderBy, nil, int64(iSkipCount), int64(iCount))
	} else {
		json.Unmarshal([]byte(sCacheData), &aResult)
	}
	return aResult
}

/*
 * 获取正在销售中的奖期
 * @param int $iLotteryId
 *
 * return string
 */
func (m *tIssues) GetOnSaleIssueInfo(iLotteryId int) string {
	sLotteryId := fmt.Sprintf("%d", iLotteryId)
	sCacheKey := CompileOnSaleIssueCacheKey(sLotteryId)
	sIssueInfo := redisClient.Redis.StringRead(sCacheKey)
	aIssues := m.GetIssues(iLotteryId, 1, time.Now().Unix(), sCacheKey)
	if len(aIssues) >= 1 {
		mIssue := aIssues[0]
		sIssueInfo = mIssue["issue"] + "," + mIssue["end_time"] + "," + mIssue["cycle"]

	}

	return sIssueInfo
}

/**
 * 专用于获取在售奖期的奖期列表缓存中取得奖期列表
 * @param int iLotteryId
 * @param int iStartTime
 * @param int iCount
 * @param sCacheKey string
 * @return array &
 */
func (c *tIssues) GetIssues(iLotteryId, iCount int, iStartTime int64, sCacheKey string) []map[string]string {
	iCacheLen := redisClient.Redis.LenList(sCacheKey)
	aIssuesFromCache, _ := c.GetDataFromRedis(iCacheLen, fmt.Sprintf("%d", iLotteryId), sCacheKey)
	aIssues := []map[string]string{}
	if len(aIssuesFromCache) < 1 {
		return aIssues
	}
	i := 0
	for _, aIssue := range aIssuesFromCache {
		iEndTime, _ := strconv.Atoi(aIssue["end_time"])
		if int64(iEndTime) <= iStartTime {
			continue
		}
		i++
		aIssues = append(aIssues, aIssue)
		if i >= iCount {
			break
		}
	}
	return aIssues
}

/*
 * 將生成獎期時緩存到數據庫的獎期取出N條
 */
func (c *tIssues) GetDataFromRedis(iCount int, sLotteryId, sCacheKey string) ([]map[string]string, int) {

	//從獎期中取出１００條緩存獎期
	aFutureIssues := redisClient.Redis.LRANGE(sCacheKey, "0", "-1")
	iNeedCount := iCount - len(aFutureIssues)
	if iNeedCount < 0 {
		iNeedCount = 0
	}
	var aIssues []map[string]string
	i := 0
	for _, sIssueInfo := range aFutureIssues {
		b := []byte(sIssueInfo)
		var result map[string]string
		json.Unmarshal(b, &result)
		aIssues = append(aIssues, result)
		if i == iCount-1 {
			break
		}
		i++
	}

	return aIssues, iNeedCount
}

/*
 * 獲取最新獎期模型
 */
func (c *tIssues) GetRecentIssues(iLotteryId, iCount int, sOnSaleEndTime string) []dao.Issues {
	var aIssues []dao.Issues
	aMoreIssues := c.GetRecentIssuesFromDb(iLotteryId, iCount, 0)
	sCacheKey := CompileRecentIssuesCacheKey(fmt.Sprintf("%d", iLotteryId))
	if len(aMoreIssues) > iCount {
		aIssues = aMoreIssues[0 : iCount-1]
	}

	iOnSaleEndTime, _ := strconv.Atoi(sOnSaleEndTime)
	b, _ := json.Marshal(aIssues)
	redisClient.Redis.StringReWrite(sCacheKey, string(b), int(int64(iOnSaleEndTime)-time.Now().Unix()))
	return aIssues
}

/*
 * 檢查某彩種的投注獎期是否合法
 */
func (c *tIssues) CheckBetIssue(oLottery *dao.Lotteries, iBetTime int64) bool {
	aTraceIssues := c.GetIssuesForBet(oLottery, 0)
	aAvailableEndTimes := common.ArrayColumn(aTraceIssues, "time") //允许的时间
	var iMin int64
	for _, sEndTime := range aAvailableEndTimes {
		EndTime := common.Strtotime(sEndTime)
		iEnd := EndTime.(int64)
		if iMin == 0 {
			iMin = iEnd
		} else {
			if iEnd < iMin {
				iMin = iEnd
			}
		}
	}

	if iBetTime > iMin {
		return false
	}
	return true
}
