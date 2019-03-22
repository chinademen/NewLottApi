package log

import (
	"common"
	"os"
	"strings"
	"time"
)

const (
	BET    = 1
	PUBLIC = 2
	USER   = 3
	ADMIN  = 4
	SQL    = 5
)

var Logpath = "logs/"
var LogType = map[int]string{
	BET:    "bet/",
	PUBLIC: "public/",
	USER:   "user/",
	ADMIN:  "admin/",
	SQL:    "sql/",
}

func AddLog(logpath, FiLeName, sMsg string, iType int) {
	LogsWithFileName(logpath, FiLeName, sMsg, iType)
}

/**
* 设置日志目录
 */
func SetLogPath(path string) {
	Logpath = path
}

/**
 * 指定文件名写入文件
 * @logpath		日志目录
 * @FiLeName	日志名称
 * @sMsg        日志內容
 * @sType       日志類型
 */
func LogsWithFileName(logpath, FiLeName, sMsg string, iType int) {

	//獲取日志根目錄
	if len(logpath) < 1 {
		logpath = Logpath
	}

	//如果目錄不存在，創建目錄
	if !common.IsDirExists(logpath) {
		os.MkdirAll(logpath, 0777)
	}

	sStr := LogType[iType]
	logpath += sStr //拼接類型層

	//如果目錄不存在，創建目錄
	if !common.IsDirExists(logpath) {
		os.MkdirAll(logpath, 0777)
	}

	sNow := time.Now().Format(common.DATE_FORMAT_YMD) //2018-01-01 00:00:00
	sMonth := sNow[0:7]
	logpath += sMonth + "/" //拼接時間層

	//如果目錄不存在，創建目錄
	if common.IsDirExists(logpath) != true {
		os.MkdirAll(logpath, 0777)
	}

	isEnd := strings.HasSuffix(logpath, "/")
	if isEnd != true {
		logpath = logpath + "/"
	}

	logFile := logpath + sNow + "-" + FiLeName + ".log" //如2018-08-01
	fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		//		fmt.Println("create log err->", err)
		return
	}

	fout.WriteString(time.Now().Format(common.DATE_FORMAT_YMDHIS) + "\r\n" + sMsg + "\r\n=====================\r\n")
	defer fout.Close()
}

/**
* 判断目录是否存在
 */
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}
