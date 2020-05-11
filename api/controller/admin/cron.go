package admin

import (
	"api/middleware"
	"api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)


type JobResource struct {
	Name    	string    	`form:"Name"`
	Hid			int    		`form:"Hid"`
	Cmd        	string 		`form:"Cmd"`
	Spec       	string 		`form:"Spec"`
	Status     	int    		`form:"Status"` // 0 暂停，1 正常
	Desc 		string    	`form:"Desc"`
}

// @Tags 任务计划
// @Description 任务列表
// @Summary  任务列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/job [get]
func GetJobList(c *gin.Context) {
}


func AddJob(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"cron-job-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data PermResource

	e 	:= c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, "Invalid Add Job Data", "", c)
		return
	}

	defer func() {
		err := recover(); if err != nil {
			util.JsonRespond(500, "Parse Spec Failed", "", c)
			return
		}
	}()
}

// 定时任务
func RunCronJob(c *cron.Cron) {
	c.Start()
	select {}
}

func CronNewJob() {
	name := cron.New()
	name.AddFunc("*/1 * * * * *", func() { fmt.Println("Every 1min ") } )
	RunCronJob(name)
}

func RemoveJobById(c *cron.Cron, id cron.EntryID)  {
	c.Remove(id)
	c.Stop()
}

