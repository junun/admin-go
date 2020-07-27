package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/setting"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"runtime"
	"strconv"
	"strings"
)

type SettingResource struct {
	Data 		[]models.Settings   `form:"Name"`
}

type RobotResource struct {
	Name    	string    	`form:"Name"`
	Type		int			`form:"Type"`
	Status		int			`form:"Status"`
	Webhook		string    	`form:"Webhook"`
	Secret		string		`form:"Secret"`
	Keyword		string		`form:"Keyword"`
	Desc 		string    	`form:"Desc"`
}

type MailResource struct {
	Server 		string    	`form:"server"`
	Port		string		`form:"port"`
	Username	string  	`form:"username"`
	Password	string		`form:"password"`
	Nickname    string		`form:"nickname"`
}

func GetSetting(c *gin.Context)  {
	var syssettings []models.Settings
	data := make(map[string]interface{})

	models.DB.Model(&models.Settings{}).
		Find(&syssettings)

	data["lists"] = syssettings

	util.JsonRespond(200, "", data, c)
}

// @Tags 系统管理
// @Description 系统设置
// @Summary  系统设置
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.SettingResource true "系统设置信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system [post]
func SettingModify(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"setting-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data SettingResource
	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(400, "Invalid Setting Data", "", c)
		return
	}

	for _, v := range data.Data {
		var set models.Settings
		models.DB.Model(&models.Settings{}).
			Where("name = ?", v.Name).
			Find(&set)

		if set.ID > 0 {
			desc := set.Desc
			if v.Desc != "" && v.Desc != desc {
				desc 	= v.Desc
			}
			e := models.DB.Model(&set).
				Updates(map[string]interface{}{"value": v.Value, "desc": desc}).
				Error

			if e != nil {
				util.JsonRespond(500, e.Error(), "", c)
				return
			}
		} else {
			e := models.DB.Save(&v).Error
			if e != nil {
				util.JsonRespond(500, e.Error(), "", c)
				return
			}
		}
	}

	util.JsonRespond(200, "操作成功", "", c)
}


func LdapTest()  {

}

// @Tags 系统管理
// @Description 邮件测试
// @Summary  邮件测试
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body MailResource true ""
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system/mail [post]
func EmailTest(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"setting-email-test") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data MailResource
	e := c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	port, _ := strconv.Atoi(data.Port)
	gd	:= models.InitDialer(data.Server, data.Username, data.Password, port)
	msg := "This is a test email！"
	m 	:= models.CreateMsg(data.Username, []string{data.Username}, msg, msg)
	if e = gd.DialAndSend(m); e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "邮件测试正常", "", c)
}

func About(c *gin.Context)  {
	data := make(map[string]interface{})
	var about models.About
	about.Golangversion = runtime.Version()
	about.SystemInfo 	= runtime.GOOS
	about.GinVersion 	= gin.Version
	data["lists"] 		= about

	util.JsonRespond(200, "", data, c)
}

// @Tags 系统管理
// @Description 机器人通道列表
// @Summary  机器人通道列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system/robot [get]
func GetRobot(c *gin.Context)  {
	var robots []models.SettingRobot
	var count int

	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	Status,_ :=  com.StrTo(c.Query("status")).Int()
	if Status != 0 {
		maps["status"] = Status
	}

	models.DB.Model(&models.SettingRobot{}).
		Where(maps).
		Offset(util.GetPage(c)).
		Limit(setting.AppSetting.PageSize).
		Find(&robots)

	models.DB.Model(&models.SettingRobot{}).Where(maps).Count(&count)

	data["lists"] = robots
	data["count"] = count
	util.JsonRespond(200, "", data, c)
}

// @Tags 系统管理
// @Description 机器人通道添加
// @Summary  机器人通道添加
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.RobotResource true "机器人通道信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system/robot [post]
func AddRobot(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"setting-robot-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data RobotResource
	var robot models.SettingRobot

	e := c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(406, "Invalid Add Robot Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 项目名唯一性检查
	models.DB.Model(&models.SettingRobot{}).
		Where("name = ?", data.Name).
		Find(&robot)

	if robot.ID > 0 {
		util.JsonRespond(500, "重复的通道标识名，请检查！", "", c)
		return
	}

	// 如果是 数字签名，必须提供 Secret
	if data.Type == 1 && data.Secret == "" {
		util.JsonRespond(500, "钉钉数字签名方式必须提供 Secret", "", c)
		return
	}

	//钉钉关键字方式必须提供 Keyword
	if data.Type == 2 && data.Keyword == "" {
		util.JsonRespond(500, "钉钉关键字方式必须提供 Keyword", "", c)
		return
	}

	robot = models.SettingRobot{
		Name: data.Name,
		Webhook: data.Webhook,
		Secret: data.Secret,
		Keyword: data.Keyword,
		Status: data.Status,
		Type: data.Type,
		Desc: data.Desc,
	}

	e 	= models.DB.Save(&robot).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加机器人通道成功", "", c)
}

// @Tags 系统管理
// @Description 机器人通道修改
// @Summary  机器人通道修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "通道ID"
// @Param Data body admin.RobotResource true "机器人通道信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system/robot/{id} [put]
func PutRobot(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"setting-robot-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data RobotResource
	var robot models.SettingRobot

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Robot Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 角色名唯一性检查
	models.DB.Model(&models.App{}).
		Where("name = ?", data.Name).
		Where("id != ?", c.Param("id")).Find(&robot)

	if robot.ID > 0 {
		util.JsonRespond(500, "重复的机器人通道标识，请检查！", "", c)
		return
	}

	models.DB.Find(&robot, c.Param("id"))

	robot.Name     	= data.Name
	robot.Webhook	= data.Webhook
	robot.Secret	= data.Secret
	robot.Keyword	= data.Keyword
	robot.Type		= data.Type
	robot.Status	= data.Status
	robot.Desc  	= data.Desc

	e := models.DB.Save(&robot).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改机器人通道成功", "", c)
}

// @Tags 系统管理
// @Description 机器人通道删除
// @Summary  机器人通道删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "通道ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system/robot/{id} [delete]
func DelRobot(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"setting-robot-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id := c.Param("id")

	// 删除之前检查是否有应用依赖
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("notify_id = ?", id).
		Find(&det)

	if det.Dtid > 0 {
		util.JsonRespond(500, "该通道被应用模板"+det.TemplateName + "引用，请先修改！", "", c)
		return
	}

	err := models.DB.Delete(models.SettingRobot{}, "id = ?", id).Error
	if err != nil {
		util.JsonRespond(500, err.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除机器人通道成功", "", c)
}

// @Tags 系统管理
// @Description 机器人通道测试
// @Summary  机器人通道测试
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "通道ID"
// @Param Data body admin.RobotResource true "机器人通道信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/system/robot/{id} [post]
func RobotTest(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"setting-robot-test") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data RobotResource
	e := c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, "Invalid Edit Robot Data", "", c)
		return
	}

	e = models.DingtalkSentTest(data.Type, data.Webhook, data.Secret, data.Keyword)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "机器人通道正常", "", c)
}



