package routers

import (
	"api/controller/admin"
	_ "api/docs"
	"api/middleware"
	"api/pkg/setting"
	"api/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	//r.Use(cors.Default())

	// 跨域问题
	r.Use(middleware.Cors())

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(middleware.MetricMiddleware())

	r.LoadHTMLGlob("templates/*")

	gin.SetMode(setting.RunMode)

	//r.GET("/debug/pprof", gin.WrapF(pprof.Index))
	//r.GET("/debug/heap", gin.WrapH(pprof.Handler("heap")))
	//r.GET("/debug/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	//r.GET("/debug/cmdline", gin.WrapH(pprof.Handler("cmdline")))
	//r.GET("/debug/profile", gin.WrapF(pprof.Profile))
	//r.GET("/debug/symbol", gin.WrapH(pprof.Handler("symbol")))
	//r.GET("/debug/trace", gin.WrapH(pprof.Handler("trace")))

	// metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// gin swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 业务初始化
	r.GET("/admin/exec/ws/:id/ssh/:token", admin.SyncApp)

	// 业务上线/回滚
	r.GET("/admin/deploy/ws/:id/ssh/:token", admin.AppDeployRedo)

	// ConsoleHost 逻辑
	r.GET("/admin/host/ssh/:id", admin.ConsoleHost)
	r.GET("/admin/ws/:id/ssh/:token", admin.SshConsumer)

	//r.POST("/upload", admin.UploadImage)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	//apiv1 := r.Group("/api")
	//{
	//}

	r.GET("/admin/test", admin.Test)

	adminv1 := r.Group("/admin", middleware.TokenAuthMiddleware())
	{
		// 用户登录
		adminv1.POST("/user/login", admin.Login)

		adminv1.POST("/user/logout", admin.Logout)
		adminv1.GET("/user/perms/:id", admin.GetUserMenu)

		adminv1.GET("/user", admin.GetUsers)
		adminv1.POST("/user", admin.PostUser)
		adminv1.PATCH("/user", admin.PatchUser)
		adminv1.PUT("/user/:id", admin.PutUser)
		adminv1.DELETE("/user/:id", admin.DeleteUser)
		adminv1.GET("/perms", admin.GetPerms)
		adminv1.POST("/perms", admin.PostPerms)
		adminv1.DELETE("/perms/:id", admin.DeletePerms)
		adminv1.PUT("/perms/:id", admin.PutPerms)
		adminv1.GET("/perms/lists", admin.GetAllPerms)

		adminv1.GET("/system", admin.GetSetting)
		adminv1.POST("/system", admin.SettingModify)
		adminv1.GET("/system/about", admin.About)
		adminv1.GET("/system/robot", admin.GetRobot)
		adminv1.POST("/system/robot", admin.AddRobot)
		adminv1.PUT("/system/robot/:id", admin.PutRobot)
		adminv1.DELETE("/system/robot/:id", admin.DelRobot)
		adminv1.POST("/system/robot/:id", admin.RobotTest)
		adminv1.POST("/system/mail", admin.EmailTest)

		adminv1.GET("/roles",admin.GetRole)
		adminv1.POST("/roles", admin.PostRole)
		adminv1.DELETE("/roles/:id", admin.DeleteRole)
		adminv1.PUT("/roles/:id", admin.PutRole)
		adminv1.GET("/roles/:id/permissions", admin.GetRolePerms)
		adminv1.POST("/roles/:id/permissions", admin.PostRolePerms)
		adminv1.GET("/roles/:id/app", admin.GetRoleApp)
		adminv1.POST("/roles/:id/app", admin.PostRoleApp)
		adminv1.GET("/roles/:id/host", admin.GetRoleHost)
		adminv1.POST("/roles/:id/host", admin.PostRoleHost)
		adminv1.GET("/env/app", admin.GetEnvApp)
		adminv1.GET("/env/host", admin.GetEnvHost)

		adminv1.GET("/menus", admin.GetMenus)
		adminv1.POST("/menus", admin.PostMenus)
		adminv1.DELETE("/menus/:id", admin.DeleteMenus)
		adminv1.PUT("/menus/:id", admin.PutMenus)

		adminv1.GET("/submenus", admin.GetSubMenu)
		adminv1.POST("/submenus", admin.PostSubMenu)
		adminv1.PUT("/submenus/:id", admin.PutSubMenus)
		adminv1.DELETE("/submenus/:id", admin.DeleteMenus)

		adminv1.GET("/domain/info", admin.GetDomainInfo)
		adminv1.POST("/domain/info", admin.AddDomainInfo)
		adminv1.PUT("/domain/info/:id", admin.PutDomainInfo)
		adminv1.DELETE("/domain/info/:id", admin.DelDomainInfo)

		adminv1.GET("/host/role", admin.GetHostRole)
		adminv1.POST("/host/role", admin.AddHostRole)
		adminv1.PUT("/host/role/:id", admin.PutHostRole)
		adminv1.DELETE("/host/role/:id", admin.DelHostRole)
		adminv1.POST("/host/import", admin.ImportHost)

		adminv1.GET("/host", admin.GetHost)
		adminv1.POST("/host", admin.AddHost)
		adminv1.PUT("/hosts/:id", admin.PutHost)
		adminv1.DELETE("/hosts/:id", admin.DelHost)

		adminv1.GET("/host/app", admin.GetHostApp)
		adminv1.POST("/host/app", admin.AddHostApp)
		adminv1.PUT("/host/app/:id", admin.PutHostApp)
		adminv1.DELETE("/host/app/:id", admin.DelHostApp)
		adminv1.GET("/host/appid", admin.GetHostByAppId)

		adminv1.GET("/config/env", admin.GetConfigEnv)
		adminv1.POST("/config/env", admin.AddConfigEnv)
		adminv1.PUT("/config/env/:id", admin.PutConfigEnv)
		adminv1.DELETE("/config/env/:id", admin.DelConfigEnv)

		adminv1.GET("/config/type", admin.GetAppType)
		adminv1.POST("/config/type", admin.AddAppType)
		adminv1.PUT("/config/type/:id", admin.PutAppType)
		adminv1.DELETE("/config/type/:id", admin.DelAppType)

		adminv1.GET("/config/app", admin.GetConfigApp)
		adminv1.POST("/config/app", admin.AddConfigApp)
		adminv1.PUT("/config/app/:id", admin.PutConfigApp)
		adminv1.DELETE("/config/app/:id", admin.DelConfigApp)
		adminv1.GET("/config/template", admin.GetAppTemplate)

		adminv1.GET("/config/value", admin.GetAppValue)
		adminv1.POST("/config/value", admin.AddAppValue)
		adminv1.PUT("/config/value/:id", admin.PutAppValue)
		adminv1.DELETE("/config/value/:id", admin.DelAppValue)

		adminv1.GET("/config/deploy", admin.GetDeployExtend)
		adminv1.POST("/config/deploy", admin.AddDeployExtend)
		adminv1.PUT("/config/deploy/:id", admin.PutDeployExtend)
		adminv1.DELETE("/config/deploy/:id", admin.DelDeployExtend)
		adminv1.GET("/sync/request/:id", admin.AppSyncRequest)

		adminv1.GET("/deploy/app", admin.GetAppDeploy)
		adminv1.POST("/deploy/app", admin.AddAppDeploy)
		adminv1.PUT("/deploy/app/:id", admin.PutAppDeploy)
		adminv1.DELETE("/deploy/app/:id", admin.DelAppDeploy)
		adminv1.PUT("/undo/confirm/:id", admin.PutUndoRequest)
		adminv1.GET("/undo/request/:id", admin.GetUndoRequest)
		adminv1.GET("/deploy/request/:id", admin.GetDeployRequest)
		adminv1.POST("/deploy/request/:id", admin.PostDeployRequest)
		adminv1.PUT("/deploy/app/:id/review", admin.PutAppDeployStatus)
		//adminv1.PUT("/deploy/app/:id/undo/:status", admin.PutAppDeployUndo)
		adminv1.GET("/deploy/app/:id/branch", admin.GetGitBranch)
		adminv1.GET("/deploy/app/:id/tag", admin.GetGitTag)
		adminv1.GET("/deploy/app/:id/commit/:branch", admin.GetGitCommit)
		//adminv1.GET("/deploy/app/:id/version", admin.GetAppVersion)


		adminv1.GET("/notify", admin.GetNotify)
		adminv1.PATCH("/notify", admin.PatchNotify)

		adminv1.GET("/schedule", admin.GetJobList)
		adminv1.PATCH("/schedule", admin.PatchJobStatus)
		adminv1.POST("/schedule", admin.AddJob)
		adminv1.PUT("/schedule/:id", admin.PutJob)
		adminv1.DELETE("/schedule/:id", admin.DelJob)
		adminv1.GET("/schedule/:id", admin.GetJobHisById)
		adminv1.GET("/schedule/:id/info", admin.GetJobInfo)
	}

	return r
}
