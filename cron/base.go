package cron

import (
	"common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
)

var debug bool = false
var logpath string

type BaseTask struct {
	Ctx *context.Context
}

/**
 * 初始化
 */
func init() {

	sDebug := beego.AppConfig.String("task_debug")

	//判断是否开启调试
	if sDebug == "on" {
		debug = true
	}

	logpath = beego.AppConfig.String("logpath")

	//开启日志
	common.SetLogPath(logpath)
}

//初始化路由开启定时计奖任务
func StartTasks() {

	/////////////////
	//定期生成奖期任务//
	/////////////////
	taskSpecMap, err := beego.AppConfig.GetSection("auto_create_issue") //每个月15,5.30分执行
	if err == nil {
		userBonusTask := &AutoCreateIssues{TaskName: taskSpecMap["task_name"], Spec: taskSpecMap["task_spec"]}
		userBonusTask.Run()
	} else {
		beego.Info(err.Error())
	}

	/////////////////////
	//定期生成奖期缓存任务//
	////////////////////
	taskSpecMap, err = beego.AppConfig.GetSection("auto_issue_cache") //每天1点正执行一次
	if err == nil {
		userBonusTask := &AutoIssueCache{TaskName: taskSpecMap["task_name"], Spec: taskSpecMap["task_spec"]}
		userBonusTask.Run()
	} else {
		beego.Info(err.Error())
	}

	toolbox.StartTask()
	defer toolbox.StopTask()
}

/**
 * 新建任务
 */
func (T *BaseTask) NewTask(taskName, taskSpec string, taskFunc toolbox.TaskFunc) *toolbox.Task {
	return toolbox.NewTask(taskName, taskSpec, taskFunc)
}

/**
 * 添加任务
 */
func (T *BaseTask) AddTask(taskName string, tasker *toolbox.Task) {
	toolbox.AddTask(taskName, tasker)
}

/**
 * 任务日志
 */
func (T *BaseTask) Log(LogMsg string) {
	logs.Info(LogMsg)
}
