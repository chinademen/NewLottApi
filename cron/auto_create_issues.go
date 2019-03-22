package cron

import (
	"lotteryJobs/issue"
)

type AutoCreateIssues struct {
	TaskName string
	Spec     string
	BaseTask
}

/**
 * 自动生成奖期
 */
func (T *AutoCreateIssues) TaskFunc() error {

	issue.OpenCreateIssues()
	return nil
}

func (T *AutoCreateIssues) Run() {
	taskName := T.TaskName
	taskSpec := T.Spec
	tasker := T.NewTask(taskName, taskSpec, T.TaskFunc)
	T.AddTask(taskName, tasker)
	T.Log("Run Task  Name:[" + taskName + "] taskSpec:[" + taskSpec + "]")
}
