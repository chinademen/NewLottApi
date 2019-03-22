package models

import (
	"common/ext/redisClient"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type tBulletin struct {
	Table
}

var Bulletin = &tBulletin{Table: Table{TableName: "bulletin"}}

type BulletInfo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

/*
 * 獲取最新的十條公告信息
 */
func (m *tBulletin) GetInfo() string {
	sCacheKey := R.ALLBulletIn
	sResult := redisClient.Redis.StringRead(sCacheKey)

	if len(sResult) < 1 {
		o := orm.NewOrm()
		sSql := fmt.Sprintf("select * from %s order by sequence desc limit 0,%d", m.TableName, 20)
		var result []BulletInfo
		o.Raw(sSql).QueryRows(&result)
		json, _ := json.Marshal(result)
		sResult = string(json)

		//緩存10分鍾
		redisClient.Redis.StringReWrite(sCacheKey, sResult, 600)
	}

	return sResult

}
