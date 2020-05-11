package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/setting"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
)

type MenuResource struct {
	Name    	string    `form:"Name"`
	Icon 		string    `form:"Icon"`
	Type 		int    	  `form:"Type"`
}

type SubMenuResource struct {
	Name    	string    `form:"Name"`
	Url 		string    `form:"Url"`
	Icon 		string    `form:"Icon"`
	Type 		int    	  `form:"Type"`
	Pid    		int 	  `form:"Pid"`
	Desc        string    `form:"Desc"`
}

// @Tags 菜单管理
// @Description 一级菜单列表
// @Summary  一级菜单列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/menus [get]
func GetMenus(c *gin.Context) {
	var menus []models.MenuPermissions
	var count int

	mtype := c.Query("type")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if mtype != "" {
		maps["type"] = mtype
	} else {
		maps["type"] = 1
	}

	maps["pid"] = 0

	models.DB.Model(&models.MenuPermissions{}).Where(maps).
		Offset(util.GetPage(c)).
		Limit(setting.AppSetting.PageSize).
		Find(&menus)

	models.DB.Model(&models.MenuPermissions{}).Where(maps).Count(&count)

	data["lists"] = menus
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 菜单管理
// @Description 一级菜单添加， type must 1
// @Summary  一级菜单添加
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.MenuResource true "一级菜单信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/menus [post]
func PostMenus(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"menu-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var dataResource MenuResource
	err := c.BindJSON(&dataResource)
	if err != nil {
		util.JsonRespond(500, "Invalid Add Menu Data", "", c)
		return
	}

	menu := models.MenuPermissions{
		Name: dataResource.Name,
		Icon: dataResource.Icon,
		Type: dataResource.Type}

	e := models.DB.Save(&menu).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	// 修改管理员菜单列表， 最简单的做法， 删除redis role 为0 的记录
	DelRedisAdminRoleRecord()

	util.JsonRespond(200, "添加一级菜单成功", "", c)
}

// @Tags 菜单管理
// @Description 一级菜单修改
// @Summary  一级菜单修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Param Data body admin.MenuResource true "一级菜单信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/menus/{id} [put]
func PutMenus(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"menu-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data MenuResource
	var menu models.MenuPermissions
	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Menu Data", "", c)
		return
	}

	models.DB.Find(&menu, c.Param("id"))

	menu.Name = data.Name
	menu.Icon = data.Icon

	e := models.DB.Save(&menu).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	// 修改管理员菜单列表， 最简单的做法， 删除redis role 为0 的记录
	DelRedisAdminRoleRecord()

	util.JsonRespond(200, "修改一级菜单成功", "", c)
}

// @Tags 菜单管理
// @Description 菜单删除
// @Summary  菜单删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/menus/{id} [delete]
func DeleteMenus(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"menu-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	e := models.DB.Delete(models.MenuPermissions{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}
	// 修改管理员菜单列表， 最简单的做法， 删除redis role 为0 的记录
	DelRedisAdminRoleRecord()

	util.JsonRespond(200, "删除成功", "", c)
}

// @Tags 菜单管理
// @Description 二级菜单列表
// @Summary  二级菜单列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/submenus [get]
func GetSubMenu(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"submenu-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var res []models.MenuPermissions
	var count int
	mtype := c.Query("type")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if mtype != "" {
		maps["type"] = mtype
	} else {
		maps["type"] = 1
	}

	models.DB.Model(&models.MenuPermissions{}).Where(maps).Where("pid > ?", 0).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&res)
	models.DB.Model(&models.MenuPermissions{}).Where(maps).Where("pid > ?", 0).
		Count(&count)

	data["lists"] = res
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 菜单管理
// @Description 二级菜单添加， type must 1
// @Summary  二级菜单添加
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.SubMenuResource true "二级菜单信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/submenus [post]
func PostSubMenu(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"submenu-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data SubMenuResource

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add SubMenu Data", "", c)
		return
	}

	menu := models.MenuPermissions{
		Name: data.Name,
		Icon: data.Icon,
		Desc: data.Desc,
		Pid:data.Pid,
		Url: data.Url,
		Type: data.Type}

	e := models.DB.Save(&menu).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	// 修改管理员菜单列表， 最简单的做法， 删除redis role 为0 的记录
	DelRedisAdminRoleRecord()

	util.JsonRespond(200, "添加二级菜单成功", "", c)
}

// @Tags 菜单管理
// @Description 二级菜单修改
// @Summary  二级菜单修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Param Data body admin.SubMenuResource true "二级菜单信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/submenus/{id} [put]
func PutSubMenus(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"submenu-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data SubMenuResource
	var menu models.MenuPermissions

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit SubMenu Data", "", c)
		return
	}

	models.DB.Find(&menu, c.Param("id"))

	menu.Name = data.Name
	menu.Icon = data.Icon
	menu.Url  = data.Url
	menu.Desc = data.Desc
	menu.Pid  = data.Pid

	e := models.DB.Save(&menu).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	// 修改管理员菜单列表， 最简单的做法， 删除redis role 为0 的记录
	DelRedisAdminRoleRecord()

	util.JsonRespond(200, "修改二级菜单成功", "", c)
}

func DelRedisAdminRoleRecord()  {
	// 硬性编码管理员的rid为0
	key := models.RoleMenuListKey+"0"
	models.Rdb.Del(key).Err()
}