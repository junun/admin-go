package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/logging"
	"api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/unknwon/com"
	"strconv"
	"strings"
	"time"
)


type JobResource struct {
	Name    	string    	`form:"Name"`
	HostIds		string    	`form:"HostIds"`
	Command     string 		`form:"Command"`
	Spec       	string 		`form:"Spec"`
	Desc 		string    	`form:"Desc"`
	TriggerType int 		`form:"TriggerType"`
	IsMore 		int 		`form:"IsMore"`
	StartTime   string		`form:"StartTime"`
	EndTime     string    	`form:"EndTime"`
}

type JobPatchResource struct {
	Name		string      `form:"Name"`
	Active     	int    		`form:"Active"` // 0 暂停，1 启用
	ID 			int	    	`form:"id"`
}

type JobInfoResp struct {
	Success int
	Failure int
	Outputs []models.TaskHistory
}

// @Tags 任务计划
// @Description 任务列表
// @Summary  任务列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule [get]
func GetJobList(c *gin.Context) {
	var task []models.Task
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.Task{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&task)
	models.DB.Model(&models.Task{}).Where(maps).Count(&count)

	data["lists"] = task
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 任务计划
// @Description 激活/禁用任务
// @Summary  激活/禁用任务
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.JobResource true "Job信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule [post]
func AddJob(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"schedule-job-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data JobResource
	var task models.Task

	e 	:= c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	// 名字唯一性检查
	models.DB.Model(&models.Task{}).Where("name = ?", data.Name).Find(&task)
	if task.ID > 0 {
		util.JsonRespond(500, "重复的Job名，请检查！", "", c)
		return
	}

	uid,_ 	:= c.Get("Uid")

	task = models.Task{
		Name: data.Name,
		HostIds: data.HostIds,
		Command: data.Command,
		TriggerType: data.TriggerType,
		Spec: data.Spec,
		Desc: data.Desc,
		IsMore: data.IsMore,
		Operator: uid.(int),
	}


	if data.StartTime != "" && data.EndTime != "" {
		task.StartTime, _ 	= time.Parse(time.RFC3339, data.StartTime)
		task.EndTime, _		= time.Parse(time.RFC3339, data.EndTime)
		if  task.StartTime.After(task.EndTime) {
			util.JsonRespond(500, "开始时间不能早于结束时间！", "", c)
			return
		}
	}


	if data.StartTime != "" {
		task.StartTime, _ 	= time.Parse(time.RFC3339, data.StartTime)
	}

	if data.EndTime != "" {
		task.EndTime, _		= time.Parse(time.RFC3339, data.EndTime)
		if  time.Now().After(task.EndTime) {
			util.JsonRespond(500, "结束时间不能晚于现在时间！", "", c)
			return
		}
	}

	e 	= models.DB.Save(&task).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加Job成功", "", c)

}

// @Tags 任务计划
// @Description 任务修改
// @Summary  任务修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "任务ID"
// @Param Data body admin.JobResource true "Job信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule/{id} [put]
func PutJob(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"schedule-job-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data JobResource
	var task models.Task

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edite Job Data", "", c)
		return
	}

	id := c.Param("id")

	// 名字唯一性检查
	models.DB.Model(&models.Task{}).
		Where("name = ?", data.Name).
		Where("id != ?", id).
		Find(&task)

	if task.ID > 0 {
		util.JsonRespond(500, "重复的Job名，请检查！", "", c)
		return
	}

	models.DB.Find(&task, id)
	uid,_ 	:= c.Get("Uid")

	isActive := task.Active

	isNeedReloadJob := task.Command != data.Command || task.Spec != data.Spec

	task.Name 		= data.Name
	task.HostIds 	= data.HostIds
	task.Command 	= data.Command
	task.TriggerType= data.TriggerType
	task.Spec 		= data.Spec
	task.Desc 		= data.Desc
	task.IsMore 	= data.IsMore
	task.Operator 	= uid.(int)

	if data.StartTime != "" && data.EndTime != "" {
		task.StartTime, _ 	= time.Parse(time.RFC3339, data.StartTime)
		task.EndTime, _		= time.Parse(time.RFC3339, data.EndTime)
		if  task.StartTime.After(task.EndTime) {
			util.JsonRespond(500, "开始时间不能早于结束时间！", "", c)
			return
		}
	}

	if data.StartTime != "" {
		task.StartTime, _ 	= time.Parse(time.RFC3339, data.StartTime)
	}

	if data.EndTime != "" {
		endTime, _		:= time.Parse(time.RFC3339, data.EndTime)
		if  time.Now().After(endTime) {
			util.JsonRespond(500, "结束时间不能晚于现在时间！", "", c)
			return
		}
		isNeedReloadJob = isNeedReloadJob || endTime.After(task.EndTime)
		task.EndTime 	= endTime

	}

	e := models.DB.Save(&task).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	if isActive == 1 && isNeedReloadJob {
		hosts := strings.Split(task.HostIds,",")
		for _,h := range hosts {
			id, _ := ReturnEntryidByName(task.Name+h)
			StopJobs(id,task.Name+h)
			AddNewJob(task.TriggerType, task.ID, task.IsMore, h, task.Name+h, task.Spec, task.Command)
		}
	}

	util.JsonRespond(200, "修改Job成功", "", c)
}

// @Tags 任务计划
// @Description 任务删除
// @Summary  任务删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "任务ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule/{id} [del]
func DelJob(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"schedule-job-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var task models.Task

	id := c.Param("id")
	// 检查任务的状态 是否为运行状态
	models.DB.Find(&task, id)

	if task.Active == 1 {
		util.JsonRespond(500, "Job激活中，请先禁用Job！", "", c)
		return
	}

	e := models.DB.Delete(models.Task{}, "id = ?", id).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除Job成功", "", c)
}

// @Tags 任务计划
// @Description 任务历史记录
// @Summary  任务历史记录
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "任务ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule/{id} [get]
func GetJobHisById(c *gin.Context)  {
	var taskhis []models.TaskHistory
	data := make(map[string]interface{})

	id := c.Param("id")

	models.DB.Model(&models.TaskHistory{}).
		Where("task_id = ?", id).
		Order("id DESC").
		Limit(util.GetPageSize(c)).
		Find(&taskhis)

	data["lists"] = taskhis

	util.JsonRespond(200, "", data, c)
}

// @Tags 任务计划
// @Description 任务详情
// @Summary  任务详情
// @Produce  json
// @Param Authorization header string true "token"
// @Param Type query string true "任务调度器类型"
// @Param id path int true "任务ID或者任务历史记录id"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule/{id}/info [get]
func GetJobInfo(c *gin.Context)  {
	var taskhis []models.TaskHistory

	data := make(map[string]interface{})

	id 			:= c.Param("id")
	mytype, _	:= com.StrTo(c.Query("Type")).Int()

	// 任务的最近一次所有主机历史信息
	if mytype == 1 {
		//sql := "select * from task_history where id in (select SUBSTRING_INDEX(group_concat(id order by id desc),',',1) from task_history where task_id=? group by host_id);"

		sql := "select a.* from task_history a join (select max(id) AS id from task_history where task_id=? group by task_id, host_id) b on a.id=b.id;"

		models.DB.Raw(sql, id).Scan(&taskhis)
	}

	// 单个历史信息
	if mytype == 2 {
		models.DB.Model(&models.TaskHistory{}).
			Where("id = ?", id).
			Find(&taskhis)
	}

	if len(taskhis) <= 0  {
		data["lists"] = nil
		util.JsonRespond(200, "", data, c)
		return
	}

	// 处理返回结果
	var res JobInfoResp
	success := 0
	faild	:= 0
	for _, v := range taskhis {
		if v.Status == 0 {
			success++
		} else {
			faild++
		}
	}

	res.Success = success
	res.Failure = faild
	res.Outputs = taskhis

	data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

// @Tags 任务计划
// @Description 激活/禁用任务
// @Summary  激活/禁用任务
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.JobPatchResource true "Job信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/schedule [patch]
func PatchJobStatus(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"schedule-job-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data 	JobPatchResource

	e 	:= c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, "Invalid Patch Job Data", "", c)
		return
	}

	var task 	models.Task
	models.DB.Model(&models.Task{}).Where("id = ?", data.ID).Find(&task)
	if task.ID <= 0 {
		util.JsonRespond(500, "任务不存在，请检查！", "", c)
		return
	}

	// 关闭定时任务 1 根据名字查找 任务id， 2 通过id 移除任务
	if data.Active == 0 {
		hosts := strings.Split(task.HostIds,",")
		for _,h := range hosts {
			id, e := ReturnEntryidByName(data.Name + h)
			if e != nil {
				msg := "没有从redis里面找到任务名字为" + data.Name + "的任务任务信息，禁用定时任务失败，请检查！"
				logging.Error(msg)
				models.MakeNotify(1, 2,"禁用定时任务异常",msg,"" )
				util.JsonRespond(500, msg, "", c)
				return
			}

			StopJobs(id, data.Name+h)
		}
	}

	// 启用任务 1:查询数据，找出任务信息  2:添加任务
	if data.Active == 1 {
		hosts := strings.Split(task.HostIds,",")
		for _,h := range hosts {
			AddNewJob(task.TriggerType, task.ID, task.IsMore, h, task.Name+h, task.Spec, task.Command)
		}
	}

	//e 	= models.DB.Model(&task).Updates(map[string]interface{}{"active": data.Active}).Error
	e 	= models.DB.Model(&models.Task{}).Where("id = ?", data.ID).Updates(map[string]interface{}{"active": data.Active}).Error
	if e != nil {
		msg := "禁用/启用定时任务异常任务名字为" + data.Name + "成功，但是修改数据库状态信息失败，请处理！"
		logging.Error(msg)
		models.MakeNotify(1, 2,"禁用/启用定时任务时修改数据库状态信息失败",msg,"" )
		util.JsonRespond(500, e.Error(), "", c)
	}

	util.JsonRespond(200, "执行成功", "", c)
}

func ReturnEntryidByName(name string) (int, error) {
	strid, e := models.Rdb.HGet(models.CronNameEntryIdKey, name).Result()
	if e != nil {
		return 0, e
	}
	id, _ := strconv.Atoi(strid)
	return id, nil
}

func StopJobs(id int, name string)  {
	models.CronMain.Remove(cron.EntryID(id))
	models.Rdb.HDel(models.CronNameEntryIdKey, name)
}

func Test(c *gin.Context)  {
	fmt.Print(models.CronMain.Entries())
	id, e := ReturnEntryidByName("test-job")
	if e != nil {
		msg := "任务Id为" + strconv.Itoa(id) + "的任务关联的主机id"  +"不存在，请检查！"
		logging.Error(msg)
		models.MakeNotify(1, 2,"定时任务异常",msg,"" )
	}
}

// 定时任务 每天检查任务是否已经过期，过期踢出任务队列，修改任务状态为0。
func CheckJobActive() {
	var tasks []models.Task
	models.DB.Model(&models.Task{}).Where("active=1").Find(&tasks)

	if len(tasks) == 0 {
		return
	}

	for _, v  := range tasks {
		if v.TriggerType == 1 {
			runTime,_ 	:= time.Parse(time.RFC3339, v.Spec)
			// 已经过了任务执行时间
			if time.Now().After(runTime) {
				// 踢出任务队列
				hosts := strings.Split(v.HostIds,",")
				for _,h := range hosts {
					id, _ := ReturnEntryidByName(v.Name + h)
					StopJobs(id, v.Name+h)
				}

				// 修改任务状态为0
				models.DB.Model(&models.Task{}).Where("id = ?", v.ID).Updates(map[string]interface{}{"active": 0})
			}
			continue
		}

		if v.TriggerType == 2 && v.EndTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
			// 已经过了任务执行规定时间
			if time.Now().After(v.EndTime) {

				// 踢出任务队列
				hosts := strings.Split(v.HostIds,",")
				for _,h := range hosts {
					id, _ := ReturnEntryidByName(v.Name + h)
					StopJobs(id, v.Name+h)
				}

				// 修改任务状态为0
				models.DB.Model(&models.Task{}).Where("id = ?", v.ID).Updates(map[string]interface{}{"active": 0})
			}
			continue
		}
	}
}

func StartCronJobsOnBoot()  {
	// 系统任务, 不能手动移除
	models.CronMain.AddFunc("@daily", CheckDomainAndCret)
	models.CronMain.AddFunc("@daily", CheckJobActive)
	//models.CronMain.AddFunc("*/1 * * * ?", CheckDomainAndCret)


	// 启动用户添加的任务,查询所以活动的task
	var task []models.Task

	// 用于 EntryID 可能不一样，所以每次重新启动服务需要先清空 redis name-EntryID Hash
	models.Rdb.Del(models.CronNameEntryIdKey)

	models.DB.Model(models.Task{}).Where("active = 1").Find(&task)
	if len(task) > 0 {
		for _, v := range task {
			if v.TriggerType == 1 {
				runTime,_ 	:= time.Parse(time.RFC3339, v.Spec)
				// 已经过了任务执行时间
				if time.Now().After(runTime) {
					continue
				}
			}

			if v.TriggerType == 2 && v.EndTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				// 已经过了任务执行规定时间
				if time.Now().After(v.EndTime) {
					continue
				}
			}

			hosts := strings.Split(v.HostIds,",")
			for _,h := range hosts {

				AddNewJob(v.TriggerType, v.ID, v.IsMore, h, v.Name+h, v.Spec, v.Command)
			}
		}
	}
}

func AddNewJob(mytype, id, isMore int, hostidstr, name, spec, cmd string)  {
	myspec := spec
	var job models.MyJob
	var taskhis models.TaskHistory
	job.Name	= name

	// 一次性任务, 需要把日期转成 UNIX cron 格式
	if mytype 	== 1 {
		mytime, _	:= time.Parse(time.RFC3339, spec)

		myspec	= strconv.Itoa(mytime.Minute()) + " " +
			strconv.Itoa(mytime.Hour()) + " " +
			strconv.Itoa(mytime.Day()) + " " +
			strconv.Itoa(int(mytime.Month())) + " ?"
	}

	switch hostidstr {
	// 本机任务
	case "0":
		job.Func = func() {
			if isMore == 0 {
				// 当前任务已经有运行了
				if !IsNeedExecuteJob(name) {
					return
				}
			}

			status	:= 0
			startTime 	:= time.Now()
			e, msg 		:= util.ExecRuntimeCmd(cmd)
			if e 	!= nil {
				status 	= 1
			}

			hostid,_ 	:= strconv.Atoi(hostidstr)
			taskhis 	= models.TaskHistory{
				Status: status,
				HostId: hostid,
				Output: msg,
				RunTime: time.Since(startTime).String(),
				TaskId: id,
				CreateTime: time.Now(),
			}

		 	models.DB.Save(&taskhis)
			if isMore == 0 {
				// 任务执行完需要删除锁定
				models.Rdb.HDel(models.CronJobOnRunKey, name)
			}
		}

	// 远端主机
	default:
		// 先检查ip地址是否为本机
		var host models.Host
		hostid,_ 	:= strconv.Atoi(hostidstr)
		models.DB.Model(&models.Host{}).
			Where("id = ?", hostid).
			Where("status = 1").
			Find(&host)

		if host.ID <= 0 {
			msg := "任务Id为" + strconv.Itoa(id) + "的任务关联的主机id" + hostidstr +"不存在，请检查！"
			logging.Error(msg)
			models.MakeNotify(1, 2,"定时任务异常",msg,"" )
			if isMore == 0 {
				// 任务执行完需要删除锁定
				models.Rdb.HDel(models.CronJobOnRunKey, name)
			}
			return
		}

		if CheckIfLocalIp(host.Addres) {
			msg := "检查发现任务运行主机为平台运行的主机，请选择本机执行！"
			logging.Error(msg)
			models.MakeNotify(1, 2,"定时任务异常",msg,"" )
			return
		}

		job.Func = func() {
			if isMore == 0 {
				// 当前任务已经有运行了
				if !IsNeedExecuteJob(name) {
					return
				}
			}

			// ReturnClientConfig
			clientConfig, e := util.ReturnClientConfig(host.Username, "")
			hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)

			Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
			if e != nil {
				msg := "连接远端主机" + hostIp +"异常， 请检查！"
				logging.Error(msg)
				models.MakeNotify(1, 2,"定时任务异常",msg,"" )
				if isMore == 0 {
					// 任务执行完需要删除锁定
					models.Rdb.HDel(models.CronJobOnRunKey, name)
				}
				return
			}

			status 	:= 0
			startTime 	:= time.Now()
			res , e := util.ExecuteCmdRemote(cmd, Scli)
			if e != nil {
				status = 1
			}

			defer Scli.Close()

			hostid,_ 	:= strconv.Atoi(hostidstr)
			taskhis 	= models.TaskHistory{
				Status: status,
				HostId: hostid,
				Output: string(res),
				RunTime: time.Since(startTime).String(),
				TaskId: id,
				CreateTime: time.Now(),
			}

			models.DB.Save(&taskhis)

			if isMore == 0 {
				// 任务执行完需要删除锁定
				models.Rdb.HDel(models.CronJobOnRunKey, name)
			}
		}
	}

	entryID, e := models.CronMain.AddJob(myspec, job)
	if e != nil {
		return
	}

	models.Rdb.HSet(models.CronNameEntryIdKey, name, strconv.Itoa(int(entryID)))
}

// 检查是否需要执行任务，用于多实例部署同时运行多个任务情况
func IsNeedExecuteJob(jobName string) bool {
	if !models.Rdb.HExists(models.CronJobOnRunKey, jobName).Val() {
		models.Rdb.HSet(models.CronJobOnRunKey, jobName, "1")
		return true
	}
	return false
}

func CheckIfLocalIp(ip string) bool {
	e 	:= models.Rdb.Exists(models.ServerLocalRunIpKey).Val()
	if e != 1 {
		ips := util.ReturnLocalIpAddress()
		for _, v := range ips {
			models.Rdb.SAdd(models.ServerLocalRunIpKey, v)
		}
	}

	return models.Rdb.SIsMember(models.ServerLocalRunIpKey, ip).Val()
}
