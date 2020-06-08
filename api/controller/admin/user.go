package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type LoginResource struct {
	Username    string    `form:"username"`
	Password 	string    `form:"password"`
}

type UserResource struct {
	Name   		string    	`form:"Name"`
	Nickname    string		`form:"Nickname"`
	Mobile      string		`form: Mobile`
	Email 		string 		`form: Email`
	Rid 		int 		`form: Rid`
	Password 	string    	`form:"password"`
	IsActive	int 		`form:IsActive`
}

// @Tags 用户管理
// @Description 用户登录
// @Summary  用户登录
// @Produce  json
// @Param Data body admin.LoginResource true "用户登录信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "用户不存在！", "data": {}}
// @Router /admin/user/login [post]
func Login(c *gin.Context)  {
	var dataResource LoginResource
	var user models.User
	var val interface{}
	var expiration  = time.Duration(86400)*time.Second

	err := c.BindJSON(&dataResource)
	if err != nil {
		util.JsonRespond(500, "Invalid Login data", "", c)
		return
	}

	username	:= dataResource.Username
	password 	:= dataResource.Password
	key 		:= username+"_login"
	e := models.DB.Where("name = ?", username).First(&user).Error
	if e != nil{
		util.JsonRespond(500, "用户不存在！", "", c)
		return
	}

	if user.IsActive == 1 {
		err := util.CheckPasswordHash(password, user.PasswordHash)
		if !err {
			// 记录用户验证失败次数
			// 检查key是否存在 1: 存在， 0: 不存在
			if models.Rdb.Exists(key).Val() == 1 {
				// 获取key的值
				val = models.GetValByKey(key)
				res := val.(int)

				// 如果验证失败次数多于3次，将锁定用户
				if res > 3 {
					util.JsonRespond(401, "用户已被禁用，请联系管理员", "", c)
					return
				}

				if err := models.SetValByKey(key, res+1, expiration); err != nil {
					panic(err)
				}
			} else {
				// 第一次登录失败
				if err := models.SetValByKey(key, 1, expiration); err != nil {
					panic(err)
				}
			}
			util.JsonRespond(401, "用户名或密码错误，连续3次错误将会被禁用", "", c)
			return
		} else {
			//生成token
			token := uuid.New().String()
			user.AccessToken = token
			user.TokenExpired = time.Now().Unix() + 86400

			//提交更改
			models.DB.Save(&user)

			// 获取用户的权限列表
			var permissions  []string
			if user.IsSupper != 1 {
				permissions = user.ReturnPermissions()
			}

			data := make(map[string]interface{})
			data["rid"]			= user.Rid
			data["token"] 		= token
			data["is_supper"] 	= user.IsSupper
			data["nickname"]	= user.Nickname
			data["permissions"]	= permissions

			// 登录成功
			if err := models.SetValByKey(key, 0, expiration); err != nil {
				panic(err)
			}

			util.JsonRespond(200, "",data, c)
			return
		}
	} else {
		util.JsonRespond(500, "用户被禁用，请联系管理员！", "", c)
		return
	}
}

// @Tags 用户管理
// @Description 用户登出
// @Summary  用户登出
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "登出失败", "data": {}}
// @Router /admin/user/logout [post]
func Logout(c *gin.Context)  {
	var user models.User

	Uid, _ := c.Get("Uid")

	models.DB.Find(&user, Uid)
	user.AccessToken = ""

	e := models.DB.Save(&user).Error
	if e != nil {
		util.JsonRespond(500, "内部错误", "", c)
		return
	}

	util.JsonRespond(200, "退出成功！", "", c)
}

// @Tags 用户管理
// @Description 用户列表
// @Summary  用户列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/user [get]
func GetUsers(c *gin.Context) {
	//if !middleware.PermissionCheckMiddleware(c,"user-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var res []models.User
	var count int

	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	models.DB.Model(&models.User{}).Where(maps).
		Where("rid > 0").
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&res)
	models.DB.Model(&models.User{}).Where(maps).
		Where("rid > 0").
		Count(&count)

	data["lists"] = res
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 用户管理
// @Description 新增用户
// @Summary 新增用户
// @Accept  application/json
// @Produce application/json
// @Param Authorization header string true "token"
// @Param Data body admin.UserResource true "用户信息"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/user [post]
func PostUser(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"user-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var user models.User
	var data UserResource

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add User Data", "", c)
		return
	}

	// 用户唯一性检查
	models.DB.Model(&models.User{}).
		Where("name = ?", data.Name).
		Find(&user)

	if user.ID > 0 {
		util.JsonRespond(500, "重复的用户名，请检查！", "", c)
		return
	}

	PasswordHash, err := util.HashPassword(data.Password)

	if err != nil {
		util.JsonRespond(500, "hash密码错误，请联系管理员！", "", c)
		return
	}

	myuser := models.User{
		Name: data.Name,
		Nickname: data.Nickname,
		Mobile: data.Mobile,
		Email:data.Email,
		IsActive:1,
		PasswordHash: PasswordHash,
		Rid: data.Rid}

	e := models.DB.Save(&myuser).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加用户成功", "", c)
}

// @Tags 用户管理
// @Description 修改用户
// @Summary 修改用户
// @Accept  application/json
// @Produce application/json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Param Data body models.User true "用户信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/user/{id} [put]
func PutUser(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"user-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var user models.User
	var data UserResource

	fmt.Println(c.PostForm("Email"))

	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
		util.JsonRespond(400, "Bad Request : Invalid Edit User Data", "", c)
		return
	}

	models.DB.Find(&user, c.Param("id"))

	user.Nickname 	= data.Nickname
	user.Mobile  	= data.Mobile
	user.Email 		= data.Email
	user.Rid  		= data.Rid
	user.IsActive	= data.IsActive

	if len(data.Password) > 0 {
		PasswordHash, err := util.HashPassword(data.Password)

		if err != nil {
			util.JsonRespond(500, "hash密码错误，请联系管理员！", "", c)
			return
		}

		user.PasswordHash = PasswordHash
	}

	e := models.DB.Save(&user).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改用户成功", "", c)
}

// @Tags 用户管理
// @Description 删除用户
// @Summary  删除用户
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/user/{id} [delete]
func DeleteUser(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"user-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	e 	:= models.DB.Delete(models.User{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除成功", "", c)
}

// @Tags 用户管理
// @Description 用户菜单列表
// @Summary  用户菜单列表
// @Produce  json
// @Param Authorization header string true "token"
// @Param id query integer true "ID"
// @Success 200 {string} string "{"code": 200, "message": "", "data": {}}"
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/user/perms/{id} [get]
func GetUserMenu(c *gin.Context)  {
	var mps []*models.MenuPermissions
	var res []models.MenuPermissions

	tmp   	:= make(map[int]*models.MenuPermissions)
	data    := make(map[string]interface{})
	rid 	:= c.Param("id")

	// 由于返回一次用户菜单列表 需要查询数据库进行各种组织操作， 故建议放到redis里面
	key 	:= models.RoleMenuListKey
	str 	:= fmt.Sprintf("%v", rid)
	key 	=  key + str
	v, _  	:= models.Rdb.Get(key).Result()

	if v != "" {
		data["lists"] = util.JsonUnmarshalFromString(v, &res)
		util.JsonRespond(200, "", data, c)
		return
	}

	if rid == "0" {
		// 超级用户直接返回所有的菜单
		models.DB.Model(&models.MenuPermissions{}).
			Where("type = 1").
			Find(&mps)

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
	} else {
		// 普通用户 根据 rid 返回菜单项
		pids := []int{}
		models.DB.Table("menu_permissions").
			Select("menu_permissions.pid").
			Joins("left join role_permission_rel on menu_permissions.id = role_permission_rel.pid").
			Where("role_permission_rel.rid = ?", rid).
			Pluck("DISTINCT menu_permissions.pid", &pids)

		models.DB.Model(&models.MenuPermissions{}).
			Where("type = ?", 1).
			Find(&mps)

		for _,v := range mps {
			for _,p := range pids {
				if _, ok := tmp[v.ID]; !ok {
					tmp[v.ID] = v
				}

				if p == v.ID {
					if x, ok  := tmp[v.Pid]; ok {
						x.Children = append(x.Children, v)
					}
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
