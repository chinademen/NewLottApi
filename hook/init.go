package hook

import (
	"NewLottApi/models"
	"common"
	"common/ext/redisClient"

	"github.com/astaxie/beego/session"
)

/**
 * 验证異地登陸,和唯一登陸
 */
func CheckLogin(sIp string, store session.Store) (int, string) {

	if store == nil {
		return 600, "系統異常"
	}

	sUsername := common.AssertionData(store.Get("username"))
	sUserId := common.AssertionData(store.Get("user_id"))
	if len(sUsername) < 1 || len(sUserId) < 1 {
		return 601, "無法獲取到用戶信息"
	}

	//判斷登陸ip與數據庫中的登陸ip是否一致
	mUser := models.Users.GetInfo(sUserId)
	if len(mUser) < 1 {
		return 602, "用戶信息錯誤"
	}

	if sIp != mUser["login_ip"] {
		return 603, "異地登錄"
	}

	//如果存在session，则比较session的id是否和reids中保存的一致
	sessionId := store.SessionID()

	//redis中緩存的session id
	sCacheKey := sUsername + "-" + sUserId
	sUserRedisSessionId := redisClient.Redis.HashReadField("session_ids", sCacheKey)
	if sessionId != sUserRedisSessionId {

		//删除登录
		store.Delete("username")
		store.Delete("user_id")
		return 604, "账号在其他地方登录"
	}

	return 200, "已登录"
}
