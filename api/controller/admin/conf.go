package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"strings"
)

type EnvResource struct {
	Name    	string    	`form:"Name"`
	Desc 		string    	`form:"Desc"`
}

type AppTypeResource struct {
	Name    	string    	`form:"Name"`
	Desc 		string    	`form:"Desc"`
}

type AppResource struct {
	Name    	string    	`form:"Name"`
	Active		int			`form:"Active"`
	Tid			int			`form:"Tid"`
	EnvId		int         `form:"EnvId"`
	EnableSync	int			`form:"EnableSync"`
	DeployType	int			`form:"DeployType"`
	Desc 		string    	`form:"Desc"`
}

type DeployExtendResource struct {
	Aid				int      `form:"Aid"`
	TemplateName    string   `form:"TemplateName"`
	EnableCheck		int      `form:"EnableCheck"`
	HostIds         string   `form:"HostIds"`
	RepoUrl         string   `form:"RepoUrl"`
	Tag             string   `form:"Tag"`
	Versions        int      `form:"Versions"`
	PreCode  		string   `form:"PreCode"`
	PostCode  		string   `form:"PostCode"`
	PreDeploy       string   `form:"PreDeploy"`
	PostDeploy      string   `form:"PostDeploy"`
	DstDir      	string   `form:"DstDir"`
	DstRepo      	string   `form:"DstRepo"`
}

type AppValueResource struct {
	Aid			int			`form:"Aid"`
	Name    	string    	`form:"Name"`
	Value    	string    	`form:"Value"`
	Desc 		string    	`form:"Desc"`
}

// @Tags 应用配置
// @Description 应用环境列表
// @Summary  应用环境列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/env [get]
func GetConfigEnv(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"config-env-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var roles []models.ConfigEnv
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.ConfigEnv{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&roles)
	models.DB.Model(&models.ConfigEnv{}).Where(maps).Count(&count)

	data["lists"] = roles
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用配置
// @Description 应用环境新增
// @Summary  应用环境新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.HostRoleResource true "应用环境信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/env [post]
func AddConfigEnv(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-env-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data EnvResource
	var role models.ConfigEnv

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add HostRole Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.ConfigEnv{}).
		Where("name = ?", data.Name).Find(&role)

	if role.ID > 0 {
		util.JsonRespond(500, "重复的角色名，请检查！", "", c)
		return
	}

	role = models.ConfigEnv{
		Name: data.Name,
		Desc: data.Desc}

	e := models.DB.Save(&role).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加主机类型成功", "", c)
}

// @Tags 应用配置
// @Description 应用环境修改
// @Summary  应用环境修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用环境ID"
// @Param Data body admin.HostRoleResource true "应用环境信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/env/{id} [put]
func PutConfigEnv(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-env-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data EnvResource
	var role models.ConfigEnv

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Role Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 角色名唯一性检查
	models.DB.Model(&models.ConfigEnv{}).
		Where("name = ?", data.Name).
		Where("id != ?", c.Param("id")).Find(&role)

	if role.ID > 0 {
		util.JsonRespond(500, "重复的主机类型，请检查！", "", c)
		return
	}

	models.DB.Find(&role, c.Param("id"))

	role.Name = data.Name
	role.Desc = data.Desc

	e := models.DB.Save(&role).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改主机类型成功", "", c)
}

// @Tags 应用配置
// @Description 应用环境删除
// @Summary  应用环境删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "主机ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/env/{id} [delete]
func DelConfigEnv(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-env-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	e := models.DB.Delete(models.ConfigEnv{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除应用类型成功", "", c)
}

// @Tags 应用配置
// @Description 应用列表
// @Summary  应用列表
// @Produce  json
// @Param active query string false "active 查询项目是否生效"
// @Param envId query string false "环境类型id"
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/app [get]
func GetConfigApp(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"config-app-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var app []models.App
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	IsActive,_ :=  com.StrTo(c.Query("active")).Int()
	if IsActive != 0 {
		maps["active"] = IsActive
	}

	EnvId,_ :=  com.StrTo(c.Query("envId")).Int()
	if EnvId != 0 {
		maps["env_id"] = EnvId
	}

	models.DB.Model(&models.App{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&app)
	models.DB.Model(&models.App{}).Where(maps).Count(&count)

	data["lists"] = app
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用配置
// @Description 应用新增
// @Summary  应用新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.AppResource true "应用信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/app [post]
func AddConfigApp(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data AppResource
	var app models.App

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(406, "Invalid Add App Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 项目名唯一性检查
	models.DB.Model(&models.App{}).
		Where("name = ?", data.Name).
		Where("env_id = ?", data.EnvId).
		Find(&app)

	if app.ID > 0 {
		util.JsonRespond(500, "重复的项目名，请检查！", "", c)
		return
	}

	app = models.App{
		Name: data.Name,
		Tid: data.Tid,
		EnvId: data.EnvId,
		Active: data.Active,
		DeployType: data.DeployType,
		EnableSync: data.EnableSync,
		Desc: data.Desc,
	}

	e := models.DB.Save(&app).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加项目成功", "", c)
}

// @Tags 应用配置
// @Description 应用修改
// @Summary  应用修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用ID"
// @Param Data body admin.AppResource true "应用信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/app/{id} [put]
func PutConfigApp(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data AppResource
	var app models.App

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit App Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 角色名唯一性检查
	models.DB.Model(&models.App{}).
		Where("name = ?", data.Name).
		Where("env_id = ?", data.EnvId).
		Where("id != ?", c.Param("id")).Find(&app)

	if app.ID > 0 {
		util.JsonRespond(500, "重复的应用名，请检查！", "", c)
		return
	}

	models.DB.Find(&app, c.Param("id"))

	app.Name     	= data.Name
	app.Tid			= data.Tid
	app.EnvId		= data.EnvId
	app.Active      = data.Active
	app.Desc  		= data.Desc
	app.EnableSync	= data.EnableSync
	app.DeployType	= data.DeployType

	e := models.DB.Save(&app).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改App成功", "", c)
}

// @Tags 应用配置
// @Description 应用删除
// @Summary  应用删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/app/{id} [delete]
func DelConfigApp(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id := c.Param("id")

	// 删除之前检查是否有主机绑定
	var hostapp models.HostApp
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", id).Find(&hostapp)

	if hostapp.ID > 0 {
		util.JsonRespond(500, "改项目仍被主机绑定，请先取消！", "", c)
		return
	}

	// 删除之前检查是否有发布模板绑定
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("aid = ?", id).Find(&det)

	if det.Dtid > 0 {
		util.JsonRespond(500, "该项目仍被发布绑定，请先取消！", "", c)
		return
	}

	err := models.DB.Delete(models.App{}, "id = ?", id).Error
	if err != nil {
		util.JsonRespond(500, err.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除应用成功", "", c)
}

// @Tags 应用配置
// @Description 应用变量列表
// @Summary  应用变量列表
// @Produce  json
// @Param Authorization header string true "token"
// @Param aid query string true "aid"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/value [get]
func GetAppValue(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"app-value-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var value []models.AppSyncValue
	var count int
	maps 	:= make(map[string]interface{})
	data 	:= make(map[string]interface{})
	aid, _ 	:= com.StrTo(c.Query("Aid")).Int()

	maps["aid"] = aid

	models.DB.Model(&models.AppSyncValue{}).Where(maps).
		Offset(util.GetPage(c)).
		Limit(util.GetPageSize(c)).
		Find(&value)

	models.DB.Model(&models.AppSyncValue{}).Where(maps).Count(&count)

	data["lists"] = value
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用配置
// @Description 应用变量新增
// @Summary  应用变量新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.AppValueResource true "应用变量信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/value [post]
func AddAppValue(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"app-value-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data AppValueResource
	var value models.AppSyncValue

	e := c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(406, e.Error(), "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 项目变量名唯一性检查
	models.DB.Model(&models.AppSyncValue{}).
		Where("name = ?", data.Name).
		Where("aid = ?", data.Aid).
		Find(&value)

	if value.ID > 0 {
		util.JsonRespond(500, "重复的应用变量名，请检查！", "", c)
		return
	}

	value = models.AppSyncValue{
		Aid: data.Aid,
		Name: data.Name,
		Value: data.Value,
		Desc: data.Desc,
	}

	e = models.DB.Save(&value).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加应用变量成功", "", c)
}

// @Tags 应用配置
// @Description 应用变量修改
// @Summary  应用变量修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用变量ID"
// @Param Data body admin.AppValueResource true "应用变量信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/value/{id} [put]
func PutAppValue(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"app-value-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data AppValueResource
	var value models.AppSyncValue

	e := c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)

	// 应用变量名唯一性检查
	models.DB.Model(&models.AppSyncValue{}).
		Where("name = ?", data.Name).
		Where("aid = ?", data.Aid).
		Where("id != ?", c.Param("id")).
		Find(&value)

	if value.ID > 0 {
		util.JsonRespond(500, "重复的应用变量名，请检查！", "", c)
		return
	}

	models.DB.Find(&value, c.Param("id"))

	value.Name  = data.Name
	value.Aid	= data.Aid
	value.Desc 	= data.Desc

	e 	= models.DB.Save(&value).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改应用变量成功", "", c)
}

// @Tags 应用配置
// @Description 应用变量删除
// @Summary  应用变量删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用变量ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/value/{id} [delete]
func DelAppValue(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"app-value-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id := c.Param("id")

	// 删除之前应该有检查逻辑
	// 待补充

	e := models.DB.Delete(models.AppSyncValue{}, "id = ?", id).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除应用变量成功", "", c)
}

// @Tags 应用配置
// @Description 应用类型列表
// @Summary  应用类型列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/type [get]
func GetAppType(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"app-type-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var atype []models.AppType
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.AppType{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&atype)
	models.DB.Model(&models.AppType{}).Where(maps).Count(&count)

	data["lists"] = atype
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用配置
// @Description 应用类型新增
// @Summary  应用类型新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.AppTypeResource true "应用类型信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/type [post]
func AddAppType(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"app-type-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data AppTypeResource
	var atype models.AppType

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add AppType Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.AppType{}).
		Where("name = ?", data.Name).Find(&atype)

	if atype.ID > 0 {
		util.JsonRespond(500, "重复的应用类型名，请检查！", "", c)
		return
	}

	atype = models.AppType{
		Name: data.Name,
		Desc: data.Desc}

	e := models.DB.Save(&atype).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加应用类型成功", "", c)
}

// @Tags 应用配置
// @Description 应用类型修改
// @Summary  应用类型修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用类型ID"
// @Param Data body admin.AppTypeResource true "应用类型信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/type/{id} [put]
func PutAppType(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"app-type-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data EnvResource
	var atype models.AppType

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit AppType Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 角色名唯一性检查
	models.DB.Model(&models.AppType{}).
		Where("name = ?", data.Name).
		Where("id != ?", c.Param("id")).Find(&atype)

	if atype.ID > 0 {
		util.JsonRespond(500, "重复的应用类型，请检查！", "", c)
		return
	}

	models.DB.Find(&atype, c.Param("id"))

	atype.Name = data.Name
	atype.Desc = data.Desc

	e := models.DB.Save(&atype).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改应用类型成功", "", c)
}

// @Tags 应用配置
// @Description 应用类型删除
// @Summary  应用类型删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用类型ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/type/{id} [delete]
func DelAppType(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"app-type-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	e := models.DB.Delete(models.AppType{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除应用类型成功", "", c)
}

// @Tags 应用配置
// @Description 应用模板列表
// @Summary  应用模板列表
// @Produce  json
// @Param Authorization header string true "token"
// @Param aid query string true "通过app id查询"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/deploy [get]
func GetDeployExtend(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	var det []models.DeployExtend
	aid, _ 	:= com.StrTo(c.Query("Aid")).Int()

	if aid > 0 {
		maps["aid"] = aid
	}

	models.DB.Model(&models.DeployExtend{}).
		Where(maps).Find(&det)

	data["lists"] = det

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用配置
// @Description 应用发布模板新增
// @Summary  应用发布模板新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.HostRoleResource true "应用发布模板信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/deploy [post]
func AddDeployExtend(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data DeployExtendResource
	var det models.DeployExtend

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add DeployExtend Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.DeployExtend{}).
		Where("template_name = ?", data.TemplateName).Find(&det)

	if det.Dtid > 0 {
		util.JsonRespond(500, "重复的发布模板名，请检查！", "", c)
		return
	}

	det = models.DeployExtend{
		Aid: data.Aid,
		TemplateName: data.TemplateName,
		EnableCheck: data.EnableCheck,
		HostIds: data.HostIds,
		RepoUrl: data.RepoUrl,
		Versions: data.Versions,
		PreCode: data.PreCode,
		PostCode: data.PostCode,
		PreDeploy: data.PreDeploy,
		PostDeploy: data.PostDeploy,
		DstDir: strings.TrimSpace(data.DstDir),
		DstRepo: strings.TrimSpace(data.DstRepo),
	}

	if data.Tag !=  "" {
		det.Tag = data.Tag
	}

	e := models.DB.Save(&det).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加应用发布模板成功", "", c)
}

// @Tags 应用配置
// @Description 应用发布模板修改
// @Summary  应用发布模板修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布模板修改ID"
// @Param Data body admin.AppTypeResource true "应用类型信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/deploy/{id} [put]
func PutDeployExtend(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"config-app-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data DeployExtendResource
	var det models.DeployExtend

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit DeployExtend Data", "", c)
		return
	}

	data.TemplateName = strings.TrimSpace(data.TemplateName)

	// 角色名唯一性检查
	models.DB.Model(&models.AppType{}).
		Where("template_name = ?", data.TemplateName).
		Where("dtid != ?", c.Param("id")).Find(&det)

	if det.Dtid > 0 {
		util.JsonRespond(500, "重复的发布模板名，请检查！", "", c)
		return
	}

	e := models.DB.Exec("update deploy_extend set aid = ?, template_name = ?," +
		"enable_check = ?, tag = ?, host_ids = ? , repo_url = ?," +
		"versions = ?, pre_code = ?, post_code = ?, pre_deploy = ? , " +
		"post_deploy = ?, dst_dir = ?, dst_repo = ? where dtid = ? ",
		data.Aid, data.TemplateName, data.EnableCheck, data.Tag, data.HostIds, data.RepoUrl,
		data.Versions, data.PreCode, data.PostCode, data.PreDeploy, data.PostDeploy,
		strings.TrimSpace(data.DstDir),strings.TrimSpace(data.DstRepo), c.Param("id")).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改发布模板成功", "", c)
}

// @Tags 应用配置
// @Description 应用发布模板删除
// @Summary  应用发布模板删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "发布模板id"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/deploy/{id} [delete]
func DelDeployExtend(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"config-app-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	e := models.DB.Delete(models.DeployExtend{}, "dtid = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除应用发布模板成功", "", c)
}

// @Tags 应用配置
// @Description 应用模板列表
// @Summary  应用模板列表
// @Produce  json
// @Param aid query string true "项目id"
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/config/template [get]
func GetAppTemplate(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"config-app-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var det []models.DeployExtend

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	aid,_ :=  com.StrTo(c.Query("aid")).Int()
	if aid != 0 {
		maps["aid"] = aid
	}

	models.DB.Model(&models.DeployExtend{}).Where(maps).Find(&det)

	data["lists"] = det


	util.JsonRespond(200, "", data, c)
}
