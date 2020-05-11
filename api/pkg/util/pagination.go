package util

import (
	"api/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {
	result 	:= 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * GetPageSize(c)
	}

	return result
}

func GetPageSize(c *gin.Context) int {
	result,_ := com.StrTo(c.Query("pagesize")).Int()

	if result == 0 {
		return setting.AppSetting.PageSize
	}

	return result
}