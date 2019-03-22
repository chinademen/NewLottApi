package models

import (
	"common"
	"common/ext/redisClient"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var Mdb *common.Mydb

type Table struct {
	TableName string
}

var StRedis *redisClient.RedisPool
var getBetCountLimitRow int = 100 //投注统计每次获取行数
var GameLobbyUrl string
var R *RedisKeys
var debug bool = false

var HookAES *common.AES //加密设置

type RedisKeys struct {
	ALLLottery      string //所有彩種
	ALLBulletIn     string //公告标题和内容
	ALLBasicMethods string //所有基礎玩法
	ALLBasicWays    string //所有基礎投注方式
	ALLSeries       string //所有系列

	OneLottery     string //一個彩種
	OneSeries      string //一個系列
	OneBasicMethod string //一個玩法
	OneBasicWay    string //一個基礎投注方式
}

/**
* 包入口，建立连接
 */
func init() {

	sqlDebug := beego.AppConfig.String("sql_debug")
	if sqlDebug == "on" { //判断是否开启调试
		debug = true
	}

	//从配置文件读取
	var mysqlAddr = beego.AppConfig.String("mysqlhost")
	var mysqlPort = beego.AppConfig.String("mysqlport")
	var mysqlUser = beego.AppConfig.String("mysqluser")
	var mysqlPwd = beego.AppConfig.String("mysqlpass")
	var mysqlDbname = beego.AppConfig.String("mysqldb")
	var mysqlCharset = beego.AppConfig.String("mysqlcharset")
	var mysqlTimeout = beego.AppConfig.String("mysqltimeout")

	var redisPrefix = beego.AppConfig.String("redis_prefix")
	var redisNetWork = beego.AppConfig.String("redis_network")
	var redisAddr = beego.AppConfig.String("redis_addr")
	var redisPort = beego.AppConfig.String("redis_port")
	var redisPwd = beego.AppConfig.String("redis_pwd")
	var redisDb = beego.AppConfig.String("redis_db")

	GameLobbyUrl = beego.AppConfig.String("game_lobby_url")

	logpath := beego.AppConfig.String("logpath")

	//缓存key赋值
	R = ReturnRedisKeys()

	//建立连接池
	conStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&timeout=%s", mysqlUser, mysqlPwd, mysqlAddr, mysqlPort, mysqlDbname, mysqlCharset, mysqlTimeout)
	Mdb = common.ConnectMysql(conStr, debug)
	orm.RegisterDataBase("default", "mysql", conStr)
	orm.SetMaxIdleConns("default", 1000)
	orm.SetMaxOpenConns("default", 2000)

	//开启redis
	StRedis = redisClient.InitRedis(redisPrefix, redisNetWork, redisAddr, redisPort, redisPwd, redisDb)

	//开启日志
	if len(logpath) < 1 {
		logpath = "/home/log/"
	}
	testbb := common.IsDirExists(logpath)
	var err error
	if testbb != true {
		err = os.MkdirAll(logpath, 0777)
		if err != nil {
			fmt.Println("日志文件夹创建失败:", err)
		}
	}
	isEnd := strings.HasSuffix(logpath, "/")
	if isEnd != true {
		logpath = logpath + "/"
	}
	logpath = logpath + "log"
	logstr := `{"filename":"` + logpath + `","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"separate":["emergency", "alert", "critical", "error", "warning"]}`
	logs.SetLogger(logs.AdapterMultiFile, logstr)

}

/**
* 生成一个表的主健id = 10位时间戳+6个随机
 */
func GetKeyId() string {
	keyid := strconv.FormatInt(time.Now().Unix(), 10)
	randStr := GetRadomRemoval(6)
	keyid = keyid + randStr
	return keyid
}

func GetRand(l int) string {
	var res string
	var list = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	lens := len(list) - 1

	randKey := rand.Intn(lens)
	//干扰随机数规律
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < r.Intn(9); j++ {
		randKey = rand.Intn(lens)
	}
	//正式生成随机数
	for i := 0; i < l; i++ {
		randKey = rand.Intn(lens)
		res = res + list[randKey]
	}
	return res
}

/**
* 使用redis的set，确保每分钟生成的随机数都不同
* @param 	length	int	随机数长度
* @param	param	string	随机数类型
* return	string		返回一个字符串
 */
func GetRadomRemoval(length int) string {
	//生成随机数，并且写入redis，成功即跳出循环
	var res string
	for i := 0; i < 100; i++ {
		res = GetRand(length)
		//将数据写入redis的set
		key := "radom_" + time.Now().Format("200601021504")
		row := redisClient.Redis.SetAddString(key, res)
		//设置key为1分钟过期
		redisClient.Redis.KeyExpire(key, 60, 1)

		if row == 1 {
			break
		}
	}
	return res
}

//Struct2Map struct 2 map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
