package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RoleResource struct {
	Name    	string    	`form:"Name"`
	Desc 		string    	`form:"Desc"`
}

type RolePermResource struct {
	Codes       []int 		`json:Codes`
}

var  menu_redis_set_key string

// @Tags 角色管理
// @Description 角色列表
// @Summary  角色列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/roles [get]
func GetRole(c *gin.Context)  {
	var roles []models.Role
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.Role{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&roles)
	models.DB.Model(&models.Role{}).Where(maps).Count(&count)

	data["lists"] = roles
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 角色管理
// @Description 新增角色
// @Summary  新增角色
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.RoleResource true "角色信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/roles [post]
func PostRole(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"role-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data RoleResource
	var role models.Role

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add Role Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.Role{}).
		Where("name = ?", data.Name).Find(&role)

	if role.ID > 0 {
		util.JsonRespond(500, "重复的角色名，请检查！", "", c)
		return
	}

	role = models.Role{
		Name: data.Name,
		Desc: data.Desc}

	e := models.DB.Save(&role).Error

	if e != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	util.JsonRespond(200, "添加角色成功", "", c)
}

// @Tags 角色管理
// @Description 修改角色
// @Summary  修改角色
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.RoleResource true "角色信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/roles/{id} [put]
func PutRole(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"role-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data RoleResource
	var role models.Role

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Role Data", "", c)
		return
	}

	// 角色名唯一性检查
	models.DB.Model(&models.Role{}).
		Where("name = ?", data.Name).
		Where("id != ?", c.Param("id")).Find(&role)

	if role.ID > 0 {
		util.JsonRespond(500, "重复的角色名，请检查！", "", c)
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

	util.JsonRespond(200, "修改角色成功", "", c)
}

// @Tags 角色管理
// @Description 删除角色
// @Summary  删除角色
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "角色id"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/roles/{id} [delete]
func DeleteRole(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"role-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	err := models.DB.Delete(models.Role{}, "id = ?", c.Param("id")).Error
	if err != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	util.JsonRespond(200, "删除成功", "", c)
}

// @Tags 角色管理
// @Description 角色权限详情
// @Summary  角色权限详情
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/roles/{id}/permissions [get]
func GetRolePerms(c *gin.Context) {
	if !middleware.PermissionCheckMiddleware(c,"role-perm-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data map[string]interface{}
	var prole []models.RolePermissionRel
	data = make(map[string]interface{})

	models.DB.Model(&models.RolePermissionRel{}).Select("pid").
		Where("rid = ?", c.Param("id")).Find(&prole)

	data["lists"] = prole

	util.JsonRespond(200, "", data, c)
}

// @Tags 角色管理
// @Description 添加/修改角色权限
// @Summary  添加/修改角色权限
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Param Data body admin.RolePermResource true "角色权限信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/roles/{id}/permissions [post]
func PostRolePerms(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"role-perm-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data  RolePermResource
	var old_prole []models.RolePermissionRel
	var rpr models.RolePermissionRel
	var mps []models.MenuPermissions

	rds 		:= make(map[int]interface{})
	prole 		:= make(map[int]interface{})
	new_prole 	:= make(map[int]interface{})
	rid,_ 		:= strconv.Atoi(c.Param("id"))

	e 	:= c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, "Invalid MenuPermissions Data", "", c)
		return
	}

	models.DB.Model(&models.RolePermissionRel{}).Select("pid").
		Where("rid = ?", rid).Find(&old_prole)

	// 可以把所有的 type=1 的菜单选项id 放到 rds队列里
	models.DB.Model(&models.MenuPermissions{}).Select("id").
		Where("type = ?", 1).Find(&mps)

	for _, v := range mps {
		rds[v.ID]	= v.ID
	}

	for _, v := range data.Codes {
		//m, _ := strconv.Atoi(v)
		if _, ok  := rds[v]; ok {
			continue
		}

		new_prole[v] = v
	}

	// 删除
	for _, k := range old_prole {
		prole[k.Pid] = k.Pid
 		if _, ok  := new_prole[k.Pid]; !ok {
			// 执行删除操作
			e 	:= models.DB.Delete(models.RolePermissionRel{}, "pid = ?", k.Pid).Error
			if e != nil {
				util.JsonRespond(500, e.Error(), "", c)
				return
			}
		}
	}

	// 新增
	for k,_ := range new_prole {
		if _, ok  := prole[k]; !ok {
			//执行新增操作，换成批量插入更好
			rpr = models.RolePermissionRel{
				Pid: k,
				Rid: rid}

			e := models.DB.Save(&rpr).Error

			if e != nil {
				util.JsonRespond(500, e.Error(), "", c)
				return
			}
		}
	}

	//更新redis里面的角色的权限集合
	key := models.RoleRermsSetKey
	str := fmt.Sprintf("%v", rid)
	key =  key + str

	// 先删除key
	models.DelKey(key)
	// 再添加
	models.SetRolePermToSet(key, rid)

	util.JsonRespond(200, "","", c)
}
