package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/unknwon/com"
	"strconv"
	"strings"
	"time"
)

type DeployApp struct {
	Name    	string    	`form:"Name"`
	Tid			int			`form:"Tid"`
	RepoBranch	string		`form:"RepoBranch"`
	RepoCommit 	string		`form:"RepoCommit"`
	Status		int			`form:"Status"`
}

type AppTemplateDeploy struct {
	ID              int
	Aid             int
	Tid				int
	Tag				string
	Name      		string
	RepoBranch 		string
	RepoCommit 		string
	Status			int
	Operator		int
	Review          int
	Deploy          int
	UpdateTime 		time.Time
}

const (
	NewDeploy		= 1
	ReviewSuccess  	= 2
	ReviewFail 		= 3
	DeployFail		= 4
	DeploySuccess 	= 5
	UndoSuccess     = 6
	UndoFail 		= 7
)


func GetGitBranch(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-git") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var det models.DeployExtend

	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", c.Param("id")).Find(&det)

	res , e := util.ReturnGitBranch(det.Aid, det.RepoUrl)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

func GetGitCommit(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"config-app-git") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var det models.DeployExtend

	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", c.Param("id")).Find(&det)

	branch := c.Param("branch")

	// 锁定项目，一个项目同时只能允许一个执行该方法
	key := models.GitAppOnWorking + det.TemplateName
	if models.GetValByKey(key) != "" {
		util.JsonRespond(500, "该项目别的用户在使用中，请稍后重试！", "", c)
		return
	}

	models.SetValByKey(key, "1",  2 * time.Second)

	res, e := util.GetGitLastTenCommitByBranch(det.Aid, det.RepoUrl, branch)
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	models.Rdb.Del(key)

	data := make(map[string]interface{})
	data["lists"] = res

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用发布
// @Description 应用发布列表
// @Summary  应用发布列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app [get]
func GetAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	data := make(map[string]interface{})

	var deploy []AppTemplateDeploy
	// 分页逻辑还没有写，有空补上。
	e := models.DB.Raw("select d.*, e.aid from app_deploy d left join deploy_extend e on d.tid=e.dtid").Scan(&deploy).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	data["lists"] = deploy
	data["total"] = len(deploy)

	util.JsonRespond(200, "", data, c)
}

// @Tags 应用发布
// @Description 应用发布提单
// @Summary 应用发布提单
// @Produce json
// @Param Authorization header string true "token"
// @Param Data body admin.DeployApp true "应用发布信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app [post]
func AddAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var data DeployApp
	var deploy models.AppDeploy

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(406, "Invalid Add Deploy Data", "", c)
		return
	}

	// 发布唯一性检查
	models.DB.Model(&models.AppDeploy{}).
		Where("tid = ?", data.Tid).
		Where("repo_branch = ?", data.RepoBranch).
		Where("repo_commit = ?", data.RepoCommit).
		Where("status <= 2").
		Find(&deploy)

	if deploy.ID > 0 {
		util.JsonRespond(500, "重复的项目上线提单，请检查！", "", c)
		return
	}

	// 检查是否需要开启审核
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", data.Tid).
		Find(&det)

	sts := NewDeploy
	if  det.EnableCheck == 0  {
		sts = ReviewSuccess
	}

	uid,_ 	:= c.Get("Uid")
	deploy = models.AppDeploy{
		Name: data.Name,
		Tid: data.Tid,
		RepoBranch: data.RepoBranch,
		RepoCommit: data.RepoCommit,
		Operator: uid.(int),
		Status: sts,
		UpdateTime: time.Now().AddDate(0,0,0),
	}

	e := models.DB.Save(&deploy).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加项目提单申请成功", "", c)
}

// @Tags 应用发布
// @Description 应用发布修改
// @Summary  应用发布修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Param Data body admin.DeployApp true "应用发布信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app/{id} [put]
func PutAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id 	:= c.Param("id")
	var data DeployApp
	var deploy models.AppDeploy

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Deploy Data", "", c)
		return
	}

	fmt.Println(data)
	// 发布唯一性检查
	models.DB.Model(&models.AppDeploy{}).
		Where("tid = ?", data.Tid).
		Where("repo_branch = ?", data.RepoBranch).
		Where("repo_commit = ?", data.RepoCommit).
		Where("status <= 2").
		Find(&deploy)

	if deploy.ID > 0 {
		util.JsonRespond(500, "重复的项目上线提单，请检查！", "", c)
		return
	}

	models.DB.Find(&deploy, id)

	uid,_ 				:= c.Get("Uid")
	deploy.Name     	= data.Name
	deploy.RepoBranch	= data.RepoBranch
	deploy.RepoCommit 	= data.RepoCommit
	deploy.Operator		= uid.(int)
	deploy.UpdateTime   = time.Now().AddDate(0,0,0)


	e := models.DB.Save(&deploy).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改Project成功", "", c)
}

// @Tags 应用发布
// @Description 应用发布删除
// @Summary  应用发布删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app/{id} [delete]
func DelAppDeploy(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id := c.Param("id")

	// 删除之前检查项目发布状态，只能删除status为1（新提单）， 2（审核成功）
	var deploy models.AppDeploy
	models.DB.Model(&models.HostApp{}).
		Where("id = ?", id).
		Where("status >= 3").
		Find(&deploy)

	if deploy.ID > 0 {
		util.JsonRespond(406, "改项目项目提单申请已经执行上线操作，无法删除！", "", c)
		return
	}

	err := models.DB.Delete(models.AppDeploy{}, "id = ?", id).Error
	if err != nil {
		util.JsonRespond(500, err.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除项目提单申请成功", "", c)
}


// @Tags 应用发布
// @Description 应用发布审核
// @Summary  应用发布审核
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "应用发布ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/deploy/app/{id}/review/{status} [put]
func PutAppDeployStatus(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"deploy-app-review") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	id 			:= c.Param("id")
	status,_  	:= com.StrTo(c.Param("status")).Int()
	uid,_ 		:= c.Get("Uid")

	rows 		:= models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Update("status", status).
		Update("review", uid.(int)).RowsAffected

	if rows == 0 {
		util.JsonRespond(500, "修改状态失败！", "", c)
		return
	}

	if status == 2 {
		util.JsonRespond(200, "审核通过", "", c)
		return
	}

	util.JsonRespond(200, "审核拒绝", "", c)
}


// 发布逻辑
func PutAppDeployRedo(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware1("deploy-app-redo", c)

	id 			:= c.Param("id")
	userid,_ 	:= c.Get("Uid")
	uid 		:= userid.(int)

	//升级get请求为webSocket协议
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("cant upgrade connection:"))
		return
	}
	defer wsConn.Close()

	// 执行发布逻辑
	var deploy models.AppDeploy
	models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Where("status = ?", ReviewSuccess).
		Find(&deploy)

	// 检查发布申请信息
	if deploy.ID <= 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("错误的发布申请，请检查！"))
		return
	}

	// 检查发布模板信息
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", deploy.Tid).
		Find(&det)

	if det.Dtid <= 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("项目发布模板已经删除，请检查！"))
		return
	}

	// 检查项目信息
	var app models.App
	models.DB.Model(&models.App{}).
		Where("id = ?", det.Aid).
		Where("active = 1").
		Find(&app)

	if app.ID <= 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("项目状态不可用或者项目已经删除，请检查！"))
		return
	}

	if app.DeployType != 0 {
		// 检查发布类型
		switch app.Tid {
		// backend jar
		case 1:
			BackendJarDeploy(id, app, det, deploy, uid, wsConn)
		case 2:
			fmt.Println(2)
		}
	} else {
		// 通用目录 copy 方式
		DirectoryCopy(id, app, det, deploy, uid, wsConn)
	}
}

func BackendJarDeploy(id string, app models.App, det models.DeployExtend, deploy models.AppDeploy, uid int, ws *websocket.Conn)  {
	// 检查项目是否初始化
	var hostapp []models.HostApp
	sts := 0
	if app.EnableSync == 1 {
		sts = 1
	}
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", app.ID).
		Where("status = ?", sts).
		Find(&hostapp)

	if len(hostapp) <= 0 && sts == 1 {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目还没有初始化，请检查！"))
		return
	}

	if len(hostapp) <= 0 && sts == 0 {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目还没有绑定到主机，请检查！"))
		return
	}

	// 锁定项目，一个项目同时只能允许一个执行该方法
	key := models.GitAppOnWorking + app.Name
	if models.GetValByKey(key) != "" {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目别的用户在发布中，请稍后重试！"))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("检查本地代码库是否有指定的分支"))
	e := util.GitCheckoutByCommit(app.ID, det.RepoUrl, deploy.RepoCommit)
	if e != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
		return
	}

	// 本地拉取代码前执行的任务
	if det.PreCode != "" {

	}

	// 本地拉取代码后要执行的任务， 如果编译打包
	gitpath := util.ReturnGitLocalPath(app.ID, det.RepoUrl)
	if det.PostCode != "" {
		e    = util.ExecRuntimeCmdToWs(det.PostCode, gitpath, ws)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}
	}

	// 本次需要发布的主机
	hosts := strings.Split(det.HostIds,",")

	// 执行部署任务
	for _, v := range hostapp {
		if !util.StringInSlice(strconv.Itoa(v.Hid), hosts) {
			continue
		}

		var host models.Host
		models.DB.Model(&models.Host{}).
			Where("id = ?", v.Hid).
			Where("status = 1").
			Find(&host)

		// sftp 通道建立
		clientConfig, _ := util.ReturnClientConfig(host.Username, "")
		hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
		Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		SftpClient, e := util.GetSftpClient(Scli)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		// 远端服务器， 发布前执行的任务
		if det.PreDeploy != "" {

		}

		path 	:= "/data/webapps/" + app.Name
		// 检查是否为首次发布
		_, e = SftpClient.Stat( path + "/lib/" +
			app.Name + "-" +
			det.Tag + ".jar")
		if e == nil {
			// 备份逻辑
			msg := "备份服务器" + host.Name
			ws.WriteMessage(websocket.TextMessage, []byte(msg))

			cmd		:= "cp " + path + "/lib/%s-" +
				det.Tag + ".jar" + " " +
				path + "/temp/%s-" +
				det.Tag + ".jar_%d"

			cmd    	= fmt.Sprintf(cmd, app.Name, app.Name, time.Now().Unix())
			fmt.Println(cmd)

			util.ExecCmdBySshToWs(cmd, Scli, ws)
		}

		// copy 逻辑
		ws.WriteMessage(websocket.TextMessage, []byte("Copy项目包到远端服务器" + host.Name + host.Addres))

		var localFile string
		localFile 	= gitpath + "/target/"  + app.Name + "-" +
			det.Tag + ".jar"
		if !util.Exists(localFile) {
			localFile 	= gitpath + "/" + app.Name +
				"/target/" + app.Name + "-" +
				det.Tag + ".jar"
		}

		e = util.PutFile(SftpClient, localFile,  path + "/lib" )
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		// 修改文件属主
		ws.WriteMessage(websocket.TextMessage, []byte("修改包的属主属性！"))
		cmd := "chown -R tomcat:tomcat " +  path + "/lib"
		util.ExecCmdBySshToWs(cmd, Scli, ws)

		// 创建软连接
		ws.WriteMessage(websocket.TextMessage, []byte("创建软连接！"))
		cmd = "ln -sv " +  path + "/lib/" + app.Name +
			"-" + det.Tag + ".jar" + " "  +
			path + "/lib/" + app.Name + ".jar"
		util.ExecCmdBySshToWs(cmd, Scli, ws)

		// 重启服务
		ws.WriteMessage(websocket.TextMessage, []byte("重启服务！"))
		cmd = "cd " + path + "/bin; bash control restart"
		util.ExecCmdBySshToWs(cmd, Scli, ws)
	}

	// 取消锁定
	models.SetValByKey(key, "1",  2 * time.Second)


	// 修改数据库信息 stauts 5
	e 	= models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Update("status", DeploySuccess).
		Update("deploy", uid).
		Error

	if e != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("项目部署完毕!"))
}

func DirectoryCopy(id string, app models.App, det models.DeployExtend, deploy models.AppDeploy, uid int, ws *websocket.Conn) {
	// 检查项目是否初始化
	var hostapp []models.HostApp
	sts := 0
	if app.EnableSync == 1 {
		sts = 1
	}
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", app.ID).
		Where("status = ?", sts).
		Find(&hostapp)

	if len(hostapp) <= 0 && sts == 1 {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目还没有初始化，请检查！"))
		return
	}

	if len(hostapp) <= 0 && sts == 0 {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目还没有绑定到主机，请检查！"))
		return
	}

	// 锁定项目，一个项目同时只能允许一个执行该方法
	key := models.GitAppOnWorking + app.Name
	if models.GetValByKey(key) != "" {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目别的用户在发布中，请稍后重试！"))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("检查本地代码库是否有指定的分支"))
	e := util.GitCheckoutByCommit(app.ID, det.RepoUrl, deploy.RepoCommit)
	if e != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
		return
	}

	// 本地拉取代码前执行的任务
	if det.PreCode != "" {

	}

	// 本地拉取代码后要执行的任务， 如果编译打包
	gitpath := util.ReturnGitLocalPath(app.ID, det.RepoUrl)
	if det.PostCode != "" {
		e    = util.ExecRuntimeCmdToWs(det.PostCode, gitpath, ws)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}
	}

	// 本次需要发布的主机
	hosts := strings.Split(det.HostIds,",")

	// 执行部署任务
	for _, v := range hostapp {
		if !util.StringInSlice(strconv.Itoa(v.Hid), hosts) {
			continue
		}

		var host models.Host
		models.DB.Model(&models.Host{}).
			Where("id = ?", v.Hid).
			Where("status = 1").
			Find(&host)

		// sftp 通道建立
		clientConfig, _ := util.ReturnClientConfig(host.Username, "")
		hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
		Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		SftpClient, e := util.GetSftpClient(Scli)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		// 远端服务器， 发布前执行的任务
		if det.PreDeploy != "" {

		}

		path 	:= det.DstDir
		backdir := det.DstRepo + "/" + strconv.Itoa(app.ID) + "/" + det.TemplateName
		cmd     := "mkdir -p " + backdir

		// 备份逻辑
		msg := "备份服务器:" + host.Name + "  业务:" + app.Name
		ws.WriteMessage(websocket.TextMessage, []byte(msg))
		util.ExecCmdBySshToWs(cmd, Scli, ws)

		cmd		= "cp -r " + path +
				" " + backdir + "/" + strconv.FormatInt(time.Now().Unix(),10)

		fmt.Println(cmd)

		util.ExecCmdBySshToWs(cmd, Scli, ws)

		// copy 逻辑
		ws.WriteMessage(websocket.TextMessage, []byte("Copy项目到远端服务器" + host.Name + host.Addres))
		e = util.PutDirectory(SftpClient, gitpath, path)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		// 执行发布后任务
		ws.WriteMessage(websocket.TextMessage, []byte("执行发布后任务 : " + det.PostDeploy))
		util.ExecCmdBySshToWs(det.PostDeploy, Scli, ws)


		// 历史版本数
		cmd 	= "cd %s && ls -1tr %s | head -n -%d | xargs -d '\\n' rm -rf --"
		cmd    	= fmt.Sprintf(cmd, backdir, backdir, det.Versions)
		fmt.Println(cmd)
		util.ExecCmdBySshToWs(cmd, Scli, ws)
	}

	// 取消锁定
	models.SetValByKey(key, "1",  2 * time.Second)


	// 修改数据库信息 stauts 5
	e 	= models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Update("status", DeploySuccess).
		Update("deploy", uid).
		Error

	if e != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("项目部署完毕!"))
}

// 回滚管理
func PutAppDeployUndo(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware1("deploy-app-undo", c)

	id 			:= c.Param("id")
	userid,_ 	:= c.Get("Uid")
	uid 		:= userid.(int)

	//升级get请求为webSocket协议
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("can not upgrade connection:"))
		return
	}
	defer wsConn.Close()

	// 执行回滚逻辑
	var deploy models.AppDeploy

	models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Where("status = ?", DeploySuccess).
		Find(&deploy)

	// 检查发布申请信息
	if deploy.ID <= 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("错误的回滚申请，请检查！"))
		return
	}

	// 检查发布模板信息
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", deploy.Tid).
		Find(&det)

	if det.Dtid <= 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("项目发布模板已经删除，请检查！"))
		return
	}

	// 检查项目信息
	var app models.App
	models.DB.Model(&models.App{}).
		Where("id = ?", det.Aid).
		Where("active = 1").
		Find(&app)

	if app.ID <= 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("项目状态不可用或者项目已经删除，请检查！"))
		return
	}

	if app.DeployType != 0 {
		// 检查发布类型
		switch app.Tid {
		// backend jar
		case 1:
			BackendJarUndoDeploy(id, app, det, uid, wsConn)
		case 2:
			fmt.Println(2)
		}
	} else {
		// 通用目录 copy 方式
		DirectoryCopyUndo(id, app, det, uid, wsConn)
	}
}

func DirectoryCopyUndo(id string, app models.App, det models.DeployExtend, uid int, ws *websocket.Conn) {
	// 检查项目是否初始化
	var hostapp []models.HostApp
	sts := 0
	if app.EnableSync == 1 {
		sts = 1
	}

	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", app.ID).
		Where("status = ?", sts).
		Find(&hostapp)

	// 本次需要回滚的主机
	hosts := strings.Split(det.HostIds,",")

	// 执行部署任务
	for _, v := range hostapp {
		if !util.StringInSlice(strconv.Itoa(v.Hid), hosts) {
			continue
		}

		var host models.Host
		models.DB.Model(&models.Host{}).
			Where("id = ?", v.Hid).
			Where("status = 1").
			Find(&host)

		// sftp 通道建立
		clientConfig, _ := util.ReturnClientConfig(host.Username, "")
		hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
		Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		path 	:= det.DstDir
		backdir := det.DstRepo + "/" + strconv.Itoa(app.ID) + "/" + det.TemplateName

		//获取最近一次备份
		cmd 	:= "cd %s && ls -1t | head -n1"
		cmd    	= fmt.Sprintf(cmd, backdir)
		res, e 	:= util.ExecuteCmd(cmd, Scli)
		res 	= strings.Replace(res, "\n", "", -1)

		if e != nil {
			// 回滚失败修改数据库信息
			models.DB.Model(&models.AppDeploy{}).
				Where("id = ?", id).
				Update("status", UndoFail).
				Update("deploy", uid)
		}

		cmd		= "cp -rf " + backdir + "/" + res + "/*" + " " + path
		util.ExecCmdBySshToWs(cmd, Scli, ws)

		// 执行发布后任务
		ws.WriteMessage(websocket.TextMessage, []byte("执行发布后任务 : " + det.PostDeploy))
		util.ExecCmdBySshToWs(det.PostDeploy, Scli, ws)
	}

	// 修改数据库信息 stauts 5
	e 	:= models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Update("status", UndoSuccess).
		Update("deploy", uid).
		Error

	if e != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("回滚成功修改数据库信息失败！"))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("项目回滚完毕!"))
}

func BackendJarUndoDeploy(id string, app models.App, det models.DeployExtend, uid int, ws *websocket.Conn)  {
	var hostapp []models.HostApp
	// 执行回滚任务
	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", app.ID).
		Where("status = 1").
		Find(&hostapp)

	if len(hostapp) <= 0 {
		ws.WriteMessage(websocket.TextMessage, []byte("该项目还没有初始化，请检查！"))
		return
	}

	// 本次需要回滚的主机
	hosts := strings.Split(det.HostIds,",")

	for _, v := range hostapp {
		if !util.StringInSlice(strconv.Itoa(v.Hid), hosts) {
			continue
		}

		var host models.Host
		models.DB.Model(&models.Host{}).
			Where("id = ?", v.Hid).
			Where("status = 1").
			Find(&host)

		// sftp 通道建立
		clientConfig, _ := util.ReturnClientConfig(host.Username, "")
		hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
		Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
		if e != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(e.Error()))
			return
		}

		path 	:= det.DstDir
		backdir := det.DstRepo + "/" + strconv.Itoa(app.ID) + "/" + det.TemplateName

		//获取最近一次备份
		cmd 	:= "cd %s && ls -1t | head -n1"
		cmd    	= fmt.Sprintf(cmd, backdir)
		res, e 	:= util.ExecuteCmd(cmd, Scli)
		res 	= strings.Replace(res, "\n", "", -1)

		ws.WriteMessage(websocket.TextMessage, []byte("回滚到上一个版本:" + res))
		if e != nil {
			// 回滚失败修改数据库信息
			models.DB.Model(&models.AppDeploy{}).
				Where("id = ?", id).
				Update("status", UndoFail).
				Update("deploy", uid)
		}

		cmd		= "cp -rf " + backdir + "/" + res + "/*" + " " + path

		fmt.Println(cmd)
		util.ExecCmdBySshToWs(cmd, Scli, ws)


		// 重启服务
		ws.WriteMessage(websocket.TextMessage, []byte("重启服务！"))
		cmd = "cd " + path + "/bin; bash control restart"
		util.ExecCmdBySshToWs(cmd, Scli, ws)
	}

	// 回滚成功修改数据库信息
	e := models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Update("status", UndoSuccess).
		Update("deploy", uid).Error

	if e != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("回滚成功修改数据库信息失败！"))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("回滚成功！"))
}
