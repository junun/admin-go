package models

import (
	"github.com/robfig/cron/v3"
	"time"
)

var (
	CronMain *cron.Cron
)

type MyJob struct {
	Name string
	Func func()
}

func (j MyJob) Run() {
	j.Func()
}

// 任务表
type Task struct {
	Model
	Name         string
	IsMore		 int
	Active       int
	TriggerType  int
	HostIds      string
	Command      string
	Spec		 string
	Operator     int
	Desc         string
	StartTime 	 time.Time
	EndTime   	 time.Time
}


// 任务执行历史信息表
type TaskHistory struct {
	Model
	TaskId			int
	HostId			int
	Status			int
	RunTime			string
	Output			string
	CreateTime 	    time.Time
}


func init()  {
	CronMain	= cron.New()
	CronMain.Start()
}
