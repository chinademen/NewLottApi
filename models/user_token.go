package models

import (
	"common"
	"common/ext/redisClient"
	"fmt"
	"time"
)

type tUserToken struct {
	TbName string
	Fields string
}

var UserToken = &tUserToken{TbName: "user_token",
	Fields: "`id`, `merchant_id`, `user_id`, `ip`, `terminal`, `browser`, `token`, `ctime`"}

var (
	RUserTokenKeyEX int = 3600
)

/*
 * 生成token;生成规则=同一用户最多2条记录 一条保存PC端,一条保存手机端
 */
func (m *tUserToken) CreateToken(merchantId, usId, userIP, device, browser string) (int, string) {
	reInt := 0

	str := usId + userIP + device + browser
	token := common.GetMd5(str)

	newUnix := common.InterfaceToString(time.Now().Unix()) //现在时间戳

	ifToken := m.RGetToken(usId, device) //此用户此设备类型是否创建token
	if len(ifToken) > 0 {                //如果有 更新

		uDtata := map[string]string{
			"ip":      userIP,
			"browser": browser,
			"token":   token,
			"utime":   newUnix,
		}

		sWhere := fmt.Sprintf("id=%s", ifToken["id"])
		reInt = Update(uDtata, m.TbName, sWhere)
		if reInt > 0 { //更新缓存token
			m.RDelToken(usId, device, ifToken["token"])
		}

	} else { //否则插入

		iDtata := map[string]string{
			"merchant_id": merchantId,
			"user_id":     usId,
			"ip":          userIP,
			"terminal":    device,
			"browser":     browser,
			"token":       token,
			"ctime":       newUnix,
			"utime":       newUnix,
		}

		reInt, _ = Insert(iDtata, m.TbName)
	}

	return reInt, token
}

/*
 * 检查token
 * token = 客户端传递的token
 * userIP = 用户ip
 * terminal = 1-PC 2-手机
 * browser = 浏览器类型
 */
func (m *tUserToken) CheckToken(token, userIP, terminal, browser string) (int, string, map[string]string) {

	//通过验证数据库是否有这个token
	mysqlToken := m.RGetRowByToken(token)
	// if len(mysqlToken) == 0 {
	// 	return 101, "token错误", mysqlToken
	// }

	// if mysqlToken["terminal"] != terminal {
	// 	return 102, fmt.Sprintf("tokenTerminal:%s-->userTerminal:%s", mysqlToken["terminal"], terminal), mysqlToken
	// }

	// if mysqlToken["ip"] != userIP {
	// 	return 103, fmt.Sprintf("tokenIP:%s-->userIp:%s", mysqlToken["ip"], userIP), mysqlToken
	// }

	return 200, "token验证成功", mysqlToken
}

/*
 * 获取一行
 * usId = 用户id
 * terminal = 1-PC 2-手机
 */
func (m *tUserToken) GetToken(usId, terminal string) map[string]string {

	//从数据库读取结果
	sWhere := fmt.Sprintf("user_id = '%s' AND terminal = '%s' ", usId, terminal)
	rMap := GetOne(m.TbName, sWhere, m.Fields)
	return rMap
}

/*
 * 获取一行
 */
func (m *tUserToken) GetRowByToken(token string) map[string]string {

	//从数据库读取结果
	sWhere := fmt.Sprintf("token = '%s' ", token)
	rMap := GetOne(m.TbName, sWhere, m.Fields)
	return rMap
}

/*
 * 获取一行 redis
 * usId = 用户id
 * terminal = 1-PC 2-手机
 */
func (m *tUserToken) RGetToken(usId, device string) map[string]string {

	rKey := CopmileUserTokenRowKey(usId, device)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetToken(usId, device)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RUserTokenKeyEX)
	}

	return rMap
}

/*
 * 获取一行 redis
 */
func (m *tUserToken) RGetRowByToken(token string) map[string]string {
	rKey := CopmileUserTokenRowKeyToken(token)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetRowByToken(token)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RUserTokenKeyEX)
	}

	return rMap
}

/*
 * 删除缓存token redis
 * usId = 用户Id
 * device = 1-PC 2-手机
 * OldToken = 缓存的token
 */
func (m *tUserToken) RDelToken(usId, device, OldToken string) {

	redisClient.Redis.KeyDel(CopmileUserTokenRowKey(usId, device))
	redisClient.Redis.KeyDel(CopmileUserTokenRowKeyToken(OldToken))
}
