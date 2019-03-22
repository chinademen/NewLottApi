package cron

import (
	"lotteryJobs/thread"
)

type AutoIssueCache struct {
	TaskName string
	Spec     string
	BaseTask
}

/**
 * 自动计奖
 */
func (T *AutoIssueCache) TaskFunc() error {
	thread.MakeIssueListCache(0)
	return nil
}

func (T *AutoIssueCache) Run() {
	taskName := T.TaskName
	taskSpec := T.Spec
	tasker := T.NewTask(taskName, taskSpec, T.TaskFunc)
	T.AddTask(taskName, tasker)
	T.Log("Run Task  Name:[" + taskName + "] taskSpec:[" + taskSpec + "]")
}
