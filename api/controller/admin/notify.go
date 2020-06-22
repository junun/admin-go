package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
)

type NotifyResource struct {
	Ids		string 	`form:"ids"`
}

// @Tags 通知管理
// @Description 通知列表
// @Summary  通知列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/notify [get]
func GetNotify(c *gin.Context)  {
	var notify []models.Notify
	data := make(map[string]interface{})

	models.DB.Model(&models.Notify{}).Where("unread=1").Find(&notify)
	data["lists"] = notify

	util.JsonRespond(200, "", data, c)
}

func PatchNotify(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"notify-read") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data 	NotifyResource
	//var notify 	models.Notify

	e 	:= c.BindJSON(&data)
	if e != nil {
		util.JsonRespond(500, "Invalid Patch Notify Data", "", c)
		return
	}

	arr := strings.Split(data.Ids, ",")
	if len(arr) > 1 {
		e := models.DB.Table("notify").Where("id IN (?)", arr).Updates(map[string]interface{}{"unread": 0}).Error
		if e != nil {
			util.JsonRespond(500, e.Error(), "", c)
		}
	} else {
		e := models.DB.Table("notify").Where("id = ?", data.Ids).Updates(map[string]interface{}{"unread": 0}).Error
		if e != nil {
			util.JsonRespond(500, e.Error(), "", c)
		}
	}

	util.JsonRespond(200, "", "", c)
}

