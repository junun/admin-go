package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/help"
	"api/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type ApproveResource struct {
	IsPass      int 		`form:"IsPass"`
	Reason 		string		`form:"Reason"`
}

type UndoConfirmResource struct {
	Version 	string 		`form:"Version"`
}

type DeployApp struct {
	ID          int         `form:"ID"`
	Name    	string    	`form:"Name"`
	Tid			int			`form:"Tid"`
	GitType		string      `form:"GitType"`
	TagBranch	string		`form:"TagBranch"`
	Commit 		string		`form:"Commit"`
	Desc 		string		`form:"Desc"`
	Status		int			`form:"Status"`
}

type AppTemplateDeploy struct {
	ID              int
	Aid             int
	Tid				int
	TemplateName    string
	GitType			string
	Name      		string
	TagBranch 		string
	Commit 			string
	Version 		string
	Reason 			string
	Desc 			string
	Status			int
	Operator		int
	Review          int
	Deploy          int
	UpdateTime 		time.Time
}

type DeployAppEnvRes struct {
	ID 				int
	Aid 			int
	Extend			int
	HostIds			string
	AppName 		string
	PreCode			string
	PreDeploy 		string
	EnableSync		int
	EnvId 			int
	EnvName 		string
}

type TargetRes struct {
	ID 		int
	Title   string
}

type LocalRes struct {
	Data  	[]string
}

type DeployRequestRes struct {
	AppName 		string
	EnvName 		string
	PreCode			string
	PreDeploy 		string
	Status 			int
	Targets			[]TargetRes
	Type			int
	Outputs			[]help.Msg
}

const (
	NewDeploy		= 1
	ReviewSuccess  	= 2
	OnDeploy		= 3
    UndoNeedDeploy  = 4
	ReviewFail 		= -1
	UndoFail		= -2
	DeployFail		= -3
	DeploySuccess 	= 5
	UndoSuccess     = 6
)

//func GetAppVersion(c *gin.Context)  {
//	if !middleware.PermissionCheckMiddleware(c,"config-app-git") {
//		util.JsonRespond(403, "请求资源被拒绝", "", c)
//		return
//	}
//	var det models.DeployExtend
//
//	models.DB.Model(&models.DeployExtend{}).
//		Where("dtid = ?", c.Param("id")).Find(&det)
//
//	if det.Dtid == 0 {
//		util.JsonRespond(500, "未找到指定发布模板", "", c)
//		return
//	}
//
//	res , e := util.FetchVersions(det.Aid, det.RepoUrl)
//	if e != nil {
//		util.JsonRespond(500, e.Error(), "", c)
//		return
//	}
//
//	data := make(map[string]interface{})
//	data["lists"] = res
//
//	util.JsonRespond(200, "", data, c)
//}

func GetGitTag(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"config-app-git") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}
	var det models.DeployExtend

	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", c.Param("id")).Find(&det)

	res , e := util.ReturnGitTagByCommand(det.Aid, det.RepoUrl)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

func GetGitBranch(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-git") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var det models.DeployExtend

	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", c.Param("id")).Find(&det)

	res , e := util.ReturnGitBranch(det.Aid, det.RepoUrl)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

func GetGitCommit(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-git") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var det models.DeployExtend

	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", c.Param("id")).Find(&det)

	branch := c.Param("branch")

	// 锁定项目，一个项目同时只能允许一个执行该方法
	key := models.GitAppOnWorking + det.TemplateName
	if models.GetValByKey(key) != "" {
		util.JsonRespond(500, "该项目别的用户在使用中，请稍后重试！", "", c)
		return
	}

	models.SetValByKey(key, "1",  2 * time.Second)

	res, e := util.GetGitLastTenCommitByBranch(det.Aid, det.RepoUrl, branch)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	models.Rdb.Del(key)

	data := make(map[string]interface{})
	data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用发布
// @Description 应用发布列表
// @Summary  应用发布列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app [get]
func GetAppDeploy(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"deploy-app-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	data := make(map[string]interface{})

	var deploy []AppTemplateDeploy
	// 分页逻辑还没有写，有空补上。
	e := models.DB.Raw("select d.*, e.aid, e.template_name from app_deploy d left join deploy_extend e on d.tid=e.dtid order by d.id DESC").Scan(&deploy).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	data["lists"] = deploy
	data["total"] = len(deploy)

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用发布
// @Description 应用发布提单
// @Summary 应用发布提单
// @Produce json
// @Param Authorization header string true "token"
// @Param Data body admin.DeployApp true "应用发布信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app [post]
func AddAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data DeployApp
	var deploy models.AppDeploy

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(406, "Invalid Add Deploy Data", "", c)
		return
	}

	maps := make(map[string]interface{})
	maps["tag_branch"] = data.TagBranch
	maps["tid"] = data.Tid

	// 发布唯一性检查
	if data.GitType == "tag" {
		if data.ID > 0  {
			models.DB.Model(&models.AppDeploy{}).
				Where(maps).
				Where("status <= 2").
				Where("id ! = ?", data.ID).
				Find(&deploy)
		} else {
			models.DB.Model(&models.AppDeploy{}).
				Where(maps).
				Where("status <= 2").
				Find(&deploy)
		}

	}

	if data.GitType == "branch" {
		if data.ID > 0  {
			models.DB.Model(&models.AppDeploy{}).
				Where(maps).
				Where("id != ?", data.ID).
				Where("commit = ?", data.Commit).
				Where("status <= 2").
				Find(&deploy)
		} else {
			models.DB.Model(&models.AppDeploy{}).
				Where(maps).
				Where("commit = ?", data.Commit).
				Where("status <= 2").
				Find(&deploy)
		}
	}

	if deploy.ID > 0 {
		util.JsonRespond(500, "重复的项目上线提单，请检查！", "", c)
		return
	}

	// 检查是否需要开启审核
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", data.Tid).
		Find(&det)

	sts := NewDeploy
	if  det.EnableCheck == 0  {
		sts = ReviewSuccess
	}

	uid,_ 	:= c.Get("Uid")
	deploy = models.AppDeploy{
		Name: data.Name,
		Tid: data.Tid,
		Desc: data.Desc,
		GitType: data.GitType,
		TagBranch: data.TagBranch,
		Commit: data.Commit,
		Operator: uid.(int),
		Status: sts,
		UpdateTime: time.Now().AddDate(0,0,0),
	}

	if data.ID > 0 {
		deploy.ID = data.ID
	}

	e := models.DB.Save(&deploy).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加项目提单申请成功", "", c)
}

// @Tags 应用发布
// @Description 应用发布修改
// @Summary  应用发布修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Param Data body admin.DeployApp true "应用发布信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app/{id} [put]
func PutAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id 	:= c.Param("id")
	var data DeployApp
	var deploy models.AppDeploy

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Deploy Data", "", c)
		return
	}

	// 发布唯一性检查
	models.DB.Model(&models.AppDeploy{}).
		Where("tid = ?", data.Tid).
		Where("repo_branch = ?", data.TagBranch).
		Where("repo_commit = ?", data.Commit).
		Where("status <= 2").
		Find(&deploy)

	if deploy.ID > 0 {
		util.JsonRespond(500, "重复的项目上线提单，请检查！", "", c)
		return
	}

	models.DB.Find(&deploy, id)

	uid,_ 				:= c.Get("Uid")
	deploy.Name     	= data.Name
	deploy.TagBranch	= data.TagBranch
	deploy.Commit 		= data.Commit
	deploy.Operator		= uid.(int)
	deploy.UpdateTime   = time.Now().AddDate(0,0,0)


	e := models.DB.Save(&deploy).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改成功", "", c)
}

// @Tags 应用发布
// @Description 应用发布删除
// @Summary  应用发布删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app/{id} [delete]
func DelAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id := c.Param("id")

	// 删除之前检查项目发布状态，只能删除status为1（新提单）， 2（审核成功）
	var deploy models.AppDeploy
	models.DB.Model(&models.HostApp{}).
		Where("id = ?", id).
		Where("status >= 3").
		Find(&deploy)

	if deploy.ID > 0 {
		util.JsonRespond(406, "改项目项目提单申请已经执行上线操作，无法删除！", "", c)
		return
	}

	err := models.DB.Delete(models.AppDeploy{}, "id = ?", id).Error
	if err != nil {
		util.JsonRespond(500, err.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除项目提单申请成功", "", c)
}


// @Tags 应用发布
// @Description 应用发布审核
// @Summary  应用发布审核
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app/{id}/review/{status} [put]
func PutAppDeployStatus(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-review") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data ApproveResource

	e 	:= c.BindJSON(&data)
	if e!= nil {
		util.JsonRespond(500, "Invalid Approve Data", "", c)
		return
	}

	id 			:= c.Param("id")
	uid,_ 		:= c.Get("Uid")
	status 		:= 2

	if data.IsPass == 0 {
		status = 3
	}

	if data.IsPass == 1 && data.Reason != "" || data.IsPass == 0 {
		e := models.DB.Model(&models.AppDeploy{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
			    "is_pass": data.IsPass,
				"status": status,
				"reason": data.Reason,
				"review": uid.(int) }).Error

		if e != nil {
			util.JsonRespond(500, e.Error(), "", c)
			return
		}
	} else {
		e := models.DB.Model(&models.AppDeploy{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"is_pass": data.IsPass,
			    "status": status,
				"review": uid.(int) }).Error

		if e != nil {
			util.JsonRespond(500, e.Error(), "", c)
			return
		}
	}

	if data.IsPass == 1 {
		util.JsonRespond(200, "审核通过", "", c)
		return
	}

	util.JsonRespond(200, "审核拒绝 :" + data.Reason, "", c)
}

// @Tags 应用发布
// @Description 应用发布查看
// @Summary  应用发布查看
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/request/{id} [get]
func GetDeployRequest(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-request") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var target TargetRes
	var targets []TargetRes

	data 	:= make(map[string]interface{})
	id 		:= c.Param("id")
	isLog 	:= c.Query("log")
	//isUndo  := c.Query("type")

	var deploy models.AppDeploy
	models.DB.Model(&models.HostApp{}).
		Where("id = ?", id).
		Find(&deploy)

	if deploy.ID == 0 {
		util.JsonRespond(500, "未找到指定发布申请", "", c)
		return
	}

	// 获取主机和app
	var res DeployAppEnvRes
	sql := "SELECT d.id, e.host_ids, e.pre_deploy, e.pre_code, e.extend, a.name as app_name, v.id as env_id, v.name as env_name FROM app_deploy d LEFT JOIN deploy_extend e ON d.tid = e.dtid LEFT JOIN app a  ON a.id = e.aid LEFT JOIN config_env v ON a.env_id = v.id WHERE a.active = 1 AND d.id = " + c.Param("id")

	e := models.DB.Raw(sql).Scan(&res).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	harr := strings.Split(res.HostIds, ",")
	var hosts []models.Host
	models.DB.Model(&models.Host{}).
		Where("id in (?)", harr).
		Find(&hosts)

	for _, v := range hosts {
		target.ID = v.ID
		target.Title = v.Name + "(" + v.Addres + string(v.Port) + ")"
		targets = append(targets, target)
	}

	var drr DeployRequestRes
	drr.AppName 	= res.AppName
	drr.EnvName 	= res.EnvName
	drr.Status		= deploy.Status
	drr.Targets 	= targets
	drr.Outputs		= []help.Msg{}

	if res.Extend == 2 {
		drr.PreCode 	= res.PreCode
		drr.PreDeploy	= res.PreDeploy
	}

	if isLog == "true" {
		msg := help.Msg{}
		var counter int64
		key 	:= models.DeployInfoKey + id
		counter = 0
		res, _ 	:= models.Rdb.LRange(key, counter, counter+9).Result()

		for {
			if len(res) > 0 {
				counter += 10
				for _, v := range res {
					if e := json.Unmarshal([]byte(v), &msg); e != nil {
						util.JsonRespond(500, e.Error(), "", c)
						return
					}
					drr.Outputs = append(drr.Outputs, msg)
				}

				res, _ 	= models.Rdb.LRange(key, counter, counter+9).Result()
				continue
			}
			break
		}
	}

    data["lists"] 	= drr
	util.JsonRespond(200, "", data, c)
}

// @Tags 应用发布
// @Description 应用发布请求
// @Summary  应用发布请求
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/request/{id} [post]
func PostDeployRequest(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-request") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id 	:= c.Param("id")
	var deploy models.AppDeploy
	models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Find(&deploy)

	if deploy.ID == 0 {
		util.JsonRespond(500, "未找到指定发布申请", "", c)
		return
	}

	if !(deploy.Status != ReviewSuccess || deploy.Status != UndoNeedDeploy) {
		util.JsonRespond(500, "该申请单当前状态还不能执行发布", "", c)
		return
	}

	// 获取主机和app
	sts := 0
	var res DeployAppEnvRes
	sql := "SELECT d.id, e.host_ids, a.id as aid, a.name as app_name, a.enable_sync, v.id as env_id, v.name as env_name FROM app_deploy d LEFT JOIN deploy_extend e ON d.tid = e.dtid LEFT JOIN app a ON a.id = e.aid LEFT JOIN config_env v ON a.env_id = v.id WHERE a.active = 1 AND d.id = " + c.Param("id")

	e := models.DB.Raw(sql).Scan(&res).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	uid,_ 	:= c.Get("Uid")
	var user models.User
	models.DB.Model(&models.User{}).
		Where("id = ?", uid).
		Find(&user)

	if user.IsSupper != 1 {
		if !middleware.RoleAppAuthMiddleware(user.Rid, res.EnvId, res.Aid, c) {
			return
		}
	}

	// 检查项目是否需要初始化
	if res.EnableSync == 1 {
		sts = 1
	}

	var hostapp []models.HostApp
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", res.Aid).
		Where("status = ?", sts).
		Find(&hostapp)


	harr := strings.Split(res.HostIds, ",")
	var hosts []models.Host
	models.DB.Model(&models.Host{}).
		Where("id in (?)", harr).
		Find(&hosts)

	if len(hostapp) <= 0 && sts == 1 {
		util.JsonRespond(500, "该项目还没有初始化，请先初始化！", "", c)
		return
	}

	if len(hostapp) <= 0 && sts == 0 {
		util.JsonRespond(500, "该项目还没有绑定到主机，请先绑定业务到主机！", "", c)
		return
	}

	data := make(map[string]interface{})
	for _, v := range hosts {
		var localData LocalRes
		localData.Data = []string{""}
		strid := strconv.Itoa(v.ID)
		data[strid] = localData
	}

	var localData LocalRes
	localData.Data = []string{util.HumanNowTime() + "建立接连...   "}

	data["local"] = localData
	util.JsonRespond(200, "", data, c)
}


// @Tags 应用发布
// @Description 应用回滚请求
// @Summary  应用回滚请求
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/undo/request/{id} [get]
func GetUndoRequest(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"undo-app-request") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	data 	:= make(map[string]interface{})
	id 		:= c.Param("id")

	var deploy models.AppDeploy
	models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Find(&deploy)

	if deploy.ID == 0 {
		util.JsonRespond(500, "未找到指定发布记录", "", c)
		return
	}

	// 检查发布申请时间,1：超过1天禁止回滚 2： 有更新发布的版本也禁止回滚
	UpdateTimeAddOneDay 	:= deploy.UpdateTime.AddDate(0, 0, 1)
	if (deploy.UpdateTime.After(UpdateTimeAddOneDay)) {
		util.JsonRespond(500, "该项目发布已经超过1天禁止回滚!", "", c)
		return
	}

	var newdeploy []models.AppDeploy
	models.DB.Model(&models.AppDeploy{}).
		Where("id > ?", id).
		Where("tid = ?", deploy.Tid).
		Where("status = ? or status = ? or status = ?", DeploySuccess, UndoSuccess, UndoNeedDeploy).
		Find(&newdeploy)

	if len(newdeploy) > 0 {
		util.JsonRespond(500, "该项目有更新的版本发布，禁止回滚该版本!", "", c)
		return
	}

	var olddeploy models.AppDeploy
	models.DB.Model(&models.AppDeploy{}).
		Where("id < ?", id).
		Where("tid = ?", deploy.Tid).
		Where("status = ? or status = ? ", DeploySuccess, UndoSuccess).
		Find(&olddeploy)

	if olddeploy.ID == 0 || olddeploy.Version == ""{
		util.JsonRespond(500, "未找到该应用可以用于回滚的版本", "", c)
		return
	}

	data["UpdateTime"] 	= olddeploy.UpdateTime
	data["Version"] 	= olddeploy.Version

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用发布
// @Description 应用回滚创建
// @Summary  应用回滚创建
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/undo/confirm/{id} [put]
func PutUndoRequest(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"undo-app-request") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data UndoConfirmResource

	id 	:= c.Param("id")
	if e 	:= c.BindJSON(&data); e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	var deploy models.AppDeploy
	var newDeploy models.AppDeploy
	models.DB.Model(&models.HostApp{}).
		Where("id = ?", id).
		Find(&deploy)

	if deploy.ID == 0 {
		util.JsonRespond(500, "未找到指定发布记录", "", c)
		return
	}

	newDeploy.Version 	= data.Version
	newDeploy.Status	= UndoNeedDeploy
	newDeploy.Tid		= deploy.Tid
	newDeploy.GitType   = deploy.GitType
	newDeploy.Name		= deploy.Name + "-回滚"
	newDeploy.IsPass	= deploy.IsPass
	newDeploy.UpdateTime= time.Now()

	if e := models.DB.Save(&newDeploy).Error; e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "回滚申请创建成功", "", c)
}