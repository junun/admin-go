package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
)

type PermResource struct {
	Name    	string    `form:"Name"`
	Permission  string    `form:"Permission"`
	Type 		int    	  `form:"Type"`
	Pid         int 	  `form:"Pid"`
	Desc        string    `form:"Desc"`
}

// 菜单权限
type Permissions struct {
	Id 				int 				`json:"id"`
	Pid				int					`json:"pid"`
	Name      		string				`json:"name"`
	Type			int					`json:"type"`
	Children    	[]*Permissions 		`json:"children"`
}

// @Tags 权限管理
// @Description 权限列表
// @Summary  权限列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/perms [get]
func GetPerms(c *gin.Context) {
	//if !middleware.PermissionCheckMiddleware(c,"perm-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var mps []models.MenuPermissions
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.MenuPermissions{}).Where(maps).Where("type = ?", 2).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&mps)
	models.DB.Model(&models.MenuPermissions{}).Where(maps).Where("type = ?", 2).Count(&count)

	data["lists"] = mps
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 权限管理
// @Description 权限添加
// @Summary  权限添加
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.PermResource true "权限信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/perms [post]
func PostPerms(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"perm-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data PermResource
	var mps  models.MenuPermissions

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add Perm Data", "", c)
		return
	}

	// 权限项唯一性检查
	models.DB.Model(&models.MenuPermissions{}).
		Where("permission = ?", data.Permission).
		Where("type = ?", 2).Find(&mps)

	if mps.ID > 0 {
		util.JsonRespond(500, "重复的标识符，请检查！", "", c)
		return
	}

	perm := models.MenuPermissions{
		Name: data.Name,
		Permission: data.Permission,
		Desc: data.Desc,
		Pid:data.Pid,
		Type: data.Type}

	e := models.DB.Save(&perm).Error

	if e != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	// 修改perm 信息， 最简单的做法， 删除redis 对应的key
	DelRedisAllPermKey()

	util.JsonRespond(200, "添加菜单权限按钮成功", "", c)

}

// @Tags 权限管理
// @Description 权限修改
// @Summary  权限修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Param Data body admin.PermResource true "权限信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/perms/{id} [put]
func PutPerms(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"perm-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data PermResource
	var perm models.MenuPermissions
	var mps  models.MenuPermissions

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Perm Data", "", c)
		return
	}

	// 权限项唯一性检查
	models.DB.Model(&models.MenuPermissions{}).
		Where("permission = ?", data.Permission).
		Where("id != ?", c.Param("id")).Find(&mps)

	if mps.ID > 0 {
		util.JsonRespond(500, "重复的标识符，请检查！", "", c)
		return
	}

	models.DB.Find(&perm, c.Param("id"))

	perm.Name = data.Name
	perm.Desc = data.Desc
	perm.Pid  = data.Pid
	perm.Permission = data.Permission

	e := models.DB.Save(&perm).Error
	if e != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	// 修改perm 信息， 最简单的做法， 删除redis 对应的key
	DelRedisAllPermKey()

	util.JsonRespond(200, "修改权限按钮成功", "", c)
}

// @Tags 权限管理
// @Description 权限删除
// @Summary  权限删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/perms/{id} [delete]
func DeletePerms(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"perm-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	err := models.DB.Delete(models.MenuPermissions{}, "id = ?", c.Param("id")).Error
	if err != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	// 修改perm 信息， 最简单的做法， 删除redis 对应的key
	DelRedisAllPermKey()

	util.JsonRespond(200, "删除权限按钮成功", "", c)
}

// @Tags 权限管理
// @Description 获取所以权限项
// @Summary  获取所以权限项
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/perms/lists [get]
func GetAllPerms(c *gin.Context)  {
	var mps []*models.MenuPermissions
	var res []models.MenuPermissions

	data   	:= make(map[string]interface{})
	tmp   	:= make(map[int]*models.MenuPermissions)

	// 所以的mod page perm组合数据 放到redis里面
	key 	:= models.AllPermsKey
	v, _  	:= models.Rdb.Get(key).Result()

	if v != "" {
		data["lists"] = util.JsonUnmarshalFromString(v, &res)
		util.JsonRespond(200, "", data, c)
		return
	}


	models.DB.Model(&models.MenuPermissions{}).Find(&mps)

	for _, p := range mps {
		if x, ok := tmp[p.ID]; ok {
			p.Children = x.Children
		}
		tmp[p.ID] = p
		if p.Pid != 0 {
			if x, ok  := tmp[p.Pid]; ok {
				x.Children = append(x.Children, p)
			} else  {
				tmp[p.Pid] = &models.MenuPermissions{
					Children: []*models.MenuPermissions{p},
				}
			}
		}
	}

    for _, v := range tmp {
		if v.Pid == 0 {
			res = append(res, *v)
		}
	}

	models.Rdb.Set(key, util.JSONMarshalToString(res), 0)

    data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

func DelRedisAllPermKey()  {
	key := models.AllPermsKey
	models.Rdb.Del(key).Err()
}