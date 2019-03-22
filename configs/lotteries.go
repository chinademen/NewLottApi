package configs

var Setting = map[string]interface{}{}

func init() {
	aLotteries := []int{1, 2, 3, 4, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 48}
	aDebug := []int{}
	bDisplayBetData := true
	iRecentlyCount := 10
	mGroups := map[string][]int{
		"recently": []int{},
		"new":      []int{59, 57, 56, 52, 53, 93, 92, 50},
		"ssc":      []int{13, 16, 1, 7, 20, 26, 5, 6, 4, 28, 36, 3, 45, 46, 49},
		"pk10":     []int{19, 10},
		"11x5":     []int{47, 43, 44, 14, 2, 8, 9, 22, 23, 24, 25, 27, 29, 32, 34, 92, 52, 53, 56, 55, 51},
		"k3":       []int{17, 15, 18, 48, 21, 30, 33, 35},
		"kl12":     []int{39, 42, 59, 41, 40, 50, 57},
		"other":    []int{11, 12, 93, 38, 37},
	}
	mSetting := map[string][]int{
		"hot":     []int{1, 13, 15, 16, 14},
		"new":     []int{52, 53, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44},
		"24h":     []int{13, 16, 19, 20, 28, 26, 17},
		"new_win": []int{26, 43},
	}
	mGames := map[string]string{
		"recently":  "您最近玩的",
		"recommand": "推荐彩种",
		"new":       "新增彩种",
		"ssc":       "时时彩",
		"11x5":      "11选5",
		"pk10":      "PK10",
		"k3":        "快三",
		"other":     "其他(低频,快乐彩)",
		"keno":      "快乐彩",
		"kl12":      "快乐十二,快乐十分",
		"low":       "低频",
		"ag":        "AG真人",
		"ga":        "GA游戏",
	}
	aInstant := []int{26, 43}
	aIndexRecommand := []int{1, 13, 16, 28, 15, 14, 2, 8}
	aAg := []int{72, 74, 78, 102}
	aGa := []int{94, 95, 96, 97, 98, 99, 100, 101, 110, 112, 113, 114}
	aFootBall := []int{31}
	mDisplaySummary := map[string]string{
		"49": "/bmc/",
		"45": "http://www.35daf.com",
	}
	Setting = map[string]interface{}{
		"lotteries":        aLotteries,
		"debug":            aDebug,
		"display_bet_data": bDisplayBetData,
		"recently_count":   iRecentlyCount,
		"groups":           mGroups,
		"settings":         mSetting,
		"games":            mGames,
		"instant":          aInstant,
		"index_recommand":  aIndexRecommand,
		"ag":               aAg,
		"ga":               aGa,
		"football":         aFootBall,
		"display-summary":  mDisplaySummary,
	}
}
