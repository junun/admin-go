package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/setting"
	"api/pkg/util"
	"fmt"
	"github.com/unknwon/com"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HostRoleResource struct {
	Name    	string    	`form:"Name"`
	Desc 		string    	`form:"Desc"`
}


type HostResource struct {
	Rid			int			`form:"Rid"`
	EnvId		int 		`form:"EnvId"`
	ZoneId		int			`form:"ZoneId"`
	Name    	string    	`form:"Name"`
	Addres    	string    	`form:"Addres"`
	Port    	int    		`form:"Port"`
	Status    	int    		`form:"Status"`
	Username    string    	`form:"Username"`
	Desc 		string    	`form:"Desc"`
	Password    string      `form:"Password"`
}

type HostAppResource struct {
	Hid			int			`form:"Hid"`
	Aid			int			`form:"Aid"`
	Desc 		string    	`form:"Desc"`
}

// @Tags 主机管理
// @Description 主机类型列表
// @Summary  主机类型列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/role [get]
func GetHostRole(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"host-role-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var roles []models.HostRole
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.HostRole{}).Where(maps).
		Offset(util.GetPage(c)).Limit(setting.AppSetting.PageSize).Find(&roles)
	models.DB.Model(&models.HostRole{}).Where(maps).Count(&count)

	data["lists"] = roles
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 主机管理
// @Description 主机类型添加
// @Summary  主机类型添加
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.HostRoleResource true "主机类型信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/role [post]
func AddHostRole(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-role-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data HostRoleResource
	var role models.HostRole

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add HostRole Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.HostRole{}).
		Where("name = ?", data.Name).Find(&role)

	if role.ID > 0 {
		util.JsonRespond(500, "重复的角色名，请检查！", "", c)
		return
	}

	role = models.HostRole{
		Name: data.Name,
		Desc: data.Desc}

	e := models.DB.Save(&role).Error

	if e != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	util.JsonRespond(200, "添加主机类型成功", "", c)
}

// @Tags 主机管理
// @Description 主机类型修改
// @Summary  主机类型修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "环境ID"
// @Param Data body admin.HostRoleResource true "主机类型信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/role/{id} [put]
func PutHostRole(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-role-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data HostRoleResource
	var role models.HostRole

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Role Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.HostRole{}).
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
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	util.JsonRespond(200, "修改主机类型成功", "", c)
}

// @Tags 主机管理
// @Description 主机类型删除
// @Summary  主机类型删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "环境ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/role/{id} [delete]
func DelHostRole(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-role-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	err := models.DB.Delete(models.HostRole{}, "id = ?", c.Param("id")).Error
	if err != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	util.JsonRespond(200, "删除主机类型成功", "", c)
}

// @Tags 主机管理
// @Description 主机列表
// @Summary  主机列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host [get]
func GetHost(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"host-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var hosts []models.Host
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})


	rid, _ 		:= com.StrTo(c.Query("Rid")).Int()
	status, _ 	:= com.StrTo(c.Query("Status")).Int()
	name 		:= c.Query("Name")

	if rid != 0 {
		maps["Rid"] = rid
	}

	if status != 0 {
		maps["Status"] = status
	}

	if name != "" {
		maps["Name"] = name
	}

	models.DB.Model(&models.Host{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&hosts)

	models.DB.Model(&models.Host{}).Where(maps).Count(&count)

	// 定时任务增加id 为0 本机选项
	if  c.Query("Source") == "schedule" {

		hostLocal := models.Host{
			Name: "本机",
			Addres: "127.0.0.1",
			Port: 22,
		}
		hosts = append(hosts, hostLocal);
		count++
	}


	data["lists"] = hosts
	data["total"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 主机管理
// @Description 主机添加
// @Summary  主机添加
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.HostResource true "主机信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host [post]
func AddHost(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data HostResource
	var host models.Host

	err 	:= c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add Host Data", "", c)
		return
	}

	// 验证主机信息并 copy 公钥
	if !util.ValidHosh(data.Addres, data.Port, data.Username, data.Password ) {
		util.JsonRespond(500, "Auth Fail", "", c)
		return
	}

	uid,_ 	:= c.Get("Uid")
	host = models.Host{
		Name: data.Name,
		Rid: data.Rid,
		EnvId: data.EnvId,
		ZoneId: data.ZoneId,
		Username: data.Username,
		Addres: data.Addres,
		Port: data.Port,
		Operator: uid.(int),
		Status: data.Status,
		Desc: data.Desc,
	}

	e := models.DB.Save(&host).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加主机成功", "", c)
}

// @Tags 主机管理
// @Description 主机修改
// @Summary  主机修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "主机ID"
// @Param Data body admin.HostRoleResource true "主机信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/hosts/{id} [put]
func PutHost(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data HostResource
	var host models.Host

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Host Data", "", c)
		return
	}

	models.DB.Find(&host, c.Param("id"))

	uid,_ 	:= c.Get("Uid")

	host.Name = data.Name
	host.Rid = data.Rid
	host.EnvId = data.EnvId
	host.ZoneId = data.ZoneId
	host.Username = data.Username
	host.Addres = data.Addres
	host.Port = data.Port
	host.Operator = uid.(int)
	host.Status = data.Status
	host.Desc = data.Desc


	e := models.DB.Save(&host).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改主机信息成功", "", c)
}

// @Tags 主机管理
// @Description 主机删除
// @Summary  主机删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "主机ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/hosts/{id} [delete]
func DelHost(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	// 检查主机是否有业务绑定
	var app models.HostApp
	models.DB.Model(&models.HostApp{}).
		Where("hid = ?", c.Param("id")).Find(&app)

	if app.ID > 0 {
		util.JsonRespond(500, "该主机有业务绑定，请先取消！", "", c)
		return
	}

	e := models.DB.Delete(models.Host{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除主机成功", "", c)
}

func ConsoleHost(c *gin.Context)  {
	token := c.Query("x-token")

	if token == "" {
		util.JsonRespond(401, "API token required", "", c)
		c.Abort()
		return
	}

	var user  models.User
	e := models.DB.Where("access_token = ?", token).First(&user).Error
	if e != nil{
		util.JsonRespond(401, "Invalid API token, please login", "", c)
		c.Abort()
	}

	if user.IsSupper != 1 {
		key := models.RoleRermsSetKey
		UserRid := user.Rid
		str := fmt.Sprintf("%v", UserRid)
		key =  key + str

		// 检查 redis 有没有该key的集合
		err := models.Rdb.Exists(key).Val()
		if err != 1 {
			rid 	:= UserRid
			models.SetRolePermToSet(key, rid)
		}

		// 检查对应的set是否有该角色权限
		if models.CheckMemberByKey(key, "host_console") {
			util.JsonRespond(403, "请求资源被拒绝", "", c)
			return
		}
	}

	var host models.Host
	models.DB.Find(&host, c.Param("id"))

	if host.ID <= 0 {
		util.JsonRespond(500, "Unknown Host！", "", c)
		return
	}

	c.HTML(http.StatusOK, "web_ssh.html", gin.H{
		"title": host.Name ,
		"id": c.Param("id"),
		"token": token,
	})
}

// @Tags 主机管理
// @Description 主机业务列表
// @Summary  主机业务列表
// @Produce  json
// @Param Authorization header string true "token"
// @Param hid query string false "通过host id 查询"
// @Param aid query string false "通过app id 查询"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/app [get]
func GetHostApp(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"host-app-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var apps []models.HostApp
	var count int
	maps 	:= make(map[string]interface{})
	data 	:= make(map[string]interface{})
	hid, _ 	:= com.StrTo(c.Query("Hid")).Int()
	aid, _ 	:= com.StrTo(c.Query("Aid")).Int()

	if hid > 0 {
		maps["hid"] = hid
	}

	if aid > 0 {
		maps["aid"] = aid
	}

	models.DB.Model(&models.HostApp{}).Where(maps).
		Offset(util.GetPage(c)).Limit(setting.AppSetting.PageSize).Find(&apps)
	models.DB.Model(&models.HostApp{}).Where(maps).Count(&count)

	data["lists"] = apps
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 主机管理
// @Description 主机业务绑定
// @Summary  主机业务绑定
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.HostAppResource true "主机业务信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/app [post]
func AddHostApp(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-app-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data HostAppResource
	var app models.HostApp

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(400, "Invalid Add HostApp Data", "", c)
		return
	}

	// 检查是否已经绑定了
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", data.Aid).
		Where("hid = ?", data.Hid).
		Find(&app)

	if app.ID > 0 {
		util.JsonRespond(500, "重复的业务绑定，请检查！", "", c)
		return
	}

	// 检查主机环境和项目环境是否一致
	var project models.App
	var host models.Host

	models.DB.Model(&models.App{}).
		Where("id = ?", data.Aid).
		Find(&project)

	models.DB.Model(&models.Host{}).
		Where("id = ?", data.Hid).
		Find(&host)

	if host.EnvId != project.EnvId {
		util.JsonRespond(500, "业务和主机环境不一致，无法绑定业务！", "", c)
		return
	}

	app = models.HostApp{
		Hid: data.Hid,
		Aid: data.Aid,
		Desc: data.Desc}

	e := models.DB.Save(&app).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "绑定主机业务成功", "", c)
}

// @Tags 主机管理
// @Description 主机业务修改
// @Summary  主机业务修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "主机业务ID"
// @Param Data body admin.HostRoleResource true "主机业务信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/app/{id} [put]
func PutHostApp(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-app-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data HostAppResource
	var app models.HostApp

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(400, "Invalid Edit HostApp Data", "", c)
		return
	}

	// 检查是否已经绑定了
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", data.Aid).
		Where("hid = ?", data.Hid).
		Where("id != ?", c.Param("id")).Find(&app)

	if app.ID > 0 {
		util.JsonRespond(500, "重复的主机业务绑定，请检查！", "", c)
		return
	}

	// 检查主机环境和项目环境是否一致
	var project models.App
	var host models.Host

	models.DB.Model(&models.App{}).
		Where("id = ?", data.Aid).
		Find(&project)

	models.DB.Model(&models.Host{}).
		Where("id = ?", data.Hid).
		Find(&host)

	if host.EnvId != project.EnvId {
		util.JsonRespond(500, "业务和主机环境不一致，无法绑定业务！", "", c)
		return
	}

	models.DB.Find(&app, c.Param("id"))
	app.Hid		= data.Hid
	app.Aid  	= data.Aid
	app.Desc 	= data.Desc

	e := models.DB.Save(&app).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改绑定业务成功", "", c)
}

// @Tags 主机管理
// @Description 主机业务删除
// @Summary  主机业务删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "主机业务ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/app/{id} [delete]
func DelHostApp(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"host-app-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	err := models.DB.Delete(models.HostApp{}, "id = ?", c.Param("id")).Error
	if err != nil {
		util.JsonRespond(500, err.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除主机业务成功", "", c)
}


// @Tags 主机管理
// @Description 通过业务查询主机列表
// @Summary  通过业务查询主机列表
// @Produce  json
// @Param Authorization header string true "token"
// @Param aid query string true "通过app id 查询"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/host/appid [get]
func GetHostByAppId(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"host-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var host []models.Host

	data 	:= make(map[string]interface{})
	aid, _ 	:= com.StrTo(c.Query("Aid")).Int()

	models.DB.Table("host").
		Joins("left join host_app on host.id = host_app.hid").
		Where("host_app.aid = ?", aid).
		Find(&host)


	data["lists"] = host
	data["count"] = len(host)

	util.JsonRespond(200, "", data, c)
}
