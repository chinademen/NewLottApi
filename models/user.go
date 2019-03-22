package models

import (
	"NewLottApi/dao"
	"common/ext/redisClient"
	"fmt"
	"strconv"
)

type tUsers struct {
	TbName string
	Fields string
}

var Users = &tUsers{TbName: "users",
	Fields: "`id`, `username`, `password`, `fund_password`, `account_id`, `merchant_id`, `prize_group`, `blocked`, `realname`, `nickname`, `email`, `mobile`, `is_tester`, `bet_multiple`, `bet_coefficient`, `login_ip`, `register_ip`, `token`, `signin_at`, `activated_at`, `register_at`, `deleted_at`, `created_at`, `updated_at`"}

var (
	RUsersOneKey   string = "users:string:%s"            //redis 字符串key
	RUsersRowKey   string = "users:row:merId_%s-name_%s" //redis 数据库一行key
	RUsersRowKeyId string = "users:row:id_%s"            //redis 数据库一行key
	RUsersKeyEX    int    = 3600
)

const (
	UserSession        = "UserSession"
	TokenKey           = "TokenKey"
	CurrentUser        = "CurrentUser"
	BLOCK_BUY          = 2
	BLOCK_FUND_OPERATE = 3
)

//TokenValue
type TokenValue struct {
	IP        string `json:"ip"`
	Device    string `json:"device"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	AccountID int    `json:"account_id"`
}

func (t *TokenValue) ToMap() map[string]string {
	result := make(map[string]string)
	result["ip"] = t.IP
	result["device"] = t.Device
	result["user_id"] = strconv.Itoa(t.UserID)
	result["user_name"] = t.UserName
	result["account_id"] = strconv.Itoa(t.AccountID)
	return result
}

//CheckUserAuth 从数据库检测用户帐户
func CheckUserAuth(username, passwrod string) (user *dao.Users, err error) {

	if username == "" {
		return nil, fmt.Errorf("用户名必填")
	}
	queryMap := make(map[string]string)
	queryMap["username"] = username
	users, _ := dao.GetAllUsers(queryMap, nil, nil, nil, 0, 1)
	if len(users) == 0 {
		return nil, fmt.Errorf("无此人")
	}
	user = users[0]
	if user.Password != passwrod {
		//区分大小写
		return nil, fmt.Errorf("帐户密码不对")

	}
	return
}

//GetUserByUsername 根据username获取user
func GetUserByUsername(username string) (user *dao.Users, err error) {

	if username == "" {
		return nil, fmt.Errorf("用户名必填")
	}
	queryMap := make(map[string]string)
	queryMap["username"] = username
	users, _ := dao.GetAllUsers(queryMap, nil, nil, nil, 0, 1)
	if len(users) == 0 {
		return nil, fmt.Errorf("无此人")
	}
	user = users[0]
	return
}

//帐号是否被人使用
func IsKyAccount(identity, username string) (int, string) {
	rInt := 141
	rMsg := "用户名不能为空"
	if len(username) < 1 {
		return rInt, rMsg
	}

	if len(identity) < 1 {
		rInt = 142
		rMsg = "商户不能为空"
		return rInt, rMsg
	}

	merchantRow := Merchants.RGetByIdentity(identity)
	if len(merchantRow) == 0 {
		rInt = 143
		rMsg = "商户不存在"
		return rInt, rMsg
	}

	userRow := Users.RGetByName(merchantRow["id"], username)
	if len(userRow) > 0 {

		rInt = 144
		rMsg = "用户名已被人使用"
		return rInt, rMsg
	}

	rInt = 200
	rMsg = "ok"
	return rInt, rMsg
}

/*
 * 核對用戶密碼
 */
func ChkLogAccount(sIdentity, sUsername, sPassword string) (int, string, map[string]string) {

	d := map[string]string{}
	if len(sUsername) < 1 {
		return 141, "用户名不能为空", d
	}

	if len(sPassword) < 1 {
		return 142, "密码不能为空", d
	}

	if len(sIdentity) < 1 {
		return 143, "商户不能为空", d
	}

	merchantRow := Merchants.RGetByIdentity(sIdentity)
	if len(merchantRow) < 1 {
		return 144, "商户不存在", d
	}

	userRow := Users.RGetByName(merchantRow["id"], sUsername)
	if len(userRow) < 1 {
		return 145, "該用户未注册", d
	}

	if userRow["blocked"] == "1" {
		return 146, "該用户已冻结", d
	}

	//解密数据库密码与前端传来密码匹配
	sDecodePwd := DecodeString(userRow["password"])
	if sDecodePwd != sPassword {
		if debug {
			fmt.Println("数据库-Pwd-->", userRow["password"])
			fmt.Println("加密-sDecodePwd-->", sDecodePwd)
			fmt.Println("接受-sPassword-->", sPassword)
		}
		return 147, "密码错误", d
	}

	return 200, "success", userRow

}

/*
 * 查询用户资料
 * merchantId	商户id
 * @param	string	username		用户名
 */
func (m *tUsers) GetByName(merchantId, username string) map[string]string {

	//从数据库读取结果
	sWhere := fmt.Sprintf("username = '%s' AND merchant_id = '%s'", username, merchantId)
	rMap := GetOne(m.TbName, sWhere, m.Fields)

	return rMap
}

/*
 * 查询用户资料 redis
 * merchantId	商户id
 * @param	string	username		用户名
 */
func (m *tUsers) RGetByName(merchantId, username string) map[string]string {

	rKey := fmt.Sprintf(RUsersRowKey, merchantId, username)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetByName(merchantId, username)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RUsersKeyEX)
	}

	return rMap
}

/*
 * flush redis
 * @merchantId	商户id
 * @param	string	username		用户名
 */
func (m *tUsers) FlushUserCache(sMerchantId, sUsername string) {
	rKey := fmt.Sprintf(RUsersRowKey, sMerchantId, sUsername)
	redisClient.Redis.KeyDel(rKey)
}

/**
 * 将插入数据还原成sql语句
 * @param		d		插入用户资料
 * @return				sql
 */
func (m *tUsers) GetAddSql(d map[string]string) string {

	sql := GetInsertSql(m.TbName, d)
	return sql
}

/**
 * 将插入数据还原成sql语句	有重复id报错
 * @param		d		插入用户资料
 * @return				sql
 */
func (m *tUsers) GetAddOnlySql(d map[string]string) string {

	sql := GetInsertTrueSql(m.TbName, d)
	return sql
}

/**
 * 插入用户主表
 * @param		d		插入用户资料
 */
func (m *tUsers) DbInsert(d map[string]string) (int, string) {

	sqlInt, sqlId := Insert(d, m.TbName)
	return sqlInt, sqlId
}

/**
 * 更新数据
 */
func (m *tUsers) DbUpdate(d map[string]string, sWhere string) int {

	sqlInt := Update(d, m.TbName, sWhere)
	return sqlInt
}

/**
 * 将更新数据还原成sql语句
 */
func (m *tUsers) DbGetUpdateSql(d map[string]string, w string) string {
	sql := Mdb.GetUpdateSql(m.TbName, d, w)
	return sql
}

/**
 * 根据id获取某条记录
 *
 * @param		id		string
 * @return		Basic Method 	map
 */
func (m *tUsers) GetInfo(id string) map[string]string {
	sWhere := fmt.Sprintf("id='%s'", id)
	return GetOne(m.TbName, sWhere, m.Fields)
}

/*
 * 查询一行 redis
 */
func (m *tUsers) RGetById(id string) map[string]string {
	rKey := fmt.Sprintf(RUsersRowKeyId, id)

	//优先读redis缓存
	rMap := redisClient.Redis.HashReadAllMap(rKey)
	if len(rMap) > 0 {
		return rMap
	}

	rMap = m.GetInfo(id)

	//将结果写入redis，缓存1小时
	if len(rMap) > 0 {
		redisClient.Redis.HashWrite(rKey, rMap, RUsersKeyEX)
	}

	return rMap
}
