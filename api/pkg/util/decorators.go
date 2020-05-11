package util

import (
	"github.com/gin-gonic/gin"
)

func RequirePermission(perms string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//if admin.G["is_supper"] != 1 {
		//	perm_list := strings.Split(perms,"|")
		//
		//	for _, item := range perm_list {
		//		if !Contains(admin.G["permissions"].([]string), item) {
		//			JsonRespond(403, "Permission denied", "", c)
		//		}
		//	}
		//}

		c.Next()
	}
}
