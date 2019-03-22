package models

import (
	"fmt"

	"github.com/astaxie/beego"
)

type All struct {
	GameId         int            `json:"gameId"`         //投注方式id
	IsTrace        int            `json:"isTrace"`        //是否追號
	TraceWinStop   int            `json:"traceWinStop"`   //追號中獎立停
	TraceStopValue int            `json:"traceStopValue"` //追號中獎值
	Amount         string         `json:"amount"`         //投注金額
	Multiple       int            `json:"multiple"`       //投注mutiple
	Orders         map[string]int //獎期對應倍數
	Balls          []*BetData
}

type BetData struct {
	JsId       int    `json:"jsId"`       //JSid
	WayId      string `json:"wayId"`      //投注方式id
	Ball       string `json:"ball"`       //投注号码
	ViewBalls  string `json:"viewBalls"`  //顯示的號碼
	Num        int    `json:"num"`        //注数
	Type       string `json:"type"`       //类型
	OnePrice   int    `json:"onePrice"`   //價格數字2
	Moneyunit  string `json:"moneyunit"`  //投注模式
	Multiple   string `json:"multiple"`   //倍数
	PrizeGroup string `json:"prizeGroup"` //奖金组
	Extra      *Ext
}

type Ext struct {
	Position string `json:"position"` //额外数据
	Seat     string `json:"seat"`     //额外数据
}

/*
 * 返回RedisKeys結構體指針
 */
func ReturnRedisKeys() *RedisKeys {
	oRedisKeys := &RedisKeys{}
	oRedisKeys.ALLLottery = beego.AppConfig.String("lotteries_list")
	oRedisKeys.ALLBulletIn = beego.AppConfig.String("bulletin")
	oRedisKeys.ALLBasicMethods = beego.AppConfig.String("basic_methods_list")
	oRedisKeys.ALLBasicWays = beego.AppConfig.String("basic_ways_list")
	oRedisKeys.ALLSeries = beego.AppConfig.String("series_list")

	oRedisKeys.OneLottery = beego.AppConfig.String("lottery_one")
	oRedisKeys.OneSeries = beego.AppConfig.String("series_one")
	oRedisKeys.OneBasicMethod = beego.AppConfig.String("basic_method-one")

	return oRedisKeys
}

/**
 * 取得销售中的奖期缓存key
 * @param  string sLotteryId 彩种Id
 * @return string sCacheKey
 */
func CompileOnSaleIssueCacheKey(sLotteryId string) string {
	return "on-sale-issue-" + sLotteryId
}

/**
 * 取得销售中的奖期缓存key
 * @param  string sLotteryId 彩种Id
 * @return string sCacheKey
 */
func CompileLastWnNumberCacheKey(sLotteryId string) string {
	return "Lottery-" + sLotteryId
}

/**
 * 取得彩種選好盤缓存key
 * @param  string sLotteryId 彩种Id
 * @return string sCacheKey
 */
func MakeCacheKeyOfLotterySelectionPlate(sLotteryId, sTerminalId string, bForBet bool) string {
	return fmt.Sprintf("selection-plate-%s-%s", sLotteryId, sTerminalId)
}

/**
 * 取得用户彩种缓存key
 * @param  string sLotteryId 彩种Id
 * @return string sCacheKey
 */
func CompileCacheKeyOfUserLottery(sLotteryId, sUserId string) string {
	return "last-wnnumber-" + sUserId + "-" + sLotteryId
}

/**
 * 取得系列獎金設置緩存key
 * @param  string sLotteryId 彩种Id
 * @return string sCacheKey
 */
func CompileCacheKeyOfOpSetting(iSeriesId int) string {
	return fmt.Sprintf("prize-groups-%d", iSeriesId)
}

/**
 * 取得用户獎金設置缓存key
 * @param  string sLotteryId 彩种Id
 * @return string sCacheKey
 */
func CompileCacheKeyOfUserPrizeSetting(sLotteryId, sUserId string) string {
	return fmt.Sprintf("prize-settings-%s-%s", sLotteryId, sUserId)
}

/**
 * 取得用户彩种缓存key
 * @param  string sGroupId 组Id
 * @return string
 */
func MakeCacheKeyOfGroupNew(sGroupId string) string {
	return "new-group-" + sGroupId
}

/**
 * 取得用户獎期缓存key
 * @param  string sLotteryId 彩种Id
 * @return string
 */
func CompileRecentIssuesCacheKey(sLotteryId string) string {
	return "recent-issues-" + sLotteryId
}

/**
 * 取得用户彩种首頁列表缓存key
 * @param  string sLotteryId 彩种Id
 * @return string
 */
func CompileUserLotteryMethodCacheKey() string {
	return "user-lottery-method"
}

/**
 * 取得奖金组缓存key
 * @param  string sClassicPrize
 * @param  string sSeriesId
 * @return string
 */
func CompilePrizeGroupCacheKey(sClassicPrize, sSeriesId string) string {
	return sSeriesId + "-" + sClassicPrize
}

/**
 * 取得historyissue-key
 * @param  string sLotteryId
 */
func CompileHistoryIssueCacheKey(sLotteryId string) string {
	return "histroy-issues-" + sLotteryId
}

/**
 * 取得sys-configs-key
 * @param  string sItem
 */
func CompileSysConfigsCacheKey(sItem string) string {
	return "sys-configs-" + sItem
}

/*
 * 取得所有系列緩存key
 */
func CompileAllSeriesCacheKey() string {
	return R.ALLSeries
}

/*
 * 取得所有系列緩存key
 */
func CompileOneSeriesCacheKey(sSeriesId string) string {
	return fmt.Sprintf(R.OneSeries, sSeriesId)
}

/*
 * 取得所有彩種緩存key
 */
func CompileAllLotteryCacheKey() string {
	return "all-lott:" + R.ALLLottery
}

/*
 * 取得某個彩種緩存key
 */
func CompileOneLotteryCacheKey(sLotteryId string) string {
	return fmt.Sprintf(R.OneLottery, sLotteryId)
}

/**
 * 取得basic-methods-key
 */
func CompileBasicMethodCacheKey() string {
	return R.ALLBasicMethods
}

/**
 * 取得某個基礎玩法緩存key
 */
func CompileOneBasicMethodCacheKey(sBasicMethodId string) string {
	return fmt.Sprintf(R.OneBasicMethod, sBasicMethodId)
}

/**
 * 取得basic-ways-key
 */
func CompileBasicWaysCacheKey() string {
	return R.ALLBasicWays
}

/**
 * 取得transactions-type-key
 */
func CompileTransactionsTypeCacheKey() string {
	return R.ALLBasicWays
}

/**
 * 取得terminal-black-list-key
 */
func CompileTerminalBlackListAllCacheKey(sLotteryId, sTerminalId, sSeriesWayId string) string {
	return fmt.Sprintf("terminal-black-list-%s-%s-%s", sLotteryId, sTerminalId, sSeriesWayId)
}

/**
 * 取得merchants-lottery-close-key
 */
func CompileMerchantsLotteryCloseCacheKey(sMerchantId string) string {
	return fmt.Sprintf("merchants-lottery-close-%s", sMerchantId)
}

/*
 * 取得獎期緩存數據issues-lottery-%d
 */
func CompileIssuesLotteryCacheKey(iLotteryId int) string {
	return fmt.Sprintf("issues-lottery-%d", iLotteryId)
}

/*
 * 取得token-key
 */
func CopmileUserTokenRowKey(sDevice, sUserId string) string {
	return fmt.Sprintf("user_token:row:uid_device_%s_%s", sUserId, sDevice)
}

/*
 * 取得token-key
 */
func CopmileUserTokenRowKeyToken(sToken string) string {
	return fmt.Sprintf("user_token:row:token_%s", sToken)
}

/*
 * 取得接入商一個ip白名單row
 */
func CopmileMerchantIpOneKey(sId, sIp string) string {
	return fmt.Sprintf("merchants_ip:string:merchantsId_%s-ip_%s", sId, sIp)
}

/*
 * 取得接入商一行
 */
func CopmileMerchantIpRowKey(sId string) string {
	return fmt.Sprintf("merchants_ip:row:%s", sId)
}
