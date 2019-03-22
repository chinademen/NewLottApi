package models

import (
	"NewLottApi/dao"
	"common"
	"fmt"
	"lotteryJobs/base/ways"
	"strings"
)

type tSeriesWays struct {
	Table
}

type SeriesWaysData struct {
	name          string
	prize         string
	max_multiple  string
	display_prize string
}

var SeriesWays = &tSeriesWays{Table: Table{TableName: "series_ways"}}
var WayClassed = []string{"1", "2", "3", "4", "5", "7", "8", "9", "10"}

/**
 * 整理投注号码，将不必要的分隔符及占位符删除
 *
 * @param  string          sBetNumber
 * @param  string          sSeriesId
 * @param  *dao.SeriesWays oSeriesWay
 * @return string
 */
func (c *tSeriesWays) CompileBetNumber(sBetNumber, sSeriesId string, oBasicWay *dao.BasicWays, oBasicMethod *dao.BasicMethods) string {
	if common.InArray(sSeriesId, WayClassed) {
		return c.CompileBetNumberNew(sBetNumber, oBasicWay, oBasicMethod)
	}
	return ""
}

/*
 * 整理投注號碼
 */
func (c *tSeriesWays) CompileBetNumberNew(sBetNumber string, oBasicWay *dao.BasicWays, oBasicMethod *dao.BasicMethods) string {

	sClass := c.GetWayClass(oBasicWay, oBasicMethod)
	obj := ways.Way

	switch sClass {
	case "WayLottoConstitutedLottoBallOddEven":
		return obj.LottoBallOddEven.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoCombin":
		return obj.LottoCombin.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoContain":
		return obj.LottoContain.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoContainMulti":
		return obj.LottoContainMulti.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoMiddle":
		return obj.LottoMiddle.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoNotContain":
		return obj.LottoNotContain.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoOddEven":
		return obj.LottoOddEven.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoSum":
		return obj.LottoSum.CompileBetNumber(sBetNumber)

	case "WayLottoConstitutedLottoSumBigMidSmall":
		return obj.LottoSumBigMidSmall.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoSumBsde":
		return obj.LottoSumBsde.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoSumOddEven":
		return obj.LottoSumOddEven.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoSumWuXing":
		return obj.LottoSumWuXing.CompileBetNumber(sBetNumber)
	case "WayLottoConstitutedLottoUpMidDown":
		return obj.LottoUpMidDown.CompileBetNumber(sBetNumber)
	case "WayLottoEqualLottoEqual":
		return obj.EqualLottoEqual.CompileBetNumber(sBetNumber)
	case "WayLottoEqualLottoCombin":
		return obj.EqualLottoCombin.CompileBetNumber(sBetNumber)

	case "WayLottoEqualLottoContain":
		return obj.EqualLottoContain.CompileBetNumber(sBetNumber)
	case "WayLottoMultiOneLottoDragon":
		return obj.MultiOneLottoDragon.CompileBetNumber(sBetNumber)
	case "WayLottoMultiOneLottoEqual":
		return obj.MultiOneLottoEqual.CompileBetNumber(sBetNumber)
	case "WayLottoNecessaryConstitutedLottoContain":
		return obj.LottoNecessaryContain.CompileBetNumber(sBetNumber)
	case "WayLottoNecessaryConstitutedLottoCombin":
		return obj.LottoNecessaryCombin.CompileBetNumber(sBetNumber)
	case "WayLottoSeparatedConstitutedLottoEqual":
		return obj.LottoSeparatedEqual.CompileBetNumber(sBetNumber)

	case "WayBigSmallOddEvenBsde":
		return obj.BigSmallOddEvenBsde.CompileBetNumber(sBetNumber)
	case "WayConstitutedBullBull":
		return obj.ConstBullBull.CompileBetNumber(sBetNumber)
	case "WayConstitutedCombin":
		return obj.ConstCombin.CompileBetNumber(sBetNumber)
	case "WayConstitutedContain":
		return obj.ConstContain.CompileBetNumber(sBetNumber)
	case "WayConstitutedDoubleAreaCombin":
		return obj.ConstDoubleAreaCombin.CompileBetNumber(sBetNumber)
	case "WayConstitutedDragon":
		return obj.ConstDragon.CompileBetNumber(sBetNumber)
	case "WayConstitutedForCombin30Combin":
		return obj.ConstForCombin30Combin.CompileBetNumber(sBetNumber)
	case "WayConstitutedNotContain":
		return obj.ConstNotContain.CompileBetNumber(sBetNumber)
	case "WayConstitutedSpecial":
		return obj.ConstSpecial.CompileBetNumber(sBetNumber)
	case "WayConstitutedSscSumBigSmallOddEven":
		return obj.ConstSumBigSmallOddEven.CompileBetNumber(sBetNumber)

	case "WayEnumCombin":
		return obj.EnumCombin.CompileBetNumber(sBetNumber)
	case "WayEnumEqual":
		return obj.EnumEqual.CompileBetNumber(sBetNumber)
	case "WayFunSeparatedConstitutedInterest":
		return obj.FunSepConstInterest.CompileBetNumber(sBetNumber)
	case "WayMixCombinCombin":
		return obj.MixCombinCombin.CompileBetNumber(sBetNumber)
	case "WayMultiOneEqual":
		return obj.MultiOneEqual.CompileBetNumber(sBetNumber)
	case "WayMultiSequencingEqual":
		return obj.MultiSeqEqual.CompileBetNumber(sBetNumber)
	case "WayNecessaryCombin":
		return obj.NecessaryCombin.CompileBetNumber(sBetNumber)

	case "WayRandomConstitutedCombin":
		return obj.RandomConstitutedCombin.CompileBetNumber(sBetNumber)
	case "WayRandomConstitutedDoubleAreaCombin":
		return obj.RandomConstDoubleAreaCombin.CompileBetNumber(sBetNumber)
	case "WayRandomEnumCombin":
		return obj.RandomEnumCombin.CompileBetNumber(sBetNumber)
	case "WayRandomEnumEqual":
		return obj.RandomEnumEqual.CompileBetNumber(sBetNumber)
	case "WayRandomSeparatedConstitutedEqual":
		return obj.RandomSepConstEqual.CompileBetNumber(sBetNumber)
	case "WayRandomSpanEqual":
		return obj.RandomSpanEqual.CompileBetNumber(sBetNumber)
	case "WayRandomSumCombin":
		return obj.RandomSumCombin.CompileBetNumber(sBetNumber)
	case "WayRandomSumEqual":
		return obj.RandomSumEqual.CompileBetNumber(sBetNumber)

	case "WaySectionalizedSeparatedConstitutedArea":
		return obj.SecSepConstArea.CompileBetNumber(sBetNumber)
	case "WaySeparatedConstitutedEqual":
		return obj.SepConstEqual.CompileBetNumber(sBetNumber)
	case "WaySpanEqual":
		return obj.SpanEqual.CompileBetNumber(sBetNumber)
	case "WaySpecialConstitutedSpecial":
		return obj.SpeConstSpecial.CompileBetNumber(sBetNumber)
	case "WaySumCombin":
		return obj.SumCombin.CompileBetNumber(sBetNumber)
	case "WaySumEqual":
		return obj.SumEqual.CompileBetNumber(sBetNumber)
	case "WaySumTailSumTail":
		return obj.SumTailSumTail.CompileBetNumber(sBetNumber)
	case "WayConstitutedK3Same":
		return obj.K3Same.CompileBetNumber(sBetNumber)
	case "WayConstitutedK3TwoSame":
		return obj.K3TwoSame.CompileBetNumber(sBetNumber)
	default:
		return sBetNumber
	}
}

/*
 * 結合算法包組合參數
 * @param *dao.BasicWays       oBasicWay 系列投注方式
 * @param *dao.BasicMethods    oBasicMethod 系列投注方式
 * @param map[string]string    mOrder   投注數據
 */
func (c *tSeriesWays) Count(oBasicWay *dao.BasicWays, oBasicMethod *dao.BasicMethods, mOrder map[string]string) (int, string, map[string]string) {
	var iCount int
	var sDisplayNumber string

	sClass := c.GetWayClass(oBasicWay, oBasicMethod)
	obj := ways.Way

	//組裝模型必須基礎玩法map
	mBasic := make(map[string]string)
	mBasic["choose_count"] = fmt.Sprintf("%d", oBasicMethod.ChooseCount)
	mBasic["min_choose_count"] = fmt.Sprintf("%d", oBasicMethod.MinChooseCount)
	mBasic["min_repeat_time"] = fmt.Sprintf("%d", oBasicMethod.MinRepeatTime)
	mBasic["max_repeat_time"] = fmt.Sprintf("%d", oBasicMethod.MaxRepeatTime)
	mBasic["buy_length"] = fmt.Sprintf("%d", oBasicMethod.BuyLength)
	mBasic["wn_length"] = fmt.Sprintf("%d", oBasicMethod.WnLength)
	mBasic["min_span"] = fmt.Sprintf("%d", oBasicMethod.MinSpan)
	mBasic["span"] = fmt.Sprintf("%d", oBasicMethod.Span)
	mBasic["unique_count"] = fmt.Sprintf("%d", oBasicMethod.UniqueCount)
	mBasic["special_count"] = fmt.Sprintf("%d", oBasicMethod.SpecialCount)
	mBasic["fixed_number"] = fmt.Sprintf("%d", oBasicMethod.FixedNumber)
	mBasic["digital_count"] = fmt.Sprintf("%d", oBasicMethod.DigitalCount)
	mBasic["valid_nums"] = oBasicMethod.ValidNums

	sBetNumber := ""
	sPosition := ""
	if sValue, ok := mOrder["bet_number"]; ok {
		sBetNumber = sValue
	}

	if sValue, ok := mOrder["position"]; ok {
		sPosition = sValue
	}
	oBasic := ways.ReturnBasicMethod(mBasic)

	//算法計算數據
	switch sClass {
	case "WayLottoConstitutedLottoBallOddEven":
		iCount = obj.LottoBallOddEven.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoBallOddEven.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoCombin":
		iCount = obj.LottoCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoContain":
		iCount = obj.LottoContain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoContain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoContainMulti":
		iCount = obj.LottoContainMulti.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoContainMulti.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoMiddle":
		iCount = obj.LottoMiddle.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoMiddle.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoNotContain":
		iCount = obj.LottoNotContain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoNotContain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoOddEven":
		iCount = obj.LottoOddEven.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoOddEven.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoSum":
		iCount = obj.LottoSum.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoSum.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoSumBigMidSmall":
		iCount = obj.LottoSumBigMidSmall.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoSumBigMidSmall.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoSumBsde":
		iCount = obj.LottoSumBsde.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoSumBsde.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoSumOddEven":
		iCount = obj.LottoSumOddEven.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoSumOddEven.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoSumWuXing":
		iCount = obj.LottoSumWuXing.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoSumWuXing.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoConstitutedLottoUpMidDown":
		iCount = obj.LottoUpMidDown.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoUpMidDown.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoEqualLottoEqual":
		iCount = obj.EqualLottoEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.EqualLottoEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoEqualLottoCombin":
		iCount = obj.EqualLottoCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.EqualLottoCombin.GetDisplayBetNumber(sBetNumber, oBasic)

	case "WayLottoEqualLottoContain":
		iCount = obj.EqualLottoContain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.EqualLottoContain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoMultiOneLottoDragon":
		iCount = obj.MultiOneLottoDragon.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.MultiOneLottoDragon.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoMultiOneLottoEqual":
		iCount = obj.MultiOneLottoEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.MultiOneLottoEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoNecessaryConstitutedLottoContain":
		iCount = obj.LottoNecessaryContain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoNecessaryContain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoNecessaryConstitutedLottoCombin":
		iCount = obj.LottoNecessaryCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoNecessaryCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayLottoSeparatedConstitutedLottoEqual":
		iCount = obj.LottoSeparatedEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.LottoSeparatedEqual.GetDisplayBetNumber(sBetNumber, oBasic)

	case "WayBigSmallOddEvenBsde":
		iCount = obj.BigSmallOddEvenBsde.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.BigSmallOddEvenBsde.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedBullBull":
		iCount = obj.ConstBullBull.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstBullBull.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedCombin":
		iCount = obj.ConstCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedContain":
		iCount = obj.ConstContain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstContain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedDoubleAreaCombin":
		iCount = obj.ConstDoubleAreaCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstDoubleAreaCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedDragon":
		iCount = obj.ConstDragon.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstDragon.GetDisplayBetNumber(sBetNumber)
	case "WayConstitutedForCombin30Combin":
		iCount = obj.ConstForCombin30Combin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstForCombin30Combin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedNotContain":
		iCount = obj.ConstNotContain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstNotContain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedSpecial":
		iCount = obj.ConstSpecial.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstSpecial.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedSscSumBigSmallOddEven":
		iCount = obj.ConstSumBigSmallOddEven.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.ConstSumBigSmallOddEven.GetDisplayBetNumber(sBetNumber, oBasic)

	case "WayEnumCombin":
		iCount = obj.EnumCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.EnumCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayEnumEqual":
		iCount = obj.EnumEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.EnumEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayFunSeparatedConstitutedInterest":
		iCount = obj.FunSepConstInterest.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.FunSepConstInterest.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayMixCombinCombin":
		iCount = obj.MixCombinCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.MixCombinCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayMultiOneEqual":
		iCount = obj.MultiOneEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.MultiOneEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayMultiSequencingEqual":
		iCount = obj.MultiSeqEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.MultiSeqEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayNecessaryCombin":
		iCount = obj.NecessaryCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.NecessaryCombin.GetDisplayBetNumber(sBetNumber, oBasic)

	case "WayRandomConstitutedCombin":
		iCount = obj.RandomConstitutedCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomConstitutedCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomConstitutedDoubleAreaCombin":
		iCount = obj.RandomConstDoubleAreaCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomConstDoubleAreaCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomEnumCombin":
		iCount = obj.RandomEnumCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomEnumCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomEnumEqual":
		iCount = obj.RandomEnumEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomEnumEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomSeparatedConstitutedEqual":
		iCount = obj.RandomSepConstEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomSepConstEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomSpanEqual":
		iCount = obj.RandomSpanEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomSpanEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomSumCombin":
		iCount = obj.RandomSumCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomSumCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRandomSumEqual":
		iCount = obj.RandomSumEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.RandomSumEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayRondamMixCombinCombin":
		iCount = obj.MixCombinCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.MixCombinCombin.GetDisplayBetNumber(sBetNumber, oBasic)

	case "WaySectionalizedSeparatedConstitutedArea":
		iCount = obj.SecSepConstArea.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SecSepConstArea.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WaySeparatedConstitutedEqual":
		iCount = obj.SepConstEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SepConstEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WaySpanEqual":
		iCount = obj.SpanEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SpanEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WaySpecialConstitutedSpecial":
		iCount = obj.SpeConstSpecial.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SpeConstSpecial.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WaySumCombin":
		iCount = obj.SumCombin.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SumCombin.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WaySumEqual":
		iCount = obj.SumEqual.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SumEqual.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WaySumTailSumTail":
		iCount = obj.SumTailSumTail.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.SumTailSumTail.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3TwoSame":
		iCount = obj.K3TwoSame.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3TwoSame.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3Sum":
		iCount = obj.K3Sum.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Sum.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3TwoSameMany":
		iCount = obj.K3TwoSameMany.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3TwoSameMany.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3Same":
		iCount = obj.K3Same.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Same.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3SameAll":
		iCount = obj.K3SameAll.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3SameAll.GetDisplayBetNumber()
	case "WayConstitutedK3Diff":
		iCount = obj.K3Diff.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Diff.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3TwoDiff":
		iCount = obj.K3TwoDiff.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3TwoDiff.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3Ordered":
		iCount = obj.K3Ordered.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Ordered.GetDisplayBetNumber()

		//江苏快3=>大小
	case "WayConstitutedK3Bs":
		iCount = obj.K3Bs.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Bs.GetDisplayBetNumber(sBetNumber, oBasic)

	case "WayConstitutedK3Oe":
		iCount = obj.K3Oe.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Oe.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3Contain":
		iCount = obj.K3Contain.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Contain.GetDisplayBetNumber(sBetNumber, oBasic)
	case "WayConstitutedK3Red":
		iCount = obj.K3Red.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Red.GetDisplayBetNumber()
	case "WayConstitutedK3Black":
		iCount = obj.K3Black.Count(sBetNumber, oBasic, sPosition)
		sDisplayNumber = obj.K3Black.GetDisplayBetNumber(sBetNumber, oBasic)
	default:
		iCount = 0
		sDisplayNumber = "no accept class"
	}

	if iCount > 0 {
		mOrder["display_bet_number"] = sDisplayNumber
	}
	return iCount, sDisplayNumber, mOrder
}

/*
 * 獲取模型名
 */
func (m *tSeriesWays) GetWayClass(oBasicWay *dao.BasicWays, oBasicMethod *dao.BasicMethods) string {
	return "Way" + common.CamelCase(oBasicWay.Function) + common.CamelCase(oBasicMethod.WnFunction)
}

/*
 * 獲取虛擬字段total_number_count
 */
func (m *tSeriesWays) GetTotalNumberCount(oSeriesWay *dao.SeriesWays) int {
	aAllCount := strings.Split(oSeriesWay.AllCount, ",")
	if oSeriesWay.BasicWayId == WAY_MULTI_SEQUENCING {
		return common.ArrayMax(aAllCount) * int(oSeriesWay.DigitalCount)
	}

	return common.ArraySum(aAllCount)
}
