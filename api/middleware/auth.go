package middleware

import (
	"api/models"
	"api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user  models.User

		// 如果是登录操作请求 不检查Token
		uri := c.Request.URL.String()
		if uri == "/admin/user/login" {
			return
		}

		token := c.Request.Header.Get("Authorization")

		if token == "" {
			util.JsonRespond(401, "API token required", "", c)
			c.Abort()
			return
		}

		e := models.DB.Where("access_token = ?", token).First(&user).Error
		if e != nil{
			util.JsonRespond(401, "Invalid API token, please login", "", c)
			c.Abort()
		}

		if user.IsActive == 1 && user.TokenExpired > time.Now().Unix() {
			user.TokenExpired = time.Now().Unix() + 86400
			models.DB.Save(&user)
		}

		c.Set("UserIsSupper", user.IsSupper)
		c.Set("UserRid", user.Rid)
		c.Set("Uid", user.ID)
		c.Next()
	}
}

func PermissionCheckMiddleware(c *gin.Context, perm string) bool{
	UserIsSupper, _ := c.Get("UserIsSupper")

	// 超级用户不做权限检查
	if UserIsSupper != 1 {
		key 		:= models.RoleRermsSetKey
		UserRid, _ 	:= c.Get("UserRid")
		str 		:= fmt.Sprintf("%v", UserRid)
		redis_key 	:=  key + str

		// 检查 redis 有没有该key的集合
		err := models.Rdb.Exists(redis_key).Val()
		if err != 1 {
			rid, _ := UserRid.(int)
			models.SetRolePermToSet(redis_key, rid)
		}

		// 检查对应的set是否有该角色权限
		return models.CheckMemberByKey(redis_key, perm)
	} else {
		return true
	}
}

//func WsTokenAuthMiddleware(pemr string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token := c.Param("token")
//
//		if token == "" {
//			util.JsonRespond(401, "API token required", "", c)
//			c.Abort()
//			return
//		}
//
//		var user  models.User
//		e := models.DB.Where("access_token = ?", token).First(&user).Error
//		if e != nil{
//			util.JsonRespond(401, "Invalid API token, please login", "", c)
//			c.Abort()
//		}
//
//		if user.IsSupper != 1 {
//			key 		:= models.RoleRermsSetKey
//			UserRid 	:= user.Rid
//			str 		:= fmt.Sprintf("%v", UserRid)
//			redis_key 	:=  key + str
//
//			// 检查 redis 有没有该key的集合
//			err := models.Rdb.Exists(redis_key).Val()
//			if err != 1 {
//				rid 	:= UserRid
//				models.SetRolePermToSet(redis_key, rid)
//			}
//
//			// 检查对应的set是否有该角色权限
//			if models.CheckMemberByKey(redis_key, pemr) {
//				util.JsonRespond(403, "请求资源被拒绝", "", c)
//				return
//			}
//		}
//
//		c.Set("Uid", user.ID)
//
//		c.Next()
//	}
//}


func WsTokenAuthMiddleware(pemr string, c *gin.Context)  {
	token := c.Param("token")
	if token == "" {
		util.JsonRespond(401, "API token required", "", c)
		c.Abort()
	}

	var user  models.User
	e := models.DB.Where("access_token = ?", token).First(&user).Error
	if e != nil{
		util.JsonRespond(401, "Invalid API token, please login", "", c)
		c.Abort()
	}

	if user.IsSupper != 1 {
		key 		:= models.RoleRermsSetKey
		UserRid 	:= user.Rid
		str 		:= fmt.Sprintf("%v", UserRid)
		RedisKey 	:= key + str

		// 检查 redis 有没有该key的集合
		err 	:= models.Rdb.Exists(RedisKey).Val()
		if err != 1 {
			rid 	:= UserRid
			models.SetRolePermToSet(RedisKey, rid)
		}

		// 检查对应的set是否有该角色权限
		if !models.CheckMemberByKey(RedisKey, pemr) {
			util.JsonRespond(403, "请求资源被拒绝", "", c)
			return
		}
	}

	c.Set("Uid", user.ID)
}

func RoleAppAuthMiddleware(rid, eid, aid int, c *gin.Context) bool {
	var app  models.RoleEnvApp
	models.DB.Model(&models.RoleEnvApp{}).
		Where("rid = ?", rid).
		Where("eid = ?", eid).
		Find(&app)

	if app.ID == 0 {
		util.JsonRespond(403, "该角色没有当前环境下该应用的操作权限！", "", c)
		return false
	}

	if !util.StrArrContains(strings.Split(app.AppIds, ","), strconv.Itoa(aid)) {
		util.JsonRespond(403, "请求应用资源被拒绝！", "", c)
		return false
	}

	return true
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}